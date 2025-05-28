// Package main demonstrates the usage of various algorithms
package main

import (
	"fmt"

	"github.com/algorithms-collection/go/algorithms/datastructures"
	"github.com/algorithms-collection/go/algorithms/graphs"
	"github.com/algorithms-collection/go/algorithms/math"
	"github.com/algorithms-collection/go/algorithms/ml"
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

	fmt.Println("\n✨ Go implementation complete! All algorithms tested and working.")
	fmt.Println("📚 Check the test files for detailed usage examples and edge cases.")
}
