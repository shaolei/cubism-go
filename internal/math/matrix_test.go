package cubismmath

import (
	"math"
	"testing"
)

const tolerance float32 = 0.001

func approxEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < float64(tolerance)
}

// --- Matrix44 Tests ---

func TestNewMatrix44IsIdentity(t *testing.T) {
	m := NewMatrix44()
	expected := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	for i := range m.Tr {
		if m.Tr[i] != expected[i] {
			t.Errorf("Tr[%d] = %f, want %f", i, m.Tr[i], expected[i])
		}
	}
}

func TestLoadIdentity(t *testing.T) {
	m := NewMatrix44()
	m.Tr[0] = 5
	m.Tr[12] = 100
	m.LoadIdentity()
	if m.Tr[0] != 1 || m.Tr[12] != 0 || m.Tr[5] != 1 {
		t.Error("LoadIdentity did not reset to identity")
	}
}

func TestGetScaleAndTranslate(t *testing.T) {
	m := NewMatrix44()
	m.Tr[0] = 2.5   // ScaleX
	m.Tr[5] = 3.0   // ScaleY
	m.Tr[12] = 10.0 // TranslateX
	m.Tr[13] = 20.0 // TranslateY

	if m.GetScaleX() != 2.5 {
		t.Errorf("GetScaleX = %f, want 2.5", m.GetScaleX())
	}
	if m.GetScaleY() != 3.0 {
		t.Errorf("GetScaleY = %f, want 3.0", m.GetScaleY())
	}
	if m.GetTranslateX() != 10.0 {
		t.Errorf("GetTranslateX = %f, want 10.0", m.GetTranslateX())
	}
	if m.GetTranslateY() != 20.0 {
		t.Errorf("GetTranslateY = %f, want 20.0", m.GetTranslateY())
	}
}

func TestTransformXY(t *testing.T) {
	m := NewMatrix44()
	m.Scale(2, 3)
	m.Translate(10, 20)

	// TransformX = ScaleX * src + TranslateX = 2*5 + 10 = 20
	if !approxEqual(m.TransformX(5), 20) {
		t.Errorf("TransformX(5) = %f, want 20", m.TransformX(5))
	}
	// TransformY = ScaleY * src + TranslateY = 3*5 + 20 = 35
	if !approxEqual(m.TransformY(5), 35) {
		t.Errorf("TransformY(5) = %f, want 35", m.TransformY(5))
	}
}

func TestInvertTransformXY(t *testing.T) {
	m := NewMatrix44()
	m.Scale(2, 3)
	m.Translate(10, 20)

	// InvertTransformX = (src - TranslateX) / ScaleX = (20 - 10) / 2 = 5
	if !approxEqual(m.InvertTransformX(20), 5) {
		t.Errorf("InvertTransformX(20) = %f, want 5", m.InvertTransformX(20))
	}
	// InvertTransformY = (src - TranslateY) / ScaleY = (35 - 20) / 3 = 5
	if !approxEqual(m.InvertTransformY(35), 5) {
		t.Errorf("InvertTransformY(35) = %f, want 5", m.InvertTransformY(35))
	}
}

func TestTranslate(t *testing.T) {
	m := NewMatrix44()
	m.Translate(5, 10)
	if m.GetTranslateX() != 5 || m.GetTranslateY() != 10 {
		t.Errorf("Translate: got (%f, %f), want (5, 10)", m.GetTranslateX(), m.GetTranslateY())
	}
}

func TestTranslateRelative(t *testing.T) {
	m := NewMatrix44()
	m.Translate(10, 20)
	m.TranslateRelative(5, 5)
	// TranslateRelative does: T(5,5) * m
	// For identity-scale matrix, result should be TranslateX = 10+5=15, TranslateY = 20+5=25
	if !approxEqual(m.GetTranslateX(), 15) || !approxEqual(m.GetTranslateY(), 25) {
		t.Errorf("TranslateRelative: got (%f, %f), want (15, 25)", m.GetTranslateX(), m.GetTranslateY())
	}
}

func TestScale(t *testing.T) {
	m := NewMatrix44()
	m.Scale(3, 4)
	if m.GetScaleX() != 3 || m.GetScaleY() != 4 {
		t.Errorf("Scale: got (%f, %f), want (3, 4)", m.GetScaleX(), m.GetScaleY())
	}
}

func TestScaleRelative(t *testing.T) {
	m := NewMatrix44()
	m.Scale(2, 3)
	m.ScaleRelative(2, 2)
	// ScaleRelative does: S(2,2) * m
	// New ScaleX = 2*2 = 4, New ScaleY = 2*3 = 6
	if !approxEqual(m.GetScaleX(), 4) || !approxEqual(m.GetScaleY(), 6) {
		t.Errorf("ScaleRelative: got (%f, %f), want (4, 6)", m.GetScaleX(), m.GetScaleY())
	}
}

func TestMultiplyIdentity(t *testing.T) {
	a := NewMatrix44()
	b := NewMatrix44()
	dst := NewMatrix44()
	Multiply(a.Tr[:], b.Tr[:], dst.Tr[:])

	for i := range dst.Tr {
		expected := float32(0)
		if i == 0 || i == 5 || i == 10 || i == 15 {
			expected = 1
		}
		if dst.Tr[i] != expected {
			t.Errorf("Identity*Identity: Tr[%d] = %f, want %f", i, dst.Tr[i], expected)
		}
	}
}

func TestMultiplyTranslateScale(t *testing.T) {
	// T(10,20) * S(2,3) means: first scale, then translate
	// Result: ScaleX=2, ScaleY=3, TranslateX=10*2=20, TranslateY=20*3=60
	// Because column-major: T * S applies S first, then T, but the translation
	// column is affected by the scale in the multiplication.
	tm := NewMatrix44()
	tm.Translate(10, 20)
	sm := NewMatrix44()
	sm.Scale(2, 3)

	result := NewMatrix44()
	Multiply(tm.Tr[:], sm.Tr[:], result.Tr[:])

	if !approxEqual(result.GetScaleX(), 2) {
		t.Errorf("ScaleX = %f, want 2", result.GetScaleX())
	}
	if !approxEqual(result.GetScaleY(), 3) {
		t.Errorf("ScaleY = %f, want 3", result.GetScaleY())
	}
	if !approxEqual(result.GetTranslateX(), 20) {
		t.Errorf("TranslateX = %f, want 20", result.GetTranslateX())
	}
	if !approxEqual(result.GetTranslateY(), 60) {
		t.Errorf("TranslateY = %f, want 60", result.GetTranslateY())
	}
}

func TestGetInvertIdentity(t *testing.T) {
	m := NewMatrix44()
	inv := m.GetInvert()
	if !approxEqual(inv.GetScaleX(), 1) || !approxEqual(inv.GetScaleY(), 1) {
		t.Errorf("Identity inverse: got scale (%f, %f)", inv.GetScaleX(), inv.GetScaleY())
	}
	if !approxEqual(inv.GetTranslateX(), 0) || !approxEqual(inv.GetTranslateY(), 0) {
		t.Errorf("Identity inverse: got translate (%f, %f)", inv.GetTranslateX(), inv.GetTranslateY())
	}
}

func TestGetInvertScaledTranslated(t *testing.T) {
	m := NewMatrix44()
	m.Scale(2, 4)
	m.Translate(10, 20)

	inv := m.GetInvert()
	// Inverse of S(2,4)*T(10,20) should have ScaleX=0.5, ScaleY=0.25
	// and TranslateX=-5, TranslateY=-5
	if !approxEqual(inv.GetScaleX(), 0.5) {
		t.Errorf("Inverse ScaleX = %f, want 0.5", inv.GetScaleX())
	}
	if !approxEqual(inv.GetScaleY(), 0.25) {
		t.Errorf("Inverse ScaleY = %f, want 0.25", inv.GetScaleY())
	}
	if !approxEqual(inv.GetTranslateX(), -5) {
		t.Errorf("Inverse TranslateX = %f, want -5", inv.GetTranslateX())
	}
	if !approxEqual(inv.GetTranslateY(), -5) {
		t.Errorf("Inverse TranslateY = %f, want -5", inv.GetTranslateY())
	}
}

func TestGetInvertDegenerate(t *testing.T) {
	m := NewMatrix44()
	m.Scale(0, 0) // Degenerate: zero scale
	inv := m.GetInvert()
	// Should return identity for non-invertible matrices
	if !approxEqual(inv.GetScaleX(), 1) || !approxEqual(inv.GetScaleY(), 1) {
		t.Errorf("Degenerate inverse: got scale (%f, %f), want (1, 1)", inv.GetScaleX(), inv.GetScaleY())
	}
}

func TestSetMatrix(t *testing.T) {
	m := NewMatrix44()
	tr := [16]float32{
		2, 0, 0, 0,
		0, 3, 0, 0,
		0, 0, 1, 0,
		10, 20, 0, 1,
	}
	m.SetMatrix(tr[:])
	if m.GetScaleX() != 2 || m.GetScaleY() != 3 {
		t.Errorf("SetMatrix: got scale (%f, %f), want (2, 3)", m.GetScaleX(), m.GetScaleY())
	}
	if m.GetTranslateX() != 10 || m.GetTranslateY() != 20 {
		t.Errorf("SetMatrix: got translate (%f, %f), want (10, 20)", m.GetTranslateX(), m.GetTranslateY())
	}
}

// --- ModelMatrix Tests ---

func TestNewModelMatrix(t *testing.T) {
	m := NewModelMatrix(100, 200)
	// SetHeight(2.0) was called, so scale = 2.0/200 = 0.01
	expectedScale := float32(2.0) / 200.0
	if !approxEqual(m.GetScaleX(), expectedScale) {
		t.Errorf("ModelMatrix ScaleX = %f, want %f", m.GetScaleX(), expectedScale)
	}
	if !approxEqual(m.GetScaleY(), expectedScale) {
		t.Errorf("ModelMatrix ScaleY = %f, want %f", m.GetScaleY(), expectedScale)
	}
}

func TestModelMatrixSetWidth(t *testing.T) {
	m := NewModelMatrix(100, 200)
	m.SetWidth(50)
	// scaleX = 50/100 = 0.5, uniform so scaleY = 0.5
	if !approxEqual(m.GetScaleX(), 0.5) {
		t.Errorf("SetWidth(50) ScaleX = %f, want 0.5", m.GetScaleX())
	}
	if !approxEqual(m.GetScaleY(), 0.5) {
		t.Errorf("SetWidth(50) ScaleY = %f, want 0.5", m.GetScaleY())
	}
}

func TestModelMatrixSetHeight(t *testing.T) {
	m := NewModelMatrix(100, 200)
	m.SetHeight(4)
	// scaleX = 4/200 = 0.02, uniform so scaleY = 0.02
	if !approxEqual(m.GetScaleX(), 0.02) {
		t.Errorf("SetHeight(4) ScaleX = %f, want 0.02", m.GetScaleX())
	}
}

func TestModelMatrixSetPosition(t *testing.T) {
	m := NewModelMatrix(100, 200)
	m.SetPosition(5, 10)
	if m.GetTranslateX() != 5 || m.GetTranslateY() != 10 {
		t.Errorf("SetPosition: got (%f, %f), want (5, 10)", m.GetTranslateX(), m.GetTranslateY())
	}
}

func TestModelMatrixCenterPosition(t *testing.T) {
	m := NewModelMatrix(100, 200)
	// After SetHeight(2.0), scale = 0.01
	// CenterX(x): TranslateX = x - width*scale/2 = x - 100*0.01/2 = x - 0.5
	m.SetCenterPosition(0, 0)
	if !approxEqual(m.GetTranslateX(), -0.5) {
		t.Errorf("CenterX(0): TranslateX = %f, want -0.5", m.GetTranslateX())
	}
	if !approxEqual(m.GetTranslateY(), -1.0) {
		t.Errorf("CenterY(0): TranslateY = %f, want -1.0", m.GetTranslateY())
	}
}

func TestModelMatrixSetupFromLayout(t *testing.T) {
	m := NewModelMatrix(100, 200)
	layout := map[string]float32{
		"width":    50,
		"center_x": 400,
		"center_y": 300,
	}
	m.SetupFromLayout(layout)
	// SetWidth(50): scale = 50/100 = 0.5
	// CenterX(400): TranslateX = 400 - 100*0.5/2 = 400 - 25 = 375
	// CenterY(300): TranslateY = 300 - 200*0.5/2 = 300 - 50 = 250
	if !approxEqual(m.GetScaleX(), 0.5) {
		t.Errorf("SetupFromLayout ScaleX = %f, want 0.5", m.GetScaleX())
	}
	if !approxEqual(m.GetTranslateX(), 375) {
		t.Errorf("SetupFromLayout TranslateX = %f, want 375", m.GetTranslateX())
	}
	if !approxEqual(m.GetTranslateY(), 250) {
		t.Errorf("SetupFromLayout TranslateY = %f, want 250", m.GetTranslateY())
	}
}

// --- ViewMatrix Tests ---

func TestNewViewMatrix(t *testing.T) {
	v := NewViewMatrix()
	if v.GetMaxScale() != 0 || v.GetMinScale() != 0 {
		t.Errorf("ViewMatrix should have zero initial scale bounds")
	}
}

func TestViewMatrixSetScreenRect(t *testing.T) {
	v := NewViewMatrix()
	v.SetScreenRect(-1, 1, -1, 1)
	if v.GetScreenLeft() != -1 || v.GetScreenRight() != 1 {
		t.Errorf("SetScreenRect: left=%f right=%f", v.GetScreenLeft(), v.GetScreenRight())
	}
}

func TestViewMatrixSetMaxMinScale(t *testing.T) {
	v := NewViewMatrix()
	v.SetMaxScale(5)
	v.SetMinScale(0.5)
	if v.GetMaxScale() != 5 || v.GetMinScale() != 0.5 {
		t.Errorf("SetMaxMinScale: max=%f min=%f", v.GetMaxScale(), v.GetMinScale())
	}
}

func TestViewMatrixIsMaxMinScale(t *testing.T) {
	v := NewViewMatrix()
	v.SetMaxScale(2)
	v.SetMinScale(0.5)
	v.Scale(2, 2)
	if !v.IsMaxScale() {
		t.Error("Should be at max scale")
	}
	v.Scale(0.5, 0.5)
	if !v.IsMinScale() {
		t.Error("Should be at min scale")
	}
}

func TestViewMatrixAdjustScale(t *testing.T) {
	v := NewViewMatrix()
	v.SetMaxScale(5)
	v.SetMinScale(0.5)
	v.Scale(1, 1)

	// Scale up by 2x around center (0,0)
	v.AdjustScale(0, 0, 2)
	if !approxEqual(v.GetScaleX(), 2) {
		t.Errorf("After AdjustScale(0,0,2): ScaleX = %f, want 2", v.GetScaleX())
	}

	// Scale down from 2x by 0.5x => should become 1x
	v.AdjustScale(0, 0, 0.5)
	if !approxEqual(v.GetScaleX(), 1) {
		t.Errorf("After AdjustScale(0,0,0.5): ScaleX = %f, want 1", v.GetScaleX())
	}
}

func TestViewMatrixAdjustScaleClampMax(t *testing.T) {
	v := NewViewMatrix()
	v.SetMaxScale(2)
	v.SetMinScale(0.5)
	v.Scale(1, 1)

	// Try to scale up to 5x, should be clamped to 2x
	v.AdjustScale(0, 0, 5)
	if !approxEqual(v.GetScaleX(), 2) {
		t.Errorf("Clamped max: ScaleX = %f, want 2", v.GetScaleX())
	}
}

func TestViewMatrixAdjustScaleClampMin(t *testing.T) {
	v := NewViewMatrix()
	v.SetMaxScale(5)
	v.SetMinScale(1)
	v.Scale(2, 2)

	// Try to scale down to 0.1x => targetScale = 0.1 * 2 = 0.2 < minScale(1)
	// scale = minScale / Tr[0] = 1/2 = 0.5
	// After AdjustScale: new ScaleX = 2 * 0.5 = 1.0 (which is the minScale)
	v.AdjustScale(0, 0, 0.1)
	if !approxEqual(v.GetScaleX(), 1.0) {
		t.Errorf("Clamped min: ScaleX = %f, want 1.0", v.GetScaleX())
	}
}

// --- TargetPoint Tests ---

func TestTargetPointSetAndGet(t *testing.T) {
	tp := NewTargetPoint()
	tp.Set(0.5, -0.3)
	// Before any Update, X and Y should still be 0
	if tp.GetX() != 0 || tp.GetY() != 0 {
		t.Errorf("Before Update: got (%f, %f), want (0, 0)", tp.GetX(), tp.GetY())
	}
}

func TestTargetPointUpdateMovesTowardTarget(t *testing.T) {
	tp := NewTargetPoint()
	tp.Set(1.0, 0.0)

	// First Update just initializes lastTimeSeconds
	tp.Update(1.0 / 30.0)

	// After enough updates, face should move toward target
	for i := 0; i < 60; i++ {
		tp.Update(1.0 / 30.0)
	}

	if tp.GetX() <= 0 {
		t.Errorf("After 60 updates toward (1,0), X = %f, expected > 0", tp.GetX())
	}
}

func TestTargetPointUpdateReachesTarget(t *testing.T) {
	tp := NewTargetPoint()
	tp.Set(1.0, 0.0)

	// First Update initializes
	tp.Update(1.0 / 30.0)

	// Run many updates to reach target
	for i := 0; i < 300; i++ {
		tp.Update(1.0 / 30.0)
	}

	if !approxEqual(tp.GetX(), 1.0) {
		t.Errorf("After 300 updates, X = %f, want ~1.0", tp.GetX())
	}
}

func TestTargetPointNoMovementWhenAtTarget(t *testing.T) {
	tp := NewTargetPoint()
	tp.Set(0, 0)

	// Initialize
	tp.Update(1.0 / 30.0)

	// Update several times — should stay at 0
	for i := 0; i < 10; i++ {
		tp.Update(1.0 / 30.0)
	}

	if !approxEqual(tp.GetX(), 0) || !approxEqual(tp.GetY(), 0) {
		t.Errorf("At target: got (%f, %f), want (0, 0)", tp.GetX(), tp.GetY())
	}
}

func TestTargetPointDirectionChange(t *testing.T) {
	tp := NewTargetPoint()
	tp.Set(1.0, 0.0)

	// Initialize and move toward positive X
	tp.Update(1.0 / 30.0)
	for i := 0; i < 60; i++ {
		tp.Update(1.0 / 30.0)
	}
	positiveX := tp.GetX()

	// Change direction to negative X
	tp.Set(-1.0, 0.0)
	for i := 0; i < 120; i++ {
		tp.Update(1.0 / 30.0)
	}

	// Should have reversed direction
	if tp.GetX() >= positiveX {
		t.Errorf("After direction change, X = %f, should be less than %f", tp.GetX(), positiveX)
	}
}
