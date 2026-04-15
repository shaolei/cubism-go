package base

import (
	"fmt"
	"os"
	"sort"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/shaolei/cubism-go/internal/core/drawable"
	"github.com/shaolei/cubism-go/internal/core/moc"
	"github.com/shaolei/cubism-go/internal/core/parameter"
	"github.com/shaolei/cubism-go/internal/strings"
	"github.com/shaolei/cubism-go/internal/utils"
)

// Funcs holds all the common Cubism Core function pointers shared across versions.
type Funcs struct {
	// Version
	CsmGetVersion func() uint32
	// Moc
	CsmReviveMocInPlace       func(uintptr, uint) uintptr
	CsmGetSizeofModel         func(uintptr) uint
	CsmInitializeModelInPlace func(uintptr, uintptr, uint) uintptr
	CsmHasMocConsistency      func(uintptr, uint) int
	// Model
	CsmUpdateModel               func(uintptr)
	CsmReadCanvasInfo            func(uintptr, uintptr, uintptr, uintptr)
	CsmResetDrawableDynamicFlags func(uintptr)
	// Parameters
	CsmGetParameterCount         func(uintptr) int
	CsmGetParameterIds           func(uintptr) uintptr
	CsmGetParameterTypes         func(uintptr) uintptr
	CsmGetParameterMinimumValues func(uintptr) uintptr
	CsmGetParameterMaximumValues func(uintptr) uintptr
	CsmGetParameterDefaultValues func(uintptr) uintptr
	CsmGetParameterValues        func(uintptr) uintptr
	// Parts
	CsmGetPartCount     func(uintptr) int
	CsmGetPartIds       func(uintptr) uintptr
	CsmGetPartOpacities func(uintptr) uintptr
	// Drawables (common)
	CsmGetDrawableCount           func(uintptr) int
	CsmGetDrawableIds             func(uintptr) uintptr
	CsmGetDrawableConstantFlags   func(uintptr) uintptr
	CsmGetDrawableDynamicFlags    func(uintptr) uintptr
	CsmGetDrawableTextureIndices  func(uintptr) uintptr
	CsmGetDrawableOpacities       func(uintptr) uintptr
	CsmGetDrawableMaskCounts      func(uintptr) uintptr
	CsmGetDrawableMasks           func(uintptr) uintptr
	CsmGetDrawableVertexCounts    func(uintptr) uintptr
	CsmGetDrawableVertexPositions func(uintptr) uintptr
	CsmGetDrawableVertexUvs       func(uintptr) uintptr
	CsmGetDrawableIndexCounts     func(uintptr) uintptr
	CsmGetDrawableIndices         func(uintptr) uintptr
	// Version-specific: function to get sorting orders (render orders or draw orders)
	CsmGetDrawableSortOrders func(uintptr) uintptr
}

// RegisterCommonFuncs registers all function pointers that are shared across versions.
func RegisterCommonFuncs(f *Funcs, lib uintptr) {
	purego.RegisterLibFunc(&f.CsmGetVersion, lib, "csmGetVersion")
	purego.RegisterLibFunc(&f.CsmReviveMocInPlace, lib, "csmReviveMocInPlace")
	purego.RegisterLibFunc(&f.CsmGetSizeofModel, lib, "csmGetSizeofModel")
	purego.RegisterLibFunc(&f.CsmInitializeModelInPlace, lib, "csmInitializeModelInPlace")
	purego.RegisterLibFunc(&f.CsmHasMocConsistency, lib, "csmHasMocConsistency")
	purego.RegisterLibFunc(&f.CsmUpdateModel, lib, "csmUpdateModel")
	purego.RegisterLibFunc(&f.CsmReadCanvasInfo, lib, "csmReadCanvasInfo")
	purego.RegisterLibFunc(&f.CsmResetDrawableDynamicFlags, lib, "csmResetDrawableDynamicFlags")
	purego.RegisterLibFunc(&f.CsmGetParameterCount, lib, "csmGetParameterCount")
	purego.RegisterLibFunc(&f.CsmGetParameterIds, lib, "csmGetParameterIds")
	purego.RegisterLibFunc(&f.CsmGetParameterTypes, lib, "csmGetParameterTypes")
	purego.RegisterLibFunc(&f.CsmGetParameterMinimumValues, lib, "csmGetParameterMinimumValues")
	purego.RegisterLibFunc(&f.CsmGetParameterMaximumValues, lib, "csmGetParameterMaximumValues")
	purego.RegisterLibFunc(&f.CsmGetParameterDefaultValues, lib, "csmGetParameterDefaultValues")
	purego.RegisterLibFunc(&f.CsmGetParameterValues, lib, "csmGetParameterValues")
	purego.RegisterLibFunc(&f.CsmGetPartCount, lib, "csmGetPartCount")
	purego.RegisterLibFunc(&f.CsmGetPartIds, lib, "csmGetPartIds")
	purego.RegisterLibFunc(&f.CsmGetPartOpacities, lib, "csmGetPartOpacities")
	purego.RegisterLibFunc(&f.CsmGetDrawableCount, lib, "csmGetDrawableCount")
	purego.RegisterLibFunc(&f.CsmGetDrawableIds, lib, "csmGetDrawableIds")
	purego.RegisterLibFunc(&f.CsmGetDrawableConstantFlags, lib, "csmGetDrawableConstantFlags")
	purego.RegisterLibFunc(&f.CsmGetDrawableDynamicFlags, lib, "csmGetDrawableDynamicFlags")
	purego.RegisterLibFunc(&f.CsmGetDrawableTextureIndices, lib, "csmGetDrawableTextureIndices")
	purego.RegisterLibFunc(&f.CsmGetDrawableOpacities, lib, "csmGetDrawableOpacities")
	purego.RegisterLibFunc(&f.CsmGetDrawableMaskCounts, lib, "csmGetDrawableMaskCounts")
	purego.RegisterLibFunc(&f.CsmGetDrawableMasks, lib, "csmGetDrawableMasks")
	purego.RegisterLibFunc(&f.CsmGetDrawableVertexCounts, lib, "csmGetDrawableVertexCounts")
	purego.RegisterLibFunc(&f.CsmGetDrawableVertexPositions, lib, "csmGetDrawableVertexPositions")
	purego.RegisterLibFunc(&f.CsmGetDrawableVertexUvs, lib, "csmGetDrawableVertexUvs")
	purego.RegisterLibFunc(&f.CsmGetDrawableIndexCounts, lib, "csmGetDrawableIndexCounts")
	purego.RegisterLibFunc(&f.CsmGetDrawableIndices, lib, "csmGetDrawableIndices")
}

// LoadMoc loads a moc3 file and returns moc.Moc
func LoadMoc(f *Funcs, path string) (m moc.Moc, err error) {
	m.MocBuffer, err = os.ReadFile(path)
	if err != nil {
		return
	}
	// Ensure MocBuffer is 64-byte aligned for SIMD operations
	// by padding and aligning the start address
	alignedBuf := make([]byte, len(m.MocBuffer)+64)
	offset := 0
	addr := uintptr(unsafe.Pointer(&alignedBuf[0]))
	if rem := addr % 64; rem != 0 {
		offset = int(64 - rem)
	}
	copy(alignedBuf[offset:], m.MocBuffer)
	m.MocBuffer = alignedBuf[offset : offset+len(m.MocBuffer)]

	consistency := f.CsmHasMocConsistency(uintptr(unsafe.Pointer(&m.MocBuffer[0])), uint(len(m.MocBuffer)))
	if consistency != 1 {
		err = fmt.Errorf("moc3 is not consistent")
		return
	}
	m.MocPtr = f.CsmReviveMocInPlace(uintptr(unsafe.Pointer(&m.MocBuffer[0])), uint(len(m.MocBuffer)))
	if m.MocPtr == 0 {
		err = fmt.Errorf("failed to revive moc3")
		return
	}
	size := f.CsmGetSizeofModel(m.MocPtr)
	if size == 0 {
		err = fmt.Errorf("failed to get size of model")
		return
	}
	// ModelBuffer also needs alignment (16-byte minimum for SSE)
	modelBuf := make([]byte, size+16)
	modelOffset := 0
	modelAddr := uintptr(unsafe.Pointer(&modelBuf[0]))
	if rem := modelAddr % 16; rem != 0 {
		modelOffset = int(16 - rem)
	}
	m.ModelBuffer = modelBuf[modelOffset : modelOffset+int(size)]

	m.ModelPtr = f.CsmInitializeModelInPlace(m.MocPtr, uintptr(unsafe.Pointer(&m.ModelBuffer[0])), size)
	if m.ModelPtr == 0 {
		err = fmt.Errorf("failed to initialize model")
		return
	}
	return
}

// GetVersion returns the Core version as a string
func GetVersion(f *Funcs) string {
	raw := f.CsmGetVersion()
	return utils.ParseVersion(raw)
}

// GetDynamicFlags returns the dynamic flags for all drawables
func GetDynamicFlags(f *Funcs, modelPtr uintptr) (rs []drawable.DynamicFlag) {
	count := f.CsmGetDrawableCount(modelPtr)
	raw := unsafe.Slice((*uint8)(unsafe.Pointer(f.CsmGetDrawableDynamicFlags(modelPtr))), count)
	for _, flag := range raw {
		rs = append(rs, drawable.ParseDynamicFlag(flag))
	}
	return
}

// GetOpacities returns the opacities for all drawables
func GetOpacities(f *Funcs, modelPtr uintptr) (rs []float32) {
	count := f.CsmGetDrawableCount(modelPtr)
	rs = unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetDrawableOpacities(modelPtr))), count)
	return
}

// GetVertexPositions returns the vertex positions for all drawables
func GetVertexPositions(f *Funcs, modelPtr uintptr) (vps [][]drawable.Vector2) {
	count := f.CsmGetDrawableCount(modelPtr)
	vertexCounts := unsafe.Slice((*int32)(unsafe.Pointer(f.CsmGetDrawableVertexCounts(modelPtr))), count)
	posPtr := f.CsmGetDrawableVertexPositions(modelPtr)
	for i := 0; i < count; i++ {
		vertexCount := vertexCounts[i]
		positions := unsafe.Slice(*(**drawable.Vector2)(unsafe.Pointer(posPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(vertexCount))
		vps = append(vps, positions)
	}
	return
}

// GetDrawables returns all drawable information.
// Since all the information is gathered, the cost is high. It is expected to be called only once initially.
func GetDrawables(f *Funcs, modelPtr uintptr) (ds []drawable.Drawable) {
	count := f.CsmGetDrawableCount(modelPtr)

	constantFlags := make([]drawable.ConstantFlag, 0, count)
	raw := unsafe.Slice((*uint8)(unsafe.Pointer(f.CsmGetDrawableConstantFlags(modelPtr))), count)
	for _, flag := range raw {
		constantFlags = append(constantFlags, drawable.ParseConstantFlag(flag))
	}

	dynamicFlags := GetDynamicFlags(f, modelPtr)

	textureIndices := unsafe.Slice((*int32)(unsafe.Pointer(f.CsmGetDrawableTextureIndices(modelPtr))), count)

	opacities := GetOpacities(f, modelPtr)

	vertexCounts := unsafe.Slice((*int32)(unsafe.Pointer(f.CsmGetDrawableVertexCounts(modelPtr))), count)

	vertexPositions := make([][]drawable.Vector2, 0, count)
	vertexUvs := make([][]drawable.Vector2, 0, count)
	posPtr := f.CsmGetDrawableVertexPositions(modelPtr)
	uvPtr := f.CsmGetDrawableVertexUvs(modelPtr)
	for i := 0; i < count; i++ {
		vertexCount := vertexCounts[i]
		positions := unsafe.Slice(*(**drawable.Vector2)(unsafe.Pointer(posPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(vertexCount))
		vertexPositions = append(vertexPositions, positions)
		uvs := unsafe.Slice(*(**drawable.Vector2)(unsafe.Pointer(uvPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(vertexCount))
		vertexUvs = append(vertexUvs, uvs)
	}

	indexCounts := unsafe.Slice((*int32)(unsafe.Pointer(f.CsmGetDrawableIndexCounts(modelPtr))), count)
	indices := make([][]uint16, 0, count)
	indicesPtr := f.CsmGetDrawableIndices(modelPtr)
	for i := 0; i < count; i++ {
		indexCount := indexCounts[i]
		indices = append(indices, unsafe.Slice(*(**uint16)(unsafe.Pointer(indicesPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(indexCount)))
	}

	maskCounts := unsafe.Slice((*int32)(unsafe.Pointer(f.CsmGetDrawableMaskCounts(modelPtr))), count)
	masks := make([][]int32, 0, count)
	maskPtr := f.CsmGetDrawableMasks(modelPtr)
	for i := 0; i < count; i++ {
		maskCount := maskCounts[i]
		masks = append(masks, unsafe.Slice(*(**int32)(unsafe.Pointer(maskPtr + uintptr(i)*unsafe.Sizeof(uintptr(0)))), int(maskCount)))
	}

	idsPtr := f.CsmGetDrawableIds(modelPtr)
	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		ids = append(ids, strings.GoString(uintptr(unsafe.Pointer(ptr))))
	}

	for i := 0; i < count; i++ {
		d := drawable.Drawable{
			Id:              ids[i],
			Texture:         textureIndices[i],
			VertexPositions: vertexPositions[i],
			VertexUvs:       vertexUvs[i],
			VertexIndices:   indices[i],
			ConstantFlag:    constantFlags[i],
			DynamicFlag:     dynamicFlags[i],
			Opacity:         opacities[i],
			Masks:           masks[i],
		}
		ds = append(ds, d)
	}
	return
}

// GetParameters returns all parameters
func GetParameters(f *Funcs, modelPtr uintptr) (parameters []parameter.Parameter) {
	count := f.CsmGetParameterCount(modelPtr)
	idsPtr := f.CsmGetParameterIds(modelPtr)
	mins := unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterMinimumValues(modelPtr))), count)
	maxs := unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterMaximumValues(modelPtr))), count)
	defs := unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterDefaultValues(modelPtr))), count)
	vals := unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterValues(modelPtr))), count)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		p := parameter.Parameter{
			Id:      strings.GoString(uintptr(unsafe.Pointer(ptr))),
			Minimum: mins[i],
			Maximum: maxs[i],
			Default: defs[i],
			Current: vals[i],
		}
		parameters = append(parameters, p)
	}
	return
}

// GetParameterValue returns the value of a parameter by ID (string lookup, O(n))
// Prefer GetParameterValueByIndex for hot-path code.
func GetParameterValue(f *Funcs, modelPtr uintptr, id string) float32 {
	count := f.CsmGetParameterCount(modelPtr)
	idsPtr := f.CsmGetParameterIds(modelPtr)
	valPtr := f.CsmGetParameterValues(modelPtr)
	vals := unsafe.Slice((*float32)(unsafe.Pointer(valPtr)), count)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		if strings.GoString(uintptr(unsafe.Pointer(ptr))) == id {
			return vals[i]
		}
	}
	return 0
}

// GetParameterValueByIndex returns the value of a parameter by its array index (O(1))
// The index should be obtained from CubismIdManager.GetParameterId().
func GetParameterValueByIndex(f *Funcs, modelPtr uintptr, index int) float32 {
	count := f.CsmGetParameterCount(modelPtr)
	if index < 0 || index >= count {
		return 0
	}
	valPtr := f.CsmGetParameterValues(modelPtr)
	vals := unsafe.Slice((*float32)(unsafe.Pointer(valPtr)), count)
	return vals[index]
}

// SetParameterValue sets the value of a parameter by ID (string lookup, O(n))
// Prefer SetParameterValueByIndex for hot-path code.
func SetParameterValue(f *Funcs, modelPtr uintptr, id string, value float32) {
	count := f.CsmGetParameterCount(modelPtr)
	idsPtr := f.CsmGetParameterIds(modelPtr)
	valPtr := f.CsmGetParameterValues(modelPtr)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		if strings.GoString(uintptr(unsafe.Pointer(ptr))) == id {
			*(*float32)(unsafe.Pointer(valPtr + uintptr(i)*unsafe.Sizeof(float32(0)))) = value
			return
		}
	}
}

// SetParameterValueByIndex sets the value of a parameter by its array index (O(1))
// The index should be obtained from CubismIdManager.GetParameterId().
func SetParameterValueByIndex(f *Funcs, modelPtr uintptr, index int, value float32) {
	count := f.CsmGetParameterCount(modelPtr)
	if index < 0 || index >= count {
		return
	}
	valPtr := f.CsmGetParameterValues(modelPtr)
	*(*float32)(unsafe.Pointer(valPtr + uintptr(index)*unsafe.Sizeof(float32(0)))) = value
}

// GetParameterIds returns all parameter ID strings in their array order.
// Used by CubismIdManager to build the ID→index mapping.
func GetParameterIds(f *Funcs, modelPtr uintptr) []string {
	count := f.CsmGetParameterCount(modelPtr)
	idsPtr := f.CsmGetParameterIds(modelPtr)
	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		ids = append(ids, strings.GoString(uintptr(unsafe.Pointer(ptr))))
	}
	return ids
}

// GetPartOpacities returns the opacities for all parts
func GetPartOpacities(f *Funcs, modelPtr uintptr) (rs []float32) {
	count := f.CsmGetPartCount(modelPtr)
	rs = unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetPartOpacities(modelPtr))), count)
	return
}

// GetPartIds returns the part IDs
func GetPartIds(f *Funcs, modelPtr uintptr) (ids []string) {
	count := f.CsmGetPartCount(modelPtr)
	idsPtr := f.CsmGetPartIds(modelPtr)
	for i := 0; i < count; i++ {
		ptr := *(**byte)(unsafe.Pointer(idsPtr + uintptr(i)*unsafe.Sizeof(uintptr(0))))
		ids = append(ids, strings.GoString(uintptr(unsafe.Pointer(ptr))))
	}
	return
}

// SetPartOpacity sets the opacity of a part by ID (string lookup)
// Prefer SetPartOpacityByIndex for hot-path code.
func SetPartOpacity(f *Funcs, modelPtr uintptr, id string, value float32) {
	ids := GetPartIds(f, modelPtr)
	ptr := f.CsmGetPartOpacities(modelPtr)
	for i, _id := range ids {
		if _id == id {
			*(*float32)(unsafe.Pointer(ptr + uintptr(i)*unsafe.Sizeof(float32(0)))) = value
			return
		}
	}
}

// SetPartOpacityByIndex sets the opacity of a part by its array index (O(1))
// The index should be obtained from CubismIdManager.GetPartId().
func SetPartOpacityByIndex(f *Funcs, modelPtr uintptr, index int, value float32) {
	count := f.CsmGetPartCount(modelPtr)
	if index < 0 || index >= count {
		return
	}
	ptr := f.CsmGetPartOpacities(modelPtr)
	*(*float32)(unsafe.Pointer(ptr + uintptr(index)*unsafe.Sizeof(float32(0)))) = value
}

// GetPartOpacityByIndex returns the opacity of a part by its array index (O(1))
// The index should be obtained from CubismIdManager.GetPartId().
func GetPartOpacityByIndex(f *Funcs, modelPtr uintptr, index int) float32 {
	count := f.CsmGetPartCount(modelPtr)
	if index < 0 || index >= count {
		return 0
	}
	opacities := unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetPartOpacities(modelPtr))), count)
	return opacities[index]
}

// GetSortedDrawableIndices returns the drawing order indices sorted by sort orders.
// Sort orders (render orders in v5, draw orders in v6) are absolute values,
// so we sort by order value and return the sorted drawable indices.
func GetSortedDrawableIndices(f *Funcs, modelPtr uintptr) (rs []int) {
	count := f.CsmGetDrawableCount(modelPtr)
	ptr := f.CsmGetDrawableSortOrders(modelPtr)
	rawOrders := unsafe.Slice((*int32)(unsafe.Pointer(ptr)), count)

	type orderEntry struct {
		index int
		order int32
	}
	entries := make([]orderEntry, count)
	for i := 0; i < count; i++ {
		entries[i] = orderEntry{index: i, order: rawOrders[i]}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].order < entries[j].order
	})

	rs = make([]int, count)
	for i, e := range entries {
		rs[i] = e.index
	}
	return
}

// GetParameterCount returns the number of parameters in the model.
func GetParameterCount(f *Funcs, modelPtr uintptr) int {
	return f.CsmGetParameterCount(modelPtr)
}

// GetParameterValues returns all parameter values as a direct float32 slice.
// This is the efficient path used by the physics engine and other subsystems
// that need bulk access to parameter values.
func GetParameterValues(f *Funcs, modelPtr uintptr) []float32 {
	count := f.CsmGetParameterCount(modelPtr)
	return unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterValues(modelPtr))), count)
}

// GetParameterMinimumValues returns all parameter minimum values as a direct float32 slice.
func GetParameterMinimumValues(f *Funcs, modelPtr uintptr) []float32 {
	count := f.CsmGetParameterCount(modelPtr)
	return unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterMinimumValues(modelPtr))), count)
}

// GetParameterMaximumValues returns all parameter maximum values as a direct float32 slice.
func GetParameterMaximumValues(f *Funcs, modelPtr uintptr) []float32 {
	count := f.CsmGetParameterCount(modelPtr)
	return unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterMaximumValues(modelPtr))), count)
}

// GetParameterDefaultValues returns all parameter default values as a direct float32 slice.
func GetParameterDefaultValues(f *Funcs, modelPtr uintptr) []float32 {
	count := f.CsmGetParameterCount(modelPtr)
	return unsafe.Slice((*float32)(unsafe.Pointer(f.CsmGetParameterDefaultValues(modelPtr))), count)
}

// GetCanvasInfo returns the canvas size, origin, and pixels per unit
func GetCanvasInfo(f *Funcs, modelPtr uintptr) (size drawable.Vector2, origin drawable.Vector2, pixelsPerUnit float32) {
	f.CsmReadCanvasInfo(modelPtr, uintptr(unsafe.Pointer(&size)), uintptr(unsafe.Pointer(&origin)), uintptr(unsafe.Pointer(&pixelsPerUnit)))
	return
}

// Update updates the model
func Update(f *Funcs, modelPtr uintptr) {
	f.CsmResetDrawableDynamicFlags(modelPtr)
	f.CsmUpdateModel(modelPtr)
}
