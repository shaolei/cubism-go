# Version Routing

## What
The `core.NewCore()` function loads the native Cubism Core DLL, queries its version via `minimum.NewCore()`, then routes to the correct version-specific implementation (`core_5_0_0` or `core_6_0_1`) based on the major version number.

## Why
Live2D Cubism SDK v5 and v6 have incompatible native APIs — notably the sorting order function differs:
- v5: `csmGetDrawableRenderOrders`
- v6: `csmGetDrawableDrawOrders`

A runtime version check ensures the correct FFI bindings are used without requiring separate builds.

## Rules
1. Every new Cubism Core major version MUST have a corresponding implementation package under `internal/core/core_X_Y_Z/`.
2. Each implementation MUST embed `base.Funcs` and call `base.RegisterCommonFuncs()` before registering version-specific functions.
3. The version-specific `CsmGetDrawableSortOrders` field in `base.Funcs` MUST be mapped to the correct native function name for that version.
4. `parseMajorVersion()` in `internal/core/core.go` MUST be updated if version string formats change.
5. `NewCore()` switch statement MUST include a default case that returns an error for unsupported versions.

## Scope
- `internal/core/core.go` (routing logic)
- `internal/core/core_5_0_0/` (v5 implementation)
- `internal/core/core_6_0_1/` (v6 implementation)
- `internal/core/minimum/` (version detection)
