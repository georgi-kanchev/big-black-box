package graphics

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/angle"

	"github.com/hajimehoshi/ebiten/v2"
)

type View internal.View

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

// private ========================================================

func (v View) matrix() ebiten.GeoM {
	var m ebiten.GeoM
	m.Translate(float64(-v.X), float64(-v.Y))
	if v.Angle != 0 {
		m.Rotate(float64(angle.ToRadians(v.Angle)))
	}
	m.Scale(float64(v.Zoom), float64(v.Zoom))
	m.Translate(float64(v.WindowArea.Width/2), float64(v.WindowArea.Height/2))
	return m
}
