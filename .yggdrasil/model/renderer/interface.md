# Renderer — Interface

## Renderer
```go
func NewRenderer(model *cubism.Model) (*Renderer, error)
func (r *Renderer) Update() error
func (r *Renderer) Draw(screen *ebiten.Image, opts ...func(*DrawOption))
func (r *Renderer) GetModel() *cubism.Model
func (r *Renderer) IsHit(x, y int, id string) (bool, error)
```

## Draw Options
```go
func WithHidden() func(*DrawOption)
func WithScale(scale float64) func(*DrawOption)
func WithPosition(x, y float64) func(*DrawOption)
func WithBackground(c color.Color) func(*DrawOption)
```

## Failure Modes
- `NewRenderer`: Texture file not found, shader compilation error
- `IsHit`: Drawable ID not found, coordinates out of bounds → returns false, nil
