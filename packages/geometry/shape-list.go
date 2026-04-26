package geometry

import (
	"slices"
)

const MaxShapes = 1024

type ShapeStore struct {
	Data [MaxShapes]Shape
	Len  int // Tracks how many active shapes are packed at the start
}

// Add inserts a shape and maintains the X-axis sort order
func (s *ShapeStore) Add(shape Shape) {
	if s.Len >= MaxShapes {
		return // Or handle overflow
	}
	// Insert at the end, then sort.
	// For better perf with many adds, sort once after a batch.
	s.Data[s.Len] = shape
	s.Len++
	s.sort()
}

// Remove finds a shape (by value or logic) and maintains packing
func (s *ShapeStore) Remove(index int) {
	if index < 0 || index >= s.Len {
		return
	}
	// Move the last element to the removed spot to keep the array contiguous
	s.Data[index] = s.Data[s.Len-1]
	s.Data[s.Len-1] = Shape{} // Clear for GC (though not strictly needed for primitives)
	s.Len--
	s.sort()
}

func (s *ShapeStore) sort() {
	// Sort by MinX (X - Width/2) for the Sweep algorithm
	active := s.Data[:s.Len]
	slices.SortFunc(active, func(a, b Shape) int {
		minA := a.X - (a.Width / 2)
		minB := b.X - (b.Width / 2)
		if minA < minB {
			return -1
		}
		if minA > minB {
			return 1
		}
		return 0
	})
}

// GetNeighbors returns active shapes that could potentially overlap the target
func (s *ShapeStore) GetNeighbors(target Shape) []Shape {
	if s.Len == 0 {
		return nil
	}

	targetMinX, targetMinY, targetMaxX, targetMaxY := target.Bounds()

	// 1. Find the starting point using Binary Search on the sorted X axis
	// We want the first shape whose MaxX >= targetMinX
	startIdx := 0
	for i := 0; i < s.Len; i++ {
		_, _, maxX, _ := s.Data[i].Bounds()
		if maxX >= targetMinX {
			startIdx = i
			break
		}
	}

	// 2. Sweep forward
	// Create a sub-slice (view) of the array to return
	// Note: In a real tight loop, you might want to provide a buffer
	// to avoid slice header overhead, but this is the idiomatic way.
	var results []Shape
	for i := startIdx; i < s.Len; i++ {
		other := s.Data[i]
		otherMinX, otherMinY, _, otherMaxY := other.Bounds()

		// If this shape's MinX is beyond our target's MaxX,
		// NO subsequent shapes can possibly collide. STOP.
		if otherMinX > targetMaxX {
			break
		}

		// Quick Y-axis pruning before doing heavy "Overlaps" math
		if otherMaxY < targetMinY || otherMinY > targetMaxY {
			continue
		}

		results = append(results, other)
	}

	return results
}
