package internal

import (
	"big-black-box/geometry"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Area struct{ X, Y, Width, Height float32 }

type Camera struct {
	X, Y, Zoom, Angle float32

	WindowArea Area // The draw area in window space. Zero value = entire window.
	MaskArea   Area // In camera space. Everything drawn outside of it is cropped. Zero value = no masking.
}

type Window struct {
	Title   string
	Mode    byte // see window.Mode
	Monitor byte

	IsMaximized, IsVsynced bool
}

type Engine struct {
	Exiting bool

	TargetTickRate int
}

type Data struct {
	Engine  Engine
	Window  Window
	Cameras [8]Camera
}

var State Data
var GameLoop func()

//=================================================================

var shapeA = geometry.Shape{X: 400, Y: 400, Width: 250, Height: 250, Angle: 20, Roundness: 1}
var shapeB = geometry.Shape{Width: 120, Height: 120, Roundness: 1}

func Init(gameLoop func()) {
	GameLoop = gameLoop
	SinCosCache()
	ShaderCache()

	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func (d *Data) Update() error {
	if d.Engine.Exiting {
		return ebiten.Termination
	}
	cacheInput()
	cacheTime()

	const speed = 5.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		shapeB.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		shapeB.X += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		shapeB.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		shapeB.Y += speed
	}

	if px, py, hit := shapeA.Collide(&shapeB); hit {
		shapeB.X += px
		shapeB.Y += py
	}

	GameLoop()
	return nil
}

func (d *Data) Draw(screen *ebiten.Image) {
	scrSize := screen.Bounds().Size()

	drawShape(screen, scrSize, shapeA, color.RGBA{100, 180, 255, 255}, []float32{1, 1, 1, 1}, 10)
	drawShape(screen, scrSize, shapeB, color.RGBA{255, 255, 255, 200}, []float32{1, 1, 1, 1}, 10)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()), 8, 8)
}

func (d *Data) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func drawShape(screen *ebiten.Image, scrSize image.Point, s geometry.Shape, fill color.RGBA, outlineColor []float32, outlineThickness float32) {
	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Center":           []float32{s.X, s.Y},
		"Size":             []float32{s.Width, s.Height},
		"Roundness":        s.Roundness,
		"Rotation":         s.Angle * (3.14159265 / 180.0),
		"OutlineThickness": outlineThickness,
		"OutlineColor":     outlineColor,
	}
	op.ColorScale.ScaleWithColor(fill)
	screen.DrawRectShader(scrSize.X, scrSize.Y, shader, op)
}
