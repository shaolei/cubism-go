# Cubism Go

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/shaolei/cubism-go.svg)](https://pkg.go.dev/github.com/shaolei/cubism-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/shaolei/cubism-go)](https://goreportcard.com/report/github.com/shaolei/cubism-go)
[![CI](https://github.com/shaolei/cubism-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/shaolei/cubism-go/actions/workflows/ci.yaml)

cubism-goは[Live2D Cubism SDK](https://www.live2d.com/sdk/about/)の非公式版のGolang実装です。[ebitengine/purego](https://github.com/ebitengine/purego)を用いてCubism Coreのネイティブライブラリを呼び出すため、CGO不要でクロスプラットフォームに利用できます。

## 特徴

- **Pure Go + purego** — CGO不要。purego経由でCubism Core動的ライブラリを呼び出し
- **マルチバージョン対応** — Cubism Core 5.x / 6.xを実行時に自動検出
- **レンダリング** — Ebitengineベースのレンダラ（マスク対応）を内蔵。カスタムレンダラも可能
- **音声再生** — プラグイン式サウンドシステム（即時読み込み、遅延読み込み、無効化、カスタム）
- **モーション＆自動瞬き** — フェードイン/アウト付きモーション再生、ループ対応、自動瞬き
- **ヒット判定** — クリック/タップ判定用のヒットエリア検出

## インストール

```bash
go get -u github.com/shaolei/cubism-go
```

### 動作に必要なもの

- Go 1.25以上
- Cubism Core動的ライブラリ（Windows: `.dll`、macOS: `.dylib`、Linux: `.so`）
  - [Live2D Cubism SDK](https://www.live2d.com/download/cubism-sdk/download-native/)から入手
- Live2Dモデル（`.model3.json`と関連ファイル）

## クイックスタート

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
    // 1. Cubism Coreライブラリのパスを指定して初期化
    csm, err := cubism.NewCubism("Live2DCubismCore.dll")
    if err != nil {
        log.Fatal(err)
    }

    // 2. 音声ローダーを設定（省略可 — 音声を無効化する場合は設定不要）
    csm.LoadSound = normal.LoadSound

    // 3. model3.jsonからモデルを読み込み
    model, err := csm.LoadModel("Resources/Haru/Haru.model3.json")
    if err != nil {
        log.Fatal(err)
    }
    defer model.Close()

    // 4. アイドルモーションを再生し、自動瞬きを有効化
    model.PlayMotion("Idle", 0, true)
    model.EnableAutoBlink()

    // 5. レンダラを作成してEbitenで実行
    r, err := renderer.NewRenderer(model)
    if err != nil {
        log.Fatal(err)
    }
    // ... Ebitenのゲームループで r.Update() と r.Draw() を使用
}
```

完全な動作例は [`example/`](example/) ディレクトリを参照してください。

## プロジェクト構成

```
cubism-go/
├── cubism.go              # Cubismエントリポイント（Cubism構造体、LoadModel）
├── model.go               # Model構造体（モーション、瞬き、更新、パラメータ）
├── drawable.go            # Drawable構造体（公開API）
├── internal/
│   ├── blink/             # 自動瞬きステートマシン
│   ├── core/              # Cubism Coreバインディング
│   │   ├── base/          # 共通Core実装（関数登録、moc読み込み）
│   │   ├── core_5_0_0/    # Cubism Core 5.xアダプタ
│   │   ├── core_6_0_1/    # Cubism Core 6.xアダプタ
│   │   ├── minimum/       # バージョン検出のみ
│   │   ├── drawable/      # Drawable型とフラグ解析
│   │   ├── moc/           # Mocリソース管理
│   │   └── parameter/     # Parameter型
│   ├── model/             # JSONモデルパーサ（model3、motion、physics等）
│   ├── motion/            # モーションマネージャと補間（linear、bezier、stepped）
│   ├── strings/           # C文字列からGo文字列への変換
│   └── utils/             # バージョンパースユーティリティ
├── renderer/
│   ├── ebitengine/        # Ebitengineレンダラ（マスクシェーダ付き）
│   └── utils/             # Normalizeユーティリティ
├── sound/
│   ├── audioutils/        # 共通音声デコード（WAV/MP3）とスピーカー初期化
│   ├── normal/            # 即時読み込み音声実装
│   ├── delay/             # 遅延読み込み音声実装
│   └── disabled/          # No-op音声実装
└── example/               # 動作例アプリケーション
```

## 音声実装

| パッケージ | 説明 |
|---|---|
| `sound/normal` | `LoadSound`時に即座に読み込み・デコード |
| `sound/delay` | `Play()`が呼ばれるまで読み込み・デコードを遅延 |
| `sound/disabled` | 何もしない実装。音声が不要な場合に使用 |
| カスタム | `sound.Sound`インターフェース（`Play() error`、`Close()`）を実装 |

## レンダリング

`renderer/ebitengine`パッケージはEbitengineベースのレンダラを提供します：

- 頂点位置からスクリーン座標への自動変換
- Kageシェーダによるマスクレンダリング
- 設定可能な描画オプション（位置、スケール、背景色、非表示モード）
- ヒット判定（`IsHit`）によるインタラクティブなクリックエリア

`Model` APIを直接利用してカスタムレンダラを実装することも可能です。

## Coreバージョン対応

ライブラリはCubism Coreのバージョンを自動検出し、適切なアダプタを選択します：

- **Cubism Core 5.x** → `csmGetDrawableRenderOrders`でソート
- **Cubism Core 6.x** → `csmGetDrawableDrawOrders`でソート

その他のAPIは`internal/core/base`パッケージで共通化されています。

## APIリファレンス

完全なAPIドキュメントは[Go Reference](https://pkg.go.dev/github.com/shaolei/cubism-go)を参照してください。

### 主要な型

- `Cubism` — エントリポイント。`NewCubism(libPath)`で初期化、`LoadModel(path)`でモデル読み込み
- `Model` — Live2Dモデル。モーション再生、自動瞬き、パラメータ取得/設定、更新サイクルをサポート
- `Drawable` — 頂点位置、UV、不透明度、フラグを持つ視覚要素

## 開発

### 前提条件

`pre-commit`フックのために[lefthook](https://github.com/evilmartians/lefthook)を利用しています：

- [staticcheck](https://staticcheck.dev)
- [typos](https://github.com/crate-ci/typos)

Homebrewでインストール：

```sh
brew install lefthook staticcheck typos-cli
lefthook install
```

### テストの実行

```sh
go test ./... -cover
```

### サンプルの実行

1. Cubism Coreライブラリ（例: `Live2DCubismCore.dll`）を`example/`ディレクトリに配置
2. モデルリソースを`example/Resources/`に配置
3. 実行：

```sh
cd example && go run main.go
```

## フォークの違い

このリポジトリは[aethiopicuschan/cubism-go](https://github.com/aethiopicuschan/cubism-go)のフォークであり、大幅な機能拡張とリファクタリングを行っています。以下に全変更内容の詳細を示します。

### フォークの理由

上流リポジトリはCubism Core 5.xのみをサポート（バージョンチェックが`if version == "5.0.0"`でハードコード）しており、新しいバージョンへの対応がありません。また、いくつかの重要なバグ（DLL二重読み込み、リソース解放の欠落、正規化でのゼロ除算）やアーキテクチャ上の問題（音声実装間のコード重複、v5コアのモノリシック構造）に対処する必要がありました。フォークの目的は以下の通りです：

1. **Cubism Core 6.x対応の追加** — 上流ではCubism Editor 5.1以降で同梱されるCore 6.xのモデルを読み込めない
2. **リソース管理の修正** — 上流には`Close()`メソッドがなく、長時間実行アプリケーションでメモリリークが発生
3. **コード重複の排除** — 音声実装でフォーマット検出とデコードロジックが重複、v5コアは300行以上が共有可能
4. **ランタイムバグの修正** — WindowsでのDLL二重読み込み、正規化でのゼロ除算、ヒット判定の初期化エラー

### 新機能

| 機能 | 説明 |
|---|---|
| **Cubism Core 6.x対応** | 新規`internal/core/core_6_0_1/`アダプタ（`csmGetDrawableDrawOrders`使用）。バージョンルーティングが`parseMajorVersion()`でメジャーバージョン5/6に対応 |
| **共通Core baseパッケージ** | 新規`internal/core/base/`パッケージ。共通FFI関数登録、moc読み込み、Drawable/Parameter/Canvas操作を抽出 — v5とv6間の約500行の重複を排除 |
| **オーディオユーティリティパッケージ** | 新規`sound/audioutils/`パッケージ。`DetectFormat()`、`DecodeAudio()`、`InitSpeaker()`を集約 — 以前は`sound/normal/`と`sound/delay/`で重複していた |
| **リソース解放API** | `Model.Close()`と`moc.Moc.Close()`による適切なリソース解放。Windows向け`core.CloseLibrary()`でDLL解放 |
| **スピーカー初期化の安全性** | `audioutils.InitSpeaker()`が`sync.Mutex`でレースコンディションを防止。上流は保護されていない`initialized`ブール変数を使用 |
| **DLLキャッシュ** | WindowsのDLL読み込みが`sync.Mutex`でキャッシュし、二重読み込みを防止 |
| **新規テスト** | `internal/blink/blink_manager_test.go`（瞬きステートマシン）、`internal/strings/strings_test.go`（C文字列変換） |

### バグ修正

| 修正 | 説明 |
|---|---|
| **バージョンパース** | `internal/utils/version.go` — ビットレイアウトの実際の動作を反映（`0x06000001` = 6.0.1）。テストカバレッジを16進リテラルとエッジケースで拡充 |
| **ゼロ除算** | `renderer/utils/normalize.go` — `Normalize()`が`n == m`の場合にNaN/Infではなく0を返すよう修正 |
| **ヒット判定の境界** | `renderer/ebitengine/renderer.go` — `IsHit()`が画面サーフェスサイズではなく最初の頂点からバウンディングボックスを初期化するよう修正 |
| **テストの独立性** | `version_test.go`と`normalize_test.go` — `testify`への依存を削除し、標準`testing`パッケージを使用 |

### リファクタリング

| 変更 | 説明 |
|---|---|
| **Core v5の委譲** | `internal/core/core_5_0_0/core.go`が約340行から約80行に削減。`base`パッケージの関数に委譲 |
| **音声の重複排除** | `sound/normal/`と`sound/delay/`が各約90行から約60行に削減。`audioutils`パッケージを使用し、重複する`nopCloser`型と`detectFormat()`関数を削除 |
| **.gitignoreの拡充** | IDEファイル、ビルド成果物、カバレッジ出力、vendorディレクトリのパターンを追加 |
| **依存関係の更新** | `purego` 0.7.1 → 0.10.0、`ebiten/v2` 2.7.10 → 2.9.9、`x/sys` 0.25.0 → 0.43.0、Go 1.22 → 1.25 |

### メンテナンスと同期戦略

- **上流のリベース**: このフォークは定期的に上流の変更をリベースします。コアのリファクタリング（baseパッケージ、v6アダプタ）により`internal/core/`でマージコンフリクトが発生する可能性がありますが、手動で解決します。
- **機能パリティ**: 上流の新機能（新しいパラメータ型、レンダラの改善など）は統合し、共通baseパッケージアーキテクチャに適応させます。
- **上流PRの予定なし**: 変更が大規模であるため、部分的なPRは困難です。将来的に上流が類似のアーキテクチャを採用した場合、収束を検討します。
- **課題管理**: フォーク固有の問題はこのリポジトリのIssueトラッカーで管理します。

## ライセンス

このプロジェクトは[MITライセンス](LICENSE)の下で公開されています。

**注意:** Cubism Coreライブラリはプロプライエタリであり、Live2D Inc.のライセンス条件に従います。このプロジェクトはCubism Coreライブラリを再配布しません。
