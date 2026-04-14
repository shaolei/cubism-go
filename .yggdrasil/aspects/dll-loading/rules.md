# DLL Loading

## What
The Cubism Core native library is loaded via platform-specific `openLibrary()` functions. On Windows, `golang.org/x/sys/windows` is used with a global DLL cache protected by a mutex. On Unix, `purego.Dlopen` is used.

## Why
- Windows requires explicit DLL management (loading, caching, releasing) for proper resource lifecycle.
- Unix platforms use `purego.Dlopen` which handles library loading differently.
- Caching prevents loading the same DLL multiple times, which would waste memory and could cause issues.

## Rules
1. `openLibrary()` MUST be defined in build-tagged files (`open_windows.go`, `open_unix.go`).
2. On Windows, loaded DLLs MUST be cached in `loadedDLLs` map with mutex protection.
3. On Windows, `CloseLibrary()` MUST be provided to release DLL resources on application exit.
4. On Unix, `purego.Dlopen` with `RTLD_NOW|RTLD_GLOBAL` flags is used — no caching needed.
5. The library handle (uintptr) is passed to all core implementation constructors.

## Scope
- `internal/core/open_windows.go`
- `internal/core/open_unix.go`
- `internal/core/core.go` (calls `openLibrary`)
