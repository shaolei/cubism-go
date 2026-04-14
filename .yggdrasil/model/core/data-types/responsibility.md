# Core Data Types — Responsibility

## Identity
The `drawable`, `moc`, and `parameter` packages define the core data structures used across the Cubism Core implementation.

## Responsibility
- **drawable.ConstantFlag**: Parsed bit flags for blend mode, double-sided rendering, inverted mask
- **drawable.DynamicFlag**: Parsed bit flags for visibility and change detection
- **drawable.Drawable**: Composite struct holding all drawable data from the native API
- **drawable.Vector2**: 2D vector (X, Y float32)
- **drawable.ParseConstantFlag / ParseDynamicFlag**: Bitwise flag parsing from uint8
- **moc.Moc**: Holds native pointers and buffers for a loaded moc3 model with close lifecycle
- **parameter.Parameter**: Holds parameter ID, min, max, default, and current values

## NOT Responsible For
- Any FFI logic or data extraction (delegated to `core/base`)

## Key Invariants
- `ConstantFlag` bits: bit0=BlendAdditive, bit1=BlendMultiplicative, bit2=IsDoubleSided, bit3=IsInvertedMask
- `DynamicFlag` bits: bit0=IsVisible, bit1=VisibilityDidChange, bit2=OpacityDidChange, bit3=DrawOrderDidChange, bit4=RenderOrderDidChange, bit5=VertexPositionsDidChange, bit6=BlendColorDidChange
- `moc.Moc.Close()` is idempotent — safe to call multiple times
