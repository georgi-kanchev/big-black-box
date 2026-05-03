package graphics

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/angle"

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

// private ========================================================

// worldToView transforms world coordinates to view space:
// centered at (0,0) for the camera's focal point, rotated and scaled by Zoom.
func (v View) worldToView() ebiten.GeoM {
	var m ebiten.GeoM
	m.Translate(float64(-v.X), float64(-v.Y))
	if v.Angle != 0 {
		m.Rotate(float64(angle.ToRadians(v.Angle)))
	}
	m.Scale(float64(v.Zoom), float64(v.Zoom))
	return m
}

// viewToWindow transforms view-space coordinates to window/screen space.
func (v View) viewToWindow() ebiten.GeoM {
	var areaX, areaY, areaW, areaH float32
	if v.WindowArea != (internal.Area{}) {
		areaX, areaY = v.WindowArea.X, v.WindowArea.Y
		areaW, areaH = v.WindowArea.Width, v.WindowArea.Height
	} else {
		var w, h = ebiten.WindowSize()
		areaW, areaH = float32(w), float32(h)
	}
	var m ebiten.GeoM
	m.Translate(float64(areaX+areaW/2), float64(areaY+areaH/2))
	return m
}

func (v View) matrix() ebiten.GeoM {
	var m = v.worldToView()
	m.Concat(v.viewToWindow())
	return m
}
