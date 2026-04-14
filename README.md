# Cubism Go

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/shaolei/cubism-go.svg)](https://pkg.go.dev/github.com/shaolei/cubism-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/shaolei/cubism-go)](https://goreportcard.com/report/github.com/shaolei/cubism-go)
[![CI](https://github.com/shaolei/cubism-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/shaolei/cubism-go/actions/workflows/ci.yaml)

cubism-go is an unofficial Go implementation of the [Live2D Cubism SDK](https://www.live2d.com/sdk/about/). It leverages [ebitengine/purego](https://github.com/ebitengine/purego) to call the Cubism Core native library without CGO, making it cross-platform and easy to integrate.

## Features

- **Pure Go + purego** — No CGO required; calls the Cubism Core dynamic library via purego
- **Multi-version support** — Supports Cubism Core 5.x and 6.x, automatically detected at runtime
- **Rendering** — Built-in Ebitengine renderer with mask support, or bring your own
- **Audio playback** — Pluggable sound system (normal, delayed loading, disabled, or custom)
- **Motion & Blink** — Motion playback with fade in/out, loop support; automatic eye blink
- **Hit detection** — Click/hit area detection for interactive applications

## Installation

```bash
go get -u github.com/shaolei/cubism-go
```

### Requirements

- Go 1.25 or later
- Cubism Core dynamic library (`Live2DCubismCore.dll` on Windows, `.dylib` on macOS, `.so` on Linux)
  - Obtain from [Live2D Cubism SDK](https://www.live2d.com/download/cubism-sdk/download-native/)
- A Live2D model (`.model3.json` with associated files)

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/shaolei/cubism-go"
    renderer "github.com/shaolei/cubism-go/renderer/ebitengine"
    "github.com/shaolei/cubism-go/sound/normal"
    "github.com/hajimehoshi/ebiten/v2"
)

func main() {
    // 1. Initialize Cubism with the Core library path
    csm, err := cubism.NewCubism("Live2DCubismCore.dll")
    if err != nil {
        log.Fatal(err)
    }

    // 2. Set audio loader (optional — omit to disable sound)
    csm.LoadSound = normal.LoadSound

    // 3. Load a model from model3.json
    model, err := csm.LoadModel("Resources/Haru/Haru.model3.json")
    if err != nil {
        log.Fatal(err)
    }
    defer model.Close()

    // 4. Play idle motion and enable auto blink
    model.PlayMotion("Idle", 0, true)
    model.EnableAutoBlink()

    // 5. Create a renderer and run with Ebiten
    r, err := renderer.NewRenderer(model)
    if err != nil {
        log.Fatal(err)
    }
    // ... use r.Update() and r.Draw() in your Ebiten game loop
}
```

See the [`example/`](example/) directory for a complete working example.

## Project Structure

```
cubism-go/
├── cubism.go              # Cubism entry point (Cubism struct, LoadModel)
├── model.go               # Model struct (motion, blink, update, parameters)
├── drawable.go            # Drawable struct (public API)
├── internal/
│   ├── blink/             # Auto-blink state machine
│   ├── core/              # Cubism Core bindings
│   │   ├── base/          # Shared core implementation (func registration, moc loading)
│   │   ├── core_5_0_0/    # Cubism Core 5.x adapter
│   │   ├── core_6_0_1/    # Cubism Core 6.x adapter
│   │   ├── minimum/       # Version detection only
│   │   ├── drawable/      # Drawable types and flag parsing
│   │   ├── moc/           # Moc resource management
│   │   └── parameter/     # Parameter types
│   ├── model/             # JSON model parsers (model3, motion, physics, etc.)
│   ├── motion/            # Motion manager and interpolation (linear, bezier, stepped)
│   ├── strings/           # C string to Go string conversion
│   └── utils/             # Version parsing utility
├── renderer/
│   ├── ebitengine/        # Ebitengine renderer with mask shader
│   └── utils/             # Normalize utility
├── sound/
│   ├── audioutils/        # Shared audio decoding (WAV/MP3) and speaker init
│   ├── normal/            # Immediate-load sound implementation
│   ├── delay/             # Lazy-load sound implementation
│   └── disabled/          # No-op sound implementation
└── example/               # Working example application
```

## Audio Implementations

| Package | Description |
|---|---|
| `sound/normal` | Loads and decodes audio immediately on `LoadSound` |
| `sound/delay` | Defers loading and decoding until `Play()` is called |
| `sound/disabled` | No-op implementation; use when audio is not needed |
| Custom | Implement the `sound.Sound` interface (`Play() error`, `Close()`) |

## Rendering

The `renderer/ebitengine` package provides an Ebitengine-based renderer with:

- Automatic vertex position to screen coordinate conversion
- Mask rendering via Kage shader
- Configurable draw options (position, scale, background color, hidden mode)
- Hit detection (`IsHit`) for interactive click areas

You can also implement your own renderer by consuming the `Model` API directly.

## Core Version Support

The library automatically detects the Cubism Core version and selects the appropriate adapter:

- **Cubism Core 5.x** → uses `csmGetDrawableRenderOrders` for sorting
- **Cubism Core 6.x** → uses `csmGetDrawableDrawOrders` for sorting

Other APIs are shared across versions via the `internal/core/base` package.

## API Reference

See the [Go Reference](https://pkg.go.dev/github.com/shaolei/cubism-go) for complete API documentation.

### Key Types

- `Cubism` — Entry point; initialize with `NewCubism(libPath)` and load models with `LoadModel(path)`
- `Model` — Live2D model; supports motion playback, auto-blink, parameter get/set, and update cycle
- `Drawable` — Visual element with vertex positions, UVs, opacity, and flags

## Development

### Prerequisites

For pre-commit hooks, we use [lefthook](https://github.com/evilmartians/lefthook):

- [staticcheck](https://staticcheck.dev)
- [typos](https://github.com/crate-ci/typos)

Install with Homebrew:

```sh
brew install lefthook staticcheck typos-cli
lefthook install
```

### Running Tests

```sh
go test ./... -cover
```

### Running the Example

1. Place the Cubism Core library (e.g., `Live2DCubismCore.dll`) in the `example/` directory
2. Place the model resources in `example/Resources/`
3. Run:

```sh
cd example && go run main.go
```

## Fork Differences

This repository is a fork of [aethiopicuschan/cubism-go](https://github.com/aethiopicuschan/cubism-go) with significant enhancements and refactoring. Below is a detailed summary of all changes.

### Why Fork

The upstream repository only supports Cubism Core 5.x (hardcoded version check `if version == "5.0.0"`), with no path to support newer versions. Additionally, several critical bugs (DLL double-loading, missing resource cleanup, division by zero in normalization) and architectural issues (duplicated code across sound implementations, monolithic v5 core) needed to be addressed. The fork was created to:

1. **Add Cubism Core 6.x support** — The upstream cannot load models built with Cubism Editor 5.1+ that ship with Core 6.x
2. **Fix resource management** — The upstream has no `Close()` methods, leading to memory leaks in long-running applications
3. **Eliminate code duplication** — Sound implementations duplicated format detection and decoding logic; core v5 had 300+ lines that could be shared
4. **Fix runtime bugs** — DLL double-loading on Windows, division by zero in normalization, hit detection initialization error

### New Features

| Feature | Description |
|---|---|
| **Cubism Core 6.x support** | New `internal/core/core_6_0_1/` adapter using `csmGetDrawableDrawOrders`; version routing now supports major versions 5 and 6 via `parseMajorVersion()` |
| **Shared core base package** | New `internal/core/base/` package extracting common FFI function registration, moc loading, drawable/parameter/canvas operations — eliminates ~500 lines of duplication between v5 and v6 |
| **Audio utilities package** | New `sound/audioutils/` package centralizing `DetectFormat()`, `DecodeAudio()`, and `InitSpeaker()` — previously duplicated in `sound/normal/` and `sound/delay/` |
| **Resource cleanup API** | `Model.Close()` and `moc.Moc.Close()` for proper resource release; `core.CloseLibrary()` for DLL unloading on Windows |
| **Speaker init safety** | `audioutils.InitSpeaker()` uses `sync.Mutex` to prevent race conditions; upstream used an unprotected `initialized` bool |
| **DLL caching** | Windows DLL loading now caches loaded DLLs with `sync.Mutex`, preventing duplicate loads |
| **New tests** | `internal/blink/blink_manager_test.go` (blink state machine), `internal/strings/strings_test.go` (C string conversion) |

### Bug Fixes

| Fix | Description |
|---|---|
| **Version parsing** | `internal/utils/version.go` — comments corrected to reflect actual bit layout (`0x06000001` = 6.0.1); test coverage expanded with hex literals and edge cases |
| **Division by zero** | `renderer/utils/normalize.go` — `Normalize()` now returns 0 when `n == m` instead of producing NaN/Inf |
| **Hit detection bounds** | `renderer/ebitengine/renderer.go` — `IsHit()` now initializes bounding box from first vertex instead of screen surface size, fixing incorrect hit areas |
| **Test independence** | `version_test.go` and `normalize_test.go` — removed dependency on `testify`, using standard `testing` package for zero external test dependencies |

### Refactoring

| Change | Description |
|---|---|
| **Core v5 delegation** | `internal/core/core_5_0_0/core.go` reduced from ~340 lines to ~80 lines by delegating to `base` package functions |
| **Sound deduplication** | `sound/normal/` and `sound/delay/` reduced from ~90 lines each to ~60 lines by using `audioutils` package; removed duplicate `nopCloser` type and `detectFormat()` function |
| **.gitignore expanded** | Added IDE files, build artifacts, coverage output, and vendor directory patterns |
| **Dependencies updated** | `purego` 0.7.1 → 0.10.0, `ebiten/v2` 2.7.10 → 2.9.9, `x/sys` 0.25.0 → 0.43.0, Go 1.22 → 1.25 |

### Maintenance & Sync Strategy

- **Upstream rebase**: This fork will periodically rebase on upstream changes. Core refactoring (base package, v6 adapter) may cause merge conflicts in `internal/core/` — these will be resolved manually.
- **Feature parity**: New upstream features (e.g., new parameter types, renderer improvements) will be integrated and adapted to the shared base package architecture.
- **No upstream PR plans**: The changes are substantial enough that a partial PR would be difficult. If the upstream adopts a similar architecture in the future, convergence will be explored.
- **Issue tracking**: Fork-specific issues are tracked in this repository's issue tracker.

## License

This project is licensed under the [MIT License](LICENSE).

**Note:** The Cubism Core library is proprietary and subject to its own license terms from Live2D Inc. This project does not distribute the Cubism Core library.
