# Purego FFI

## What
All native Cubism Core C functions are registered and called using `github.com/ebitengine/purego`, enabling CGO-free dynamic library interop.

## Why
- CGO introduces build complexity, cross-compilation issues, and runtime overhead.
- Purego allows calling C functions from Go without CGO by using platform-specific dynamic linking mechanisms.
- This is critical for a library that needs to ship as a single Go binary with an optional native DLL.

## Rules
1. ALL native function registrations MUST use `purego.RegisterLibFunc()`.
2. Native function pointers MUST be stored in `base.Funcs` struct.
3. New native functions MUST be registered in `base.RegisterCommonFuncs()` (shared) or in the version-specific `NewCore()` (version-specific).
4. Pointer arithmetic for reading native data MUST use `unsafe.Pointer` and `unsafe.Slice`.
5. Memory alignment requirements:
   - MocBuffer: 64-byte alignment (SIMD)
   - ModelBuffer: 16-byte alignment (SSE)
6. C strings from native code MUST be converted via `strings.GoString()` — NOT `C.GoString`.

## Scope
- `internal/core/base/core.go` (function registration and data reading)
- `internal/core/core_5_0_0/core.go` (v5-specific registration)
- `internal/core/core_6_0_1/core.go` (v6-specific registration)
- `internal/core/minimum/core.go` (version query)
- `internal/strings/strings.go` (C string conversion)
