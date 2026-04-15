package cubismmath

// ModelMatrix handles 4x4 matrices for setting model coordinates.
// Extends Matrix44 with width/height-aware positioning methods.
// Matches the official SDK's CubismModelMatrix.
type ModelMatrix struct {
	*Matrix44
	width  float32
	height float32
}

// NewModelMatrix creates a ModelMatrix with the given model width and height.
// The height is set to 2.0 (matching the SDK's default scaling where the model
// occupies a vertical range of [-1, 1] in NDC space).
func NewModelMatrix(w, h float32) *ModelMatrix {
	m := &ModelMatrix{
		Matrix44: NewMatrix44(),
		width:    w,
		height:   h,
	}
	m.SetHeight(2.0)
	return m
}

// SetWidth sets the model scale so that the model has the given width in logical coordinates.
// The scale is uniform (aspect ratio preserved): scaleY = scaleX.
func (m *ModelMatrix) SetWidth(w float32) {
	scaleX := w / m.width
	m.Scale(scaleX, scaleX)
}

// SetHeight sets the model scale so that the model has the given height in logical coordinates.
// The scale is uniform (aspect ratio preserved): scaleX = scaleY.
func (m *ModelMatrix) SetHeight(h float32) {
	scaleX := h / m.height
	m.Scale(scaleX, scaleX)
}

// SetPosition sets the absolute position of the model's top-left corner.
func (m *ModelMatrix) SetPosition(x, y float32) {
	m.Translate(x, y)
}

// SetCenterPosition sets the model so its center is at (x, y).
// Must be called after SetWidth/SetHeight to ensure correct scaling.
func (m *ModelMatrix) SetCenterPosition(x, y float32) {
	m.CenterX(x)
	m.CenterY(y)
}

// Top positions the model's top edge at y.
func (m *ModelMatrix) Top(y float32) { m.SetY(y) }

// Bottom positions the model's bottom edge at y.
func (m *ModelMatrix) Bottom(y float32) {
	h := m.height * m.GetScaleY()
	m.TranslateY(y - h)
}

// Left positions the model's left edge at x.
func (m *ModelMatrix) Left(x float32) { m.SetX(x) }

// Right positions the model's right edge at x.
func (m *ModelMatrix) Right(x float32) {
	w := m.width * m.GetScaleX()
	m.TranslateX(x - w)
}

// CenterX positions the model's horizontal center at x.
func (m *ModelMatrix) CenterX(x float32) {
	w := m.width * m.GetScaleX()
	m.TranslateX(x - w/2.0)
}

// CenterY positions the model's vertical center at y.
func (m *ModelMatrix) CenterY(y float32) {
	h := m.height * m.GetScaleY()
	m.TranslateY(y - h/2.0)
}

// SetX sets the absolute X position.
func (m *ModelMatrix) SetX(x float32) { m.TranslateX(x) }

// SetY sets the absolute Y position.
func (m *ModelMatrix) SetY(y float32) { m.TranslateY(y) }

// SetupFromLayout configures the model matrix from layout information.
// Layout keys: "width", "height", "x", "y", "center_x", "center_y",
// "top", "bottom", "left", "right".
// Width/height are applied first, then positioning (matching the SDK's two-pass order).
func (m *ModelMatrix) SetupFromLayout(layout map[string]float32) {
	// First pass: apply width and height (scaling)
	for key, value := range layout {
		switch key {
		case "width":
			m.SetWidth(value)
		case "height":
			m.SetHeight(value)
		}
	}

	// Second pass: apply positioning
	for key, value := range layout {
		switch key {
		case "x":
			m.SetX(value)
		case "y":
			m.SetY(value)
		case "center_x":
			m.CenterX(value)
		case "center_y":
			m.CenterY(value)
		case "top":
			m.Top(value)
		case "bottom":
			m.Bottom(value)
		case "left":
			m.Left(value)
		case "right":
			m.Right(value)
		}
	}
}
