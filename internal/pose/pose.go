package pose

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/model"
)

// Epsilon is the threshold for determining if a part should be visible
const Epsilon = 0.001

// PartData represents a single part entry in a pose group
type PartData struct {
	PartId         string
	ParameterIndex int // index into the parameter array (-1 if not found)
	PartIndex      int // index into the part array (-1 if not found)
	Link           []PartData
}

// PoseManager manages pose groups and applies part opacity transitions
type PoseManager struct {
	groups         [][]PartData
	fadeInTime     float64
	lastDelta      float64
	partGroupCounts []int
}

// NewPoseManager creates a new PoseManager from pose3.json data
func NewPoseManager(poseJson model.PoseJson, core core.Core, modelPtr uintptr) *PoseManager {
	pm := &PoseManager{
		fadeInTime: poseJson.FadeInTime,
	}

	// Get parameter IDs and part IDs for index lookup
	allParamIds := getParameterIds(core, modelPtr)
	partIds := core.GetPartIds(modelPtr)

	// Build pose groups
	pm.groups = make([][]PartData, len(poseJson.Groups))
	pm.partGroupCounts = make([]int, len(poseJson.Groups))

	for i, group := range poseJson.Groups {
		pm.groups[i] = make([]PartData, len(group))
		pm.partGroupCounts[i] = len(group)

		for j, entry := range group {
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
			pm.groups[i][j] = pd
		}
	}

	// Initialize: set first part in each group to visible
	for _, group := range pm.groups {
		if len(group) > 0 {
			// Set default: first part visible
			for j := range group {
				if j == 0 {
					core.SetPartOpacity(modelPtr, group[j].PartId, 1.0)
					for _, link := range group[j].Link {
						core.SetPartOpacity(modelPtr, link.PartId, 1.0)
					}
				}
			}
		}
	}

	return pm
}

// Update applies pose logic — should be called after motion update but before model.Update()
// However, since model.Update() is called in the main update loop, we apply pose
// after core.Update() to override part opacities based on the current state.
func (pm *PoseManager) Update(core core.Core, modelPtr uintptr, deltaTime float64) {
	pm.lastDelta = deltaTime

	// Process each pose group
	for _, group := range pm.groups {
		// Find which part in this group should be visible based on parameter values
		visibleIndex := -1
		maxOpacity := float32(-1.0)

		for j, part := range group {
			if part.ParameterIndex >= 0 {
				// Check if this part's parameter value indicates visibility
				paramValue := core.GetParameterValue(modelPtr, part.PartId)
				if paramValue > Epsilon {
					visibleIndex = j
				}
			}
			// Also check current part opacity
			currentOpacity := core.GetPartOpacities(modelPtr)
			if part.PartIndex >= 0 && part.PartIndex < len(currentOpacity) {
				if currentOpacity[part.PartIndex] > maxOpacity {
					maxOpacity = currentOpacity[part.PartIndex]
				}
			}
		}

		// If no parameter explicitly set visibility, use the first part as default
		if visibleIndex == -1 {
			visibleIndex = 0
		}

		// Apply fade to the group
		pm.doFade(core, modelPtr, group, visibleIndex, deltaTime)
	}
}

// doFade applies fade-in/fade-out transitions to parts in a group
func (pm *PoseManager) doFade(core core.Core, modelPtr uintptr, group []PartData, visibleIndex int, deltaTime float64) {
	fadeInTime := pm.fadeInTime
	if fadeInTime <= 0 {
		fadeInTime = 0.5 // default fade time
	}

	for j, part := range group {
		if part.PartIndex < 0 {
			continue
		}

		currentOpacities := core.GetPartOpacities(modelPtr)
		if part.PartIndex >= len(currentOpacities) {
			continue
		}
		currentOpacity := currentOpacities[part.PartIndex]

		var newOpacity float32
		if j == visibleIndex {
			// Fade in
			if currentOpacity < 1.0 {
				newOpacity = currentOpacity + float32(deltaTime/fadeInTime)
				if newOpacity > 1.0 {
					newOpacity = 1.0
				}
			} else {
				newOpacity = 1.0
			}
		} else {
			// Fade out
			if currentOpacity > 0 {
				newOpacity = currentOpacity - float32(deltaTime/fadeInTime)
				if newOpacity < 0 {
					newOpacity = 0
				}
			} else {
				newOpacity = 0
			}
		}

		core.SetPartOpacity(modelPtr, part.PartId, newOpacity)

		// Copy opacity to linked parts
		for _, link := range part.Link {
			if link.PartIndex >= 0 {
				core.SetPartOpacity(modelPtr, link.PartId, newOpacity)
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
func getParameterIds(core core.Core, modelPtr uintptr) []string {
	params := core.GetParameters(modelPtr)
	ids := make([]string, len(params))
	for i, p := range params {
		ids[i] = p.Id
	}
	return ids
}
