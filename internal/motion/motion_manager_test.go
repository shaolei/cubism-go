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
	mu       sync.Mutex
	params   map[string]float32
	partOpas map[string]float32
}

func newMockMotionCore() *mockCore {
	return &mockCore{
		params:   make(map[string]float32),
		partOpas: make(map[string]float32),
	}
}

func (m *mockCore) LoadMoc(_ string) (moc.Moc, error)            { return moc.Moc{}, nil }
func (m *mockCore) GetVersion() string                             { return "" }
func (m *mockCore) GetDynamicFlags(_ uintptr) []drawable.DynamicFlag { return nil }
func (m *mockCore) GetOpacities(_ uintptr) []float32               { return nil }
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
func (m *mockCore) GetPartIds(_ uintptr) []string       { return nil }
func (m *mockCore) GetPartOpacities(_ uintptr) []float32 { return nil }
func (m *mockCore) SetPartOpacity(_ uintptr, id string, value float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.partOpas[id] = value
}
func (m *mockCore) GetParameterIds(_ uintptr) []string             { return nil }
func (m *mockCore) GetParameterValueByIndex(_ uintptr, _ int) float32 { return 0 }
func (m *mockCore) SetParameterValueByIndex(_ uintptr, _ int, _ float32) {}
func (m *mockCore) GetParameterCount(_ uintptr) int                          { return 0 }
func (m *mockCore) GetParameterValues(_ uintptr) []float32                   { return nil }
func (m *mockCore) GetParameterMinimumValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetParameterMaximumValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetParameterDefaultValues(_ uintptr) []float32            { return nil }
func (m *mockCore) SetPartOpacityByIndex(_ uintptr, _ int, _ float32)    {}
func (m *mockCore) GetPartOpacityByIndex(_ uintptr, _ int) float32      { return 0 }
func (m *mockCore) GetSortedDrawableIndices(_ uintptr) []int { return nil }
func (m *mockCore) GetCanvasInfo(_ uintptr) (drawable.Vector2, drawable.Vector2, float32) {
	return drawable.Vector2{}, drawable.Vector2{}, 0
}
func (m *mockCore) Update(_ uintptr) {}

func TestCubismMotionManagerStart(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{
		File: "test.motion3.json",
		Meta: Meta{Duration: 5.0},
	}

	id := mm.StartMotionWithPriority(m, false, PriorityNormal)
	if id != 1 {
		t.Errorf("first motion id = %v, want 1", id)
	}

	id2 := mm.StartMotionWithPriority(m, false, PriorityNormal)
	if id2 != 2 {
		t.Errorf("second motion id = %v, want 2", id2)
	}
}

func TestCubismMotionManagerClose(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	id := mm.StartMotionWithPriority(m, false, PriorityNormal)

	mm.Close(id)

	if len(mm.entries) != 0 {
		t.Errorf("entries should be empty after close, has %d items", len(mm.entries))
	}
}

func TestCubismMotionManagerCloseNonExistent(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	mm.StartMotionWithPriority(m, false, PriorityNormal)

	// Closing a non-existent id should not crash or remove other entries
	mm.Close(999)
	if len(mm.entries) != 1 {
		t.Errorf("entries should still have 1 item, has %d", len(mm.entries))
	}
}

func TestCubismMotionManagerReset(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	id := mm.StartMotionWithPriority(m, false, PriorityNormal)

	// Simulate some progress
	mm.Update(3.0)

	mm.Reset(id)
	// After reset, the entry's startTime should be updated so local time resets
	entry := mm.GetEntryById(id)
	if entry == nil {
		t.Fatal("entry should exist after reset")
	}
	localTime := entry.GetLocalTime(mm.userTime)
	if localTime != 0 {
		t.Errorf("localTime after reset = %v, want 0", localTime)
	}
}

func TestCubismMotionManagerResetNonExistent(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	mm.StartMotionWithPriority(m, false, PriorityNormal)

	// Should not crash
	mm.Reset(999)
	if len(mm.entries) != 1 {
		t.Errorf("entries should still have 1 item, has %d", len(mm.entries))
	}
}

func TestCubismMotionManagerUpdateEmpty(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	// Should not panic on empty queue
	mm.Update(0.016)
}

func TestCubismMotionManagerUpdateFinished(t *testing.T) {
	core := newMockMotionCore()
	var finishedId int
	mm := NewCubismMotionManager(core, 0, func(id int) { finishedId = id })

	m := Motion{Meta: Meta{Duration: 2.0}, FadeInTime: 0, FadeOutTime: 0}
	id := mm.StartMotionWithPriority(m, false, PriorityNormal)

	// First update starts the motion (entry starts at userTime=0.001)
	mm.Update(0.001)
	// Second update advances past duration
	mm.Update(3.0)

	if finishedId != id {
		t.Errorf("finishedId = %v, want %v", finishedId, id)
	}
}

func TestCubismMotionManagerUpdateWithParameterCurve(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	// Set initial parameter value
	core.SetParameterValue(0, "ParamAngleX", 0.0)

	m := Motion{
		Meta:        Meta{Duration: 5.0},
		FadeInTime:  0,
		FadeOutTime: 0,
		Curves: []Curve{
			{
				Target:      "Parameter",
				Id:          "ParamAngleX",
				FadeInTime:  -1.0,
				FadeOutTime: -1.0,
				Segments: []Segment{
					{
						Type:   Linear,
						Points: []Point{{Time: 0, Value: 0}, {Time: 5, Value: 30}},
					},
				},
			},
		},
	}

	mm.StartMotionWithPriority(m, false, PriorityNormal)
	// First update starts the motion
	mm.Update(0.001)
	// Second update advances to t≈2.5
	mm.Update(2.5)

	// At t≈2.5 in a linear 0->30 curve over 0->5, value should be ~15
	got := core.GetParameterValue(0, "ParamAngleX")
	if got < 14.5 || got > 15.5 {
		t.Errorf("ParamAngleX = %v, want ~15", got)
	}
}

func TestCubismMotionManagerUpdateWithPartOpacity(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

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

	mm.StartMotionWithPriority(m, false, PriorityNormal)
	// First update starts the motion
	mm.Update(0.001)
	// Second update advances to t≈2.5
	mm.Update(2.5)

	got := core.partOpas["PartArm"]
	if got < 0.45 || got > 0.55 {
		t.Errorf("PartArm opacity = %v, want ~0.5", got)
	}
}

func TestCubismMotionManagerPrioritySystem(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}

	// Start with idle priority
	id1 := mm.StartMotionWithPriority(m, false, PriorityIdle)
	if id1 == -1 {
		t.Error("should be able to start idle motion")
	}

	// Normal priority should be able to start
	id2 := mm.StartMotionWithPriority(m, false, PriorityNormal)
	if id2 == -1 {
		t.Error("normal priority should be able to start over idle")
	}

	// Force priority should always be able to start
	id3 := mm.StartMotionWithPriority(m, false, PriorityForce)
	if id3 == -1 {
		t.Error("force priority should always be able to start")
	}
}

func TestCubismMotionManagerPriorityResetOnFinish(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 2.0}, FadeInTime: 0, FadeOutTime: 0}
	mm.StartMotionWithPriority(m, false, PriorityNormal)

	if mm.GetCurrentPriority() != PriorityNormal {
		t.Errorf("current priority = %v, want %v", mm.GetCurrentPriority(), PriorityNormal)
	}

	// First update starts the motion
	mm.Update(0.001)
	// Second update advances past duration
	mm.Update(3.0)

	if mm.GetCurrentPriority() != PriorityNone {
		t.Errorf("priority should reset to none after motion finishes, got %v", mm.GetCurrentPriority())
	}
}

func TestCubismMotionManagerFadeOutOnNewMotion(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m1 := Motion{Meta: Meta{Duration: 5.0}, FadeInTime: 0, FadeOutTime: 1.0}
	m2 := Motion{Meta: Meta{Duration: 5.0}, FadeInTime: 0, FadeOutTime: 0}

	id1 := mm.StartMotionWithPriority(m1, false, PriorityNormal)
	mm.Update(1.0) // Advance to t=1.0

	// Start a second motion - should trigger fade-out on the first
	mm.StartMotionWithPriority(m2, false, PriorityNormal)

	entry := mm.GetEntryById(id1)
	if entry == nil {
		t.Fatal("first entry should still exist")
	}
	if !entry.isTriggeredFadeOut {
		t.Error("first entry should have fade-out triggered when new motion starts")
	}
}

func TestCubismMotionManagerStopAllMotions(t *testing.T) {
	core := newMockMotionCore()
	mm := NewCubismMotionManager(core, 0, func(int) {})

	m := Motion{Meta: Meta{Duration: 5.0}}
	mm.StartMotionWithPriority(m, false, PriorityNormal)
	mm.StartMotionWithPriority(m, false, PriorityNormal)

	mm.StopAllMotions()

	if !mm.IsFinished() {
		t.Error("all motions should be finished after StopAllMotions")
	}
}
