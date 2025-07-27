// Package dp implements various dynamic programming algorithms
package dp

import (
	"errors"
	"math"
)

// LongestIncreasingSubsequence finds the length of the longest increasing subsequence
// Time Complexity: O(n²), Space Complexity: O(n)
func LongestIncreasingSubsequence(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	n := len(arr)
	dp := make([]int, n)
	
	// Initialize all lengths as 1
	for i := range dp {
		dp[i] = 1
	}

	// Compute optimized LIS values in bottom-up manner
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if arr[i] > arr[j] && dp[i] < dp[j]+1 {
				dp[i] = dp[j] + 1
			}
		}
	}

	// Find maximum value in dp array
	maxLength := dp[0]
	for i := 1; i < n; i++ {
		if dp[i] > maxLength {
			maxLength = dp[i]
		}
	}

	return maxLength
}

// LISOptimized finds LIS using binary search optimization
// Time Complexity: O(n log n), Space Complexity: O(n)
func LISOptimized(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	// tails[i] stores the smallest tail of all increasing subsequences of length i+1
	tails := make([]int, 0, len(arr))

	for _, num := range arr {
		// Binary search for the position to insert/replace
		left, right := 0, len(tails)
		for left < right {
			mid := (left + right) / 2
			if tails[mid] < num {
				left = mid + 1
			} else {
				right = mid
			}
		}

		// If num is larger than all elements in tails, append it
		if left == len(tails) {
			tails = append(tails, num)
		} else {
			// Replace the first element that is >= num
			tails[left] = num
		}
	}

	return len(tails)
}

// KnapsackItem represents an item with weight and value
type KnapsackItem struct {
	Weight int
	Value  int
}

// Knapsack01 solves the 0/1 knapsack problem
// Time Complexity: O(n*W), Space Complexity: O(n*W)
func Knapsack01(items []KnapsackItem, capacity int) (int, []int) {
	if len(items) == 0 || capacity <= 0 {
		return 0, []int{}
	}

	n := len(items)
	// dp[i][w] = maximum value that can be obtained with weight less than or equal to w using items up to i
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	// Build table dp[][] in bottom-up manner
	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			// If weight of the nth item is more than capacity w, it cannot be included
			if items[i-1].Weight > w {
				dp[i][w] = dp[i-1][w]
			} else {
				// Return the maximum of two cases:
				// 1. nth item included
				// 2. not included
				include := items[i-1].Value + dp[i-1][w-items[i-1].Weight]
				exclude := dp[i-1][w]
				if include > exclude {
					dp[i][w] = include
				} else {
					dp[i][w] = exclude
				}
			}
		}
	}

	// Backtrack to find which items to include
	selectedItems := make([]int, 0)
	w := capacity
	for i := n; i > 0 && w > 0; i-- {
		// If value comes from the top (dp[i-1][w]), then the item is not included
		if dp[i][w] != dp[i-1][w] {
			selectedItems = append(selectedItems, i-1) // Add item index
			w -= items[i-1].Weight
		}
	}

	// Reverse the selectedItems slice
	for i, j := 0, len(selectedItems)-1; i < j; i, j = i+1, j-1 {
		selectedItems[i], selectedItems[j] = selectedItems[j], selectedItems[i]
	}

	return dp[n][capacity], selectedItems
}

// CoinChange finds the minimum number of coins needed to make the amount
// Time Complexity: O(amount * n), Space Complexity: O(amount)
func CoinChange(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}
	if len(coins) == 0 {
		return -1
	}

	// dp[i] represents the minimum number of coins needed to make amount i
	dp := make([]int, amount+1)
	
	// Initialize with a value larger than any possible answer
	for i := 1; i <= amount; i++ {
		dp[i] = amount + 1
	}

	for i := 1; i <= amount; i++ {
		for _, coin := range coins {
			if coin <= i {
				if dp[i-coin]+1 < dp[i] {
					dp[i] = dp[i-coin] + 1
				}
			}
		}
	}

	if dp[amount] > amount {
		return -1 // No solution
	}
	return dp[amount]
}

// CoinChangeWays finds the number of ways to make the amount using given coins
// Time Complexity: O(amount * n), Space Complexity: O(amount)
func CoinChangeWays(coins []int, amount int) int {
	if amount == 0 {
		return 1
	}
	if len(coins) == 0 {
		return 0
	}

	// dp[i] represents the number of ways to make amount i
	dp := make([]int, amount+1)
	dp[0] = 1 // One way to make amount 0

	for _, coin := range coins {
		for i := coin; i <= amount; i++ {
			dp[i] += dp[i-coin]
		}
	}

	return dp[amount]
}

// MaxSubarraySum finds the maximum sum of a contiguous subarray (Kadane's Algorithm)
// Time Complexity: O(n), Space Complexity: O(1)
func MaxSubarraySum(arr []int) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("array is empty")
	}

	maxSoFar := arr[0]
	maxEndingHere := arr[0]

	for i := 1; i < len(arr); i++ {
		// Either extend the existing subarray or start a new one
		if maxEndingHere+arr[i] > arr[i] {
			maxEndingHere = maxEndingHere + arr[i]
		} else {
			maxEndingHere = arr[i]
		}

		// Update the maximum sum found so far
		if maxEndingHere > maxSoFar {
			maxSoFar = maxEndingHere
		}
	}

	return maxSoFar, nil
}

// MaxSubarrayWithIndices returns the maximum sum and the start/end indices
// Time Complexity: O(n), Space Complexity: O(1)
func MaxSubarrayWithIndices(arr []int) (int, int, int, error) {
	if len(arr) == 0 {
		return 0, 0, 0, errors.New("array is empty")
	}

	maxSoFar := arr[0]
	maxEndingHere := arr[0]
	start, end, tempStart := 0, 0, 0

	for i := 1; i < len(arr); i++ {
		if maxEndingHere < 0 {
			maxEndingHere = arr[i]
			tempStart = i
		} else {
			maxEndingHere += arr[i]
		}

		if maxEndingHere > maxSoFar {
			maxSoFar = maxEndingHere
			start = tempStart
			end = i
		}
	}

	return maxSoFar, start, end, nil
}

// HouseRobber finds the maximum amount that can be robbed from houses in a line
// Time Complexity: O(n), Space Complexity: O(1)
func HouseRobber(houses []int) int {
	if len(houses) == 0 {
		return 0
	}
	if len(houses) == 1 {
		return houses[0]
	}

	// prev2 = maximum money robbed up to house i-2
	// prev1 = maximum money robbed up to house i-1
	prev2, prev1 := 0, houses[0]

	for i := 1; i < len(houses); i++ {
		// Current max is either:
		// 1. Rob current house + max from i-2
		// 2. Don't rob current house, take max from i-1
		current := int(math.Max(float64(houses[i]+prev2), float64(prev1)))
		prev2 = prev1
		prev1 = current
	}

	return prev1
}

// HouseRobberCircular solves house robber problem where houses are in a circle
// Time Complexity: O(n), Space Complexity: O(1)
func HouseRobberCircular(houses []int) int {
	if len(houses) == 0 {
		return 0
	}
	if len(houses) == 1 {
		return houses[0]
	}
	if len(houses) == 2 {
		return int(math.Max(float64(houses[0]), float64(houses[1])))
	}

	// Case 1: Rob houses 0 to n-2 (exclude last house)
	case1 := robLinear(houses[:len(houses)-1])
	
	// Case 2: Rob houses 1 to n-1 (exclude first house)
	case2 := robLinear(houses[1:])

	return int(math.Max(float64(case1), float64(case2)))
}

// robLinear is a helper function for linear house robber
func robLinear(houses []int) int {
	if len(houses) == 0 {
		return 0
	}
	if len(houses) == 1 {
		return houses[0]
	}

	prev2, prev1 := 0, houses[0]
	for i := 1; i < len(houses); i++ {
		current := int(math.Max(float64(houses[i]+prev2), float64(prev1)))
		prev2 = prev1
		prev1 = current
	}

	return prev1
}

// EditDistance calculates the minimum edit distance between two strings
// Time Complexity: O(m*n), Space Complexity: O(m*n)
func EditDistance(str1, str2 string) int {
	m, n := len(str1), len(str2)
	
	// Create a DP table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize base cases
	for i := 0; i <= m; i++ {
		dp[i][0] = i // Delete all characters from str1
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j // Insert all characters to str1
	}

	// Fill the DP table
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if str1[i-1] == str2[j-1] {
				dp[i][j] = dp[i-1][j-1] // No operation needed
			} else {
				// Take minimum of three operations
				insert := dp[i][j-1] + 1
				delete := dp[i-1][j] + 1
				replace := dp[i-1][j-1] + 1
				
				dp[i][j] = int(math.Min(float64(insert), math.Min(float64(delete), float64(replace))))
			}
		}
	}

	return dp[m][n]
}
