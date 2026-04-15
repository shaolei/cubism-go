package blink

import (
	"testing"

	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
)

// mockCore implements core.Core for testing BlinkManager
type mockCore struct {
	setValues map[string]float32
}

func newMockCore() *mockCore {
	return &mockCore{setValues: make(map[string]float32)}
}

func (m *mockCore) SetParameterValue(_ uintptr, id string, value float32) {
	m.setValues[id] = value
}
func (m *mockCore) SetParameterValueByIndex(_ uintptr, _ int, _ float32) {}
func (m *mockCore) GetParameterValueByIndex(_ uintptr, _ int) float32    { return 0 }
func (m *mockCore) GetParameterIds(_ uintptr) []string                   { return nil }
func (m *mockCore) SetPartOpacityByIndex(_ uintptr, _ int, _ float32)    {}
func (m *mockCore) GetPartOpacityByIndex(_ uintptr, _ int) float32      { return 0 }
func (m *mockCore) LoadMoc(_ string) (moc.Moc, error)                    { return moc.Moc{}, nil }
func (m *mockCore) GetVersion() string                                    { return "" }
func (m *mockCore) GetDynamicFlags(_ uintptr) []drawable.DynamicFlag      { return nil }
func (m *mockCore) GetOpacities(_ uintptr) []float32                     { return nil }
func (m *mockCore) GetVertexPositions(_ uintptr) [][]drawable.Vector2    { return nil }
func (m *mockCore) GetDrawables(_ uintptr) []drawable.Drawable           { return nil }
func (m *mockCore) GetParameters(_ uintptr) []parameter.Parameter        { return nil }
func (m *mockCore) GetParameterValue(_ uintptr, _ string) float32        { return 0 }
func (m *mockCore) GetParameterCount(_ uintptr) int                          { return 0 }
func (m *mockCore) GetParameterValues(_ uintptr) []float32                   { return nil }
func (m *mockCore) GetParameterMinimumValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetParameterMaximumValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetParameterDefaultValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetPartIds(_ uintptr) []string                        { return nil }
func (m *mockCore) GetPartOpacities(_ uintptr) []float32                 { return nil }
func (m *mockCore) SetPartOpacity(_ uintptr, _ string, _ float32)        {}
func (m *mockCore) GetSortedDrawableIndices(_ uintptr) []int             { return nil }
func (m *mockCore) GetCanvasInfo(_ uintptr) (drawable.Vector2, drawable.Vector2, float32) {
	return drawable.Vector2{}, drawable.Vector2{}, 0
}
func (m *mockCore) Update(_ uintptr) {}

func TestNewBlinkManager(t *testing.T) {
	t.Parallel()

	core := newMockCore()
	ids := []string{"ParamEyeLOpen", "ParamEyeROpen"}
	b := NewBlinkManager(core, 0, ids)

	if b.state != EyeStateFirst {
		t.Errorf("initial state = %v, want EyeStateFirst", b.state)
	}
	if b.interval != 4.0 {
		t.Errorf("interval = %v, want 4.0", b.interval)
	}
	if b.closing != 0.1 {
		t.Errorf("closing = %v, want 0.1", b.closing)
	}
	if b.opening != 0.15 {
		t.Errorf("opening = %v, want 0.15", b.opening)
	}
}

func TestBlinkFirstStateTransition(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"ParamEyeLOpen"})

	b.Update(0.016)

	if b.state != EyeStateInterval {
		t.Errorf("after first update, state = %v, want EyeStateInterval", b.state)
	}
	if core.setValues["ParamEyeLOpen"] != 1.0 {
		t.Errorf("first state value = %v, want 1.0", core.setValues["ParamEyeLOpen"])
	}
}

func TestBlinkIntervalBeforeNextBlink(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"ParamEyeLOpen"})

	// First update
	b.Update(0.016)
	// Set next blink far in the future
	b.nextBlinkingTime = b.currentTime + 100.0

	core.setValues = make(map[string]float32)
	b.Update(0.016)

	if core.setValues["ParamEyeLOpen"] != 1.0 {
		t.Errorf("interval value = %v, want 1.0", core.setValues["ParamEyeLOpen"])
	}
	if b.state != EyeStateInterval {
		t.Errorf("should still be in interval, got %v", b.state)
	}
}

func TestBlinkClosingTransition(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"ParamEyeLOpen"})
	b.closing = 0.1

	// Transition to interval
	b.Update(0.016)
	// Force blink now
	b.nextBlinkingTime = b.currentTime

	b.Update(0.016)
	if b.state != EyeStateClosing {
		t.Errorf("should be closing, got %v", b.state)
	}

	// Complete closing
	b.Update(b.closing)
	if b.state != EyeStateClosed {
		t.Errorf("should be closed, got %v", b.state)
	}
}

func TestBlinkFullCycle(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"ParamEyeLOpen"})
	b.closing = 0.1
	b.opening = 0.15

	// First -> Interval
	b.Update(0.016)
	// Trigger closing
	b.nextBlinkingTime = b.currentTime
	b.Update(0.016)
	if b.state != EyeStateClosing {
		t.Fatalf("should be closing, got %v", b.state)
	}

	// Closing -> Closed
	b.Update(b.closing)
	if b.state != EyeStateClosed {
		t.Fatalf("should be closed, got %v", b.state)
	}

	// Closed -> Opening
	b.Update(b.closing)
	if b.state != EyeStateOpening {
		t.Fatalf("should be opening, got %v", b.state)
	}

	// Opening -> Interval
	b.Update(b.opening)
	if b.state != EyeStateInterval {
		t.Errorf("should be back in interval, got %v", b.state)
	}
}

func TestBlinkClosedValueIsZero(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"ParamEyeLOpen"})
	b.closing = 0.1

	// Get to closed state
	b.Update(0.016)
	b.nextBlinkingTime = b.currentTime
	b.Update(0.016)         // -> closing
	b.Update(b.closing)     // -> closed

	core.setValues = make(map[string]float32)
	b.Update(0.001) // in closed state
	if core.setValues["ParamEyeLOpen"] != 0.0 {
		t.Errorf("closed state value = %v, want 0.0", core.setValues["ParamEyeLOpen"])
	}
}

func TestBlinkMultipleIds(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"L", "R"})

	b.Update(0.016)

	if core.setValues["L"] != 1.0 {
		t.Errorf("L = %v, want 1.0", core.setValues["L"])
	}
	if core.setValues["R"] != 1.0 {
		t.Errorf("R = %v, want 1.0", core.setValues["R"])
	}
}

func TestDetermineNextBlinkingTiming(t *testing.T) {
	core := newMockCore()
	b := NewBlinkManager(core, 0, []string{"L"})
	b.currentTime = 10.0

	nextTime := b.DetermineNextBlinkingTiming()

	minExpected := 10.0
	maxExpected := 10.0 + (2.0*4.0 - 1.0) // 17.0
	if nextTime < minExpected || nextTime > maxExpected {
		t.Errorf("nextBlinkingTime = %v, want between %v and %v", nextTime, minExpected, maxExpected)
	}
}
