package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shaolei/cubism-go"
	"github.com/shaolei/cubism-go/internal/breath"
	"github.com/shaolei/cubism-go/internal/look"
	renderer "github.com/shaolei/cubism-go/renderer/ebitengine"
	"github.com/shaolei/cubism-go/sound/normal"
)

const (
	Width  = 2880
	Height = 1800
)

type Game struct {
	ow, oh      int
	tapId       int
	renderer    *renderer.Renderer
	expressions []string
	model       *cubism.Model
	lookEnabled bool
}

func (g *Game) Update() (err error) {
	g.renderer.Update()

	// Look/eye-tracking: follow the mouse cursor
	if g.lookEnabled {
		x, y := ebiten.CursorPosition()
		// Normalize cursor position to [-1, 1] range
		dragX := float32(2.0*float64(x)/float64(g.ow) - 1.0)
		dragY := float32(-(2.0*float64(y)/float64(g.oh) - 1.0))
		g.model.SetLookTarget(dragX, dragY)
	}

	// Hit area detection for click interaction
	x, y := ebiten.CursorPosition()
	if x < 0 || y < 0 || x > g.ow || y > g.oh {
		return
	}
	if !ebiten.IsFocused() {
		return
	}
	hitareas := g.renderer.GetModel().GetHitAreas()
	hitted := false
	for _, hitarea := range hitareas {
		hit, err := g.renderer.IsHit(x, y, hitarea.Id)
		if err != nil {
			return err
		}
		if hit {
			hitted = true
		}
	}
	if hitted {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.model.StopMotion(g.tapId)
			g.tapId = g.model.PlayMotion("TapBody", 0, false)
		}
	} else if ebiten.CursorShape() == ebiten.CursorShapePointer {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	// Expression switching with number keys (1-9 = expression index, 0 = stop)
	keys := []ebiten.Key{ebiten.Key0, ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5, ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9}
	for i, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			if i == 0 {
				g.model.StopExpression()
				fmt.Println("Expression stopped")
			} else if i <= len(g.expressions) {
				name := g.expressions[i-1]
				g.model.PlayExpression(name)
				fmt.Printf("Playing expression [%d]: %s\n", i, name)
			}
		}
	}

	// L key: toggle look/eye-tracking
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.lookEnabled = !g.lookEnabled
		if g.lookEnabled {
			fmt.Println("Look tracking enabled")
		} else {
			g.model.SetLookTarget(0, 0)
			fmt.Println("Look tracking disabled (reset to center)")
		}
	}

	// B key: toggle breath effect
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.model.EnableBreath()
		fmt.Println("Breath effect enabled")
	}

	// P key: toggle physics
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.model.DisablePhysics()
		fmt.Println("Physics disabled (press R to re-enable)")
	}

	// R key: reset physics
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.model.SetPhysicsOptions(0, 1, 0, 0)
		fmt.Println("Physics reset to default gravity/wind")
	}

	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.renderer.Draw(screen, renderer.WithBackground(color.RGBA{0, 255, 0, 255}))

	// HUD: show TPS/FPS and status info
	dragX, dragY := g.model.GetLookTarget()
	status := fmt.Sprintf("TPS: %0.2f  FPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	status += fmt.Sprintf("\nLook: %v (drag: %.2f, %.2f)", g.lookEnabled, dragX, dragY)
	status += fmt.Sprintf("\nExpression: %s", g.model.GetCurrentExpression())
	status += "\n\nControls:"
	status += "\n  0-9: Switch/stop expressions"
	status += "\n  L: Toggle look tracking"
	status += "\n  B: Enable breath"
	status += "\n  P: Disable physics  R: Reset physics"
	status += "\n  Click hit areas: TapBody motion"
	ebitenutil.DebugPrint(screen, status)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.ow, g.oh = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	//modelPath := flag.String("model", "example/Resources/Haru/Haru.model3.json", "path to .model3.json file")
	modelPath := flag.String("model", "example/Resources/河原木桃香 · 标准模式/cat.model3.json", "path to .model3.json file")
	libPath := flag.String("lib", "example/Live2DCubismCore.dll", "path to Cubism Core library")
	flag.Parse()

	fmt.Printf("Loading model: %s\n", *modelPath)

	// Initialize the Cubism engine with the Core library
	csm, err := cubism.NewCubism(*libPath)
	if err != nil {
		log.Fatal(err)
	}

	// Set sound loader for motion audio playback
	csm.LoadSound = normal.LoadSound

	// Load the model from model3.json
	model, err := csm.LoadModel(*modelPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Model loaded successfully. Drawables: %d, Parameters: %d\n", len(model.GetDrawables()), len(model.GetParameters()))

	// Play idle motion (if available)
	motions := model.GetMotionGroups()
	if len(motions) > 0 {
		// Try "Idle" first, then fall back to the first available group
		groupName := ""
		for _, g := range motions {
			if g == "Idle" {
				groupName = g
				break
			}
		}
		if groupName == "" {
			groupName = motions[0]
		}
		model.PlayMotion(groupName, 0, true)
		fmt.Printf("Playing motion group: %s\n", groupName)
	} else {
		fmt.Println("No motions available for this model")
	}

	// Enable auto blink
	model.EnableAutoBlink()

	// Enable look/eye-tracking with standard parameters
	model.EnableLook([]look.LookParameterData{
		{ParameterId: "ParamAngleX", FactorX: 30.0, FactorY: 0.0, FactorXY: 0.0},
		{ParameterId: "ParamAngleY", FactorX: 0.0, FactorY: 30.0, FactorXY: 0.0},
		{ParameterId: "ParamAngleZ", FactorX: 0.0, FactorY: 0.0, FactorXY: 10.0},
		{ParameterId: "ParamBodyAngleX", FactorX: 10.0, FactorY: 0.0, FactorXY: 0.0},
		{ParameterId: "ParamEyeBallX", FactorX: 1.0, FactorY: 0.0, FactorXY: 0.0},
		{ParameterId: "ParamEyeBallY", FactorX: 0.0, FactorY: 1.0, FactorXY: 0.0},
	})

	// Enable breathing effect with default parameters
	model.EnableBreathWithParameters(breath.DefaultBreathParameters())

	// Print available expressions and play the first one
	expressions := model.GetExpressionNames()
	fmt.Printf("Available expressions (%d): %v\n", len(expressions), expressions)
	if len(expressions) > 0 {
		model.PlayExpression(expressions[0])
		fmt.Printf("Playing expression: %s\n", expressions[0])
	}

	// Print physics info
	gravityX, gravityY, windX, windY := model.GetPhysicsOptions()
	fmt.Printf("Physics: gravity=(%.1f, %.1f), wind=(%.1f, %.1f)\n", gravityX, gravityY, windX, windY)

	// Create renderer
	r, err := renderer.NewRenderer(model)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Renderer created, starting game loop...")

	g := &Game{
		renderer:    r,
		model:       model,
		expressions: expressions,
		lookEnabled: true,
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
