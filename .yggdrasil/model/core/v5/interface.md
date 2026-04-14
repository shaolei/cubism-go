# Core v5 — Interface

Implements the `Core` interface (see core router interface.md) for Cubism SDK v5.x.

## Constructor
```go
func NewCore(lib uintptr) (*Core, error)
```

## Version-Specific Detail
- Maps `CsmGetDrawableSortOrders` to native function `csmGetDrawableRenderOrders`
- All other methods delegate to `base` package functions
