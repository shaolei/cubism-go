package look

import (
	"testing"

	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
)

type mockCore struct {
	params map[string]float32
}

func newMockCore() *mockCore {
	return &mockCore{params: make(map[string]float32)}
}

func (m *mockCore) LoadMoc(_ string) (moc.Moc, error)                              { return moc.Moc{}, nil }
func (m *mockCore) GetVersion() string                                              { return "" }
func (m *mockCore) GetDynamicFlags(_ uintptr) []drawable.DynamicFlag                { return nil }
func (m *mockCore) GetOpacities(_ uintptr) []float32                                { return nil }
func (m *mockCore) GetVertexPositions(_ uintptr) [][]drawable.Vector2               { return nil }
func (m *mockCore) GetDrawables(_ uintptr) []drawable.Drawable                      { return nil }
func (m *mockCore) GetParameters(_ uintptr) []parameter.Parameter                   { return nil }
func (m *mockCore) GetParameterIds(_ uintptr) []string                              { return nil }
func (m *mockCore) GetParameterValue(_ uintptr, id string) float32   { return m.params[id] }
func (m *mockCore) SetParameterValue(_ uintptr, id string, v float32) { m.params[id] = v }
func (m *mockCore) GetParameterValueByIndex(_ uintptr, _ int) float32 { return 0 }
func (m *mockCore) SetParameterValueByIndex(_ uintptr, _ int, _ float32) {}
func (m *mockCore) GetParameterCount(_ uintptr) int                          { return 0 }
func (m *mockCore) GetParameterValues(_ uintptr) []float32                   { return nil }
func (m *mockCore) GetParameterMinimumValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetParameterMaximumValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetParameterDefaultValues(_ uintptr) []float32            { return nil }
func (m *mockCore) GetPartIds(_ uintptr) []string                       { return nil }
func (m *mockCore) GetPartOpacities(_ uintptr) []float32                { return nil }
func (m *mockCore) SetPartOpacity(_ uintptr, _ string, _ float32)       {}
func (m *mockCore) SetPartOpacityByIndex(_ uintptr, _ int, _ float32)   {}
func (m *mockCore) GetPartOpacityByIndex(_ uintptr, _ int) float32      { return 0 }
func (m *mockCore) GetSortedDrawableIndices(_ uintptr) []int            { return nil }
func (m *mockCore) GetCanvasInfo(_ uintptr) (drawable.Vector2, drawable.Vector2, float32) {
	return drawable.Vector2{}, drawable.Vector2{}, 0
}
func (m *mockCore) Update(_ uintptr) {}

func TestLookSetTarget(t *testing.T) {
	c := newMockCore()
	l := NewLookManager(c, 0)

	l.SetTarget(0.5, -0.3)
	x, y := l.GetTarget()
	if x != 0.5 || y != -0.3 {
		t.Errorf("target = (%v, %v), want (0.5, -0.3)", x, y)
	}
}

func TestLookUpdateWithFactorX(t *testing.T) {
	c := newMockCore()
	l := NewLookManager(c, 0)

	l.SetParameters([]LookParameterData{
		{ParameterId: "ParamAngleX", FactorX: 30.0, FactorY: 0.0, FactorXY: 0.0},
	})

	l.SetTarget(1.0, 0.0)
	l.Update(0.016)

	// Expected: FactorX * dragX = 30.0 * 1.0 = 30.0
	if c.params["ParamAngleX"] != 30.0 {
		t.Errorf("ParamAngleX = %v, want 30.0", c.params["ParamAngleX"])
	}
}

func TestLookUpdateWithFactorY(t *testing.T) {
	c := newMockCore()
	l := NewLookManager(c, 0)

	l.SetParameters([]LookParameterData{
		{ParameterId: "ParamAngleY", FactorX: 0.0, FactorY: 30.0, FactorXY: 0.0},
	})

	l.SetTarget(0.0, 1.0)
	l.Update(0.016)

	// Expected: FactorY * dragY = 30.0 * 1.0 = 30.0
	if c.params["ParamAngleY"] != 30.0 {
		t.Errorf("ParamAngleY = %v, want 30.0", c.params["ParamAngleY"])
	}
}

func TestLookUpdateWithFactorXY(t *testing.T) {
	c := newMockCore()
	l := NewLookManager(c, 0)

	l.SetParameters([]LookParameterData{
		{ParameterId: "ParamTest", FactorX: 1.0, FactorY: 1.0, FactorXY: 10.0},
	})

	l.SetTarget(0.5, 0.5)
	l.Update(0.016)

	// Expected: FactorX*0.5 + FactorY*0.5 + FactorXY*(0.5*0.5)
	// = 0.5 + 0.5 + 10.0*0.25 = 1.0 + 2.5 = 3.5
	if c.params["ParamTest"] != 3.5 {
		t.Errorf("ParamTest = %v, want 3.5", c.params["ParamTest"])
	}
}

func TestLookUpdateWithZeroDrag(t *testing.T) {
	c := newMockCore()
	l := NewLookManager(c, 0)

	l.SetParameters([]LookParameterData{
		{ParameterId: "ParamAngleX", FactorX: 30.0, FactorY: 0.0, FactorXY: 0.0},
	})

	l.SetTarget(0.0, 0.0)
	l.Update(0.016)

	// Expected: 0.0 (no offset from center)
	if c.params["ParamAngleX"] != 0.0 {
		t.Errorf("ParamAngleX = %v, want 0.0", c.params["ParamAngleX"])
	}
}

func TestLookDefaultTarget(t *testing.T) {
	c := newMockCore()
	l := NewLookManager(c, 0)

	x, y := l.GetTarget()
	if x != 0 || y != 0 {
		t.Errorf("default target = (%v, %v), want (0, 0)", x, y)
	}
}
