// --- Directions
// Implement Dijkstra's shortest path algorithm.
// Dijkstra's algorithm finds the shortest paths between nodes in a graph.
// It works for graphs with non-negative edge weights.
// --- Examples
//   const graph = {
//     'A': { 'B': 4, 'C': 2 },
//     'B': { 'A': 4, 'C': 1, 'D': 5 },
//     'C': { 'A': 2, 'B': 1, 'D': 8, 'E': 10 },
//     'D': { 'B': 5, 'C': 8, 'E': 2 },
//     'E': { 'C': 10, 'D': 2 }
//   };
//   dijkstra(graph, 'A') --> { A: 0, B: 3, C: 2, D: 8, E: 10 }
// --- Complexity
//   Time: O((V + E) log V) with binary heap
//   Space: O(V)

/**
 * Priority Queue implementation using binary heap
 */
class PriorityQueue {
  constructor() {
    this.values = [];
  }

  enqueue(val, priority) {
    this.values.push({ val, priority });
    this.bubbleUp();
  }

  dequeue() {
    const min = this.values[0];
    const end = this.values.pop();
    if (this.values.length > 0) {
      this.values[0] = end;
      this.sinkDown();
    }
    return min;
  }

  bubbleUp() {
    let idx = this.values.length - 1;
    const element = this.values[idx];

    while (idx > 0) {
      let parentIdx = Math.floor((idx - 1) / 2);
      let parent = this.values[parentIdx];

      if (element.priority >= parent.priority) break;

      this.values[parentIdx] = element;
      this.values[idx] = parent;
      idx = parentIdx;
    }
  }

  sinkDown() {
    let idx = 0;
    const length = this.values.length;
    const element = this.values[0];

    while (true) {
      let leftChildIdx = 2 * idx + 1;
      let rightChildIdx = 2 * idx + 2;
      let leftChild, rightChild;
      let swap = null;

      if (leftChildIdx < length) {
        leftChild = this.values[leftChildIdx];
        if (leftChild.priority < element.priority) {
          swap = leftChildIdx;
        }
      }

      if (rightChildIdx < length) {
        rightChild = this.values[rightChildIdx];
        if (
          (swap === null && rightChild.priority < element.priority) ||
          (swap !== null && rightChild.priority < leftChild.priority)
        ) {
          swap = rightChildIdx;
        }
      }

      if (swap === null) break;

      this.values[idx] = this.values[swap];
      this.values[swap] = element;
      idx = swap;
    }
  }

  isEmpty() {
    return this.values.length === 0;
  }
}

/**
 * Dijkstra's shortest path algorithm
 * @param {Object} graph - Adjacency list representation of graph
 * @param {string} start - Starting vertex
 * @returns {Object} Object with shortest distances to all vertices
 */
function dijkstra(graph, start) {
  const distances = {};
  const previous = {};
  const pq = new PriorityQueue();

  // Initialize distances and previous
  for (let vertex in graph) {
    if (vertex === start) {
      distances[vertex] = 0;
      pq.enqueue(vertex, 0);
    } else {
      distances[vertex] = Infinity;
      pq.enqueue(vertex, Infinity);
    }
    previous[vertex] = null;
  }

  while (!pq.isEmpty()) {
    let smallest = pq.dequeue().val;

    if (distances[smallest] === Infinity) break;

    for (let neighbor in graph[smallest]) {
      let candidate = distances[smallest] + graph[smallest][neighbor];

      if (candidate < distances[neighbor]) {
        distances[neighbor] = candidate;
        previous[neighbor] = smallest;
        pq.enqueue(neighbor, candidate);
      }
    }
  }

  return distances;
}

/**
 * Get shortest path between two vertices
 * @param {Object} graph - Adjacency list representation of graph
 * @param {string} start - Starting vertex
 * @param {string} end - Ending vertex
 * @returns {Object} Object with path array and total distance
 */
function shortestPath(graph, start, end) {
  const distances = {};
  const previous = {};
  const pq = new PriorityQueue();
  const path = [];

  // Initialize
  for (let vertex in graph) {
    if (vertex === start) {
      distances[vertex] = 0;
      pq.enqueue(vertex, 0);
    } else {
      distances[vertex] = Infinity;
      pq.enqueue(vertex, Infinity);
    }
    previous[vertex] = null;
  }

  while (!pq.isEmpty()) {
    let smallest = pq.dequeue().val;

    if (smallest === end) {
      // Build path
      while (smallest !== null) {
        path.push(smallest);
        smallest = previous[smallest];
      }
      break;
    }

    if (distances[smallest] === Infinity) break;

    for (let neighbor in graph[smallest]) {
      let candidate = distances[smallest] + graph[smallest][neighbor];

      if (candidate < distances[neighbor]) {
        distances[neighbor] = candidate;
        previous[neighbor] = smallest;
        pq.enqueue(neighbor, candidate);
      }
    }
  }

  path.reverse();

  // If no path was found (distance is Infinity), return empty path
  if (distances[end] === Infinity) {
    return {
      path: [],
      distance: -1
    };
  }

  return {
    path: path.length > 0 ? path : [],
    distance: distances[end]
  };
}

/**
 * Dijkstra for weighted matrix representation
 * @param {number[][]} matrix - Adjacency matrix (0 means no edge)
 * @param {number} start - Starting vertex index
 * @returns {number[]} Array of shortest distances
 */
function dijkstraMatrix(matrix, start) {
  const n = matrix.length;
  const distances = new Array(n).fill(Infinity);
  const visited = new Array(n).fill(false);

  distances[start] = 0;

  for (let count = 0; count < n - 1; count++) {
    // Find minimum distance vertex
    let u = -1;
    for (let v = 0; v < n; v++) {
      if (!visited[v] && (u === -1 || distances[v] < distances[u])) {
        u = v;
      }
    }

    visited[u] = true;

    // Update distances
    for (let v = 0; v < n; v++) {
      if (!visited[v] && matrix[u][v] !== 0 &&
          distances[u] !== Infinity &&
          distances[u] + matrix[u][v] < distances[v]) {
        distances[v] = distances[u] + matrix[u][v];
      }
    }
  }

  return distances;
}

/**
 * Check if graph has negative weights (Dijkstra doesn't work with negative weights)
 * @param {Object} graph - Adjacency list representation
 * @returns {boolean} True if graph has negative weights
 */
function hasNegativeWeights(graph) {
  for (let vertex in graph) {
    for (let neighbor in graph[vertex]) {
      if (graph[vertex][neighbor] < 0) {
        return true;
      }
    }
  }
  return false;
}

/**
 * All pairs shortest path using Dijkstra
 * @param {Object} graph - Adjacency list representation
 * @returns {Object} Object with distances from each vertex to all others
 */
function allPairsShortestPath(graph) {
  const result = {};

  for (let vertex in graph) {
    result[vertex] = dijkstra(graph, vertex);
  }

  return result;
}

module.exports = {
  dijkstra,
  shortestPath,
  dijkstraMatrix,
  hasNegativeWeights,
  allPairsShortestPath,
  PriorityQueue
};
