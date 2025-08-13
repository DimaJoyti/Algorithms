# Algorithms Collection - Go Implementation

A complete, production-ready collection of algorithms and data structures implemented in Go. This library contains over **100+ algorithms** across **10 major categories**, making it one of the most comprehensive algorithm collections available for Go developers.

## 🌟 **Key Features**

- ✅ **100+ Algorithms** across 10 major categories
- ✅ **Production-ready code** with comprehensive error handling
- ✅ **Extensive test coverage** with edge cases and benchmarks
- ✅ **Complete documentation** with time/space complexity analysis
- ✅ **Clean, idiomatic Go** following best practices
- ✅ **Zero external dependencies** (uses only Go standard library)
- ✅ **Educational examples** with detailed explanations

## � **Algorithm Categories**

### 🔄 **Sorting Algorithms** (15+ algorithms)
- **Basic**: Bubble, Selection, Insertion, Shell, Cocktail, Gnome, Pancake, Comb
- **Advanced**: Merge, Quick, Heap, Counting, Radix, Bucket, Tim Sort
- **Novelty**: Bogo Sort (for educational purposes)

### 🔍 **Searching Algorithms** (8+ algorithms)
- **Linear & Binary**: Linear Search, Binary Search (iterative & recursive)
- **Advanced**: Ternary, Exponential, Interpolation, Fibonacci, Jump Search

### 🌐 **Graph Algorithms** (15+ algorithms)
- **Traversal**: BFS, DFS, Cycle Detection
- **Shortest Path**: Dijkstra, Bellman-Ford, Floyd-Warshall, A*
- **Minimum Spanning Tree**: Kruskal's, Prim's algorithms
- **Connectivity**: Strongly Connected Components (Tarjan's, Kosaraju's)
- **Critical Elements**: Articulation Points, Bridges
- **Flow Networks**: Ford-Fulkerson, Edmonds-Karp, Min-Cut
- **Topological**: Topological Sort

### 🧮 **Dynamic Programming** (15+ algorithms)
- **Classic**: Fibonacci, Knapsack (0/1 & Unbounded), Coin Change
- **Sequences**: Longest Increasing/Common Subsequence, Edit Distance
- **Strings**: Longest Palindromic Subsequence, Palindrome Partitioning
- **Optimization**: Matrix Chain Multiplication, Rod Cutting, Egg Dropping
- **Path Problems**: Unique Paths, Minimum Path Sum, Climbing Stairs
- **Advanced**: Word Break, Subset Sum, House Robber

### 🔄 **Backtracking Algorithms** (6+ algorithms)
- **Constraint Satisfaction**: N-Queens Problem, Sudoku Solver
- **Path Finding**: Knight's Tour, Hamiltonian Path
- **Combinatorial**: Graph Coloring, Subset Generation

### 💰 **Greedy Algorithms** (8+ algorithms)
- **Scheduling**: Activity Selection, Job Scheduling, Interval Scheduling
- **Optimization**: Fractional Knapsack, Minimum Coins, Gas Station
- **Compression**: Huffman Coding
- **Resource Allocation**: Minimum Platforms

### 📐 **Computational Geometry** (8+ algorithms)
- **Convex Hull**: Graham Scan algorithm
- **Distance**: Closest Pair of Points (divide & conquer)
- **Intersection**: Line Segment Intersection
- **Containment**: Point in Polygon (ray casting)
- **Measurements**: Polygon Area (shoelace formula), Convex Polygon Diameter
- **Properties**: Convexity Testing

### 🔤 **String Algorithms** (15+ algorithms)
- **Pattern Matching**: KMP, Rabin-Karp, Boyer-Moore, Z-Algorithm
- **Multiple Patterns**: Aho-Corasick Algorithm
- **Palindromes**: Manacher's Algorithm, Palindromic Tree, All Palindromes
- **Advanced**: Suffix Array, Longest Repeated Substring
- **Utilities**: String Rotation, Longest Common Prefix, Anagram Detection

### 🔢 **Number Theory** (15+ algorithms)
- **Primality**: Miller-Rabin Test, Trial Division
- **Factorization**: Pollard's Rho, Trial Division
- **Modular Arithmetic**: Extended Euclidean, Chinese Remainder Theorem
- **Advanced**: Fast Fourier Transform, Discrete Logarithm
- **Properties**: Euler's Totient, Primitive Roots, Jacobi Symbol, Quadratic Residues

### 🏗️ **Data Structures** (15+ structures)
- **Linear**: Linked Lists, Stacks, Queues, Deques
- **Trees**: Binary Search Tree, AVL Tree, Red-Black Tree, B-Tree, Splay Tree
- **Heaps**: Min/Max Heap, Priority Queue
- **Advanced**: Trie, Segment Tree, Fenwick Tree, Skip List
- **Sets**: Disjoint Set (Union-Find) with path compression
- **Hashing**: Hash Table, Bloom Filter

### 🤖 **Machine Learning** (3+ algorithms)
- **Regression**: Linear Regression
- **Clustering**: K-Means
- **Classification**: K-Nearest Neighbors

## 📁 **Project Structure**

```
go/
├── algorithms/
│   ├── sorting/          # 15+ sorting algorithms
│   ├── searching/        # 8+ search algorithms
│   ├── graphs/          # 15+ graph algorithms
│   ├── dp/              # 15+ dynamic programming algorithms
│   ├── backtracking/    # 6+ backtracking algorithms
│   ├── greedy/          # 8+ greedy algorithms
│   ├── geometry/        # 8+ computational geometry algorithms
│   ├── strings/         # 15+ string algorithms
│   ├── numbertheory/    # 15+ number theory algorithms
│   ├── datastructures/  # 15+ data structure implementations
│   ├── trees/           # Advanced tree structures
│   ├── math/            # Mathematical algorithms
│   └── ml/              # Machine learning algorithms
├── main.go              # Comprehensive demonstrations
└── README.md            # This documentation
```

## 🛠️ **Installation & Quick Start**

```bash
# Clone the repository
git clone https://github.com/DimaJoyti/Algorithms.git
cd Algorithms/go

# Initialize Go module (if needed)
go mod init github.com/algorithms-collection/go

# Run the comprehensive demonstration
go run main.go

# Run all tests
go test ./...

# Run benchmarks
go test -bench=. ./...

# Test specific algorithm category
go test ./algorithms/graphs -v
go test ./algorithms/backtracking -v
go test ./algorithms/geometry -v
```

## 🎯 **Quick Examples**

### Graph Algorithms
```go
import "github.com/algorithms-collection/go/algorithms/graphs"

// Create a weighted graph
g := graphs.NewWeightedGraph(5)
g.AddEdge(0, 1, 10)
g.AddEdge(0, 4, 5)
g.AddEdge(1, 2, 1)
g.AddEdge(1, 4, 2)

// Find shortest paths
distances := g.Dijkstra(0)
fmt.Printf("Shortest distances: %v\n", distances)

// Find minimum spanning tree
mstEdges, weight := g.KruskalMST()
fmt.Printf("MST weight: %d\n", weight)

// Find strongly connected components
sccs := g.TarjanSCC()
fmt.Printf("SCCs: %v\n", sccs)
```

### Backtracking Algorithms
```go
import "github.com/algorithms-collection/go/algorithms/backtracking"

// Solve N-Queens problem
solutions := backtracking.NQueens(8)
fmt.Printf("8-Queens has %d solutions\n", len(solutions))

// Solve Sudoku
sudoku := [][]int{
    {5, 3, 0, 0, 7, 0, 0, 0, 0},
    // ... more rows
}
if backtracking.SudokuSolver(sudoku) {
    fmt.Println("Sudoku solved!")
}
```

### Greedy Algorithms
```go
import "github.com/algorithms-collection/go/algorithms/greedy"

// Activity selection
activities := []greedy.Activity{
    {Start: 1, Finish: 4, Index: 0},
    {Start: 3, Finish: 5, Index: 1},
    // ... more activities
}
selected := greedy.ActivitySelection(activities)

// Huffman coding
frequencies := map[rune]int{'a': 5, 'b': 9, 'c': 12}
codes, tree := greedy.HuffmanCoding(frequencies)
```

### Computational Geometry
```go
import "github.com/algorithms-collection/go/algorithms/geometry"

// Find convex hull
points := []geometry.Point{
    {X: 0, Y: 3}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 4, Y: 4},
}
hull := geometry.ConvexHullGrahamScan(points)

// Check point in polygon
polygon := []geometry.Point{{0, 0}, {4, 0}, {4, 4}, {0, 4}}
inside := geometry.PointInPolygon(geometry.Point{X: 2, Y: 2}, polygon)
```

### Number Theory
```go
import "github.com/algorithms-collection/go/algorithms/numbertheory"

// Primality testing
isPrime := numbertheory.MillerRabinPrimality(982451653, 10)

// Integer factorization
factor := numbertheory.PollardRho(1403) // Returns 23 or 61

// Chinese Remainder Theorem
remainders := []int64{2, 3, 2}
moduli := []int64{3, 5, 7}
solution, valid := numbertheory.ChineseRemainderTheorem(remainders, moduli)
```

## 🧪 **Testing & Quality Assurance**

### Run All Tests
```bash
# Run all tests across all packages
go test ./...

# Run tests with coverage report
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Test Specific Categories
```bash
# Test individual algorithm categories
go test ./algorithms/graphs -v
go test ./algorithms/backtracking -v
go test ./algorithms/greedy -v
go test ./algorithms/geometry -v
go test ./algorithms/numbertheory -v
go test ./algorithms/dp -v
go test ./algorithms/strings -v
```

### Performance Benchmarks
```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific category benchmarks
go test -bench=. ./algorithms/sorting
go test -bench=. ./algorithms/graphs

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem ./...
```

## 📊 **Algorithm Complexity Reference**

### **Sorting Algorithms**
| Algorithm | Best Case | Average Case | Worst Case | Space |
|-----------|-----------|--------------|------------|-------|
| Quick Sort | O(n log n) | O(n log n) | O(n²) | O(log n) |
| Merge Sort | O(n log n) | O(n log n) | O(n log n) | O(n) |
| Heap Sort | O(n log n) | O(n log n) | O(n log n) | O(1) |
| Tim Sort | O(n) | O(n log n) | O(n log n) | O(n) |
| Counting Sort | O(n + k) | O(n + k) | O(n + k) | O(k) |

### **Graph Algorithms**
| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Dijkstra's | O((V + E) log V) | O(V) |
| Bellman-Ford | O(VE) | O(V) |
| Floyd-Warshall | O(V³) | O(V²) |
| Kruskal's MST | O(E log E) | O(V) |
| Prim's MST | O((V + E) log V) | O(V) |
| Tarjan's SCC | O(V + E) | O(V) |
| Ford-Fulkerson | O(E × max_flow) | O(V²) |

### **Dynamic Programming**
| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Knapsack 0/1 | O(nW) | O(nW) |
| LCS | O(mn) | O(mn) |
| Edit Distance | O(mn) | O(mn) |
| Matrix Chain | O(n³) | O(n²) |
| Egg Dropping | O(nk²) | O(nk) |

### **String Algorithms**
| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| KMP Search | O(n + m) | O(m) |
| Rabin-Karp | O(n + m) avg | O(1) |
| Boyer-Moore | O(n/m) best, O(nm) worst | O(σ) |
| Aho-Corasick | O(n + m + z) | O(m) |
| Suffix Array | O(n log n) | O(n) |

### **Geometric Algorithms**
| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Convex Hull | O(n log n) | O(n) |
| Closest Pair | O(n log n) | O(n) |
| Point in Polygon | O(n) | O(1) |
| Line Intersection | O(1) | O(1) |

### **Number Theory**
| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Miller-Rabin | O(k log³ n) | O(1) |
| Pollard's Rho | O(n^(1/4)) | O(1) |
| FFT | O(n log n) | O(n) |
| Discrete Log | O(√p) | O(√p) |

## 🎯 **Use Cases**

### **Educational**
- 🎓 **Computer Science Students**: Learn algorithm implementation and analysis
- 📚 **Interview Preparation**: Practice coding problems and algorithm design
- 🔬 **Research**: Reference implementations for algorithm comparison

### **Professional Development**
- 💼 **Backend Development**: Graph algorithms for social networks, routing
- 🔐 **Cryptography**: Number theory algorithms for security applications
- 🎮 **Game Development**: Pathfinding, collision detection, procedural generation
- 📊 **Data Science**: String processing, geometric computations, optimization

### **Production Applications**
- 🌐 **Web Services**: Efficient data processing and optimization
- 🤖 **Machine Learning**: Feature engineering and data preprocessing
- 📱 **Mobile Apps**: Local data processing and algorithms
- 🏗️ **System Design**: Building scalable, efficient systems

## 🤝 **Contributing**

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### **Development Setup**
```bash
# Clone and setup
git clone https://github.com/DimaJoyti/Algorithms.git
cd Algorithms/go

# Run tests to ensure everything works
go test ./...

# Run the demo
go run main.go
```

### **Adding New Algorithms**
1. Choose appropriate package (or create new one)
2. Implement algorithm with proper documentation
3. Add comprehensive tests with edge cases
4. Add benchmarks for performance measurement
5. Update main.go with demonstration
6. Update README.md with algorithm description

## 📝 **License**

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🔗 **Related Projects**

- [JavaScript Implementation](../javascript/)
- [Python Implementation](../python/)
- [Algorithm Visualizations](../visualizations/)

## 🏆 **Acknowledgments**

This collection implements algorithms from various sources including:
- Classic computer science textbooks
- Research papers and academic publications
- Competitive programming resources
- Open source algorithm libraries

---

**⭐ Star this repository if you find it helpful!**

**🐛 Found a bug or want to suggest an improvement? [Open an issue](https://github.com/DimaJoyti/Algorithms/issues)**
