# Core v6 — Interface

Implements the `Core` interface (see core router interface.md) for Cubism SDK v6.x.

## Constructor
```go
func NewCore(lib uintptr) (*Core, error)
```

## Version-Specific Detail
- Maps `CsmGetDrawableSortOrders` to native function `csmGetDrawableDrawOrders`
- All other methods delegate to `base` package functions
