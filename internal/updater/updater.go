package updater

// UpdateOrder defines execution priority constants for model subsystem updaters.
// Lower values execute first. Matches the official SDK's CubismUpdateOrder enum.
const (
	UpdateOrderEyeBlink   = 200
	UpdateOrderExpression = 300
	UpdateOrderLook       = 400
	UpdateOrderBreath     = 500
	UpdateOrderPhysics    = 600
	UpdateOrderLipSync    = 700
	UpdateOrderPose       = 800
	UpdateOrderMax        = 1<<31 - 1 // math.MaxInt32
)

// Updater is the interface for model subsystem updates.
// Each subsystem (blink, expression, breath, physics, look, pose) implements
// this interface and is registered with the UpdateScheduler at a specific
// execution order priority.
type Updater interface {
	// OnLateUpdate is called each frame to update the subsystem.
	// deltaTimeSeconds is the elapsed time since the last frame.
	OnLateUpdate(deltaTimeSeconds float64)
}

// FuncUpdater wraps a function as an Updater, similar to http.HandlerFunc.
type FuncUpdater struct {
	ExecutionOrder int
	Fn            func(deltaTimeSeconds float64)
}

// NewFuncUpdater creates a FuncUpdater with the given execution order and function.
func NewFuncUpdater(executionOrder int, fn func(float64)) *FuncUpdater {
	return &FuncUpdater{
		ExecutionOrder: executionOrder,
		Fn:            fn,
	}
}

// OnLateUpdate calls the wrapped function.
func (f *FuncUpdater) OnLateUpdate(deltaTimeSeconds float64) {
	if f.Fn != nil {
		f.Fn(deltaTimeSeconds)
	}
}

// GetExecutionOrder returns the execution order for sorting.
func (f *FuncUpdater) GetExecutionOrder() int {
	return f.ExecutionOrder
}

// orderedUpdater pairs an Updater with its execution order for sorting.
type orderedUpdater struct {
	order   int
	updater Updater
}

// UpdateScheduler manages a list of Updaters sorted by execution order.
// When OnLateUpdate is called, each Updater is invoked in ascending order.
// Matches the official SDK's CubismUpdateScheduler design.
type UpdateScheduler struct {
	updaters  []orderedUpdater
	needsSort bool
}

// NewUpdateScheduler creates a new UpdateScheduler.
func NewUpdateScheduler() *UpdateScheduler {
	return &UpdateScheduler{
		updaters:  make([]orderedUpdater, 0),
		needsSort: false,
	}
}

// AddUpdater adds an Updater with the given execution order.
// The list is marked for sorting before the next update.
func (s *UpdateScheduler) AddUpdater(executionOrder int, u Updater) {
	if u == nil {
		return
	}
	s.updaters = append(s.updaters, orderedUpdater{
		order:   executionOrder,
		updater: u,
	})
	s.needsSort = true
}

// RemoveUpdater removes an Updater from the scheduler.
func (s *UpdateScheduler) RemoveUpdater(u Updater) {
	for i, ou := range s.updaters {
		if ou.updater == u {
			s.updaters = append(s.updaters[:i], s.updaters[i+1:]...)
			return
		}
	}
}

// Clear removes all updaters.
func (s *UpdateScheduler) Clear() {
	s.updaters = s.updaters[:0]
	s.needsSort = false
}

// OnLateUpdate sorts the list if needed, then calls OnLateUpdate on each
// Updater in ascending execution order.
func (s *UpdateScheduler) OnLateUpdate(deltaTimeSeconds float64) {
	if s.needsSort {
		s.sortUpdaters()
	}
	for _, ou := range s.updaters {
		ou.updater.OnLateUpdate(deltaTimeSeconds)
	}
}

// sortUpdaters sorts the updater list by execution order using merge sort,
// matching the official SDK's use of csmVectorSort::MergeSort.
func (s *UpdateScheduler) sortUpdaters() {
	if len(s.updaters) <= 1 {
		s.needsSort = false
		return
	}
	// Merge sort (stable)
	s.mergeSort(0, len(s.updaters))
	s.needsSort = false
}

func (s *UpdateScheduler) mergeSort(begin, end int) {
	if begin+1 >= end {
		return
	}
	mid := (begin + end) / 2
	s.mergeSort(begin, mid)
	s.mergeSort(mid, end)
	s.merge(begin, mid, end)
}

func (s *UpdateScheduler) merge(begin, mid, end int) {
	left := make([]orderedUpdater, mid-begin)
	right := make([]orderedUpdater, end-mid)
	copy(left, s.updaters[begin:mid])
	copy(right, s.updaters[mid:end])

	i, j, k := 0, 0, begin
	for i < len(left) && j < len(right) {
		if left[i].order <= right[j].order {
			s.updaters[k] = left[i]
			i++
		} else {
			s.updaters[k] = right[j]
			j++
		}
		k++
	}
	for i < len(left) {
		s.updaters[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		s.updaters[k] = right[j]
		j++
		k++
	}
}

// Len returns the number of registered updaters.
func (s *UpdateScheduler) Len() int {
	return len(s.updaters)
}
