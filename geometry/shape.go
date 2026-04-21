package geometry

import (
	"big-black-box/internal"
	"big-black-box/utility/number"
	"big-black-box/utility/point"
)

type Shape internal.Shape

// Negative result means inside, zero means on the edge, positive means outside.
func (s Shape) DistanceToPoint(x, y float32) float32 {
	var px = x - s.X
	var py = y - s.Y

	var sinR, cosR = internal.SinCos(-s.Angle)
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	var qx = number.Absolute(lx) - (hx - r)
	var qy = number.Absolute(ly) - (hy - r)
	return number.SquareRoot(max(qx, 0)*max(qx, 0)+max(qy, 0)*max(qy, 0)) + min(max(qx, qy), 0) - r
}
func (s Shape) Contains(x, y float32) bool {
	return s.DistanceToPoint(x, y) <= 0
}
func (s Shape) ClosestPointOnEdge(x, y float32) (float32, float32) {
	var px = x - s.X
	var py = y - s.Y

	var sinR, cosR = internal.SinCos(-s.Angle)
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	var cx = max(-(hx - r), min(hx-r, lx))
	var cy = max(-(hy - r), min(hy-r, ly))

	var dx = lx - cx
	var dy = ly - cy
	var dist = number.SquareRoot(dx*dx + dy*dy)

	var bx, by float32
	if dist > 1e-6 {
		bx = cx + r*dx/dist
		by = cy + r*dy/dist
	} else {
		// inside inner box: push to nearest flat face
		var dPosX = hx - lx
		var dNegX = hx + lx
		var dPosY = hy - ly
		var dNegY = hy + ly
		var minD = min(min(dPosX, dNegX), min(dPosY, dNegY))
		switch minD {
		case dPosX:
			bx, by = hx, ly
		case dNegX:
			bx, by = -hx, ly
		case dPosY:
			bx, by = lx, hy
		default:
			bx, by = lx, -hy
		}
	}

	return bx*cosR + by*sinR + s.X, -bx*sinR + by*cosR + s.Y
}
func (s Shape) Raycast(x, y, angle, length float32) (float32, float32) {
	var lx, ly = point.MoveAtAngle(x, y, angle, length*0.5)
	if !s.Overlap(Shape{X: lx, Y: ly, Width: length, Angle: angle}) {
		return number.NaN(), number.NaN()
	}
	var dy, dx = internal.SinCos(angle)
	var t = float32(0)
	for range 64 {
		var cx, cy = x + t*dx, y + t*dy
		var dist = s.DistanceToPoint(cx, cy)
		if dist <= 1e-4 {
			return s.ClosestPointOnEdge(cx, cy)
		}
		t += dist
	}
	return number.NaN(), number.NaN()
}
func (s Shape) Bounds() (minX, minY, maxX, maxY float32) {
	var sinR, cosR = internal.SinCos(s.Angle)
	sinR, cosR = number.Absolute(sinR), number.Absolute(cosR)
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	var extentX = (hx-r)*cosR + (hy-r)*sinR + r
	var extentY = (hx-r)*sinR + (hy-r)*cosR + r
	return s.X - extentX, s.Y - extentY, s.X + extentX, s.Y + extentY
}
func (s Shape) Overlap(other Shape) bool {
	{ // AABB broadphase
		var sMinX, sMinY, sMaxX, sMaxY = s.Bounds()
		var oMinX, oMinY, oMaxX, oMaxY = other.Bounds()
		if sMaxX < oMinX || oMaxX < sMinX || sMaxY < oMinY || oMaxY < sMinY {
			return false
		}
	}

	var dx, dy = other.X - s.X, other.Y - s.Y
	var sSin, sCos = internal.SinCos(s.Angle)
	var oSin, oCos = internal.SinCos(other.Angle)

	{ // s local X
		var proj = number.Absolute(dx*sCos + dy*sSin)
		if proj > s.support(sCos, sSin, sCos, sSin)+other.support(sCos, sSin, oCos, oSin) {
			return false
		}
	}
	{ // s local Y
		var proj = number.Absolute(dx*-sSin + dy*sCos)
		if proj > s.support(-sSin, sCos, sCos, sSin)+other.support(-sSin, sCos, oCos, oSin) {
			return false
		}
	}
	{ // other local X
		var proj = number.Absolute(dx*oCos + dy*oSin)
		if proj > s.support(oCos, oSin, sCos, sSin)+other.support(oCos, oSin, oCos, oSin) {
			return false
		}
	}
	{ // other local Y
		var proj = number.Absolute(dx*-oSin + dy*oCos)
		if proj > s.support(-oSin, oCos, sCos, sSin)+other.support(-oSin, oCos, oCos, oSin) {
			return false
		}
	}
	{ // corner axis
		var pAx, pAy = s.nearestInnerBoxPoint(other.X, other.Y, sCos, sSin)
		var pBx, pBy = other.nearestInnerBoxPoint(s.X, s.Y, oCos, oSin)
		var cax, cay = pBx - pAx, pBy - pAy
		var l = number.SquareRoot(cax*cax + cay*cay)
		if l > 1e-6 {
			cax, cay = cax/l, cay/l
			var proj = number.Absolute(dx*cax + dy*cay)
			if proj > s.support(cax, cay, sCos, sSin)+other.support(cax, cay, oCos, oSin) {
				return false
			}
		}
	}

	return true
}
func (s Shape) Collide(other Shape) Shape {
	if !s.Overlap(other) {
		return other
	}

	var dx, dy = other.X - s.X, other.Y - s.Y
	var sSin, sCos = internal.SinCos(s.Angle)
	var oSin, oCos = internal.SinCos(other.Angle)

	var minDepth = number.ValueBiggest[float32]()
	var minAx0, minAx1 float32

	{ // s local X
		var ax0, ax1 = sCos, sSin
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + other.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // s local Y
		var ax0, ax1 = -sSin, sCos
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + other.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // other local X
		var ax0, ax1 = oCos, oSin
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + other.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // other local Y
		var ax0, ax1 = -oSin, oCos
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + other.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // corner axis
		var pAx, pAy = s.nearestInnerBoxPoint(other.X, other.Y, sCos, sSin)
		var pBx, pBy = other.nearestInnerBoxPoint(s.X, s.Y, oCos, oSin)
		var ax0, ax1 = pBx - pAx, pBy - pAy
		var l = number.SquareRoot(ax0*ax0 + ax1*ax1)
		if l > 1e-6 {
			ax0, ax1 = ax0/l, ax1/l
			var d = dx*ax0 + dy*ax1
			if d < 0 {
				d, ax0, ax1 = -d, -ax0, -ax1
			}
			var depth = s.support(ax0, ax1, sCos, sSin) + other.support(ax0, ax1, oCos, oSin) - d
			if depth < minDepth {
				minDepth, minAx0, minAx1 = depth, ax0, ax1
			}
		}
	}

	other.X += minAx0 * minDepth
	other.Y += minAx1 * minDepth
	return other
}

// private ========================================================

func (s Shape) support(ax0, ax1, cosR, sinR float32) float32 {
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	var dX = number.Absolute(ax0*cosR + ax1*sinR)
	var dY = number.Absolute(-ax0*sinR + ax1*cosR)
	return (hx-r)*dX + (hy-r)*dY + r
}
func (s Shape) nearestInnerBoxPoint(px, py, cosR, sinR float32) (float32, float32) {
	var rx, ry = px - s.X, py - s.Y
	var lx = rx*cosR + ry*sinR
	var ly = -rx*sinR + ry*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	lx = max(-(hx - r), min(hx-r, lx))
	ly = max(-(hy - r), min(hy-r, ly))

	return s.X + lx*cosR - ly*sinR, s.Y + lx*sinR + ly*cosR
}
