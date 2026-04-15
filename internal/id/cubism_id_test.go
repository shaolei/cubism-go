package id

import "testing"

func TestCubismIdHandleIsValid(t *testing.T) {
	t.Parallel()

	if InvalidHandle.IsValid() {
		t.Error("InvalidHandle should not be valid")
	}

	h := CubismIdHandle(0)
	if !h.IsValid() {
		t.Error("handle 0 should be valid")
	}

	h = CubismIdHandle(5)
	if !h.IsValid() {
		t.Error("handle 5 should be valid")
	}

	h = CubismIdHandle(-1)
	if h.IsValid() {
		t.Error("handle -1 should not be valid")
	}
}

func TestNewCubismIdManager(t *testing.T) {
	t.Parallel()

	paramIds := []string{"ParamA", "ParamB", "ParamC"}
	partIds := []string{"PartA", "PartB"}

	m := NewCubismIdManager(paramIds, partIds)

	if m.GetParameterCount() != 3 {
		t.Errorf("parameter count = %d, want 3", m.GetParameterCount())
	}
	if m.GetPartCount() != 2 {
		t.Errorf("part count = %d, want 2", m.GetPartCount())
	}
}

func TestGetParameterId(t *testing.T) {
	t.Parallel()

	paramIds := []string{"ParamEyeLOpen", "ParamEyeROpen", "ParamMouthOpen"}
	m := NewCubismIdManager(paramIds, nil)

	tests := []struct {
		id       string
		expected CubismIdHandle
	}{
		{"ParamEyeLOpen", 0},
		{"ParamEyeROpen", 1},
		{"ParamMouthOpen", 2},
		{"NonExistent", InvalidHandle},
	}

	for _, tt := range tests {
		got := m.GetParameterId(tt.id)
		if got != tt.expected {
			t.Errorf("GetParameterId(%q) = %d, want %d", tt.id, got, tt.expected)
		}
	}
}

func TestGetPartId(t *testing.T) {
	t.Parallel()

	partIds := []string{"PartArmL", "PartArmR"}
	m := NewCubismIdManager(nil, partIds)

	if m.GetPartId("PartArmL") != 0 {
		t.Errorf("PartArmL handle = %d, want 0", m.GetPartId("PartArmL"))
	}
	if m.GetPartId("PartArmR") != 1 {
		t.Errorf("PartArmR handle = %d, want 1", m.GetPartId("PartArmR"))
	}
	if m.GetPartId("NonExistent") != InvalidHandle {
		t.Error("non-existent part should return InvalidHandle")
	}
}

func TestEmptyManager(t *testing.T) {
	t.Parallel()

	m := NewCubismIdManager(nil, nil)

	if m.GetParameterCount() != 0 {
		t.Errorf("parameter count = %d, want 0", m.GetParameterCount())
	}
	if m.GetPartCount() != 0 {
		t.Errorf("part count = %d, want 0", m.GetPartCount())
	}
	if m.GetParameterId("anything") != InvalidHandle {
		t.Error("empty manager should return InvalidHandle for any lookup")
	}
}

func TestHandleAsIndex(t *testing.T) {
	t.Parallel()

	paramIds := []string{"A", "B", "C", "D", "E"}
	m := NewCubismIdManager(paramIds, nil)

	// Verify handles can be used directly as array indices
	for i, id := range paramIds {
		handle := m.GetParameterId(id)
		if int(handle) != i {
			t.Errorf("handle for %q = %d, want index %d", id, handle, i)
		}
	}
}
