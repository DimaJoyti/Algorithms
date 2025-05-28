package searching

import (
	"testing"
)

func TestLinearSearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Single element not found", []int{5}, 3, -1},
		{"Element at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"Element at end", []int{1, 2, 3, 4, 5}, 5, 4},
		{"Element in middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"Element not found", []int{1, 2, 3, 4, 5}, 6, -1},
		{"Duplicates", []int{1, 2, 2, 3, 4}, 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LinearSearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("LinearSearch(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Single element not found", []int{5}, 3, -1},
		{"Element at beginning", []int{1, 2, 3, 4, 5}, 1, 0},
		{"Element at end", []int{1, 2, 3, 4, 5}, 5, 4},
		{"Element in middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"Element not found", []int{1, 2, 3, 4, 5}, 6, -1},
		{"Large array", []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}, 13, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BinarySearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("BinarySearch(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestBinarySearchRecursive(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Element in middle", []int{1, 2, 3, 4, 5}, 3, 2},
		{"Element not found", []int{1, 2, 3, 4, 5}, 6, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BinarySearchRecursive(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("BinarySearchRecursive(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestSearchRotatedArray(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"No rotation", []int{1, 2, 3, 4, 5}, 3, 2},
		{"Rotated once", []int{5, 1, 2, 3, 4}, 3, 3},
		{"Rotated multiple", []int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{"Target at pivot", []int{4, 5, 6, 7, 0, 1, 2}, 4, 0},
		{"Target not found", []int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{"Single element", []int{1}, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchRotatedArray(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("SearchRotatedArray(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestFindPivot(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected int
	}{
		{"No rotation", []int{1, 2, 3, 4, 5}, 0},
		{"Rotated once", []int{5, 1, 2, 3, 4}, 1},
		{"Rotated multiple", []int{4, 5, 6, 7, 0, 1, 2}, 4},
		{"Single element", []int{1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindPivot(tt.arr)
			if result != tt.expected {
				t.Errorf("FindPivot(%v) = %d, want %d", tt.arr, result, tt.expected)
			}
		})
	}
}

func TestInterpolationSearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Uniform distribution", []int{10, 20, 30, 40, 50}, 30, 2},
		{"Element not found", []int{10, 20, 30, 40, 50}, 25, -1},
		{"Target out of range", []int{10, 20, 30, 40, 50}, 60, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InterpolationSearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("InterpolationSearch(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestExponentialSearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Element at beginning", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 0},
		{"Element in middle", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5, 4},
		{"Element at end", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, 9},
		{"Element not found", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 11, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExponentialSearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("ExponentialSearch(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestJumpSearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Element found", []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610}, 55, 10},
		{"Element not found", []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610}, 4, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JumpSearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("JumpSearch(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestTernarySearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Element in middle", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 5, 4},
		{"Element not found", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TernarySearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("TernarySearch(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestFindFirst(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"No duplicates", []int{1, 2, 3, 4, 5}, 3, 2},
		{"With duplicates", []int{1, 2, 2, 2, 3, 4, 5}, 2, 1},
		{"All same", []int{2, 2, 2, 2, 2}, 2, 0},
		{"Not found", []int{1, 2, 3, 4, 5}, 6, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindFirst(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("FindFirst(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestFindLast(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"No duplicates", []int{1, 2, 3, 4, 5}, 3, 2},
		{"With duplicates", []int{1, 2, 2, 2, 3, 4, 5}, 2, 3},
		{"All same", []int{2, 2, 2, 2, 2}, 2, 4},
		{"Not found", []int{1, 2, 3, 4, 5}, 6, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindLast(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("FindLast(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestCountOccurrences(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"No duplicates", []int{1, 2, 3, 4, 5}, 3, 1},
		{"With duplicates", []int{1, 2, 2, 2, 3, 4, 5}, 2, 3},
		{"All same", []int{2, 2, 2, 2, 2}, 2, 5},
		{"Not found", []int{1, 2, 3, 4, 5}, 6, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountOccurrences(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("CountOccurrences(%v, %d) = %d, want %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestSearchMatrix(t *testing.T) {
	matrix := [][]int{
		{1, 4, 7, 11},
		{2, 5, 8, 12},
		{3, 6, 9, 16},
		{10, 13, 14, 17},
	}

	tests := []struct {
		name      string
		matrix    [][]int
		target    int
		expectRow int
		expectCol int
		expectErr bool
	}{
		{"Element found", matrix, 5, 1, 1, false},
		{"Element at corner", matrix, 1, 0, 0, false},
		{"Element not found", matrix, 15, -1, -1, true},
		{"Empty matrix", [][]int{}, 5, -1, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, err := SearchMatrix(tt.matrix, tt.target)
			if (err != nil) != tt.expectErr {
				t.Errorf("SearchMatrix() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if row != tt.expectRow || col != tt.expectCol {
				t.Errorf("SearchMatrix() = (%d, %d), want (%d, %d)", row, col, tt.expectRow, tt.expectCol)
			}
		})
	}
}

// Benchmark tests
func BenchmarkLinearSearch(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = i
	}
	target := 500

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearSearch(arr, target)
	}
}

func BenchmarkBinarySearch(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = i
	}
	target := 500

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BinarySearch(arr, target)
	}
}

func BenchmarkInterpolationSearch(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = i * 10 // Uniform distribution
	}
	target := 5000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InterpolationSearch(arr, target)
	}
}
