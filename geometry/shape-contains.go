package geometry

import (
	"big-black-box/internal"
	"big-black-box/utility/number"
	"big-black-box/utility/point"
)

func (s Shape) pointContains(shape Shape) bool {
	if shape.IsPoint() {
		var x, y = s.Center()
		var sx, sy = shape.Center()
		return number.IsWithin(x, sx, 0.001) && number.IsWithin(y, sy, 0.001)
	}
	return false
}
func (s Shape) lineContains(shape Shape) bool {
	// A Line (Capsule) contains another shape if the furthest distance
	// from any point in 'shape' to the line segment AB is <= s.Roundness()
	var r1 = s.Roundness()

	if shape.IsPoint() {
		var px, py = shape.Center()
		return distToSegment(px, py, s) <= r1
	}

	if shape.IsCircle() {
		var cx, cy = shape.Center()
		var r2, _ = shape.Size() // For circles, this is the radius
		// The circle is inside if its center is close enough that
		// its entire radius fits within the capsule's radius.
		return distToSegment(cx, cy, s)+r2 <= r1
	}

	if shape.IsCapsule() {
		// A capsule is inside another if its two end-points (plus its radius)
		// are both within the parent capsule.
		var ax, ay = shape.A()
		var bx, by = shape.B()
		var r2 = shape.Roundness()

		return (distToSegment(ax, ay, s)+r2 <= r1) && (distToSegment(bx, by, s)+r2 <= r1)
	}

	if shape.IsRectangle() {
		var r1 = s.Roundness()
		var r2 = shape.Roundness()

		// If the rectangle's internal thickness is greater than the capsule's,
		// it can't possibly fit inside.
		if r2 > r1 {
			return false
		}

		// Get the 4 corners of the inner rectangle.
		// Note: We use the dimensions MINUS the roundness, because the
		// roundness extends outward from this core box.
		var w, h = shape.Size()
		w = (w / 2) - r2
		h = (h / 2) - r2
		if w < 0 {
			w = 0
		}
		if h < 0 {
			h = 0
		}

		var cx, cy = shape.Center()
		var sinA, cosA = internal.SinCos(shape.Angle())

		// The 4 offsets for an axis-aligned box
		var offsets = [4][2]float32{{w, h}, {w, -h}, {-w, h}, {-w, -h}}

		for _, off := range offsets {
			// Rotate the corner by the rectangle's angle
			var px = cx + (off[0]*cosA - off[1]*sinA)
			var py = cy + (off[0]*sinA + off[1]*cosA)

			// Check if this corner (expanded by the rect's own roundness)
			// fits inside the capsule's radius.
			if distToSegment(px, py, s)+r2 > r1 {
				return false
			}
		}

		return true
	}
	return false
}
func (s Shape) rectangleContains(shape Shape) bool {
	return false
}
func (s Shape) coneContains(shape Shape) bool {
	return false
}

//=================================================================

func distToSegment(px, py float32, s Shape) float32 {
	var ax, ay = s.A()
	var bx, by = s.B()
	var dx = bx - ax
	var dy = by - ay
	var l2 = dx*dx + dy*dy
	if l2 == 0 {
		return point.DistanceToPoint(px, py, ax, ay)
	}

	var t = ((px-ax)*dx + (py-ay)*dy) / l2 // Project point onto the line, clamping to the [0, 1] range of the segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	var projX = ax + t*dx
	var projY = ay + t*dy
	return point.DistanceToPoint(px, py, projX, projY)
}
