// Package sorting implements various sorting algorithms
package sorting

import (
	"math/rand"
	"time"
)

// BubbleSort implements the bubble sort algorithm
// Time Complexity: O(n²), Space Complexity: O(1)
func BubbleSort(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	n := len(result)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
				swapped = true
			}
		}
		// If no swapping occurred, array is already sorted
		if !swapped {
			break
		}
	}
	return result
}

// SelectionSort implements the selection sort algorithm
// Time Complexity: O(n²), Space Complexity: O(1)
func SelectionSort(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	n := len(result)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if result[j] < result[minIdx] {
				minIdx = j
			}
		}
		result[i], result[minIdx] = result[minIdx], result[i]
	}
	return result
}

// InsertionSort implements the insertion sort algorithm
// Time Complexity: O(n²), Space Complexity: O(1)
func InsertionSort(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	for i := 1; i < len(result); i++ {
		key := result[i]
		j := i - 1
		
		// Move elements greater than key one position ahead
		for j >= 0 && result[j] > key {
			result[j+1] = result[j]
			j--
		}
		result[j+1] = key
	}
	return result
}

// MergeSort implements the merge sort algorithm
// Time Complexity: O(n log n), Space Complexity: O(n)
func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	
	result := make([]int, len(arr))
	copy(result, arr)
	
	mergeSortHelper(result, 0, len(result)-1)
	return result
}

func mergeSortHelper(arr []int, left, right int) {
	if left < right {
		mid := left + (right-left)/2
		
		mergeSortHelper(arr, left, mid)
		mergeSortHelper(arr, mid+1, right)
		merge(arr, left, mid, right)
	}
}

func merge(arr []int, left, mid, right int) {
	// Create temporary arrays for left and right subarrays
	leftArr := make([]int, mid-left+1)
	rightArr := make([]int, right-mid)
	
	copy(leftArr, arr[left:mid+1])
	copy(rightArr, arr[mid+1:right+1])
	
	i, j, k := 0, 0, left
	
	// Merge the temporary arrays back into arr[left..right]
	for i < len(leftArr) && j < len(rightArr) {
		if leftArr[i] <= rightArr[j] {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}
	
	// Copy remaining elements
	for i < len(leftArr) {
		arr[k] = leftArr[i]
		i++
		k++
	}
	
	for j < len(rightArr) {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

// QuickSort implements the quick sort algorithm
// Time Complexity: O(n log n) average, O(n²) worst, Space Complexity: O(log n)
func QuickSort(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	quickSortHelper(result, 0, len(result)-1)
	return result
}

func quickSortHelper(arr []int, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		
		quickSortHelper(arr, low, pi-1)
		quickSortHelper(arr, pi+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// QuickSortRandomized implements quick sort with random pivot selection
func QuickSortRandomized(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	rand.Seed(time.Now().UnixNano())
	quickSortRandomizedHelper(result, 0, len(result)-1)
	return result
}

func quickSortRandomizedHelper(arr []int, low, high int) {
	if low < high {
		// Randomize pivot
		randomIndex := low + rand.Intn(high-low+1)
		arr[randomIndex], arr[high] = arr[high], arr[randomIndex]
		
		pi := partition(arr, low, high)
		
		quickSortRandomizedHelper(arr, low, pi-1)
		quickSortRandomizedHelper(arr, pi+1, high)
	}
}

// HeapSort implements the heap sort algorithm
// Time Complexity: O(n log n), Space Complexity: O(1)
func HeapSort(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	
	n := len(result)
	
	// Build max heap
	for i := n/2 - 1; i >= 0; i-- {
		heapify(result, n, i)
	}
	
	// Extract elements from heap one by one
	for i := n - 1; i > 0; i-- {
		result[0], result[i] = result[i], result[0]
		heapify(result, i, 0)
	}
	
	return result
}

func heapify(arr []int, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2
	
	if left < n && arr[left] > arr[largest] {
		largest = left
	}
	
	if right < n && arr[right] > arr[largest] {
		largest = right
	}
	
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

// CountingSort implements counting sort for non-negative integers
// Time Complexity: O(n + k), Space Complexity: O(k)
func CountingSort(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	
	// Find the maximum element
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
		if v < 0 {
			// Counting sort works only for non-negative integers
			return arr
		}
	}
	
	// Create counting array
	count := make([]int, max+1)
	
	// Count occurrences
	for _, v := range arr {
		count[v]++
	}
	
	// Build result array
	result := make([]int, 0, len(arr))
	for i, c := range count {
		for j := 0; j < c; j++ {
			result = append(result, i)
		}
	}
	
	return result
}

// RadixSort implements radix sort for non-negative integers
// Time Complexity: O(d × (n + k)), Space Complexity: O(n + k)
func RadixSort(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	
	result := make([]int, len(arr))
	copy(result, arr)
	
	// Find maximum number to know number of digits
	max := result[0]
	for _, v := range result {
		if v > max {
			max = v
		}
		if v < 0 {
			// Radix sort works only for non-negative integers
			return result
		}
	}
	
	// Do counting sort for every digit
	for exp := 1; max/exp > 0; exp *= 10 {
		countingSortByDigit(result, exp)
	}
	
	return result
}

func countingSortByDigit(arr []int, exp int) {
	n := len(arr)
	output := make([]int, n)
	count := make([]int, 10)
	
	// Count occurrences of each digit
	for i := 0; i < n; i++ {
		count[(arr[i]/exp)%10]++
	}
	
	// Change count[i] to actual position
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}
	
	// Build output array
	for i := n - 1; i >= 0; i-- {
		output[count[(arr[i]/exp)%10]-1] = arr[i]
		count[(arr[i]/exp)%10]--
	}
	
	// Copy output array to arr
	copy(arr, output)
}

// IsSorted checks if an array is sorted in ascending order
func IsSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

// Reverse reverses an array in place
func Reverse(arr []int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
