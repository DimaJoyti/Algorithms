package geometry

import (
	"math"
	"testing"
)

func TestPointDistance(t *testing.T) {
	p1 := Point{X: 0, Y: 0}
	p2 := Point{X: 3, Y: 4}

	expected := 5.0
	actual := p1.Distance(p2)

	if math.Abs(actual-expected) > 1e-9 {
		t.Errorf("Distance(%v, %v) = %f, expected %f", p1, p2, actual, expected)
	}
}

func TestOrientation(t *testing.T) {
	tests := []struct {
		p, q, r  Point
		expected int
	}{
		{Point{0, 0}, Point{4, 4}, Point{1, 2}, 2}, // Counterclockwise
		{Point{0, 0}, Point{4, 4}, Point{1, 1}, 0}, // Collinear
		{Point{0, 0}, Point{4, 4}, Point{2, 1}, 1}, // Clockwise
	}

	for _, test := range tests {
		actual := Orientation(test.p, test.q, test.r)
		if actual != test.expected {
			t.Errorf("Orientation(%v, %v, %v) = %d, expected %d",
				test.p, test.q, test.r, actual, test.expected)
		}
	}
}

func TestConvexHullGrahamScan(t *testing.T) {
	// Test case 1: Simple square
	points := []Point{
		{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0.5, 0.5}, // Interior point
	}

	hull := ConvexHullGrahamScan(points)

	// Should have 4 points for a square
	if len(hull) != 4 {
		t.Errorf("ConvexHull of square should have 4 points, got %d", len(hull))
	}

	// Test case 2: Collinear points
	collinear := []Point{
		{0, 0}, {1, 1}, {2, 2},
	}

	hullCollinear := ConvexHullGrahamScan(collinear)
	// Graham scan may include all collinear points on the boundary
	if len(hullCollinear) < 2 {
		t.Errorf("ConvexHull of collinear points should have at least 2 points, got %d", len(hullCollinear))
	}

	// Test case 3: Less than 3 points
	twoPoints := []Point{{0, 0}, {1, 1}}
	hullTwo := ConvexHullGrahamScan(twoPoints)
	if len(hullTwo) != 2 {
		t.Errorf("ConvexHull of 2 points should return 2 points, got %d", len(hullTwo))
	}
}

func TestClosestPairOfPoints(t *testing.T) {
	// Test case 1: Simple case
	points := []Point{
		{2, 3}, {12, 30}, {40, 50}, {5, 1}, {12, 10}, {3, 4},
	}

	p1, p2, dist := ClosestPairOfPoints(points)

	// The closest pair should be (2,3) and (3,4) with distance sqrt(2)
	expectedDist := math.Sqrt(2)
	if math.Abs(dist-expectedDist) > 1e-9 {
		t.Errorf("ClosestPair distance = %f, expected %f", dist, expectedDist)
	}

	// Verify the points are correct
	actualDist := p1.Distance(p2)
	if math.Abs(actualDist-dist) > 1e-9 {
		t.Errorf("Returned points distance %f doesn't match returned distance %f", actualDist, dist)
	}

	// Test case 2: Two points
	twoPoints := []Point{{0, 0}, {1, 1}}
	_, _, distTwo := ClosestPairOfPoints(twoPoints)
	expectedDistTwo := math.Sqrt(2)
	if math.Abs(distTwo-expectedDistTwo) > 1e-9 {
		t.Errorf("ClosestPair of two points = %f, expected %f", distTwo, expectedDistTwo)
	}

	// Test case 3: One point
	onePoint := []Point{{0, 0}}
	_, _, distOne := ClosestPairOfPoints(onePoint)
	if !math.IsInf(distOne, 1) {
		t.Errorf("ClosestPair of one point should return infinity, got %f", distOne)
	}
}

func TestLineSegmentIntersection(t *testing.T) {
	tests := []struct {
		seg1, seg2 LineSegment
		expected   bool
	}{
		// Intersecting segments
		{
			LineSegment{Point{1, 1}, Point{10, 1}},
			LineSegment{Point{1, 2}, Point{10, 2}},
			false,
		},
		// Intersecting segments
		{
			LineSegment{Point{10, 0}, Point{0, 10}},
			LineSegment{Point{0, 0}, Point{10, 10}},
			true,
		},
		// Non-intersecting segments
		{
			LineSegment{Point{-5, -5}, Point{0, 0}},
			LineSegment{Point{1, 1}, Point{10, 10}},
			false,
		},
	}

	for i, test := range tests {
		actual := test.seg1.DoesIntersect(test.seg2)
		if actual != test.expected {
			t.Errorf("Test %d: DoesIntersect = %v, expected %v", i, actual, test.expected)
		}
	}
}

func TestPointInPolygon(t *testing.T) {
	// Square polygon
	square := []Point{
		{0, 0}, {4, 0}, {4, 4}, {0, 4},
	}

	tests := []struct {
		point    Point
		expected bool
	}{
		{Point{2, 2}, true},   // Inside
		{Point{0, 0}, true},   // On vertex
		{Point{2, 0}, true},   // On edge
		{Point{5, 5}, false},  // Outside
		{Point{-1, 2}, false}, // Outside
	}

	for i, test := range tests {
		actual := PointInPolygon(test.point, square)
		if actual != test.expected {
			t.Errorf("Test %d: PointInPolygon(%v) = %v, expected %v",
				i, test.point, actual, test.expected)
		}
	}
}

func TestPolygonArea(t *testing.T) {
	// Test case 1: Square with side length 4
	square := []Point{
		{0, 0}, {4, 0}, {4, 4}, {0, 4},
	}

	expectedArea := 16.0
	actualArea := PolygonArea(square)

	if math.Abs(actualArea-expectedArea) > 1e-9 {
		t.Errorf("PolygonArea of square = %f, expected %f", actualArea, expectedArea)
	}

	// Test case 2: Triangle
	triangle := []Point{
		{0, 0}, {4, 0}, {2, 3},
	}

	expectedTriangleArea := 6.0
	actualTriangleArea := PolygonArea(triangle)

	if math.Abs(actualTriangleArea-expectedTriangleArea) > 1e-9 {
		t.Errorf("PolygonArea of triangle = %f, expected %f", actualTriangleArea, expectedTriangleArea)
	}

	// Test case 3: Less than 3 points
	line := []Point{{0, 0}, {1, 1}}
	areaLine := PolygonArea(line)
	if areaLine != 0 {
		t.Errorf("PolygonArea of line should be 0, got %f", areaLine)
	}
}

func TestIsConvex(t *testing.T) {
	// Test case 1: Convex square
	square := []Point{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
	}

	if !IsConvex(square) {
		t.Error("Square should be convex")
	}

	// Test case 2: Non-convex polygon (star shape)
	star := []Point{
		{0, 0}, {2, 1}, {1, 2}, {2, 3}, {0, 2}, {-1, 3}, {-1, 1}, {-2, 1},
	}

	if IsConvex(star) {
		t.Error("Star shape should not be convex")
	}

	// Test case 3: Triangle (always convex)
	triangle := []Point{
		{0, 0}, {1, 0}, {0.5, 1},
	}

	if !IsConvex(triangle) {
		t.Error("Triangle should be convex")
	}

	// Test case 4: Less than 3 points
	line := []Point{{0, 0}, {1, 1}}
	if IsConvex(line) {
		t.Error("Line should not be convex")
	}
}

func TestConvexPolygonDiameter(t *testing.T) {
	// Test case 1: Square
	square := []Point{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
	}

	_, _, diameter := ConvexPolygonDiameter(square)
	expectedDiameter := math.Sqrt(2) // Diagonal of unit square

	if math.Abs(diameter-expectedDiameter) > 1e-9 {
		t.Errorf("ConvexPolygonDiameter of square = %f, expected %f", diameter, expectedDiameter)
	}

	// Test case 2: Two points
	twoPoints := []Point{{0, 0}, {3, 4}}
	_, _, diameterTwo := ConvexPolygonDiameter(twoPoints)
	expectedDiameterTwo := 5.0

	if math.Abs(diameterTwo-expectedDiameterTwo) > 1e-9 {
		t.Errorf("ConvexPolygonDiameter of two points = %f, expected %f", diameterTwo, expectedDiameterTwo)
	}

	// Test case 3: One point
	onePoint := []Point{{0, 0}}
	_, _, diameterOne := ConvexPolygonDiameter(onePoint)
	if diameterOne != 0 {
		t.Errorf("ConvexPolygonDiameter of one point should be 0, got %f", diameterOne)
	}
}

// Benchmark tests
func BenchmarkConvexHullGrahamScan(b *testing.B) {
	points := make([]Point, 1000)
	for i := 0; i < 1000; i++ {
		points[i] = Point{X: float64(i), Y: float64(i * i % 100)}
	}

	for i := 0; i < b.N; i++ {
		ConvexHullGrahamScan(points)
	}
}

func BenchmarkClosestPairOfPoints(b *testing.B) {
	points := make([]Point, 1000)
	for i := 0; i < 1000; i++ {
		points[i] = Point{X: float64(i), Y: float64(i * 7 % 100)}
	}

	for i := 0; i < b.N; i++ {
		ClosestPairOfPoints(points)
	}
}

func BenchmarkPointInPolygon(b *testing.B) {
	// Create a polygon with many vertices
	polygon := make([]Point, 100)
	for i := 0; i < 100; i++ {
		angle := 2 * math.Pi * float64(i) / 100
		polygon[i] = Point{X: math.Cos(angle), Y: math.Sin(angle)}
	}

	testPoint := Point{X: 0.5, Y: 0.5}

	for i := 0; i < b.N; i++ {
		PointInPolygon(testPoint, polygon)
	}
}

func BenchmarkPolygonArea(b *testing.B) {
	polygon := make([]Point, 1000)
	for i := 0; i < 1000; i++ {
		angle := 2 * math.Pi * float64(i) / 1000
		polygon[i] = Point{X: math.Cos(angle), Y: math.Sin(angle)}
	}

	for i := 0; i < b.N; i++ {
		PolygonArea(polygon)
	}
}
