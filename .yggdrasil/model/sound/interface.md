# Sound — Interface

## Sound Interface
```go
type Sound interface {
    Play() error
    Close()
}
```

## LoadSound Functions
```go
// Eager decoding
func normal.LoadSound(fp string) (sound.Sound, error)
// Lazy decoding
func delay.LoadSound(fp string) (sound.Sound, error)
// No-op
func disabled.LoadSound(fp string) (sound.Sound, error)
```

## Audio Utilities
```go
func audioutils.DetectFormat(fp string) (string, error)
func audioutils.DecodeAudio(format string, buf []byte) (beep.StreamSeekCloser, beep.Format, error)
func audioutils.InitSpeaker(format beep.Format) error
```

## Failure Modes
- `DetectFormat`: Unsupported file extension → error
- `DecodeAudio`: Unsupported format → error; corrupt audio data → error
- `InitSpeaker`: Audio driver unavailable → error
- `disabled.Play()`: Always returns nil (no-op)
