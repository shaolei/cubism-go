# Motion — Interface

## MotionManager
```go
func NewMotionManager(core core.Core, modelPtr uintptr, onFinished func(int)) *MotionManager
func (mm *MotionManager) Start(motion Motion) int      // Returns motion ID
func (mm *MotionManager) Close(id int)
func (mm *MotionManager) Reset(id int)
func (mm *MotionManager) Update(deltaTime float64)
```

## Data Types
```go
type Motion struct {
    File        string
    FadeInTime  float64
    FadeOutTime float64
    Sound       string
    LoadedSound sound.Sound
    Meta        Meta
    Curves      []Curve
}

type Curve struct {
    Target      string   // "Model", "Parameter", "PartOpacity"
    Id          string
    FadeInTime  float64  // -1 = use motion fade
    FadeOutTime float64  // -1 = use motion fade
    Segments    []Segment
}

type Segment struct {
    Points []Point
    Type   int  // Linear, Bezier, Stepped, InverseStepped
    Value  float64
}

type Point struct { Time, Value float64 }
type Meta struct {
    Duration             float64
    Loop                 bool
    AreBeziersRestricted bool
}
```

## Failure Modes
- `Start` with invalid motion data: no error, but playback may produce zero values
- `Close` with non-existent ID: no-op
- Sound playback failure: `LoadedSound.Play()` error is silently ignored
