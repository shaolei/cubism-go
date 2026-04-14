# Core Minimum — Interface

## Core
```go
type Core struct { ... }

func NewCore(lib uintptr) (Core, error)
func (c Core) GetVersion() string
```

## Failure Modes
- `NewCore`: DLL handle invalid → purego registration may panic
- `GetVersion`: Returns version string like "5.0.0" or "6.0.1"
