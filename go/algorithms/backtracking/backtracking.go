// Package backtracking implements various backtracking algorithms
package backtracking

import (
	"errors"
	"fmt"
)

// NQueens solves the N-Queens problem and returns all solutions
// Time Complexity: O(N!), Space Complexity: O(N²)
func NQueens(n int) [][][]string {
	if n <= 0 {
		return [][][]string{}
	}

	solutions := make([][][]string, 0)
	board := make([][]string, n)
	for i := range board {
		board[i] = make([]string, n)
		for j := range board[i] {
			board[i][j] = "."
		}
	}

	solveNQueens(board, 0, &solutions)
	return solutions
}

func solveNQueens(board [][]string, row int, solutions *[][][]string) {
	n := len(board)
	if row == n {
		// Found a solution, make a copy
		solution := make([][]string, n)
		for i := range board {
			solution[i] = make([]string, n)
			copy(solution[i], board[i])
		}
		*solutions = append(*solutions, solution)
		return
	}

	for col := 0; col < n; col++ {
		if isSafeQueen(board, row, col) {
			board[row][col] = "Q"
			solveNQueens(board, row+1, solutions)
			board[row][col] = "." // Backtrack
		}
	}
}

func isSafeQueen(board [][]string, row, col int) bool {
	n := len(board)

	// Check column
	for i := 0; i < row; i++ {
		if board[i][col] == "Q" {
			return false
		}
	}

	// Check diagonal (top-left to bottom-right)
	for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if board[i][j] == "Q" {
			return false
		}
	}

	// Check diagonal (top-right to bottom-left)
	for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
		if board[i][j] == "Q" {
			return false
		}
	}

	return true
}

// SudokuSolver solves a 9x9 Sudoku puzzle
// Time Complexity: O(9^(n*n)), Space Complexity: O(n*n)
func SudokuSolver(board [][]int) bool {
	return solveSudoku(board)
}

func solveSudoku(board [][]int) bool {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if board[row][col] == 0 {
				for num := 1; num <= 9; num++ {
					if isSafeSudoku(board, row, col, num) {
						board[row][col] = num

						if solveSudoku(board) {
							return true
						}

						board[row][col] = 0 // Backtrack
					}
				}
				return false
			}
		}
	}
	return true
}

func isSafeSudoku(board [][]int, row, col, num int) bool {
	// Check row
	for j := 0; j < 9; j++ {
		if board[row][j] == num {
			return false
		}
	}

	// Check column
	for i := 0; i < 9; i++ {
		if board[i][col] == num {
			return false
		}
	}

	// Check 3x3 box
	startRow := row - row%3
	startCol := col - col%3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i+startRow][j+startCol] == num {
				return false
			}
		}
	}

	return true
}

// KnightsTour finds a knight's tour on an n×n chessboard
// Time Complexity: O(8^(n²)), Space Complexity: O(n²)
func KnightsTour(n int) ([][]int, error) {
	if n < 5 {
		return nil, errors.New("knight's tour not possible for board size less than 5")
	}

	board := make([][]int, n)
	for i := range board {
		board[i] = make([]int, n)
		for j := range board[i] {
			board[i][j] = -1
		}
	}

	// Knight moves
	moveX := []int{2, 1, -1, -2, -2, -1, 1, 2}
	moveY := []int{1, 2, 2, 1, -1, -2, -2, -1}

	// Start from (0, 0)
	board[0][0] = 0

	if solveKnightsTour(board, 0, 0, 1, moveX, moveY, n) {
		return board, nil
	}

	return nil, errors.New("knight's tour solution does not exist")
}

func solveKnightsTour(board [][]int, x, y, moveCount int, moveX, moveY []int, n int) bool {
	if moveCount == n*n {
		return true
	}

	for i := 0; i < 8; i++ {
		nextX := x + moveX[i]
		nextY := y + moveY[i]

		if isSafeKnight(board, nextX, nextY, n) {
			board[nextX][nextY] = moveCount
			if solveKnightsTour(board, nextX, nextY, moveCount+1, moveX, moveY, n) {
				return true
			}
			board[nextX][nextY] = -1 // Backtrack
		}
	}

	return false
}

func isSafeKnight(board [][]int, x, y, n int) bool {
	return x >= 0 && x < n && y >= 0 && y < n && board[x][y] == -1
}

// GraphColoring solves the graph coloring problem
// Time Complexity: O(m^V), Space Complexity: O(V)
func GraphColoring(graph [][]int, numColors int) ([]int, bool) {
	vertices := len(graph)
	if vertices == 0 {
		return []int{}, true
	}

	colors := make([]int, vertices)
	for i := range colors {
		colors[i] = -1
	}

	if solveGraphColoring(graph, numColors, colors, 0) {
		return colors, true
	}

	return []int{}, false
}

func solveGraphColoring(graph [][]int, numColors int, colors []int, vertex int) bool {
	vertices := len(graph)
	if vertex == vertices {
		return true
	}

	for color := 0; color < numColors; color++ {
		if isSafeColor(graph, colors, vertex, color) {
			colors[vertex] = color

			if solveGraphColoring(graph, numColors, colors, vertex+1) {
				return true
			}

			colors[vertex] = -1 // Backtrack
		}
	}

	return false
}

func isSafeColor(graph [][]int, colors []int, vertex, color int) bool {
	for i := 0; i < len(graph); i++ {
		if graph[vertex][i] == 1 && colors[i] == color {
			return false
		}
	}
	return true
}

// HamiltonianPath finds a Hamiltonian path in the graph
// Time Complexity: O(N!), Space Complexity: O(N)
func HamiltonianPath(graph [][]int) ([]int, bool) {
	vertices := len(graph)
	if vertices == 0 {
		return []int{}, true
	}

	path := make([]int, vertices)
	visited := make([]bool, vertices)

	// Try starting from each vertex
	for start := 0; start < vertices; start++ {
		// Reset for each starting vertex
		for i := range visited {
			visited[i] = false
		}
		for i := range path {
			path[i] = -1
		}

		path[0] = start
		visited[start] = true

		if solveHamiltonianPath(graph, path, visited, 1) {
			return path, true
		}
	}

	return []int{}, false
}

func solveHamiltonianPath(graph [][]int, path []int, visited []bool, pos int) bool {
	vertices := len(graph)
	if pos == vertices {
		return true
	}

	for v := 0; v < vertices; v++ {
		if isSafeHamiltonian(graph, path, visited, pos, v) {
			path[pos] = v
			visited[v] = true

			if solveHamiltonianPath(graph, path, visited, pos+1) {
				return true
			}

			visited[v] = false // Backtrack
		}
	}

	return false
}

func isSafeHamiltonian(graph [][]int, path []int, visited []bool, pos, v int) bool {
	// Check if vertex v is adjacent to the previously added vertex
	if graph[path[pos-1]][v] == 0 {
		return false
	}

	// Check if vertex v is already included in the path
	if visited[v] {
		return false
	}

	return true
}

// GenerateSubsets generates all subsets of a given set
// Time Complexity: O(2^n), Space Complexity: O(2^n)
func GenerateSubsets(nums []int) [][]int {
	result := make([][]int, 0)
	current := make([]int, 0)
	generateSubsetsHelper(nums, 0, current, &result)
	return result
}

func generateSubsetsHelper(nums []int, index int, current []int, result *[][]int) {
	// Add current subset to result
	subset := make([]int, len(current))
	copy(subset, current)
	*result = append(*result, subset)

	// Generate subsets by including elements from index onwards
	for i := index; i < len(nums); i++ {
		current = append(current, nums[i])
		generateSubsetsHelper(nums, i+1, current, result)
		current = current[:len(current)-1] // Backtrack
	}
}

// PrintNQueensSolution prints a single N-Queens solution in a readable format
func PrintNQueensSolution(solution [][]string) {
	for _, row := range solution {
		fmt.Println(row)
	}
	fmt.Println()
}

// PrintSudoku prints a Sudoku board in a readable format
func PrintSudoku(board [][]int) {
	for i, row := range board {
		if i%3 == 0 && i != 0 {
			fmt.Println("------+-------+------")
		}
		for j, val := range row {
			if j%3 == 0 && j != 0 {
				fmt.Print("| ")
			}
			if val == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", val)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
