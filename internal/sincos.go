package internal

import "math"

// optimized for speed via power-of-two lookup table and bitwise masking
func SinCos(degrees float32) (sin, cos float32) {
	// Scale degrees to the 4096 range (roughly 11.377 units per degree)
	var index = int(degrees * (4096.0 / 360.0))
	index &= 4095 // Fast bitwise wrap-around (replaces modulo)

	// Sine is direct lookup, cosine is sine shifted by 90 degrees (1024 indices)
	return sineTable[index], sineTable[(index+1024)&4095]
}

func SinCosCache() {
	for i := range 4096 {
		// Convert index back to radians for the initial calculation
		// (i / 4096.0) * 2 * Pi
		var rad = (float64(i) / 4096.0) * (2.0 * math.Pi)
		sineTable[i] = float32(math.Sin(rad))
	}
}

// private ========================================================
// 4096 * 4 bytes = 16KB (Fits comfortably in L1 Cache)
var sineTable [4096]float32
