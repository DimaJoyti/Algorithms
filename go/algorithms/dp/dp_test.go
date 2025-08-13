package dp

import (
	"testing"
)

func TestLongestIncreasingSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"Empty array", []int{}, 0},
		{"Single element", []int{5}, 1},
		{"Increasing sequence", []int{1, 2, 3, 4, 5}, 5},
		{"Decreasing sequence", []int{5, 4, 3, 2, 1}, 1},
		{"Mixed sequence", []int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{"With duplicates", []int{1, 3, 6, 7, 9, 4, 10, 5, 6}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LongestIncreasingSubsequence(tt.input)
			if result != tt.expected {
				t.Errorf("LongestIncreasingSubsequence(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLISOptimized(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"Empty array", []int{}, 0},
		{"Single element", []int{5}, 1},
		{"Increasing sequence", []int{1, 2, 3, 4, 5}, 5},
		{"Mixed sequence", []int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LISOptimized(tt.input)
			if result != tt.expected {
				t.Errorf("LISOptimized(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestKnapsack01(t *testing.T) {
	tests := []struct {
		name        string
		items       []KnapsackItem
		capacity    int
		expectedVal int
		expectedLen int
	}{
		{
			"Empty items",
			[]KnapsackItem{},
			10,
			0,
			0,
		},
		{
			"Zero capacity",
			[]KnapsackItem{{Weight: 1, Value: 1}},
			0,
			0,
			0,
		},
		{
			"Classic example",
			[]KnapsackItem{
				{Weight: 10, Value: 60},
				{Weight: 20, Value: 100},
				{Weight: 30, Value: 120},
			},
			50,
			220,
			2,
		},
		{
			"Single item fits",
			[]KnapsackItem{{Weight: 5, Value: 10}},
			10,
			10,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, items := Knapsack01(tt.items, tt.capacity)
			if value != tt.expectedVal {
				t.Errorf("Knapsack01() value = %d, want %d", value, tt.expectedVal)
			}
			if len(items) != tt.expectedLen {
				t.Errorf("Knapsack01() items length = %d, want %d", len(items), tt.expectedLen)
			}
		})
	}
}

func TestCoinChange(t *testing.T) {
	tests := []struct {
		name     string
		coins    []int
		amount   int
		expected int
	}{
		{"Zero amount", []int{1, 3, 4}, 0, 0},
		{"No coins", []int{}, 5, -1},
		{"Impossible", []int{2}, 3, -1},
		{"Classic example", []int{1, 3, 4}, 6, 2},
		{"Single coin", []int{5}, 10, 2},
		{"Multiple solutions", []int{1, 2, 5}, 11, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CoinChange(tt.coins, tt.amount)
			if result != tt.expected {
				t.Errorf("CoinChange(%v, %d) = %d, want %d", tt.coins, tt.amount, result, tt.expected)
			}
		})
	}
}

func TestCoinChangeWays(t *testing.T) {
	tests := []struct {
		name     string
		coins    []int
		amount   int
		expected int
	}{
		{"Zero amount", []int{1, 2, 3}, 0, 1},
		{"No coins", []int{}, 5, 0},
		{"Single way", []int{2}, 3, 0},
		{"Multiple ways", []int{1, 2, 3}, 4, 4},
		{"Classic example", []int{2, 3, 5}, 8, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CoinChangeWays(tt.coins, tt.amount)
			if result != tt.expected {
				t.Errorf("CoinChangeWays(%v, %d) = %d, want %d", tt.coins, tt.amount, result, tt.expected)
			}
		})
	}
}

func TestMaxSubarraySum(t *testing.T) {
	tests := []struct {
		name        string
		input       []int
		expected    int
		shouldError bool
	}{
		{"Empty array", []int{}, 0, true},
		{"Single positive", []int{5}, 5, false},
		{"Single negative", []int{-5}, -5, false},
		{"All negative", []int{-2, -3, -1, -5}, -1, false},
		{"Mixed array", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6, false},
		{"All positive", []int{1, 2, 3, 4, 5}, 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MaxSubarraySum(tt.input)
			if tt.shouldError && err == nil {
				t.Errorf("MaxSubarraySum(%v) should return error", tt.input)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("MaxSubarraySum(%v) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldError && result != tt.expected {
				t.Errorf("MaxSubarraySum(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaxSubarrayWithIndices(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		expectedSum   int
		expectedStart int
		expectedEnd   int
		shouldError   bool
	}{
		{"Empty array", []int{}, 0, 0, 0, true},
		{"Single element", []int{5}, 5, 0, 0, false},
		{"Mixed array", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6, 3, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, start, end, err := MaxSubarrayWithIndices(tt.input)
			if tt.shouldError && err == nil {
				t.Errorf("MaxSubarrayWithIndices(%v) should return error", tt.input)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("MaxSubarrayWithIndices(%v) unexpected error: %v", tt.input, err)
			}
			if !tt.shouldError {
				if sum != tt.expectedSum {
					t.Errorf("MaxSubarrayWithIndices(%v) sum = %d, want %d", tt.input, sum, tt.expectedSum)
				}
				if start != tt.expectedStart {
					t.Errorf("MaxSubarrayWithIndices(%v) start = %d, want %d", tt.input, start, tt.expectedStart)
				}
				if end != tt.expectedEnd {
					t.Errorf("MaxSubarrayWithIndices(%v) end = %d, want %d", tt.input, end, tt.expectedEnd)
				}
			}
		})
	}
}

func TestHouseRobber(t *testing.T) {
	tests := []struct {
		name     string
		houses   []int
		expected int
	}{
		{"Empty houses", []int{}, 0},
		{"Single house", []int{5}, 5},
		{"Two houses", []int{2, 7}, 7},
		{"Three houses", []int{2, 1, 1}, 3},
		{"Classic example", []int{2, 7, 9, 3, 1}, 12},
		{"Another example", []int{1, 2, 3, 1}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HouseRobber(tt.houses)
			if result != tt.expected {
				t.Errorf("HouseRobber(%v) = %d, want %d", tt.houses, result, tt.expected)
			}
		})
	}
}

func TestHouseRobberCircular(t *testing.T) {
	tests := []struct {
		name     string
		houses   []int
		expected int
	}{
		{"Empty houses", []int{}, 0},
		{"Single house", []int{5}, 5},
		{"Two houses", []int{2, 3}, 3},
		{"Three houses", []int{2, 3, 2}, 3},
		{"Four houses", []int{1, 2, 3, 1}, 4},
		{"Five houses", []int{2, 7, 9, 3, 1}, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HouseRobberCircular(tt.houses)
			if result != tt.expected {
				t.Errorf("HouseRobberCircular(%v) = %d, want %d", tt.houses, result, tt.expected)
			}
		})
	}
}

func TestEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		str1     string
		str2     string
		expected int
	}{
		{"Both empty", "", "", 0},
		{"One empty", "abc", "", 3},
		{"Other empty", "", "abc", 3},
		{"Same strings", "abc", "abc", 0},
		{"Classic example", "kitten", "sitting", 3},
		{"Another example", "sunday", "saturday", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EditDistance(tt.str1, tt.str2)
			if result != tt.expected {
				t.Errorf("EditDistance(%s, %s) = %d, want %d", tt.str1, tt.str2, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkLongestIncreasingSubsequence(b *testing.B) {
	arr := []int{10, 9, 2, 5, 3, 7, 101, 18, 1, 4, 6, 8, 11, 12, 13, 14, 15}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LongestIncreasingSubsequence(arr)
	}
}

func BenchmarkLISOptimized(b *testing.B) {
	arr := []int{10, 9, 2, 5, 3, 7, 101, 18, 1, 4, 6, 8, 11, 12, 13, 14, 15}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LISOptimized(arr)
	}
}

func BenchmarkKnapsack01(b *testing.B) {
	items := []KnapsackItem{
		{Weight: 10, Value: 60},
		{Weight: 20, Value: 100},
		{Weight: 30, Value: 120},
		{Weight: 15, Value: 80},
		{Weight: 25, Value: 90},
	}
	capacity := 50
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Knapsack01(items, capacity)
	}
}

func BenchmarkMaxSubarraySum(b *testing.B) {
	arr := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4, 3, -1, 2, -4, 5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MaxSubarraySum(arr)
	}
}
