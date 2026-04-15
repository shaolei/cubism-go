package cubismmath

import "math"

// Epsilon is a small value used for floating-point comparisons.
const Epsilon float32 = 0.00001

// Pi is the value of π as a float32.
const Pi float32 = 3.1415926535897932384626433832795

// Matrix44 represents a 4x4 transformation matrix in column-major order.
// The 16-element array stores columns sequentially:
//
//	_tr[0]  _tr[4]  _tr[8]   _tr[12]     | ScaleX    0       0      TranslateX |
//	_tr[1]  _tr[5]  _tr[9]   _tr[13]     |   0    ScaleY     0      TranslateY |
//	_tr[2]  _tr[6]  _tr[10]  _tr[14]     |   0       0       1         TZ       |
//	_tr[3]  _tr[7]  _tr[11]  _tr[15]     |   0       0       0         1        |
//
// Matches the official SDK's CubismMatrix44 layout.
type Matrix44 struct {
	Tr [16]float32
}

// NewMatrix44 creates a new Matrix44 initialized to the identity matrix.
func NewMatrix44() *Matrix44 {
	m := &Matrix44{}
	m.LoadIdentity()
	return m
}

// LoadIdentity sets the matrix to the 4x4 identity matrix.
func (m *Matrix44) LoadIdentity() {
	m.Tr = [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Multiply computes dst = a * b (standard matrix multiplication).
// All three arrays must have at least 16 elements.
func Multiply(a, b, dst []float32) {
	var c [16]float32
	const n = 4
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				c[j+i*4] += a[k+i*4] * b[j+k*4]
			}
		}
	}
	copy(dst, c[:])
}

// GetArray returns the matrix elements as a slice.
func (m *Matrix44) GetArray() []float32 {
	return m.Tr[:]
}

// SetMatrix copies 16 floats from tr into the matrix.
func (m *Matrix44) SetMatrix(tr []float32) {
	copy(m.Tr[:], tr)
}

// GetScaleX returns the scaling factor along the X-axis.
func (m *Matrix44) GetScaleX() float32 { return m.Tr[0] }

// GetScaleY returns the scaling factor along the Y-axis.
func (m *Matrix44) GetScaleY() float32 { return m.Tr[5] }

// GetTranslateX returns the translation along the X-axis.
func (m *Matrix44) GetTranslateX() float32 { return m.Tr[12] }

// GetTranslateY returns the translation along the Y-axis.
func (m *Matrix44) GetTranslateY() float32 { return m.Tr[13] }

// TransformX transforms an X coordinate: ScaleX * src + TranslateX.
func (m *Matrix44) TransformX(src float32) float32 {
	return m.Tr[0]*src + m.Tr[12]
}

// TransformY transforms a Y coordinate: ScaleY * src + TranslateY.
func (m *Matrix44) TransformY(src float32) float32 {
	return m.Tr[5]*src + m.Tr[13]
}

// InvertTransformX computes the inverse X transform: (src - TranslateX) / ScaleX.
func (m *Matrix44) InvertTransformX(src float32) float32 {
	return (src - m.Tr[12]) / m.Tr[0]
}

// InvertTransformY computes the inverse Y transform: (src - TranslateY) / ScaleY.
func (m *Matrix44) InvertTransformY(src float32) float32 {
	return (src - m.Tr[13]) / m.Tr[5]
}

// TranslateRelative applies a relative translation by multiplying
// a translation matrix T on the left: m = T * m.
func (m *Matrix44) TranslateRelative(x, y float32) {
	tr1 := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, 0, 1,
	}
	Multiply(tr1[:], m.Tr[:], m.Tr[:])
}

// Translate sets the absolute translation values.
func (m *Matrix44) Translate(x, y float32) {
	m.Tr[12] = x
	m.Tr[13] = y
}

// TranslateX sets the absolute X translation.
func (m *Matrix44) TranslateX(x float32) { m.Tr[12] = x }

// TranslateY sets the absolute Y translation.
func (m *Matrix44) TranslateY(y float32) { m.Tr[13] = y }

// ScaleRelative applies a relative scale by multiplying
// a scaling matrix S on the left: m = S * m.
func (m *Matrix44) ScaleRelative(x, y float32) {
	tr1 := [16]float32{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	Multiply(tr1[:], m.Tr[:], m.Tr[:])
}

// Scale sets the absolute scale values.
func (m *Matrix44) Scale(x, y float32) {
	m.Tr[0] = x
	m.Tr[5] = y
}

// MultiplyByMatrix multiplies the current matrix by m2 on the left: m = m2 * m.
func (m *Matrix44) MultiplyByMatrix(m2 *Matrix44) {
	Multiply(m2.Tr[:], m.Tr[:], m.Tr[:])
}

// GetInvert returns the inverse of the matrix.
// If the matrix is not invertible (determinant near zero), returns the identity matrix.
func (m *Matrix44) GetInvert() *Matrix44 {
	r00 := m.Tr[0]
	r10 := m.Tr[1]
	r20 := m.Tr[2]
	r01 := m.Tr[4]
	r11 := m.Tr[5]
	r21 := m.Tr[6]
	r02 := m.Tr[8]
	r12 := m.Tr[9]
	r22 := m.Tr[10]

	tx := m.Tr[12]
	ty := m.Tr[13]
	tz := m.Tr[14]

	det := r00*(r11*r22-r12*r21) -
		r01*(r10*r22-r12*r20) +
		r02*(r10*r21-r11*r20)

	dst := NewMatrix44()

	if absF(det) < Epsilon {
		dst.LoadIdentity()
		return dst
	}

	invDet := 1.0 / det

	inv00 := (r11*r22 - r12*r21) * invDet
	inv01 := -(r01*r22 - r02*r21) * invDet
	inv02 := (r01*r12 - r02*r11) * invDet
	inv10 := -(r10*r22 - r12*r20) * invDet
	inv11 := (r00*r22 - r02*r20) * invDet
	inv12 := -(r00*r12 - r02*r10) * invDet
	inv20 := (r10*r21 - r11*r20) * invDet
	inv21 := -(r00*r21 - r01*r20) * invDet
	inv22 := (r00*r11 - r01*r10) * invDet

	dst.Tr[0] = inv00
	dst.Tr[1] = inv10
	dst.Tr[2] = inv20
	dst.Tr[3] = 0
	dst.Tr[4] = inv01
	dst.Tr[5] = inv11
	dst.Tr[6] = inv21
	dst.Tr[7] = 0
	dst.Tr[8] = inv02
	dst.Tr[9] = inv12
	dst.Tr[10] = inv22
	dst.Tr[11] = 0

	dst.Tr[12] = -(inv00*tx + inv01*ty + inv02*tz)
	dst.Tr[13] = -(inv10*tx + inv11*ty + inv12*tz)
	dst.Tr[14] = -(inv20*tx + inv21*ty + inv22*tz)
	dst.Tr[15] = 1

	return dst
}

// absF returns the absolute value of a float32.
func absF(x float32) float32 {
	return float32(math.Abs(float64(x)))
}
