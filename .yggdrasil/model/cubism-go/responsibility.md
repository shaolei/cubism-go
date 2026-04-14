# cubism-go — Root Module

## Identity
The root package (`github.com/shaolei/cubism-go`) provides the public API for loading, animating, and rendering Live2D Cubism models. It is the single entry point for consumers of this library.

## Responsibility
- **Cubism struct**: Initializes the native Core DLL via `core.NewCore()`.
- **Model loading**: `Cubism.LoadModel(path)` reads a `model3.json`, loads all referenced resources (moc3, textures, physics, pose, expressions, motions, userdata), and returns a fully populated `Model`.
- **Model struct**: Holds all model data (drawables, parameters, motions, textures) and provides methods for parameter manipulation, motion playback, blink control, and per-frame updates.
- **Drawable struct**: Public-facing representation of a drawable with vertex positions, UVs, indices, flags, opacity, and texture path.

## NOT Responsible For
- Direct native FFI calls (delegated to `core` hierarchy)
- JSON struct definitions (delegated to `internal/model`)
- Rendering (delegated to `renderer`)
- Audio playback (delegated to `sound`)

## Key Invariants
- `Cubism.LoadModel()` MUST be called after `NewCubism()` successfully loads the DLL.
- `Model.Update(delta)` MUST be called each frame before reading drawable state.
- `Model.Close()` MUST be called to release native moc3 resources when done.
- If `Cubism.LoadSound` is nil, motion sounds use the `disabled` implementation (no-op).
