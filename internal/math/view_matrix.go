package cubismmath

// ViewMatrix handles 4x4 matrices for camera/projection control.
// Extends Matrix44 with screen-boundary-aware translation and scaling.
// Matches the official SDK's CubismViewMatrix.
type ViewMatrix struct {
	*Matrix44
	screenLeft   float32
	screenRight  float32
	screenTop    float32
	screenBottom float32
	maxLeft      float32
	maxRight     float32
	maxTop       float32
	maxBottom    float32
	maxScale     float32
	minScale     float32
}

// NewViewMatrix creates a new ViewMatrix with all bounds initialized to zero.
func NewViewMatrix() *ViewMatrix {
	return &ViewMatrix{
		Matrix44: NewMatrix44(),
	}
}

// AdjustTranslate applies a relative translation (x, y) clamped to the
// maximum screen boundaries. The movement is adjusted so that the model
// never moves beyond the allowed range.
func (v *ViewMatrix) AdjustTranslate(x, y float32) {
	// Clamp X translation
	if v.Tr[0]*v.maxLeft+(v.Tr[12]+x) > v.screenLeft {
		x = v.screenLeft - v.Tr[0]*v.maxLeft - v.Tr[12]
	}
	if v.Tr[0]*v.maxRight+(v.Tr[12]+x) < v.screenRight {
		x = v.screenRight - v.Tr[0]*v.maxRight - v.Tr[12]
	}

	// Clamp Y translation
	if v.Tr[5]*v.maxTop+(v.Tr[13]+y) < v.screenTop {
		y = v.screenTop - v.Tr[5]*v.maxTop - v.Tr[13]
	}
	if v.Tr[5]*v.maxBottom+(v.Tr[13]+y) > v.screenBottom {
		y = v.screenBottom - v.Tr[5]*v.maxBottom - v.Tr[13]
	}

	// Apply the translation: m = T * m
	tr1 := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, 0, 1,
	}
	Multiply(tr1[:], v.Tr[:], v.Tr[:])
}

// AdjustScale scales the view around the center point (cx, cy).
// The scale factor is clamped to the [minScale, maxScale] range.
// The transform is: m = T(cx,cy) * S(scale) * T(-cx,-cy) * m
func (v *ViewMatrix) AdjustScale(cx, cy, scale float32) {
	targetScale := scale * v.Tr[0]

	if targetScale < v.minScale {
		if v.Tr[0] > 0 {
			scale = v.minScale / v.Tr[0]
		}
	} else if targetScale > v.maxScale {
		if v.Tr[0] > 0 {
			scale = v.maxScale / v.Tr[0]
		}
	}

	// Translate to center
	tr1 := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		cx, cy, 0, 1,
	}
	// Scale
	tr2 := [16]float32{
		scale, 0, 0, 0,
		0, scale, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	// Translate back
	tr3 := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		-cx, -cy, 0, 1,
	}

	Multiply(tr3[:], v.Tr[:], v.Tr[:])
	Multiply(tr2[:], v.Tr[:], v.Tr[:])
	Multiply(tr1[:], v.Tr[:], v.Tr[:])
}

// SetScreenRect sets the logical screen boundaries.
func (v *ViewMatrix) SetScreenRect(left, right, bottom, top float32) {
	v.screenLeft = left
	v.screenRight = right
	v.screenTop = top
	v.screenBottom = bottom
}

// SetMaxScreenRect sets the maximum movable screen boundaries.
func (v *ViewMatrix) SetMaxScreenRect(left, right, bottom, top float32) {
	v.maxLeft = left
	v.maxRight = right
	v.maxTop = top
	v.maxBottom = bottom
}

// SetMaxScale sets the maximum scale factor.
func (v *ViewMatrix) SetMaxScale(maxScale float32) { v.maxScale = maxScale }

// SetMinScale sets the minimum scale factor.
func (v *ViewMatrix) SetMinScale(minScale float32) { v.minScale = minScale }

// GetMaxScale returns the maximum scale factor.
func (v *ViewMatrix) GetMaxScale() float32 { return v.maxScale }

// GetMinScale returns the minimum scale factor.
func (v *ViewMatrix) GetMinScale() float32 { return v.minScale }

// IsMaxScale returns true if the current scale is at or above the maximum.
func (v *ViewMatrix) IsMaxScale() bool { return v.GetScaleX() >= v.maxScale }

// IsMinScale returns true if the current scale is at or below the minimum.
func (v *ViewMatrix) IsMinScale() bool { return v.GetScaleX() <= v.minScale }

// GetScreenLeft returns the left edge of the screen rect.
func (v *ViewMatrix) GetScreenLeft() float32 { return v.screenLeft }

// GetScreenRight returns the right edge of the screen rect.
func (v *ViewMatrix) GetScreenRight() float32 { return v.screenRight }

// GetScreenBottom returns the bottom edge of the screen rect.
func (v *ViewMatrix) GetScreenBottom() float32 { return v.screenBottom }

// GetScreenTop returns the top edge of the screen rect.
func (v *ViewMatrix) GetScreenTop() float32 { return v.screenTop }

// GetMaxLeft returns the left edge of the max screen rect.
func (v *ViewMatrix) GetMaxLeft() float32 { return v.maxLeft }

// GetMaxRight returns the right edge of the max screen rect.
func (v *ViewMatrix) GetMaxRight() float32 { return v.maxRight }

// GetMaxBottom returns the bottom edge of the max screen rect.
func (v *ViewMatrix) GetMaxBottom() float32 { return v.maxBottom }

// GetMaxTop returns the top edge of the max screen rect.
func (v *ViewMatrix) GetMaxTop() float32 { return v.maxTop }
