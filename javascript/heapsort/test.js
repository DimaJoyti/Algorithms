const {
  heapSort,
  heapify,
  swap,
  buildMaxHeap,
  extractMax,
  insertIntoHeap,
  isMaxHeap
} = require('./index');

describe('Heap Sort Algorithm', () => {
  describe('heapSort function', () => {
    test('heapSort function exists', () => {
      expect(typeof heapSort).toEqual('function');
    });

    test('sorts an array of numbers', () => {
      const arr = [3, 1, 4, 1, 5, 9, 2, 6];
      const result = heapSort([...arr]);
      expect(result).toEqual([1, 1, 2, 3, 4, 5, 6, 9]);
    });

    test('sorts a reverse sorted array', () => {
      const arr = [5, 4, 3, 2, 1];
      const result = heapSort([...arr]);
      expect(result).toEqual([1, 2, 3, 4, 5]);
    });

    test('handles already sorted array', () => {
      const arr = [1, 2, 3, 4, 5];
      const result = heapSort([...arr]);
      expect(result).toEqual([1, 2, 3, 4, 5]);
    });

    test('handles array with single element', () => {
      const arr = [42];
      const result = heapSort([...arr]);
      expect(result).toEqual([42]);
    });

    test('handles empty array', () => {
      const arr = [];
      const result = heapSort([...arr]);
      expect(result).toEqual([]);
    });

    test('handles array with duplicate elements', () => {
      const arr = [3, 1, 3, 1, 3, 1];
      const result = heapSort([...arr]);
      expect(result).toEqual([1, 1, 1, 3, 3, 3]);
    });

    test('handles array with negative numbers', () => {
      const arr = [-3, 1, -4, 1, 5, -9, 2, 6];
      const result = heapSort([...arr]);
      expect(result).toEqual([-9, -4, -3, 1, 1, 2, 5, 6]);
    });

    test('modifies original array', () => {
      const arr = [3, 1, 4, 1, 5];
      heapSort(arr);
      expect(arr).toEqual([1, 1, 3, 4, 5]);
    });
  });

  describe('heapify function', () => {
    test('heapify function exists', () => {
      expect(typeof heapify).toEqual('function');
    });

    test('heapifies a subtree correctly', () => {
      const arr = [1, 3, 6, 5, 2, 4];
      heapify(arr, arr.length, 0);
      
      // Root should be the largest among root, left, and right
      expect(arr[0]).toBeGreaterThanOrEqual(arr[1]);
      expect(arr[0]).toBeGreaterThanOrEqual(arr[2]);
    });

    test('handles single element', () => {
      const arr = [5];
      heapify(arr, 1, 0);
      expect(arr).toEqual([5]);
    });
  });

  describe('buildMaxHeap function', () => {
    test('buildMaxHeap function exists', () => {
      expect(typeof buildMaxHeap).toEqual('function');
    });

    test('builds a valid max heap', () => {
      const arr = [3, 1, 4, 1, 5, 9, 2, 6];
      const heap = buildMaxHeap(arr);
      
      expect(isMaxHeap(heap)).toBe(true);
      expect(heap.length).toBe(arr.length);
    });

    test('handles empty array', () => {
      const heap = buildMaxHeap([]);
      expect(heap).toEqual([]);
    });

    test('handles single element', () => {
      const heap = buildMaxHeap([42]);
      expect(heap).toEqual([42]);
    });
  });

  describe('extractMax function', () => {
    test('extractMax function exists', () => {
      expect(typeof extractMax).toEqual('function');
    });

    test('extracts maximum from heap', () => {
      const heap = [9, 5, 6, 1, 3, 4, 2];
      const max = extractMax(heap);
      
      expect(max).toBe(9);
      expect(heap.length).toBe(6);
      expect(isMaxHeap(heap)).toBe(true);
    });

    test('handles empty heap', () => {
      const heap = [];
      const max = extractMax(heap);
      expect(max).toBeNull();
    });

    test('handles single element heap', () => {
      const heap = [42];
      const max = extractMax(heap);
      expect(max).toBe(42);
      expect(heap.length).toBe(0);
    });

    test('extracts all elements in descending order', () => {
      const heap = buildMaxHeap([3, 1, 4, 1, 5, 9, 2, 6]);
      const extracted = [];
      
      while (heap.length > 0) {
        extracted.push(extractMax(heap));
      }
      
      expect(extracted).toEqual([9, 6, 5, 4, 3, 2, 1, 1]);
    });
  });

  describe('insertIntoHeap function', () => {
    test('insertIntoHeap function exists', () => {
      expect(typeof insertIntoHeap).toEqual('function');
    });

    test('inserts element into heap maintaining heap property', () => {
      const heap = [9, 5, 6, 1, 3, 4, 2];
      insertIntoHeap(heap, 8);
      
      expect(heap.length).toBe(8);
      expect(isMaxHeap(heap)).toBe(true);
      expect(heap).toContain(8);
    });

    test('inserts into empty heap', () => {
      const heap = [];
      insertIntoHeap(heap, 42);
      
      expect(heap).toEqual([42]);
      expect(isMaxHeap(heap)).toBe(true);
    });

    test('inserts largest element', () => {
      const heap = [5, 3, 4, 1, 2];
      insertIntoHeap(heap, 10);
      
      expect(heap[0]).toBe(10);
      expect(isMaxHeap(heap)).toBe(true);
    });
  });

  describe('isMaxHeap function', () => {
    test('isMaxHeap function exists', () => {
      expect(typeof isMaxHeap).toEqual('function');
    });

    test('returns true for valid max heap', () => {
      const heap = [9, 5, 6, 1, 3, 4, 2];
      expect(isMaxHeap(heap)).toBe(true);
    });

    test('returns false for invalid max heap', () => {
      const notHeap = [1, 5, 6, 9, 3, 4, 2];
      expect(isMaxHeap(notHeap)).toBe(false);
    });

    test('returns true for empty array', () => {
      expect(isMaxHeap([])).toBe(true);
    });

    test('returns true for single element', () => {
      expect(isMaxHeap([42])).toBe(true);
    });

    test('returns true for two elements in correct order', () => {
      expect(isMaxHeap([5, 3])).toBe(true);
    });

    test('returns false for two elements in wrong order', () => {
      expect(isMaxHeap([3, 5])).toBe(false);
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
  });

  describe('Performance and Edge Cases', () => {
    test('handles large arrays efficiently', () => {
      const size = 1000;
      const arr = Array.from({ length: size }, () => Math.floor(Math.random() * 1000));
      const start = Date.now();
      const result = heapSort([...arr]);
      const end = Date.now();
      
      expect(result.length).toBe(size);
      expect(end - start).toBeLessThan(100); // Should complete in reasonable time
      
      // Verify it's sorted
      for (let i = 1; i < result.length; i++) {
        expect(result[i]).toBeGreaterThanOrEqual(result[i - 1]);
      }
    });

    test('maintains stability for heap operations', () => {
      const heap = buildMaxHeap([1, 2, 3, 4, 5]);
      const originalLength = heap.length;
      
      insertIntoHeap(heap, 6);
      expect(heap.length).toBe(originalLength + 1);
      expect(isMaxHeap(heap)).toBe(true);
      
      const max = extractMax(heap);
      expect(max).toBe(6);
      expect(heap.length).toBe(originalLength);
      expect(isMaxHeap(heap)).toBe(true);
    });
  });
});
