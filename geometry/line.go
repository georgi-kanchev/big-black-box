package geometry

import "big-black-box/utility/number"

func LinesCrossPoint(line1, line2 Shape) (x, y float32) {
	var dx1, dy1 = line1[2] - line1[0], line1[3] - line1[1]
	var dx2, dy2 = line2[2] - line2[0], line2[3] - line2[1]
	var det = dx1*dy2 - dy1*dx2

	if det > -0.001 && det < 0.001 {
		return number.NaN(), number.NaN()
	}

	var s = ((line1[1]-line2[1])*dx2 - (line1[0]-line2[0])*dy2) / det
	var t = ((line1[1]-line2[1])*dx1 - (line1[0]-line2[0])*dy1) / det

	if s < 0 || s > 1 || t < 0 || t > 1 {
		return number.NaN(), number.NaN()
	}

	return line1[0] + s*dx1, line1[1] + s*dy1
}
func LinesCross(line1, line2 Shape) bool {
	var d1 = (line2[2]-line2[0])*(line1[1]-line2[1]) - (line2[3]-line2[1])*(line1[0]-line2[0])
	var d2 = (line2[2]-line2[0])*(line1[3]-line2[1]) - (line2[3]-line2[1])*(line1[2]-line2[0])
	var d3 = (line1[2]-line1[0])*(line2[1]-line1[1]) - (line1[3]-line1[1])*(line2[0]-line1[0])
	var d4 = (line1[2]-line1[0])*(line2[3]-line1[1]) - (line1[3]-line1[1])*(line2[2]-line1[0])
	return d1*d2 < 0 && d3*d4 < 0
}

func LineClosestPoint(line Shape, pointX, pointY float32) (x, y float32) {
	var apx, apy = pointX - line[0], pointY - line[1]
	var abx, aby = line[2] - line[0], line[3] - line[1]
	var magSq = abx*abx + aby*aby

	if magSq == 0 {
		return line[0], line[1]
	}

	var dot = apx*abx + apy*aby
	var dist = dot / magSq

	if dist < 0 {
		return line[0], line[1]
	}
	if dist > 1 {
		return line[2], line[3]
	}

	return line[0] + abx*dist, line[1] + aby*dist
}
func LineIsLeftOfPoint(line Shape, pointX, pointY float32) bool {
	return (line[2]-line[0])*(pointY-line[1])-(line[3]-line[1])*(pointX-line[0]) < 0
}
func LineContainsPoint(line Shape, pointX, pointY, tolerance float32) bool {
	var cx, cy = LineClosestPoint(line, pointX, pointY)
	var dx, dy = cx - pointX, cy - pointY
	return dx*dx+dy*dy <= tolerance*tolerance
}
