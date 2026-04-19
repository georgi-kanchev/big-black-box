package geometry

import (
	"big-black-box/internal"
	"big-black-box/utility/number"
)

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
func (s Shape) ellipseContains(shape Shape) bool {
	var edy, edx = internal.SinCos(s.Angle())
	var ox, oy = s.Position()
	var a, b = s.Size()
	a, b = a/2, b/2
	var sx, sy = shape.Position()
	var vx, vy = sx - ox, sy - oy
	var lx = vx*edx + vy*edy // shape center in ellipse local frame
	var ly = -vx*edy + vy*edx
	if shape.IsPoint() {
		var nx, ny = lx / a, ly / b
		return nx*nx+ny*ny <= 1
	}
	if shape.IsLineSegment() {
		var p2x, p2y = shape.Point2()
		return s.ellipseContains(Point(sx, sy)) && s.ellipseContains(Point(p2x, p2y))
	}
	if shape.IsCircle() {
		var r, _ = shape.Size()
		// check if the circle (lx,ly) radius r is inside ellipse (a,b):
		// maximize ((lx + r·cos t)/a)² + ((ly + r·sin t)/b)² in normalized space
		return maxNormSqOnCurve(lx/a, ly/b, r/a, 0, 0, r/b) <= 1
	}
	if shape.IsRectangle() { // all 4 corners must be inside the ellipse (exact by convexity)
		var w2, h2 = shape.Size()
		var rdy2, rdx2 = internal.SinCos(shape.Angle())
		var cp = rdx2*edx + rdy2*edy  // cos of relative angle
		var sp = -rdx2*edy + rdy2*edx // sin of relative angle
		var hw2, hh2 = w2 / 2, h2 / 2
		for _, c := range [4][2]float32{{hw2, hh2}, {-hw2, hh2}, {hw2, -hh2}, {-hw2, -hh2}} {
			var cx2 = lx + c[0]*cp - c[1]*sp
			var cy2 = ly + c[0]*sp + c[1]*cp
			var nx, ny = cx2 / a, cy2 / b
			if nx*nx+ny*ny > 1 {
				return false
			}
		}
		return true
	}
	if shape.IsEllipse() {
		var a2, b2 = shape.Size()
		a2, b2 = a2/2, b2/2
		var edy2, edx2 = internal.SinCos(shape.Angle())
		var cp = edx2*edx + edy2*edy  // cos of relative angle
		var sp = -edx2*edy + edy2*edx // sin of relative angle
		// inner ellipse in outer's normalized space traces: (uc + P·cos t + Q·sin t, vc + R·cos t + S·sin t)
		var P = a2 * cp / a
		var Q = -b2 * sp / a
		var R = a2 * sp / b
		var S = b2 * cp / b
		return maxNormSqOnCurve(lx/a, ly/b, P, Q, R, S) <= 1
	}
	if shape.IsCapsule() { // the capsule is inside the ellipse if both end-cap circles are inside the ellipse
		var w2, h2 = shape.Size()
		var capRadius = h2 / 2
		var halfLen = number.Limit(w2/2-h2/2, float32(0), w2/2)
		var cdy2, cdx2 = internal.SinCos(shape.Angle())
		var capDirX = cdx2*edx + cdy2*edy // capsule axis direction in ellipse local frame
		var capDirY = -cdx2*edy + cdy2*edx
		var c1x, c1y = lx + halfLen*capDirX, ly + halfLen*capDirY
		var c2x, c2y = lx - halfLen*capDirX, ly - halfLen*capDirY
		return maxNormSqOnCurve(c1x/a, c1y/b, capRadius/a, 0, 0, capRadius/b) <= 1 &&
			maxNormSqOnCurve(c2x/a, c2y/b, capRadius/a, 0, 0, capRadius/b) <= 1
	}
	return false
}
func (s Shape) capsuleContains(shape Shape) bool {
	var cdy, cdx = internal.SinCos(s.Angle())
	var ox, oy = s.Position()
	var w, h = s.Size()
	var capRadius = h / 2
	var halfLen = number.Limit(w/2-h/2, float32(0), w/2)
	var sx, sy = shape.Position()
	var vx, vy = sx - ox, sy - oy
	var lx = vx*cdx + vy*cdy // shape center in capsule local frame
	var ly = -vx*cdy + vy*cdx
	if shape.IsPoint() {
		var nearestX = number.Limit(lx, -halfLen, halfLen)
		var dx, dy = lx - nearestX, ly
		return dx*dx+dy*dy <= capRadius*capRadius
	}
	if shape.IsLineSegment() {
		var p2x, p2y = shape.Point2()
		return s.capsuleContains(Point(sx, sy)) && s.capsuleContains(Point(p2x, p2y))
	}
	if shape.IsCircle() {
		var r, _ = shape.Size()
		var nearestX = number.Limit(lx, -halfLen, halfLen)
		var dx, dy = lx - nearestX, ly
		return number.SquareRoot(dx*dx+dy*dy)+r <= capRadius
	}
	if shape.IsRectangle() { // check all 4 corners (exact by convexity)
		var w2, h2 = shape.Size()
		var rdy2, rdx2 = internal.SinCos(shape.Angle())
		var cp = rdx2*cdx + rdy2*cdy
		var sp = -rdx2*cdy + rdy2*cdx
		var hw2, hh2 = w2 / 2, h2 / 2
		for _, c := range [4][2]float32{{hw2, hh2}, {-hw2, hh2}, {hw2, -hh2}, {-hw2, -hh2}} {
			var cx2 = lx + c[0]*cp - c[1]*sp
			var cy2 = ly + c[0]*sp + c[1]*cp
			var nearestX = number.Limit(cx2, -halfLen, halfLen)
			var dx, dy = cx2 - nearestX, cy2
			if dx*dx+dy*dy > capRadius*capRadius {
				return false
			}
		}
		return true
	}
	if shape.IsEllipse() {
		var a2, b2 = shape.Size() // ellipse boundary in capsule local frame: (lx + A·cos t + B·sin t, ly + C·cos t + D·sin t)
		a2, b2 = a2/2, b2/2       // maximize dist²(P, segment) = (Px − clamp(Px, −halfLen, halfLen))² + Py²
		var edy2, edx2 = internal.SinCos(shape.Angle())
		var cp = edx2*cdx + edy2*cdy
		var sp = -edx2*cdy + edy2*cdx
		var A, B, C, D = a2 * cp, -b2 * sp, a2 * sp, b2 * cp
		const radToDeg = float32(57.295779513)
		var maxF = float32(0)
		for _, t0 := range [8]float32{0, 45, 90, 135, 180, 225, 270, 315} {
			var t = t0
			for range 10 {
				var sn, cs = internal.SinCos(t)
				var Px = lx + A*cs + B*sn
				var Py = ly + C*cs + D*sn
				var Pxp = -A*sn + B*cs
				var Pyp = -C*sn + D*cs
				var dx, dxp float32
				if Px > halfLen {
					dx, dxp = Px-halfLen, Pxp
				} else if Px < -halfLen {
					dx, dxp = Px+halfLen, Pxp
				}
				var Pxpp = -A*cs - B*sn
				var Pypp = -C*cs - D*sn
				var dxpp float32
				if Px > halfLen || Px < -halfLen {
					dxpp = Pxpp
				}
				var fp = dx*dxp + Py*Pyp
				var fpp = dxp*dxp + dx*dxpp + Pyp*Pyp + Py*Pypp
				if number.Absolute(fpp) < 1e-6 {
					break
				}
				t -= fp / fpp * radToDeg
			}
			var sn, cs = internal.SinCos(t)
			var Px = lx + A*cs + B*sn
			var Py = ly + C*cs + D*sn
			var nearestX = number.Limit(Px, -halfLen, halfLen)
			var dx, dy = Px - nearestX, Py
			maxF = max(maxF, dx*dx+dy*dy)
		}
		return maxF <= capRadius*capRadius
	}
	if shape.IsCapsule() { // inner capsule ⊆ outer capsule iff both inner cap circles fit inside the outer capsule
		var w2, h2 = shape.Size()
		var innerCapRadius = h2 / 2
		var innerHalfLen = number.Limit(w2/2-h2/2, float32(0), w2/2)
		var c2dy, c2dx = internal.SinCos(shape.Angle())
		var capDirX = c2dx*cdx + c2dy*cdy // inner capsule axis in outer local frame
		var capDirY = -c2dx*cdy + c2dy*cdx
		var c1x, c1y = lx + innerHalfLen*capDirX, ly + innerHalfLen*capDirY
		var c2x, c2y = lx - innerHalfLen*capDirX, ly - innerHalfLen*capDirY
		for _, cap := range [2][2]float32{{c1x, c1y}, {c2x, c2y}} {
			var nearestX = number.Limit(cap[0], -halfLen, halfLen)
			var dx, dy = cap[0] - nearestX, cap[1]
			if number.SquareRoot(dx*dx+dy*dy)+innerCapRadius > capRadius {
				return false
			}
		}
		return true
	}
	return false
}

// maxNormSqOnCurve finds max ||(uc + P·cos t + Q·sin t, vc + R·cos t + S·sin t)||²
// using Newton's method from 4 starting angles (t in degrees throughout)
func maxNormSqOnCurve(uc, vc, P, Q, R, S float32) float32 {
	const radToDeg = float32(57.295779513)
	var maxF = float32(0)
	for _, t0 := range [4]float32{0, 90, 180, 270} {
		var t = t0
		for range 10 {
			var sn, cs = internal.SinCos(t)
			var X = uc + P*cs + Q*sn
			var Y = vc + R*cs + S*sn
			var Xp = -P*sn + Q*cs
			var Yp = -R*sn + S*cs
			var fp = X*Xp + Y*Yp                                      // ½ · d/dt_rad [X² + Y²]
			var fpp = Xp*Xp + X*(-P*cs-Q*sn) + Yp*Yp + Y*(-R*cs-S*sn) // ½ · d²/dt_rad² [X² + Y²]
			if number.Absolute(fpp) < 1e-6 {
				break
			}
			t -= fp / fpp * radToDeg
		}
		var sn, cs = internal.SinCos(t)
		var X = uc + P*cs + Q*sn
		var Y = vc + R*cs + S*sn
		maxF = max(maxF, X*X+Y*Y)
	}
	return maxF
}
