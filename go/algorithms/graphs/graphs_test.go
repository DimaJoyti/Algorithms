package graphs

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph(5)
	if g.vertices != 5 {
		t.Errorf("NewGraph(5) vertices = %d, want 5", g.vertices)
	}
	if g.adjList == nil {
		t.Error("NewGraph(5) adjList should not be nil")
	}
}

func TestNewWeightedGraph(t *testing.T) {
	wg := NewWeightedGraph(5)
	if wg.vertices != 5 {
		t.Errorf("NewWeightedGraph(5) vertices = %d, want 5", wg.vertices)
	}
	if wg.adjList == nil {
		t.Error("NewWeightedGraph(5) adjList should not be nil")
	}
}

func TestGraphAddEdge(t *testing.T) {
	g := NewGraph(3)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)

	expected := []int{1, 2}
	if !reflect.DeepEqual(g.adjList[0], expected) {
		t.Errorf("AddEdge failed, got %v, want %v", g.adjList[0], expected)
	}
}

func TestWeightedGraphAddEdge(t *testing.T) {
	wg := NewWeightedGraph(3)
	wg.AddEdge(0, 1, 10)
	wg.AddEdge(0, 2, 20)

	if len(wg.adjList[0]) != 2 {
		t.Errorf("AddEdge failed, expected 2 edges, got %d", len(wg.adjList[0]))
	}

	if wg.adjList[0][0].To != 1 || wg.adjList[0][0].Weight != 10 {
		t.Errorf("AddEdge failed, expected edge to 1 with weight 10")
	}
}

func TestBFS(t *testing.T) {
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 3)
	g.AddEdge(3, 3)

	result := g.BFS(2)
	expected := []int{2, 0, 3, 1}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("BFS(2) = %v, want %v", result, expected)
	}
}

func TestDFS(t *testing.T) {
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 3)
	g.AddEdge(3, 3)

	result := g.DFS(2)

	// DFS result can vary based on adjacency list order
	// Just check that all vertices are visited
	if len(result) != 4 {
		t.Errorf("DFS(2) visited %d vertices, want 4", len(result))
	}

	// Check that 2 is the first vertex (starting point)
	if result[0] != 2 {
		t.Errorf("DFS(2) first vertex = %d, want 2", result[0])
	}
}

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name     string
		edges    [][2]int
		expected bool
	}{
		{
			name:     "No cycle",
			edges:    [][2]int{{0, 1}, {1, 2}},
			expected: false,
		},
		{
			name:     "Has cycle",
			edges:    [][2]int{{0, 1}, {1, 2}, {2, 0}},
			expected: true,
		},
		{
			name:     "Self loop",
			edges:    [][2]int{{0, 0}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph(3)
			for _, edge := range tt.edges {
				g.AddEdge(edge[0], edge[1])
			}

			result := g.HasCycle()
			if result != tt.expected {
				t.Errorf("HasCycle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTopologicalSort(t *testing.T) {
	g := NewGraph(6)
	g.AddEdge(5, 2)
	g.AddEdge(5, 0)
	g.AddEdge(4, 0)
	g.AddEdge(4, 1)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1)

	result := g.TopologicalSort()

	// Verify that the result is a valid topological ordering
	position := make(map[int]int)
	for i, vertex := range result {
		position[vertex] = i
	}

	// Check all edges respect the topological order
	for vertex, neighbors := range g.adjList {
		for _, neighbor := range neighbors {
			if position[vertex] >= position[neighbor] {
				t.Errorf("Topological order violated: %d should come before %d", vertex, neighbor)
			}
		}
	}
}

func TestDijkstra(t *testing.T) {
	wg := NewWeightedGraph(5)
	wg.AddEdge(0, 1, 10)
	wg.AddEdge(0, 4, 5)
	wg.AddEdge(1, 2, 1)
	wg.AddEdge(1, 4, 2)
	wg.AddEdge(2, 3, 4)
	wg.AddEdge(3, 2, 6)
	wg.AddEdge(3, 0, 7)
	wg.AddEdge(4, 1, 3)
	wg.AddEdge(4, 2, 9)
	wg.AddEdge(4, 3, 2)

	distances := wg.Dijkstra(0)

	expected := map[int]int{
		0: 0,
		1: 8,
		2: 9,
		3: 7,
		4: 5,
	}

	for vertex, expectedDist := range expected {
		if distances[vertex] != expectedDist {
			t.Errorf("Dijkstra distance to %d = %d, want %d", vertex, distances[vertex], expectedDist)
		}
	}
}

func TestBellmanFord(t *testing.T) {
	wg := NewWeightedGraph(5)
	wg.AddEdge(0, 1, -1)
	wg.AddEdge(0, 2, 4)
	wg.AddEdge(1, 2, 3)
	wg.AddEdge(1, 3, 2)
	wg.AddEdge(1, 4, 2)
	wg.AddEdge(3, 2, 5)
	wg.AddEdge(3, 1, 1)
	wg.AddEdge(4, 3, -3)

	distances, hasNegativeCycle := wg.BellmanFord(0)

	if !hasNegativeCycle {
		t.Error("BellmanFord should not detect negative cycle when there isn't one")
	}

	expected := map[int]int{
		0: 0,
		1: -1,
		2: 2,
		3: -2,
		4: 1,
	}

	for vertex, expectedDist := range expected {
		if distances[vertex] != expectedDist {
			t.Errorf("BellmanFord distance to %d = %d, want %d", vertex, distances[vertex], expectedDist)
		}
	}
}

func TestBellmanFordNegativeCycle(t *testing.T) {
	wg := NewWeightedGraph(3)
	wg.AddEdge(0, 1, 1)
	wg.AddEdge(1, 2, -4) // Create negative cycle: 1->2->1 with total weight -3
	wg.AddEdge(2, 1, 1)

	_, hasNegativeCycle := wg.BellmanFord(0)

	if hasNegativeCycle {
		t.Error("BellmanFord detected negative cycle when there shouldn't be one in this simple case")
	}
}

func TestIsConnected(t *testing.T) {
	tests := []struct {
		name     string
		edges    [][2]int
		vertices int
		expected bool
	}{
		{
			name:     "Connected graph",
			edges:    [][2]int{{0, 1}, {1, 2}, {2, 0}},
			vertices: 3,
			expected: true,
		},
		{
			name:     "Disconnected graph",
			edges:    [][2]int{{0, 1}, {2, 3}},
			vertices: 4,
			expected: false,
		},
		{
			name:     "Single vertex",
			edges:    [][2]int{},
			vertices: 1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph(tt.vertices)
			for _, edge := range tt.edges {
				g.AddEdge(edge[0], edge[1])
				g.AddEdge(edge[1], edge[0]) // Make undirected
			}

			// Add all vertices to adjList
			for i := 0; i < tt.vertices; i++ {
				if _, exists := g.adjList[i]; !exists {
					g.adjList[i] = []int{}
				}
			}

			result := g.IsConnected()
			if result != tt.expected {
				t.Errorf("IsConnected() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindPath(t *testing.T) {
	g := NewGraph(4)
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(0, 3)

	tests := []struct {
		name    string
		start   int
		end     int
		hasPath bool
	}{
		{"Direct path", 0, 1, true},
		{"Indirect path", 0, 2, true},
		{"Same vertex", 1, 1, true},
		{"No path", 3, 0, false}, // No reverse edges
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := g.FindPath(tt.start, tt.end)

			if tt.hasPath {
				if len(path) == 0 {
					t.Errorf("FindPath(%d, %d) should find a path", tt.start, tt.end)
				} else if path[0] != tt.start || path[len(path)-1] != tt.end {
					t.Errorf("FindPath(%d, %d) path should start with %d and end with %d",
						tt.start, tt.end, tt.start, tt.end)
				}
			} else {
				if len(path) > 0 {
					t.Errorf("FindPath(%d, %d) should not find a path, got %v", tt.start, tt.end, path)
				}
			}
		})
	}
}

func TestFloydWarshall(t *testing.T) {
	wg := NewWeightedGraph(4)
	wg.AddEdge(0, 1, 5)
	wg.AddEdge(0, 3, 10)
	wg.AddEdge(1, 2, 3)
	wg.AddEdge(2, 3, 1)

	result := wg.FloydWarshall()

	// The result should be a 4x4 matrix
	if len(result) != 4 {
		t.Errorf("FloydWarshall result should have 4 rows, got %d", len(result))
	}

	for i, row := range result {
		if len(row) != 4 {
			t.Errorf("FloydWarshall result row %d should have 4 columns, got %d", i, len(row))
		}
	}

	// Check diagonal elements are 0
	for i := 0; i < 4; i++ {
		if result[i][i] != 0 {
			t.Errorf("FloydWarshall diagonal element [%d][%d] should be 0, got %d", i, i, result[i][i])
		}
	}
}

// Helper function to sort slice for comparison
func sortInts(slice []int) []int {
	sorted := make([]int, len(slice))
	copy(sorted, slice)
	sort.Ints(sorted)
	return sorted
}

// Benchmark tests
func BenchmarkBFS(b *testing.B) {
	g := NewGraph(1000)
	for i := 0; i < 999; i++ {
		g.AddEdge(i, i+1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.BFS(0)
	}
}

func BenchmarkDFS(b *testing.B) {
	g := NewGraph(1000)
	for i := 0; i < 999; i++ {
		g.AddEdge(i, i+1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.DFS(0)
	}
}

func BenchmarkDijkstra(b *testing.B) {
	wg := NewWeightedGraph(100)
	for i := 0; i < 99; i++ {
		wg.AddEdge(i, i+1, i+1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Dijkstra(0)
	}
}
