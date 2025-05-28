package sorting

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestBubbleSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Duplicates", []int{3, 1, 3, 1, 3}, []int{1, 1, 3, 3, 3}},
		{"Negative numbers", []int{-3, 1, -4, 1, 5}, []int{-4, -3, 1, 1, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BubbleSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BubbleSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSelectionSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectionSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SelectionSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInsertionSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertionSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertionSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMergeSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestQuickSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuickSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("QuickSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestQuickSortRandomized(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuickSortRandomized(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("QuickSortRandomized(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHeapSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HeapSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("HeapSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountingSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Small range", []int{4, 2, 2, 8, 3, 3, 1}, []int{1, 2, 2, 3, 3, 4, 8}},
		{"With zeros", []int{0, 5, 3, 0, 2}, []int{0, 0, 2, 3, 5}},
		{"Negative numbers", []int{-1, 2, 3}, []int{-1, 2, 3}}, // Should return unchanged
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountingSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CountingSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRadixSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Multi-digit", []int{170, 45, 75, 90, 2, 802, 24, 66}, []int{2, 24, 45, 66, 75, 90, 170, 802}},
		{"Single digits", []int{4, 2, 2, 8, 3, 3, 1}, []int{1, 2, 2, 3, 3, 4, 8}},
		{"Negative numbers", []int{-1, 2, 3}, []int{-1, 2, 3}}, // Should return unchanged
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RadixSort(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RadixSort(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{"Empty array", []int{}, true},
		{"Single element", []int{5}, true},
		{"Sorted ascending", []int{1, 2, 3, 4, 5}, true},
		{"Not sorted", []int{5, 4, 3, 2, 1}, false},
		{"Partially sorted", []int{1, 2, 4, 3, 5}, false},
		{"Duplicates sorted", []int{1, 1, 2, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSorted(tt.input)
			if result != tt.expected {
				t.Errorf("IsSorted(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty array", []int{}, []int{}},
		{"Single element", []int{5}, []int{5}},
		{"Multiple elements", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"Even length", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]int, len(tt.input))
			copy(input, tt.input)
			Reverse(input)
			if !reflect.DeepEqual(input, tt.expected) {
				t.Errorf("Reverse(%v) = %v, want %v", tt.input, input, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkBubbleSort(b *testing.B) {
	arr := generateRandomArray(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BubbleSort(arr)
	}
}

func BenchmarkQuickSort(b *testing.B) {
	arr := generateRandomArray(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(arr)
	}
}

func BenchmarkMergeSort(b *testing.B) {
	arr := generateRandomArray(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeSort(arr)
	}
}

func BenchmarkHeapSort(b *testing.B) {
	arr := generateRandomArray(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HeapSort(arr)
	}
}

func generateRandomArray(size int) []int {
	rand.Seed(time.Now().UnixNano())
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(1000)
	}
	return arr
}
