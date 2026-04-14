# Blink — Responsibility

## Identity
The `internal/blink` package implements an automatic eye-blink state machine that periodically drives eye-openness parameters on a Live2D model.

## Responsibility
- **BlinkManager**: State machine with 5 states (First, Interval, Closing, Closed, Opening)
- Cycles through blink states with configurable timing
- Sets eye parameter values: 1.0 = fully open, 0.0 = fully closed
- Randomized blink interval based on `interval` constant

## NOT Responsible For
- Determining which parameters are eye-blink parameters (provided by caller via `ids`)
- Parameter value clamping (Core handles this)

## State Machine
```
First → Interval → Closing → Closed → Opening → Interval → ...
```
- **First**: Initial state, transitions immediately to Interval
- **Interval**: Eyes open (value=1.0), waits for next blink time
- **Closing**: Eyes closing (value=1.0→0.0), duration=closing (0.1s)
- **Closed**: Eyes closed (value=0.0), duration=closing (0.1s)
- **Opening**: Eyes opening (value=0.0→1.0), duration=opening (0.15s)

## Key Invariants
- Default interval: 4.0 seconds, closing: 0.1s, opening: 0.15s
- Next blink time = currentTime + random × (2×interval - 1)
- All specified IDs are set to the same value each frame
