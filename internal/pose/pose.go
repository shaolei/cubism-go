package pose

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/model"
)

const (
	Epsilon                = 0.001
	DefaultFadeInSeconds   = 0.5
	Phi                    = 0.5   // Cross-fade intersection point
	BackOpacityThreshold   = 0.15  // Maximum background visibility ratio
)

// PartData represents a single part entry in a pose group.
// Matches the official SDK's CubismPose::PartData structure.
type PartData struct {
	PartId         string
	ParameterIndex int // index into the parameter array (-1 if not found)
	PartIndex      int // index into the part array (-1 if not found)
	Link           []PartData
}

// PoseManager manages pose groups and applies part opacity transitions.
// Matches the official SDK's CubismPose design with:
//   - Flat partGroups array with partGroupCounts for grouping
//   - DoFade with Phi/BackOpacityThreshold cross-fade algorithm
//   - CopyPartOpacities for Link handling (separate from DoFade)
//   - Reset for initial state setup
type PoseManager struct {
	partGroups      []PartData // flat array of all parts across all groups
	partGroupCounts []int      // count of parts in each group
	fadeTimeSeconds float64
	lastModelPtr    uintptr // track model pointer for Reset detection
	idManager       *id.CubismIdManager
}

// NewPoseManager creates a new PoseManager from pose3.json data.
// Matches the official SDK's CubismPose::Create parsing logic.
func NewPoseManager(poseJson model.PoseJson, c core.Core, modelPtr uintptr) *PoseManager {
	pm := &PoseManager{
		fadeTimeSeconds: poseJson.FadeInTime,
	}

	// Validate fade time
	if pm.fadeTimeSeconds < 0.0 {
		pm.fadeTimeSeconds = DefaultFadeInSeconds
	}

	// Get parameter IDs and part IDs for index lookup
	allParamIds := getParameterIds(c, modelPtr)
	partIds := c.GetPartIds(modelPtr)

	// Build flat partGroups array and partGroupCounts
	// Matches official SDK's flat structure with group counts
	for _, group := range poseJson.Groups {
		groupCount := 0
		for _, entry := range group {
			pd := PartData{
				PartId:         entry.Id,
				ParameterIndex: findIndex(allParamIds, entry.Id),
				PartIndex:      findIndex(partIds, entry.Id),
			}
			// Build link data
			for _, linkId := range entry.Link {
				linkPd := PartData{
					PartId:         linkId,
					ParameterIndex: findIndex(allParamIds, linkId),
					PartIndex:      findIndex(partIds, linkId),
				}
				pd.Link = append(pd.Link, linkPd)
			}
			pm.partGroups = append(pm.partGroups, pd)
			groupCount++
		}
		pm.partGroupCounts = append(pm.partGroupCounts, groupCount)
	}

	// Initialize parts
	pm.Reset(c, modelPtr)

	return pm
}

// SetIdManager sets the CubismIdManager for fast parameter/part access by index
func (pm *PoseManager) SetIdManager(idMgr *id.CubismIdManager) {
	pm.idManager = idMgr
}

// Reset initializes the pose state, setting the first part in each group to visible.
// Matches the official SDK's CubismPose::Reset.
func (pm *PoseManager) Reset(c core.Core, modelPtr uintptr) {
	beginIndex := 0

	for i := 0; i < len(pm.partGroupCounts); i++ {
		groupCount := pm.partGroupCounts[i]

		for j := beginIndex; j < beginIndex+groupCount; j++ {
			pd := &pm.partGroups[j]

			// Initialize parameter value to 1 for this part
			if pd.ParameterIndex >= 0 {
				if pm.idManager != nil {
					c.SetParameterValueByIndex(modelPtr, pd.ParameterIndex, 1.0)
				} else {
					c.SetParameterValue(modelPtr, pd.PartId, 1.0)
				}
			}

			// Set part opacity: first in group = 1.0, others = 0.0
			if pd.PartIndex >= 0 {
				opacity := float32(0.0)
				if j == beginIndex {
					opacity = 1.0
				}
				if pm.idManager != nil {
					c.SetPartOpacityByIndex(modelPtr, pd.PartIndex, opacity)
				} else {
					c.SetPartOpacity(modelPtr, pd.PartId, opacity)
				}
			}

			// Initialize links
			for k := range pd.Link {
				link := &pd.Link[k]
				if link.ParameterIndex >= 0 {
					if pm.idManager != nil {
						c.SetParameterValueByIndex(modelPtr, link.ParameterIndex, 1.0)
					} else {
						c.SetParameterValue(modelPtr, link.PartId, 1.0)
					}
				}
			}
		}

		beginIndex += groupCount
	}
}

// Update applies pose logic.
// Matches the official SDK's CubismPose::UpdateParameters:
// 1. If model changed, Reset
// 2. For each group: DoFade
// 3. CopyPartOpacities for Link handling
func (pm *PoseManager) Update(c core.Core, modelPtr uintptr, deltaTimeSeconds float64) {
	// If the model has changed, reset
	if modelPtr != pm.lastModelPtr {
		pm.Reset(c, modelPtr)
	}
	pm.lastModelPtr = modelPtr

	// Clamp negative delta time
	if deltaTimeSeconds < 0.0 {
		deltaTimeSeconds = 0.0
	}

	beginIndex := 0
	for i := 0; i < len(pm.partGroupCounts); i++ {
		partGroupCount := pm.partGroupCounts[i]
		pm.doFade(c, modelPtr, deltaTimeSeconds, beginIndex, partGroupCount)
		beginIndex += partGroupCount
	}

	pm.copyPartOpacities(c, modelPtr)
}

// doFade applies the cross-fade algorithm for a single pose group.
// Matches the official SDK's CubismPose::DoFade exactly:
//   - Find the first part whose parameter > Epsilon as visible
//   - Fade in the visible part linearly
//   - For hidden parts, calculate opacity using Phi intersection and BackOpacityThreshold
func (pm *PoseManager) doFade(c core.Core, modelPtr uintptr, deltaTimeSeconds float64, beginIndex, partGroupCount int) {
	visiblePartIndex := -1
	newOpacity := float32(1.0)

	// Find the currently visible part
	for i := beginIndex; i < beginIndex+partGroupCount; i++ {
		pd := &pm.partGroups[i]
		partIndex := pd.PartIndex
		paramIndex := pd.ParameterIndex

		// Get parameter value by index
		var paramValue float32
		if paramIndex >= 0 {
			if pm.idManager != nil {
				paramValue = c.GetParameterValueByIndex(modelPtr, paramIndex)
			} else {
				paramValue = c.GetParameterValue(modelPtr, pd.PartId)
			}
		}

		if paramValue > Epsilon {
			if visiblePartIndex >= 0 {
				break // Only the first visible part is found
			}

			visiblePartIndex = i

			if pm.fadeTimeSeconds == 0.0 {
				newOpacity = 1.0
				continue
			}

			// Get current opacity and calculate fade-in
			var currentOpacity float32
			if pm.idManager != nil {
				currentOpacity = c.GetPartOpacityByIndex(modelPtr, partIndex)
			} else {
				opacities := c.GetPartOpacities(modelPtr)
				if partIndex >= 0 && partIndex < len(opacities) {
					currentOpacity = opacities[partIndex]
				}
			}

			newOpacity = currentOpacity + float32(deltaTimeSeconds/pm.fadeTimeSeconds)
			if newOpacity > 1.0 {
				newOpacity = 1.0
			}
		}
	}

	// Default to first part if none visible
	if visiblePartIndex < 0 {
		visiblePartIndex = beginIndex
		newOpacity = 1.0
	}

	// Set opacity for visible and hidden parts
	for i := beginIndex; i < beginIndex+partGroupCount; i++ {
		pd := &pm.partGroups[i]
		partsIndex := pd.PartIndex

		if partsIndex < 0 {
			continue
		}

		if visiblePartIndex == i {
			// Visible part: set new opacity directly
			if pm.idManager != nil {
				c.SetPartOpacityByIndex(modelPtr, partsIndex, newOpacity)
			} else {
				c.SetPartOpacity(modelPtr, pd.PartId, newOpacity)
			}
		} else {
			// Hidden part: calculate opacity using cross-fade algorithm
			var opacity float32
			if pm.idManager != nil {
				opacity = c.GetPartOpacityByIndex(modelPtr, partsIndex)
			} else {
				opacities := c.GetPartOpacities(modelPtr)
				if partsIndex < len(opacities) {
					opacity = opacities[partsIndex]
				}
			}

			var a1 float32 // calculated opacity boundary
			if newOpacity < Phi {
				// Line through (0,1) and (Phi, Phi)
				a1 = newOpacity*(Phi-1.0)/Phi + 1.0
			} else {
				// Line through (1,0) and (Phi, Phi)
				a1 = (1.0 - newOpacity) * Phi / (1.0 - Phi)
			}

			// Limit background visibility ratio
			backOpacity := (1.0 - a1) * (1.0 - newOpacity)
			if backOpacity > BackOpacityThreshold {
				a1 = 1.0 - BackOpacityThreshold/(1.0-newOpacity)
			}

			// Only reduce opacity, never increase for hidden parts
			if opacity > a1 {
				opacity = a1
			}

			if pm.idManager != nil {
				c.SetPartOpacityByIndex(modelPtr, partsIndex, opacity)
			} else {
				c.SetPartOpacity(modelPtr, pd.PartId, opacity)
			}
		}
	}
}

// copyPartOpacities copies opacity from parent parts to their linked parts.
// Matches the official SDK's CubismPose::CopyPartOpacities.
// This is called after DoFade for all groups, separate from the fade logic.
func (pm *PoseManager) copyPartOpacities(c core.Core, modelPtr uintptr) {
	for i := range pm.partGroups {
		pd := &pm.partGroups[i]

		if len(pd.Link) == 0 {
			continue
		}

		// Get this part's opacity
		var opacity float32
		if pm.idManager != nil {
			opacity = c.GetPartOpacityByIndex(modelPtr, pd.PartIndex)
		} else {
			opacities := c.GetPartOpacities(modelPtr)
			if pd.PartIndex >= 0 && pd.PartIndex < len(opacities) {
				opacity = opacities[pd.PartIndex]
			}
		}

		// Copy to linked parts
		for j := range pd.Link {
			linkPart := &pd.Link[j]
			if linkPart.PartIndex < 0 {
				continue
			}
			if pm.idManager != nil {
				c.SetPartOpacityByIndex(modelPtr, linkPart.PartIndex, opacity)
			} else {
				c.SetPartOpacity(modelPtr, linkPart.PartId, opacity)
			}
		}
	}
}

// findIndex returns the index of id in the slice, or -1 if not found
func findIndex(ids []string, id string) int {
	for i, s := range ids {
		if s == id {
			return i
		}
	}
	return -1
}

// getParameterIds gets all parameter IDs from the core
func getParameterIds(c core.Core, modelPtr uintptr) []string {
	return c.GetParameterIds(modelPtr)
}
