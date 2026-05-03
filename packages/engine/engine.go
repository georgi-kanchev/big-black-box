package engine

import (
	"big-black-box/packages/internal"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func Run(targetTPS int, gameLoop func()) {
	internal.Init(gameLoop)

	ebiten.SetTPS(max(targetTPS, 1))
	var err = ebiten.RunGame(&internal.Engine)
	if err != nil {
		log.Fatal(err)
	}
}
func Quit() {
	internal.IsExiting = true
}

// Ticks per second, provided in:
//
//	engine.Run(targetTPS, gameLoop)
func TPS() int {
	return ebiten.TPS()
}
