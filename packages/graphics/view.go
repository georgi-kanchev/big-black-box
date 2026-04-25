package graphics

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/number"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type View byte

func NewView() View {
	var slotIndex = -1
	for i, cam := range internal.State.Views {
		if i == 0 { // 0 is for invalid
			continue
		}
		if cam == (internal.View{}) {
			slotIndex = i
			break
		}
	}

	if slotIndex == -1 {
		log.Println("Cannot create more than 7 Cameras!")
		return 0 // invalid
	}
	internal.State.Views[slotIndex] = internal.View{Zoom: 1}
	return View(slotIndex)
}

//=================================================================

func (v View) SetPosition(x, y float32) {
	if v.isValid() {
		internal.State.Views[v].X = x
		internal.State.Views[v].Y = y
	}
}
func (v View) SetZoom(zoom float32) {
	if v.isValid() {
		internal.State.Views[v].Zoom = zoom
	}
}
func (v View) SetAngle(angle float32) {
	if v.isValid() {
		internal.State.Views[v].Angle = angle
	}
}
func (v View) SetWindowArea(x, y, width, height float32) {
	if v.isValid() {
		internal.State.Views[v].WindowArea = internal.Area{X: x, Y: y, Width: width, Height: height}
	}
}
func (v View) SetMaskArea(x, y, width, height float32) {
	if v.isValid() {
		internal.State.Views[v].MaskArea = internal.Area{X: x, Y: y, Width: width, Height: height}
	}
}

//=================================================================

func (v View) Position() (x, y float32) {
	if v.isValid() {
		return internal.State.Views[v].X, internal.State.Views[v].Y
	}
	return number.NaN(), number.NaN()
}
func (v View) Zoom() float32 {
	if v.isValid() {
		return internal.State.Views[v].Zoom
	}
	return number.NaN()
}
func (v View) Angle() float32 {
	if v.isValid() {
		return internal.State.Views[v].Angle
	}
	return number.NaN()
}
func (v View) WindowArea() (x, y, width, height float32) {
	if v.isValid() {
		var area = internal.State.Views[v].WindowArea
		if area == (internal.Area{}) {
			var ww, wh = ebiten.WindowSize()
			return 0, 0, float32(ww), float32(wh)
		}
		return area.X, area.Y, area.Width, area.Height
	}
	return number.NaN(), number.NaN(), number.NaN(), number.NaN()
}
func (v View) MaskArea() (x, y, width, height float32) {
	if v.isValid() {
		var area = internal.State.Views[v].MaskArea
		if area == (internal.Area{}) {
			var ww, wh = ebiten.WindowSize()
			return 0, 0, float32(ww), float32(wh)
		}
		return area.X, area.Y, area.Width, area.Height
	}
	return number.NaN(), number.NaN(), number.NaN(), number.NaN()
}
func (v View) Size() (width, height float32) {
	var _, _, sw, sh = v.WindowArea()
	var zoom = v.Zoom()
	return sw / zoom, sh / zoom
}

// private ========================================================

func (v View) isValid() bool {
	return v > 0 && int(v) < len(internal.State.Views)
}

func (v View) matrix() ebiten.GeoM {
	if !v.isValid() {
		return ebiten.GeoM{}
	}

	var data = internal.State.Views[v] // copy the struct from the array once to the stack

	var m ebiten.GeoM
	m.Translate(float64(-data.X), float64(-data.Y))

	if data.Angle != 0 {
		m.Rotate(float64(angle.ToRadians(data.Angle)))
	}

	m.Scale(float64(data.Zoom), float64(data.Zoom))

	var _, _, ww, wh = v.WindowArea()
	m.Translate(float64(ww/2), float64(wh/2))
	return m
}
