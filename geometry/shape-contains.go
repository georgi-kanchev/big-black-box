package geometry

// func (s Shape) pointContains(shape Shape) bool {
// 	var x, y = s.Position()
// 	var sx, sy = shape.Position()
// 	if shape.IsPoint() {
// 		return number.IsWithin(x, sx, 0.001) && number.IsWithin(y, sy, 0.001)
// 	}
// 	if shape.IsLine() {
// 		var bx, by = shape.Point2()
// 		var crossProduct = (y-sy)*(bx-sx) - (x-sx)*(by-sy)
// 		if crossProduct != 0 {
// 			return false
// 		}
// 		return x >= min(sx, bx) && x <= max(sx, bx) && y >= min(sy, by) && y <= max(sy, by)
// 	}
// 	if shape.IsCircle() {
// 		var radius, _ = shape.Size()
// 		var distX, distY = x - sx, y - sy
// 		return distX*distX+distY*distY <= radius*radius
// 	}
// 	if shape.IsRectangle() {
// 		var dy, dx = internal.SinCos(shape.Angle())
// 		var vx, vy = x - sx, y - sy
// 		var localX = vx*dx + vy*dy
// 		var localY = -vx*dy + vy*dx
// 		var w, h = shape.Size()
// 		return number.Absolute(localX) <= w/2 && number.Absolute(localY) <= h/2
// 	}
// 	if shape.IsCapsule() {
// 		var dy, dx = internal.SinCos(shape.Angle())
// 		var vx, vy = x - sx, y - sy
// 		var localX = vx*dx + vy*dy
// 		var localY = -vx*dy + vy*dx
// 		var w, h = shape.Size()
// 		var halfLen = number.Limit(w/2-h/2, float32(0), w/2)   // inner half-length of the rectangular body
// 		var nearestX = number.Limit(localX, -halfLen, halfLen) // nearest point on the centerline
// 		var distX, distY = localX - nearestX, localY
// 		return distX*distX+distY*distY <= (h/2)*(h/2)
// 	}
// 	return false
// }
// func (s Shape) lineSegmentContains(shape Shape) bool {
// 	var ax, ay = s.Position()
// 	var bx, by = s.Point2()
// 	var sx, sy = shape.Position()
// 	if shape.IsPoint() {
// 		var crossProduct = (sx-ax)*(by-ay) - (sy-ay)*(bx-ax)
// 		if number.Absolute(crossProduct) > 1e-9 {
// 			return false
// 		}
// 		return sx >= min(ax, bx) && sx <= max(ax, bx) && sy >= min(ay, by) && sy <= max(ay, by)
// 	}
// 	if shape.IsLine() {
// 		var ex, ey = shape.Point2()
// 		return s.lineSegmentContains(Point(sx, sy)) && s.lineSegmentContains(Point(ex, ey))
// 	}
// 	return false // other shapes return false by default
// }
// func (s Shape) circleContains(shape Shape) bool {
// 	var cx, cy = s.Position()
// 	var r, _ = s.Size()
// 	var sx, sy = shape.Position()
// 	if shape.IsPoint() {
// 		var dx, dy = sx - cx, sy - cy
// 		return dx*dx+dy*dy <= r*r
// 	}
// 	if shape.IsLine() {
// 		var ex, ey = shape.Point2()
// 		return s.circleContains(Point(sx, sy)) && s.circleContains(Point(ex, ey))
// 	}
// 	if shape.IsCircle() {
// 		var r2, _ = shape.Size()
// 		var dx, dy = sx - cx, sy - cy
// 		return number.SquareRoot(dx*dx+dy*dy)+r2 <= r
// 	}
// 	if shape.IsRectangle() {
// 		var w, h = shape.Size()
// 		var sin, cos = internal.SinCos(shape.Angle())
// 		var hw, hh = w / 2, h / 2
// 		for _, c := range [4][2]float32{{hw, hh}, {-hw, hh}, {hw, -hh}, {-hw, -hh}} {
// 			if !s.circleContains(Point(sx+c[0]*cos-c[1]*sin, sy+c[0]*sin+c[1]*cos)) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	if shape.IsCapsule() {
// 		// a capsule's farthest point from any external location is on one of the two end-cap circles
// 		var w, h = shape.Size()
// 		var capRadius = h / 2
// 		var halfLen = number.Limit(w/2-h/2, float32(0), w/2)
// 		var cdy, cdx = internal.SinCos(shape.Angle())
// 		var e1x, e1y = sx + halfLen*cdx, sy + halfLen*cdy
// 		var e2x, e2y = sx - halfLen*cdx, sy - halfLen*cdy
// 		var d1x, d1y = e1x - cx, e1y - cy
// 		var d2x, d2y = e2x - cx, e2y - cy
// 		return number.SquareRoot(d1x*d1x+d1y*d1y)+capRadius <= r &&
// 			number.SquareRoot(d2x*d2x+d2y*d2y)+capRadius <= r
// 	}
// 	return false
// }
// func (s Shape) rectangleContains(shape Shape) bool {
// 	var rdy, rdx = internal.SinCos(s.Angle())
// 	var rx, ry = s.Position()
// 	var w, h = s.Size()
// 	var sx, sy = shape.Position()
// 	var vx, vy = sx - rx, sy - ry
// 	var lx = vx*rdx + vy*rdy // shape center in rectangle local frame
// 	var ly = -vx*rdy + vy*rdx
// 	if shape.IsPoint() {
// 		return number.Absolute(lx) <= w/2 && number.Absolute(ly) <= h/2
// 	}
// 	if shape.IsLine() {
// 		var ex, ey = shape.Point2()
// 		return s.rectangleContains(Point(sx, sy)) && s.rectangleContains(Point(ex, ey))
// 	}
// 	if shape.IsCircle() {
// 		var r2, _ = shape.Size()
// 		return number.Absolute(lx)+r2 <= w/2 && number.Absolute(ly)+r2 <= h/2
// 	}
// 	if shape.IsRectangle() {
// 		var w2, h2 = shape.Size()
// 		var edy, edx = internal.SinCos(shape.Angle())
// 		var cp = edx*rdx + edy*rdy  // cos of angle between rectangles
// 		var sp = -edx*rdy + edy*rdx // sin of angle between rectangles
// 		return number.Absolute(lx)+w2/2*number.Absolute(cp)+h2/2*number.Absolute(sp) <= w/2 &&
// 			number.Absolute(ly)+w2/2*number.Absolute(sp)+h2/2*number.Absolute(cp) <= h/2
// 	}
// 	if shape.IsCapsule() {
// 		var w2, h2 = shape.Size()
// 		var capRadius = h2 / 2
// 		var halfLen = number.Limit(w2/2-h2/2, float32(0), w2/2)
// 		var cdy, cdx = internal.SinCos(shape.Angle())
// 		var capDirX = cdx*rdx + cdy*rdy // capsule axis direction in rectangle local frame
// 		var capDirY = -cdx*rdy + cdy*rdx
// 		var c1x, c1y = lx + halfLen*capDirX, ly + halfLen*capDirY // cap centers in local frame
// 		var c2x, c2y = lx - halfLen*capDirX, ly - halfLen*capDirY
// 		return number.Absolute(c1x)+capRadius <= w/2 && number.Absolute(c1y)+capRadius <= h/2 &&
// 			number.Absolute(c2x)+capRadius <= w/2 && number.Absolute(c2y)+capRadius <= h/2
// 	}
// 	return false
// }
// func (s Shape) capsuleContains(shape Shape) bool {
// 	var cdy, cdx = internal.SinCos(s.Angle())
// 	var ox, oy = s.Position()
// 	var w, h = s.Size()
// 	var capRadius = h / 2
// 	var halfLen = number.Limit(w/2-h/2, float32(0), w/2)
// 	var sx, sy = shape.Position()
// 	var vx, vy = sx - ox, sy - oy
// 	var lx = vx*cdx + vy*cdy // shape center in capsule local frame
// 	var ly = -vx*cdy + vy*cdx
// 	if shape.IsPoint() {
// 		var nearestX = number.Limit(lx, -halfLen, halfLen)
// 		var dx, dy = lx - nearestX, ly
// 		return dx*dx+dy*dy <= capRadius*capRadius
// 	}
// 	if shape.IsLine() {
// 		var p2x, p2y = shape.Point2()
// 		return s.capsuleContains(Point(sx, sy)) && s.capsuleContains(Point(p2x, p2y))
// 	}
// 	if shape.IsCircle() {
// 		var r, _ = shape.Size()
// 		var nearestX = number.Limit(lx, -halfLen, halfLen)
// 		var dx, dy = lx - nearestX, ly
// 		return number.SquareRoot(dx*dx+dy*dy)+r <= capRadius
// 	}
// 	if shape.IsRectangle() { // check all 4 corners (exact by convexity)
// 		var w2, h2 = shape.Size()
// 		var rdy2, rdx2 = internal.SinCos(shape.Angle())
// 		var cp = rdx2*cdx + rdy2*cdy
// 		var sp = -rdx2*cdy + rdy2*cdx
// 		var hw2, hh2 = w2 / 2, h2 / 2
// 		for _, c := range [4][2]float32{{hw2, hh2}, {-hw2, hh2}, {hw2, -hh2}, {-hw2, -hh2}} {
// 			var cx2 = lx + c[0]*cp - c[1]*sp
// 			var cy2 = ly + c[0]*sp + c[1]*cp
// 			var nearestX = number.Limit(cx2, -halfLen, halfLen)
// 			var dx, dy = cx2 - nearestX, cy2
// 			if dx*dx+dy*dy > capRadius*capRadius {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	if shape.IsCapsule() { // inner capsule ⊆ outer capsule iff both inner cap circles fit inside the outer capsule
// 		var w2, h2 = shape.Size()
// 		var innerCapRadius = h2 / 2
// 		var innerHalfLen = number.Limit(w2/2-h2/2, float32(0), w2/2)
// 		var c2dy, c2dx = internal.SinCos(shape.Angle())
// 		var capDirX = c2dx*cdx + c2dy*cdy // inner capsule axis in outer local frame
// 		var capDirY = -c2dx*cdy + c2dy*cdx
// 		var c1x, c1y = lx + innerHalfLen*capDirX, ly + innerHalfLen*capDirY
// 		var c2x, c2y = lx - innerHalfLen*capDirX, ly - innerHalfLen*capDirY
// 		for _, cap := range [2][2]float32{{c1x, c1y}, {c2x, c2y}} {
// 			var nearestX = number.Limit(cap[0], -halfLen, halfLen)
// 			var dx, dy = cap[0] - nearestX, cap[1]
// 			if number.SquareRoot(dx*dx+dy*dy)+innerCapRadius > capRadius {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	return false
// }
