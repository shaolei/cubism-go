# Utils — Interface

## internal/utils
```go
func ParseVersion(v uint32) string
```

## internal/strings
```go
func GoString(p uintptr) string
```

## renderer/utils
```go
func Normalize(x, n, m float32) float32
```

## Failure Modes
- `GoString` with nil pointer: returns empty string (no error)
- `Normalize` with n == m: returns 0 (no division by zero)
