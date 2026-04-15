# Example

A complete working example demonstrating cubism-go with the Ebitengine renderer.

## Directory Structure

```
example/
├── README.md
├── main.go
├── Live2DCubismCore.dll   # Cubism Core library (not included)
└── Resources/              # Live2D model files (not included)
    └── Haru/
        ├── Haru.model3.json
        ├── Haru.moc3
        ├── Haru.physics3.json
        ├── *.png
        └── ...
```

## Setup

1. Download the [Cubism SDK for Native](https://www.live2d.com/download/cubism-sdk/download-native/) and place the Core dynamic library in this directory:
   - Windows: `Live2DCubismCore.dll`
   - macOS: `libLive2DCubismCore.dylib`
   - Linux: `libLive2DCubismCore.so`

2. Place the Live2D model resources in the `Resources/` directory.

## Running

```sh
cd example
go run main.go
```

### Command Line Flags

- `-model` — Path to the `.model3.json` file (default: `example/Resources/Haru/Haru.model3.json`)
- `-lib` — Path to the Cubism Core library (default: `example/Live2DCubismCore.dll`)

## Controls

| Key | Action |
|-----|--------|
| **Click** | Click on the model's hit areas to trigger the `TapBody` motion |
| **0** | Stop the current expression |
| **1-9** | Switch to the corresponding expression |
| **L** | Toggle look/eye-tracking (follows mouse cursor) |
| **B** | Enable the breathing effect |
| **P** | Disable physics simulation |
| **R** | Reset physics to default gravity/wind |

The cursor changes to a pointer when hovering over a hit area.

## Architecture

The example demonstrates the core features of cubism-go:

- **CubismModel / CubismUserModel** — The model is split into a pure data layer (`CubismModel`) and a composition manager (`CubismUserModel`). The public `Model` type acts as a facade.
- **Motion** — Idle motion plays on startup; tap motions are triggered by clicking hit areas.
- **Expression** — Switch between expressions using number keys.
- **Look** — Eye-tracking follows the mouse cursor in real-time.
- **Breath** — Subtle body movement via sine wave oscillation.
- **Physics** — Hair and clothing physics with configurable gravity/wind.
- **Auto Blink** — Automatic eye blinking that pauses during motion playback.
