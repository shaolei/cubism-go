package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	"github.com/shaolei/cubism-go"
	renderer "github.com/shaolei/cubism-go/renderer/ebitengine"
	"github.com/shaolei/cubism-go/sound/normal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
}

func (g *Game) Update() (err error) {
	g.renderer.Update()
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
			g.renderer.GetModel().StopMotion(g.tapId)
			g.tapId = g.renderer.GetModel().PlayMotion("TapBody", 0, false)
		}
	} else if ebiten.CursorShape() == ebiten.CursorShapePointer {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	// Expression switching with number keys (1-9 = expression index, 0 = stop)
	keys := []ebiten.Key{ebiten.Key0, ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5, ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9}
	for i, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			m := g.renderer.GetModel()
			if i == 0 {
				m.StopExpression()
				fmt.Println("Expression stopped")
			} else if i <= len(g.expressions) {
				name := g.expressions[i-1]
				m.PlayExpression(name)
				fmt.Printf("Playing expression [%d]: %s\n", i, name)
			}
		}
	}
	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.renderer.Draw(screen, renderer.WithBackground(color.RGBA{0, 255, 0, 255}))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.ow, g.oh = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	modelPath := flag.String("model", "Resources/Haru/Haru.model3.json", "path to .model3.json file")
	libPath := flag.String("lib", "Live2DCubismCore.dll", "path to Cubism Core library")
	flag.Parse()

	fmt.Printf("Loading model: %s\n", *modelPath)
	csm, err := cubism.NewCubism(*libPath)
	if err != nil {
		log.Fatal(err)
	}
	// Set function for playing sound
	csm.LoadSound = normal.LoadSound
	model, err := csm.LoadModel(*modelPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Model loaded successfully. Drawables: %d, Parameters: %d\n", len(model.GetDrawables()), len(model.GetParameters()))
	// Set parameters for model-specific features
	model.SetParameterValue("ParamMouseLeftDown", 1.0) // 显示左手
	model.SetParameterValue("Param7", 1.0)             // 切换河原木桃香发型
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
	model.EnableAutoBlink()
	// Print available expressions and play the first one
	expressions := model.GetExpressionNames()
	fmt.Printf("Available expressions (%d): %v\n", len(expressions), expressions)
	if len(expressions) > 0 {
		model.PlayExpression(expressions[0])
		fmt.Printf("Playing expression: %s\n", expressions[0])
	}
	renderer, err := renderer.NewRenderer(model)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Renderer created, starting game loop...")
	g := &Game{
		renderer:    renderer,
		expressions: expressions,
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
