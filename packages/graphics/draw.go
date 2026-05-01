package graphics

import (
	"big-black-box/packages/assets"
	"big-black-box/packages/debug"
	"big-black-box/packages/geometry"
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/color"
	"big-black-box/packages/utility/color/palette"
	"big-black-box/packages/utility/text"
	"big-black-box/packages/utility/time"
	"big-black-box/packages/utility/time/unit"

	"github.com/hajimehoshi/ebiten/v2"
)

func (v View) DrawShape(shape geometry.Shape) {
	var item = internal.DrawItem{
		View:         internal.View(v),
		Kind:         internal.KindShape,
		Shape:        internal.Shape(shape),
		Color:        palette.White,
		OutlineSize:  1,
		OutlineColor: palette.Green,
	}
	item.Verts = buildVerts(item)
	internal.Queue(internal.LayerDefault, item)
}

func (v View) DrawFPS() {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		View:  internal.View(v),
		Kind:  internal.KindText,
		Shape: internal.Shape{X: 5, Y: 5},
		Text:  text.Start().String("FPS: ").Int(int(internal.FPS)).End(),
		Color: palette.White,
	})
}
func (v View) DrawDebugInfo() {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		View:  internal.View(v),
		Kind:  internal.KindText,
		Shape: internal.Shape{X: 5, Y: 5},
		Text: text.Start().String("FPS: ").Int(int(internal.FPS)).String(" TPS: ").Int(int(internal.TPS)).String("\n\n").
			String("Runtime: ").
			String(time.AsClock12(internal.Runtime, ":", unit.Hour|unit.Minute|unit.Second|unit.Millisecond, false)).
			String("\n\n").
			String(debug.MemoryUsage()).
			End(),
		Color: palette.White,
	})
}

func (v View) DrawImage(shape geometry.Shape, image assets.Image) {
	var item = internal.DrawItem{
		View:         internal.View(v),
		Kind:         internal.KindImage,
		Shape:        internal.Shape(shape),
		Color:        palette.White,
		Image:        internal.Image(image),
		OutlineSize:  1,
		OutlineColor: palette.Red,
	}
	item.Verts = buildVerts(item)
	internal.Queue(internal.LayerDefault, item)
}

func (v View) DrawText(shape geometry.Shape, font assets.Font) {
	var item = internal.DrawItem{
		View:  internal.View(v),
		Kind:  internal.KindText,
		Shape: internal.Shape(shape),
		Color: palette.White,
		Font:  internal.Font(font),
	}
	item.Verts = buildVerts(item)
	internal.Queue(internal.LayerDefault, item)
}

// private ========================================================

func buildVerts(item internal.DrawItem) [4]ebiten.Vertex {
	const pivotX, pivotY float32 = 0.5, 0.5
	var pad float32 = 1.5
	if item.OutlineSize > 0 {
		pad += item.OutlineSize
	}

	var origW, origH = item.Shape.Width, item.Shape.Height
	var w, h = origW + pad*2, origH + pad*2
	var x, y = item.Shape.X, item.Shape.Y
	var px, py = w * pivotX, h * pivotY
	var sin, cos = internal.SinCos(item.Shape.Angle)
	var wx, wy = w * cos, w * sin
	var hx, hy = -h * sin, h * cos
	var x0, y0 = -px*cos + py*sin, -px*sin - py*cos

	var verts [4]ebiten.Vertex
	for i := range verts {
		verts[i].ColorR = origW
		verts[i].ColorG = origH
		verts[i].ColorB = angle.ToRadians(item.Shape.Angle)
		verts[i].ColorA = packColor(item.Color)
	}

	verts[0].DstX, verts[0].DstY = x+x0, y+y0
	verts[1].DstX, verts[1].DstY = x+x0+wx, y+y0+wy
	verts[2].DstX, verts[2].DstY = x+x0+hx, y+y0+hy
	verts[3].DstX, verts[3].DstY = x+x0+wx+hx, y+y0+wy+hy

	var img *ebiten.Image
	switch item.Kind {
	case internal.KindShape:
		img = internal.White1x1
		for i := range verts {
			verts[i].Custom0 = item.Shape.Roundness
			verts[i].Custom1 = packColor(item.OutlineColor)
			verts[i].Custom2 = item.OutlineSize
			verts[i].Custom3 = float32(item.Kind)
		}
	case internal.KindImage:
		img = internal.Images[item.Image-1]
		for i := range verts {
			verts[i].Custom0 = item.Shape.Roundness
			verts[i].Custom1 = packColor(item.OutlineColor)
			verts[i].Custom2 = item.OutlineSize
			verts[i].Custom3 = float32(item.Kind)
		}
	case internal.KindText:
		if item.Font == 0 {
			return verts
		}
		img = internal.Fonts[item.Font-1]
		for i := range verts {
			verts[i].Custom0 = 0
			verts[i].Custom1 = 0
			verts[i].Custom2 = 0
			verts[i].Custom3 = packTextLayout(item.Kind, 0, false)
		}
	default:
		return verts
	}

	var imgW, imgH = float32(img.Bounds().Dx()), float32(img.Bounds().Dy())
	var subW, subH = item.ImageWidth, item.ImageHeight
	if subW == 0 {
		subW = 1
	}
	if subH == 0 {
		subH = 1
	}
	var srcX, srcY = item.ImageX * imgW, item.ImageY * imgH
	var srcW, srcH = subW * imgW, subH * imgH

	var padU, padV float32
	if origW > 0 {
		padU = pad * (srcW / origW)
	}
	if origH > 0 {
		padV = pad * (srcH / origH)
	}

	verts[0].SrcX, verts[0].SrcY = srcX-padU, srcY-padV
	verts[1].SrcX, verts[1].SrcY = srcX+srcW+padU, srcY-padV
	verts[2].SrcX, verts[2].SrcY = srcX-padU, srcY+srcH+padV
	verts[3].SrcX, verts[3].SrcY = srcX+srcW+padU, srcY+srcH+padV

	var view = View(item.View)

	// World → view space (position, zoom, angle).
	var wtv = view.worldToView()
	for i := range verts {
		var dx, dy = wtv.Apply(float64(verts[i].DstX), float64(verts[i].DstY))
		verts[i].DstX, verts[i].DstY = float32(dx), float32(dy)
	}

	// Clip to MaskArea in view space.
	if item.View.MaskArea != (internal.Area{}) {
		verts = cropVerts(verts, item.View.MaskArea)
	}

	// View → screen space (window area offset).
	var vts = view.viewToScreen()
	for i := range verts {
		var dx, dy = vts.Apply(float64(verts[i].DstX), float64(verts[i].DstY))
		verts[i].DstX, verts[i].DstY = float32(dx), float32(dy)
	}

	// Clip to WindowArea in window/screen space.
	if item.View.WindowArea != (internal.Area{}) {
		verts = cropVerts(verts, item.View.WindowArea)
	}

	return verts
}

// cropVerts clips a quad's bounding box to area, interpolating SrcX/SrcY proportionally.
// Correct for axis-aligned quads; for rotated quads it operates on the screen-space AABB.
func cropVerts(verts [4]ebiten.Vertex, area internal.Area) [4]ebiten.Vertex {
	var dstMinX = min(min(verts[0].DstX, verts[1].DstX), min(verts[2].DstX, verts[3].DstX))
	var dstMaxX = max(max(verts[0].DstX, verts[1].DstX), max(verts[2].DstX, verts[3].DstX))
	var dstMinY = min(min(verts[0].DstY, verts[1].DstY), min(verts[2].DstY, verts[3].DstY))
	var dstMaxY = max(max(verts[0].DstY, verts[1].DstY), max(verts[2].DstY, verts[3].DstY))

	var dstW, dstH = dstMaxX - dstMinX, dstMaxY - dstMinY
	if dstW <= 0 || dstH <= 0 {
		return verts
	}

	var newMinX = max(dstMinX, area.X)
	var newMaxX = min(dstMaxX, area.X+area.Width)
	var newMinY = max(dstMinY, area.Y)
	var newMaxY = min(dstMaxY, area.Y+area.Height)

	if newMinX >= newMaxX || newMinY >= newMaxY {
		for i := range verts {
			verts[i].DstX, verts[i].DstY = newMinX, newMinY
		}
		return verts
	}

	var srcMinX = min(min(verts[0].SrcX, verts[1].SrcX), min(verts[2].SrcX, verts[3].SrcX))
	var srcMaxX = max(max(verts[0].SrcX, verts[1].SrcX), max(verts[2].SrcX, verts[3].SrcX))
	var srcMinY = min(min(verts[0].SrcY, verts[1].SrcY), min(verts[2].SrcY, verts[3].SrcY))
	var srcMaxY = max(max(verts[0].SrcY, verts[1].SrcY), max(verts[2].SrcY, verts[3].SrcY))

	var t0x = (newMinX - dstMinX) / dstW
	var t1x = (newMaxX - dstMinX) / dstW
	var t0y = (newMinY - dstMinY) / dstH
	var t1y = (newMaxY - dstMinY) / dstH
	var srcW, srcH = srcMaxX - srcMinX, srcMaxY - srcMinY

	verts[0].DstX, verts[0].DstY = newMinX, newMinY
	verts[1].DstX, verts[1].DstY = newMaxX, newMinY
	verts[2].DstX, verts[2].DstY = newMinX, newMaxY
	verts[3].DstX, verts[3].DstY = newMaxX, newMaxY

	verts[0].SrcX, verts[0].SrcY = srcMinX+t0x*srcW, srcMinY+t0y*srcH
	verts[1].SrcX, verts[1].SrcY = srcMinX+t1x*srcW, srcMinY+t0y*srcH
	verts[2].SrcX, verts[2].SrcY = srcMinX+t0x*srcW, srcMinY+t1y*srcH
	verts[3].SrcX, verts[3].SrcY = srcMinX+t1x*srcW, srcMinY+t1y*srcH

	return verts
}

// packColor packs an RGBA uint into a 24-bit float (6 bits per channel).
// Layout: R6[23:18] G6[17:12] B6[11:6] A6[5:0].
func packColor(c uint) float32 {
	r, g, b, a := color.Channels(c)
	return float32(uint32(r>>2)<<18 | uint32(g>>2)<<12 | uint32(b>>2)<<6 | uint32(a>>2))
}

// packTextLayout packs kind (0–7), wordwrap, and align (0–7) into custom.w.
// Layout: Align[6:4] WordWrap[3] Kind[2:0].
func packTextLayout(kind internal.Kind, align byte, wordWrap bool) float32 {
	v := uint32(kind & 0x7)
	if wordWrap {
		v |= 1 << 3
	}
	v |= uint32(align&0x7) << 4
	return float32(v)
}
