package physics

import (
	"math"
	"testing"

	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
	"github.com/shaolei/cubism-go/internal/id"
	"github.com/shaolei/cubism-go/internal/model"
)

// mockCore implements core.Core for physics testing.
// It stores parameter values in a slice for direct access (matching the physics engine's
// usage of GetParameterValues/SetParameterValueByIndex).
type mockCore struct {
	paramIds    []string
	paramValues []float32
	paramMin    []float32
	paramMax    []float32
	paramDef    []float32
	partOpas    []float32
}

func newMockPhysicsCore() *mockCore {
	ids := []string{"ParamAngleX", "ParamAngleY", "ParamAngleZ", "ParamBodyAngleX", "ParamBreath"}
	count := len(ids)
	c := &mockCore{
		paramIds:    ids,
		paramValues: make([]float32, count),
		paramMin:    make([]float32, count),
		paramMax:    make([]float32, count),
		paramDef:    make([]float32, count),
	}
	for i := range ids {
		c.paramMin[i] = -30.0
		c.paramMax[i] = 30.0
		c.paramDef[i] = 0.0
	}
	return c
}

func (m *mockCore) LoadMoc(_ string) (moc.Moc, error)           { return moc.Moc{}, nil }
func (m *mockCore) GetVersion() string                           { return "" }
func (m *mockCore) GetDynamicFlags(_ uintptr) []drawable.DynamicFlag { return nil }
func (m *mockCore) GetOpacities(_ uintptr) []float32             { return nil }
func (m *mockCore) GetVertexPositions(_ uintptr) [][]drawable.Vector2 { return nil }
func (m *mockCore) GetDrawables(_ uintptr) []drawable.Drawable   { return nil }
func (m *mockCore) GetParameters(_ uintptr) []parameter.Parameter {
	ps := make([]parameter.Parameter, len(m.paramIds))
	for i, id := range m.paramIds {
		ps[i] = parameter.Parameter{Id: id, Current: m.paramValues[i], Minimum: m.paramMin[i], Maximum: m.paramMax[i], Default: m.paramDef[i]}
	}
	return ps
}
func (m *mockCore) GetParameterIds(_ uintptr) []string   { return m.paramIds }
func (m *mockCore) GetParameterValue(_ uintptr, id string) float32 {
	for i, pid := range m.paramIds {
		if pid == id {
			return m.paramValues[i]
		}
	}
	return 0
}
func (m *mockCore) SetParameterValue(_ uintptr, id string, v float32) {
	for i, pid := range m.paramIds {
		if pid == id {
			m.paramValues[i] = v
			return
		}
	}
}
func (m *mockCore) GetParameterValueByIndex(_ uintptr, idx int) float32 {
	if idx >= 0 && idx < len(m.paramValues) {
		return m.paramValues[idx]
	}
	return 0
}
func (m *mockCore) SetParameterValueByIndex(_ uintptr, idx int, v float32) {
	if idx >= 0 && idx < len(m.paramValues) {
		m.paramValues[idx] = v
	}
}
func (m *mockCore) GetParameterCount(_ uintptr) int               { return len(m.paramIds) }
func (m *mockCore) GetParameterValues(_ uintptr) []float32        { return m.paramValues }
func (m *mockCore) GetParameterMinimumValues(_ uintptr) []float32 { return m.paramMin }
func (m *mockCore) GetParameterMaximumValues(_ uintptr) []float32 { return m.paramMax }
func (m *mockCore) GetParameterDefaultValues(_ uintptr) []float32 { return m.paramDef }
func (m *mockCore) GetPartIds(_ uintptr) []string                 { return nil }
func (m *mockCore) GetPartOpacities(_ uintptr) []float32          { return m.partOpas }
func (m *mockCore) SetPartOpacity(_ uintptr, _ string, _ float32) {}
func (m *mockCore) SetPartOpacityByIndex(_ uintptr, _ int, _ float32) {}
func (m *mockCore) GetPartOpacityByIndex(_ uintptr, _ int) float32 { return 0 }
func (m *mockCore) GetSortedDrawableIndices(_ uintptr) []int     { return nil }
func (m *mockCore) GetCanvasInfo(_ uintptr) (drawable.Vector2, drawable.Vector2, float32) {
	return drawable.Vector2{}, drawable.Vector2{}, 0
}
func (m *mockCore) Update(_ uintptr) {}

// createTestPhysicsJson creates a minimal physics3.json structure for testing.
func createTestPhysicsJson() model.PhysicsJson {
	return model.PhysicsJson{
		Version: 3,
		Meta: struct {
			PhysicsSettingCount int     `json:"PhysicsSettingCount"`
			TotalInputCount     int     `json:"TotalInputCount"`
			TotalOutputCount    int     `json:"TotalOutputCount"`
			VertexCount         int     `json:"VertexCount"`
			Fps                 float64 `json:"Fps"`
			EffectiveForces     struct {
				Gravity struct {
					X float64 `json:"X"`
					Y float64 `json:"Y"`
				} `json:"Gravity"`
				Wind struct {
					X float64 `json:"X"`
					Y float64 `json:"Y"`
				} `json:"Wind"`
			} `json:"EffectiveForces"`
			PhysicsDictionary []struct {
				Id   string `json:"Id"`
				Name string `json:"Name"`
			} `json:"PhysicsDictionary"`
		}{
			PhysicsSettingCount: 1,
			TotalInputCount:     1,
			TotalOutputCount:    1,
			VertexCount:         3,
		},
		PhysicsSettings: []struct {
			Id    string `json:"Id"`
			Input []struct {
				Source struct {
					Target string `json:"Target"`
					Id     string `json:"Id"`
				} `json:"Source"`
				Weight  float64 `json:"Weight"`
				Type    string  `json:"Type"`
				Reflect bool    `json:"Reflect"`
			} `json:"Input"`
			Output []struct {
				Destination struct {
					Target string `json:"Target"`
					Id     string `json:"Id"`
				} `json:"Destination"`
				VertexIndex int     `json:"VertexIndex"`
				Scale       float64 `json:"Scale"`
				Weight      float64 `json:"Weight"`
				Type        string  `json:"Type"`
				Reflect     bool    `json:"Reflect"`
			} `json:"Output"`
			Vertices []struct {
				Position struct {
					X float64 `json:"X"`
					Y float64 `json:"Y"`
				} `json:"Position"`
				Mobility     float64 `json:"Mobility"`
				Delay        float64 `json:"Delay"`
				Acceleration float64 `json:"Acceleration"`
				Radius       float64 `json:"Radius"`
			} `json:"Vertices"`
			Normalization struct {
				Position struct {
					Minimum float64 `json:"Minimum"`
					Default float64 `json:"Default"`
					Maximum float64 `json:"Maximum"`
				} `json:"Position"`
				Angle struct {
					Minimum float64 `json:"Minimum"`
					Default float64 `json:"Default"`
					Maximum float64 `json:"Maximum"`
				} `json:"Angle"`
			} `json:"Normalization"`
		}{
			{
				Id: "PhysicsSetting1",
				Input: []struct {
					Source struct {
						Target string `json:"Target"`
						Id     string `json:"Id"`
					} `json:"Source"`
					Weight  float64 `json:"Weight"`
					Type    string  `json:"Type"`
					Reflect bool    `json:"Reflect"`
				}{
					{
						Source:  struct {
							Target string `json:"Target"`
							Id     string `json:"Id"`
						}{Target: "Parameter", Id: "ParamAngleX"},
						Weight:  1.0,
						Type:    "X",
						Reflect: false,
					},
				},
				Output: []struct {
					Destination struct {
						Target string `json:"Target"`
						Id     string `json:"Id"`
					} `json:"Destination"`
					VertexIndex int     `json:"VertexIndex"`
					Scale       float64 `json:"Scale"`
					Weight      float64 `json:"Weight"`
					Type        string  `json:"Type"`
					Reflect     bool    `json:"Reflect"`
				}{
					{
						Destination: struct {
							Target string `json:"Target"`
							Id     string `json:"Id"`
						}{Target: "Parameter", Id: "ParamAngleZ"},
						VertexIndex: 2,
						Scale:       1.0,
						Weight:      1.0,
						Type:        "Angle",
						Reflect:     false,
					},
				},
				Vertices: []struct {
					Position struct {
						X float64 `json:"X"`
						Y float64 `json:"Y"`
					} `json:"Position"`
					Mobility     float64 `json:"Mobility"`
					Delay        float64 `json:"Delay"`
					Acceleration float64 `json:"Acceleration"`
					Radius       float64 `json:"Radius"`
				}{
					{Position: struct {
						X float64 `json:"X"`
						Y float64 `json:"Y"`
					}{X: 0, Y: 0}, Mobility: 1.0, Delay: 0.5, Acceleration: 1.0, Radius: 10.0},
					{Position: struct {
						X float64 `json:"X"`
						Y float64 `json:"Y"`
					}{X: 0, Y: 10}, Mobility: 1.0, Delay: 0.5, Acceleration: 1.0, Radius: 10.0},
					{Position: struct {
						X float64 `json:"X"`
						Y float64 `json:"Y"`
					}{X: 0, Y: 20}, Mobility: 1.0, Delay: 0.5, Acceleration: 1.0, Radius: 10.0},
				},
				Normalization: struct {
					Position struct {
						Minimum float64 `json:"Minimum"`
						Default float64 `json:"Default"`
						Maximum float64 `json:"Maximum"`
					} `json:"Position"`
					Angle struct {
						Minimum float64 `json:"Minimum"`
						Default float64 `json:"Default"`
						Maximum float64 `json:"Maximum"`
					} `json:"Angle"`
				}{
					Position: struct {
						Minimum float64 `json:"Minimum"`
						Default float64 `json:"Default"`
						Maximum float64 `json:"Maximum"`
					}{Minimum: -1, Default: 0, Maximum: 1},
					Angle: struct {
						Minimum float64 `json:"Minimum"`
						Default float64 `json:"Default"`
						Maximum float64 `json:"Maximum"`
					}{Minimum: -10, Default: 0, Maximum: 10},
				},
			},
		},
	}
}

// Fix the Meta fields that couldn't be initialized inline
func init() {
	json := createTestPhysicsJson()
	_ = json
}

func TestParsePhysicsJson(t *testing.T) {
	json := createTestPhysicsJson()
	rig := ParsePhysicsJson(json)

	if rig.SubRigCount != 1 {
		t.Errorf("SubRigCount = %d, want 1", rig.SubRigCount)
	}
	if len(rig.Inputs) != 1 {
		t.Errorf("Inputs count = %d, want 1", len(rig.Inputs))
	}
	if len(rig.Outputs) != 1 {
		t.Errorf("Outputs count = %d, want 1", len(rig.Outputs))
	}
	if len(rig.Particles) != 3 {
		t.Errorf("Particles count = %d, want 3", len(rig.Particles))
	}

	// Check input type
	if rig.Inputs[0].Type != PhysicsSourceX {
		t.Errorf("Input type = %v, want PhysicsSourceX", rig.Inputs[0].Type)
	}
	if rig.Inputs[0].SourceId != "ParamAngleX" {
		t.Errorf("Input source ID = %s, want ParamAngleX", rig.Inputs[0].SourceId)
	}

	// Check output type
	if rig.Outputs[0].Type != PhysicsSourceAngle {
		t.Errorf("Output type = %v, want PhysicsSourceAngle", rig.Outputs[0].Type)
	}
	if rig.Outputs[0].DestinationId != "ParamAngleZ" {
		t.Errorf("Output destination ID = %s, want ParamAngleZ", rig.Outputs[0].DestinationId)
	}

	// Check sub rig indices
	if rig.Settings[0].BaseInputIndex != 0 {
		t.Errorf("BaseInputIndex = %d, want 0", rig.Settings[0].BaseInputIndex)
	}
	if rig.Settings[0].BaseOutputIndex != 0 {
		t.Errorf("BaseOutputIndex = %d, want 0", rig.Settings[0].BaseOutputIndex)
	}
	if rig.Settings[0].BaseParticleIndex != 0 {
		t.Errorf("BaseParticleIndex = %d, want 0", rig.Settings[0].BaseParticleIndex)
	}
	if rig.Settings[0].ParticleCount != 3 {
		t.Errorf("ParticleCount = %d, want 3", rig.Settings[0].ParticleCount)
	}
}

func TestNewPhysicsManager(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()

	pm := NewPhysicsManager(c, 0, json)
	if pm == nil {
		t.Fatal("NewPhysicsManager returned nil")
	}

	if pm.rig == nil {
		t.Error("rig should not be nil")
	}
	if pm.rig.SubRigCount != 1 {
		t.Errorf("SubRigCount = %d, want 1", pm.rig.SubRigCount)
	}
}

func TestPhysicsManagerEvaluateDoesNotCrash(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()
	pm := NewPhysicsManager(c, 0, json)

	// Set up id manager
	paramIds := c.GetParameterIds(0)
	idMgr := id.NewCubismIdManager(paramIds, nil)
	pm.SetIdManager(idMgr)

	// Should not panic
	pm.Evaluate(0.016)
}

func TestPhysicsManagerEvaluateZeroDelta(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()
	pm := NewPhysicsManager(c, 0, json)

	// Zero or negative delta should be a no-op
	pm.Evaluate(0.0)
	pm.Evaluate(-0.1)
}

func TestPhysicsManagerReset(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()
	pm := NewPhysicsManager(c, 0, json)

	// Run some physics
	pm.Evaluate(0.016)

	// Reset should not panic
	pm.Reset()

	// After reset, options should be back to defaults
	opts := pm.GetOptions()
	if opts.Gravity.Y != -1.0 {
		t.Errorf("after reset, gravity Y = %v, want -1.0", opts.Gravity.Y)
	}
}

func TestPhysicsManagerSetOptions(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()
	pm := NewPhysicsManager(c, 0, json)

	pm.SetOptions(Options{
		Gravity: Vector2{X: 0.5, Y: -2.0},
		Wind:    Vector2{X: 1.0, Y: 0.5},
	})

	opts := pm.GetOptions()
	if opts.Gravity.X != 0.5 || opts.Gravity.Y != -2.0 {
		t.Errorf("gravity = (%v, %v), want (0.5, -2.0)", opts.Gravity.X, opts.Gravity.Y)
	}
	if opts.Wind.X != 1.0 || opts.Wind.Y != 0.5 {
		t.Errorf("wind = (%v, %v), want (1.0, 0.5)", opts.Wind.X, opts.Wind.Y)
	}
}

func TestPhysicsManagerStabilization(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()
	pm := NewPhysicsManager(c, 0, json)

	paramIds := c.GetParameterIds(0)
	idMgr := id.NewCubismIdManager(paramIds, nil)
	pm.SetIdManager(idMgr)

	// Stabilization should not panic
	pm.Stabilization()
}

func TestVector2Operations(t *testing.T) {
	v1 := Vector2{X: 3, Y: 4}
	v2 := Vector2{X: 1, Y: 2}

	// Add
	result := v1.Add(v2)
	if result.X != 4 || result.Y != 6 {
		t.Errorf("Add = (%v, %v), want (4, 6)", result.X, result.Y)
	}

	// Sub
	result = v1.Sub(v2)
	if result.X != 2 || result.Y != 2 {
		t.Errorf("Sub = (%v, %v), want (2, 2)", result.X, result.Y)
	}

	// Scale
	result = v1.Scale(2)
	if result.X != 6 || result.Y != 8 {
		t.Errorf("Scale = (%v, %v), want (6, 8)", result.X, result.Y)
	}

	// Length
	length := v1.Length()
	if math.Abs(float64(length)-5.0) > 0.001 {
		t.Errorf("Length = %v, want 5.0", length)
	}

	// Normalize
	normalized := v1.Normalize()
	if math.Abs(float64(normalized.X)-0.6) > 0.001 || math.Abs(float64(normalized.Y)-0.8) > 0.001 {
		t.Errorf("Normalize = (%v, %v), want (0.6, 0.8)", normalized.X, normalized.Y)
	}

	// Zero vector normalize
	zero := Vector2{}
	zeroNorm := zero.Normalize()
	if zeroNorm.X != 0 || zeroNorm.Y != 0 {
		t.Errorf("Zero normalize = (%v, %v), want (0, 0)", zeroNorm.X, zeroNorm.Y)
	}
}

func TestMathHelpers(t *testing.T) {
	// directionToRadian: two same-direction vectors should give 0
	from := Vector2{X: 0, Y: 1}
	to := Vector2{X: 0, Y: 1}
	rad := directionToRadian(from, to)
	if absFloat32(rad) > 0.001 {
		t.Errorf("same direction rad = %v, want 0", rad)
	}

	// radianToDirection: 0 radian should point up (sin(0)=0, cos(0)=1)
	dir := radianToDirection(0)
	if absFloat32(dir.X) > 0.001 || absFloat32(dir.Y-1.0) > 0.001 {
		t.Errorf("radianToDirection(0) = (%v, %v), want (0, 1)", dir.X, dir.Y)
	}

	// degreesToRadian
	rad = degreesToRadian(180)
	if absFloat32(rad-float32(math.Pi)) > 0.001 {
		t.Errorf("degreesToRadian(180) = %v, want π", rad)
	}
}

func TestNormalizeParameterValue(t *testing.T) {
	// Test middle value (default)
	result := normalizeParameterValue(0.0, -30.0, 30.0, 0.0, -1.0, 1.0, 0.0, false)
	// At default value, the result should be Default of normalized range * -1.0 (not inverted)
	// Actually: value=0, middleValue=0, paramValue=0 → sign=0 → result = middleNormValue = 0 → * -1.0 = 0
	if absFloat32(result) > 0.001 {
		t.Errorf("normalize at default = %v, want 0", result)
	}

	// Test max value
	result = normalizeParameterValue(30.0, -30.0, 30.0, 0.0, -1.0, 1.0, 0.0, false)
	// value=30, maxValue=30, middleValue=0, paramValue=30, sign=1
	// nLength = 1.0 - 0.0 = 1.0, pLength = 30 - 0 = 30
	// result = 30 * (1.0/30.0) + 0 = 1.0, not inverted → 1.0 * -1.0 = -1.0
	if absFloat32(result-(-1.0)) > 0.001 {
		t.Errorf("normalize at max (not inverted) = %v, want -1.0", result)
	}

	// Test max value with invert
	result = normalizeParameterValue(30.0, -30.0, 30.0, 0.0, -1.0, 1.0, 0.0, true)
	if absFloat32(result-1.0) > 0.001 {
		t.Errorf("normalize at max (inverted) = %v, want 1.0", result)
	}
}

func TestPhysicsManagerMultipleEvaluates(t *testing.T) {
	c := newMockPhysicsCore()
	json := createTestPhysicsJson()
	pm := NewPhysicsManager(c, 0, json)

	paramIds := c.GetParameterIds(0)
	idMgr := id.NewCubismIdManager(paramIds, nil)
	pm.SetIdManager(idMgr)

	// Set input parameter to a non-default value
	c.paramValues[0] = 10.0 // ParamAngleX = 10

	// Run multiple frames
	for i := 0; i < 100; i++ {
		pm.Evaluate(0.016)
	}
	// Should not crash or produce NaN
	for i, v := range c.paramValues {
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			t.Errorf("paramValues[%d] = %v, want finite number", i, v)
		}
	}
}
