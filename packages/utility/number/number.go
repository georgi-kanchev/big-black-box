// Everything can be represented by numbers, making this one of the most important packages.
// It has many useful helper functions, operating on all sorts of numbers, representing all sorts of things.
// Contains functionalities such as:
//   - Formatting numbers into text.
//   - Numbers interacting with ranges.
//   - Shortcuts for finding a certain number in a collection.
//   - Transforming, wrapping, mapping numbers.
//   - Wrapping universally known math functions to make them more digestible and clarify their API.
//   - Number checks.
//   - Index conversion.
//   - etc
package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Number interface{ Float | Integer }
type Float interface{ ~float32 | ~float64 }
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Format[T Number](number T) string {
	switch v := any(number).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprint(v) // fallback
	}
}
func SeparateThousands[T Number](number T) string {
	var str = Format(number)
	var parts = strings.SplitN(str, ".", 2)
	var intPart = parts[0]
	var result = ""
	var n = len(intPart)
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result += " "
		}
		result += string(c)
	}

	if len(parts) == 2 {
		result += "."
		result += parts[1]
	}
	return result
}

func Limit[T Number](number, a, b T) T {
	if a > b {
		a, b = b, a
	}
	return max(a, min(number, b))
}
func Map[T Number](number, fromA, fromB, toA, toB T) T {
	// Convert inputs to float64 for high-precision math
	var n, fa, fb, ta, tb = float64(number), float64(fromA), float64(fromB), float64(toA), float64(toB)
	var deltaFrom = fb - fa

	if math.Abs(deltaFrom) < 1e-9 { // is zero
		return T((ta + tb) / 2)
	}

	var value = ((n-fa)/deltaFrom)*(tb-ta) + ta
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return toA
	}

	return T(value)
}
func Wrap[T Number](number, a, b T) T {
	if a > b {
		a, b = b, a
	}

	switch any(number).(type) {
	case float32, float64:
		var d = float64(b - a)
		if d < 0.001 {
			return a
		}
		var num = float64(number - a)
		var wrapped = math.Mod(math.Mod(num, d)+d, d) + float64(a)
		return T(wrapped)
	default: // integer types
		var d = int64(b - a)
		if d == 0 {
			return a
		}
		var num = int64(number - a)
		var wrapped = ((num % d) + d) % d
		return T(wrapped + int64(a))
	}
}

func Absolute[T Number](number T) T {
	if number < 0 {
		return -number
	}
	return number
}
func Unsign[T Number](number T) T {
	return Absolute(number)
}
func Snap[T Float](number, interval T) T {
	var n, i = float64(number), float64(interval)

	if i == 0 || math.IsNaN(i) || math.IsInf(n, 0) {
		return number
	}

	var remainder = math.Mod(n, i)
	var halfway = i / 2.0

	if remainder < 0 {
		if remainder > -halfway {
			return T(n - remainder)
		}
		return T(n - (remainder + i))
	}

	if remainder < halfway {
		return T(n - remainder)
	}
	return T(n + (i - remainder))
}
func Power[T Number](number, power T) T {
	return T(math.Pow(float64(number), float64(power)))
}
func SquareRoot[T Number](number T) T {
	return T(math.Sqrt(float64(number)))
}

func IsWithin[T Number](number, target, distance T) bool {
	return target-distance <= number && number <= target+distance
}

func ValueSmallest[T Number]() T {
	var zero T
	switch any(&zero).(type) {
	case *int:
		return any(math.MinInt).(T)
	case *int8:
		return any(int8(math.MinInt8)).(T)
	case *int16:
		return any(int16(math.MinInt16)).(T)
	case *int32:
		return any(int32(math.MinInt32)).(T)
	case *int64:
		return any(int64(math.MinInt64)).(T)
	case *uint, *uint8, *uint16, *uint32, *uint64:
		return zero // zero is already 0 for these types
	case *float32:
		return any(float32(-math.MaxFloat32)).(T)
	case *float64:
		return any(-math.MaxFloat64).(T)
	default:
		return zero
	}
}
func ValueBiggest[T Number]() T {
	var zero T
	switch any(&zero).(type) {
	case *int:
		return any(math.MaxInt).(T)
	case *int8:
		return any(int8(math.MaxInt8)).(T)
	case *int16:
		return any(int16(math.MaxInt16)).(T)
	case *int32:
		return any(int32(math.MaxInt32)).(T)
	case *int64:
		return any(int64(math.MaxInt64)).(T)
	case *uint:
		return any(uint(math.MaxUint)).(T)
	case *uint8:
		return any(uint8(math.MaxUint8)).(T)
	case *uint16:
		return any(uint16(math.MaxUint16)).(T)
	case *uint32:
		return any(uint32(math.MaxUint32)).(T)
	case *uint64:
		return any(uint64(math.MaxUint64)).(T)
	case *float32:
		return any(float32(math.MaxFloat32)).(T)
	case *float64:
		return any(math.MaxFloat64).(T)
	default:
		return zero
	}
}

//=================================================================

func Animate[T Float](value, target, rate T) T {
	var result T
	var factor float64 = 1.0 - math.Pow(2.0, -float64(rate))
	var delta float64 = float64(target - value)

	result = T(float64(value) + delta*factor)

	if IsWithin(float64(result), float64(target), float64(0.001)) {
		return target
	}

	return result
}

func DivisionRemainder[T Float](number, target T) T {
	return T(math.Mod(float64(number), float64(target)))
}
func Sine[T Float](number T) T {
	// var sin, _ = internal.SinCos(float32(number))
	// return T(sin)
	return T(math.Sin(float64(number)))
}
func Cosine[T Float](number T) T {
	// var _, cos = internal.SinCos(float32(number))
	// return T(cos)
	return T(math.Cos(float64(number)))
}
func Precision[T Float](number T) int {
	for i := range 9 {
		if math.Abs(float64(number)-math.Round(float64(number))) < 1e-6 {
			return i
		}
		number *= 10
	}
	return 0
}
func Exponential[T Float](number T) T {
	return T(math.Exp(float64(number)))
}

func Round[T Float](number T) T {
	return T(math.Round(float64(number)))
}
func RoundUp[T Float](number T) T {
	return T(math.Ceil(float64(number)))
}
func RoundDown[T Float](number T) T {
	return T(math.Floor(float64(number)))
}
func RoundFraction[T Float](number T, precision int) T {
	var p = getPow(precision)
	return T(math.Round(float64(number)*p) / p)
}
func RoundUpFraction[T Float](number T, precision int) T {
	var p = getPow(precision)
	return T(math.Ceil(float64(number)*p) / p)
}
func RoundDownFraction[T Float](number T, precision int) T {
	var p = getPow(precision)
	return T(math.Floor(float64(number)*p) / p)
}

func PositiveInfinity() float32 {
	return float32(math.Inf(1))
}
func NegativeInfinity() float32 {
	return float32(math.Inf(-1))
}
func NaN() float32 {
	return float32(math.NaN())
}

func IsNaN(number float32) bool {
	return number != number
}
func IsInfinity(number float32) bool {
	return number > math.MaxFloat32 || number < -math.MaxFloat32
}
func IsPositiveInfinity(number float32) bool {
	return number > math.MaxFloat32
}
func IsNegativeInfinity(number float32) bool {
	return number < -math.MaxFloat32
}

//=================================================================

func Indexes2DToIndex1D[T Integer](x, y, width, height T) T {
	var result = x*width + y
	var max = width * height
	if result < 0 {
		return 0
	} else if result > max {
		return max
	}
	return result
}
func Index1DToIndexes2D[T Integer](index, width, height T) (x, y T) {
	var max = width * height
	if index < 0 {
		index = 0
	} else if index > max {
		index = max
	}
	x = index % width
	y = index / width
	return x, y
}

// private ========================================================

var pow10 = [7]float64{1, 10, 100, 1000, 10000, 100000, 1000000}

func getPow(precision int) float64 {
	if precision >= 0 && precision < len(pow10) {
		return pow10[precision]
	}
	return math.Pow(10, float64(precision))
}
