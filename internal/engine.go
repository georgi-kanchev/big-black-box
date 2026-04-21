package internal

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DrawItem struct {
	Shape        Shape
	TextureID    int
	Color        color.RGBA
	OutlineColor [4]float32
	OutlineSize  float32
}

type Area struct{ X, Y, Width, Height float32 }
type Shape struct{ X, Y, Width, Height, Angle, Roundness float32 }

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

	PixelScale     float32
	TargetTickRate int
}

type Data struct {
	Engine  Engine
	Window  Window
	Cameras [8]Camera
}

var State Data
var GameLoop func()

var DrawQueue = make([]DrawItem, 0, 1024)
var DrawCount = 0

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

	DrawCount = 0
	GameLoop()
	return nil
}

func (d Data) Draw(screen *ebiten.Image) {
	scrSize := screen.Bounds().Size()
	for i := range DrawCount {
		drawShape(screen, scrSize, DrawQueue[i].Shape, DrawQueue[i].Color, DrawQueue[i].OutlineColor, DrawQueue[i].OutlineSize)
	}

	// // 1. Draw the actual shapes
	// drawShape(screen, scrSize, shapeA, color.RGBA{100, 180, 255, 255}, []float32{1, 1, 1, 1}, 10)
	// drawShape(screen, scrSize, shapeB, color.RGBA{255, 255, 255, 200}, []float32{1, 1, 1, 1}, 10)

	// // 2. Draw the Bounds of ShapeA
	// x1, y1, x2, y2 := shapeA.Bounds()
	// bw, bh := x2-x1, y2-y1
	// // A thin cyan rectangle to show the AABB
	// vector.StrokeRect(screen, x1, y1, bw, bh, 3, color.RGBA{0, 255, 255, 100}, true)

	// // 4. Debug Prints
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()), 8, 8)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.0f", ebiten.ActualTPS()), 8, 20)
}

func (d Data) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float32(outsideWidth) / d.Engine.PixelScale), int(float32(outsideHeight) / d.Engine.PixelScale)
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
