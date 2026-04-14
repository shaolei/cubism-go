# Motion — Internals

## Logic

### Update Flow
1. If queue is empty → return
2. Update last entry's currentTime by deltaTime
3. If finished → call onFinished callback
4. Load saved parameters (restore to pre-motion state)
5. If entry just started (currentTime == deltaTime) and has sound → play sound
6. Calculate fade weights (motion-level + curve-level)
7. For each curve → for each segment:
   - Check if segment intersects current time
   - Interpolate value at current time
   - Apply to model based on target type:
     - "Parameter": blend with fade weight and set parameter value
     - "PartOpacity": set part opacity directly
     - "Model": TODO (not implemented)
8. Save current parameters for next frame

### Interpolation Algorithms
- **Linear**: Simple linear interpolation between two points
- **Bezier**: Cubic Bezier via De Casteljau's algorithm (3 levels of lerp)
- **Stepped**: Returns start point value (holds until segment end time)
- **InverseStepped**: Returns end point value (holds from segment start time)

### Fade Calculation
- Motion fade: `getEasingSine(t / fadeInTime)` × `getEasingSine((duration - t) / fadeOutTime)`
- Curve fade (overrides motion fade if set): same formula with per-curve times
- `getEasingSine`: `0.5 - 0.5 * cos(value * π)`, clamped to [0, 1]

## Decisions
- Chose De Casteljau over Bernstein polynomial for Bezier: simpler implementation, numerically stable
- Chose last-entry-only evaluation: matches Live2D specification where later motions override earlier ones
- Chose save/restore parameter pattern: prevents parameter drift when motions overlap
