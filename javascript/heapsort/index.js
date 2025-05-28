// --- Directions
// Implement the Heap Sort algorithm.
// Heap Sort is a comparison-based sorting algorithm that uses a binary heap data structure.
// It divides its input into a sorted and an unsorted region, and iteratively shrinks the 
// unsorted region by extracting the largest element and moving that to the sorted region.
// --- Examples
//   heapSort([3, 1, 4, 1, 5, 9, 2, 6]) --> [1, 1, 2, 3, 4, 5, 6, 9]
//   heapSort([5, 4, 3, 2, 1]) --> [1, 2, 3, 4, 5]
//   heapSort([1]) --> [1]
//   heapSort([]) --> []
// --- Complexity
//   Time: O(n log n) in all cases
//   Space: O(1) - in-place sorting

/**
 * Heap Sort implementation
 * @param {number[]} arr - Array to be sorted
 * @returns {number[]} Sorted array
 */
function heapSort(arr) {
  if (arr.length <= 1) {
    return arr;
  }

  const n = arr.length;

  // Build max heap (rearrange array)
  for (let i = Math.floor(n / 2) - 1; i >= 0; i--) {
    heapify(arr, n, i);
  }

  // Extract elements from heap one by one
  for (let i = n - 1; i > 0; i--) {
    // Move current root to end
    swap(arr, 0, i);

    // Call max heapify on the reduced heap
    heapify(arr, i, 0);
  }

  return arr;
}

/**
 * Heapify a subtree rooted with node i which is an index in arr[]
 * @param {number[]} arr - Array to heapify
 * @param {number} n - Size of heap
 * @param {number} i - Root index
 */
function heapify(arr, n, i) {
  let largest = i; // Initialize largest as root
  const left = 2 * i + 1; // Left child
  const right = 2 * i + 2; // Right child

  // If left child is larger than root
  if (left < n && arr[left] > arr[largest]) {
    largest = left;
  }

  // If right child is larger than largest so far
  if (right < n && arr[right] > arr[largest]) {
    largest = right;
  }

  // If largest is not root
  if (largest !== i) {
    swap(arr, i, largest);

    // Recursively heapify the affected sub-tree
    heapify(arr, n, largest);
  }
}

/**
 * Swap two elements in an array
 * @param {number[]} arr - Array containing elements to swap
 * @param {number} i - Index of first element
 * @param {number} j - Index of second element
 */
function swap(arr, i, j) {
  const temp = arr[i];
  arr[i] = arr[j];
  arr[j] = temp;
}

/**
 * Build a max heap from an array
 * @param {number[]} arr - Array to build heap from
 * @returns {number[]} Max heap array
 */
function buildMaxHeap(arr) {
  const n = arr.length;
  const heap = [...arr];

  // Start from the last non-leaf node and heapify each node
  for (let i = Math.floor(n / 2) - 1; i >= 0; i--) {
    heapify(heap, n, i);
  }

  return heap;
}

/**
 * Extract maximum element from heap
 * @param {number[]} heap - Max heap array
 * @returns {number|null} Maximum element or null if heap is empty
 */
function extractMax(heap) {
  if (heap.length === 0) {
    return null;
  }

  if (heap.length === 1) {
    return heap.pop();
  }

  // Store the maximum value
  const max = heap[0];

  // Move the last element to root
  heap[0] = heap.pop();

  // Heapify the root
  heapify(heap, heap.length, 0);

  return max;
}

/**
 * Insert a new element into the heap
 * @param {number[]} heap - Max heap array
 * @param {number} key - Element to insert
 */
function insertIntoHeap(heap, key) {
  // Add the new key at the end
  heap.push(key);

  // Fix the max heap property if it's violated
  let i = heap.length - 1;
  while (i > 0 && heap[Math.floor((i - 1) / 2)] < heap[i]) {
    swap(heap, i, Math.floor((i - 1) / 2));
    i = Math.floor((i - 1) / 2);
  }
}

/**
 * Check if an array represents a valid max heap
 * @param {number[]} arr - Array to check
 * @returns {boolean} True if array is a valid max heap
 */
function isMaxHeap(arr) {
  const n = arr.length;

  for (let i = 0; i <= Math.floor((n - 2) / 2); i++) {
    const left = 2 * i + 1;
    const right = 2 * i + 2;

    // Check if parent is greater than left child
    if (left < n && arr[i] < arr[left]) {
      return false;
    }

    // Check if parent is greater than right child
    if (right < n && arr[i] < arr[right]) {
      return false;
    }
  }

  return true;
}

module.exports = {
  heapSort,
  heapify,
  swap,
  buildMaxHeap,
  extractMax,
  insertIntoHeap,
  isMaxHeap
};
