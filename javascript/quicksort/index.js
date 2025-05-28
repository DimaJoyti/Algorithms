// --- Directions
// Implement the Quick Sort algorithm.
// Quick Sort is a divide-and-conquer algorithm that works by selecting a 'pivot' 
// element from the array and partitioning the other elements into two sub-arrays, 
// according to whether they are less than or greater than the pivot.
// --- Examples
//   quickSort([3, 1, 4, 1, 5, 9, 2, 6]) --> [1, 1, 2, 3, 4, 5, 6, 9]
//   quickSort([5, 4, 3, 2, 1]) --> [1, 2, 3, 4, 5]
//   quickSort([1]) --> [1]
//   quickSort([]) --> []
// --- Complexity
//   Time: O(n log n) average case, O(n²) worst case
//   Space: O(log n) average case due to recursion stack

/**
 * Quick Sort implementation using Lomuto partition scheme
 * @param {number[]} arr - Array to be sorted
 * @param {number} low - Starting index (default: 0)
 * @param {number} high - Ending index (default: arr.length - 1)
 * @returns {number[]} Sorted array
 */
function quickSort(arr, low = 0, high = arr.length - 1) {
  // Base case: if array has 0 or 1 element, it's already sorted
  if (low >= high) {
    return arr;
  }
  
  // Partition the array and get the pivot index
  const pivotIndex = partition(arr, low, high);
  
  // Recursively sort elements before and after partition
  quickSort(arr, low, pivotIndex - 1);
  quickSort(arr, pivotIndex + 1, high);
  
  return arr;
}

/**
 * Partition function using Lomuto partition scheme
 * Places the pivot element at its correct position in sorted array
 * and places all smaller elements to left of pivot and all greater to right
 * @param {number[]} arr - Array to partition
 * @param {number} low - Starting index
 * @param {number} high - Ending index
 * @returns {number} Index of the pivot element after partitioning
 */
function partition(arr, low, high) {
  // Choose the rightmost element as pivot
  const pivot = arr[high];
  
  // Index of smaller element indicates the right position of pivot found so far
  let i = low - 1;
  
  for (let j = low; j < high; j++) {
    // If current element is smaller than or equal to pivot
    if (arr[j] <= pivot) {
      i++; // increment index of smaller element
      swap(arr, i, j);
    }
  }
  
  // Place pivot at correct position
  swap(arr, i + 1, high);
  return i + 1;
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
 * Quick Sort with random pivot for better average case performance
 * @param {number[]} arr - Array to be sorted
 * @param {number} low - Starting index (default: 0)
 * @param {number} high - Ending index (default: arr.length - 1)
 * @returns {number[]} Sorted array
 */
function quickSortRandomized(arr, low = 0, high = arr.length - 1) {
  if (low >= high) {
    return arr;
  }
  
  // Randomize pivot to improve average case performance
  const randomIndex = Math.floor(Math.random() * (high - low + 1)) + low;
  swap(arr, randomIndex, high);
  
  const pivotIndex = partition(arr, low, high);
  
  quickSortRandomized(arr, low, pivotIndex - 1);
  quickSortRandomized(arr, pivotIndex + 1, high);
  
  return arr;
}

/**
 * Quick Sort that returns a new array (non-mutating version)
 * @param {number[]} arr - Array to be sorted
 * @returns {number[]} New sorted array
 */
function quickSortPure(arr) {
  // Base case
  if (arr.length <= 1) {
    return arr;
  }
  
  // Choose pivot (middle element for better average case)
  const pivotIndex = Math.floor(arr.length / 2);
  const pivot = arr[pivotIndex];
  
  // Partition array into three parts
  const less = [];
  const equal = [];
  const greater = [];
  
  for (let i = 0; i < arr.length; i++) {
    if (arr[i] < pivot) {
      less.push(arr[i]);
    } else if (arr[i] > pivot) {
      greater.push(arr[i]);
    } else {
      equal.push(arr[i]);
    }
  }
  
  // Recursively sort and combine
  return [
    ...quickSortPure(less),
    ...equal,
    ...quickSortPure(greater)
  ];
}

module.exports = {
  quickSort,
  quickSortRandomized,
  quickSortPure,
  partition,
  swap
};
