# 🔥 Pattern Quick Reference

> Essential patterns and when to use them

---

## 1. Arrays & Hashing

**Use when:** Looking for pairs, duplicates, frequency counting

**Key insight:** O(1) lookups with HashMap/HashSet

```python
# Frequency counter
freq = {}
for num in arr:
    freq[num] = freq.get(num, 0) + 1

# Two sum with hash
seen = {}
for i, num in enumerate(arr):
    if target - num in seen:
        return [seen[target - num], i]
    seen[num] = i
```

**Time:** O(n) | **Space:** O(n)

---

## 2. Two Pointers

**Use when:** Sorted array, pairs, palindromes

**Key insight:** Eliminate impossible cases

```python
# Sorted array pair sum
l, r = 0, len(arr) - 1
while l < r:
    s = arr[l] + arr[r]
    if s == target:
        return [l, r]
    elif s < target:
        l += 1
    else:
        r -= 1
```

**Time:** O(n) | **Space:** O(1)

---

## 3. Sliding Window

**Use when:** Contiguous subarray, substring, fixed/variable window

**Key insight:** Expand right, shrink left when invalid

```python
# Variable window
l = 0
for r in range(len(arr)):
    # Add arr[r] to window
    while invalid_condition:
        # Remove arr[l] from window
        l += 1
    # Update result
```

**Time:** O(n) | **Space:** O(k) or O(1)

---

## 4. Stack

**Use when:** Parentheses, next greater/smaller, monotonic

**Key insight:** LIFO for matching/reversing

```python
# Monotonic stack (decreasing)
stack = []
for num in arr:
    while stack and stack[-1] < num:
        stack.pop()
    stack.append(num)
```

**Time:** O(n) | **Space:** O(n)

---

## 5. Binary Search

**Use when:** Sorted array, find target/insert position

**Key insight:** Eliminate half each step

```python
l, r = 0, len(arr) - 1
while l <= r:
    mid = (l + r) // 2
    if arr[mid] == target:
        return mid
    elif arr[mid] < target:
        l = mid + 1
    else:
        r = mid - 1
```

**Time:** O(log n) | **Space:** O(1)

---

## 6. Linked List

**Use when:** In-place reversal, cycle detection, merge

**Key insight:** Slow/fast pointers, dummy head

```python
# Floyd's cycle detection
slow = fast = head
while fast and fast.next:
    slow = slow.next
    fast = fast.next.next
    if slow == fast:
        return True  # cycle
```

**Time:** O(n) | **Space:** O(1)

---

## 7. Trees - DFS

**Use when:** Path problems, depth, subtree

**Key insight:** Recursive or iterative with stack

```python
def dfs(node):
    if not node:
        return 0
    left = dfs(node.left)
    right = dfs(node.right)
    return max(left, right) + 1
```

**Time:** O(n) | **Space:** O(h) where h = height

---

## 8. Trees - BFS

**Use when:** Level order, shortest path in tree

**Key insight:** Queue for level-by-level

```python
from collections import deque
q = deque([root])
while q:
    level = []
    for _ in range(len(q)):
        node = q.popleft()
        level.append(node.val)
        if node.left:
            q.append(node.left)
        if node.right:
            q.append(node.right)
```

**Time:** O(n) | **Space:** O(w) where w = max width

---

## 9. Heap / Priority Queue

**Use when:** Top K, merge K sorted, median

**Key insight:** Min heap for largest K, max heap for smallest K

```python
import heapq

# K largest: min heap of size K
heap = []
for num in arr:
    heapq.heappush(heap, num)
    if len(heap) > k:
        heapq.heappop(heap)
```

**Time:** O(n log k) | **Space:** O(k)

---

## 10. Backtracking

**Use when:** All permutations/combinations, subsets

**Key insight:** Try all options, backtrack on failure

```python
def backtrack(path, options):
    if valid(path):
        result.append(path[:])
        return
    for i, option in enumerate(options):
        path.append(option)
        backtrack(path, options[i+1:])
        path.pop()
```

**Time:** O(n!) or O(2^n) | **Space:** O(n)

---

## 11. Graphs - DFS

**Use when:** Count connected components, cycle detection

**Key insight:** Mark visited, explore neighbors

```python
def dfs(node, visited, graph):
    visited.add(node)
    for neighbor in graph[node]:
        if neighbor not in visited:
            dfs(neighbor, visited, graph)
```

**Time:** O(V + E) | **Space:** O(V)

---

## 12. Graphs - BFS

**Use when:** Shortest path in unweighted graph

**Key insight:** Level by level exploration

```python
from collections import deque
visited = set()
q = deque([start])
while q:
    node = q.popleft()
    for neighbor in graph[node]:
        if neighbor not in visited:
            visited.add(neighbor)
            q.append(neighbor)
```

**Time:** O(V + E) | **Space:** O(V)

---

## 13. Union Find

**Use when:** Connected components, cycle detection

**Key insight:** Union by rank + path compression

```python
parent = list(range(n))
rank = [0] * n

def find(x):
    if parent[x] != x:
        parent[x] = find(parent[x])  # path compression
    return parent[x]

def union(x, y):
    px, py = find(x), find(y)
    if px == py:
        return False
    if rank[px] < rank[py]:
        px, py = py, px
    parent[py] = px
    rank[px] += rank[px] == rank[py]
    return True
```

**Time:** O(α(n)) ≈ O(1) | **Space:** O(n)

---

## 14. Topological Sort

**Use when:** Course schedule, build order, dependencies

**Key insight:** Process nodes with 0 in-degree

```python
from collections import deque, defaultdict

in_degree = defaultdict(int)
graph = defaultdict(list)
# Build graph and in_degree...

q = deque([n for n in range(n) if in_degree[n] == 0])
result = []
while q:
    node = q.popleft()
    result.append(node)
    for neighbor in graph[node]:
        in_degree[neighbor] -= 1
        if in_degree[neighbor] == 0:
            q.append(neighbor)
```

**Time:** O(V + E) | **Space:** O(V)

---

## 15. Trie

**Use when:** Word search, autocomplete, prefix matching

**Key insight:** Character-based tree

```python
class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_end = False

class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word):
        node = self.root
        for c in word:
            if c not in node.children:
                node.children[c] = TrieNode()
            node = node.children[c]
        node.is_end = True
```

**Time:** O(m) for m-length word | **Space:** O(26^m)

---

## 16. Dynamic Programming - 1D

**Use when:** Overlapping subproblems, optimal substructure

**Key insight:** Cache results, build from base cases

```python
# Bottom-up
dp = [0] * (n + 1)
dp[0], dp[1] = 0, 1
for i in range(2, n + 1):
    dp[i] = dp[i-1] + dp[i-2]
return dp[n]

# Space optimized
a, b = 0, 1
for _ in range(n):
    a, b = b, a + b
return a
```

**Time:** O(n) | **Space:** O(n) or O(1)

---

## 17. Dynamic Programming - 2D

**Use when:** Grid problems, string comparison, knapsack

**Key insight:** dp[i][j] from dp[i-1][j], dp[i][j-1], dp[i-1][j-1]

```python
# LCS
dp = [[0] * (n+1) for _ in range(m+1)]
for i in range(1, m+1):
    for j in range(1, n+1):
        if s1[i-1] == s2[j-1]:
            dp[i][j] = dp[i-1][j-1] + 1
        else:
            dp[i][j] = max(dp[i-1][j], dp[i][j-1])
```

**Time:** O(mn) | **Space:** O(mn)

---

## 18. Greedy

**Use when:** Local optimal leads to global optimal

**Key insight:** Sort + make best choice at each step

```python
# Jump Game II
jumps = current_end = farthest = 0
for i in range(len(nums) - 1):
    farthest = max(farthest, i + nums[i])
    if i == current_end:
        jumps += 1
        current_end = farthest
```

**Time:** O(n) or O(n log n) | **Space:** O(1)

---

## Quick Decision Tree

```
Problem Type?
├── Sorted array → Binary Search / Two Pointers
├── Subarray/Substring → Sliding Window
├── Pairs/Sum → HashMap / Two Pointers
├── Tree → DFS / BFS
├── Graph → DFS / BFS / Union Find
├── Shortest Path → BFS (unweighted) / Dijkstra (weighted)
├── All combinations → Backtracking
├── Overlapping subproblems → DP
├── Top K / Median → Heap
├── Prefix matching → Trie
└── Parentheses/Matching → Stack
```

---

## Complexity Quick Reference

| Data Structure | Access | Search | Insert | Delete |
|---------------|--------|--------|--------|--------|
| Array | O(1) | O(n) | O(n) | O(n) |
| HashMap | O(1)* | O(1)* | O(1)* | O(1)* |
| HashSet | - | O(1)* | O(1)* | O(1)* |
| Heap | O(1) peek | O(n) | O(log n) | O(log n) |
| BST | O(log n)* | O(log n)* | O(log n)* | O(log n)* |
| Trie | O(m) | O(m) | O(m) | O(m) |

*Average case

---

## Common Time Complexities

| Complexity | n=10 | n=100 | n=1000 | n=10^6 |
|-----------|------|-------|--------|--------|
| O(1) | instant | instant | instant | instant |
| O(log n) | 3 | 7 | 10 | 20 |
| O(n) | 10 | 100 | 1000 | 10^6 |
| O(n log n) | 30 | 700 | 10^4 | 2×10^7 |
| O(n²) | 100 | 10^4 | 10^6 | 10^12 |
| O(2^n) | 1024 | 10^30 | ∞ | ∞ |

---

*Last Updated: YYYY-MM-DD*
