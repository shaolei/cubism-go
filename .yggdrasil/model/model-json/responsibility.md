# Model JSON — Responsibility

## Identity
The `internal/model` package defines JSON struct types for all Live2D Cubism model configuration files and provides the `ToMotion()` converter.

## Responsibility
- **ModelJson**: `model3.json` — the root model descriptor
- **ExpJson**: `exp3.json` — expression definitions
- **MotionJson**: `motion3.json` — motion curve data with `ToMotion()` converter
- **PhysicsJson**: `physics3.json` — physics simulation settings
- **PoseJson**: `pose3.json` — pose group definitions
- **CdiJson**: `cdi3.json` — parameter/part display names
- **UserDataJson**: `userdata3.json` — custom user data
- **Group/HitArea**: Shared types for model groups and hit areas

## NOT Responsible For
- File I/O or model loading (handled by root `cubism.go`)
- Motion playback (handled by `motion` package)

## Key Invariants
- `ToMotion()` converts the flat `Segments []float64` array into typed `Segment` structs by iterating with type-dependent stride (Linear: +3, Bezier: +7, Stepped: +3, InverseStepped: +3)
- Null fade times in JSON are converted to -1.0 in the domain model (signal to use motion-level fade)
