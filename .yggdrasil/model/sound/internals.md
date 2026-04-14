# Sound — Internals

## Logic

### Normal Sound Flow
1. Read entire file into memory
2. Detect format from extension
3. Decode audio stream immediately
4. Create `beep.Ctrl` wrapper
5. Initialize speaker (once)
6. On Play: seek to 0, unpause, play via speaker
7. On Close: pause, seek to 0

### Delay Sound Flow
1. Store file path only
2. On first Play():
   - Read file from disk
   - Detect format, decode, init speaker
   - Create Ctrl wrapper
3. Subsequent Play(): same as normal

### Speaker Initialization
- `InitSpeaker()` uses sync.Mutex to ensure only one call to `speaker.Init()`
- Sample rate from decoded audio format
- Buffer size: SampleRate.N(time.Second/10) = 1/10 second

## Decisions
- Chose beep/oto over PortAudio: pure Go, no CGO, simpler build
- Chose delay strategy for default: avoids loading all sounds at model load time, reducing memory and startup time
- Chose mutex-guarded speaker init: `speaker.Init()` must be called exactly once; concurrent calls from multiple goroutines would panic
