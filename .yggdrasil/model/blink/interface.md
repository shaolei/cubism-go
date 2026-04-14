# Blink — Interface

## BlinkManager
```go
func NewBlinkManager(core core.Core, modelPtr uintptr, ids []string) *BlinkManager
func (b *BlinkManager) Update(delta float64)
func (b *BlinkManager) DetermineNextBlinkingTiming() float64
```

## State Constants
```go
const (
    EyeStateFirst    = iota
    EyeStateInterval
    EyeStateClosing
    EyeStateClosed
    EyeStateOpening
)
```

## Parameters
- `ids`: Parameter IDs for eye blink (typically from "EyeBlink" group in model3.json)
- `delta`: Time elapsed since last frame (seconds)
