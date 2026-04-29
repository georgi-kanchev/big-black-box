package internal

import (
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/color"
	"cmp"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Font byte
type Image uint16

type Layer byte
type Kind byte

const KindNone, KindShape, KindImage, KindText = 0, 1, 2, 3

type DrawItem struct {
	Kind  Kind
	Shape Shape
	Image Image
	Font  Font

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

var op = &ebiten.DrawTrianglesShaderOptions{Uniforms: map[string]any{}}
var vertices = make([]ebiten.Vertex, 4)
var indices = []uint16{0, 1, 2, 1, 2, 3}

func drawItem(screen *ebiten.Image, item DrawItem) {
	var br, bg, bb, ba = color.Channels(item.Color)
	var r, g, b, a = float32(br) / 255, float32(bg) / 255, float32(bb) / 255, float32(ba) / 255
	var scrSize = screen.Bounds().Size()
	var or, og, ob, oa = color.Channels(item.OutlineColor)

	for i := range vertices {
		vertices[i].ColorR, vertices[i].ColorG, vertices[i].ColorB, vertices[i].ColorA = r, g, b, a
	}
	vertices[1].DstX = float32(scrSize.X)
	vertices[2].DstY = float32(scrSize.Y)
	vertices[3].DstX = float32(scrSize.X)
	vertices[3].DstY = float32(scrSize.Y)

	op.Uniforms["CenterX"] = item.Shape.X
	op.Uniforms["CenterY"] = item.Shape.Y
	op.Uniforms["Width"] = item.Shape.Width
	op.Uniforms["Height"] = item.Shape.Height
	op.Uniforms["Rotation"] = angle.ToRadians(item.Shape.Angle)
	op.Uniforms["Roundness"] = max(min(item.Shape.Roundness, 1), 0)
	op.Uniforms["OutlineSize"] = item.OutlineSize
	op.Uniforms["OutlineColorR"] = float32(or) / 255
	op.Uniforms["OutlineColorG"] = float32(og) / 255
	op.Uniforms["OutlineColorB"] = float32(ob) / 255
	op.Uniforms["OutlineColorA"] = float32(oa) / 255

	switch item.Kind {
	case KindShape:
		op.Images[0] = White1x1

	case KindImage:
		op.Images[0] = Images[item.Image-1]
	case KindText:
		if item.Font == 0 {
			ebitenutil.DebugPrintAt(screen, item.Text, int(item.Shape.X), int(item.Shape.Y))
		} else {
			// screen.DrawImage(Fonts[item.Font-1], nil)
		}
	}
	screen.DrawTrianglesShader(vertices, indices, shader, op)
}
