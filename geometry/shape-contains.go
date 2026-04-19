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
