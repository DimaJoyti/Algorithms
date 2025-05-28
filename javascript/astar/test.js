const {
  astar,
  astarWithCost,
  manhattanDistance,
  euclideanDistance,
  diagonalDistance,
  Node,
  getNeighbors,
  reconstructPath
} = require('./index');

describe('A* Algorithm', () => {
  const simpleGrid = [
    [0, 0, 0, 0],
    [1, 1, 0, 1],
    [0, 0, 0, 0],
    [0, 1, 1, 0]
  ];

  describe('Node class', () => {
    test('Node class exists', () => {
      expect(typeof Node).toEqual('function');
    });

    test('creates node with correct properties', () => {
      const node = new Node(1, 2, 5, 3);
      expect(node.x).toBe(1);
      expect(node.y).toBe(2);
      expect(node.g).toBe(5);
      expect(node.h).toBe(3);
      expect(node.f).toBe(8);
      expect(node.parent).toBeNull();
    });

    test('equals method works correctly', () => {
      const node1 = new Node(1, 2);
      const node2 = new Node(1, 2);
      const node3 = new Node(2, 1);

      expect(node1.equals(node2)).toBe(true);
      expect(node1.equals(node3)).toBe(false);
    });

    test('toString method works correctly', () => {
      const node = new Node(3, 4);
      expect(node.toString()).toBe('3,4');
    });
  });

  describe('Distance functions', () => {
    test('manhattanDistance function exists', () => {
      expect(typeof manhattanDistance).toEqual('function');
    });

    test('calculates Manhattan distance correctly', () => {
      const node1 = new Node(0, 0);
      const node2 = new Node(3, 4);
      expect(manhattanDistance(node1, node2)).toBe(7);
    });

    test('euclideanDistance function exists', () => {
      expect(typeof euclideanDistance).toEqual('function');
    });

    test('calculates Euclidean distance correctly', () => {
      const node1 = new Node(0, 0);
      const node2 = new Node(3, 4);
      expect(euclideanDistance(node1, node2)).toBe(5);
    });

    test('diagonalDistance function exists', () => {
      expect(typeof diagonalDistance).toEqual('function');
    });

    test('calculates diagonal distance correctly', () => {
      const node1 = new Node(0, 0);
      const node2 = new Node(3, 4);
      expect(diagonalDistance(node1, node2)).toBe(4);
    });
  });

  describe('getNeighbors function', () => {
    test('getNeighbors function exists', () => {
      expect(typeof getNeighbors).toEqual('function');
    });

    test('gets valid neighbors without diagonal movement', () => {
      const node = new Node(1, 1);
      const neighbors = getNeighbors(simpleGrid, node, false);

      // Position (1,1) in simpleGrid has neighbors: up(0,1), right(1,2), down(2,1)
      // Left(1,0) is blocked (value 1)
      expect(neighbors.length).toBe(3);
      expect(neighbors.some(n => n.x === 0 && n.y === 1)).toBe(true); // up
      expect(neighbors.some(n => n.x === 1 && n.y === 2)).toBe(true); // right
      expect(neighbors.some(n => n.x === 2 && n.y === 1)).toBe(true); // down
    });

    test('gets valid neighbors with diagonal movement', () => {
      const node = new Node(2, 2);
      const neighbors = getNeighbors(simpleGrid, node, true);

      expect(neighbors.length).toBeGreaterThan(0);
      // Should include diagonal neighbors where possible
    });

    test('handles edge cases', () => {
      const cornerNode = new Node(0, 0);
      const neighbors = getNeighbors(simpleGrid, cornerNode, false);

      expect(neighbors.length).toBe(1); // Only one valid neighbor
      expect(neighbors[0].x === 0 && neighbors[0].y === 1).toBe(true);
    });
  });

  describe('reconstructPath function', () => {
    test('reconstructPath function exists', () => {
      expect(typeof reconstructPath).toEqual('function');
    });

    test('reconstructs path correctly', () => {
      const start = new Node(0, 0);
      const middle = new Node(0, 1, 0, 0, start);
      const end = new Node(0, 2, 0, 0, middle);

      const path = reconstructPath(end);
      expect(path).toEqual([[0, 0], [0, 1], [0, 2]]);
    });

    test('handles single node path', () => {
      const node = new Node(1, 1);
      const path = reconstructPath(node);
      expect(path).toEqual([[1, 1]]);
    });
  });

  describe('astar function', () => {
    test('astar function exists', () => {
      expect(typeof astar).toEqual('function');
    });

    test('finds path in simple grid', () => {
      const path = astar(simpleGrid, [0, 0], [3, 3]);

      expect(path).not.toBeNull();
      expect(path.length).toBeGreaterThan(0);
      expect(path[0]).toEqual([0, 0]);
      expect(path[path.length - 1]).toEqual([3, 3]);
    });

    test('returns null for blocked path', () => {
      const blockedGrid = [
        [0, 1, 0],
        [0, 1, 0],
        [0, 1, 0]
      ];

      const path = astar(blockedGrid, [0, 0], [0, 2]);
      expect(path).toBeNull();
    });

    test('handles same start and goal', () => {
      const path = astar(simpleGrid, [0, 0], [0, 0]);
      expect(path).toEqual([[0, 0]]);
    });

    test('returns null for out of bounds start', () => {
      const path = astar(simpleGrid, [-1, 0], [3, 3]);
      expect(path).toBeNull();
    });

    test('returns null for out of bounds goal', () => {
      const path = astar(simpleGrid, [0, 0], [10, 10]);
      expect(path).toBeNull();
    });

    test('returns null for blocked start', () => {
      const path = astar(simpleGrid, [1, 0], [3, 3]);
      expect(path).toBeNull();
    });

    test('returns null for blocked goal', () => {
      const path = astar(simpleGrid, [0, 0], [1, 1]);
      expect(path).toBeNull();
    });

    test('handles empty grid', () => {
      const path = astar([], [0, 0], [1, 1]);
      expect(path).toBeNull();
    });

    test('works with different heuristics', () => {
      const pathManhattan = astar(simpleGrid, [0, 0], [3, 3], {
        heuristic: manhattanDistance
      });

      const pathEuclidean = astar(simpleGrid, [0, 0], [3, 3], {
        heuristic: euclideanDistance
      });

      expect(pathManhattan).not.toBeNull();
      expect(pathEuclidean).not.toBeNull();
    });

    test('works with diagonal movement', () => {
      const path = astar(simpleGrid, [0, 0], [3, 3], {
        allowDiagonal: true
      });

      expect(path).not.toBeNull();
      expect(path.length).toBeGreaterThan(0);
    });

    test('works with weighted heuristic', () => {
      const path = astar(simpleGrid, [0, 0], [3, 3], {
        weight: 2
      });

      expect(path).not.toBeNull();
      expect(path.length).toBeGreaterThan(0);
    });
  });

  describe('astarWithCost function', () => {
    test('astarWithCost function exists', () => {
      expect(typeof astarWithCost).toEqual('function');
    });

    test('works with custom cost function', () => {
      const costFunction = (from, to, grid) => {
        // Simple cost function - higher cost for certain cells
        return grid[to.x][to.y] === 0 ? 1 : 10;
      };

      const path = astarWithCost(simpleGrid, [0, 0], [3, 3], costFunction);
      expect(path).not.toBeNull();
    });

    test('handles same start and goal with cost function', () => {
      const costFunction = () => 1;
      const path = astarWithCost(simpleGrid, [0, 0], [0, 0], costFunction);
      expect(path).toEqual([[0, 0]]);
    });
  });

  describe('Performance and Edge Cases', () => {
    test('handles larger grids efficiently', () => {
      const largeGrid = Array(20).fill().map(() => Array(20).fill(0));

      const start = Date.now();
      const path = astar(largeGrid, [0, 0], [19, 19]);
      const end = Date.now();

      expect(path).not.toBeNull();
      expect(end - start).toBeLessThan(100); // Should complete reasonably fast
    });

    test('handles grid with many obstacles', () => {
      const obstacleGrid = [
        [0, 1, 0, 0, 0],
        [0, 1, 0, 1, 0],
        [0, 0, 0, 1, 0],
        [1, 1, 0, 1, 0],
        [0, 0, 0, 0, 0]
      ];

      const path = astar(obstacleGrid, [0, 0], [4, 4]);
      expect(path).not.toBeNull();
    });

    test('finds optimal path', () => {
      const grid = [
        [0, 0, 0],
        [0, 0, 0],
        [0, 0, 0]
      ];

      const path = astar(grid, [0, 0], [2, 2]);
      expect(path).not.toBeNull();
      expect(path.length).toBe(5); // Optimal path length without diagonal
    });
  });
});
