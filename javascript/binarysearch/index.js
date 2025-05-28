// --- Directions
// Implement Binary Search algorithm.
// Binary Search is a search algorithm that finds the position of a target value 
// within a sorted array. It compares the target value to the middle element of 
// the array and eliminates half of the search space in each iteration.
// --- Examples
//   binarySearch([1, 2, 3, 4, 5], 3) --> 2
//   binarySearch([1, 2, 3, 4, 5], 6) --> -1
//   binarySearch([1, 3, 5, 7, 9], 1) --> 0
//   binarySearch([1, 3, 5, 7, 9], 9) --> 4
// --- Complexity
//   Time: O(log n)
//   Space: O(1) for iterative, O(log n) for recursive

/**
 * Binary Search - Iterative implementation
 * @param {number[]} arr - Sorted array to search in
 * @param {number} target - Value to search for
 * @returns {number} Index of target if found, -1 otherwise
 */
function binarySearch(arr, target) {
  let left = 0;
  let right = arr.length - 1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] === target) {
      return mid;
    } else if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return -1; // Target not found
}

/**
 * Binary Search - Recursive implementation
 * @param {number[]} arr - Sorted array to search in
 * @param {number} target - Value to search for
 * @param {number} left - Left boundary (default: 0)
 * @param {number} right - Right boundary (default: arr.length - 1)
 * @returns {number} Index of target if found, -1 otherwise
 */
function binarySearchRecursive(arr, target, left = 0, right = arr.length - 1) {
  if (left > right) {
    return -1; // Base case: target not found
  }

  const mid = Math.floor((left + right) / 2);

  if (arr[mid] === target) {
    return mid;
  } else if (arr[mid] < target) {
    return binarySearchRecursive(arr, target, mid + 1, right);
  } else {
    return binarySearchRecursive(arr, target, left, mid - 1);
  }
}

/**
 * Find the first occurrence of target in a sorted array with duplicates
 * @param {number[]} arr - Sorted array to search in
 * @param {number} target - Value to search for
 * @returns {number} Index of first occurrence, -1 if not found
 */
function binarySearchFirst(arr, target) {
  let left = 0;
  let right = arr.length - 1;
  let result = -1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] === target) {
      result = mid;
      right = mid - 1; // Continue searching in left half
    } else if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return result;
}

/**
 * Find the last occurrence of target in a sorted array with duplicates
 * @param {number[]} arr - Sorted array to search in
 * @param {number} target - Value to search for
 * @returns {number} Index of last occurrence, -1 if not found
 */
function binarySearchLast(arr, target) {
  let left = 0;
  let right = arr.length - 1;
  let result = -1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] === target) {
      result = mid;
      left = mid + 1; // Continue searching in right half
    } else if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return result;
}

/**
 * Find the insertion point for target to maintain sorted order
 * @param {number[]} arr - Sorted array
 * @param {number} target - Value to find insertion point for
 * @returns {number} Index where target should be inserted
 */
function binarySearchInsertionPoint(arr, target) {
  let left = 0;
  let right = arr.length;

  while (left < right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] < target) {
      left = mid + 1;
    } else {
      right = mid;
    }
  }

  return left;
}

/**
 * Search in a rotated sorted array
 * @param {number[]} arr - Rotated sorted array
 * @param {number} target - Value to search for
 * @returns {number} Index of target if found, -1 otherwise
 */
function searchRotatedArray(arr, target) {
  let left = 0;
  let right = arr.length - 1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] === target) {
      return mid;
    }

    // Check which half is sorted
    if (arr[left] <= arr[mid]) {
      // Left half is sorted
      if (target >= arr[left] && target < arr[mid]) {
        right = mid - 1;
      } else {
        left = mid + 1;
      }
    } else {
      // Right half is sorted
      if (target > arr[mid] && target <= arr[right]) {
        left = mid + 1;
      } else {
        right = mid - 1;
      }
    }
  }

  return -1;
}

/**
 * Find peak element in array (element greater than its neighbors)
 * @param {number[]} arr - Array to search in
 * @returns {number} Index of a peak element
 */
function findPeakElement(arr) {
  let left = 0;
  let right = arr.length - 1;

  while (left < right) {
    const mid = Math.floor((left + right) / 2);

    if (arr[mid] > arr[mid + 1]) {
      right = mid;
    } else {
      left = mid + 1;
    }
  }

  return left;
}

/**
 * Search in a 2D matrix where each row and column is sorted
 * @param {number[][]} matrix - 2D sorted matrix
 * @param {number} target - Value to search for
 * @returns {boolean} True if target is found, false otherwise
 */
function searchMatrix(matrix, target) {
  if (!matrix || matrix.length === 0 || matrix[0].length === 0) {
    return false;
  }

  let row = 0;
  let col = matrix[0].length - 1;

  while (row < matrix.length && col >= 0) {
    if (matrix[row][col] === target) {
      return true;
    } else if (matrix[row][col] > target) {
      col--;
    } else {
      row++;
    }
  }

  return false;
}

module.exports = {
  binarySearch,
  binarySearchRecursive,
  binarySearchFirst,
  binarySearchLast,
  binarySearchInsertionPoint,
  searchRotatedArray,
  findPeakElement,
  searchMatrix
};
