package motion

import (
	"github.com/shaolei/cubism-go/internal/core"
	"github.com/shaolei/cubism-go/internal/id"
)

// CubismMotionQueueManager manages a queue of motion entries,
// matching the official SDK's CubismMotionQueueManager design.
// It handles the motion queue lifecycle, fade-out triggering,
// and parameter application.
type CubismMotionQueueManager struct {
	core            core.Core
	modelPtr        uintptr
	idManager       *id.CubismIdManager
	entries         []*CubismMotionQueueEntry
	lastId          int
	userTime        float64 // Accumulated global time
	onFinished      func(int)
	savedParameters map[string]float32
}

// NewCubismMotionQueueManager creates a new queue-based motion manager
func NewCubismMotionQueueManager(c core.Core, modelPtr uintptr, onFinished func(int)) *CubismMotionQueueManager {
	return &CubismMotionQueueManager{
		core:       c,
		modelPtr:   modelPtr,
		entries:    make([]*CubismMotionQueueEntry, 0),
		onFinished: onFinished,
	}
}

// SetIdManager sets the CubismIdManager for fast parameter access by index
func (m *CubismMotionQueueManager) SetIdManager(idMgr *id.CubismIdManager) {
	m.idManager = idMgr
}

// StartMotion adds a new motion to the queue.
// If isExpress is false (normal motion), existing entries will have their
// fade-out triggered. Returns the motion queue entry ID.
func (m *CubismMotionQueueManager) StartMotion(mtn Motion, loop bool, isExpress bool) int {
	m.lastId++
	entry := newCubismMotionQueueEntry(mtn, m.lastId, loop)
	m.entries = append(m.entries, entry)

	// For non-expression motions, trigger fade-out on all existing entries
	if !isExpress {
		for _, existing := range m.entries {
			if existing != entry && existing.IsAvailable() && !existing.IsFinished() {
				fadeOutSeconds := existing.GetMotion().FadeOutTime
				if fadeOutSeconds < 0.0 {
					fadeOutSeconds = 0.0
				}
				existing.StartFadeout(fadeOutSeconds, m.userTime)
			}
		}
	}

	return m.lastId
}

// DoUpdateMotion updates all active queue entries and applies their parameters.
// This matches the official SDK's DoUpdateMotion pattern:
// 1. Load saved parameters
// 2. For each entry: update fade weight, evaluate curves, apply parameters
// 3. Save parameters
// 4. Remove finished entries
func (m *CubismMotionQueueManager) DoUpdateMotion(deltaTime float64) {
	m.userTime += deltaTime

	// Remove unavailable entries first
	m.removeUnavailable()

	if len(m.entries) == 0 {
		return
	}

	// Load saved parameters (restore to pre-motion state)
	m.loadParameters()

	// Process each entry
	for _, entry := range m.entries {
		if !entry.IsAvailable() {
			continue
		}

		mtn := entry.GetMotion()

		// Start entry on first update
		if !entry.IsStarted() {
			entry.Start(m.userTime)
			// Play sound at start
			if mtn.Sound != "" {
				mtn.LoadedSound.Play()
			}
		}

		// Update fade weight
		entry.UpdateFadeWeight(m.userTime)

		// Calculate local time within the motion
		localTime := entry.GetLocalTime(m.userTime)
		duration := mtn.Meta.Duration

		// Check if motion has finished
		if !entry.IsLoop() && duration > 0.0 && localTime >= duration {
			entry.finished = true
			if m.onFinished != nil {
				m.onFinished(entry.GetId())
			}
			continue
		}

		// Handle looping: wrap local time
		if entry.IsLoop() && duration > 0.0 {
			for localTime >= duration {
				localTime -= duration
				// Restart entry timing for correct fade calculation on loop
				entry.Restart(m.userTime - localTime)
			}
		}

		// Clamp local time to valid range
		if localTime < 0.0 {
			localTime = 0.0
		}

		// Apply motion curves
		fadeWeight := entry.GetFadeWeight()
		m.applyMotionParameters(mtn, localTime, fadeWeight)
	}

	// Save parameters for next frame
	m.saveParameters()

	// Clean up finished/unavailable entries
	m.removeUnavailable()
}

// applyMotionParameters evaluates all curves in a motion and applies them to the model
func (m *CubismMotionQueueManager) applyMotionParameters(mtn *Motion, localTime float64, fadeWeight float64) {
	for _, curve := range mtn.Curves {
		for _, seg := range curve.Segments {
			if !segmentIntersects(seg, localTime) {
				continue
			}
			value := segmentInterpolate(seg, localTime)

			switch curve.Target {
			case "Parameter":
				m.applyParameterCurve(curve, mtn, value, localTime, fadeWeight)
			case "PartOpacity":
				m.core.SetPartOpacity(m.modelPtr, curve.Id, float32(value))
			case "Model":
				// TODO: implement Model target curves
			}
		}
	}
}

// applyParameterCurve applies a parameter curve value with proper fade weighting,
// matching the official SDK's per-curve fade behavior
func (m *CubismMotionQueueManager) applyParameterCurve(curve Curve, mtn *Motion, value float64, localTime float64, fadeWeight float64) {
	var sourceValue float32
	if m.idManager != nil {
		handle := m.idManager.GetParameterId(curve.Id)
		if !handle.IsValid() {
			return
		}
		sourceValue = m.core.GetParameterValueByIndex(m.modelPtr, int(handle))
	} else {
		sourceValue = m.core.GetParameterValue(m.modelPtr, curve.Id)
	}

	// If curve has its own fade settings, use those; otherwise use motion-level fade
	if curve.FadeInTime < 0.0 && curve.FadeOutTime < 0.0 {
		// No per-curve fade: use motion fade weight
		v := sourceValue + (float32(value)-sourceValue)*float32(fadeWeight)
		if m.idManager != nil {
			handle := m.idManager.GetParameterId(curve.Id)
			m.core.SetParameterValueByIndex(m.modelPtr, int(handle), v)
		} else {
			m.core.SetParameterValue(m.modelPtr, curve.Id, v)
		}
		return
	}

	// Per-curve fade calculation
	var fadeInWeight float64 = 1.0
	var fadeOutWeight float64 = 1.0

	if curve.FadeInTime >= 0.0 {
		if curve.FadeInTime == 0.0 {
			fadeInWeight = 1.0
		} else {
			fadeInWeight = getEasingSine(localTime / curve.FadeInTime)
		}
	} else {
		fadeInWeight = 1.0 // Use motion-level fade-in (already in fadeWeight)
	}

	if curve.FadeOutTime >= 0.0 {
		if curve.FadeOutTime == 0.0 {
			fadeOutWeight = 1.0
		} else {
			fadeOutWeight = getEasingSine((mtn.Meta.Duration - localTime) / curve.FadeOutTime)
		}
	} else {
		fadeOutWeight = 1.0 // Use motion-level fade-out (already in fadeWeight)
	}

	paramWeight := fadeInWeight * fadeOutWeight
	v := sourceValue + (float32(value)-sourceValue)*float32(paramWeight)
	if m.idManager != nil {
		handle := m.idManager.GetParameterId(curve.Id)
		m.core.SetParameterValueByIndex(m.modelPtr, int(handle), v)
	} else {
		m.core.SetParameterValue(m.modelPtr, curve.Id, v)
	}
}

// IsFinished returns whether all queue entries have finished
func (m *CubismMotionQueueManager) IsFinished() bool {
	for _, entry := range m.entries {
		if entry.IsAvailable() && !entry.IsFinished() {
			return false
		}
	}
	return true
}

// StopAllMotions stops all motions in the queue
func (m *CubismMotionQueueManager) StopAllMotions() {
	for _, entry := range m.entries {
		if entry.IsAvailable() && !entry.IsFinished() {
			entry.available = false
			if m.onFinished != nil {
				m.onFinished(entry.GetId())
			}
		}
	}
}

// GetEntryById finds an entry by its ID
func (m *CubismMotionQueueManager) GetEntryById(id int) *CubismMotionQueueEntry {
	for _, entry := range m.entries {
		if entry.GetId() == id {
			return entry
		}
	}
	return nil
}

// Close removes an entry from the queue by ID
func (m *CubismMotionQueueManager) Close(id int) {
	for i, entry := range m.entries {
		if entry.GetId() == id {
			entry.available = false
			// Close sound if any
			if entry.GetMotion().Sound != "" {
				entry.GetMotion().LoadedSound.Close()
			}
			// Remove from slice
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return
		}
	}
}

// Reset restarts an entry by ID (for looping)
func (m *CubismMotionQueueManager) Reset(id int) {
	entry := m.GetEntryById(id)
	if entry != nil {
		entry.Restart(m.userTime)
		entry.finished = false
	}
}

// SetOnFinished sets the callback for when motions finish
func (m *CubismMotionQueueManager) SetOnFinished(fn func(int)) {
	m.onFinished = fn
}

// removeUnavailable removes all unavailable or finished entries from the queue
func (m *CubismMotionQueueManager) removeUnavailable() {
	active := make([]*CubismMotionQueueEntry, 0, len(m.entries))
	for _, entry := range m.entries {
		if entry.IsAvailable() && !entry.IsFinished() {
			active = append(active, entry)
		} else {
			// Close sound for removed entries
			if entry.GetMotion().Sound != "" {
				entry.GetMotion().LoadedSound.Close()
			}
		}
	}
	m.entries = active
}

// saveParameters saves the current parameter values
func (m *CubismMotionQueueManager) saveParameters() {
	parameters := m.core.GetParameters(m.modelPtr)
	if m.savedParameters == nil {
		m.savedParameters = make(map[string]float32, len(parameters))
	}
	for _, parameter := range parameters {
		m.savedParameters[parameter.Id] = parameter.Current
	}
}

// loadParameters restores the saved parameter values
func (m *CubismMotionQueueManager) loadParameters() {
	if m.savedParameters == nil {
		return
	}
	for id, value := range m.savedParameters {
		if m.idManager != nil {
			handle := m.idManager.GetParameterId(id)
			if handle.IsValid() {
				m.core.SetParameterValueByIndex(m.modelPtr, int(handle), value)
			}
		} else {
			m.core.SetParameterValue(m.modelPtr, id, value)
		}
	}
}

// GetSavedParameters returns the saved parameter map (for ExpressionManager access)
func (m *CubismMotionQueueManager) GetSavedParameters() map[string]float32 {
	return m.savedParameters
}
