package geometry

import (
	"big-black-box/internal"
	"big-black-box/utility/number"
)

func (s Shape) Contains(shape Shape) bool {
	if s.IsPoint() {
		return s.pointContains(shape)
	}
	if s.IsLineSegment() {
		return s.lineSegmentContains(shape)
	}
	if s.IsInfiniteRay() {
		return s.infiniteRayContains(shape)
	}
	if s.IsInfinitePlane() {
		return s.infinitePlaneContains(shape)
	}
	if s.IsCircle() {
		return s.circleContains(shape)
	}
	if s.IsRectangle() {
		return s.rectangleContains(shape)
	}

	return false
}

// private ========================================================

func (s Shape) pointContains(shape Shape) bool {
	var x, y = s.Position()
	var sx, sy = shape.Position()
	if shape.IsPoint() {
		return number.IsWithin(x, sx, 0.001) && number.IsWithin(y, sy, 0.001)
	}
	if shape.IsLineSegment() {
		var bx, by = shape.Point2()
		var crossProduct = (y-sy)*(bx-sx) - (x-sx)*(by-sy)
		if crossProduct != 0 {
			return false
		}
		return x >= min(sx, bx) && x <= max(sx, bx) && y >= min(sy, by) && y <= max(sy, by)
	}
	if shape.IsInfiniteRay() {
		var dy, dx = internal.SinCos(shape.Angle())
		var vx, vy = x - sx, y - sy
		var crossProduct = vx*dy - vy*dx
		if number.Absolute(crossProduct) > 1e-9 {
			return false
		}
		var dotProduct = vx*dx + vy*dy
		return dotProduct >= 0 // the dot product must be positive for the point to be "in front" of the origin
	}
	if shape.IsInfinitePlane() {
		var dy, dx = internal.SinCos(shape.Angle())
		var vx, vy = x - sx, y - sy
		var crossProduct = vx*dy - vy*dx
		return number.Absolute(crossProduct) <= 1e-9
	}
	if shape.IsCircle() {
		var radius, _ = shape.Size()
		var distX, distY = x - sx, y - sy
		return distX*distX+distY*distY <= radius*radius
	}
	if shape.IsRectangle() {
		var dy, dx = internal.SinCos(shape.Angle())
		var vx, vy = x - sx, y - sy
		var localX = vx*dx + vy*dy
		var localY = -vx*dy + vy*dx
		var w, h = shape.Size()
		return number.Absolute(localX) <= w/2 && number.Absolute(localY) <= h/2
	}
	if shape.IsEllipse() {
		var dy, dx = internal.SinCos(shape.Angle())
		var vx, vy = x - sx, y - sy
		var localX = vx*dx + vy*dy
		var localY = -vx*dy + vy*dx
		var w, h = shape.Size()
		var nx, ny = localX / (w / 2), localY / (h / 2)
		return nx*nx+ny*ny <= 1
	}
	if shape.IsCapsule() {
		var dy, dx = internal.SinCos(shape.Angle())
		var vx, vy = x - sx, y - sy
		var localX = vx*dx + vy*dy
		var localY = -vx*dy + vy*dx
		var w, h = shape.Size()
		var halfLen = number.Limit(w/2-h/2, float32(0), w/2)   // inner half-length of the rectangular body
		var nearestX = number.Limit(localX, -halfLen, halfLen) // nearest point on the centerline
		var distX, distY = localX - nearestX, localY
		return distX*distX+distY*distY <= (h/2)*(h/2)
	}
	return false
}
func (s Shape) lineSegmentContains(shape Shape) bool {
	var ax, ay = s.Position()
	var bx, by = s.Point2()
	var sx, sy = shape.Position()
	if shape.IsPoint() {
		var crossProduct = (sx-ax)*(by-ay) - (sy-ay)*(bx-ax)
		if number.Absolute(crossProduct) > 1e-9 {
			return false
		}
		return sx >= min(ax, bx) && sx <= max(ax, bx) && sy >= min(ay, by) && sy <= max(ay, by)
	}
	if shape.IsLineSegment() {
		var ex, ey = shape.Point2()
		return s.lineSegmentContains(Point(sx, sy)) && s.lineSegmentContains(Point(ex, ey))
	}
	return false // other shapes return false by default
}
func (s Shape) infiniteRayContains(shape Shape) bool {
	var dy, dx = internal.SinCos(s.Angle())
	var ax, ay = s.Position()
	var sx, sy = shape.Position()
	if shape.IsPoint() {
		var vx, vy = sx - ax, sy - ay
		var crossProduct = vx*dy - vy*dx
		if number.Absolute(crossProduct) > 1e-9 {
			return false
		}
		var dotProduct = vx*dx + vy*dy
		return dotProduct >= 0
	}
	if shape.IsLineSegment() {
		var ex, ey = shape.Point2()
		return s.infiniteRayContains(Point(sx, sy)) && s.infiniteRayContains(Point(ex, ey))
	}
	if shape.IsInfiniteRay() {
		var bdy, bdx = internal.SinCos(shape.Angle())
		var dirCross = bdx*dy - bdy*dx // directions must be parallel and same orientation
		if number.Absolute(dirCross) > 1e-9 {
			return false
		}
		if bdx*dx+bdy*dy < 0 { // opposite direction
			return false
		}
		// shape's origin must lie on s (collinear and at or past s's origin)
		var vx, vy = sx - ax, sy - ay
		var crossProduct = vx*dy - vy*dx
		if number.Absolute(crossProduct) > 1e-9 {
			return false
		}
		var dotProduct = vx*dx + vy*dy
		return dotProduct >= 0
	}
	return false
}
func (s Shape) infinitePlaneContains(shape Shape) bool {
	var dy, dx = internal.SinCos(s.Angle())
	var ax, ay = s.Position()
	var sx, sy = shape.Position()
	if shape.IsPoint() {
		var vx, vy = sx - ax, sy - ay
		var crossProduct = vx*dy - vy*dx
		return number.Absolute(crossProduct) <= 1e-9
	}
	if shape.IsLineSegment() {
		var ex, ey = shape.Point2()
		return s.infinitePlaneContains(Point(sx, sy)) && s.infinitePlaneContains(Point(ex, ey))
	}
	if shape.IsInfiniteRay() {
		var bdy, bdx = internal.SinCos(shape.Angle())
		var dirCross = bdx*dy - bdy*dx // direction must be parallel (same or opposite is fine, plane has no orientation)
		if number.Absolute(dirCross) > 1e-9 {
			return false
		}
		// origin must lie on s
		var vx, vy = sx - ax, sy - ay
		var crossProduct = vx*dy - vy*dx
		return number.Absolute(crossProduct) <= 1e-9
	}
	if shape.IsInfinitePlane() {
		var bdy, bdx = internal.SinCos(shape.Angle())
		var dirCross = bdx*dy - bdy*dx // direction must be parallel (same or opposite is fine, plane has no orientation)
		if number.Absolute(dirCross) > 1e-9 {
			return false
		}
		// origin must lie on s
		var vx, vy = sx - ax, sy - ay
		var crossProduct = vx*dy - vy*dx
		return number.Absolute(crossProduct) <= 1e-9
	}
	return false
}
func (s Shape) circleContains(shape Shape) bool {
	var cx, cy = s.Position()
	var r, _ = s.Size()
	var sx, sy = shape.Position()
	if shape.IsPoint() {
		var dx, dy = sx - cx, sy - cy
		return dx*dx+dy*dy <= r*r
	}
	if shape.IsLineSegment() {
		var ex, ey = shape.Point2()
		return s.circleContains(Point(sx, sy)) && s.circleContains(Point(ex, ey))
	}
	if shape.IsCircle() {
		var r2, _ = shape.Size()
		var dx, dy = sx - cx, sy - cy
		return number.SquareRoot(dx*dx+dy*dy)+r2 <= r
	}
	if shape.IsRectangle() {
		var w, h = shape.Size()
		var sin, cos = internal.SinCos(shape.Angle())
		var hw, hh = w / 2, h / 2
		for _, c := range [4][2]float32{{hw, hh}, {-hw, hh}, {hw, -hh}, {-hw, -hh}} {
			if !s.circleContains(Point(sx+c[0]*cos-c[1]*sin, sy+c[0]*sin+c[1]*cos)) {
				return false
			}
		}
		return true
	}
	if shape.IsEllipse() {
		// transform circle center into ellipse's local frame
		var edy, edx = internal.SinCos(shape.Angle())
		var vx, vy = sx - cx, sy - cy
		var u = vx*edx + vy*edy  // circle center in ellipse local frame (along major axis)
		var v = -vx*edy + vy*edx // circle center in ellipse local frame (along minor axis)
		var a, b = shape.Size()
		a, b = a/2, b/2
		// maximize f(t) = (a·cos t - u)² + (b·sin t - v)² using Newton's method (t in degrees)
		// f'_rad  = (b²-a²)·sin(2t) + 2au·sin t - 2bv·cos t
		// f''_rad = 2(b²-a²)·cos(2t) + 2au·cos t + 2bv·sin t
		// newton step in degrees: Δt = -(f'_rad / f''_rad) · (180/π)
		const radToDeg = float32(57.295779513)
		var maxD2 = float32(0)
		for _, t0 := range [4]float32{0, 90, 180, 270} {
			var t = t0
			for range 10 {
				var sn, cs = internal.SinCos(t)
				var sn2, cs2 = internal.SinCos(2 * t)
				var fp = (b*b-a*a)*sn2 + 2*a*u*sn - 2*b*v*cs
				var fpp = 2*(b*b-a*a)*cs2 + 2*a*u*cs + 2*b*v*sn
				if number.Absolute(fpp) < 1e-6 {
					break
				}
				t -= fp / fpp * radToDeg
			}
			var sn, cs = internal.SinCos(t)
			var d2 = (a*cs-u)*(a*cs-u) + (b*sn-v)*(b*sn-v)
			maxD2 = max(maxD2, d2)
		}
		return maxD2 <= r*r
	}
	if shape.IsCapsule() {
		// a capsule's farthest point from any external location is on one of the two end-cap circles
		var w, h = shape.Size()
		var capRadius = h / 2
		var halfLen = number.Limit(w/2-h/2, float32(0), w/2)
		var cdy, cdx = internal.SinCos(shape.Angle())
		var e1x, e1y = sx + halfLen*cdx, sy + halfLen*cdy
		var e2x, e2y = sx - halfLen*cdx, sy - halfLen*cdy
		var d1x, d1y = e1x - cx, e1y - cy
		var d2x, d2y = e2x - cx, e2y - cy
		return number.SquareRoot(d1x*d1x+d1y*d1y)+capRadius <= r &&
			number.SquareRoot(d2x*d2x+d2y*d2y)+capRadius <= r
	}
	return false
}
func (s Shape) rectangleContains(shape Shape) bool {
	var rdy, rdx = internal.SinCos(s.Angle())
	var rx, ry = s.Position()
	var w, h = s.Size()
	var sx, sy = shape.Position()
	var vx, vy = sx - rx, sy - ry
	var lx = vx*rdx + vy*rdy // shape center in rectangle local frame
	var ly = -vx*rdy + vy*rdx
	if shape.IsPoint() {
		return number.Absolute(lx) <= w/2 && number.Absolute(ly) <= h/2
	}
	if shape.IsLineSegment() {
		var ex, ey = shape.Point2()
		return s.rectangleContains(Point(sx, sy)) && s.rectangleContains(Point(ex, ey))
	}
	if shape.IsCircle() {
		var r2, _ = shape.Size()
		return number.Absolute(lx)+r2 <= w/2 && number.Absolute(ly)+r2 <= h/2
	}
	if shape.IsRectangle() {
		var w2, h2 = shape.Size()
		var edy, edx = internal.SinCos(shape.Angle())
		var cp = edx*rdx + edy*rdy  // cos of angle between rectangles
		var sp = -edx*rdy + edy*rdx // sin of angle between rectangles
		return number.Absolute(lx)+w2/2*number.Absolute(cp)+h2/2*number.Absolute(sp) <= w/2 &&
			number.Absolute(ly)+w2/2*number.Absolute(sp)+h2/2*number.Absolute(cp) <= h/2
	}
	if shape.IsEllipse() {
		var a, b = shape.Size()
		a, b = a/2, b/2
		var edy, edx = internal.SinCos(shape.Angle())
		var cp = edx*rdx + edy*rdy                                  // cos of angle between ellipse and rectangle
		var sp = -edx*rdy + edy*rdx                                 // sin of angle between ellipse and rectangle
		var extX = number.SquareRoot((a*cp)*(a*cp) + (b*sp)*(b*sp)) // ellipse reach along rectangle's x-axis
		var extY = number.SquareRoot((a*sp)*(a*sp) + (b*cp)*(b*cp)) // ellipse reach along rectangle's y-axis
		return number.Absolute(lx)+extX <= w/2 && number.Absolute(ly)+extY <= h/2
	}
	if shape.IsCapsule() {
		var w2, h2 = shape.Size()
		var capRadius = h2 / 2
		var halfLen = number.Limit(w2/2-h2/2, float32(0), w2/2)
		var cdy, cdx = internal.SinCos(shape.Angle())
		var capDirX = cdx*rdx + cdy*rdy // capsule axis direction in rectangle local frame
		var capDirY = -cdx*rdy + cdy*rdx
		var c1x, c1y = lx + halfLen*capDirX, ly + halfLen*capDirY // cap centers in local frame
		var c2x, c2y = lx - halfLen*capDirX, ly - halfLen*capDirY
		return number.Absolute(c1x)+capRadius <= w/2 && number.Absolute(c1y)+capRadius <= h/2 &&
			number.Absolute(c2x)+capRadius <= w/2 && number.Absolute(c2y)+capRadius <= h/2
	}
	return false
}
