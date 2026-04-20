package internal

import (
	"big-black-box/geometry"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

var shapeA = func() geometry.Shape {
	var s = geometry.Shape{X: 400, Y: 400, Width: 450, Height: 250}
	s.SetAR(20, 0.5)
	return s
}()
var shapeB = func() geometry.Shape {
	var s = geometry.Shape{Width: 120, Height: 120}
	s.SetAR(0, 1)
	return s
}()

func Init(gameLoop func()) {
	GameLoop = gameLoop
	SinCosCache()
	ShaderCache()

	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

var clX, clY float32 // Closest point coordinates
var mx, my float32   // Mouse coordinates

func (d *Data) Update() error {
	if d.Engine.Exiting {
		return ebiten.Termination
	}
	cacheInput()
	cacheTime()

	// 1. Handle Movement for Shape B
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

	// 2. Handle Rotation for Shape A
	currentAngle := shapeA.Angle()
	currentRoundness := shapeA.Roundness()
	// Increment angle (e.g., 1 degree per frame)
	newAngle := currentAngle + 1.0
	shapeA.SetAR(newAngle, currentRoundness)

	// 3. Closest Point & Collision Logic
	imx, imy := ebiten.CursorPosition()
	mx, my = float32(imx), float32(imy)
	clX, clY = shapeA.ClosestPointOnEdge(mx, my)

	shapeB = shapeA.Collide(shapeB)

	GameLoop()
	return nil
}

func (d *Data) Draw(screen *ebiten.Image) {
	scrSize := screen.Bounds().Size()

	// 1. Draw the actual shapes
	drawShape(screen, scrSize, shapeA, color.RGBA{100, 180, 255, 255}, []float32{1, 1, 1, 1}, 10)
	drawShape(screen, scrSize, shapeB, color.RGBA{255, 255, 255, 200}, []float32{1, 1, 1, 1}, 10)

	// 2. Draw the Bounds of ShapeA
	x1, y1, x2, y2 := shapeA.Bounds()
	bw, bh := x2-x1, y2-y1
	// A thin cyan rectangle to show the AABB
	vector.StrokeRect(screen, x1, y1, bw, bh, 3, color.RGBA{0, 255, 255, 100}, true)

	// 3. Draw the Closest Point visualization
	vector.StrokeLine(screen, mx, my, clX, clY, 12, color.RGBA{255, 255, 0, 150}, true)
	vector.DrawFilledCircle(screen, clX, clY, 6, color.RGBA{255, 255, 0, 255}, true)

	// 4. Debug Prints
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()), 8, 8)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bounds: (%.1f, %.1f) to (%.1f, %.1f)", x1, y1, x2, y2), 8, 24)
}

func (d *Data) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func drawShape(screen *ebiten.Image, scrSize image.Point, s geometry.Shape, fill color.RGBA, outlineColor []float32, outlineThickness float32) {
	var angle, roundness = s.Angle(), s.Roundness()
	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Center":           []float32{s.X, s.Y},
		"Size":             []float32{s.Width, s.Height},
		"Roundness":        roundness,
		"Rotation":         angle * (3.14159265 / 180.0),
		"OutlineThickness": outlineThickness,
		"OutlineColor":     outlineColor,
	}
	op.ColorScale.ScaleWithColor(fill)
	screen.DrawRectShader(scrSize.X, scrSize.Y, shader, op)
}
