package internal

import (
	"fmt"

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

func (d *Data) Update() error {
	if d.Engine.Exiting {
		return ebiten.Termination
	}
	cacheInput()
	cacheTime()

	GameLoop()
	return nil
}
func (d *Data) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %f", ebiten.ActualTPS()), 0, 16)
}
func (d *Data) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 4, outsideHeight / 4
}
