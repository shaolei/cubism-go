# Core Router — Responsibility

## Identity
The `internal/core` package defines the `Core` interface and provides the `NewCore()` factory function that loads the native DLL and routes to the correct version-specific implementation.

## Responsibility
- Define the `Core` interface that all version implementations must satisfy
- Load the native Cubism Core DLL/SO via platform-specific `openLibrary()`
- Query the DLL version via `minimum.NewCore()` to determine the major version
- Route to `core_5_0_0.NewCore()` (v5) or `core_6_0_1.NewCore()` (v6) based on major version
- Provide `CloseLibrary()` on Windows for DLL resource cleanup

## NOT Responsible For
- Actual FFI function implementation (delegated to child nodes)
- Data type definitions (delegated to `core/data-types`)

## Key Invariants
- `NewCore()` MUST return an error for unsupported major versions (not v5 or v6)
- The `Core` interface is the contract between all consumers and any version-specific implementation
- `openLibrary()` is platform-specific and build-tagged
