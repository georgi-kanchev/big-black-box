package internal

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Area struct{ X, Y, Width, Height float32 }
type Shape struct{ X, Y, Width, Height, Angle, Roundness float32 }

type View struct {
	X, Y, Zoom, Angle float32

	WindowArea Area // The draw area in window space. Zero value = entire window.
	MaskArea   Area // In view space. Everything drawn outside of it is cropped. Zero value = no masking.
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
	Engine Engine
	Window Window
	Views  [8]View
}

var State Data
var GameLoop func()

var Fonts []*ebiten.Image = make([]*ebiten.Image, 0, 64)
var Images []*ebiten.Image = make([]*ebiten.Image, 0, 64)

var White1x1 = ebiten.NewImage(1, 1)

//=================================================================

func Init(gameLoop func()) {
	GameLoop = gameLoop
	SinCosCache()
	ShaderCache()
	White1x1.Set(0, 0, color.White)

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
