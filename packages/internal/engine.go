// lifetime, updates and events of the engine and game

package internal

import (
	"big-black-box/packages/window"

	"github.com/hajimehoshi/ebiten/v2"
)

type Events struct{}

var IsExiting bool
var TargetTPS int
var Engine Events
var GameLoop func()

//=================================================================

func Init(gameLoop func()) {
	GameLoop = gameLoop
	sinCosCache()
	initAssets()

	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	readWindowValues()
}
func (e Events) Update() error {
	if IsExiting {
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
