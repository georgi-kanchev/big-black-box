package geometry

import (
	"big-black-box/internal"
	"big-black-box/utility/angle"
	"big-black-box/utility/number"
	"big-black-box/utility/point"
)

// Types of shapes:
//
//	[x, y, NaN, NaN, NaN]                       // Point
//	[x, y, x2, y2, NegInf]                      // Line Segment
//	[x, y, PosInf, NaN, angle]                  // Infinite Ray
//	[x, y, NaN, NaN, angle]                	    // Infinite Plane
//	[centerX, centerY, width, height, angle]    // Rectangle
//	[centerX, centerY, -width, height, angle]   // Ellipse
//	[centerX, centerY, width, -height, angle]   // Capsule
//	[centerX, centerY, -radius, -spread, angle] // Cone
//	[centerX, centerY, radius, radius, PosInf]  // Circle
type Shape [5]float32

func Point(x, y float32) Shape {
	return [5]float32{x, y, number.NaN(), number.NaN(), number.NaN()}
}
func LineSegment(x, y, x2, y2 float32) Shape {
	return [5]float32{x, y, x2, y2, number.NegativeInfinity()}
}
func InfiniteRay(x, y, angle float32) Shape {
	return [5]float32{x, y, number.PositiveInfinity(), number.NaN(), angle}
}
func InfinitePlane(x, y, angle float32) Shape {
	return [5]float32{x, y, number.NaN(), number.NaN(), angle}
}
func Rectangle(centerX, centerY, width, height, angle float32) Shape {
	return [5]float32{centerX, centerY, width, height, angle}
}
func Ellipse(centerX, centerY, width, height, angle float32) Shape {
	return [5]float32{centerX, centerY, -width, height, angle}
}
func Capsule(centerX, centerY, width, height, angle float32) Shape {
	return [5]float32{centerX, centerY, width, -height, angle}
}
func Cone(centerX, centerY, radius, spread, angle float32) Shape {
	return [5]float32{centerX, centerY, -radius, -spread, angle}
}
func Circle(centerX, centerY, radius float32) Shape {
	return [5]float32{centerX, centerY, radius, radius, number.PositiveInfinity()}
}

//=================================================================

func (s Shape) IsPoint() bool {
	return number.IsNaN(s[2]) && number.IsNaN(s[3]) && number.IsNaN(s[4])
}
func (s Shape) IsLineSegment() bool {
	return number.IsNegativeInfinity(s[4])
}
func (s Shape) IsInfiniteRay() bool {
	return number.IsPositiveInfinity(s[2]) && number.IsNaN(s[3]) && !number.IsNaN(s[4]) && !number.IsInfinity(s[4])
}
func (s Shape) IsInfinitePlane() bool {
	return number.IsNaN(s[2]) && number.IsNaN(s[3]) && !number.IsNaN(s[4]) && !number.IsInfinity(s[4])
}
func (s Shape) IsRectangle() bool {
	return s[2] >= 0 && s[3] >= 0 && !number.IsNaN(s[4]) && !number.IsInfinity(s[4])
}
func (s Shape) IsEllipse() bool {
	return s[2] < 0 && s[3] >= 0 && !number.IsInfinity(s[4])
}
func (s Shape) IsCapsule() bool {
	return s[2] >= 0 && s[3] < 0 && !number.IsInfinity(s[4])
}
func (s Shape) IsCone() bool {
	return s[2] < 0 && s[3] < 0 && !number.IsInfinity(s[4])
}
func (s Shape) IsCircle() bool {
	return number.IsInfinity(s[4])
}

//=================================================================

func (s Shape) Position() (x, y float32) {
	return s[0], s[1]
}
func (s Shape) Point2() (x2, y2 float32) {
	if s.IsLineSegment() {
		return s[2], s[3]
	}
	return number.NaN(), number.NaN()
}
func (s Shape) Size() (width, height float32) {
	width, height = number.NaN(), number.NaN()

	if s.IsCircle() {
		return number.Absolute(s[2]), number.Absolute(s[3])
	}

	if s.IsLineSegment() {
		width = point.DistanceToPoint(s[0], s[1], s[2], s[3])
	} else if !s.IsPoint() && !s.IsInfinitePlane() {
		width = number.Absolute(s[2])
	}
	if !s.IsPoint() && !s.IsLineSegment() && !s.IsInfiniteRay() && !s.IsInfinitePlane() {
		height = number.Absolute(s[3])
	}
	return width, height
}
func (s Shape) Angle() float32 {
	if s.IsLineSegment() {
		return angle.BetweenPoints(s[0], s[1], s[2], s[3])
	}
	if s.IsPoint() || s.IsCircle() {
		return number.NaN()
	}
	return s[4]
}

//=================================================================

func (s *Shape) SetPosition(x, y float32) {
	s[0], s[1] = x, y
}
func (s *Shape) SetPoint2(x2, y2 float32) {
	if s.IsLineSegment() {
		s[2], s[3] = x2, y2
	}
}
func (s *Shape) SetSize(width, height float32) {
	width = number.Absolute(width)
	height = number.Absolute(height)

	if s.IsCircle() && width == height {
		s[2], s[3] = width, height
		return
	}

	if s.IsEllipse() || s.IsCone() {
		s[2] = -width
	} else if !s.IsPoint() && !s.IsInfinitePlane() {
		s[2] = width
	}
	if s.IsCapsule() || s.IsCone() {
		s[3] = -height
	} else if !s.IsPoint() && !s.IsLineSegment() && !s.IsInfiniteRay() && !s.IsInfinitePlane() {
		s[3] = height
	}
}
func (s *Shape) SetAngle(angle float32) {
	if s.IsLineSegment() {
		var length = point.DistanceToPoint(s[0], s[1], s[2], s[3])
		var sin, cos = internal.SinCos(angle)
		s[2] = s[0] + cos*length
		s[3] = s[1] + sin*length
		return
	}
	if !s.IsPoint() && !s.IsCircle() && !number.IsInfinity(angle) && !number.IsNaN(angle) {
		s[4] = angle // protect the tags for Points and Circles
	}
}

//=================================================================

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
	if s.IsEllipse() {
		return s.ellipseContains(shape)
	}
	if s.IsCapsule() {
		return s.capsuleContains(shape)
	}

	return false
}
