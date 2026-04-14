# Core Base — Internals

## Logic

### LoadMoc Flow
1. Read moc3 file into buffer
2. Allocate aligned buffer (64-byte for SIMD) and copy data
3. Check consistency via `CsmHasMocConsistency()` — MUST return 1
4. Revive moc in place via `CsmReviveMocInPlace()`
5. Get model size via `CsmGetSizeofModel()`
6. Allocate aligned model buffer (16-byte for SSE)
7. Initialize model in place via `CsmInitializeModelInPlace()`

### GetDrawables Flow
1. Get drawable count
2. Parse constant flags (bitwise) and dynamic flags
3. Read texture indices, opacities, vertex counts
4. Iterate drawables to read vertex positions, UVs, indices
5. Read mask counts and masks
6. Read IDs via C string conversion
7. Assemble `drawable.Drawable` structs

### GetSortedDrawableIndices Flow
1. Get drawable count
2. Read raw sort orders (render orders in v5, draw orders in v6)
3. Sort entries by order value
4. Return sorted drawable indices

## Constraints
- All pointer arithmetic uses `unsafe.Pointer` and `unsafe.Slice` — no bounds checking
- Native memory is NOT owned by Go — the model buffer is managed by the native DLL
- `SetParameterValue` writes directly to native memory via pointer arithmetic

## Decisions
- Chose manual memory alignment over `malloc` alignment: gives Go control over the buffer lifecycle, enabling `moc.Moc.Close()` to nil the buffers
- Chose `unsafe.Slice` over reflect.SliceHeader: more idiomatic in modern Go (1.17+), cleaner and safer
- Chose sort-based index calculation for draw order: the native API returns absolute order values, not sorted indices
