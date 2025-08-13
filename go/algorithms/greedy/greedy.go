// Package greedy implements various greedy algorithms
package greedy

import (
	"container/heap"
	"errors"
	"sort"
)

// Activity represents an activity with start and finish times
type Activity struct {
	Start  int
	Finish int
	Index  int
}

// ActivitySelection solves the activity selection problem
// Time Complexity: O(n log n), Space Complexity: O(n)
func ActivitySelection(activities []Activity) []Activity {
	if len(activities) == 0 {
		return []Activity{}
	}

	// Sort activities by finish time
	sortedActivities := make([]Activity, len(activities))
	copy(sortedActivities, activities)
	sort.Slice(sortedActivities, func(i, j int) bool {
		return sortedActivities[i].Finish < sortedActivities[j].Finish
	})

	selected := []Activity{sortedActivities[0]}
	lastFinish := sortedActivities[0].Finish

	for i := 1; i < len(sortedActivities); i++ {
		if sortedActivities[i].Start >= lastFinish {
			selected = append(selected, sortedActivities[i])
			lastFinish = sortedActivities[i].Finish
		}
	}

	return selected
}

// HuffmanNode represents a node in the Huffman tree
type HuffmanNode struct {
	Char      rune
	Frequency int
	Left      *HuffmanNode
	Right     *HuffmanNode
}

// HuffmanHeap implements heap.Interface for HuffmanNode
type HuffmanHeap []*HuffmanNode

func (h HuffmanHeap) Len() int           { return len(h) }
func (h HuffmanHeap) Less(i, j int) bool { return h[i].Frequency < h[j].Frequency }
func (h HuffmanHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *HuffmanHeap) Push(x interface{}) {
	*h = append(*h, x.(*HuffmanNode))
}

func (h *HuffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// HuffmanCoding generates Huffman codes for given character frequencies
// Time Complexity: O(n log n), Space Complexity: O(n)
func HuffmanCoding(frequencies map[rune]int) (map[rune]string, *HuffmanNode) {
	if len(frequencies) == 0 {
		return map[rune]string{}, nil
	}

	if len(frequencies) == 1 {
		// Special case: only one character
		codes := make(map[rune]string)
		for char := range frequencies {
			codes[char] = "0"
		}
		return codes, &HuffmanNode{Char: 0, Frequency: 0}
	}

	// Create min heap
	h := &HuffmanHeap{}
	heap.Init(h)

	// Add all characters to heap
	for char, freq := range frequencies {
		heap.Push(h, &HuffmanNode{Char: char, Frequency: freq})
	}

	// Build Huffman tree
	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffmanNode)
		right := heap.Pop(h).(*HuffmanNode)

		merged := &HuffmanNode{
			Char:      0, // Internal node
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}

		heap.Push(h, merged)
	}

	root := heap.Pop(h).(*HuffmanNode)

	// Generate codes
	codes := make(map[rune]string)
	generateCodes(root, "", codes)

	return codes, root
}

func generateCodes(node *HuffmanNode, code string, codes map[rune]string) {
	if node == nil {
		return
	}

	// Leaf node
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = code
		return
	}

	generateCodes(node.Left, code+"0", codes)
	generateCodes(node.Right, code+"1", codes)
}

// Job represents a job with deadline and profit
type Job struct {
	ID       int
	Deadline int
	Profit   int
}

// JobScheduling solves the job scheduling problem to maximize profit
// Time Complexity: O(n²), Space Complexity: O(n)
func JobScheduling(jobs []Job) ([]Job, int) {
	if len(jobs) == 0 {
		return []Job{}, 0
	}

	// Sort jobs by profit in descending order
	sortedJobs := make([]Job, len(jobs))
	copy(sortedJobs, jobs)
	sort.Slice(sortedJobs, func(i, j int) bool {
		return sortedJobs[i].Profit > sortedJobs[j].Profit
	})

	// Find maximum deadline
	maxDeadline := 0
	for _, job := range sortedJobs {
		if job.Deadline > maxDeadline {
			maxDeadline = job.Deadline
		}
	}

	// Create time slots
	timeSlots := make([]bool, maxDeadline)
	selectedJobs := make([]Job, 0)
	totalProfit := 0

	// Schedule jobs
	for _, job := range sortedJobs {
		// Find a free slot for this job (from deadline-1 to 0)
		for slot := job.Deadline - 1; slot >= 0; slot-- {
			if !timeSlots[slot] {
				timeSlots[slot] = true
				selectedJobs = append(selectedJobs, job)
				totalProfit += job.Profit
				break
			}
		}
	}

	return selectedJobs, totalProfit
}

// Item represents an item for fractional knapsack
type Item struct {
	Value  int
	Weight int
	Index  int
}

// FractionalKnapsack solves the fractional knapsack problem
// Time Complexity: O(n log n), Space Complexity: O(n)
func FractionalKnapsack(items []Item, capacity int) (float64, []float64) {
	if len(items) == 0 || capacity <= 0 {
		return 0, []float64{}
	}

	// Create items with value-to-weight ratio
	type ItemRatio struct {
		Item  Item
		Ratio float64
	}

	itemRatios := make([]ItemRatio, len(items))
	for i, item := range items {
		if item.Weight == 0 {
			itemRatios[i] = ItemRatio{Item: item, Ratio: 0}
		} else {
			itemRatios[i] = ItemRatio{Item: item, Ratio: float64(item.Value) / float64(item.Weight)}
		}
	}

	// Sort by ratio in descending order
	sort.Slice(itemRatios, func(i, j int) bool {
		return itemRatios[i].Ratio > itemRatios[j].Ratio
	})

	totalValue := 0.0
	remainingCapacity := capacity
	fractions := make([]float64, len(items))

	for _, itemRatio := range itemRatios {
		item := itemRatio.Item
		if remainingCapacity >= item.Weight {
			// Take the whole item
			fractions[item.Index] = 1.0
			totalValue += float64(item.Value)
			remainingCapacity -= item.Weight
		} else if remainingCapacity > 0 {
			// Take fraction of the item
			fraction := float64(remainingCapacity) / float64(item.Weight)
			fractions[item.Index] = fraction
			totalValue += fraction * float64(item.Value)
			remainingCapacity = 0
			break
		}
	}

	return totalValue, fractions
}

// MinimumCoins finds minimum number of coins needed to make change
// Time Complexity: O(n), Space Complexity: O(1) for greedy approach
func MinimumCoins(coins []int, amount int) (int, []int, error) {
	if amount < 0 {
		return 0, []int{}, errors.New("amount cannot be negative")
	}
	if amount == 0 {
		return 0, []int{}, nil
	}

	// Sort coins in descending order
	sortedCoins := make([]int, len(coins))
	copy(sortedCoins, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedCoins)))

	result := make([]int, 0)
	remaining := amount

	for _, coin := range sortedCoins {
		count := remaining / coin
		for i := 0; i < count; i++ {
			result = append(result, coin)
		}
		remaining %= coin
	}

	if remaining > 0 {
		return 0, []int{}, errors.New("cannot make exact change with given coins")
	}

	return len(result), result, nil
}

// GasStation solves the gas station problem
// Time Complexity: O(n), Space Complexity: O(1)
func GasStation(gas []int, cost []int) int {
	if len(gas) != len(cost) {
		return -1
	}

	totalGas := 0
	totalCost := 0
	currentGas := 0
	start := 0

	for i := 0; i < len(gas); i++ {
		totalGas += gas[i]
		totalCost += cost[i]
		currentGas += gas[i] - cost[i]

		// If current gas becomes negative, we can't reach next station
		if currentGas < 0 {
			start = i + 1
			currentGas = 0
		}
	}

	// If total gas is less than total cost, impossible to complete circuit
	if totalGas < totalCost {
		return -1
	}

	return start
}

// IntervalScheduling represents an interval with start and end times
type Interval struct {
	Start int
	End   int
}

// IntervalScheduling finds maximum number of non-overlapping intervals
// Time Complexity: O(n log n), Space Complexity: O(n)
func IntervalScheduling(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return []Interval{}
	}

	// Sort intervals by end time
	sortedIntervals := make([]Interval, len(intervals))
	copy(sortedIntervals, intervals)
	sort.Slice(sortedIntervals, func(i, j int) bool {
		return sortedIntervals[i].End < sortedIntervals[j].End
	})

	selected := []Interval{sortedIntervals[0]}
	lastEnd := sortedIntervals[0].End

	for i := 1; i < len(sortedIntervals); i++ {
		if sortedIntervals[i].Start >= lastEnd {
			selected = append(selected, sortedIntervals[i])
			lastEnd = sortedIntervals[i].End
		}
	}

	return selected
}

// MinimumPlatforms finds minimum number of platforms needed for train schedule
// Time Complexity: O(n log n), Space Complexity: O(n)
func MinimumPlatforms(arrivals []int, departures []int) int {
	if len(arrivals) != len(departures) {
		return 0
	}

	n := len(arrivals)
	if n == 0 {
		return 0
	}

	// Sort arrival and departure times
	sortedArrivals := make([]int, n)
	sortedDepartures := make([]int, n)
	copy(sortedArrivals, arrivals)
	copy(sortedDepartures, departures)
	sort.Ints(sortedArrivals)
	sort.Ints(sortedDepartures)

	platformsNeeded := 1
	maxPlatforms := 1
	i, j := 1, 0

	// Use two pointers technique
	for i < n && j < n {
		if sortedArrivals[i] <= sortedDepartures[j] {
			platformsNeeded++
			i++
		} else {
			platformsNeeded--
			j++
		}

		if platformsNeeded > maxPlatforms {
			maxPlatforms = platformsNeeded
		}
	}

	return maxPlatforms
}
