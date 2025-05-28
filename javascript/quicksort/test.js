const {
  quickSort,
  quickSortRandomized,
  quickSortPure,
  partition,
  swap
} = require('./index');

describe('Quick Sort Algorithm', () => {
  describe('quickSort function', () => {
    test('quickSort function exists', () => {
      expect(typeof quickSort).toEqual('function');
    });

    test('sorts an array of numbers', () => {
      const arr = [3, 1, 4, 1, 5, 9, 2, 6];
      const result = quickSort([...arr]);
      expect(result).toEqual([1, 1, 2, 3, 4, 5, 6, 9]);
    });

    test('sorts a reverse sorted array', () => {
      const arr = [5, 4, 3, 2, 1];
      const result = quickSort([...arr]);
      expect(result).toEqual([1, 2, 3, 4, 5]);
    });

    test('handles already sorted array', () => {
      const arr = [1, 2, 3, 4, 5];
      const result = quickSort([...arr]);
      expect(result).toEqual([1, 2, 3, 4, 5]);
    });

    test('handles array with single element', () => {
      const arr = [42];
      const result = quickSort([...arr]);
      expect(result).toEqual([42]);
    });

    test('handles empty array', () => {
      const arr = [];
      const result = quickSort([...arr]);
      expect(result).toEqual([]);
    });

    test('handles array with duplicate elements', () => {
      const arr = [3, 1, 3, 1, 3, 1];
      const result = quickSort([...arr]);
      expect(result).toEqual([1, 1, 1, 3, 3, 3]);
    });

    test('handles array with negative numbers', () => {
      const arr = [-3, 1, -4, 1, 5, -9, 2, 6];
      const result = quickSort([...arr]);
      expect(result).toEqual([-9, -4, -3, 1, 1, 2, 5, 6]);
    });

    test('handles array with zeros', () => {
      const arr = [0, -1, 0, 1, 0];
      const result = quickSort([...arr]);
      expect(result).toEqual([-1, 0, 0, 0, 1]);
    });

    test('modifies original array', () => {
      const arr = [3, 1, 4, 1, 5];
      quickSort(arr);
      expect(arr).toEqual([1, 1, 3, 4, 5]);
    });
  });

  describe('quickSortRandomized function', () => {
    test('quickSortRandomized function exists', () => {
      expect(typeof quickSortRandomized).toEqual('function');
    });

    test('sorts an array correctly with randomized pivot', () => {
      const arr = [3, 1, 4, 1, 5, 9, 2, 6];
      const result = quickSortRandomized([...arr]);
      expect(result).toEqual([1, 1, 2, 3, 4, 5, 6, 9]);
    });

    test('handles large array with randomized pivot', () => {
      const arr = Array.from({ length: 100 }, () => Math.floor(Math.random() * 100));
      const original = [...arr];
      const result = quickSortRandomized(arr);
      const expected = original.sort((a, b) => a - b);
      expect(result).toEqual(expected);
    });
  });

  describe('quickSortPure function', () => {
    test('quickSortPure function exists', () => {
      expect(typeof quickSortPure).toEqual('function');
    });

    test('sorts array without modifying original', () => {
      const arr = [3, 1, 4, 1, 5, 9, 2, 6];
      const original = [...arr];
      const result = quickSortPure(arr);
      
      expect(result).toEqual([1, 1, 2, 3, 4, 5, 6, 9]);
      expect(arr).toEqual(original); // Original array unchanged
    });

    test('handles empty array', () => {
      const result = quickSortPure([]);
      expect(result).toEqual([]);
    });

    test('handles single element array', () => {
      const result = quickSortPure([42]);
      expect(result).toEqual([42]);
    });

    test('handles array with all same elements', () => {
      const result = quickSortPure([5, 5, 5, 5, 5]);
      expect(result).toEqual([5, 5, 5, 5, 5]);
    });
  });

  describe('partition function', () => {
    test('partition function exists', () => {
      expect(typeof partition).toEqual('function');
    });

    test('partitions array correctly', () => {
      const arr = [3, 1, 4, 1, 5, 9, 2, 6];
      const pivotIndex = partition(arr, 0, arr.length - 1);
      
      // Check that pivot is in correct position
      expect(typeof pivotIndex).toBe('number');
      expect(pivotIndex).toBeGreaterThanOrEqual(0);
      expect(pivotIndex).toBeLessThan(arr.length);
      
      // Check that all elements to the left are <= pivot
      // and all elements to the right are > pivot
      const pivot = arr[pivotIndex];
      for (let i = 0; i < pivotIndex; i++) {
        expect(arr[i]).toBeLessThanOrEqual(pivot);
      }
      for (let i = pivotIndex + 1; i < arr.length; i++) {
        expect(arr[i]).toBeGreaterThan(pivot);
      }
    });

    test('handles single element partition', () => {
      const arr = [42];
      const pivotIndex = partition(arr, 0, 0);
      expect(pivotIndex).toBe(0);
      expect(arr[0]).toBe(42);
    });
  });

  describe('swap function', () => {
    test('swap function exists', () => {
      expect(typeof swap).toEqual('function');
    });

    test('swaps two elements in array', () => {
      const arr = [1, 2, 3, 4, 5];
      swap(arr, 0, 4);
      expect(arr).toEqual([5, 2, 3, 4, 1]);
    });

    test('swaps same element (no change)', () => {
      const arr = [1, 2, 3, 4, 5];
      swap(arr, 2, 2);
      expect(arr).toEqual([1, 2, 3, 4, 5]);
    });

    test('swaps adjacent elements', () => {
      const arr = [1, 2, 3, 4, 5];
      swap(arr, 1, 2);
      expect(arr).toEqual([1, 3, 2, 4, 5]);
    });
  });

  describe('Performance and Edge Cases', () => {
    test('handles large arrays efficiently', () => {
      const size = 1000;
      const arr = Array.from({ length: size }, (_, i) => size - i);
      const start = Date.now();
      const result = quickSort([...arr]);
      const end = Date.now();
      
      expect(result).toEqual(Array.from({ length: size }, (_, i) => i + 1));
      expect(end - start).toBeLessThan(100); // Should complete in reasonable time
    });

    test('handles array with floating point numbers', () => {
      const arr = [3.14, 1.41, 2.71, 1.73, 0.57];
      const result = quickSortPure(arr);
      expect(result).toEqual([0.57, 1.41, 1.73, 2.71, 3.14]);
    });

    test('maintains stability for equal elements (relative order)', () => {
      // Note: Quick sort is not stable, but we test that equal elements are grouped together
      const arr = [3, 1, 3, 1, 3, 1];
      const result = quickSort([...arr]);
      expect(result).toEqual([1, 1, 1, 3, 3, 3]);
    });
  });
});
