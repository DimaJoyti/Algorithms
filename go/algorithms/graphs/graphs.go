// Package graphs implements various graph algorithms
package graphs

import (
	"container/heap"
	"math"
)

// Graph represents an unweighted graph using adjacency list
type Graph struct {
	vertices int
	adjList  map[int][]int
}

// WeightedGraph represents a weighted graph using adjacency list
type WeightedGraph struct {
	vertices int
	adjList  map[int][]Edge
}

// Edge represents a weighted edge
type Edge struct {
	To     int
	Weight int
}

// PriorityQueueItem represents an item in priority queue
type PriorityQueueItem struct {
	Vertex   int
	Distance int
	Index    int
}

// PriorityQueue implements a min-heap for Dijkstra's algorithm
type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityQueueItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// NewGraph creates a new unweighted graph
func NewGraph(vertices int) *Graph {
	return &Graph{
		vertices: vertices,
		adjList:  make(map[int][]int),
	}
}

// NewWeightedGraph creates a new weighted graph
func NewWeightedGraph(vertices int) *WeightedGraph {
	return &WeightedGraph{
		vertices: vertices,
		adjList:  make(map[int][]Edge),
	}
}

// AddEdge adds an edge to unweighted graph
func (g *Graph) AddEdge(from, to int) {
	g.adjList[from] = append(g.adjList[from], to)
}

// AddEdge adds a weighted edge to weighted graph
func (wg *WeightedGraph) AddEdge(from, to, weight int) {
	wg.adjList[from] = append(wg.adjList[from], Edge{To: to, Weight: weight})
}

// BFS performs Breadth-First Search
// Time Complexity: O(V + E), Space Complexity: O(V)
func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	result := []int{}

	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

// DFS performs Depth-First Search
// Time Complexity: O(V + E), Space Complexity: O(V)
func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	result := []int{}

	g.dfsHelper(start, visited, &result)
	return result
}

func (g *Graph) dfsHelper(vertex int, visited map[int]bool, result *[]int) {
	visited[vertex] = true
	*result = append(*result, vertex)

	for _, neighbor := range g.adjList[vertex] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited, result)
		}
	}
}

// HasCycle detects if the graph has a cycle using DFS
// Time Complexity: O(V + E), Space Complexity: O(V)
func (g *Graph) HasCycle() bool {
	visited := make(map[int]bool)
	recStack := make(map[int]bool)

	for vertex := range g.adjList {
		if !visited[vertex] {
			if g.hasCycleHelper(vertex, visited, recStack) {
				return true
			}
		}
	}

	return false
}

func (g *Graph) hasCycleHelper(vertex int, visited, recStack map[int]bool) bool {
	visited[vertex] = true
	recStack[vertex] = true

	for _, neighbor := range g.adjList[vertex] {
		if !visited[neighbor] {
			if g.hasCycleHelper(neighbor, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			return true
		}
	}

	recStack[vertex] = false
	return false
}

// TopologicalSort performs topological sorting using DFS
// Time Complexity: O(V + E), Space Complexity: O(V)
func (g *Graph) TopologicalSort() []int {
	visited := make(map[int]bool)
	stack := []int{}

	for vertex := range g.adjList {
		if !visited[vertex] {
			g.topologicalSortHelper(vertex, visited, &stack)
		}
	}

	// Reverse the stack
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}

	return stack
}

func (g *Graph) topologicalSortHelper(vertex int, visited map[int]bool, stack *[]int) {
	visited[vertex] = true

	for _, neighbor := range g.adjList[vertex] {
		if !visited[neighbor] {
			g.topologicalSortHelper(neighbor, visited, stack)
		}
	}

	*stack = append(*stack, vertex)
}

// Dijkstra implements Dijkstra's shortest path algorithm
// Time Complexity: O((V + E) log V), Space Complexity: O(V)
func (wg *WeightedGraph) Dijkstra(start int) map[int]int {
	distances := make(map[int]int)
	pq := &PriorityQueue{}

	// Initialize distances
	for vertex := range wg.adjList {
		distances[vertex] = math.MaxInt32
	}
	distances[start] = 0

	heap.Init(pq)
	heap.Push(pq, &PriorityQueueItem{Vertex: start, Distance: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*PriorityQueueItem)

		if current.Distance > distances[current.Vertex] {
			continue
		}

		for _, edge := range wg.adjList[current.Vertex] {
			newDistance := distances[current.Vertex] + edge.Weight

			if newDistance < distances[edge.To] {
				distances[edge.To] = newDistance
				heap.Push(pq, &PriorityQueueItem{Vertex: edge.To, Distance: newDistance})
			}
		}
	}

	return distances
}

// BellmanFord implements Bellman-Ford algorithm for shortest paths with negative weights
// Time Complexity: O(VE), Space Complexity: O(V)
func (wg *WeightedGraph) BellmanFord(start int) (map[int]int, bool) {
	distances := make(map[int]int)

	// Get all vertices
	allVertices := make(map[int]bool)
	allVertices[start] = true
	for vertex := range wg.adjList {
		allVertices[vertex] = true
		for _, edge := range wg.adjList[vertex] {
			allVertices[edge.To] = true
		}
	}

	// Initialize distances
	for vertex := range allVertices {
		distances[vertex] = math.MaxInt32
	}
	distances[start] = 0

	// Relax edges V-1 times
	for i := 0; i < wg.vertices-1; i++ {
		for vertex, edges := range wg.adjList {
			if distances[vertex] == math.MaxInt32 {
				continue
			}

			for _, edge := range edges {
				newDistance := distances[vertex] + edge.Weight
				if newDistance < distances[edge.To] {
					distances[edge.To] = newDistance
				}
			}
		}
	}

	// Check for negative cycles
	for vertex, edges := range wg.adjList {
		if distances[vertex] == math.MaxInt32 {
			continue
		}

		for _, edge := range edges {
			if distances[vertex]+edge.Weight < distances[edge.To] {
				return distances, false // Negative cycle detected
			}
		}
	}

	return distances, true
}

// FloydWarshall implements Floyd-Warshall algorithm for all-pairs shortest paths
// Time Complexity: O(V³), Space Complexity: O(V²)
func (wg *WeightedGraph) FloydWarshall() [][]int {
	// Get all vertices
	allVertices := make(map[int]bool)
	for vertex := range wg.adjList {
		allVertices[vertex] = true
		for _, edge := range wg.adjList[vertex] {
			allVertices[edge.To] = true
		}
	}

	vertices := make([]int, 0, len(allVertices))
	vertexIndex := make(map[int]int)

	// Create vertex mapping
	i := 0
	for vertex := range allVertices {
		vertices = append(vertices, vertex)
		vertexIndex[vertex] = i
		i++
	}

	n := len(vertices)
	dist := make([][]int, n)

	// Initialize distance matrix
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = math.MaxInt32
			}
		}
	}

	// Fill direct edges
	for vertex, edges := range wg.adjList {
		i := vertexIndex[vertex]
		for _, edge := range edges {
			j := vertexIndex[edge.To]
			dist[i][j] = edge.Weight
		}
	}

	// Floyd-Warshall algorithm
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 {
					if dist[i][k]+dist[k][j] < dist[i][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
					}
				}
			}
		}
	}

	return dist
}

// IsConnected checks if the graph is connected
func (g *Graph) IsConnected() bool {
	if len(g.adjList) == 0 {
		return true
	}

	// Get any vertex as starting point
	var start int
	for vertex := range g.adjList {
		start = vertex
		break
	}

	visited := g.BFS(start)
	return len(visited) == len(g.adjList)
}

// FindPath finds a path between two vertices using BFS
func (g *Graph) FindPath(start, end int) []int {
	if start == end {
		return []int{start}
	}

	visited := make(map[int]bool)
	parent := make(map[int]int)
	queue := []int{start}

	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.adjList[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = vertex
				queue = append(queue, neighbor)

				if neighbor == end {
					// Reconstruct path
					path := []int{}
					current := end
					for current != start {
						path = append([]int{current}, path...)
						current = parent[current]
					}
					path = append([]int{start}, path...)
					return path
				}
			}
		}
	}

	return []int{} // No path found
}

// TarjanSCC finds strongly connected components using Tarjan's algorithm
// Time Complexity: O(V + E), Space Complexity: O(V)
func (g *Graph) TarjanSCC() [][]int {
	index := 0
	stack := make([]int, 0)
	indices := make(map[int]int)
	lowlinks := make(map[int]int)
	onStack := make(map[int]bool)
	sccs := make([][]int, 0)

	var strongConnect func(int)
	strongConnect = func(v int) {
		// Set the depth index for v to the smallest unused index
		indices[v] = index
		lowlinks[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		// Consider successors of v
		for _, w := range g.adjList[v] {
			if _, exists := indices[w]; !exists {
				// Successor w has not yet been visited; recurse on it
				strongConnect(w)
				if lowlinks[w] < lowlinks[v] {
					lowlinks[v] = lowlinks[w]
				}
			} else if onStack[w] {
				// Successor w is in stack and hence in the current SCC
				if indices[w] < lowlinks[v] {
					lowlinks[v] = indices[w]
				}
			}
		}

		// If v is a root node, pop the stack and create an SCC
		if lowlinks[v] == indices[v] {
			scc := make([]int, 0)
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			sccs = append(sccs, scc)
		}
	}

	// Find all vertices
	allVertices := make(map[int]bool)
	for vertex := range g.adjList {
		allVertices[vertex] = true
		for _, neighbor := range g.adjList[vertex] {
			allVertices[neighbor] = true
		}
	}

	for vertex := range allVertices {
		if _, exists := indices[vertex]; !exists {
			strongConnect(vertex)
		}
	}

	return sccs
}

// UnionFind data structure for Kruskal's algorithm
type UnionFind struct {
	parent []int
	rank   []int
}

// NewUnionFind creates a new Union-Find structure
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank}
}

// Find finds the root of element x with path compression
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

// Union unites two sets containing x and y
func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // Already in same set
	}

	// Union by rank
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	return true
}

// EdgeForMST represents an edge for MST algorithms
type EdgeForMST struct {
	From   int
	To     int
	Weight int
}

// KruskalMST finds Minimum Spanning Tree using Kruskal's algorithm
// Time Complexity: O(E log E), Space Complexity: O(V)
func (wg *WeightedGraph) KruskalMST() ([]EdgeForMST, int) {
	// Collect all edges
	edges := make([]EdgeForMST, 0)
	for from, adjEdges := range wg.adjList {
		for _, edge := range adjEdges {
			edges = append(edges, EdgeForMST{From: from, To: edge.To, Weight: edge.Weight})
		}
	}

	// Sort edges by weight (simple bubble sort for demonstration)
	for i := 0; i < len(edges)-1; i++ {
		for j := 0; j < len(edges)-i-1; j++ {
			if edges[j].Weight > edges[j+1].Weight {
				edges[j], edges[j+1] = edges[j+1], edges[j]
			}
		}
	}

	uf := NewUnionFind(wg.vertices)
	mst := make([]EdgeForMST, 0)
	totalWeight := 0

	for _, edge := range edges {
		if uf.Union(edge.From, edge.To) {
			mst = append(mst, edge)
			totalWeight += edge.Weight
			if len(mst) == wg.vertices-1 {
				break
			}
		}
	}

	return mst, totalWeight
}

// PrimMST finds Minimum Spanning Tree using Prim's algorithm
// Time Complexity: O(E log V), Space Complexity: O(V)
func (wg *WeightedGraph) PrimMST(start int) ([]EdgeForMST, int) {
	mst := make([]EdgeForMST, 0)
	totalWeight := 0
	visited := make(map[int]bool)

	// Priority queue for edges
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Start with the given vertex
	visited[start] = true

	// Add all edges from start vertex to priority queue
	for _, edge := range wg.adjList[start] {
		heap.Push(pq, &PriorityQueueItem{
			Vertex:   edge.To,
			Distance: edge.Weight,
		})
	}

	for pq.Len() > 0 && len(mst) < wg.vertices-1 {
		item := heap.Pop(pq).(*PriorityQueueItem)
		vertex := item.Vertex
		weight := item.Distance

		if visited[vertex] {
			continue
		}

		visited[vertex] = true

		// Find the parent vertex (the one that led to this vertex with minimum weight)
		parent := -1
		for from, edges := range wg.adjList {
			if visited[from] {
				for _, edge := range edges {
					if edge.To == vertex && edge.Weight == weight {
						parent = from
						break
					}
				}
			}
			if parent != -1 {
				break
			}
		}

		if parent != -1 {
			mst = append(mst, EdgeForMST{From: parent, To: vertex, Weight: weight})
			totalWeight += weight
		}

		// Add all edges from the new vertex to priority queue
		for _, edge := range wg.adjList[vertex] {
			if !visited[edge.To] {
				heap.Push(pq, &PriorityQueueItem{
					Vertex:   edge.To,
					Distance: edge.Weight,
				})
			}
		}
	}

	return mst, totalWeight
}
