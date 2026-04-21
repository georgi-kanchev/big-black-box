package internal

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Area struct{ X, Y, Width, Height float32 }
type Shape struct{ X, Y, Width, Height, Angle, Roundness float32 }

type Camera struct {
	X, Y, Zoom, Angle float32

	WindowArea Area // The draw area in window space. Zero value = entire window.
	MaskArea   Area // In camera space. Everything drawn outside of it is cropped. Zero value = no masking.
}

type Window struct {
	Title      string
	Mode       byte // see window.Mode
	Monitor    byte
	PixelScale float32

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

func Init(gameLoop func()) {
	GameLoop = gameLoop
	SinCosCache()
	ShaderCache()

	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func (d Data) Update() error {
	if d.Engine.Exiting {
		return ebiten.Termination
	}
	cacheInput()
	cacheTime()

	d.DrawStart()
	GameLoop()
	d.DrawEnd()
	return nil
}

func (d Data) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float32(outsideWidth) / d.Window.PixelScale), int(float32(outsideHeight) / d.Window.PixelScale)
}

func drawShape(screen *ebiten.Image, scrSize image.Point, s Shape, fill color.RGBA, outlineColor [4]float32, outlineSize float32) {
	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Center":           []float32{s.X, s.Y},
		"Size":             []float32{s.Width, s.Height},
		"Roundness":        s.Roundness,
		"Rotation":         s.Angle * (3.14159265 / 180.0),
		"OutlineThickness": outlineSize,
		"OutlineColor":     outlineColor,
	}
	op.ColorScale.ScaleWithColor(fill)
	screen.DrawRectShader(scrSize.X, scrSize.Y, shader, op)
}
