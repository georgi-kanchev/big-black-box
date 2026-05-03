package graphics

import (
	"big-black-box/packages/internal"

	"github.com/hajimehoshi/ebiten/v2"
)

type View internal.View
type Area = internal.Area

func NewView() View {
	return View{Zoom: 1}
}

func (v View) Size() (width, height float32) {
	var w, h = ebiten.WindowSize()
	if v.WindowArea != (internal.Area{}) {
		w, h = int(v.WindowArea.Width), int(v.WindowArea.Height)
	}
	return float32(w) / v.Zoom, float32(h) / v.Zoom
}
