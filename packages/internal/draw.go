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

const KindNone, KindShape, KindImage, KindText, KindTilemap = 0, 1, 2, 3, 4

type DrawItem struct {
	Kind  Kind
	Shape Shape
	Image Image
	Font  Font

	Color, OutlineColor uint
	OutlineSize, Z      float32
	Text                string

	ImageX, ImageY, ImageWidth, ImageHeight float32
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

// packColor packs an RGBA uint into a 24-bit float (6 bits per channel).
// Layout: R6[23:18] G6[17:12] B6[11:6] A6[5:0].
func packColor(c uint) float32 {
	r, g, b, a := color.Channels(c)
	return float32(uint32(r>>2)<<18 | uint32(g>>2)<<12 | uint32(b>>2)<<6 | uint32(a>>2))
}

// packTextLayout packs a 3-bit align value (0–7) and a 1-bit wordwrap flag.
// Layout: WordWrap[3] Align[2:0].
func packTextLayout(align byte, wordWrap bool) float32 {
	v := uint32(align & 0x7)
	if wordWrap {
		v |= 1 << 3
	}
	return float32(v)
}

func drawItem(screen *ebiten.Image, item DrawItem) {
	const pivotX, pivotY float32 = 0.5, 0.5
	var pad float32 = 1.5
	if item.OutlineSize > 0 {
		pad += item.OutlineSize
	}

	var origW, origH = item.Shape.Width, item.Shape.Height
	var w, h = origW + (pad * 2), origH + (pad * 2)
	var x, y = item.Shape.X, item.Shape.Y

	var px, py = w * pivotX, h * pivotY
	var sin, cos = SinCos(item.Shape.Angle)
	var wx, wy = w * cos, w * sin
	var hx, hy = -h * sin, h * cos
	var x0, y0 = -px*cos + py*sin, -px*sin - py*cos

	for i := range vertices {
		vertices[i].ColorR = origW
		vertices[i].ColorG = origH
		vertices[i].ColorB = angle.ToRadians(item.Shape.Angle)
		vertices[i].ColorA = packColor(item.Color)
	}

	vertices[0].DstX, vertices[0].DstY = x+x0, y+y0             // Top-Left
	vertices[1].DstX, vertices[1].DstY = x+x0+wx, y+y0+wy       // Top-Right
	vertices[2].DstX, vertices[2].DstY = x+x0+hx, y+y0+hy       // Bottom-Left
	vertices[3].DstX, vertices[3].DstY = x+x0+wx+hx, y+y0+wy+hy // Bottom-Right

	switch item.Kind {
	case KindShape:
		op.Images[0] = White1x1
		for i := range vertices {
			vertices[i].Custom0 = item.Shape.Roundness
			vertices[i].Custom1 = packColor(item.OutlineColor)
			vertices[i].Custom2 = item.OutlineSize
			vertices[i].Custom3 = 0
		}
	case KindImage:
		item.ImageX = 0.5
		op.Images[0] = Images[item.Image-1]
		for i := range vertices {
			vertices[i].Custom0 = item.Shape.Roundness
			vertices[i].Custom1 = packColor(item.OutlineColor)
			vertices[i].Custom2 = item.OutlineSize
			vertices[i].Custom3 = 0
		}
	case KindText:
		if item.Font == 0 {
			ebitenutil.DebugPrintAt(screen, item.Text, int(item.Shape.X), int(item.Shape.Y))
			return
		}
		for i := range vertices {
			vertices[i].Custom0 = 0 // SymbolGap
			vertices[i].Custom1 = 0 // LineGap
			vertices[i].Custom2 = 0 // LineHeight
			vertices[i].Custom3 = packTextLayout(0, false)
		}
	}

	var imgW, imgH = float32(op.Images[0].Bounds().Dx()), float32(op.Images[0].Bounds().Dy())

	// Resolve normalized sub-region (default 0,0,0,0 → 0,0,1,1 = full image).
	var subW, subH = item.ImageWidth, item.ImageHeight
	if subW == 0 {
		subW = 1
	}
	if subH == 0 {
		subH = 1
	}
	var srcX, srcY = item.ImageX * imgW, item.ImageY * imgH
	var srcW, srcH = subW * imgW, subH * imgH

	// Calculate UV padding proportionally to the sub-region size.
	var padU, padV float32 = 0, 0
	if origW > 0 {
		padU = pad * (srcW / origW)
	}
	if origH > 0 {
		padV = pad * (srcH / origH)
	}

	// Assign expanded UVs. The shader will sample outside the sub-region in the padding area,
	// but shapeAlpha = 0 will mask it out, allowing only the custom outline color to show.
	vertices[0].SrcX, vertices[0].SrcY = srcX-padU, srcY-padV
	vertices[1].SrcX, vertices[1].SrcY = srcX+srcW+padU, srcY-padV
	vertices[2].SrcX, vertices[2].SrcY = srcX-padU, srcY+srcH+padV
	vertices[3].SrcX, vertices[3].SrcY = srcX+srcW+padU, srcY+srcH+padV

	screen.DrawTrianglesShader(vertices, indices, shader, op)
}
