const {
  binarySearch,
  binarySearchRecursive,
  binarySearchFirst,
  binarySearchLast,
  binarySearchInsertionPoint,
  searchRotatedArray,
  findPeakElement,
  searchMatrix
} = require('./index');

describe('Binary Search Algorithms', () => {
  describe('binarySearch function', () => {
    test('binarySearch function exists', () => {
      expect(typeof binarySearch).toEqual('function');
    });

    test('finds element in middle of array', () => {
      expect(binarySearch([1, 2, 3, 4, 5], 3)).toBe(2);
    });

    test('finds element at beginning of array', () => {
      expect(binarySearch([1, 2, 3, 4, 5], 1)).toBe(0);
    });

    test('finds element at end of array', () => {
      expect(binarySearch([1, 2, 3, 4, 5], 5)).toBe(4);
    });

    test('returns -1 for element not in array', () => {
      expect(binarySearch([1, 2, 3, 4, 5], 6)).toBe(-1);
    });

    test('handles empty array', () => {
      expect(binarySearch([], 1)).toBe(-1);
    });

    test('handles single element array - found', () => {
      expect(binarySearch([42], 42)).toBe(0);
    });

    test('handles single element array - not found', () => {
      expect(binarySearch([42], 1)).toBe(-1);
    });

    test('works with negative numbers', () => {
      expect(binarySearch([-5, -3, -1, 0, 2, 4], -1)).toBe(2);
    });
  });

  describe('binarySearchRecursive function', () => {
    test('binarySearchRecursive function exists', () => {
      expect(typeof binarySearchRecursive).toEqual('function');
    });

    test('finds element recursively', () => {
      expect(binarySearchRecursive([1, 2, 3, 4, 5], 3)).toBe(2);
    });

    test('returns -1 for element not found recursively', () => {
      expect(binarySearchRecursive([1, 2, 3, 4, 5], 6)).toBe(-1);
    });

    test('handles empty array recursively', () => {
      expect(binarySearchRecursive([], 1)).toBe(-1);
    });

    test('works with large arrays recursively', () => {
      const arr = Array.from({ length: 1000 }, (_, i) => i * 2);
      expect(binarySearchRecursive(arr, 500)).toBe(250);
    });
  });

  describe('binarySearchFirst function', () => {
    test('binarySearchFirst function exists', () => {
      expect(typeof binarySearchFirst).toEqual('function');
    });

    test('finds first occurrence of duplicate elements', () => {
      expect(binarySearchFirst([1, 2, 2, 2, 3, 4], 2)).toBe(1);
    });

    test('finds single occurrence', () => {
      expect(binarySearchFirst([1, 2, 3, 4, 5], 3)).toBe(2);
    });

    test('returns -1 for element not found', () => {
      expect(binarySearchFirst([1, 2, 3, 4, 5], 6)).toBe(-1);
    });

    test('handles array with all same elements', () => {
      expect(binarySearchFirst([2, 2, 2, 2, 2], 2)).toBe(0);
    });
  });

  describe('binarySearchLast function', () => {
    test('binarySearchLast function exists', () => {
      expect(typeof binarySearchLast).toEqual('function');
    });

    test('finds last occurrence of duplicate elements', () => {
      expect(binarySearchLast([1, 2, 2, 2, 3, 4], 2)).toBe(3);
    });

    test('finds single occurrence', () => {
      expect(binarySearchLast([1, 2, 3, 4, 5], 3)).toBe(2);
    });

    test('returns -1 for element not found', () => {
      expect(binarySearchLast([1, 2, 3, 4, 5], 6)).toBe(-1);
    });

    test('handles array with all same elements', () => {
      expect(binarySearchLast([2, 2, 2, 2, 2], 2)).toBe(4);
    });
  });

  describe('binarySearchInsertionPoint function', () => {
    test('binarySearchInsertionPoint function exists', () => {
      expect(typeof binarySearchInsertionPoint).toEqual('function');
    });

    test('finds insertion point for element smaller than all', () => {
      expect(binarySearchInsertionPoint([1, 3, 5, 7], 0)).toBe(0);
    });

    test('finds insertion point for element larger than all', () => {
      expect(binarySearchInsertionPoint([1, 3, 5, 7], 10)).toBe(4);
    });

    test('finds insertion point for element in middle', () => {
      expect(binarySearchInsertionPoint([1, 3, 5, 7], 4)).toBe(2);
    });

    test('finds insertion point for existing element', () => {
      expect(binarySearchInsertionPoint([1, 3, 5, 7], 5)).toBe(2);
    });

    test('handles empty array', () => {
      expect(binarySearchInsertionPoint([], 5)).toBe(0);
    });
  });

  describe('searchRotatedArray function', () => {
    test('searchRotatedArray function exists', () => {
      expect(typeof searchRotatedArray).toEqual('function');
    });

    test('finds element in rotated array', () => {
      expect(searchRotatedArray([4, 5, 6, 7, 0, 1, 2], 0)).toBe(4);
    });

    test('finds element in left part of rotated array', () => {
      expect(searchRotatedArray([4, 5, 6, 7, 0, 1, 2], 6)).toBe(2);
    });

    test('finds element in right part of rotated array', () => {
      expect(searchRotatedArray([4, 5, 6, 7, 0, 1, 2], 1)).toBe(5);
    });

    test('returns -1 for element not in rotated array', () => {
      expect(searchRotatedArray([4, 5, 6, 7, 0, 1, 2], 3)).toBe(-1);
    });

    test('handles non-rotated array', () => {
      expect(searchRotatedArray([1, 2, 3, 4, 5], 3)).toBe(2);
    });
  });

  describe('findPeakElement function', () => {
    test('findPeakElement function exists', () => {
      expect(typeof findPeakElement).toEqual('function');
    });

    test('finds peak element in array', () => {
      const arr = [1, 2, 3, 1];
      const peak = findPeakElement(arr);
      expect(peak).toBe(2);
      expect(arr[peak]).toBe(3);
    });

    test('finds peak when peak is at beginning', () => {
      const arr = [3, 2, 1];
      const peak = findPeakElement(arr);
      expect(peak).toBe(0);
    });

    test('finds peak when peak is at end', () => {
      const arr = [1, 2, 3];
      const peak = findPeakElement(arr);
      expect(peak).toBe(2);
    });

    test('handles single element', () => {
      expect(findPeakElement([42])).toBe(0);
    });
  });

  describe('searchMatrix function', () => {
    test('searchMatrix function exists', () => {
      expect(typeof searchMatrix).toEqual('function');
    });

    test('finds element in 2D matrix', () => {
      const matrix = [
        [1,  4,  7,  11],
        [2,  5,  8,  12],
        [3,  6,  9,  16],
        [10, 13, 14, 17]
      ];
      expect(searchMatrix(matrix, 5)).toBe(true);
    });

    test('returns false for element not in matrix', () => {
      const matrix = [
        [1,  4,  7,  11],
        [2,  5,  8,  12],
        [3,  6,  9,  16],
        [10, 13, 14, 17]
      ];
      expect(searchMatrix(matrix, 15)).toBe(false);
    });

    test('handles empty matrix', () => {
      expect(searchMatrix([], 1)).toBe(false);
      expect(searchMatrix([[]], 1)).toBe(false);
    });

    test('finds element at corners', () => {
      const matrix = [
        [1,  4,  7],
        [2,  5,  8],
        [3,  6,  9]
      ];
      expect(searchMatrix(matrix, 1)).toBe(true); // Top-left
      expect(searchMatrix(matrix, 9)).toBe(true); // Bottom-right
    });
  });

  describe('Performance Tests', () => {
    test('handles large arrays efficiently', () => {
      const size = 100000;
      const arr = Array.from({ length: size }, (_, i) => i);
      
      const start = Date.now();
      const result = binarySearch(arr, size - 1);
      const end = Date.now();
      
      expect(result).toBe(size - 1);
      expect(end - start).toBeLessThan(10); // Should be very fast
    });

    test('recursive version handles reasonable depth', () => {
      const arr = Array.from({ length: 1000 }, (_, i) => i * 2);
      expect(binarySearchRecursive(arr, 1000)).toBe(500);
    });
  });
});
