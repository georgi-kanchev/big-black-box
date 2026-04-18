package internal

import "math"

// optimized for speed via lookup table
func SinCos(degrees float32) (sin, cos float32) {
	var index = int(degrees * 10)          // convert to index (0.1 degree precision)
	index = ((index % 3600) + 3600) % 3600 // and wrap 0-3599
	// sine is direct lookup, cosine is sine shifted by 90 degrees (900 indices)
	return sineTable[index], sineTable[(index+900)%3600]
}

func SinCosCache() {
	for i := range 3600 {
		var rad = float64(i) * math.Pi / 1800.0 // convert index to radians (i / 10.0 * Pi / 180.0)
		sineTable[i] = float32(math.Sin(rad))
	}
}

// private =================================================================
var sineTable [3600]float32
