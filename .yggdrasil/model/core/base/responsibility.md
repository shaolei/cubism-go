# Core Base — Responsibility

## Identity
The `internal/core/base` package contains all shared logic for FFI function registration, moc3 loading, and data extraction from the native Cubism Core model. It is used by all version-specific implementations.

## Responsibility
- `Funcs` struct: holds all native function pointers (common + version-specific)
- `RegisterCommonFuncs()`: registers all FFI function pointers shared across v5 and v6
- `LoadMoc()`: loads a moc3 file with proper memory alignment (64-byte for SIMD, 16-byte for SSE)
- Data extraction functions: `GetDynamicFlags`, `GetOpacities`, `GetVertexPositions`, `GetDrawables`, `GetParameters`, `GetParameterValue`, `SetParameterValue`, `GetPartIds`, `SetPartOpacity`, `GetSortedDrawableIndices`, `GetCanvasInfo`
- `Update()`: resets dynamic flags and updates the model
- `GetVersion()`: queries and formats the native DLL version

## NOT Responsible For
- Version-specific function registration (delegated to v5/v6 nodes)
- Platform-specific DLL loading (delegated to core router)

## Key Invariants
- MocBuffer MUST be 64-byte aligned for SIMD operations
- ModelBuffer MUST be 16-byte aligned for SSE operations
- `GetDrawables()` is expensive — should only be called once during initial load
- `GetSortedDrawableIndices()` sorts by render orders (v5) or draw orders (v6) depending on which function `CsmGetDrawableSortOrders` points to
- C string conversion MUST use `strings.GoString()`, never `C.GoString()`
