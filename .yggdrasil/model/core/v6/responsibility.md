# Core v6 — Responsibility

## Identity
The `internal/core/core_6_0_1` package implements the `Core` interface for Cubism SDK version 6.x.

## Responsibility
- Register all common FFI functions via `base.RegisterCommonFuncs()`
- Register v6-specific `csmGetDrawableDrawOrders` as the sort order function
- Delegate all Core interface methods to `base` package functions

## NOT Responsible For
- Any logic beyond FFI registration and delegation (all logic is in `base`)

## Key Difference from v5
- Maps `CsmGetDrawableSortOrders` to `csmGetDrawableDrawOrders`
