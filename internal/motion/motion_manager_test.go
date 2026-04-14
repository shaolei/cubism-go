package motion

import (
	"sync"
	"testing"

	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
)

// mockCore implements core.Core for testing
type mockCore struct {
	mu        sync.Mutex
	params    map[string]float32
	partOpas  map[string]float32
}

func newMockMotionCore() *mockCore {
	return &mockCore{
		params:   make(map[string]float32),
		partOpas: make(map[string]float32),
	}
}

func (m *mockCore) LoadMoc(_ string) (moc.Moc, error)        { return moc.Moc{}, nil }
func (m *mockCore) GetVersion() string                         { return "" }
func (m *mockCore) GetDynamicFlags(_ uintptr) []drawable.DynamicFlag { return nil }
func (m *mockCore) GetOpacities(_ uintptr) []float32           { return nil }
func (m *mockCore) GetVertexPositions(_ uintptr) [][]drawable.Vector2 {
	return nil
}
func (m *mockCore) GetDrawables(_ uintptr) []drawable.Drawable { return nil }
func (m *mockCore) GetParameters(_ uintptr) []parameter.Parameter {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ps []parameter.Parameter
	for id, v := range m.params {
		ps = append(ps, parameter.Parameter{Id: id, Current: v})
	}
	return ps
}
func (m *mockCore) GetParameterValue(_ uintptr, id string) float32 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.params[id]
}
func (m *mockCore) SetParameterValue(_ uintptr, id string, value float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.params[id] = value
}
func (m *mockCore) GetPartIds(_ uintptr) []string { return nil }
func (m *mockCore) GetPartOpacities(_ uintptr) []float32 { return nil }
func (m *mockCore) SetPartOpacity(_ uintptr, id string, value float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.partOpas[id] = value
}
func (m *mockCore) GetSortedDrawableIndices(_ uintptr) []int      { return nil }
func (m *mockCore) GetCanvasInfo(_ uintptr) (drawable.Vector2, drawable.Vector2, float32) {
	return drawable.Vector2{}, drawable.Vector2{}, 0
}
func (m *mockCore) Update(_ uintptr) {}

func TestMotionManagerStart(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	m := Motion{
		File: "test.motion3.json",
		Meta: Meta{Duration: 5.0},
	}

	id := mm.Start(m)
	if id != 1 {
		t.Errorf("first motion id = %v, want 1", id)
	}

	id2 := mm.Start(m)
	if id2 != 2 {
		t.Errorf("second motion id = %v, want 2", id2)
	}
}

func TestMotionManagerClose(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	id := mm.Start(m)

	mm.Close(id)

	if len(mm.queue) != 0 {
		t.Errorf("queue should be empty after close, has %d items", len(mm.queue))
	}
}

func TestMotionManagerCloseNonExistent(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	mm.Start(m)

	// Closing a non-existent id should not crash or remove other entries
	mm.Close(999)
	if len(mm.queue) != 1 {
		t.Errorf("queue should still have 1 item, has %d", len(mm.queue))
	}
}

func TestMotionManagerReset(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	id := mm.Start(m)

	// Simulate some progress
	mm.queue[0].currentTime = 3.0

	mm.Reset(id)
	if mm.queue[0].currentTime != 0 {
		t.Errorf("currentTime after reset = %v, want 0", mm.queue[0].currentTime)
	}
}

func TestMotionManagerResetNonExistent(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	mm.Start(m)

	// Should not crash
	mm.Reset(999)
	if len(mm.queue) != 1 {
		t.Errorf("queue should still have 1 item, has %d", len(mm.queue))
	}
}

func TestMotionManagerUpdateEmpty(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	// Should not panic on empty queue
	mm.Update(0.016)
}

func TestMotionManagerUpdateFinished(t *testing.T) {
	core := newMockMotionCore()
	var finishedId int
	mm := NewMotionManager(core, 0, func(id int) { finishedId = id })

	m := Motion{Meta: Meta{Duration: 2.0}}
	id := mm.Start(m)

	// Advance past duration
	mm.Update(3.0)

	if finishedId != id {
		t.Errorf("finishedId = %v, want %v", finishedId, id)
	}
}

func TestMotionManagerUpdateWithParameterCurve(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	// Set initial parameter value
	core.SetParameterValue(0, "ParamAngleX", 0.0)

	m := Motion{
		Meta:        Meta{Duration: 5.0},
		FadeInTime:  0,
		FadeOutTime: 0,
		Curves: []Curve{
			{
				Target: "Parameter",
				Id:     "ParamAngleX",
				Segments: []Segment{
					{
						Type:   Linear,
						Points: []Point{{Time: 0, Value: 0}, {Time: 5, Value: 30}},
					},
				},
			},
		},
	}

	mm.Start(m)
	mm.Update(2.5)

	// At t=2.5 in a linear 0->30 curve over 0->5, value should be 15
	got := core.GetParameterValue(0, "ParamAngleX")
	if got < 14.9 || got > 15.1 {
		t.Errorf("ParamAngleX = %v, want ~15", got)
	}
}

func TestMotionManagerUpdateWithPartOpacity(t *testing.T) {
	core := newMockMotionCore()
	mm := NewMotionManager(core, 0, func(int) {})

	m := Motion{
		Meta:        Meta{Duration: 5.0},
		FadeInTime:  0,
		FadeOutTime: 0,
		Curves: []Curve{
			{
				Target: "PartOpacity",
				Id:     "PartArm",
				Segments: []Segment{
					{
						Type:   Linear,
						Points: []Point{{Time: 0, Value: 0}, {Time: 5, Value: 1}},
					},
				},
			},
		},
	}

	mm.Start(m)
	mm.Update(2.5)

	got := core.partOpas["PartArm"]
	if got < 0.49 || got > 0.51 {
		t.Errorf("PartArm opacity = %v, want ~0.5", got)
	}
}
