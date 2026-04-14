# Core Router — Interface

## Core Interface
```go
type Core interface {
    LoadMoc(path string) (moc.Moc, error)
    GetVersion() string
    GetDynamicFlags(uintptr) []drawable.DynamicFlag
    GetOpacities(uintptr) []float32
    GetVertexPositions(uintptr) [][]drawable.Vector2
    GetDrawables(uintptr) []drawable.Drawable
    GetParameters(uintptr) []parameter.Parameter
    GetParameterValue(uintptr, string) float32
    SetParameterValue(uintptr, string, float32)
    GetPartIds(uintptr) []string
    SetPartOpacity(uintptr, string, float32)
    GetSortedDrawableIndices(uintptr) []int
    GetCanvasInfo(uintptr) (drawable.Vector2, drawable.Vector2, float32)
    Update(uintptr)
}
```

## Factory
```go
func NewCore(lib string) (Core, error)
```

## Platform Functions
```go
// Windows only
func CloseLibrary(name string) error
```

## Failure Modes
- `NewCore`: DLL load failure, version parse failure, unsupported version
- All Core methods accept `modelPtr uintptr` — passing 0 or invalid pointer causes undefined behavior (native crash)
