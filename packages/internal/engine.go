package internal

import (
	"big-black-box/packages/window"
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

type Events struct{}

var Exiting bool
var TargetTPS int
var State Events
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
func (e Events) Update() error {
	if Exiting {
		return ebiten.Termination
	}
	readWindowValues()
	cacheInput()
	cacheTime()

	BeforeGameLoop()
	GameLoop()
	AfterGameLoop()
	return nil
}
func (e Events) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float32(outsideWidth) / window.PixelScale), int(float32(outsideHeight) / window.PixelScale)
}
