# Renderer — Responsibility

## Identity
The `renderer/ebitengine` package renders Live2D Cubism models using the Ebitengine game framework, including mask rendering via Kage shaders and coordinate-based hit detection.

## Responsibility
- **Renderer**: Manages framebuffers, textures, vertex transformation, and mask rendering
- **NewRenderer()**: Loads textures, creates framebuffers, compiles mask shader
- **Update()**: Refreshes vertex positions from model drawables each frame
- **Draw()**: Renders all visible drawables in sorted order with mask support
- **IsHit()**: Hit detection by normalizing screen coordinates to model space
- Draw options: `WithHidden()`, `WithScale()`, `WithPosition()`, `WithBackground()`

## NOT Responsible For
- Model updates (calls `model.Update()` internally but doesn't control timing)
- Audio playback

## Key Invariants
- Vertex transformation: model space [-1,1] → screen space via `(pos+1)/2 × size`
- UV Y-axis is flipped: `(1 - uv.Y)` because Live2D uses bottom-left origin
- Mask rendering uses a Kage shader that composites mask and fragment buffers
- Drawables are rendered in sorted index order (from `GetSortedIndices()`)
- Opacity is applied per-drawable via `colorm.ColorM.Scale()`
