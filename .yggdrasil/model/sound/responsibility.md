# Sound — Responsibility

## Identity
The `sound` package hierarchy provides audio playback for Live2D motion sounds with three strategy implementations.

## Responsibility
- **sound.Sound interface**: `Play() error` + `Close()` — contract for all audio implementations
- **sound/normal**: Eager decoding — loads and decodes audio file immediately on `LoadSound()`
- **sound/delay**: Lazy decoding — stores file path, decodes only on first `Play()` call
- **sound/disabled**: No-op — used when `Cubism.LoadSound` is nil
- **sound/audioutils**: Shared utilities — format detection, audio decoding, speaker initialization

## NOT Responsible For
- When to play sounds (controlled by MotionManager)
- Which sound strategy to use (configured by caller via `Cubism.LoadSound`)

## Key Invariants
- Speaker is initialized only once (guarded by `speakerInitMu` mutex)
- Supported formats: WAV (.wav, .wave) and MP3 (.mp3)
- All implementations satisfy `sound.Sound` interface
