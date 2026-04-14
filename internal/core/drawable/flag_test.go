package drawable

import "testing"

func TestParseConstantFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		flag                 uint8
		blendAdditive        bool
		blendMultiplicative  bool
		isDoubleSided        bool
		isInvertedMask       bool
	}{
		{"zero flag", 0, false, false, false, false},
		{"blend additive only", 1, true, false, false, false},
		{"blend multiplicative only", 2, false, true, false, false},
		{"double sided only", 4, false, false, true, false},
		{"inverted mask only", 8, false, false, false, true},
		{"additive + double sided", 5, true, false, true, false},
		{"all flags", 15, true, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseConstantFlag(tt.flag)
			if got.BlendAdditive != tt.blendAdditive {
				t.Errorf("BlendAdditive = %v, want %v", got.BlendAdditive, tt.blendAdditive)
			}
			if got.BlendMultiplicative != tt.blendMultiplicative {
				t.Errorf("BlendMultiplicative = %v, want %v", got.BlendMultiplicative, tt.blendMultiplicative)
			}
			if got.IsDoubleSided != tt.isDoubleSided {
				t.Errorf("IsDoubleSided = %v, want %v", got.IsDoubleSided, tt.isDoubleSided)
			}
			if got.IsInvertedMask != tt.isInvertedMask {
				t.Errorf("IsInvertedMask = %v, want %v", got.IsInvertedMask, tt.isInvertedMask)
			}
		})
	}
}

func TestParseDynamicFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		flag                    uint8
		isVisible               bool
		visibilityDidChange     bool
		opacityDidChange        bool
		drawOrderDidChange      bool
		renderOrderDidChange    bool
		vertexPositionsDidChange bool
		blendColorDidChange     bool
	}{
		{"zero flag - invisible", 0, false, false, false, false, false, false, false},
		{"visible only", 1, true, false, false, false, false, false, false},
		{"visibility changed", 2, false, true, false, false, false, false, false},
		{"opacity changed", 4, false, false, true, false, false, false, false},
		{"draw order changed", 8, false, false, false, true, false, false, false},
		{"render order changed", 16, false, false, false, false, true, false, false},
		{"vertex positions changed", 32, false, false, false, false, false, true, false},
		{"blend color changed", 64, false, false, false, false, false, false, true},
		{"visible + opacity + vertex", 37, true, false, true, false, false, true, false},
		{"all flags", 127, true, true, true, true, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseDynamicFlag(tt.flag)
			if got.IsVisible != tt.isVisible {
				t.Errorf("IsVisible = %v, want %v", got.IsVisible, tt.isVisible)
			}
			if got.VisibilityDidChange != tt.visibilityDidChange {
				t.Errorf("VisibilityDidChange = %v, want %v", got.VisibilityDidChange, tt.visibilityDidChange)
			}
			if got.OpacityDidChange != tt.opacityDidChange {
				t.Errorf("OpacityDidChange = %v, want %v", got.OpacityDidChange, tt.opacityDidChange)
			}
			if got.DrawOrderDidChange != tt.drawOrderDidChange {
				t.Errorf("DrawOrderDidChange = %v, want %v", got.DrawOrderDidChange, tt.drawOrderDidChange)
			}
			if got.RenderOrderDidChange != tt.renderOrderDidChange {
				t.Errorf("RenderOrderDidChange = %v, want %v", got.RenderOrderDidChange, tt.renderOrderDidChange)
			}
			if got.VertexPositionsDidChange != tt.vertexPositionsDidChange {
				t.Errorf("VertexPositionsDidChange = %v, want %v", got.VertexPositionsDidChange, tt.vertexPositionsDidChange)
			}
			if got.BlendColorDidChange != tt.blendColorDidChange {
				t.Errorf("BlendColorDidChange = %v, want %v", got.BlendColorDidChange, tt.blendColorDidChange)
			}
		})
	}
}
