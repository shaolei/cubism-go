# Renderer — Internals

## Logic

### Draw Flow
1. Calculate final transform: scale to screen → apply user scale → center with position offset
2. Clear surface with background color
3. For each drawable in sorted order:
   - Skip if not visible
   - If drawable has masks:
     - Clear mask buffer and fragment buffer
     - Draw each mask into mask buffer (black, alpha-only)
     - Draw drawable into fragment buffer
     - Composite mask + fragment via shader → surface
   - If no masks:
     - Draw triangles directly to surface with opacity
4. Draw surface to screen with final transform

### Mask Shader (mask.kage)
- Takes two images: mask buffer (image0) and fragment buffer (image1)
- For each pixel: if mask alpha > 0, use fragment color; otherwise transparent

### IsHit Flow
1. Check if coordinates are within the renderer's final bounding rect
2. Normalize screen coordinates to [-1, 1] model space
3. Flip Y axis (model Y is up, screen Y is down)
4. Check if normalized point falls within the drawable's bounding box (min/max of all vertex positions)

## Decisions
- Chose three-buffer approach (surface + fb + mb) for mask rendering: allows compositing masks before drawing to the final surface, preventing alpha blending artifacts
- Chose Kage shader over pure Go image manipulation: GPU-accelerated mask compositing is significantly faster
- Chose bounding-box hit detection over pixel-perfect: simpler, faster, sufficient for most use cases
