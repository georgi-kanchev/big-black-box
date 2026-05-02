package engine

import (
	"big-black-box/packages/internal"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func Run(targetTPS uint, gameLoop func()) {
	internal.Init(gameLoop)

	ebiten.SetTPS(int(targetTPS))
	var err = ebiten.RunGame(&internal.State)
	if err != nil {
		log.Fatal(err)
	}
}
func Quit() {
	internal.Exiting = true
}
