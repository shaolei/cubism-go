# Core v5 — Responsibility

## Identity
The `internal/core/core_5_0_0` package implements the `Core` interface for Cubism SDK version 5.x.

## Responsibility
- Register all common FFI functions via `base.RegisterCommonFuncs()`
- Register v5-specific `csmGetDrawableRenderOrders` as the sort order function
- Delegate all Core interface methods to `base` package functions

## NOT Responsible For
- Any logic beyond FFI registration and delegation (all logic is in `base`)

## Key Difference from v6
- Maps `CsmGetDrawableSortOrders` to `csmGetDrawableRenderOrders`
