package graphics

import (
	"big-black-box/internal"
	"big-black-box/utility/angle"
	"big-black-box/utility/number"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera byte

func NewCamera() Camera {
	var slotIndex = -1
	for i, cam := range internal.State.Cameras {
		if i == 0 { // 0 is for invalid
			continue
		}
		if cam == (internal.Camera{}) {
			slotIndex = i
			break
		}
	}

	if slotIndex == -1 {
		return 0 // invalid
	}
	internal.State.Cameras[slotIndex] = internal.Camera{Zoom: 1}
	return Camera(slotIndex)
}

//=================================================================

func (c Camera) SetPosition(x, y float32) {
	if c.isValid() {
		internal.State.Cameras[c].X = x
		internal.State.Cameras[c].Y = y
	}
}
func (c Camera) SetZoom(zoom float32) {
	if c.isValid() {
		internal.State.Cameras[c].Zoom = zoom
	}
}
func (c Camera) SetAngle(angle float32) {
	if c.isValid() {
		internal.State.Cameras[c].Angle = angle
	}
}
func (c Camera) SetWindowArea(x, y, width, height float32) {
	if c.isValid() {
		internal.State.Cameras[c].WindowArea = internal.Area{X: x, Y: y, Width: width, Height: height}
	}
}
func (c Camera) SetMaskArea(x, y, width, height float32) {
	if c.isValid() {
		internal.State.Cameras[c].MaskArea = internal.Area{X: x, Y: y, Width: width, Height: height}
	}
}

//=================================================================

func (c Camera) Position() (x, y float32) {
	if c.isValid() {
		return internal.State.Cameras[c].X, internal.State.Cameras[c].Y
	}
	return number.NaN(), number.NaN()
}
func (c Camera) Zoom() float32 {
	if c.isValid() {
		return internal.State.Cameras[c].Zoom
	}
	return number.NaN()
}
func (c Camera) Angle() float32 {
	if c.isValid() {
		return internal.State.Cameras[c].Angle
	}
	return number.NaN()
}
func (c Camera) WindowArea() (x, y, width, height float32) {
	if c.isValid() {
		var area = internal.State.Cameras[c].WindowArea
		if area == (internal.Area{}) {
			var ww, wh = ebiten.WindowSize()
			return 0, 0, float32(ww), float32(wh)
		}
		return area.X, area.Y, area.Width, area.Height
	}
	return number.NaN(), number.NaN(), number.NaN(), number.NaN()
}
func (c Camera) MaskArea() (x, y, width, height float32) {
	if c.isValid() {
		var area = internal.State.Cameras[c].MaskArea
		if area == (internal.Area{}) {
			var ww, wh = ebiten.WindowSize()
			return 0, 0, float32(ww), float32(wh)
		}
		return area.X, area.Y, area.Width, area.Height
	}
	return number.NaN(), number.NaN(), number.NaN(), number.NaN()
}
func (c Camera) Size() (width, height float32) {
	var _, _, sw, sh = c.WindowArea()
	var zoom = c.Zoom()
	return sw / zoom, sh / zoom
}

// private ========================================================

func (c Camera) isValid() bool {
	return c > 0 && int(c) < len(internal.State.Cameras)
}

func (c Camera) matrix() ebiten.GeoM {
	if !c.isValid() {
		return ebiten.GeoM{}
	}

	var data = internal.State.Cameras[c] // copy the struct from the array once to the stack

	var m ebiten.GeoM
	m.Translate(float64(-data.X), float64(-data.Y))

	if data.Angle != 0 {
		m.Rotate(float64(angle.ToRadians(data.Angle)))
	}

	m.Scale(float64(data.Zoom), float64(data.Zoom))

	var _, _, ww, wh = c.WindowArea()
	m.Translate(float64(ww/2), float64(wh/2))
	return m
}
