const {
  dijkstra,
  shortestPath,
  dijkstraMatrix,
  hasNegativeWeights,
  allPairsShortestPath,
  PriorityQueue
} = require('./index');

describe('Dijkstra Algorithm', () => {
  const sampleGraph = {
    'A': { 'B': 4, 'C': 2 },
    'B': { 'A': 4, 'C': 1, 'D': 5 },
    'C': { 'A': 2, 'B': 1, 'D': 8, 'E': 10 },
    'D': { 'B': 5, 'C': 8, 'E': 2 },
    'E': { 'C': 10, 'D': 2 }
  };

  describe('PriorityQueue', () => {
    test('PriorityQueue exists', () => {
      expect(typeof PriorityQueue).toEqual('function');
    });

    test('enqueue and dequeue work correctly', () => {
      const pq = new PriorityQueue();
      pq.enqueue('A', 3);
      pq.enqueue('B', 1);
      pq.enqueue('C', 2);

      expect(pq.dequeue().val).toBe('B');
      expect(pq.dequeue().val).toBe('C');
      expect(pq.dequeue().val).toBe('A');
    });

    test('isEmpty works correctly', () => {
      const pq = new PriorityQueue();
      expect(pq.isEmpty()).toBe(true);

      pq.enqueue('A', 1);
      expect(pq.isEmpty()).toBe(false);

      pq.dequeue();
      expect(pq.isEmpty()).toBe(true);
    });
  });

  describe('dijkstra function', () => {
    test('dijkstra function exists', () => {
      expect(typeof dijkstra).toEqual('function');
    });

    test('finds shortest distances from starting vertex', () => {
      const result = dijkstra(sampleGraph, 'A');

      expect(result).toEqual({
        'A': 0,
        'B': 3,
        'C': 2,
        'D': 8,
        'E': 10
      });
    });

    test('works with different starting vertex', () => {
      const result = dijkstra(sampleGraph, 'B');

      expect(result['B']).toBe(0);
      expect(result['A']).toBe(3); // B->C->A = 1+2 = 3, not B->A = 4
      expect(result['C']).toBe(1);
      expect(result['D']).toBe(5);
      expect(result['E']).toBe(7);
    });

    test('handles single vertex graph', () => {
      const singleGraph = { 'A': {} };
      const result = dijkstra(singleGraph, 'A');

      expect(result).toEqual({ 'A': 0 });
    });

    test('handles disconnected graph', () => {
      const disconnectedGraph = {
        'A': { 'B': 1 },
        'B': { 'A': 1 },
        'C': { 'D': 2 },
        'D': { 'C': 2 }
      };

      const result = dijkstra(disconnectedGraph, 'A');

      expect(result['A']).toBe(0);
      expect(result['B']).toBe(1);
      expect(result['C']).toBe(Infinity);
      expect(result['D']).toBe(Infinity);
    });
  });

  describe('shortestPath function', () => {
    test('shortestPath function exists', () => {
      expect(typeof shortestPath).toEqual('function');
    });

    test('finds shortest path between two vertices', () => {
      const result = shortestPath(sampleGraph, 'A', 'E');

      expect(result.path).toEqual(['A', 'C', 'B', 'D', 'E']);
      expect(result.distance).toBe(10);
    });

    test('finds direct path when it exists', () => {
      const result = shortestPath(sampleGraph, 'A', 'C');

      expect(result.path).toEqual(['A', 'C']);
      expect(result.distance).toBe(2);
    });

    test('returns empty path for unreachable destination', () => {
      const disconnectedGraph = {
        'A': { 'B': 1 },
        'B': { 'A': 1 },
        'C': {}
      };

      const result = shortestPath(disconnectedGraph, 'A', 'C');

      expect(result.path).toEqual([]);
      expect(result.distance).toBe(-1);
    });

    test('handles same start and end vertex', () => {
      const result = shortestPath(sampleGraph, 'A', 'A');

      expect(result.path).toEqual(['A']);
      expect(result.distance).toBe(0);
    });
  });

  describe('dijkstraMatrix function', () => {
    test('dijkstraMatrix function exists', () => {
      expect(typeof dijkstraMatrix).toEqual('function');
    });

    test('works with adjacency matrix', () => {
      const matrix = [
        [0, 4, 2, 0, 0],  // A
        [4, 0, 1, 5, 0],  // B
        [2, 1, 0, 8, 10], // C
        [0, 5, 8, 0, 2],  // D
        [0, 0, 10, 2, 0]  // E
      ];

      const result = dijkstraMatrix(matrix, 0); // Start from A (index 0)

      expect(result).toEqual([0, 3, 2, 8, 10]);
    });

    test('handles single vertex matrix', () => {
      const matrix = [[0]];
      const result = dijkstraMatrix(matrix, 0);

      expect(result).toEqual([0]);
    });

    test('handles disconnected matrix', () => {
      const matrix = [
        [0, 1, 0],
        [1, 0, 0],
        [0, 0, 0]
      ];

      const result = dijkstraMatrix(matrix, 0);

      expect(result[0]).toBe(0);
      expect(result[1]).toBe(1);
      expect(result[2]).toBe(Infinity);
    });
  });

  describe('hasNegativeWeights function', () => {
    test('hasNegativeWeights function exists', () => {
      expect(typeof hasNegativeWeights).toEqual('function');
    });

    test('returns false for graph with no negative weights', () => {
      expect(hasNegativeWeights(sampleGraph)).toBe(false);
    });

    test('returns true for graph with negative weights', () => {
      const negativeGraph = {
        'A': { 'B': -1, 'C': 2 },
        'B': { 'A': 1 },
        'C': { 'A': 2 }
      };

      expect(hasNegativeWeights(negativeGraph)).toBe(true);
    });

    test('handles empty graph', () => {
      expect(hasNegativeWeights({})).toBe(false);
    });
  });

  describe('allPairsShortestPath function', () => {
    test('allPairsShortestPath function exists', () => {
      expect(typeof allPairsShortestPath).toEqual('function');
    });

    test('computes shortest paths from all vertices', () => {
      const smallGraph = {
        'A': { 'B': 1, 'C': 4 },
        'B': { 'A': 1, 'C': 2 },
        'C': { 'A': 4, 'B': 2 }
      };

      const result = allPairsShortestPath(smallGraph);

      expect(result['A']['A']).toBe(0);
      expect(result['A']['B']).toBe(1);
      expect(result['A']['C']).toBe(3);

      expect(result['B']['A']).toBe(1);
      expect(result['B']['B']).toBe(0);
      expect(result['B']['C']).toBe(2);

      expect(result['C']['A']).toBe(3);
      expect(result['C']['B']).toBe(2);
      expect(result['C']['C']).toBe(0);
    });
  });

  describe('Edge Cases and Performance', () => {
    test('handles large graph efficiently', () => {
      // Create a larger graph
      const largeGraph = {};
      const size = 100;

      for (let i = 0; i < size; i++) {
        largeGraph[i.toString()] = {};
        if (i < size - 1) {
          largeGraph[i.toString()][(i + 1).toString()] = 1;
        }
      }

      const start = Date.now();
      const result = dijkstra(largeGraph, '0');
      const end = Date.now();

      expect(result['0']).toBe(0);
      expect(result['99']).toBe(99);
      expect(end - start).toBeLessThan(100); // Should complete reasonably fast
    });

    test('handles graph with self-loops', () => {
      const selfLoopGraph = {
        'A': { 'A': 0, 'B': 1 },
        'B': { 'A': 1, 'B': 0 }
      };

      const result = dijkstra(selfLoopGraph, 'A');

      expect(result['A']).toBe(0);
      expect(result['B']).toBe(1);
    });

    test('handles graph with multiple edges of same weight', () => {
      const multiGraph = {
        'A': { 'B': 2, 'C': 2 },
        'B': { 'A': 2, 'D': 1 },
        'C': { 'A': 2, 'D': 1 },
        'D': { 'B': 1, 'C': 1 }
      };

      const result = dijkstra(multiGraph, 'A');

      expect(result['A']).toBe(0);
      expect(result['B']).toBe(2);
      expect(result['C']).toBe(2);
      expect(result['D']).toBe(3);
    });
  });
});
