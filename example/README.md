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

## Controls

- Click on the model's hit areas to trigger the `TapBody` motion
- The cursor changes to a pointer when hovering over a hit area
- Idle motion plays automatically on startup
