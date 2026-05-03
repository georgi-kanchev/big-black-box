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
	queue(&internal.DrawItem{
		View:         internal.View(v),
		Kind:         internal.KindShape,
		Shape:        internal.Shape(shape),
		Color:        palette.White,
		OutlineSize:  2,
		OutlineColor: palette.Green,
	})

}
func (v View) DrawImage(shape geometry.Shape, image assets.Image) {
	queue(&internal.DrawItem{
		View:         internal.View(v),
		Kind:         internal.KindImage,
		Shape:        internal.Shape(shape),
		Color:        palette.White,
		Image:        internal.Image(image),
		OutlineSize:  1,
		OutlineColor: palette.Red,
	})
}
func (v View) DrawText(shape geometry.Shape, font assets.Font) {
	queue(&internal.DrawItem{
		View:  internal.View(v),
		Kind:  internal.KindText,
		Shape: internal.Shape(shape),
		Color: palette.White,
		Font:  internal.Font(font),
	})
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

// private ========================================================

func queue(item *internal.DrawItem) {
	var mask = intersectAreas(item.MaskArea, item.View.MaskArea)
	var hasMask = item.MaskArea != (internal.Area{}) || item.View.MaskArea != (internal.Area{})
	var isCulled = hasMask && (mask.Width <= 0 || mask.Height <= 0)
	if isCulled {
		return
	}

	const pivotX, pivotY float32 = 0.5, 0.5
	var pad float32 = 2
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

	var poly [4]ebiten.Vertex
	var packedCol = packColor(item.Color)
	var rads = angle.ToRadians(item.Shape.Angle)

	for i := range poly {
		poly[i].ColorR = origW
		poly[i].ColorG = origH
		poly[i].ColorB = rads
		poly[i].ColorA = packedCol
	}

	poly[0].DstX, poly[0].DstY = x+x0, y+y0
	poly[1].DstX, poly[1].DstY = x+x0+wx, y+y0+wy
	poly[2].DstX, poly[2].DstY = x+x0+wx+hx, y+y0+wy+hy
	poly[3].DstX, poly[3].DstY = x+x0+hx, y+y0+hy

	var img *ebiten.Image
	switch item.Kind {
	case internal.KindShape:
		img = internal.White1x1
		for i := range poly {
			poly[i].Custom0 = item.Shape.Roundness
			poly[i].Custom1 = packColor(item.OutlineColor)
			poly[i].Custom2 = item.OutlineSize
			poly[i].Custom3 = float32(item.Kind)
		}
	case internal.KindImage:
		img = internal.Images[item.Image-1]
		for i := range poly {
			poly[i].Custom0 = item.Shape.Roundness
			poly[i].Custom1 = packColor(item.OutlineColor)
			poly[i].Custom2 = item.OutlineSize
			poly[i].Custom3 = float32(item.Kind)
		}
	case internal.KindText:
		if item.Font == 0 {
			return
		}
		img = internal.Fonts[item.Font-1]
		for i := range poly {
			poly[i].Custom0 = 0
			poly[i].Custom1 = 0
			poly[i].Custom2 = 0
			poly[i].Custom3 = packTextLayout(item.Kind, 0, false)
		}
	}

	var imgW, imgH = float32(img.Bounds().Dx()), float32(img.Bounds().Dy())
	var subW, subH = item.TexWidth, item.TexHeight
	if subW == 0 {
		subW = 1
	}
	if subH == 0 {
		subH = 1
	}

	var srcX, srcY = item.TexX * imgW, item.TexY * imgH
	var srcW, srcH = subW * imgW, subH * imgH

	var padU, padV float32
	if origW > 0 {
		padU = pad * (srcW / origW)
	}
	if origH > 0 {
		padV = pad * (srcH / origH)
	}

	poly[0].SrcX, poly[0].SrcY = srcX-padU, srcY-padV
	poly[1].SrcX, poly[1].SrcY = srcX+srcW+padU, srcY-padV
	poly[2].SrcX, poly[2].SrcY = srcX+srcW+padU, srcY+srcH+padV
	poly[3].SrcX, poly[3].SrcY = srcX-padU, srcY+srcH+padV

	var count = 4
	var workingVerts = poly[:]

	if hasMask {
		var clipDest [12]ebiten.Vertex
		var clipTemp [12]ebiten.Vertex
		count = clipPolygonAABB(workingVerts, clipDest[:], clipTemp[:], mask)
		if count < 3 {
			return
		}
		workingVerts = clipDest[:count]
	}

	var windowArea = item.View.WindowArea
	if windowArea == (internal.Area{}) {
		var w, h = internal.Engine.Layout(ebiten.WindowSize())
		windowArea.Width, windowArea.Height = float32(w), float32(h)
	}

	var vSin, vCos = internal.SinCos(-item.View.Angle)
	var centerX = windowArea.X + windowArea.Width*0.5
	var centerY = windowArea.Y + windowArea.Height*0.5

	for i := 0; i < count; i++ {
		var dx = workingVerts[i].DstX - item.View.X
		var dy = workingVerts[i].DstY - item.View.Y

		if item.View.Angle != 0 {
			var rx = dx*vCos - dy*vSin
			var ry = dx*vSin + dy*vCos
			dx, dy = rx, ry
		}

		workingVerts[i].DstX = (dx * item.View.Zoom) + centerX
		workingVerts[i].DstY = (dy * item.View.Zoom) + centerY
	}

	var winDest [12]ebiten.Vertex
	var winTemp [12]ebiten.Vertex
	count = clipPolygonAABB(workingVerts, winDest[:], winTemp[:], windowArea)
	if count < 3 {
		return
	}
	workingVerts = winDest[:count]

	item.VertCount = count
	for i := 0; i < count; i++ {
		item.Verts[i] = workingVerts[i]
	}

	internal.Queue(internal.LayerDefault, *item)
}

func clipPolygonAABB(poly, outBuf, tempBuf []ebiten.Vertex, mask Area) int {
	var minX, maxX = mask.X, mask.X + mask.Width
	var minY, maxY = mask.Y, mask.Y + mask.Height
	var count = clipPolyEdge(poly, tempBuf, true, minX, true)
	if count == 0 {
		return 0
	}

	count = clipPolyEdge(tempBuf[:count], outBuf, true, maxX, false)
	if count == 0 {
		return 0
	}

	count = clipPolyEdge(outBuf[:count], tempBuf, false, minY, true)
	if count == 0 {
		return 0
	}

	count = clipPolyEdge(tempBuf[:count], outBuf, false, maxY, false)
	return count
}
func clipPolyEdge(in, out []ebiten.Vertex, isX bool, edgeVal float32, keepGreater bool) int {
	var outCount = 0
	if len(in) == 0 {
		return 0
	}

	var prev = in[len(in)-1]
	var prevVal float32
	if isX {
		prevVal = prev.DstX
	} else {
		prevVal = prev.DstY
	}
	var prevInside = (keepGreater && prevVal >= edgeVal) || (!keepGreater && prevVal <= edgeVal)

	for _, curr := range in {
		var currVal float32
		if isX {
			currVal = curr.DstX
		} else {
			currVal = curr.DstY
		}
		var currInside = (keepGreater && currVal >= edgeVal) || (!keepGreater && currVal <= edgeVal)

		if currInside != prevInside {
			var t float32
			if isX {
				t = (edgeVal - prev.DstX) / (curr.DstX - prev.DstX)
			} else {
				t = (edgeVal - prev.DstY) / (curr.DstY - prev.DstY)
			}

			out[outCount] = ebiten.Vertex{
				DstX:   prev.DstX + t*(curr.DstX-prev.DstX),
				DstY:   prev.DstY + t*(curr.DstY-prev.DstY),
				SrcX:   prev.SrcX + t*(curr.SrcX-prev.SrcX),
				SrcY:   prev.SrcY + t*(curr.SrcY-prev.SrcY),
				ColorR: prev.ColorR, ColorG: prev.ColorG, ColorB: prev.ColorB, ColorA: prev.ColorA,
				Custom0: prev.Custom0, Custom1: prev.Custom1, Custom2: prev.Custom2, Custom3: prev.Custom3,
			}
			outCount++
		}

		if currInside {
			out[outCount] = curr
			outCount++
		}

		prev = curr
		prevInside = currInside
	}

	return outCount
}

func packColor(c uint) float32 {
	r, g, b, a := color.Channels(c)
	return float32(uint32(r>>2)<<18 | uint32(g>>2)<<12 | uint32(b>>2)<<6 | uint32(a>>2))
}
func packTextLayout(kind internal.Kind, align byte, wordWrap bool) float32 {
	v := uint32(kind & 0x7)
	if wordWrap {
		v |= 1 << 3
	}
	v |= uint32(align&0x7) << 4
	return float32(v)
}

func intersectAreas(a, b Area) Area {
	if a == (Area{}) {
		return b
	}
	if b == (Area{}) {
		return a
	}

	var x1 = max(a.X, b.X)
	var y1 = max(a.Y, b.Y)
	var x2 = min(a.X+a.Width, b.X+b.Width)
	var y2 = min(a.Y+a.Height, b.Y+b.Height)
	return Area{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
}
