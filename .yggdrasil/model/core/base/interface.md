# Core Base — Interface

## Funcs Struct
```go
type Funcs struct {
    CsmGetVersion                func() uint32
    CsmReviveMocInPlace          func(uintptr, uint) uintptr
    CsmGetSizeofModel            func(uintptr) uint
    CsmInitializeModelInPlace    func(uintptr, uintptr, uint) uintptr
    CsmHasMocConsistency         func(uintptr, uint) int
    CsmUpdateModel               func(uintptr)
    CsmReadCanvasInfo            func(uintptr, uintptr, uintptr, uintptr)
    CsmResetDrawableDynamicFlags func(uintptr)
    CsmGetParameterCount         func(uintptr) int
    // ... (see source for full list)
    CsmGetDrawableSortOrders     func(uintptr) uintptr  // Version-specific: render orders or draw orders
}
```

## Functions
```go
func RegisterCommonFuncs(f *Funcs, lib uintptr)
func LoadMoc(f *Funcs, path string) (moc.Moc, error)
func GetVersion(f *Funcs) string
func GetDynamicFlags(f *Funcs, modelPtr uintptr) []drawable.DynamicFlag
func GetOpacities(f *Funcs, modelPtr uintptr) []float32
func GetVertexPositions(f *Funcs, modelPtr uintptr) [][]drawable.Vector2
func GetDrawables(f *Funcs, modelPtr uintptr) []drawable.Drawable
func GetParameters(f *Funcs, modelPtr uintptr) []parameter.Parameter
func GetParameterValue(f *Funcs, modelPtr uintptr, id string) float32
func SetParameterValue(f *Funcs, modelPtr uintptr, id string, value float32)
func GetPartIds(f *Funcs, modelPtr uintptr) []string
func SetPartOpacity(f *Funcs, modelPtr uintptr, id string, value float32)
func GetSortedDrawableIndices(f *Funcs, modelPtr uintptr) []int
func GetCanvasInfo(f *Funcs, modelPtr uintptr) (drawable.Vector2, drawable.Vector2, float32)
func Update(f *Funcs, modelPtr uintptr)
```

## Failure Modes
- `LoadMoc`: File not found, consistency check failed, revive failed, model size is 0, initialization failed
- `GetParameterValue`: Parameter ID not found → returns 0
- `SetParameterValue`: Parameter ID not found → silently ignored
