package engine

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/window"
	"big-black-box/packages/window/mode"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func Run(gameLoop func()) {
	internal.Init(gameLoop)

	window.SetMonitor(0)
	window.SetTitle("game")
	window.SetVsync(true)
	window.SetPixelScale(1)
	window.SetMode(mode.Maximized)
	SetTargetTickRate(60)

	var err = ebiten.RunGame(&internal.State)
	if err != nil {
		log.Fatal(err)
	}
}
func Quit() {
	internal.State.Engine.Exiting = true
}

//=================================================================

func GetTargetTickRate() int {
	return internal.State.Engine.TargetTickRate
}
func SetTargetTickRate(ticksPerSecond int) {
	internal.State.Engine.TargetTickRate = ticksPerSecond
	ebiten.SetTPS(ticksPerSecond)
}
