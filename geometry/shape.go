package geometry

import (
	"big-black-box/utility/angle"
	"big-black-box/utility/number"
	"big-black-box/utility/point"
)

// Types of shapes:
//
//	[x, y, NaN, NaN, NaN]                              // Point
//	[x, y, x2, y2, NegInf]                             // Line
//	[centerX, centerY, width, height, angle+roundness] // Rectangle
//	[centerX, centerY, -width, height, angle]          // Capsule
//	[centerX, centerY, radius, -spread, angle]         // Cone
//	[centerX, centerY, radius, radius, PosInf]         // Circle
type Shape [5]float32

func Point(x, y float32) Shape {
	return [5]float32{x, y, number.NaN(), number.NaN(), number.NaN()}
}
func Line(x, y, x2, y2 float32) Shape {
	return [5]float32{x, y, x2, y2, number.NegativeInfinity()}
}
func Rectangle(centerX, centerY, width, height, angle float32) Shape {
	return [5]float32{centerX, centerY, width, height, rrPack(angle, 0)}
}
func RoundedRectangle(centerX, centerY, width, height, angle, cornerRadius float32) Shape {
	return [5]float32{centerX, centerY, width, height, rrPack(angle, cornerRadius)}
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
func (s Shape) IsLine() bool {
	return number.IsNegativeInfinity(s[4])
}
func (s Shape) IsRectangle() bool {
	return s[2] >= 0 && s[3] >= 0 && s[4] <= 0 && !number.IsNegativeInfinity(s[4])
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
	if s.IsLine() {
		return s[2], s[3]
	}
	return number.NaN(), number.NaN()
}
func (s Shape) Size() (width, height float32) {
	width, height = number.NaN(), number.NaN()

	if s.IsCircle() {
		return number.Absolute(s[2]), number.Absolute(s[3])
	}

	if s.IsLine() {
		width = point.DistanceToPoint(s[0], s[1], s[2], s[3])
	} else if !s.IsPoint() {
		width = number.Absolute(s[2])
	}
	if !s.IsPoint() && !s.IsLine() {
		height = number.Absolute(s[3])
	}
	return width, height
}
func (s Shape) Angle() float32 {
	if s.IsLine() {
		return angle.BetweenPoints(s[0], s[1], s[2], s[3])
	}
	if s.IsPoint() || s.IsCircle() {
		return number.NaN()
	}
	if s.IsRectangle() {
		var idx, _ = rrUnpack(s[4])
		return float32(idx) * 0.1
	}
	return s[4]
}
func (s Shape) CornerRadius() float32 {
	if s.IsRectangle() {
		var _, cr = rrUnpack(s[4])
		return float32(cr)
	}
	return number.NaN()
}

//=================================================================

func (s Shape) Contains(shape Shape) bool {
	if s.IsPoint() {
		return s.pointContains(shape)
	}
	if s.IsLine() {
		return s.lineSegmentContains(shape)
	}
	if s.IsCircle() {
		return s.circleContains(shape)
	}
	if s.IsRectangle() {
		return s.rectangleContains(shape)
	}
	if s.IsCapsule() {
		return s.capsuleContains(shape)
	}

	return false
}

// private ========================================================

func rrPack(angle, cornerRadius float32) float32 {
	// rrPack encodes angle (0.1° precision) and cornerRadius (integer, max 4659) into a single negative float32.
	// Encoding: -(crInt * 3600 + angleIdx) where angleIdx = round(angle*10) mod 3600.
	var idx = ((int(angle*10) % 3600) + 3600) % 3600
	var cr = min(int(number.Absolute(cornerRadius)), 4659)
	return -float32(cr*3600 + idx)
}
func rrUnpack(packed float32) (angleIdx, crInt int) {
	var abs = int(-packed)
	return abs % 3600, abs / 3600
}
