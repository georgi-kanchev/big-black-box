package geometry

import (
	"big-black-box/utility/number"
	"math"
)

type Shape struct {
	X, Y, Width, Height float32
	Angle               float32 // in degrees
	Roundness           float32 // 0 to 1
}

// DistanceToPoint returns the signed distance from the shape boundary to (x, y).
// Negative means inside, zero means on the edge, positive means outside.
func (s *Shape) DistanceToPoint(x, y float32) float32 {
	var px = x - s.X
	var py = y - s.Y

	var rad = float64(-s.Angle) * (math.Pi / 180.0)
	var cosR = float32(math.Cos(rad))
	var sinR = float32(math.Sin(rad))
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	var qx = number.Absolute(lx) - (hx - r)
	var qy = number.Absolute(ly) - (hy - r)
	return number.SquareRoot(max(qx, 0)*max(qx, 0)+max(qy, 0)*max(qy, 0)) + min(max(qx, qy), 0) - r
}

func (s *Shape) Contains(x, y float32) bool {
	return s.DistanceToPoint(x, y) <= 0
}

// Overlap reports whether s and other intersect.
func (s *Shape) Overlap(other *Shape) bool {
	var dx, dy = other.X - s.X, other.Y - s.Y

	var ra = number.SquareRoot(s.Width*s.Width+s.Height*s.Height) * 0.5
	var rb = number.SquareRoot(other.Width*other.Width+other.Height*other.Height) * 0.5
	if dx*dx+dy*dy > (ra+rb)*(ra+rb) {
		return false
	}

	var sRad = float64(s.Angle) * (math.Pi / 180.0)
	var sCos, sSin = float32(math.Cos(sRad)), float32(math.Sin(sRad))
	var oRad = float64(other.Angle) * (math.Pi / 180.0)
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

// Collide returns the minimum vector to move other out of s, and whether they collide.
// Add the returned (px, py) to other.X / other.Y to resolve the collision.
func (s *Shape) Collide(other *Shape) (px, py float32, collided bool) {
	var dx, dy = other.X - s.X, other.Y - s.Y

	var ra = number.SquareRoot(s.Width*s.Width+s.Height*s.Height) * 0.5
	var rb = number.SquareRoot(other.Width*other.Width+other.Height*other.Height) * 0.5
	if dx*dx+dy*dy > (ra+rb)*(ra+rb) {
		return 0, 0, false
	}

	var sRad = float64(s.Angle) * (math.Pi / 180.0)
	var sCos, sSin = float32(math.Cos(sRad)), float32(math.Sin(sRad))
	var oRad = float64(other.Angle) * (math.Pi / 180.0)
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
		if depth <= 0 {
			return 0, 0, false
		}
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
		if depth <= 0 {
			return 0, 0, false
		}
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
		if depth <= 0 {
			return 0, 0, false
		}
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
		if depth <= 0 {
			return 0, 0, false
		}
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
			if depth <= 0 {
				return 0, 0, false
			}
			if depth < minDepth {
				minDepth, minAx0, minAx1 = depth, ax0, ax1
			}
		}
	}

	return minAx0 * minDepth, minAx1 * minDepth, true
}

// private ========================================================

func support(s *Shape, ax0, ax1, cosR, sinR float32) float32 {
	// support returns how far shape s extends along axis (ax0, ax1).
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	var dX = number.Absolute(ax0*cosR + ax1*sinR)
	var dY = number.Absolute(-ax0*sinR + ax1*cosR)
	return (hx-r)*dX + (hy-r)*dY + r
}
func nearestInnerBoxPoint(s *Shape, px, py, cosR, sinR float32) (float32, float32) {
	// nearestInnerBoxPoint returns the closest point on s's inner box to world point (px, py).
	var rx, ry = px - s.X, py - s.Y
	var lx = rx*cosR + ry*sinR
	var ly = -rx*sinR + ry*cosR

	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.Roundness * min(hx, hy)
	lx = max(-(hx - r), min(hx-r, lx))
	ly = max(-(hy - r), min(hy-r, ly))

	return s.X + lx*cosR - ly*sinR, s.Y + lx*sinR + ly*cosR
}
