package internal

import (
	"cmp"
	"fmt"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Layer byte

type DrawItem struct {
	Shape          Shape
	TextureID      int
	Color          color.RGBA
	OutlineColor   [4]float32
	OutlineSize, Z float32
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
	scrSize := screen.Bounds().Size()
	for i := range DrawQueues[0] {
		drawShape(screen, scrSize, DrawQueues[0][i].Shape, DrawQueues[0][i].Color, DrawQueues[0][i].OutlineColor, DrawQueues[0][i].OutlineSize)
	}
	for i := range DrawQueues[1] {
		drawShape(screen, scrSize, DrawQueues[1][i].Shape, DrawQueues[1][i].Color, DrawQueues[1][i].OutlineColor, DrawQueues[1][i].OutlineSize)
	}
	for i := range DrawQueues[2] {
		drawShape(screen, scrSize, DrawQueues[2][i].Shape, DrawQueues[2][i].Color, DrawQueues[2][i].OutlineColor, DrawQueues[2][i].OutlineSize)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()), 8, 8)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.0f", ebiten.ActualTPS()), 8, 20)
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
