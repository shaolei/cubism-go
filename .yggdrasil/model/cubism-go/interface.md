# cubism-go — Interface

## Public Types

### Cubism
```go
type Cubism struct { ... }

func NewCubism(lib string) (Cubism, error)
func (c *Cubism) LoadModel(path string) (*Model, error)
```
- `lib`: Path to the native Cubism Core DLL/SO.
- Returns error if DLL cannot be loaded or model files are missing/invalid.

### Model
```go
type Model struct { ... }

func (m *Model) GetVersion() int
func (m *Model) GetCore() core.Core
func (m *Model) GetMoc() moc.Moc
func (m *Model) GetOpacity() float32
func (m *Model) GetTextures() []string
func (m *Model) GetSortedIndices() []int
func (m *Model) GetDrawables() []Drawable
func (m *Model) GetDrawable(id string) (Drawable, error)
func (m *Model) GetHitAreas() []model.HitArea
func (m *Model) GetParameters() []parameter.Parameter
func (m *Model) GetParameterValue(id string) float32
func (m *Model) SetParameterValue(id string, value float32)
func (m *Model) GetMotionGroupNames() []string
func (m *Model) GetMotions(groupName string) []motion.Motion
func (m *Model) PlayMotion(groupName string, index int, loop bool) int
func (m *Model) StopMotion(id int)
func (m *Model) EnableAutoBlink()
func (m *Model) DisableAutoBlink()
func (m *Model) Update(delta float64)
func (m *Model) Close()
```

### Drawable
```go
type Drawable struct {
    Id              string
    Texture         string
    VertexPositions []drawable.Vector2
    VertexUvs       []drawable.Vector2
    VertexIndices   []uint16
    ConstantFlag    drawable.ConstantFlag
    DynamicFlag     drawable.DynamicFlag
    Opacity         float32
    Masks           []int32
}
```

## Failure Modes
- `NewCubism`: DLL not found, DLL incompatible, version unsupported
- `LoadModel`: File not found, JSON parse error, moc3 consistency check failed, moc3 initialization failed
- `GetDrawable`: Drawable ID not found → returns error
- `PlayMotion`: Group name not found → panic (index out of range)
