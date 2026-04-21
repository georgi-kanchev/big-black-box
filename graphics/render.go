package graphics

import (
	"big-black-box/geometry"
	"big-black-box/internal"
	"image/color"
)

func (c Camera) RenderShape(shape geometry.Shape) {
	if len(internal.DrawQueue) == internal.DrawCount {
		internal.DrawQueue = append(internal.DrawQueue, internal.DrawItem{})
	}
	internal.DrawQueue[internal.DrawCount].Shape = internal.Shape(shape)
	internal.DrawQueue[internal.DrawCount].Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	internal.DrawCount++
}
