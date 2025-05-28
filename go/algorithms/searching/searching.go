// Package searching implements various search algorithms
package searching

import "errors"

// LinearSearch performs linear search on an array
// Time Complexity: O(n), Space Complexity: O(1)
func LinearSearch(arr []int, target int) int {
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

// BinarySearch performs binary search on a sorted array
// Time Complexity: O(log n), Space Complexity: O(1)
func BinarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	
	for left <= right {
		mid := left + (right-left)/2
		
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	
	return -1
}

// BinarySearchRecursive performs recursive binary search
// Time Complexity: O(log n), Space Complexity: O(log n)
func BinarySearchRecursive(arr []int, target int) int {
	return binarySearchHelper(arr, target, 0, len(arr)-1)
}

func binarySearchHelper(arr []int, target, left, right int) int {
	if left > right {
		return -1
	}
	
	mid := left + (right-left)/2
	
	if arr[mid] == target {
		return mid
	} else if arr[mid] < target {
		return binarySearchHelper(arr, target, mid+1, right)
	} else {
		return binarySearchHelper(arr, target, left, mid-1)
	}
}

// SearchRotatedArray searches in a rotated sorted array
// Time Complexity: O(log n), Space Complexity: O(1)
func SearchRotatedArray(arr []int, target int) int {
	left, right := 0, len(arr)-1
	
	for left <= right {
		mid := left + (right-left)/2
		
		if arr[mid] == target {
			return mid
		}
		
		// Check if left half is sorted
		if arr[left] <= arr[mid] {
			if target >= arr[left] && target < arr[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// Right half is sorted
			if target > arr[mid] && target <= arr[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	
	return -1
}

// FindPivot finds the pivot point in a rotated sorted array
func FindPivot(arr []int) int {
	left, right := 0, len(arr)-1
	
	for left < right {
		mid := left + (right-left)/2
		
		if arr[mid] > arr[right] {
			left = mid + 1
		} else {
			right = mid
		}
	}
	
	return left
}

// InterpolationSearch performs interpolation search on a uniformly distributed sorted array
// Time Complexity: O(log log n) average, O(n) worst, Space Complexity: O(1)
func InterpolationSearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	
	for left <= right && target >= arr[left] && target <= arr[right] {
		if left == right {
			if arr[left] == target {
				return left
			}
			return -1
		}
		
		// Calculate position using interpolation formula
		pos := left + ((target-arr[left])*(right-left))/(arr[right]-arr[left])
		
		if arr[pos] == target {
			return pos
		} else if arr[pos] < target {
			left = pos + 1
		} else {
			right = pos - 1
		}
	}
	
	return -1
}

// ExponentialSearch performs exponential search
// Time Complexity: O(log n), Space Complexity: O(1)
func ExponentialSearch(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}
	
	if arr[0] == target {
		return 0
	}
	
	// Find range for binary search
	bound := 1
	for bound < len(arr) && arr[bound] < target {
		bound *= 2
	}
	
	// Perform binary search in the found range
	left := bound / 2
	right := bound
	if right >= len(arr) {
		right = len(arr) - 1
	}
	
	return binarySearchRange(arr, target, left, right)
}

func binarySearchRange(arr []int, target, left, right int) int {
	for left <= right {
		mid := left + (right-left)/2
		
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	
	return -1
}

// JumpSearch performs jump search
// Time Complexity: O(√n), Space Complexity: O(1)
func JumpSearch(arr []int, target int) int {
	n := len(arr)
	if n == 0 {
		return -1
	}
	
	// Calculate optimal jump size
	step := int(sqrt(float64(n)))
	prev := 0
	
	// Jump through the array
	for arr[min(step, n)-1] < target {
		prev = step
		step += int(sqrt(float64(n)))
		if prev >= n {
			return -1
		}
	}
	
	// Linear search in the identified block
	for arr[prev] < target {
		prev++
		if prev == min(step, n) {
			return -1
		}
	}
	
	if arr[prev] == target {
		return prev
	}
	
	return -1
}

// TernarySearch performs ternary search on a sorted array
// Time Complexity: O(log₃ n), Space Complexity: O(1)
func TernarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	
	for left <= right {
		mid1 := left + (right-left)/3
		mid2 := right - (right-left)/3
		
		if arr[mid1] == target {
			return mid1
		}
		if arr[mid2] == target {
			return mid2
		}
		
		if target < arr[mid1] {
			right = mid1 - 1
		} else if target > arr[mid2] {
			left = mid2 + 1
		} else {
			left = mid1 + 1
			right = mid2 - 1
		}
	}
	
	return -1
}

// FindFirst finds the first occurrence of target in a sorted array with duplicates
func FindFirst(arr []int, target int) int {
	left, right := 0, len(arr)-1
	result := -1
	
	for left <= right {
		mid := left + (right-left)/2
		
		if arr[mid] == target {
			result = mid
			right = mid - 1 // Continue searching in left half
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	
	return result
}

// FindLast finds the last occurrence of target in a sorted array with duplicates
func FindLast(arr []int, target int) int {
	left, right := 0, len(arr)-1
	result := -1
	
	for left <= right {
		mid := left + (right-left)/2
		
		if arr[mid] == target {
			result = mid
			left = mid + 1 // Continue searching in right half
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	
	return result
}

// CountOccurrences counts occurrences of target in a sorted array
func CountOccurrences(arr []int, target int) int {
	first := FindFirst(arr, target)
	if first == -1 {
		return 0
	}
	
	last := FindLast(arr, target)
	return last - first + 1
}

// SearchMatrix searches for target in a 2D matrix where each row and column is sorted
// Time Complexity: O(m + n), Space Complexity: O(1)
func SearchMatrix(matrix [][]int, target int) (int, int, error) {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return -1, -1, errors.New("empty matrix")
	}
	
	rows, cols := len(matrix), len(matrix[0])
	row, col := 0, cols-1
	
	for row < rows && col >= 0 {
		if matrix[row][col] == target {
			return row, col, nil
		} else if matrix[row][col] > target {
			col--
		} else {
			row++
		}
	}
	
	return -1, -1, errors.New("target not found")
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}
	
	// Newton's method for square root
	guess := x
	for i := 0; i < 10; i++ {
		guess = (guess + x/guess) / 2
	}
	return guess
}
