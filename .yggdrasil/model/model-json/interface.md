# Model JSON — Interface

## ModelJson
```go
type ModelJson struct {
    Version        int
    FileReferences struct {
        Moc         string
        Textures    []string
        Physics     string
        Pose        string
        DisplayInfo string
        Expressions []struct { Name, File string }
        Motions     map[string][]struct { File string; FadeInTime, FadeOutTime float64; Sound, MotionSync string }
        UserData    string
    }
    Groups   []Group
    HitAreas []HitArea
}
```

## ExpJson
```go
type ExpJson struct {
    Name       string
    Type       string
    Parameters []struct { Id string; Value float64; Blend string }
}
```

## MotionJson
```go
type MotionJson struct {
    Version  int
    Meta     Meta
    UserData []struct { Time float64; Value string }
    Curves   []struct { Target, Id string; FadeInTime, FadeOutTime *float64; Segments []float64 }
}
func (m *MotionJson) ToMotion(fp string, fadein, fadeout float64, sound string) motion.Motion
```

## Supporting Types
```go
type Group struct { Target, Name string; Ids []string }
type HitArea struct { Id, Name string }
type Meta struct { Duration float64; Loop, AreBeziersRestricted bool }
type PhysicsJson struct { ... }
type PoseJson struct { ... }
type CdiJson struct { ... }
type UserDataJson struct { ... }
```
