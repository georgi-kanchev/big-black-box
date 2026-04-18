package internal

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var Clock, Runtime, FPS, TPS, prevClock float32
var FrameCount uint64

func cacheTime() {
	var now = time.Now()
	var y, m, d = now.Date()
	var midnight = time.Date(y, m, d, 0, 0, 0, 0, now.Location())

	var secondsSinceMidnight = float32(now.Sub(midnight).Seconds())
	Clock = secondsSinceMidnight

	// handle midnight wrap-around
	if prevClock == 0 {
		prevClock = Clock
	}
	if prevClock > Clock {
		prevClock = Clock
	}

	FPS = float32(ebiten.ActualFPS())
	TPS = float32(ebiten.ActualTPS())

	Runtime += 1.0 / float32(ebiten.TPS())

	FrameCount = uint64(ebiten.Tick())
	prevClock = Clock
}
