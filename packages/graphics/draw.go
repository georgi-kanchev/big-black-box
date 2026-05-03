package graphics

import (
	"big-black-box/packages/assets"
	"big-black-box/packages/geometry"
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/color"
	"big-black-box/packages/utility/color/palette"
	"big-black-box/packages/utility/debug"
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
	item.Verts, item.VertCount = buildVerts(item)
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
	item.Verts, item.VertCount = buildVerts(item)
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
	item.Verts, item.VertCount = buildVerts(item)
	internal.Queue(internal.LayerDefault, item)
}

// private ========================================================

func buildVerts(item internal.DrawItem) (poly [8]ebiten.Vertex, n int) {
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

	// TL, TR, BL, BR
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
			return poly, 0
		}
		img = internal.Fonts[item.Font-1]
		for i := range verts {
			verts[i].Custom0 = 0
			verts[i].Custom1 = 0
			verts[i].Custom2 = 0
			verts[i].Custom3 = packTextLayout(item.Kind, 0, false)
		}
	default:
		return poly, 0
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

	// TL, TR, BL, BR
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

	// Reorder TL,TR,BL,BR → TL,TR,BR,BL (CW polygon winding for Sutherland-Hodgman).
	poly[0], poly[1], poly[2], poly[3] = verts[0], verts[1], verts[3], verts[2]
	n = 4

	// Clip to MaskArea in view space.
	if item.View.MaskArea != (internal.Area{}) {
		n = clipPoly(&poly, n, item.View.MaskArea)
		if n < 3 {
			return poly, 0
		}
	}

	// View → screen space (window area offset).
	var vts = view.viewToScreen()
	for i := range n {
		var dx, dy = vts.Apply(float64(poly[i].DstX), float64(poly[i].DstY))
		poly[i].DstX, poly[i].DstY = float32(dx), float32(dy)
	}

	// Clip to WindowArea in window/screen space.
	if item.View.WindowArea != (internal.Area{}) {
		n = clipPoly(&poly, n, item.View.WindowArea)
		if n < 3 {
			return poly, 0
		}
	}

	return poly, n
}

// clipPoly clips a convex polygon against an axis-aligned rectangle using Sutherland-Hodgman.
// poly holds up to 8 vertices; n is the input count. Returns the new vertex count.
func clipPoly(poly *[8]ebiten.Vertex, n int, area internal.Area) int {
	n = clipPlane(poly, n, true, false, area.X)             // left:   x >= area.X
	n = clipPlane(poly, n, true, true, area.X+area.Width)   // right:  x <= area.X+W
	n = clipPlane(poly, n, false, false, area.Y)            // top:    y >= area.Y
	n = clipPlane(poly, n, false, true, area.Y+area.Height) // bottom: y <= area.Y+H
	return n
}

// clipPlane clips the polygon against one half-plane.
// xAxis=true clips on X, xAxis=false on Y.
// maxSide=false keeps vertices >= val; maxSide=true keeps vertices <= val.
func clipPlane(poly *[8]ebiten.Vertex, n int, xAxis, maxSide bool, val float32) int {
	if n == 0 {
		return 0
	}
	var scratch [8]ebiten.Vertex
	var out int

	coord := func(v ebiten.Vertex) float32 {
		if xAxis {
			return v.DstX
		}
		return v.DstY
	}
	inside := func(v ebiten.Vertex) bool {
		c := coord(v)
		if maxSide {
			return c <= val
		}
		return c >= val
	}

	for i := range n {
		cur := poly[i]
		prev := poly[(i+n-1)%n]
		curIn := inside(cur)
		prevIn := inside(prev)

		if curIn != prevIn {
			// Edge crosses the clip boundary — compute intersection.
			var t float32
			dc, dp := coord(cur), coord(prev)
			if dc != dp {
				t = (val - dp) / (dc - dp)
			}
			scratch[out] = lerpVert(prev, cur, t)
			out++
		}
		if curIn {
			scratch[out] = cur
			out++
		}
	}

	copy(poly[:], scratch[:out])
	return out
}

// lerpVert linearly interpolates Dst and Src between a and b at parameter t.
// Color and Custom fields are copied from a (they are uniform across the polygon).
func lerpVert(a, b ebiten.Vertex, t float32) ebiten.Vertex {
	return ebiten.Vertex{
		DstX:   a.DstX + t*(b.DstX-a.DstX),
		DstY:   a.DstY + t*(b.DstY-a.DstY),
		SrcX:   a.SrcX + t*(b.SrcX-a.SrcX),
		SrcY:   a.SrcY + t*(b.SrcY-a.SrcY),
		ColorR: a.ColorR, ColorG: a.ColorG, ColorB: a.ColorB, ColorA: a.ColorA,
		Custom0: a.Custom0, Custom1: a.Custom1, Custom2: a.Custom2, Custom3: a.Custom3,
	}
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
