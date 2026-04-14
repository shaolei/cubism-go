# Core Data Types — Interface

## drawable.ConstantFlag
```go
type ConstantFlag struct {
    BlendAdditive       bool
    BlendMultiplicative bool
    IsDoubleSided       bool
    IsInvertedMask      bool
}
func ParseConstantFlag(flag uint8) ConstantFlag
```

## drawable.DynamicFlag
```go
type DynamicFlag struct {
    IsVisible                bool
    VisibilityDidChange      bool
    OpacityDidChange         bool
    DrawOrderDidChange       bool
    RenderOrderDidChange     bool
    VertexPositionsDidChange bool
    BlendColorDidChange      bool
}
func ParseDynamicFlag(flag uint8) DynamicFlag
```

## drawable.Drawable
```go
type Drawable struct {
    Id              string
    Texture         int32
    VertexPositions []Vector2
    VertexUvs       []Vector2
    VertexIndices   []uint16
    ConstantFlag    ConstantFlag
    DynamicFlag     DynamicFlag
    Opacity         float32
    Masks           []int32
}
```

## drawable.Vector2
```go
type Vector2 struct { X, Y float32 }
```

## moc.Moc
```go
type Moc struct {
    MocPtr      uintptr
    MocBuffer   []byte
    ModelPtr    uintptr
    ModelBuffer []byte
}
func (m *Moc) Close()
```

## parameter.Parameter
```go
type Parameter struct {
    Id      string
    Minimum float32
    Maximum float32
    Default float32
    Current float32
}
```
