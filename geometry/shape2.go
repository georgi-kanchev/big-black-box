package geometry

import (
	"big-black-box/utility/number"
	"math"
)

type Shape struct {
	X, Y, Width, Height float32
	ar                  float32
}

func (s Shape) Angle() float32                  { var a, _ = unpackAR(s.ar); return a }
func (s Shape) Roundness() float32              { var _, r = unpackAR(s.ar); return r }
func (s *Shape) SetAR(angle, roundness float32) { s.ar = packAR(angle, roundness) }

func packAR(angle, roundness float32) float32 {
	var idx = ((int(angle*10) % 3600) + 3600) % 3600
	var r = min(int(roundness*4659), 4659)
	return -float32(r*3600 + idx + 1)
}

func unpackAR(ar float32) (angle, roundness float32) {
	var abs = int(-ar) - 1
	if abs < 0 {
		return 0, 0
	}
	return float32(abs%3600) * 0.1, float32(abs/3600) / 4659
}

// DistanceToPoint returns the signed distance from the shape boundary to (x, y).
// Negative means inside, zero means on the edge, positive means outside.
func (s Shape) DistanceToPoint(x, y float32) float32 {
	var angle, roundness = unpackAR(s.ar)
	var px = x - s.X
	var py = y - s.Y

	var rad = float64(-angle) * (math.Pi / 180.0)
	var cosR = float32(math.Cos(rad))
	var sinR = float32(math.Sin(rad))
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = roundness * min(hx, hy)
	var qx = number.Absolute(lx) - (hx - r)
	var qy = number.Absolute(ly) - (hy - r)
	return number.SquareRoot(max(qx, 0)*max(qx, 0)+max(qy, 0)*max(qy, 0)) + min(max(qx, qy), 0) - r
}

func (s Shape) Contains(x, y float32) bool {
	return s.DistanceToPoint(x, y) <= 0
}

// ClosestPointOnEdge returns the nearest point on the shape boundary to (x, y).
func (s Shape) ClosestPointOnEdge(x, y float32) (float32, float32) {
	var angle, roundness = unpackAR(s.ar)
	var px = x - s.X
	var py = y - s.Y

	var rad = float64(-angle) * (math.Pi / 180.0)
	var cosR = float32(math.Cos(rad))
	var sinR = float32(math.Sin(rad))
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = roundness * min(hx, hy)
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

// func (s Shape) Raycast(ox, oy, dirX, dirY float32) (x, y, t float32, hit bool)

// Bounds returns the tight axis-aligned bounding box of the shape.
func (s Shape) Bounds() (minX, minY, maxX, maxY float32) {
	var angle, roundness = unpackAR(s.ar)
	var rad = float64(angle) * (math.Pi / 180.0)
	var cosR = number.Absolute(float32(math.Cos(rad)))
	var sinR = number.Absolute(float32(math.Sin(rad)))
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = roundness * min(hx, hy)
	var extentX = (hx-r)*cosR + (hy-r)*sinR + r
	var extentY = (hx-r)*sinR + (hy-r)*cosR + r
	return s.X - extentX, s.Y - extentY, s.X + extentX, s.Y + extentY
}

// Overlap reports whether s and other intersect.
func (s Shape) Overlap(other Shape) bool {
	{ // AABB broadphase
		var sMinX, sMinY, sMaxX, sMaxY = s.Bounds()
		var oMinX, oMinY, oMaxX, oMaxY = other.Bounds()
		if sMaxX < oMinX || oMaxX < sMinX || sMaxY < oMinY || oMaxY < sMinY {
			return false
		}
	}

	var dx, dy = other.X - s.X, other.Y - s.Y
	var sAngle, _ = unpackAR(s.ar)
	var sRad = float64(sAngle) * (math.Pi / 180.0)
	var sCos, sSin = float32(math.Cos(sRad)), float32(math.Sin(sRad))
	var oAngle, _ = unpackAR(other.ar)
	var oRad = float64(oAngle) * (math.Pi / 180.0)
	var oCos, oSin = float32(math.Cos(oRad)), float32(math.Sin(oRad))

	{ // s local X
		var proj = number.Absolute(dx*sCos + dy*sSin)
		if proj > support(s, sCos, sSin, sCos, sSin)+support(other, sCos, sSin, oCos, oSin) {
			return false
		}
	}
	{ // s local Y
		var proj = number.Absolute(dx*-sSin + dy*sCos)
		if proj > support(s, -sSin, sCos, sCos, sSin)+support(other, -sSin, sCos, oCos, oSin) {
			return false
		}
	}
	{ // other local X
		var proj = number.Absolute(dx*oCos + dy*oSin)
		if proj > support(s, oCos, oSin, sCos, sSin)+support(other, oCos, oSin, oCos, oSin) {
			return false
		}
	}
	{ // other local Y
		var proj = number.Absolute(dx*-oSin + dy*oCos)
		if proj > support(s, -oSin, oCos, sCos, sSin)+support(other, -oSin, oCos, oCos, oSin) {
			return false
		}
	}
	{ // corner axis
		var pAx, pAy = nearestInnerBoxPoint(s, other.X, other.Y, sCos, sSin)
		var pBx, pBy = nearestInnerBoxPoint(other, s.X, s.Y, oCos, oSin)
		var cax, cay = pBx - pAx, pBy - pAy
		var l = number.SquareRoot(cax*cax + cay*cay)
		if l > 1e-6 {
			cax, cay = cax/l, cay/l
			var proj = number.Absolute(dx*cax + dy*cay)
			if proj > support(s, cax, cay, sCos, sSin)+support(other, cax, cay, oCos, oSin) {
				return false
			}
		}
	}

	return true
}

// Collide returns other moved out of s by the minimum translation vector.
func (s Shape) Collide(other Shape) Shape {
	if !s.Overlap(other) {
		return other
	}

	var dx, dy = other.X - s.X, other.Y - s.Y
	var sAngle, _ = unpackAR(s.ar)
	var sRad = float64(sAngle) * (math.Pi / 180.0)
	var sCos, sSin = float32(math.Cos(sRad)), float32(math.Sin(sRad))
	var oAngle, _ = unpackAR(other.ar)
	var oRad = float64(oAngle) * (math.Pi / 180.0)
	var oCos, oSin = float32(math.Cos(oRad)), float32(math.Sin(oRad))

	var minDepth = number.ValueBiggest[float32]()
	var minAx0, minAx1 float32

	{ // s local X
		var ax0, ax1 = sCos, sSin
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = support(s, ax0, ax1, sCos, sSin) + support(other, ax0, ax1, oCos, oSin) - d
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
		var depth = support(s, ax0, ax1, sCos, sSin) + support(other, ax0, ax1, oCos, oSin) - d
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
		var depth = support(s, ax0, ax1, sCos, sSin) + support(other, ax0, ax1, oCos, oSin) - d
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
		var depth = support(s, ax0, ax1, sCos, sSin) + support(other, ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // corner axis
		var pAx, pAy = nearestInnerBoxPoint(s, other.X, other.Y, sCos, sSin)
		var pBx, pBy = nearestInnerBoxPoint(other, s.X, s.Y, oCos, oSin)
		var ax0, ax1 = pBx - pAx, pBy - pAy
		var l = number.SquareRoot(ax0*ax0 + ax1*ax1)
		if l > 1e-6 {
			ax0, ax1 = ax0/l, ax1/l
			var d = dx*ax0 + dy*ax1
			if d < 0 {
				d, ax0, ax1 = -d, -ax0, -ax1
			}
			var depth = support(s, ax0, ax1, sCos, sSin) + support(other, ax0, ax1, oCos, oSin) - d
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

func support(s Shape, ax0, ax1, cosR, sinR float32) float32 {
	var _, roundness = unpackAR(s.ar)
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = roundness * min(hx, hy)
	var dX = number.Absolute(ax0*cosR + ax1*sinR)
	var dY = number.Absolute(-ax0*sinR + ax1*cosR)
	return (hx-r)*dX + (hy-r)*dY + r
}
func nearestInnerBoxPoint(s Shape, px, py, cosR, sinR float32) (float32, float32) {
	var _, roundness = unpackAR(s.ar)
	var rx, ry = px - s.X, py - s.Y
	var lx = rx*cosR + ry*sinR
	var ly = -rx*sinR + ry*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = roundness * min(hx, hy)
	lx = max(-(hx - r), min(hx-r, lx))
	ly = max(-(hy - r), min(hy-r, ly))

	return s.X + lx*cosR - ly*sinR, s.Y + lx*sinR + ly*cosR
}
