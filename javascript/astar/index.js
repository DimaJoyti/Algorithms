// --- Directions
// Implement A* (A-star) pathfinding algorithm.
// A* is a graph traversal and path search algorithm that finds the shortest path
// between nodes using heuristics to guide the search.
// --- Examples
//   const grid = [
//     [0, 0, 0, 0],
//     [1, 1, 0, 1],
//     [0, 0, 0, 0],
//     [0, 1, 1, 0]
//   ];
//   astar(grid, [0, 0], [3, 3]) --> [[0,0], [0,1], [0,2], [1,2], [2,2], [2,3], [3,3]]
// --- Complexity
//   Time: O(b^d) where b is branching factor and d is depth
//   Space: O(b^d)

/**
 * Node class for A* algorithm
 */
class Node {
  constructor(x, y, g = 0, h = 0, parent = null) {
    this.x = x;
    this.y = y;
    this.g = g; // Distance from start
    this.h = h; // Heuristic (estimated distance to goal)
    this.f = g + h; // Total cost
    this.parent = parent;
  }

  equals(other) {
    return this.x === other.x && this.y === other.y;
  }

  toString() {
    return `${this.x},${this.y}`;
  }
}

/**
 * Manhattan distance heuristic
 * @param {Node} node - Current node
 * @param {Node} goal - Goal node
 * @returns {number} Manhattan distance
 */
function manhattanDistance(node, goal) {
  return Math.abs(node.x - goal.x) + Math.abs(node.y - goal.y);
}

/**
 * Euclidean distance heuristic
 * @param {Node} node - Current node
 * @param {Node} goal - Goal node
 * @returns {number} Euclidean distance
 */
function euclideanDistance(node, goal) {
  return Math.sqrt(Math.pow(node.x - goal.x, 2) + Math.pow(node.y - goal.y, 2));
}

/**
 * Diagonal distance heuristic (Chebyshev distance)
 * @param {Node} node - Current node
 * @param {Node} goal - Goal node
 * @returns {number} Diagonal distance
 */
function diagonalDistance(node, goal) {
  const dx = Math.abs(node.x - goal.x);
  const dy = Math.abs(node.y - goal.y);
  return Math.max(dx, dy);
}

/**
 * Get valid neighbors for a node in a grid
 * @param {number[][]} grid - 2D grid (0 = walkable, 1 = obstacle)
 * @param {Node} node - Current node
 * @param {boolean} allowDiagonal - Allow diagonal movement
 * @returns {Node[]} Array of valid neighbor nodes
 */
function getNeighbors(grid, node, allowDiagonal = false) {
  const neighbors = [];
  const directions = allowDiagonal 
    ? [[-1,-1], [-1,0], [-1,1], [0,-1], [0,1], [1,-1], [1,0], [1,1]]
    : [[-1,0], [1,0], [0,-1], [0,1]];

  for (const [dx, dy] of directions) {
    const x = node.x + dx;
    const y = node.y + dy;

    if (x >= 0 && x < grid.length && 
        y >= 0 && y < grid[0].length && 
        grid[x][y] === 0) {
      
      // For diagonal movement, check if path is not blocked
      if (allowDiagonal && dx !== 0 && dy !== 0) {
        if (grid[node.x + dx][node.y] === 1 || grid[node.x][node.y + dy] === 1) {
          continue; // Skip diagonal if adjacent cells are blocked
        }
      }
      
      neighbors.push(new Node(x, y));
    }
  }

  return neighbors;
}

/**
 * Reconstruct path from goal to start
 * @param {Node} node - Goal node with parent chain
 * @returns {number[][]} Path as array of [x, y] coordinates
 */
function reconstructPath(node) {
  const path = [];
  let current = node;
  
  while (current) {
    path.unshift([current.x, current.y]);
    current = current.parent;
  }
  
  return path;
}

/**
 * A* pathfinding algorithm
 * @param {number[][]} grid - 2D grid (0 = walkable, 1 = obstacle)
 * @param {number[]} start - Start position [x, y]
 * @param {number[]} goal - Goal position [x, y]
 * @param {Object} options - Algorithm options
 * @returns {number[][]|null} Path as array of coordinates or null if no path
 */
function astar(grid, start, goal, options = {}) {
  const {
    heuristic = manhattanDistance,
    allowDiagonal = false,
    weight = 1
  } = options;

  if (!grid || grid.length === 0 || grid[0].length === 0) {
    return null;
  }

  const startNode = new Node(start[0], start[1]);
  const goalNode = new Node(goal[0], goal[1]);

  // Check if start or goal is out of bounds or blocked
  if (start[0] < 0 || start[0] >= grid.length || 
      start[1] < 0 || start[1] >= grid[0].length ||
      goal[0] < 0 || goal[0] >= grid.length || 
      goal[1] < 0 || goal[1] >= grid[0].length ||
      grid[start[0]][start[1]] === 1 || 
      grid[goal[0]][goal[1]] === 1) {
    return null;
  }

  // If start equals goal
  if (startNode.equals(goalNode)) {
    return [[start[0], start[1]]];
  }

  const openSet = [startNode];
  const closedSet = new Set();
  const openSetMap = new Map();
  openSetMap.set(startNode.toString(), startNode);

  while (openSet.length > 0) {
    // Find node with lowest f score
    let current = openSet[0];
    let currentIndex = 0;
    
    for (let i = 1; i < openSet.length; i++) {
      if (openSet[i].f < current.f) {
        current = openSet[i];
        currentIndex = i;
      }
    }

    // Remove current from open set
    openSet.splice(currentIndex, 1);
    openSetMap.delete(current.toString());
    closedSet.add(current.toString());

    // Check if we reached the goal
    if (current.equals(goalNode)) {
      return reconstructPath(current);
    }

    // Check all neighbors
    const neighbors = getNeighbors(grid, current, allowDiagonal);
    
    for (const neighbor of neighbors) {
      const neighborKey = neighbor.toString();
      
      if (closedSet.has(neighborKey)) {
        continue;
      }

      // Calculate movement cost
      const movementCost = allowDiagonal && 
        (Math.abs(neighbor.x - current.x) + Math.abs(neighbor.y - current.y) === 2) 
        ? Math.sqrt(2) : 1;
      
      const tentativeG = current.g + movementCost;
      
      const existingNode = openSetMap.get(neighborKey);
      
      if (!existingNode) {
        // New node
        neighbor.g = tentativeG;
        neighbor.h = heuristic(neighbor, goalNode) * weight;
        neighbor.f = neighbor.g + neighbor.h;
        neighbor.parent = current;
        
        openSet.push(neighbor);
        openSetMap.set(neighborKey, neighbor);
      } else if (tentativeG < existingNode.g) {
        // Better path found
        existingNode.g = tentativeG;
        existingNode.f = existingNode.g + existingNode.h;
        existingNode.parent = current;
      }
    }
  }

  return null; // No path found
}

/**
 * A* with custom cost function
 * @param {number[][]} grid - 2D grid with terrain costs
 * @param {number[]} start - Start position [x, y]
 * @param {number[]} goal - Goal position [x, y]
 * @param {Function} costFunction - Function to calculate movement cost
 * @returns {number[][]|null} Path or null if no path found
 */
function astarWithCost(grid, start, goal, costFunction) {
  const startNode = new Node(start[0], start[1]);
  const goalNode = new Node(goal[0], goal[1]);

  if (startNode.equals(goalNode)) {
    return [[start[0], start[1]]];
  }

  const openSet = [startNode];
  const closedSet = new Set();
  const openSetMap = new Map();
  openSetMap.set(startNode.toString(), startNode);

  while (openSet.length > 0) {
    let current = openSet.reduce((min, node) => node.f < min.f ? node : min);
    const currentIndex = openSet.indexOf(current);
    
    openSet.splice(currentIndex, 1);
    openSetMap.delete(current.toString());
    closedSet.add(current.toString());

    if (current.equals(goalNode)) {
      return reconstructPath(current);
    }

    const neighbors = getNeighbors(grid, current, true);
    
    for (const neighbor of neighbors) {
      const neighborKey = neighbor.toString();
      
      if (closedSet.has(neighborKey)) continue;

      const movementCost = costFunction(current, neighbor, grid);
      const tentativeG = current.g + movementCost;
      
      const existingNode = openSetMap.get(neighborKey);
      
      if (!existingNode) {
        neighbor.g = tentativeG;
        neighbor.h = manhattanDistance(neighbor, goalNode);
        neighbor.f = neighbor.g + neighbor.h;
        neighbor.parent = current;
        
        openSet.push(neighbor);
        openSetMap.set(neighborKey, neighbor);
      } else if (tentativeG < existingNode.g) {
        existingNode.g = tentativeG;
        existingNode.f = existingNode.g + existingNode.h;
        existingNode.parent = current;
      }
    }
  }

  return null;
}

module.exports = {
  astar,
  astarWithCost,
  manhattanDistance,
  euclideanDistance,
  diagonalDistance,
  Node,
  getNeighbors,
  reconstructPath
};
