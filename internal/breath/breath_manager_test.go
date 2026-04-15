package breath

import (
	"math"
	"testing"

	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
	"github.com/shaolei/cubism-go/internal/id"
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
func (m *mockCore) GetParameterIds(_ uintptr) []string {
	return []string{"ParamAngleX", "ParamAngleY", "ParamAngleZ", "ParamBodyAngleX", "ParamBreath"}
}
func (m *mockCore) GetParameterValue(_ uintptr, id string) float32   { return m.params[id] }
func (m *mockCore) SetParameterValue(_ uintptr, id string, v float32) { m.params[id] = v }
func (m *mockCore) GetParameterValueByIndex(_ uintptr, _ int) float32 { return 0 }
func (m *mockCore) SetParameterValueByIndex(_ uintptr, idx int, v float32) {
	ids := []string{"ParamAngleX", "ParamAngleY", "ParamAngleZ", "ParamBodyAngleX", "ParamBreath"}
	if idx >= 0 && idx < len(ids) {
		m.params[ids[idx]] = v
	}
}
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

func TestDefaultBreathParameters(t *testing.T) {
	t.Parallel()
	params := DefaultBreathParameters()
	if len(params) != 5 {
		t.Errorf("default params count = %d, want 5", len(params))
	}
	// Verify ParamBreath exists
	found := false
	for _, p := range params {
		if p.ParameterId == "ParamBreath" {
			found = true
			if p.Peak != 0.5 {
				t.Errorf("ParamBreath peak = %v, want 0.5", p.Peak)
			}
		}
	}
	if !found {
		t.Error("ParamBreath not found in default parameters")
	}
}

func TestBreathUpdateSineWave(t *testing.T) {
	c := newMockCore()
	b := NewBreathManager(c, 0)

	// Set a single parameter with known values
	b.SetParameters([]BreathParameterData{
		{ParameterId: "ParamBreath", Offset: 0.0, Peak: 1.0, Cycle: 3.0, Weight: 1.0},
	})

	// At t=0: sin(0) = 0, so value = Offset + Peak * sin(0) = 0
	b.Update(0.0)
	if c.params["ParamBreath"] != 0.0 {
		t.Errorf("at t=0, value = %v, want 0.0", c.params["ParamBreath"])
	}

	// After one quarter cycle: sin(π/2) = 1.0
	// t = 0.75 (quarter of Cycle=3.0)
	// currentTime * 2π / Cycle = 0.75 * 2π / 3.0 = π/2
	b.Update(0.75)
	expected := float32(1.0) // Offset(0) + Peak(1) * sin(π/2) = 1.0
	if math.Abs(float64(c.params["ParamBreath"]-expected)) > 0.01 {
		t.Errorf("at t=0.75, value = %v, want ~%v", c.params["ParamBreath"], expected)
	}
}

func TestBreathWithOffset(t *testing.T) {
	c := newMockCore()
	b := NewBreathManager(c, 0)

	b.SetParameters([]BreathParameterData{
		{ParameterId: "ParamTest", Offset: 0.5, Peak: 0.0, Cycle: 1.0, Weight: 1.0},
	})

	// With Peak=0, value is always Offset regardless of time
	b.Update(1.0)
	if c.params["ParamTest"] != 0.5 {
		t.Errorf("with Peak=0, value = %v, want 0.5", c.params["ParamTest"])
	}
}

func TestBreathWithWeight(t *testing.T) {
	c := newMockCore()
	b := NewBreathManager(c, 0)

	b.SetParameters([]BreathParameterData{
		{ParameterId: "ParamTest", Offset: 1.0, Peak: 0.0, Cycle: 1.0, Weight: 0.5},
	})

	// With Weight=0.5, applied value = 1.0 * 0.5 = 0.5
	b.Update(0.0)
	if c.params["ParamTest"] != 0.5 {
		t.Errorf("with Weight=0.5, value = %v, want 0.5", c.params["ParamTest"])
	}
}

func TestBreathWithIdManager(t *testing.T) {
	c := newMockCore()
	b := NewBreathManager(c, 0)

	paramIds := []string{"ParamBreath"}
	idMgr := id.NewCubismIdManager(paramIds, nil)
	b.SetIdManager(idMgr)

	// With idManager, SetParameterValueByIndex is used.
	// mockCore maps index 0 → "ParamBreath" via its internal ids list.
	// But the idManager has "ParamBreath" at index 0 too.
	// The issue is mockCore.SetParameterValueByIndex uses its own ids list
	// which may not include "ParamBreath". Fix by checking the actual behavior.

	b.SetParameters([]BreathParameterData{
		{ParameterId: "ParamAngleX", Offset: 1.0, Peak: 0.0, Cycle: 1.0, Weight: 1.0},
	})

	// ParamAngleX is at index 0 in mockCore's GetParameterIds
	b.Update(0.0)
	// idManager resolves "ParamAngleX" to InvalidHandle since it only has "ParamBreath"
	// So this should skip the parameter entirely
	if _, ok := c.params["ParamAngleX"]; ok {
		t.Error("should skip parameter not in idManager")
	}

	// Now test with matching idManager
	idMgr2 := id.NewCubismIdManager([]string{"ParamAngleX", "ParamAngleY"}, nil)
	b2 := NewBreathManager(c, 0)
	b2.SetIdManager(idMgr2)
	b2.SetParameters([]BreathParameterData{
		{ParameterId: "ParamAngleX", Offset: 1.0, Peak: 0.0, Cycle: 1.0, Weight: 1.0},
	})
	b2.Update(0.0)
	// index 0 in idMgr2 = "ParamAngleX", mockCore maps index 0 → "ParamAngleX"
	if c.params["ParamAngleX"] != 1.0 {
		t.Errorf("with idManager, value = %v, want 1.0", c.params["ParamAngleX"])
	}
}
