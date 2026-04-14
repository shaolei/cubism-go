# Motion — Responsibility

## Identity
The `internal/motion` package manages animation playback for Live2D models, including curve interpolation, fade blending, and a motion queue.

## Responsibility
- **MotionManager**: Manages a queue of `Entry` instances, applies parameter curves each frame
- **Entry**: Tracks playback time for a single motion
- **Motion/Curve/Segment/Point/Meta**: Data types for motion definitions
- **math.go**: Easing functions (`getEasingSine`), segment intersection/interpolation (Linear, Bezier, Stepped, InverseStepped), fade calculation

## NOT Responsible For
- Motion JSON parsing (delegated to `model-json`)
- Sound playback (uses `sound.Sound` interface)

## Key Invariants
- Only the LAST entry in the queue is actively evaluated
- Parameters are saved before evaluation and restored before applying new values (prevents cross-motion interference)
- When a motion finishes: loop motions are reset, non-loop motions are closed
- Fade weights multiply: motion-level fade × per-curve fade
