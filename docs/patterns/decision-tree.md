# Pattern Decision Tree

> Quick reference for selecting the right algorithm pattern

---

## Main Decision Flow

```
START: Read the problem
    │
    ▼
┌─────────────────────────────────────────────────────────────┐
│ Q1: What is the INPUT TYPE?                                  │
└─────────────────────────────────────────────────────────────┘
    │
    ├─── ARRAY / STRING
    │        │
    │        ├─── Sorted? ─────────────────────────────────────┐
    │        │       │                                         │
    │        │       ├── YES → Binary Search / Two Pointers    │
    │        │       │                                         │
    │        │       └── NO → Continue below                   │
    │        │                                                 │
    │        ├─── Find subarray/substring? ────────────────────┤
    │        │       │                                         │
    │        │       ├── Fixed size → Sliding Window (fixed)   │
    │        │       ├── Variable size → Sliding Window (var)  │
    │        │       └── With constraint → Hash + Window       │
    │        │                                                 │
    │        ├─── Find pairs/count? ───────────────────────────┤
    │        │       │                                         │
    │        │       ├── Exists? → HashSet                     │
    │        │       ├── Count all? → HashMap                  │
    │        │       └── In sorted → Two Pointers              │
    │        │                                                 │
    │        ├─── Need prefix/suffix? ─────────────────────────┤
    │        │       │                                         │
    │        │       └── Prefix Sum / Product                  │
    │        │                                                 │
    │        └─── Next greater/smaller? ───────────────────────┤
    │                │                                         │
    │                └── Monotonic Stack                       │
    │                                                          │
    ├─── TREE
    │        │
    │        ├─── Need depth/path? ──────→ DFS                 │
    │        ├─── Need level order? ────→ BFS                  │
    │        ├─── BST search? ──────────→ BST Properties       │
    │        ├─── Serialize? ───────────→ DFS with markers     │
    │        └─── LCA? ─────────────────→ DFS + Backtracking   │
    │                                                          │
    ├─── GRAPH
    │        │
    │        ├─── Connected components? ─→ DFS / Union Find    │
    │        ├─── Shortest path?                                      │
    │        │       ├── Unweighted → BFS                           │
    │        │       ├── Positive weights → Dijkstra                │
    │        │       └── Negative weights → Bellman-Ford            │
    │        ├─── Dependencies? ──────→ Topological Sort            │
    │        ├─── MST? ───────────────→ Prim's / Kruskal's          │
    │        └─── Cycle detection? ───→ DFS / Union Find            │
    │                                                          │
    ├─── LINKED LIST
    │        │
    │        ├─── Cycle? ────────────────→ Floyd's (Fast/Slow)      │
    │        ├─── Reverse? ──────────────→ Iterative               │
    │        ├─── Merge? ────────────────→ Merge Pattern           │
    │        └─── Find nth from end? ────→ Two Pointers            │
    │                                                          │
    ├─── MATRIX / 2D GRID
    │        │
    │        ├─── Find path? ────────────→ DFS / BFS              │
    │        ├─── Count islands? ────────→ DFS / Union Find       │
    │        ├─── Shortest path? ────────→ BFS                    │
    │        └─── Dynamic programming? ──→ 2D DP                   │
    │                                                          │
    └─── NUMBER / MATH
             │
             ├─── Count ways? ────────────→ DP / Combinatorics    │
             ├─── Find optimal? ──────────→ DP / Greedy           │
             ├─── Bit operations? ────────→ Bit Manipulation      │
             └─── Very large n? ──────────→ Math / Logarithmic    │
```

---

## Secondary Questions

### Q2: What is the GOAL?

```
┌─────────────────────────────────────────────────────────────┐
│ GOAL: What do you need to find/do?                           │
└─────────────────────────────────────────────────────────────┘
    │
    ├─── Find ONE element ─────────────→ Binary Search
    ├─── Find ALL elements ────────────→ DFS / Backtracking
    ├─── Find MAXIMUM/MINIMUM ─────────→ Greedy / DP / Heap
    ├─── Find SHORTEST/LONGEST ────────→ BFS / DP / Sliding Window
    ├─── Count ways ───────────────────→ DP
    ├─── Check if possible ────────────→ Greedy / DFS / DP
    ├─── All permutations ─────────────→ Backtracking
    ├─── All combinations ─────────────→ Backtracking / DP
    └─── Optimize value ───────────────→ DP / Greedy / Binary Search
```

### Q3: What are the CONSTRAINTS?

```
┌─────────────────────────────────────────────────────────────┐
│ CONSTRAINTS: Time and space requirements                     │
└─────────────────────────────────────────────────────────────┘
    │
    ├─── O(log n) required ────────────→ Binary Search
    ├─── O(n) required ────────────────→ Two Pointers / Sliding Window
    ├─── O(n log n) allowed ───────────→ Sorting + Greedy / Heap
    ├─── O(1) space required ──────────→ Two Pointers / In-place
    ├─── O(n) space allowed ───────────→ HashMap / DP
    └─── Large input (10^6+) ──────────→ O(n) or better required
```

---

## Pattern Combinations

### Common Combos

```
SORTED + Find target = Binary Search
SORTED + Find pairs = Two Pointers
SORTED + Rotated = Modified Binary Search

CONTIGUOUS + Fixed size = Sliding Window (fixed)
CONTIGUOUS + Variable size = Sliding Window (variable)
CONTIGUOUS + Sum/Product = Prefix Sum/Product

TREE + Depth = DFS
TREE + Level = BFS
TREE + Path Sum = DFS + Backtracking

GRAPH + Shortest (unweighted) = BFS
GRAPH + Shortest (weighted) = Dijkstra
GRAPH + Dependencies = Topological Sort
GRAPH + Connected = DFS / Union Find

STRING + Prefix = Trie
STRING + Pattern = KMP / Rolling Hash
STRING + Parentheses = Stack

OPTIMIZATION + Overlapping subproblems = DP
OPTIMIZATION + All combinations = Backtracking
OPTIMIZATION + Greedy choice property = Greedy
OPTIMIZATION + Priority = Heap

TOP K = Heap (min-heap for largest, max-heap for smallest)
MEDIAN = Two Heaps
LRU/LFU = HashMap + LinkedList
```

---

## Difficulty Indicators

### Easy Problems

```
✓ Single pass through array
✓ Simple HashMap usage
✓ Basic recursion
✓ Standard pattern application
```

### Medium Problems

```
✓ Combine 2+ patterns
✓ Modified standard pattern
✓ Careful edge case handling
✓ Space optimization needed
```

### Hard Problems

```
✓ Multiple patterns combined
✓ Custom data structure
✓ Complex state management
✓ Advanced optimization
✓ Mathematical insight
```

---

## Quick Pattern Codes

| Code | Pattern | When to Use |
|------|---------|-------------|
| `2P` | Two Pointers | Sorted array, pairs, palindromes |
| `SW` | Sliding Window | Contiguous subarray, substring |
| `BS` | Binary Search | Sorted data, O(log n) needed |
| `DFS` | Depth-First Search | Tree depth, graph exploration |
| `BFS` | Breadth-First Search | Level order, shortest path |
| `DP1` | 1D Dynamic Programming | Linear optimization |
| `DP2` | 2D Dynamic Programming | Grid, string comparison |
| `BT` | Backtracking | All combinations, permutations |
| `HP` | Heap | Top K, priority-based |
| `UF` | Union Find | Connected components |
| `MS` | Monotonic Stack | Next greater, histogram |
| `TR` | Trie | Prefix matching |

---

## Decision Tree Summary

```
1. INPUT TYPE → Narrows to 2-3 patterns
2. GOAL → Further narrows to 1-2 patterns
3. CONSTRAINTS → Confirms pattern choice
4. COMBINE → Apply secondary patterns if needed
```

---

## Common Mistakes

| Problem Type | Wrong Choice | Right Choice |
|--------------|--------------|--------------|
| Sorted array | Linear scan | Binary Search |
| Contiguous subarray | Brute force O(n²) | Sliding Window O(n) |
| Find pairs | Nested loops | Two Pointers / HashMap |
| Shortest path | DFS | BFS (unweighted) |
| All combinations | Iteration | Backtracking |
| Top K elements | Sort all | Heap of size K |
| Prefix queries | Linear scan | Trie |

---

*Use this decision tree at the start of every problem to quickly identify the right approach.*
