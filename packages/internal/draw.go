// keeps track of a draw call queue that gets appended by the game in the Update phase,
// it gets sorted afterwards and drawn during ebiten's Draw phase

package internal

import (
	"cmp"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Layer byte
type Kind byte

type Area struct{ X, Y, Width, Height float32 }
type Shape struct{ X, Y, Width, Height, Angle, Roundness float32 }

type View struct {
	X, Y, Zoom, Angle float32

	WindowArea Area // The draw area in window space. Zero value = entire window.
	MaskArea   Area // In view space. Everything drawn outside of it is cropped. Zero value = no masking.
}

type DrawItem struct {
	Kind  Kind
	Shape Shape
	Image Image
	Font  Font

	Color, OutlineColor uint
	OutlineSize, Z      float32
	Text                string

	View View

	TexX, TexY, TexWidth, TexHeight float32

	Verts     [8]ebiten.Vertex
	VertCount int
}

const KindNone, KindShape, KindImage, KindText, KindTilemap = 0, 1, 2, 3, 4

const LayerBelow, LayerDefault, LayerAbove, LayerCount Layer = 0, 1, 2, 3

var DrawQueues [LayerCount][]DrawItem
var DrawCounts [LayerCount]int

func BeforeGameLoop() {
	for i := range DrawCounts {
		DrawCounts[i] = 0
	}
}
func Queue(layer Layer, item DrawItem) {
	if DrawCounts[layer] < len(DrawQueues[layer]) {
		DrawQueues[layer][DrawCounts[layer]] = item
	} else {
		DrawQueues[layer] = append(DrawQueues[layer], item)
	}
	DrawCounts[layer]++
}
func AfterGameLoop() {
	var below = DrawQueues[LayerBelow][:DrawCounts[LayerBelow]]
	slices.SortStableFunc(below, func(a, b DrawItem) int {
		return cmp.Compare(a.Z, b.Z)
	})
	var above = DrawQueues[LayerAbove][:DrawCounts[LayerAbove]]
	slices.SortStableFunc(above, func(a, b DrawItem) int {
		return cmp.Compare(a.Z, b.Z)
	})
}
func (e Events) Draw(screen *ebiten.Image) {
	for i := range LayerCount {
		drawLayer(screen, DrawQueues[i][:DrawCounts[i]])
	}
}

// private ========================================================

var op = &ebiten.DrawTrianglesShaderOptions{Uniforms: map[string]any{}}
var vertices = make([]ebiten.Vertex, 0, 1024*4)
var indices = make([]uint16, 0, 1024*6)

func drawLayer(screen *ebiten.Image, items []DrawItem) {
	vertices = vertices[:0]
	indices = indices[:0]
	var currentImage *ebiten.Image

	for _, item := range items {
		if item.Kind == KindText && item.Font == 0 {
			draw(screen, currentImage)
			ebitenutil.DebugPrintAt(screen, item.Text, int(item.Shape.X), int(item.Shape.Y))
			continue
		}

		var img *ebiten.Image
		switch item.Kind {
		case KindImage:
			img = Images[item.Image-1]
		case KindText:
			img = Fonts[item.Font-1]
		default:
			img = White1x1
		}

		if currentImage != nil && img != currentImage {
			draw(screen, currentImage)
		}
		currentImage = img

		if item.VertCount < 3 {
			continue
		}
		var base = uint16(len(vertices))
		vertices = append(vertices, item.Verts[:item.VertCount]...)
		for i := 1; i < item.VertCount-1; i++ {
			indices = append(indices, base, base+uint16(i), base+uint16(i+1))
		}
	}
	draw(screen, currentImage)
}

func draw(screen, currentImage *ebiten.Image) {
	if len(vertices) == 0 {
		return
	}
	op.Images[0] = currentImage
	op.Uniforms["Smoothing"] = 1.0
	screen.DrawTrianglesShader(vertices, indices, shader, op)
	vertices = vertices[:0]
	indices = indices[:0]
}
