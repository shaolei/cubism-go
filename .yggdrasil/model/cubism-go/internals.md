# cubism-go — Internals

## Logic
### LoadModel Flow
1. Read `model3.json` → parse into `model.ModelJson`
2. Load moc3 file via `core.LoadMoc()` → get `moc.Moc`
3. Get all drawables via `core.GetDrawables()` → map texture indices to absolute paths
4. Build `drawablesMap` (id → Drawable) for O(1) lookup
5. Get sorted indices via `core.GetSortedDrawableIndices()`
6. Conditionally load: physics, pose, display info, expressions, motions, userdata
7. For motions with sound: use `LoadSound` callback or `disabled.LoadSound`

### Update Flow
1. If motionManager exists → `motionManager.Update(delta)` (applies parameter curves)
2. If blinkManager exists → `blinkManager.Update(delta)` (applies eye blink parameter)
3. Call `core.Update(modelPtr)` (native model evaluation)
4. Read dynamic flags → detect what changed (draw order, opacity, vertex positions)
5. Conditionally refresh: sorted indices, opacities, vertex positions

## Constraints
- MotionManager is lazily initialized on first `PlayMotion()` call
- Loop motions: on finish callback, the motion is reset (not closed)
- Non-loop motions: on finish callback, the motion is closed and removed from queue
- Only the last entry in the motion queue is actively evaluated

## Decisions
- Chose lazy MotionManager init over eager init: avoids allocating motion infrastructure when no motions are played
- Chose separate disabled sound package over nil checks: provides a clean Sound interface implementation that no-ops, avoiding nil pointer panics
- Chose map for drawables (O(1) lookup by ID) over slice scan: `GetDrawable()` is called per-frame for hit detection
