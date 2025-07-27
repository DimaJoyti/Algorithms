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

// BucketSort sorts an array using bucket sort algorithm
// Time Complexity: O(n + k), Space Complexity: O(n + k) where k is the range
// Works best when input is uniformly distributed
func BucketSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	result := make([]int, len(arr))
	copy(result, arr)

	// Find min and max values
	minVal, maxVal := result[0], result[0]
	for _, val := range result {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	// Handle edge case where all elements are the same
	if minVal == maxVal {
		return result
	}

	// Create buckets
	bucketCount := len(result)
	buckets := make([][]int, bucketCount)
	for i := range buckets {
		buckets[i] = make([]int, 0)
	}

	// Distribute elements into buckets
	for _, val := range result {
		bucketIndex := (val - minVal) * (bucketCount - 1) / (maxVal - minVal)
		buckets[bucketIndex] = append(buckets[bucketIndex], val)
	}

	// Sort individual buckets and concatenate
	index := 0
	for _, bucket := range buckets {
		if len(bucket) > 0 {
			// Use insertion sort for small buckets
			sortedBucket := InsertionSort(bucket)
			for _, val := range sortedBucket {
				result[index] = val
				index++
			}
		}
	}

	return result
}

// PancakeSort sorts an array using pancake sort algorithm
// Time Complexity: O(n²), Space Complexity: O(1)
// Interesting algorithm that only uses flip operations
func PancakeSort(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)

	n := len(result)
	for currSize := n; currSize > 1; currSize-- {
		// Find index of maximum element in result[0..currSize-1]
		maxIdx := findMaxIndex(result, currSize)

		// If maximum element is not at the end, move it to the end
		if maxIdx != currSize-1 {
			// First move maximum element to beginning if it's not already there
			if maxIdx != 0 {
				flip(result, maxIdx)
			}
			// Now move the maximum element from beginning to end
			flip(result, currSize-1)
		}
	}

	return result
}

// findMaxIndex finds the index of maximum element in arr[0..n-1]
func findMaxIndex(arr []int, n int) int {
	maxIdx := 0
	for i := 1; i < n; i++ {
		if arr[i] > arr[maxIdx] {
			maxIdx = i
		}
	}
	return maxIdx
}

// flip reverses arr[0..i]
func flip(arr []int, i int) {
	start := 0
	for start < i {
		arr[start], arr[i] = arr[i], arr[start]
		start++
		i--
	}
}

// TimSort implements a simplified version of TimSort algorithm
// Time Complexity: O(n log n), Space Complexity: O(n)
// This is a hybrid stable sorting algorithm derived from merge sort and insertion sort
func TimSort(arr []int) []int {
	const MIN_MERGE = 32

	result := make([]int, len(arr))
	copy(result, arr)
	n := len(result)

	if n < 2 {
		return result
	}

	// Sort individual subarrays of size MIN_MERGE using insertion sort
	for i := 0; i < n; i += MIN_MERGE {
		end := i + MIN_MERGE - 1
		if end > n-1 {
			end = n - 1
		}
		insertionSortRange(result, i, end)
	}

	// Start merging from size MIN_MERGE
	size := MIN_MERGE
	for size < n {
		// Pick starting point of left sub array
		for start := 0; start < n; start += size * 2 {
			// Calculate mid and end points
			mid := start + size - 1
			end := start + size*2 - 1

			// Merge subarrays if mid is smaller than end
			if mid < n-1 {
				if end > n-1 {
					end = n - 1
				}
				mergeRange(result, start, mid, end)
			}
		}
		size *= 2
	}

	return result
}

// insertionSortRange performs insertion sort on arr[left..right]
func insertionSortRange(arr []int, left, right int) {
	for i := left + 1; i <= right; i++ {
		key := arr[i]
		j := i - 1
		for j >= left && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// mergeRange merges two sorted subarrays arr[left..mid] and arr[mid+1..right]
func mergeRange(arr []int, left, mid, right int) {
	// Create temporary arrays for left and right subarrays
	leftArr := make([]int, mid-left+1)
	rightArr := make([]int, right-mid)

	// Copy data to temporary arrays
	for i := 0; i < len(leftArr); i++ {
		leftArr[i] = arr[left+i]
	}
	for i := 0; i < len(rightArr); i++ {
		rightArr[i] = arr[mid+1+i]
	}

	// Merge the temporary arrays back into arr[left..right]
	i, j, k := 0, 0, left
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
