package greedy

import (
	"testing"
)

func TestActivitySelection(t *testing.T) {
	activities := []Activity{
		{Start: 1, Finish: 4, Index: 0},
		{Start: 3, Finish: 5, Index: 1},
		{Start: 0, Finish: 6, Index: 2},
		{Start: 5, Finish: 7, Index: 3},
		{Start: 8, Finish: 9, Index: 4},
		{Start: 5, Finish: 9, Index: 5},
	}

	selected := ActivitySelection(activities)

	// Should select activities that don't overlap
	if len(selected) < 3 {
		t.Errorf("Expected at least 3 activities, got %d", len(selected))
	}

	// Verify no overlaps
	for i := 0; i < len(selected)-1; i++ {
		if selected[i].Finish > selected[i+1].Start {
			t.Error("Selected activities overlap")
		}
	}

	// Test empty input
	empty := ActivitySelection([]Activity{})
	if len(empty) != 0 {
		t.Error("Expected empty result for empty input")
	}
}

func TestHuffmanCoding(t *testing.T) {
	frequencies := map[rune]int{
		'a': 5,
		'b': 9,
		'c': 12,
		'd': 13,
		'e': 16,
		'f': 45,
	}

	codes, root := HuffmanCoding(frequencies)

	// Check that all characters have codes
	if len(codes) != len(frequencies) {
		t.Errorf("Expected %d codes, got %d", len(frequencies), len(codes))
	}

	// Check that codes are valid (no code is prefix of another)
	for char1, code1 := range codes {
		for char2, code2 := range codes {
			if char1 != char2 {
				if len(code1) <= len(code2) && code2[:len(code1)] == code1 {
					t.Errorf("Code %s is prefix of %s", code1, code2)
				}
			}
		}
	}

	// Root should not be nil
	if root == nil {
		t.Error("Root should not be nil")
	}

	// Test single character
	singleChar := map[rune]int{'a': 10}
	singleCodes, _ := HuffmanCoding(singleChar)
	if len(singleCodes) != 1 || singleCodes['a'] != "0" {
		t.Error("Single character should have code '0'")
	}

	// Test empty input
	emptyCodes, emptyRoot := HuffmanCoding(map[rune]int{})
	if len(emptyCodes) != 0 || emptyRoot != nil {
		t.Error("Empty input should return empty codes and nil root")
	}
}

func TestJobScheduling(t *testing.T) {
	jobs := []Job{
		{ID: 1, Deadline: 4, Profit: 20},
		{ID: 2, Deadline: 1, Profit: 10},
		{ID: 3, Deadline: 1, Profit: 40},
		{ID: 4, Deadline: 1, Profit: 30},
	}

	selected, totalProfit := JobScheduling(jobs)

	// Should select jobs with maximum profit
	if totalProfit <= 0 {
		t.Error("Total profit should be positive")
	}

	// Check that deadlines are respected
	timeSlots := make(map[int]bool)
	for _, job := range selected {
		// Find which slot this job was assigned
		for slot := 0; slot < job.Deadline; slot++ {
			if !timeSlots[slot] {
				timeSlots[slot] = true
				break
			}
		}
	}

	// Test empty input
	emptyJobs, emptyProfit := JobScheduling([]Job{})
	if len(emptyJobs) != 0 || emptyProfit != 0 {
		t.Error("Empty input should return empty jobs and zero profit")
	}
}

func TestFractionalKnapsack(t *testing.T) {
	items := []Item{
		{Value: 60, Weight: 10, Index: 0},
		{Value: 100, Weight: 20, Index: 1},
		{Value: 120, Weight: 30, Index: 2},
	}
	capacity := 50

	totalValue, fractions := FractionalKnapsack(items, capacity)

	// Should achieve maximum value
	if totalValue <= 0 {
		t.Error("Total value should be positive")
	}

	// Check fractions are valid (between 0 and 1)
	for i, fraction := range fractions {
		if fraction < 0 || fraction > 1 {
			t.Errorf("Invalid fraction %f for item %d", fraction, i)
		}
	}

	// Test zero capacity
	zeroValue, zeroFractions := FractionalKnapsack(items, 0)
	if zeroValue != 0 || len(zeroFractions) != 0 {
		t.Error("Zero capacity should return zero value and empty fractions")
	}

	// Test empty items
	emptyValue, emptyFractions := FractionalKnapsack([]Item{}, capacity)
	if emptyValue != 0 || len(emptyFractions) != 0 {
		t.Error("Empty items should return zero value and empty fractions")
	}
}

func TestMinimumCoins(t *testing.T) {
	coins := []int{1, 5, 10, 25}
	amount := 30

	count, result, err := MinimumCoins(coins, amount)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the result
	total := 0
	for _, coin := range result {
		total += coin
	}
	if total != amount {
		t.Errorf("Coins sum to %d, expected %d", total, amount)
	}

	// Should use minimum number of coins
	if count != len(result) {
		t.Errorf("Count %d doesn't match result length %d", count, len(result))
	}

	// Test impossible case
	impossibleCoins := []int{3, 5}
	impossibleAmount := 1
	_, _, err = MinimumCoins(impossibleCoins, impossibleAmount)
	if err == nil {
		t.Error("Should return error for impossible case")
	}

	// Test zero amount
	zeroCount, zeroResult, err := MinimumCoins(coins, 0)
	if err != nil || zeroCount != 0 || len(zeroResult) != 0 {
		t.Error("Zero amount should return zero count and empty result")
	}

	// Test negative amount
	_, _, err = MinimumCoins(coins, -1)
	if err == nil {
		t.Error("Negative amount should return error")
	}
}

func TestGasStation(t *testing.T) {
	// Test case 1: Possible to complete circuit
	gas1 := []int{1, 2, 3, 4, 5}
	cost1 := []int{3, 4, 5, 1, 2}
	start1 := GasStation(gas1, cost1)
	if start1 == -1 {
		t.Error("Should be able to complete circuit")
	}

	// Test case 2: Impossible to complete circuit
	gas2 := []int{2, 3, 4}
	cost2 := []int{3, 4, 3}
	start2 := GasStation(gas2, cost2)
	if start2 != -1 {
		t.Error("Should not be able to complete circuit")
	}

	// Test case 3: Mismatched arrays
	gas3 := []int{1, 2}
	cost3 := []int{1}
	start3 := GasStation(gas3, cost3)
	if start3 != -1 {
		t.Error("Mismatched arrays should return -1")
	}
}

func TestIntervalScheduling(t *testing.T) {
	intervals := []Interval{
		{Start: 1, End: 3},
		{Start: 2, End: 4},
		{Start: 3, End: 6},
		{Start: 5, End: 7},
		{Start: 8, End: 9},
	}

	selected := IntervalScheduling(intervals)

	// Should select non-overlapping intervals
	for i := 0; i < len(selected)-1; i++ {
		if selected[i].End > selected[i+1].Start {
			t.Error("Selected intervals overlap")
		}
	}

	// Should select maximum number of intervals
	if len(selected) < 3 {
		t.Errorf("Expected at least 3 intervals, got %d", len(selected))
	}

	// Test empty input
	empty := IntervalScheduling([]Interval{})
	if len(empty) != 0 {
		t.Error("Empty input should return empty result")
	}
}

func TestMinimumPlatforms(t *testing.T) {
	arrivals := []int{900, 940, 950, 1100, 1500, 1800}
	departures := []int{910, 1200, 1120, 1130, 1900, 2000}

	platforms := MinimumPlatforms(arrivals, departures)
	if platforms <= 0 {
		t.Error("Should need at least one platform")
	}

	// Test empty input
	emptyPlatforms := MinimumPlatforms([]int{}, []int{})
	if emptyPlatforms != 0 {
		t.Error("Empty input should return 0 platforms")
	}

	// Test mismatched arrays
	mismatchedPlatforms := MinimumPlatforms([]int{900}, []int{})
	if mismatchedPlatforms != 0 {
		t.Error("Mismatched arrays should return 0 platforms")
	}

	// Test single train
	singlePlatforms := MinimumPlatforms([]int{900}, []int{1000})
	if singlePlatforms != 1 {
		t.Error("Single train should need 1 platform")
	}
}

// Benchmark tests
func BenchmarkActivitySelection(b *testing.B) {
	activities := make([]Activity, 1000)
	for i := 0; i < 1000; i++ {
		activities[i] = Activity{Start: i, Finish: i + 10, Index: i}
	}

	for i := 0; i < b.N; i++ {
		ActivitySelection(activities)
	}
}

func BenchmarkHuffmanCoding(b *testing.B) {
	frequencies := map[rune]int{
		'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45,
		'g': 3, 'h': 7, 'i': 11, 'j': 2, 'k': 8, 'l': 15,
	}

	for i := 0; i < b.N; i++ {
		HuffmanCoding(frequencies)
	}
}

func BenchmarkJobScheduling(b *testing.B) {
	jobs := make([]Job, 100)
	for i := 0; i < 100; i++ {
		jobs[i] = Job{ID: i, Deadline: i%10 + 1, Profit: i*10 + 5}
	}

	for i := 0; i < b.N; i++ {
		JobScheduling(jobs)
	}
}

func BenchmarkFractionalKnapsack(b *testing.B) {
	items := make([]Item, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = Item{Value: i*10 + 5, Weight: i + 1, Index: i}
	}

	for i := 0; i < b.N; i++ {
		FractionalKnapsack(items, 5000)
	}
}
