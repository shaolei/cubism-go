package id

// CubismIdHandle represents a handle to a parameter/part ID,
// avoiding repeated string comparisons during parameter access.
// This matches the official SDK's CubismIdHandle design pattern.
//
// In the official SDK, CubismIdManager maps string IDs to integer handles.
// Here we use the parameter index directly as the handle, which allows
// O(1) access to parameter values instead of O(n) string comparison per call.
type CubismIdHandle int

// InvalidHandle represents an uninitialized or invalid ID handle.
const InvalidHandle CubismIdHandle = -1

// IsValid returns true if the handle is valid (non-negative index)
func (h CubismIdHandle) IsValid() bool {
	return h >= 0
}

// CubismIdManager manages the mapping from string IDs to CubismIdHandle values.
// It is initialized once when a model is loaded and provides fast lookups
// for the hot path (parameter get/set operations called every frame).
//
// The manager maintains two maps:
// - parameterIds: parameter string ID → parameter array index
// - partIds: part string ID → part array index
type CubismIdManager struct {
	parameterIds map[string]CubismIdHandle
	partIds      map[string]CubismIdHandle
}

// NewCubismIdManager creates a new ID manager with the given parameter and part IDs.
// parameterIdList and partIdList are the ordered lists of IDs as returned by the core.
// The index of each ID in the list becomes its CubismIdHandle value.
func NewCubismIdManager(parameterIdList []string, partIdList []string) *CubismIdManager {
	m := &CubismIdManager{
		parameterIds: make(map[string]CubismIdHandle, len(parameterIdList)),
		partIds:      make(map[string]CubismIdHandle, len(partIdList)),
	}

	for i, id := range parameterIdList {
		m.parameterIds[id] = CubismIdHandle(i)
	}

	for i, id := range partIdList {
		m.partIds[id] = CubismIdHandle(i)
	}

	return m
}

// GetParameterId returns the handle for a parameter ID string.
// Returns InvalidHandle (-1) if the parameter is not found.
func (m *CubismIdManager) GetParameterId(id string) CubismIdHandle {
	if h, ok := m.parameterIds[id]; ok {
		return h
	}
	return InvalidHandle
}

// GetPartId returns the handle for a part ID string.
// Returns InvalidHandle (-1) if the part is not found.
func (m *CubismIdManager) GetPartId(id string) CubismIdHandle {
	if h, ok := m.partIds[id]; ok {
		return h
	}
	return InvalidHandle
}

// GetParameterCount returns the number of registered parameter IDs
func (m *CubismIdManager) GetParameterCount() int {
	return len(m.parameterIds)
}

// GetPartCount returns the number of registered part IDs
func (m *CubismIdManager) GetPartCount() int {
	return len(m.partIds)
}
