package graphics

import (
	"big-black-box/geometry"
	"big-black-box/internal"
	"image/color"
)

func (c Camera) DrawShape(shape geometry.Shape) {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		Shape: internal.Shape(shape),
		Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	})
}
