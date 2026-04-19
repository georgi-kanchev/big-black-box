package geometry

import (
	"big-black-box/internal"
	"big-black-box/utility/angle"
	"big-black-box/utility/number"
	"big-black-box/utility/point"
)

// Types of shapes:
//
//	[x, y, NaN, NaN, NaN]                               // Point
//	[x, y, x2, y2, roundness]                           // Line + Capsule
//	[centerX, centerY, -width, height, angle+roundness] // Rectangle + Rounded
//	[centerX, centerY, radius, -spread, angle]          // Cone + Circle
type Shape [5]float32

func Point(x, y float32) Shape {
	return [5]float32{x, y, number.NaN(), number.NaN(), number.NaN()}
}
func Line(ax, ay, bx, by float32) Shape {
	return [5]float32{ax, ay, bx, by, number.NaN()}
}
func LineRounded(ax, ay, bx, by, roundness float32) Shape {
	return [5]float32{ax, ay, bx, by, number.Absolute(roundness)}
}
func Rectangle(x, y, width, height, angle float32) Shape {
	return [5]float32{x, y, -number.Absolute(width), number.Absolute(height), rrPack(angle, 0)}
}
func RectangleRounded(x, y, width, height, angle, roundness float32) Shape {
	return [5]float32{x, y, -number.Absolute(width), number.Absolute(height), rrPack(angle, roundness)}
}
func Capsule(x, y, width, height, angle float32) Shape {
	return [5]float32{x, y, number.Absolute(width), number.Absolute(height), rrPack(angle, 0)}
}
func Cone(x, y, radius, spread, angle float32) Shape {
	var s = number.Absolute(spread)
	if s >= 360 {
		return Circle(x, y, radius)
	}
	return [5]float32{x, y, number.Absolute(radius), -s, angle}
}
func Circle(x, y, radius float32) Shape {
	return [5]float32{x, y, number.Absolute(radius), -360, number.NaN()}
}

//=================================================================

func (s Shape) IsPoint() bool {
	return number.IsNaN(s[2]) && number.IsNaN(s[3]) && number.IsNaN(s[4])
}
func (s Shape) IsLine() bool {
	return !number.IsNaN(s[2]) && s[3] >= 0 && (number.IsNaN(s[4]) || s[4] >= 0)
}
func (s Shape) IsRectangle() bool {
	return s[2] < 0 && s[3] >= 0 && s[4] < 0
}
func (s Shape) IsCapsule() bool {
	return s[2] >= 0 && s[3] >= 0 && s[4] < 0
}
func (s Shape) IsCone() bool {
	return s[3] < 0 && s[3] > -360 && !number.IsNaN(s[4])
}
func (s Shape) IsCircle() bool {
	return !number.IsNaN(s[2]) && s[3] <= -360 && number.IsNaN(s[4])
}

//=================================================================

func (s Shape) A() (ax, ay float32) {
	if s.IsLine() {
		return s[0], s[1]
	}
	return number.NaN(), number.NaN()
}
func (s Shape) B() (bx, by float32) {
	if s.IsLine() {
		return s[2], s[3]
	}
	return number.NaN(), number.NaN()
}
func (s Shape) Size() (width, height float32) {
	if s.IsPoint() {
		return 0, 0
	}
	if s.IsCircle() {
		var r = number.Absolute(s[2])
		return r * 2, r * 2
	}
	if s.IsCone() {
		var r = number.Absolute(s[2])
		return r, r
	}
	if s.IsLine() {
		var length = point.DistanceToPoint(s[0], s[1], s[2], s[3])
		var thickness = float32(0)
		if !number.IsNaN(s[4]) {
			thickness = s[4]
		}
		return length, thickness
	}
	return number.Absolute(s[2]), number.Absolute(s[3])
}
func (s Shape) Angle() float32 {
	if s.IsLine() {
		return angle.BetweenPoints(s[0], s[1], s[2], s[3])
	}
	if s.IsRectangle() || s.IsCapsule() {
		angleIdx, _ := rrUnpack(s[4])
		return float32(angleIdx) * 0.1
	}
	if s.IsCone() {
		return s[4]
	}
	return number.NaN()
}
func (s Shape) Roundness() float32 {
	if s.IsLine() {
		if number.IsNaN(s[4]) {
			return 0
		}
		return s[4]
	}
	if s.IsRectangle() || s.IsCapsule() {
		_, r := rrUnpack(s[4])
		return float32(r)
	}
	if s.IsCone() || s.IsCircle() {
		return s.Spread()
	}
	return number.NaN()
}
func (s Shape) Spread() float32 {
	if s.IsCircle() {
		return 360
	}
	if s.IsCone() {
		return number.Absolute(s[3])
	}
	return number.NaN()
}
func (s Shape) Center() (x, y float32) {
	if s.IsLine() {
		return (s[0] + s[2]) * 0.5, (s[1] + s[3]) * 0.5
	}
	return s[0], s[1]
}
func (s Shape) Bounds() (minX, minY, maxX, maxY float32) {
	var x, y = s.A()
	if s.IsPoint() {
		return x, y, x, y
	}
	if s.IsCircle() {
		var r = number.Absolute(s[2])
		return x - r, y - r, x + r, y + r
	}

	if s.IsLine() {
		var x2, y2 = s.B()
		minX, maxX = min(x, x2), max(x, x2)
		minY, maxY = min(y, y2), max(y, y2)
		var r = s.Roundness()
		if r > 0 { // factor in thickness/roundness if it exists
			return minX - r, minY - r, maxX + r, maxY + r
		}
		return minX, minY, maxX, maxY
	}

	var w, h = s.Size() // Rectangles, Capsules, and Cones (Rotated around center)
	var hw, hh = w * 0.5, h * 0.5

	if s.IsCapsule() { // Capsules extend further by the radius on the "height" ends
		hh += hw
	}

	var sin, cos = internal.SinCos(s.Angle())
	var ex = hw*cos + hh*sin
	var ey = hw*sin + hh*cos
	return x - ex, y - ey, x + ex, y + ey
}
func (s Shape) String() string {
	switch {
	case s.IsPoint():
		return "Point"
	case s.IsLine():
		return "Line"
	case s.IsRectangle():
		return "Rectangle"
	case s.IsCapsule():
		return "Capsule"
	case s.IsCircle():
		return "Circle"
	case s.IsCone():
		return "Cone"
	}
	return ""
}

// private ========================================================

func rrPack(angle, roundness float32) float32 {
	// rrPack encodes angle (0.1° precision) and roundness into a single negative float32.
	// Encoding: -(roundnessInt * 3600 + angleIdx + 1) to ensure it's always < 0.
	var idx = ((int(angle*10) % 3600) + 3600) % 3600
	var r = min(int(number.Absolute(roundness)), 4659)
	return -float32(r*3600 + idx + 1)
}
func rrUnpack(packed float32) (angleIdx, rInt int) {
	var abs = int(-packed) - 1
	if abs < 0 {
		return 0, 0
	}
	return abs % 3600, abs / 3600
}
