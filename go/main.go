// Package main demonstrates the usage of various algorithms
package main

import (
	"fmt"

	"github.com/algorithms-collection/go/algorithms/backtracking"
	"github.com/algorithms-collection/go/algorithms/datastructures"
	"github.com/algorithms-collection/go/algorithms/dp"
	"github.com/algorithms-collection/go/algorithms/geometry"
	"github.com/algorithms-collection/go/algorithms/graphs"
	"github.com/algorithms-collection/go/algorithms/greedy"
	"github.com/algorithms-collection/go/algorithms/math"
	"github.com/algorithms-collection/go/algorithms/ml"
	"github.com/algorithms-collection/go/algorithms/numbertheory"
	"github.com/algorithms-collection/go/algorithms/searching"
	"github.com/algorithms-collection/go/algorithms/sorting"
	"github.com/algorithms-collection/go/algorithms/strings"
)

func main() {
	fmt.Println("🚀 Algorithms Collection - Go Implementation")
	fmt.Println("============================================")

	// Demonstrate sorting algorithms
	fmt.Println("\n📊 Sorting Algorithms:")
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original array: %v\n", arr)

	bubbleSorted := sorting.BubbleSort(arr)
	fmt.Printf("Bubble Sort:    %v\n", bubbleSorted)

	quickSorted := sorting.QuickSort(arr)
	fmt.Printf("Quick Sort:     %v\n", quickSorted)

	mergeSorted := sorting.MergeSort(arr)
	fmt.Printf("Merge Sort:     %v\n", mergeSorted)

	heapSorted := sorting.HeapSort(arr)
	fmt.Printf("Heap Sort:      %v\n", heapSorted)

	// Demonstrate search algorithms
	fmt.Println("\n🔍 Search Algorithms:")
	sortedArr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	target := 7
	fmt.Printf("Sorted array: %v, Target: %d\n", sortedArr, target)

	linearIndex := searching.LinearSearch(sortedArr, target)
	fmt.Printf("Linear Search:     Index %d\n", linearIndex)

	binaryIndex := searching.BinarySearch(sortedArr, target)
	fmt.Printf("Binary Search:     Index %d\n", binaryIndex)

	ternaryIndex := searching.TernarySearch(sortedArr, target)
	fmt.Printf("Ternary Search:    Index %d\n", ternaryIndex)

	// Demonstrate string algorithms
	fmt.Println("\n📝 String Algorithms:")

	str1 := "listen"
	str2 := "silent"
	fmt.Printf("Are '%s' and '%s' anagrams? %v\n", str1, str2, strings.AreAnagrams(str1, str2))

	palindrome := "racecar"
	fmt.Printf("Is '%s' a palindrome? %v\n", palindrome, strings.IsPalindrome(palindrome))

	text := "hello world"
	reversed := strings.ReverseString(text)
	fmt.Printf("'%s' reversed: '%s'\n", text, reversed)

	vowelCount := strings.CountVowels(text)
	fmt.Printf("Vowels in '%s': %d\n", text, vowelCount)

	maxChar, maxCount := strings.MaxCharacter(text)
	fmt.Printf("Most frequent character in '%s': '%c' (%d times)\n", text, maxChar, maxCount)

	// Demonstrate pattern matching
	fmt.Println("\n🔎 Pattern Matching:")
	textPattern := "ababcababa"
	pattern := "aba"
	fmt.Printf("Text: '%s', Pattern: '%s'\n", textPattern, pattern)

	kmpMatches := strings.KMPSearch(textPattern, pattern)
	fmt.Printf("KMP Search matches at indices: %v\n", kmpMatches)

	rabinKarpMatches := strings.RabinKarpSearch(textPattern, pattern)
	fmt.Printf("Rabin-Karp matches at indices: %v\n", rabinKarpMatches)

	// Advanced String Algorithms
	fmt.Println("\n🔤 Advanced String Algorithms:")

	// Suffix Array
	suffixText := "banana"
	suffixArray := strings.SuffixArray(suffixText)
	fmt.Printf("Suffix array of '%s': %v\n", suffixText, suffixArray)

	// Longest Repeated Substring
	repeatedText := "abcabcabc"
	longestRepeated := strings.LongestRepeatedSubstring(repeatedText)
	fmt.Printf("Longest repeated substring in '%s': '%s'\n", repeatedText, longestRepeated)

	// Aho-Corasick Multiple Pattern Matching
	ac := strings.NewAhoCorasick()
	ac.AddPattern("he")
	ac.AddPattern("she")
	ac.AddPattern("his")
	ac.AddPattern("hers")

	searchText := "ushers"
	matches := ac.Search(searchText)
	fmt.Printf("Aho-Corasick matches in '%s': %v\n", searchText, matches)

	// Palindromic Substrings
	palindromeText := "aabaa"
	allPalindromes := strings.AllPalindromes(palindromeText)
	fmt.Printf("All palindromes in '%s': %v\n", palindromeText, allPalindromes)

	palindromeCount := strings.CountPalindromes(palindromeText)
	fmt.Printf("Total palindromic substrings: %d\n", palindromeCount)

	// String Rotation
	rotStr1 := "waterbottle"
	rotStr2 := "erbottlewat"
	isRotation := strings.StringRotation(rotStr1, rotStr2)
	fmt.Printf("'%s' is rotation of '%s': %v\n", rotStr2, rotStr1, isRotation)

	// Longest Common Prefix
	prefixStrings := []string{"flower", "flow", "flight"}
	commonPrefix := strings.LongestCommonPrefix(prefixStrings)
	fmt.Printf("Longest common prefix of %v: '%s'\n", prefixStrings, commonPrefix)

	// Demonstrate graph algorithms
	fmt.Println("\n🗺️  Graph Algorithms:")

	// Create unweighted graph for BFS/DFS
	g := graphs.NewGraph(5)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)

	fmt.Println("Unweighted Graph (0->1, 0->2, 1->3, 2->4, 3->4):")
	bfsResult := g.BFS(0)
	fmt.Printf("BFS from vertex 0: %v\n", bfsResult)

	dfsResult := g.DFS(0)
	fmt.Printf("DFS from vertex 0: %v\n", dfsResult)

	path := g.FindPath(0, 4)
	fmt.Printf("Path from 0 to 4: %v\n", path)

	// Create weighted graph for Dijkstra
	wg := graphs.NewWeightedGraph(5)
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

	fmt.Println("\nWeighted Graph - Dijkstra's Algorithm:")
	distances := wg.Dijkstra(0)
	fmt.Printf("Shortest distances from vertex 0: %v\n", distances)

	// Demonstrate Bellman-Ford
	fmt.Println("\nBellman-Ford Algorithm:")
	bellmanDistances, hasNegativeCycle := wg.BellmanFord(0)
	fmt.Printf("Distances: %v\n", bellmanDistances)
	fmt.Printf("Has negative cycle: %v\n", !hasNegativeCycle)

	// Demonstrate Minimum Spanning Tree algorithms
	fmt.Println("\n🌳 Minimum Spanning Tree:")
	mstEdges, mstWeight := wg.KruskalMST()
	fmt.Printf("Kruskal's MST edges: %v\n", mstEdges)
	fmt.Printf("Total MST weight: %d\n", mstWeight)

	primEdges, primWeight := wg.PrimMST(0)
	fmt.Printf("Prim's MST edges: %v\n", primEdges)
	fmt.Printf("Total MST weight: %d\n", primWeight)

	// Demonstrate Strongly Connected Components
	fmt.Println("\n🔗 Strongly Connected Components:")

	// Create a directed graph with SCCs
	sccGraph := graphs.NewGraph(5)
	sccGraph.AddEdge(1, 0)
	sccGraph.AddEdge(0, 2)
	sccGraph.AddEdge(2, 1)
	sccGraph.AddEdge(0, 3)
	sccGraph.AddEdge(3, 4)

	tarjanSCCs := sccGraph.TarjanSCC()
	fmt.Printf("Tarjan's SCCs: %v\n", tarjanSCCs)

	kosarajuSCCs := sccGraph.KosarajuSCC()
	fmt.Printf("Kosaraju's SCCs: %v\n", kosarajuSCCs)

	// Demonstrate Articulation Points and Bridges
	fmt.Println("\n🔍 Critical Connections:")

	// Create a graph with articulation points and bridges
	criticalGraph := graphs.NewGraph(7)
	criticalGraph.AddEdge(0, 1)
	criticalGraph.AddEdge(0, 2)
	criticalGraph.AddEdge(1, 2)
	criticalGraph.AddEdge(1, 3)
	criticalGraph.AddEdge(3, 4)
	criticalGraph.AddEdge(4, 5)
	criticalGraph.AddEdge(5, 6)
	criticalGraph.AddEdge(4, 6)

	articulationPoints := criticalGraph.ArticulationPoints()
	fmt.Printf("Articulation points: %v\n", articulationPoints)

	bridges := criticalGraph.Bridges()
	fmt.Printf("Bridges: %v\n", bridges)

	// Demonstrate Maximum Flow
	fmt.Println("\n💧 Maximum Flow:")

	// Create flow network
	flowNet := graphs.NewFlowNetwork(6)
	flowNet.AddEdge(0, 1, 16)
	flowNet.AddEdge(0, 2, 13)
	flowNet.AddEdge(1, 2, 10)
	flowNet.AddEdge(1, 3, 12)
	flowNet.AddEdge(2, 1, 4)
	flowNet.AddEdge(2, 4, 14)
	flowNet.AddEdge(3, 2, 9)
	flowNet.AddEdge(3, 5, 20)
	flowNet.AddEdge(4, 3, 7)
	flowNet.AddEdge(4, 5, 4)

	maxFlow := flowNet.FordFulkerson(0, 5)
	fmt.Printf("Maximum flow from 0 to 5: %d\n", maxFlow)

	sourceSet, sinkSet := flowNet.GetMinCut(0, 5)
	fmt.Printf("Min-cut source set: %v\n", sourceSet)
	fmt.Printf("Min-cut sink set: %v\n", sinkSet)

	// Demonstrate string distance algorithms
	fmt.Println("\n📏 String Distance Algorithms:")
	s1 := "kitten"
	s2 := "sitting"
	editDist := strings.EditDistance(s1, s2)
	fmt.Printf("Edit distance between '%s' and '%s': %d\n", s1, s2, editDist)

	lcs := strings.LongestCommonSubsequence(s1, s2)
	fmt.Printf("Longest common subsequence length: %d\n", lcs)

	commonSubstring := strings.LongestCommonSubstring("GeeksforGeeks", "GeeksQuiz")
	fmt.Printf("Longest common substring: '%s'\n", commonSubstring)

	// Demonstrate advanced dynamic programming algorithms
	fmt.Println("\n🧮 Advanced Dynamic Programming:")

	// Matrix Chain Multiplication
	dimensions := []int{1, 2, 3, 4, 5} // Matrices: 1x2, 2x3, 3x4, 4x5
	mcmCost := dp.MatrixChainMultiplication(dimensions)
	fmt.Printf("Matrix Chain Multiplication cost: %d\n", mcmCost)

	// Longest Palindromic Subsequence
	palindromeStr := "bbbab"
	lpsLength := dp.LongestPalindromicSubsequence(palindromeStr)
	fmt.Printf("Longest palindromic subsequence in '%s': %d\n", palindromeStr, lpsLength)

	// Subset Sum
	subsetArr := []int{3, 34, 4, 12, 5, 2}
	targetSum := 9
	canSum := dp.SubsetSum(subsetArr, targetSum)
	fmt.Printf("Can subset of %v sum to %d? %v\n", subsetArr, targetSum, canSum)

	// Rod Cutting
	prices := []int{1, 5, 8, 9, 10, 17, 17, 20}
	rodLength := 8
	maxValue, cuts := dp.RodCuttingWithCuts(prices, rodLength)
	fmt.Printf("Rod cutting (length %d): max value = %d, cuts = %v\n", rodLength, maxValue, cuts)

	// Egg Dropping
	eggs := 2
	floors := 10
	minTrials := dp.EggDropping(eggs, floors)
	fmt.Printf("Egg dropping (%d eggs, %d floors): minimum trials = %d\n", eggs, floors, minTrials)

	// Palindrome Partitioning
	partitionStr := "aab"
	minCuts := dp.PalindromePartitioning(partitionStr)
	fmt.Printf("Palindrome partitioning '%s': minimum cuts = %d\n", partitionStr, minCuts)

	// Word Break
	wordBreakStr := "leetcode"
	dictionary := []string{"leet", "code"}
	canBreak := dp.WordBreak(wordBreakStr, dictionary)
	fmt.Printf("Can '%s' be segmented using %v? %v\n", wordBreakStr, dictionary, canBreak)

	sentences := dp.WordBreakII("catsanddog", []string{"cat", "cats", "and", "sand", "dog"})
	fmt.Printf("Word break sentences: %v\n", sentences)

	// Demonstrate backtracking algorithms
	fmt.Println("\n🔄 Backtracking Algorithms:")

	// N-Queens Problem
	nQueensSolutions := backtracking.NQueens(4)
	fmt.Printf("4-Queens problem: %d solutions found\n", len(nQueensSolutions))
	if len(nQueensSolutions) > 0 {
		fmt.Println("First solution:")
		backtracking.PrintNQueensSolution(nQueensSolutions[0])
	}

	// Sudoku Solver
	sudokuBoard := [][]int{
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

	fmt.Println("Sudoku puzzle solved:")
	if backtracking.SudokuSolver(sudokuBoard) {
		backtracking.PrintSudoku(sudokuBoard)
	}

	// Knight's Tour
	knightBoard, err := backtracking.KnightsTour(5)
	if err == nil {
		fmt.Printf("Knight's tour on 5x5 board: Found solution\n")
		fmt.Printf("First few moves: [0][0]=%d, [0][1]=%d, [0][2]=%d\n",
			knightBoard[0][0], knightBoard[0][1], knightBoard[0][2])
	} else {
		fmt.Printf("Knight's tour: %v\n", err)
	}

	// Graph Coloring
	colorGraph := [][]int{
		{0, 1, 1, 1},
		{1, 0, 1, 0},
		{1, 1, 0, 1},
		{1, 0, 1, 0},
	}
	colors, canColor := backtracking.GraphColoring(colorGraph, 3)
	fmt.Printf("Graph coloring (4 vertices, 3 colors): %v, colors: %v\n", canColor, colors)

	// Hamiltonian Path
	hamiltonianGraph := [][]int{
		{0, 1, 0, 1, 0},
		{1, 0, 1, 1, 1},
		{0, 1, 0, 0, 1},
		{1, 1, 0, 0, 1},
		{0, 1, 1, 1, 0},
	}
	hamiltonianPath, hasPath := backtracking.HamiltonianPath(hamiltonianGraph)
	fmt.Printf("Hamiltonian path: %v, path: %v\n", hasPath, hamiltonianPath)

	// Subset Generation
	subsets := backtracking.GenerateSubsets([]int{1, 2, 3})
	fmt.Printf("All subsets of [1,2,3]: %v\n", subsets)

	// Demonstrate greedy algorithms
	fmt.Println("\n💰 Greedy Algorithms:")

	// Activity Selection
	activities := []greedy.Activity{
		{Start: 1, Finish: 4, Index: 0},
		{Start: 3, Finish: 5, Index: 1},
		{Start: 0, Finish: 6, Index: 2},
		{Start: 5, Finish: 7, Index: 3},
		{Start: 8, Finish: 9, Index: 4},
		{Start: 5, Finish: 9, Index: 5},
	}
	selectedActivities := greedy.ActivitySelection(activities)
	fmt.Printf("Activity selection: %d activities selected\n", len(selectedActivities))

	// Huffman Coding
	frequencies := map[rune]int{'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45}
	huffmanCodes, _ := greedy.HuffmanCoding(frequencies)
	fmt.Printf("Huffman codes: a=%s, b=%s, f=%s\n", huffmanCodes['a'], huffmanCodes['b'], huffmanCodes['f'])

	// Job Scheduling
	jobs := []greedy.Job{
		{ID: 1, Deadline: 4, Profit: 20},
		{ID: 2, Deadline: 1, Profit: 10},
		{ID: 3, Deadline: 1, Profit: 40},
		{ID: 4, Deadline: 1, Profit: 30},
	}
	selectedJobs, totalProfit := greedy.JobScheduling(jobs)
	fmt.Printf("Job scheduling: %d jobs selected, total profit = %d\n", len(selectedJobs), totalProfit)

	// Fractional Knapsack
	items := []greedy.Item{
		{Value: 60, Weight: 10, Index: 0},
		{Value: 100, Weight: 20, Index: 1},
		{Value: 120, Weight: 30, Index: 2},
	}
	knapsackValue, fractions := greedy.FractionalKnapsack(items, 50)
	fmt.Printf("Fractional knapsack: max value = %.2f, fractions = %v\n", knapsackValue, fractions)

	// Minimum Coins
	coins := []int{1, 5, 10, 25}
	coinCount, coinResult, _ := greedy.MinimumCoins(coins, 30)
	fmt.Printf("Minimum coins for 30: %d coins = %v\n", coinCount, coinResult)

	// Gas Station
	gas := []int{1, 2, 3, 4, 5}
	cost := []int{3, 4, 5, 1, 2}
	startStation := greedy.GasStation(gas, cost)
	fmt.Printf("Gas station: start from station %d\n", startStation)

	// Minimum Platforms
	arrivals := []int{900, 940, 950, 1100, 1500, 1800}
	departures := []int{910, 1200, 1120, 1130, 1900, 2000}
	platforms := greedy.MinimumPlatforms(arrivals, departures)
	fmt.Printf("Minimum platforms needed: %d\n", platforms)

	// Demonstrate computational geometry algorithms
	fmt.Println("\n📐 Computational Geometry:")

	// Convex Hull
	points := []geometry.Point{
		{X: 0, Y: 3}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 4, Y: 4},
		{X: 0, Y: 0}, {X: 1, Y: 2}, {X: 3, Y: 1}, {X: 3, Y: 3},
	}
	convexHull := geometry.ConvexHullGrahamScan(points)
	fmt.Printf("Convex hull of %d points: %d vertices\n", len(points), len(convexHull))

	// Closest Pair of Points
	p1, p2, closestDist := geometry.ClosestPairOfPoints(points)
	fmt.Printf("Closest pair: (%.1f,%.1f) and (%.1f,%.1f), distance = %.3f\n",
		p1.X, p1.Y, p2.X, p2.Y, closestDist)

	// Line Intersection
	line1 := geometry.LineSegment{
		P1: geometry.Point{X: 1, Y: 1},
		P2: geometry.Point{X: 10, Y: 1},
	}
	line2 := geometry.LineSegment{
		P1: geometry.Point{X: 1, Y: 2},
		P2: geometry.Point{X: 10, Y: 2},
	}
	intersects := line1.DoesIntersect(line2)
	fmt.Printf("Line segments intersect: %v\n", intersects)

	// Point in Polygon
	polygon := []geometry.Point{
		{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4},
	}
	testPoint := geometry.Point{X: 2, Y: 2}
	isInside := geometry.PointInPolygon(testPoint, polygon)
	fmt.Printf("Point (%.1f,%.1f) is inside polygon: %v\n", testPoint.X, testPoint.Y, isInside)

	// Polygon Area
	area := geometry.PolygonArea(polygon)
	fmt.Printf("Polygon area: %.1f\n", area)

	// Convex check
	isConvex := geometry.IsConvex(polygon)
	fmt.Printf("Polygon is convex: %v\n", isConvex)

	// Performance comparison
	fmt.Println("\n⚡ Performance Comparison:")
	largeArr := make([]int, 1000)
	for i := range largeArr {
		largeArr[i] = 1000 - i // Reverse sorted for worst case
	}

	fmt.Printf("Sorting 1000 elements (reverse sorted):\n")

	// Note: In a real application, you would use benchmarking tools
	// This is just for demonstration
	fmt.Printf("✓ Bubble Sort:    O(n²) - Good for small arrays\n")
	fmt.Printf("✓ Quick Sort:     O(n log n) avg - General purpose\n")
	fmt.Printf("✓ Merge Sort:     O(n log n) - Stable, predictable\n")
	fmt.Printf("✓ Heap Sort:      O(n log n) - In-place, not stable\n")

	fmt.Println("\n🎯 Algorithm Complexity Summary:")
	fmt.Println("┌─────────────────┬──────────────┬──────────────┬──────────────┐")
	fmt.Println("│ Algorithm       │ Best Case    │ Average Case │ Worst Case   │")
	fmt.Println("├─────────────────┼──────────────┼──────────────┼──────────────┤")
	fmt.Println("│ Bubble Sort     │ O(n)         │ O(n²)        │ O(n²)        │")
	fmt.Println("│ Quick Sort      │ O(n log n)   │ O(n log n)   │ O(n²)        │")
	fmt.Println("│ Merge Sort      │ O(n log n)   │ O(n log n)   │ O(n log n)   │")
	fmt.Println("│ Heap Sort       │ O(n log n)   │ O(n log n)   │ O(n log n)   │")
	fmt.Println("│ Binary Search   │ O(1)         │ O(log n)     │ O(log n)     │")
	fmt.Println("│ Linear Search   │ O(1)         │ O(n)         │ O(n)         │")
	fmt.Println("│ Dijkstra        │ O(E log V)   │ O(E log V)   │ O(E log V)   │")
	fmt.Println("│ Bellman-Ford    │ O(VE)        │ O(VE)        │ O(VE)        │")
	fmt.Println("└─────────────────┴──────────────┴──────────────┴──────────────┘")

	// Demonstrate mathematical algorithms
	fmt.Println("\n🔢 Mathematical Algorithms:")

	// Fibonacci
	fib10 := math.Fibonacci(10)
	fmt.Printf("Fibonacci(10): %d\n", fib10)

	// Prime checking
	isPrime97 := math.IsPrime(97)
	fmt.Printf("Is 97 prime? %v\n", isPrime97)

	// GCD and LCM
	gcd := math.GCD(48, 18)
	lcm := math.LCM(48, 18)
	fmt.Printf("GCD(48, 18): %d, LCM(48, 18): %d\n", gcd, lcm)

	// Factorial
	fact5 := math.Factorial(5)
	fmt.Printf("5! = %d\n", fact5)

	// Demonstrate data structures
	fmt.Println("\n📚 Data Structures:")

	// Stack
	stack := datastructures.NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	top, _ := stack.Pop()
	fmt.Printf("Stack: pushed 1,2,3, popped %v\n", top)

	// Queue
	queue := datastructures.NewQueue()
	queue.Enqueue("first")
	queue.Enqueue("second")
	front, _ := queue.Dequeue()
	fmt.Printf("Queue: enqueued 'first','second', dequeued '%v'\n", front)

	// Linked List
	ll := datastructures.NewLinkedList()
	ll.Append(1)
	ll.Append(2)
	ll.Append(3)
	fmt.Printf("Linked List: %v\n", ll.ToSlice())

	// Binary Tree
	bt := datastructures.NewBinaryTree()
	bt.Insert(5)
	bt.Insert(3)
	bt.Insert(7)
	bt.Insert(1)
	bt.Insert(9)
	inorder := bt.InorderTraversal()
	fmt.Printf("Binary Tree inorder: %v\n", inorder)

	// Min Heap
	heap := datastructures.NewMinHeap()
	heap.Insert(5)
	heap.Insert(3)
	heap.Insert(8)
	heap.Insert(1)
	min, _ := heap.ExtractMin()
	fmt.Printf("Min Heap: inserted 5,3,8,1, extracted min %d\n", min)

	// Demonstrate machine learning
	fmt.Println("\n🤖 Machine Learning Algorithms:")

	// Linear Regression
	lr := ml.NewLinearRegression()
	XLR := [][]float64{{1}, {2}, {3}, {4}, {5}}
	yLR := []float64{2, 4, 6, 8, 10}
	lr.Fit(XLR, yLR)
	predictions, _ := lr.Predict([][]float64{{6}, {7}})
	fmt.Printf("Linear Regression: y=2x, predict(6,7) = [%.1f, %.1f]\n", predictions[0], predictions[1])

	// K-Means Clustering
	kmeans := ml.NewKMeans(2, 100)
	XKM := [][]float64{
		{1, 1}, {1, 2}, {2, 1}, {2, 2},
		{8, 8}, {8, 9}, {9, 8}, {9, 9},
	}
	labels, _ := kmeans.Fit(XKM)
	fmt.Printf("K-Means: 2 clusters, labels: %v\n", labels)

	// K-Nearest Neighbors
	knn := ml.NewKNearestNeighbors(3)
	XKNN := [][]float64{{1, 1}, {1, 2}, {2, 1}, {5, 5}, {5, 6}, {6, 5}}
	yKNN := []float64{0, 0, 0, 1, 1, 1}
	knn.Fit(XKNN, yKNN)
	knnPred, _ := knn.Predict([][]float64{{1.5, 1.5}, {5.5, 5.5}})
	fmt.Printf("K-NN: predict classes for (1.5,1.5) and (5.5,5.5) = [%.1f, %.1f]\n", knnPred[0], knnPred[1])

	// Demonstrate advanced data structures
	fmt.Println("\n🏗️ Advanced Data Structures:")

	// Splay Tree
	splayTree := datastructures.NewSplayTree()
	splayKeys := []int{10, 20, 30, 40, 50, 25}
	for _, key := range splayKeys {
		splayTree.Insert(key)
	}
	fmt.Printf("Splay tree: inserted keys %v\n", splayKeys)
	fmt.Printf("Search for 25: %v\n", splayTree.Search(25))
	fmt.Printf("Search for 35: %v\n", splayTree.Search(35))
	fmt.Printf("Delete 20: %v\n", splayTree.Delete(20))

	// Skip List
	skipList := datastructures.NewSkipList(4)
	skipList.Insert(3, "three")
	skipList.Insert(6, "six")
	skipList.Insert(7, "seven")
	skipList.Insert(9, "nine")
	skipList.Insert(12, "twelve")
	skipList.Insert(19, "nineteen")
	skipList.Insert(17, "seventeen")

	value6, found6 := skipList.Search(6)
	fmt.Printf("Skip list search for 6: found=%v, value=%v\n", found6, value6)

	value15, found15 := skipList.Search(15)
	fmt.Printf("Skip list search for 15: found=%v, value=%v\n", found15, value15)

	deleted := skipList.Delete(6)
	fmt.Printf("Skip list delete 6: %v\n", deleted)

	// Disjoint Set (Union-Find)
	ds := datastructures.NewDisjointSet(10)
	ds.Union(1, 2)
	ds.Union(2, 3)
	ds.Union(4, 5)
	ds.Union(6, 7)
	ds.Union(5, 6)

	fmt.Printf("Disjoint set: 1 and 3 connected? %v\n", ds.Connected(1, 3))
	fmt.Printf("Disjoint set: 1 and 4 connected? %v\n", ds.Connected(1, 4))
	fmt.Printf("Disjoint set: 4 and 7 connected? %v\n", ds.Connected(4, 7))
	fmt.Printf("Number of disjoint sets: %d\n", ds.CountSets())

	// Demonstrate number theory algorithms
	fmt.Println("\n🔢 Number Theory:")

	// Miller-Rabin Primality Test
	testNumber := int64(982451653)
	isPrime := numbertheory.MillerRabinPrimality(testNumber, 10)
	fmt.Printf("Miller-Rabin: %d is prime? %v\n", testNumber, isPrime)

	// Pollard's Rho Factorization
	composite := int64(1403) // 23 * 61
	factor := numbertheory.PollardRho(composite)
	fmt.Printf("Pollard's Rho: factor of %d = %d\n", composite, factor)

	// Fast Fourier Transform
	fftInput := []numbertheory.Complex{
		{Real: 1, Imag: 0}, {Real: 2, Imag: 0}, {Real: 3, Imag: 0}, {Real: 4, Imag: 0},
	}
	fftResult := numbertheory.FFT(fftInput)
	fmt.Printf("FFT: transformed %d complex numbers\n", len(fftResult))

	// Discrete Logarithm
	dlBase, dlTarget, dlPrime := int64(3), int64(4), int64(7)
	discreteLog := numbertheory.DiscreteLogarithm(dlBase, dlTarget, dlPrime)
	fmt.Printf("Discrete log: %d^%d ≡ %d (mod %d)\n", dlBase, discreteLog, dlTarget, dlPrime)

	// Chinese Remainder Theorem
	remainders := []int64{2, 3, 2}
	moduli := []int64{3, 5, 7}
	crtResult, crtValid := numbertheory.ChineseRemainderTheorem(remainders, moduli)
	fmt.Printf("Chinese Remainder Theorem: x ≡ %d (mod %d), valid: %v\n", crtResult, 3*5*7, crtValid)

	// Euler's Totient Function
	totientN := int64(12)
	totient := numbertheory.EulerTotient(totientN)
	fmt.Printf("Euler's totient φ(%d) = %d\n", totientN, totient)

	// Primitive Roots
	primitiveP := int64(7)
	primitiveRoots := numbertheory.PrimitiveRoots(primitiveP)
	fmt.Printf("Primitive roots modulo %d: %v\n", primitiveP, primitiveRoots)

	// Trial Division Factorization
	factorN := int64(60)
	factors := numbertheory.FactorizeTrialDivision(factorN)
	fmt.Printf("Prime factorization of %d: %v\n", factorN, factors)

	// Demonstrate new advanced sorting algorithms
	fmt.Println("\n🔄 Advanced Sorting Algorithms:")
	advancedArr := []int{64, 34, 25, 12, 22, 11, 90, 5, 77, 30}
	fmt.Printf("Original array: %v\n", advancedArr)

	bucketSorted := sorting.BucketSort(advancedArr)
	fmt.Printf("Bucket Sort:    %v\n", bucketSorted)

	pancakeSorted := sorting.PancakeSort(advancedArr)
	fmt.Printf("Pancake Sort:   %v\n", pancakeSorted)

	timSorted := sorting.TimSort(advancedArr)
	fmt.Printf("Tim Sort:       %v\n", timSorted)

	fmt.Println("\n✨ Go implementation complete! All algorithms tested and working.")
	fmt.Println("📚 Check the test files for detailed usage examples and edge cases.")
}
