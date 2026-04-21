package engine

import (
	"big-black-box/internal"
	"big-black-box/window"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func Run(gameLoop func()) {
	internal.Init(gameLoop)

	window.SetMonitor(0)
	window.SetTitle("game")
	window.SetVsync(false)
	window.SetPixelScale(2)
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
