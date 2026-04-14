# Live2D Model Rendering

## Business Context
A user (game developer or application integrator) wants to display animated Live2D characters in their Go application. They provide a native Cubism Core DLL and a model3.json file, and the library handles all loading, animation, and rendering.

## Trigger
User calls `Cubism.LoadModel(path)` after initializing `Cubism` with a DLL path.

## Goal
Display a fully animated Live2D model on screen with motion playback, eye blinking, and audio, updating at 60 FPS.

## Participants
1. **cubism-go**: Orchestrates model loading, provides public API
2. **core**: Routes to correct native implementation, manages DLL lifecycle
3. **core/base**: Reads native model data, executes model updates
4. **core/data-types**: Provides typed data structures for native data
5. **motion**: Applies animation curves to model parameters each frame
6. **blink**: Drives periodic eye blink animation
7. **model-json**: Parses JSON configuration files into domain types
8. **renderer**: Draws the model to screen using Ebitengine
9. **sound**: Plays audio files associated with motions

## Paths

### Happy Path
1. User calls `NewCubism("Live2DCubismCore.dll")` → DLL loaded, version detected, correct Core implementation selected
2. User calls `Cubism.LoadModel("model.model3.json")` → all resources loaded and parsed
3. User creates `Renderer` with the loaded model → textures loaded, framebuffers created
4. Each frame:
   a. Renderer calls `model.Update(delta)` → motion/blink parameters applied, native model updated
   b. Renderer calls `renderer.Update()` → vertex positions refreshed
   c. Renderer calls `renderer.Draw(screen)` → drawables rendered in sorted order
5. User optionally calls `model.PlayMotion("Idle", 0, true)` → idle animation loops
6. User optionally calls `model.EnableAutoBlink()` → periodic blinking

### Error Paths
- DLL not found or incompatible → `NewCubism` returns error
- Model file not found → `LoadModel` returns error
- Moc3 consistency check fails → `LoadModel` returns error
- Texture file missing → `NewRenderer` returns error
- Audio file missing → `LoadModel` returns error (if LoadSound is set)

## Invariants
- `Model.Update()` MUST be called before reading drawable state each frame
- Only one motion is actively evaluated at a time (last in queue)
- Core DLL MUST remain loaded for the entire lifetime of the model
- `Model.Close()` MUST be called when done to release native resources
