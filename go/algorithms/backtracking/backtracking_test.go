package backtracking

import (
	"reflect"
	"testing"
)

func TestNQueens(t *testing.T) {
	tests := []struct {
		n        int
		expected int // number of solutions
	}{
		{0, 0},
		{1, 1},
		{2, 0},
		{3, 0},
		{4, 2},
		{8, 92},
	}

	for _, test := range tests {
		solutions := NQueens(test.n)
		if len(solutions) != test.expected {
			t.Errorf("NQueens(%d) = %d solutions, expected %d", test.n, len(solutions), test.expected)
		}
	}
}

func TestSudokuSolver(t *testing.T) {
	// Test case 1: Solvable Sudoku
	board1 := [][]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	expected1 := [][]int{
		{5, 3, 4, 6, 7, 8, 9, 1, 2},
		{6, 7, 2, 1, 9, 5, 3, 4, 8},
		{1, 9, 8, 3, 4, 2, 5, 6, 7},
		{8, 5, 9, 7, 6, 1, 4, 2, 3},
		{4, 2, 6, 8, 5, 3, 7, 9, 1},
		{7, 1, 3, 9, 2, 4, 8, 5, 6},
		{9, 6, 1, 5, 3, 7, 2, 8, 4},
		{2, 8, 7, 4, 1, 9, 6, 3, 5},
		{3, 4, 5, 2, 8, 6, 1, 7, 9},
	}

	if !SudokuSolver(board1) {
		t.Error("SudokuSolver failed to solve a solvable puzzle")
	}

	if !reflect.DeepEqual(board1, expected1) {
		t.Error("SudokuSolver did not produce the expected solution")
	}

	// Test case 2: Unsolvable Sudoku
	board2 := [][]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 8}, // Invalid: two 8s in last row
	}

	if SudokuSolver(board2) {
		t.Error("SudokuSolver should have failed on unsolvable puzzle")
	}
}

func TestKnightsTour(t *testing.T) {
	tests := []struct {
		n           int
		shouldSolve bool
	}{
		{4, false}, // Too small
		{5, true},  // Should have solution
		{6, true},  // Should have solution
		{8, true},  // Should have solution
	}

	for _, test := range tests {
		solution, err := KnightsTour(test.n)
		if test.shouldSolve {
			if err != nil {
				t.Errorf("KnightsTour(%d) should have found a solution, got error: %v", test.n, err)
			}
			if solution == nil {
				t.Errorf("KnightsTour(%d) returned nil solution", test.n)
			}
			// Verify solution has correct number of moves
			if solution != nil && len(solution) != test.n {
				t.Errorf("KnightsTour(%d) solution has wrong dimensions", test.n)
			}
		} else {
			if err == nil {
				t.Errorf("KnightsTour(%d) should have failed", test.n)
			}
		}
	}
}

func TestGraphColoring(t *testing.T) {
	// Test case 1: Simple triangle graph (3 vertices, all connected)
	graph1 := [][]int{
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}

	colors1, solved1 := GraphColoring(graph1, 3)
	if !solved1 {
		t.Error("GraphColoring should solve triangle graph with 3 colors")
	}
	if len(colors1) != 3 {
		t.Error("GraphColoring should return 3 colors for 3 vertices")
	}

	// Verify no adjacent vertices have same color
	for i := 0; i < len(graph1); i++ {
		for j := 0; j < len(graph1[i]); j++ {
			if graph1[i][j] == 1 && colors1[i] == colors1[j] {
				t.Error("Adjacent vertices have same color")
			}
		}
	}

	// Test case 2: Triangle graph with only 2 colors (should fail)
	colors2, solved2 := GraphColoring(graph1, 2)
	if solved2 {
		t.Error("GraphColoring should not solve triangle graph with only 2 colors")
	}
	if len(colors2) != 0 {
		t.Error("GraphColoring should return empty slice when no solution exists")
	}
}

func TestHamiltonianPath(t *testing.T) {
	// Test case 1: Simple path graph
	graph1 := [][]int{
		{0, 1, 0, 0},
		{1, 0, 1, 0},
		{0, 1, 0, 1},
		{0, 0, 1, 0},
	}

	path1, found1 := HamiltonianPath(graph1)
	if !found1 {
		t.Error("HamiltonianPath should find path in simple path graph")
	}
	if len(path1) != 4 {
		t.Error("HamiltonianPath should return path of length 4")
	}

	// Test case 2: Disconnected graph (should fail)
	graph2 := [][]int{
		{0, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 0, 1},
		{0, 0, 1, 0},
	}

	path2, found2 := HamiltonianPath(graph2)
	if found2 {
		t.Error("HamiltonianPath should not find path in disconnected graph")
	}
	if len(path2) != 0 {
		t.Error("HamiltonianPath should return empty path when no solution exists")
	}
}

func TestGenerateSubsets(t *testing.T) {
	tests := []struct {
		input    []int
		expected int // number of subsets
	}{
		{[]int{}, 1},       // Empty set has 1 subset (empty subset)
		{[]int{1}, 2},      // {}, {1}
		{[]int{1, 2}, 4},   // {}, {1}, {2}, {1,2}
		{[]int{1, 2, 3}, 8}, // 2^3 = 8 subsets
	}

	for _, test := range tests {
		subsets := GenerateSubsets(test.input)
		if len(subsets) != test.expected {
			t.Errorf("GenerateSubsets(%v) = %d subsets, expected %d", test.input, len(subsets), test.expected)
		}
	}

	// Test specific case
	subsets := GenerateSubsets([]int{1, 2})
	expectedSubsets := [][]int{
		{},
		{1},
		{1, 2},
		{2},
	}

	if len(subsets) != len(expectedSubsets) {
		t.Errorf("GenerateSubsets([1,2]) returned %d subsets, expected %d", len(subsets), len(expectedSubsets))
	}

	// Check if all expected subsets are present (order may vary)
	found := make(map[string]bool)
	for _, subset := range subsets {
		key := ""
		for _, num := range subset {
			key += string(rune(num + '0'))
		}
		found[key] = true
	}

	expectedKeys := []string{"", "1", "12", "2"}
	for _, key := range expectedKeys {
		if !found[key] {
			t.Errorf("Expected subset %s not found", key)
		}
	}
}

// Benchmark tests
func BenchmarkNQueens4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NQueens(4)
	}
}

func BenchmarkNQueens8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NQueens(8)
	}
}

func BenchmarkSudokuSolver(b *testing.B) {
	board := [][]int{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	for i := 0; i < b.N; i++ {
		// Make a copy for each iteration
		testBoard := make([][]int, 9)
		for j := range testBoard {
			testBoard[j] = make([]int, 9)
			copy(testBoard[j], board[j])
		}
		SudokuSolver(testBoard)
	}
}

func BenchmarkGenerateSubsets(b *testing.B) {
	nums := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		GenerateSubsets(nums)
	}
}
