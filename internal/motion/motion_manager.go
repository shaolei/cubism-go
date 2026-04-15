package motion

import (
	"github.com/shaolei/cubism-go/internal/core"
)

// Priority levels for motion playback
const (
	PriorityNone       = 0
	PriorityIdle       = 1
	PriorityNormal     = 2
	PriorityForce      = 3
)

// CubismMotionManager extends CubismMotionQueueManager with a priority system,
// matching the official SDK's CubismMotionManager design.
// Only motions with higher or equal priority than the current reservation
// can be started.
type CubismMotionManager struct {
	*CubismMotionQueueManager
	currentPriority int
	reservePriority int
}

// NewCubismMotionManager creates a new priority-based motion manager
func NewCubismMotionManager(c core.Core, modelPtr uintptr, onFinished func(int)) *CubismMotionManager {
	return &CubismMotionManager{
		CubismMotionQueueManager: NewCubismMotionQueueManager(c, modelPtr, onFinished),
		currentPriority:          PriorityNone,
		reservePriority:          PriorityNone,
	}
}

// GetCurrentPriority returns the priority of the currently playing motion
func (mm *CubismMotionManager) GetCurrentPriority() int {
	return mm.currentPriority
}

// GetReservePriority returns the reserved priority for the next motion
func (mm *CubismMotionManager) GetReservePriority() int {
	return mm.reservePriority
}

// SetReservePriority sets the reserved priority for the next motion
func (mm *CubismMotionManager) SetReservePriority(priority int) {
	mm.reservePriority = priority
}

// CanStartMotion checks whether a motion with the given priority can be started
func (mm *CubismMotionManager) CanStartMotion(priority int) bool {
	if priority == PriorityForce {
		return true
	}
	// Can start if reserve priority allows it
	if mm.reservePriority == PriorityNone {
		return true
	}
	return priority >= mm.reservePriority
}

// StartMotion starts a motion with the given priority.
// Returns -1 if the motion cannot be started due to priority constraints.
func (mm *CubismMotionManager) StartMotionWithPriority(mtn Motion, loop bool, priority int) int {
	if !mm.CanStartMotion(priority) {
		return -1
	}

	mm.currentPriority = priority
	mm.reservePriority = PriorityNone

	id := mm.CubismMotionQueueManager.StartMotion(mtn, loop, false)
	return id
}

// Update is the main update loop for the motion manager
func (mm *CubismMotionManager) Update(deltaTime float64) {
	mm.DoUpdateMotion(deltaTime)

	// Reset priority when all motions have finished
	if mm.IsFinished() {
		mm.currentPriority = PriorityNone
	}
}
