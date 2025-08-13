// Package geometry implements various computational geometry algorithms
package geometry

import (
	"math"
	"sort"
)

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// Distance calculates the Euclidean distance between two points
func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// CrossProduct calculates the cross product of vectors OA and OB
func CrossProduct(O, A, B Point) float64 {
	return (A.X-O.X)*(B.Y-O.Y) - (A.Y-O.Y)*(B.X-O.X)
}

// Orientation determines the orientation of ordered triplet (p, q, r)
// Returns 0 if collinear, 1 if clockwise, 2 if counterclockwise
func Orientation(p, q, r Point) int {
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	if math.Abs(val) < 1e-9 {
		return 0 // Collinear
	}
	if val > 0 {
		return 1 // Clockwise
	}
	return 2 // Counterclockwise
}

// ConvexHullGrahamScan finds the convex hull using Graham Scan algorithm
// Time Complexity: O(n log n), Space Complexity: O(n)
func ConvexHullGrahamScan(points []Point) []Point {
	n := len(points)
	if n < 3 {
		return points
	}

	// Find the bottom-most point (or left most in case of tie)
	l := 0
	for i := 1; i < n; i++ {
		if points[i].Y < points[l].Y {
			l = i
		} else if points[i].Y == points[l].Y && points[i].X < points[l].X {
			l = i
		}
	}

	// Swap the bottom-most point to position 0
	points[0], points[l] = points[l], points[0]

	// Sort points by polar angle with respect to first point
	pivot := points[0]
	sort.Slice(points[1:], func(i, j int) bool {
		i++ // Adjust for slice starting at index 1
		j++

		o := Orientation(pivot, points[i], points[j])
		if o == 0 {
			// If collinear, put the closer point first
			return pivot.Distance(points[i]) < pivot.Distance(points[j])
		}
		return o == 2 // Counterclockwise
	})

	// Create an empty stack and push first three points
	stack := make([]Point, 0, n)
	stack = append(stack, points[0], points[1], points[2])

	// Process remaining points
	for i := 3; i < n; i++ {
		// Keep removing top while the angle formed by points next-to-top,
		// top, and points[i] makes a clockwise turn
		for len(stack) > 1 && Orientation(stack[len(stack)-2], stack[len(stack)-1], points[i]) != 2 {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, points[i])
	}

	return stack
}

// ClosestPairOfPoints finds the closest pair of points using divide and conquer
// Time Complexity: O(n log n), Space Complexity: O(n)
func ClosestPairOfPoints(points []Point) (Point, Point, float64) {
	n := len(points)
	if n < 2 {
		return Point{}, Point{}, math.Inf(1)
	}
	if n == 2 {
		return points[0], points[1], points[0].Distance(points[1])
	}

	// Sort points by X coordinate
	sortedPoints := make([]Point, len(points))
	copy(sortedPoints, points)
	sort.Slice(sortedPoints, func(i, j int) bool {
		return sortedPoints[i].X < sortedPoints[j].X
	})

	return closestPairRec(sortedPoints)
}

func closestPairRec(points []Point) (Point, Point, float64) {
	n := len(points)

	// Base case for small arrays
	if n <= 3 {
		return bruteForceClosest(points)
	}

	// Divide
	mid := n / 2
	midPoint := points[mid]

	left := points[:mid]
	right := points[mid:]

	// Conquer
	p1, q1, dl := closestPairRec(left)
	p2, q2, dr := closestPairRec(right)

	// Find minimum of the two halves
	var minP, minQ Point
	minDist := dl
	if dr < dl {
		minP, minQ = p2, q2
		minDist = dr
	} else {
		minP, minQ = p1, q1
	}

	// Create array of points close to the line dividing the two halves
	strip := make([]Point, 0)
	for _, point := range points {
		if math.Abs(point.X-midPoint.X) < minDist {
			strip = append(strip, point)
		}
	}

	// Find the closest points in strip
	stripP, stripQ, stripDist := closestInStrip(strip, minDist)
	if stripDist < minDist {
		return stripP, stripQ, stripDist
	}

	return minP, minQ, minDist
}

func bruteForceClosest(points []Point) (Point, Point, float64) {
	n := len(points)
	if n < 2 {
		return Point{}, Point{}, math.Inf(1)
	}

	minDist := math.Inf(1)
	var p1, p2 Point

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dist := points[i].Distance(points[j])
			if dist < minDist {
				minDist = dist
				p1, p2 = points[i], points[j]
			}
		}
	}

	return p1, p2, minDist
}

func closestInStrip(strip []Point, d float64) (Point, Point, float64) {
	minDist := d
	var p1, p2 Point

	// Sort strip by Y coordinate
	sort.Slice(strip, func(i, j int) bool {
		return strip[i].Y < strip[j].Y
	})

	// Find closest points in strip
	for i := 0; i < len(strip); i++ {
		j := i + 1
		for j < len(strip) && (strip[j].Y-strip[i].Y) < minDist {
			dist := strip[i].Distance(strip[j])
			if dist < minDist {
				minDist = dist
				p1, p2 = strip[i], strip[j]
			}
			j++
		}
	}

	return p1, p2, minDist
}

// LineSegment represents a line segment
type LineSegment struct {
	P1, P2 Point
}

// DoesIntersect checks if two line segments intersect
func (ls LineSegment) DoesIntersect(other LineSegment) bool {
	p1, q1 := ls.P1, ls.P2
	p2, q2 := other.P1, other.P2

	o1 := Orientation(p1, q1, p2)
	o2 := Orientation(p1, q1, q2)
	o3 := Orientation(p2, q2, p1)
	o4 := Orientation(p2, q2, q1)

	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special cases
	// p1, q1 and p2 are collinear and p2 lies on segment p1q1
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	// p1, q1 and q2 are collinear and q2 lies on segment p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	// p2, q2 and p1 are collinear and p1 lies on segment p2q2
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}

	// p2, q2 and q1 are collinear and q1 lies on segment p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}

// onSegment checks if point q lies on segment pr
func onSegment(p, q, r Point) bool {
	return q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) &&
		q.Y <= math.Max(p.Y, r.Y) && q.Y >= math.Min(p.Y, r.Y)
}

// PointInPolygon checks if a point is inside a polygon using ray casting
// Time Complexity: O(n), Space Complexity: O(1)
func PointInPolygon(point Point, polygon []Point) bool {
	n := len(polygon)
	if n < 3 {
		return false
	}

	// Create a point at infinity
	extreme := Point{X: 10000, Y: point.Y}

	// Count intersections of the line from point to extreme with sides of polygon
	count := 0
	i := 0
	for {
		next := (i + 1) % n

		// Check if the line segment from point to extreme intersects
		// with the side from polygon[i] to polygon[next]
		seg1 := LineSegment{P1: polygon[i], P2: polygon[next]}
		seg2 := LineSegment{P1: point, P2: extreme}
		if seg1.DoesIntersect(seg2) {
			// If the point is collinear with line segment, check if it lies on segment
			if Orientation(polygon[i], point, polygon[next]) == 0 {
				return onSegment(polygon[i], point, polygon[next])
			}
			count++
		}
		i = next
		if i == 0 {
			break
		}
	}

	// Return true if count is odd
	return count%2 == 1
}

// PolygonArea calculates the area of a polygon using the shoelace formula
// Time Complexity: O(n), Space Complexity: O(1)
func PolygonArea(polygon []Point) float64 {
	n := len(polygon)
	if n < 3 {
		return 0
	}

	area := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += polygon[i].X * polygon[j].Y
		area -= polygon[j].X * polygon[i].Y
	}

	return math.Abs(area) / 2.0
}

// ConvexPolygonDiameter finds the diameter (maximum distance) of a convex polygon
// Time Complexity: O(n), Space Complexity: O(1)
func ConvexPolygonDiameter(convexPolygon []Point) (Point, Point, float64) {
	n := len(convexPolygon)
	if n < 2 {
		return Point{}, Point{}, 0
	}
	if n == 2 {
		return convexPolygon[0], convexPolygon[1], convexPolygon[0].Distance(convexPolygon[1])
	}

	maxDist := 0.0
	var p1, p2 Point

	// Use rotating calipers algorithm
	i := 0
	j := 1

	// Find the farthest pair
	for k := 0; k < n; k++ {
		// Find the farthest point from edge i to i+1
		for {
			next := (j + 1) % n
			if CrossProduct(convexPolygon[i], convexPolygon[(i+1)%n], convexPolygon[j]) <
				CrossProduct(convexPolygon[i], convexPolygon[(i+1)%n], convexPolygon[next]) {
				j = next
			} else {
				break
			}
		}

		// Update maximum distance
		dist := convexPolygon[i].Distance(convexPolygon[j])
		if dist > maxDist {
			maxDist = dist
			p1, p2 = convexPolygon[i], convexPolygon[j]
		}

		dist = convexPolygon[(i+1)%n].Distance(convexPolygon[j])
		if dist > maxDist {
			maxDist = dist
			p1, p2 = convexPolygon[(i+1)%n], convexPolygon[j]
		}

		i = (i + 1) % n
	}

	return p1, p2, maxDist
}

// IsConvex checks if a polygon is convex
// Time Complexity: O(n), Space Complexity: O(1)
func IsConvex(polygon []Point) bool {
	n := len(polygon)
	if n < 3 {
		return false
	}

	sign := 0
	for i := 0; i < n; i++ {
		o := Orientation(polygon[i], polygon[(i+1)%n], polygon[(i+2)%n])
		if o != 0 {
			if sign == 0 {
				sign = o
			} else if sign != o {
				return false
			}
		}
	}

	return true
}
