// A few helper functions, related to randomness. Also provides a way to use controlled randomness
// in the form of seeds, as well as combining seeds. Has a few functions that act upon a collection -
// shuffling it or choosing an item.
package random

import (
	"big-black-box/packages/utility/number"
	"math/rand/v2"
	"time"
)

func AnySeed() float32 {
	seedCounter++
	var n = time.Now().UnixNano()
	var h = hashSeed(uint64(n), seedCounter)
	return float32(h) / 18446744073709551615.0
}
func CombineSeeds[T number.Number](seed1, seed2 T) T {
	var c1, c2 uint64
	switch any(seed1).(type) {
	case int, int8, int16, int32, int64:
		c1 = uint64(int64(seed1))
		c2 = uint64(int64(seed2))
	case uint, uint8, uint16, uint32, uint64:
		c1 = uint64(seed1)
		c2 = uint64(seed2)
	case float32, float64:
		c1 = uint64(float64(seed1) * 1e9)
		c2 = uint64(float64(seed2) * 1e9)
	}

	var out = hashSeed(hashSeed(uint64(2654435769), c1), c2)
	var zero T
	switch any(zero).(type) {
	case int:
		return T(int(out))
	case int8:
		return T(int8(out))
	case int16:
		return T(int16(out))
	case int32:
		return T(int32(out))
	case int64:
		return T(int64(out))
	case uint:
		return T(uint(out))
	case uint8:
		return T(uint8(out))
	case uint16:
		return T(uint16(out))
	case uint32:
		return T(uint32(out))
	case uint64:
		return T(uint64(out))
	case float32:
		return T(float32(out))
	case float64:
		return T(float64(out))
	}
	return zero
}
func Range[T number.Number](min, max T, seed float32) T {
	switch any(min).(type) {
	case int, int8, int16, int32, int64:
		return T(rangeInt(int64(min), int64(max), seed))
	case uint, uint8, uint16, uint32, uint64:
		return T(rangeUint(uint64(min), uint64(max), seed))
	case float32, float64:
		return T(rangeFloat(float64(min), float64(max), seed))
	}
	var zero T
	return zero
}
func HasChance(percent, seed float32) bool {
	if percent <= 0 {
		return false
	}
	return Range(float32(0), 100, seed) <= min(100, percent)
}
func Shuffle[T any](items []T, seed float32) []T {
	for i := len(items) - 1; i > 0; i-- {
		var j = int(Range(0, i, seed))
		items[i], items[j] = items[j], items[i]
	}
	return items
}
func PickFrom[T any](items []T, seed float32) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	return items[int(Range(0, len(items)-1, seed))]
}

// private ========================================================

var seedCounter uint64

func hashSeed(seed, value uint64) uint64 {
	seed ^= value
	seed = (seed ^ (seed >> 16)) * 2246822519
	seed = (seed ^ (seed >> 13)) * 3266489917
	seed ^= seed >> 16
	return seed
}
func rangeInt(val1, val2 int64, seed float32) int64 {
	var ua, ub uint64
	ua, ub = uint64(val1), uint64(val2)

	if ua == ub {
		return val1
	}
	if ua > ub {
		ua, ub = ub, ua
	}

	var diff = ub - ua
	if number.IsNaN(seed) {
		seed = rand.Float32()
	}
	var s = uint64(seed * 2147483647)
	s = (1103515245*s + 12345) % 2147483647
	var result = ua + (s*diff)/2147483647
	return int64(result)
}
func rangeUint(ua, ub uint64, seed float32) uint64 {
	if ua == ub {
		return ua
	}
	if ua > ub {
		ua, ub = ub, ua
	}

	var diff = ub - ua
	if number.IsNaN(seed) {
		seed = rand.Float32()
	}
	var s = uint64(seed * 2147483647)
	s = (1103515245*s + 12345) % 2147483647
	var result = ua + (s*diff)/2147483647
	return result
}
func rangeFloat(fa, fb float64, seed float32) float64 {
	if fa == fb {
		return fa
	}
	if fa > fb {
		fa, fb = fb, fa
	}

	if number.IsNaN(seed) {
		seed = rand.Float32()
	}

	var s = int(seed * 2147483647)
	s = (1103515245*s + 12345) % 2147483647
	var normalized = float64(s) / 2147483647.0
	var r = fa + (fb-fa)*normalized
	return r
}
