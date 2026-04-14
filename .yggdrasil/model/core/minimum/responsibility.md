# Core Minimum — Responsibility

## Identity
The `internal/core/minimum` package provides a minimal Core wrapper that only registers the `csmGetVersion` native function, used to determine which full Core implementation to instantiate.

## Responsibility
- Register only `csmGetVersion` FFI function
- Query version string from native DLL

## NOT Responsible For
- Any model operations (moc loading, parameter access, etc.)
- Determining which version implementation to use (handled by Core Router)

## Key Invariants
- Only used temporarily during initialization to detect the DLL version
- After version routing, this instance is discarded
