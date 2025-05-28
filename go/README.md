# Algorithms Collection - Go Implementation

This directory contains Go implementations of various algorithms including sorting, searching, graph algorithms, and machine learning basics.

## 🚀 Features

- **Sorting Algorithms**: Bubble Sort, Quick Sort, Merge Sort, Heap Sort
- **Search Algorithms**: Binary Search, Linear Search
- **Graph Algorithms**: Dijkstra's Algorithm, A* Pathfinding
- **Machine Learning**: Linear Regression, K-Means, K-NN
- **Data Structures**: Linked Lists, Stacks, Queues, Trees
- **String Algorithms**: Pattern matching, manipulation
- **Mathematical Algorithms**: Fibonacci, Prime numbers, GCD

## 📁 Project Structure

```
go/
├── algorithms/
│   ├── sorting/          # Sorting algorithms
│   ├── searching/        # Search algorithms
│   ├── graphs/          # Graph algorithms
│   ├── ml/              # Machine learning algorithms
│   ├── strings/         # String algorithms
│   ├── math/            # Mathematical algorithms
│   └── datastructures/  # Data structure implementations
├── examples/            # Usage examples
├── benchmarks/          # Performance benchmarks
├── utils/              # Utility functions
└── tests/              # Integration tests
```

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/algorithms-collection/algorithms.git
cd algorithms/go

# Download dependencies
go mod download

# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./...
```

## 📖 Usage

### Sorting Algorithms

```go
package main

import (
    "fmt"
    "github.com/algorithms-collection/go/algorithms/sorting"
)

func main() {
    arr := []int{64, 34, 25, 12, 22, 11, 90}
    
    // Bubble Sort
    bubbleSorted := sorting.BubbleSort(arr)
    fmt.Println("Bubble Sort:", bubbleSorted)
    
    // Quick Sort
    quickSorted := sorting.QuickSort(arr)
    fmt.Println("Quick Sort:", quickSorted)
    
    // Merge Sort
    mergeSorted := sorting.MergeSort(arr)
    fmt.Println("Merge Sort:", mergeSorted)
}
```

### Search Algorithms

```go
package main

import (
    "fmt"
    "github.com/algorithms-collection/go/algorithms/searching"
)

func main() {
    arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
    target := 7
    
    // Binary Search
    index := searching.BinarySearch(arr, target)
    fmt.Printf("Binary Search: Element %d found at index %d\n", target, index)
    
    // Linear Search
    index = searching.LinearSearch(arr, target)
    fmt.Printf("Linear Search: Element %d found at index %d\n", target, index)
}
```

### Graph Algorithms

```go
package main

import (
    "fmt"
    "github.com/algorithms-collection/go/algorithms/graphs"
)

func main() {
    // Create a weighted graph
    graph := graphs.NewWeightedGraph()
    graph.AddEdge("A", "B", 4)
    graph.AddEdge("A", "C", 2)
    graph.AddEdge("B", "C", 1)
    graph.AddEdge("B", "D", 5)
    graph.AddEdge("C", "D", 8)
    
    // Find shortest path using Dijkstra's algorithm
    distances := graph.Dijkstra("A")
    fmt.Println("Shortest distances from A:", distances)
}
```

### Machine Learning

```go
package main

import (
    "fmt"
    "github.com/algorithms-collection/go/algorithms/ml"
)

func main() {
    // Linear Regression
    X := [][]float64{{1}, {2}, {3}, {4}, {5}}
    y := []float64{2, 4, 6, 8, 10}
    
    lr := ml.NewLinearRegression()
    lr.Fit(X, y)
    
    predictions := lr.Predict([][]float64{{6}, {7}})
    fmt.Println("Linear Regression predictions:", predictions)
    
    // K-Means Clustering
    data := [][]float64{
        {1, 1}, {1, 2}, {2, 1}, {2, 2},
        {8, 8}, {8, 9}, {9, 8}, {9, 9},
    }
    
    kmeans := ml.NewKMeans(2, 100)
    clusters := kmeans.Fit(data)
    fmt.Println("K-Means clusters:", clusters)
}
```

## 🧪 Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run specific package tests:
```bash
go test ./algorithms/sorting
go test ./algorithms/searching
go test ./algorithms/graphs
```

## 📊 Benchmarks

Run performance benchmarks:
```bash
go test -bench=. ./...
```

Run specific benchmarks:
```bash
go test -bench=BenchmarkBubbleSort ./algorithms/sorting
go test -bench=BenchmarkQuickSort ./algorithms/sorting
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🔗 Related

- [JavaScript Implementation](../javascript/)
- [Python Implementation](../python/)
- [Algorithm Visualizations](../javascript/visualizations/)

## 📚 Algorithm Complexity Reference

| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Bubble Sort | O(n²) | O(1) |
| Quick Sort | O(n log n) avg, O(n²) worst | O(log n) |
| Merge Sort | O(n log n) | O(n) |
| Heap Sort | O(n log n) | O(1) |
| Binary Search | O(log n) | O(1) |
| Dijkstra's | O((V + E) log V) | O(V) |
| A* | O(b^d) | O(b^d) |
| Linear Regression | O(n) | O(1) |
| K-Means | O(n × k × i) | O(n + k) |
| K-NN | O(n × d) | O(n × d) |

Where:
- n = number of elements
- V = number of vertices
- E = number of edges
- b = branching factor
- d = depth/dimensions
- k = number of clusters
- i = number of iterations
