package internal

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Data struct {
	Engine Engine
	Window Window
}

var State Data
var GameLoop func()

func (d *Data) Update() error {
	if d.Engine.Exiting {
		return ebiten.Termination
	}
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
