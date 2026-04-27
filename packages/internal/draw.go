package internal

import (
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/color"
	"cmp"
	"slices"

	col "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Font byte
type Texture uint16

type Layer byte
type Kind byte

const KindNone, KindShape, KindSprite, KindText = 0, 1, 2, 3

type DrawItem struct {
	Kind                Kind
	Shape               Shape
	Texture             Texture
	Font                Font
	Color, OutlineColor uint
	OutlineSize, Z      float32
	Text                string
}

const LayerBelow, LayerDefault, LayerAbove Layer = 0, 1, 2

var DrawQueues [3][]DrawItem
var DrawCounts [3]int

func Queue(layer Layer, item DrawItem) {
	if DrawCounts[layer] < len(DrawQueues[layer]) {
		DrawQueues[layer][DrawCounts[layer]] = item
	} else {
		DrawQueues[layer] = append(DrawQueues[layer], item)
	}
	DrawCounts[layer]++
}

func (d Data) DrawStart() {
	for i := range DrawCounts {
		DrawCounts[i] = 0
	}
}
func (d Data) Draw(screen *ebiten.Image) {
	for i := range DrawCounts[0] {
		drawItem(screen, DrawQueues[0][i])
	}
	for i := range DrawCounts[1] {
		drawItem(screen, DrawQueues[1][i])
	}
	for i := range DrawCounts[2] {
		drawItem(screen, DrawQueues[2][i])
	}
}
func (d Data) DrawEnd() {
	var below = DrawQueues[LayerBelow][:DrawCounts[LayerBelow]]
	slices.SortStableFunc(below, func(a, b DrawItem) int {
		return cmp.Compare(a.Z, b.Z)
	})
	var above = DrawQueues[LayerAbove][:DrawCounts[LayerAbove]]
	slices.SortStableFunc(above, func(a, b DrawItem) int {
		return cmp.Compare(a.Z, b.Z)
	})
}

// private ========================================================

var op = &ebiten.DrawRectShaderOptions{Uniforms: map[string]any{}}

func drawItem(screen *ebiten.Image, item DrawItem) {
	var r, g, b, a = color.Channels(item.Color)

	switch item.Kind {
	case KindShape:
		var scrSize = screen.Bounds().Size()
		var or, og, ob, oa = color.Channels(item.OutlineColor)
		op.Uniforms["CenterX"] = item.Shape.X
		op.Uniforms["CenterY"] = item.Shape.Y
		op.Uniforms["Width"] = item.Shape.Width
		op.Uniforms["Height"] = item.Shape.Height
		op.Uniforms["Roundness"] = item.Shape.Roundness
		op.Uniforms["Rotation"] = angle.ToRadians(item.Shape.Angle)
		op.Uniforms["OutlineSize"] = item.OutlineSize
		op.Uniforms["OutlineColorR"] = float32(or) / 255
		op.Uniforms["OutlineColorG"] = float32(og) / 255
		op.Uniforms["OutlineColorB"] = float32(ob) / 255
		op.Uniforms["OutlineColorA"] = float32(oa) / 255
		op.ColorScale.ScaleWithColor(col.RGBA{R: r, G: g, B: b, A: a})

		screen.DrawRectShader(scrSize.X, scrSize.Y, shader, op)
	case KindSprite:
	case KindText:
		if item.Font == 0 {
			ebitenutil.DebugPrintAt(screen, item.Text, int(item.Shape.X), int(item.Shape.Y))
		} else {
			screen.DrawImage(Fonts[item.Font-1], nil)
		}
	}
}
