# Advanced Pattern Templates

> Extended patterns for medium-hard problems

---

## 1. Trie (Prefix Tree)

```python
class TrieNode:
    def __init__(self):
        self.children = {}
        self.is_end = False

class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
        node.is_end = True

    def search(self, word: str) -> bool:
        node = self._find(word)
        return node is not None and node.is_end

    def startsWith(self, prefix: str) -> bool:
        return self._find(prefix) is not None

    def _find(self, prefix: str) -> TrieNode:
        node = self.root
        for char in prefix:
            if char not in node.children:
                return None
            node = node.children[char]
        return node

    # With wildcard support (e.g., "a.le" matches "apple")
    def searchWildcard(self, word: str) -> bool:
        def dfs(node, i):
            if i == len(word):
                return node.is_end
            char = word[i]
            if char == '.':
                return any(dfs(child, i + 1) for child in node.children.values())
            return char in node.children and dfs(node.children[char], i + 1)
        return dfs(self.root, 0)
```

**Use When:**
- Prefix matching
- Autocomplete
- Word dictionary
- "Starts with" queries
- Pattern matching with wildcards

**Time:** O(m) per operation (m = word length)
**Space:** O(ALPHABET_SIZE × N × M) worst case

---

## 2. Graph - Dijkstra's Algorithm

```python
import heapq
from collections import defaultdict

def dijkstra(n: int, edges: list, start: int) -> dict:
    """
    Find shortest path from start to all nodes.
    edges: [(u, v, weight), ...]
    """
    # Build adjacency list
    graph = defaultdict(list)
    for u, v, w in edges:
        graph[u].append((v, w))
        graph[v].append((u, w))  # Remove if directed

    # Initialize distances
    dist = {i: float('inf') for i in range(n)}
    dist[start] = 0

    # Priority queue: (distance, node)
    pq = [(0, start)]

    while pq:
        d, node = heapq.heappop(pq)

        # Skip if we've found a better path
        if d > dist[node]:
            continue

        for neighbor, weight in graph[node]:
            new_dist = dist[node] + weight
            if new_dist < dist[neighbor]:
                dist[neighbor] = new_dist
                heapq.heappush(pq, (new_dist, neighbor))

    return dist

# Variation: With at most K stops
def dijkstra_k_stops(n, flights, src, dst, k):
    graph = defaultdict(list)
    for u, v, w in flights:
        graph[u].append((v, w))

    # (cost, node, stops)
    pq = [(0, src, 0)]
    visited = {}  # node -> min stops to reach

    while pq:
        cost, node, stops = heapq.heappop(pq)

        if node == dst:
            return cost

        if stops > k:
            continue

        if node in visited and visited[node] < stops:
            continue
        visited[node] = stops

        for neighbor, price in graph[node]:
            heapq.heappush(pq, (cost + price, neighbor, stops + 1))

    return -1
```

**Use When:**
- Shortest path (weighted graph)
- Minimum cost to reach target
- Network routing

**Time:** O((V + E) log V)
**Space:** O(V)

---

## 3. Graph - Bellman-Ford

```python
def bellman_ford(n: int, edges: list, start: int) -> list:
    """
    Find shortest paths (handles negative weights).
    Detects negative cycles.
    """
    dist = [float('inf')] * n
    dist[start] = 0

    # Relax all edges n-1 times
    for _ in range(n - 1):
        for u, v, w in edges:
            if dist[u] != float('inf') and dist[u] + w < dist[v]:
                dist[v] = dist[u] + w

    # Check for negative cycle
    for u, v, w in edges:
        if dist[u] != float('inf') and dist[u] + w < dist[v]:
            return None  # Negative cycle detected

    return dist
```

**Use When:**
- Negative edge weights
- Detect negative cycles
- "Cheapest flights within K stops"

**Time:** O(V × E)
**Space:** O(V)

---

## 4. Graph - Topological Sort

```python
from collections import deque, defaultdict

def topological_sort(n: int, edges: list) -> list:
    """
    Return topological order or [] if cycle exists.
    edges: [(u, v), ...] means u must come before v
    """
    graph = defaultdict(list)
    in_degree = [0] * n

    for u, v in edges:
        graph[u].append(v)
        in_degree[v] += 1

    # Start with nodes having no dependencies
    queue = deque([i for i in range(n) if in_degree[i] == 0])
    result = []

    while queue:
        node = queue.popleft()
        result.append(node)

        for neighbor in graph[node]:
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)

    return result if len(result) == n else []

# Alien Dictionary variant
def alien_order(words: list) -> str:
    """
    Given sorted words in alien language, return character order.
    """
    # Build graph from adjacent word comparisons
    chars = set(''.join(words))
    graph = defaultdict(set)
    in_degree = {c: 0 for c in chars}

    for i in range(len(words) - 1):
        w1, w2 = words[i], words[i + 1]
        min_len = min(len(w1), len(w2))

        # Invalid: "abc" before "ab"
        if len(w1) > len(w2) and w1[:min_len] == w2[:min_len]:
            return ""

        for j in range(min_len):
            if w1[j] != w2[j]:
                if w2[j] not in graph[w1[j]]:
                    graph[w1[j]].add(w2[j])
                    in_degree[w2[j]] += 1
                break

    # Topological sort
    queue = deque([c for c in chars if in_degree[c] == 0])
    result = []

    while queue:
        char = queue.popleft()
        result.append(char)
        for neighbor in graph[char]:
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)

    return ''.join(result) if len(result) == len(chars) else ""
```

**Use When:**
- Course schedule / dependencies
- Build order
- Alien dictionary
- Any "must come before" relationships

**Time:** O(V + E)
**Space:** O(V + E)

---

## 5. Graph - Prim's MST

```python
import heapq
from collections import defaultdict

def prim_mst(n: int, edges: list) -> int:
    """
    Find minimum spanning tree cost.
    edges: [(u, v, weight), ...]
    """
    graph = defaultdict(list)
    for u, v, w in edges:
        graph[u].append((v, w))
        graph[v].append((u, w))

    visited = set()
    min_heap = [(0, 0)]  # (weight, node)
    total_cost = 0

    while min_heap and len(visited) < n:
        weight, node = heapq.heappop(min_heap)

        if node in visited:
            continue

        visited.add(node)
        total_cost += weight

        for neighbor, w in graph[node]:
            if neighbor not in visited:
                heapq.heappush(min_heap, (w, neighbor))

    return total_cost if len(visited) == n else -1
```

**Use When:**
- Connect all nodes with minimum cost
- Network design
- Minimum spanning tree

**Time:** O(E log V)
**Space:** O(V + E)

---

## 6. Graph - Kruskal's MST (Union Find)

```python
class UnionFind:
    def __init__(self, n):
        self.parent = list(range(n))
        self.rank = [0] * n

    def find(self, x):
        if self.parent[x] != x:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        px, py = self.find(x), self.find(y)
        if px == py:
            return False
        if self.rank[px] < self.rank[py]:
            px, py = py, px
        self.parent[py] = px
        self.rank[px] += self.rank[px] == self.rank[py]
        return True

def kruskal_mst(n: int, edges: list) -> int:
    """
    Find minimum spanning tree cost using Union Find.
    """
    edges.sort(key=lambda x: x[2])  # Sort by weight
    uf = UnionFind(n)
    total_cost = 0
    edges_used = 0

    for u, v, w in edges:
        if uf.union(u, v):
            total_cost += w
            edges_used += 1
            if edges_used == n - 1:
                break

    return total_cost if edges_used == n - 1 else -1
```

**Use When:**
- MST with Union Find
- Detect cycles in undirected graph
- Connected components

**Time:** O(E log E) for sorting
**Space:** O(V)

---

## 7. Advanced DP - Kadane's Variants

```python
# Maximum Subarray Sum
def max_subarray(nums):
    max_sum = curr = nums[0]
    for num in nums[1:]:
        curr = max(num, curr + num)
        max_sum = max(max_sum, curr)
    return max_sum

# Maximum Circular Subarray Sum
def max_circular_subarray(nums):
    def kadane(arr):
        max_sum = curr = arr[0]
        for num in arr[1:]:
            curr = max(num, curr + num)
            max_sum = max(max_sum, curr)
        return max_sum

    # Case 1: Non-circular (standard Kadane)
    max_normal = kadane(nums)

    # Case 2: Circular (total - minimum subarray)
    total = sum(nums)
    inverted = [-x for x in nums]
    min_subarray = -kadane(inverted)

    if min_subarray == total:  # All negative
        return max_normal

    return max(max_normal, total - min_subarray)

# Maximum Product Subarray
def max_product_subarray(nums):
    max_prod = min_prod = result = nums[0]

    for num in nums[1:]:
        if num < 0:
            max_prod, min_prod = min_prod, max_prod

        max_prod = max(num, max_prod * num)
        min_prod = min(num, min_prod * num)

        result = max(result, max_prod)

    return result
```

---

## 8. Advanced DP - LCS Variants

```python
# Longest Common Subsequence
def lcs(text1: str, text2: str) -> int:
    m, n = len(text1), len(text2)
    dp = [[0] * (n + 1) for _ in range(m + 1)]

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if text1[i-1] == text2[j-1]:
                dp[i][j] = dp[i-1][j-1] + 1
            else:
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])

    return dp[m][n]

# Space-optimized LCS
def lcs_optimized(text1: str, text2: str) -> int:
    if len(text1) < len(text2):
        text1, text2 = text2, text1

    prev = [0] * (len(text2) + 1)

    for c1 in text1:
        curr = [0] * (len(text2) + 1)
        for j, c2 in enumerate(text2, 1):
            if c1 == c2:
                curr[j] = prev[j-1] + 1
            else:
                curr[j] = max(prev[j], curr[j-1])
        prev = curr

    return prev[-1]

# Longest Palindromic Subsequence
def longest_palindrome_subseq(s: str) -> int:
    return lcs(s, s[::-1])

# Edit Distance
def edit_distance(word1: str, word2: str) -> int:
    m, n = len(word1), len(word2)
    dp = [[0] * (n + 1) for _ in range(m + 1)]

    # Base cases
    for i in range(m + 1):
        dp[i][0] = i
    for j in range(n + 1):
        dp[0][j] = j

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if word1[i-1] == word2[j-1]:
                dp[i][j] = dp[i-1][j-1]
            else:
                dp[i][j] = 1 + min(
                    dp[i-1][j],    # delete
                    dp[i][j-1],    # insert
                    dp[i-1][j-1]   # replace
                )

    return dp[m][n]
```

---

## 9. Advanced DP - Interval DP

```python
# Burst Balloons
def burst_balloons(nums: list) -> int:
    nums = [1] + nums + [1]
    n = len(nums)
    dp = [[0] * n for _ in range(n)]

    for length in range(2, n):
        for left in range(n - length):
            right = left + length
            for k in range(left + 1, right):
                dp[left][right] = max(
                    dp[left][right],
                    nums[left] * nums[k] * nums[right] + dp[left][k] + dp[k][right]
                )

    return dp[0][n-1]

# Matrix Chain Multiplication
def matrix_chain_mult(dimensions: list) -> int:
    n = len(dimensions) - 1
    dp = [[0] * n for _ in range(n)]

    for length in range(2, n + 1):
        for i in range(n - length + 1):
            j = i + length - 1
            dp[i][j] = float('inf')
            for k in range(i, j):
                cost = dp[i][k] + dp[k+1][j] + dimensions[i] * dimensions[k+1] * dimensions[j+1]
                dp[i][j] = min(dp[i][j], cost)

    return dp[0][n-1]

# Predict the Winner (Optimal Game Strategy)
def predict_winner(nums: list) -> bool:
    n = len(nums)
    dp = [[0] * n for _ in range(n)]

    for i in range(n):
        dp[i][i] = nums[i]

    for length in range(2, n + 1):
        for i in range(n - length + 1):
            j = i + length - 1
            dp[i][j] = max(nums[i] - dp[i+1][j], nums[j] - dp[i][j-1])

    return dp[0][n-1] >= 0
```

---

## 10. Advanced DP - Digit DP

```python
# Count numbers with digit sum in range [low, high]
def count_numbers_with_digit_sum(n: str, target: int) -> int:
    """
    Count numbers from 0 to n where digit sum equals target.
    """
    from functools import lru_cache

    @lru_cache(None)
    def dp(pos: int, tight: bool, curr_sum: int) -> int:
        if pos == len(n):
            return 1 if curr_sum == target else 0

        limit = int(n[pos]) if tight else 9
        count = 0

        for digit in range(limit + 1):
            new_tight = tight and digit == limit
            count += dp(pos + 1, new_tight, curr_sum + digit)

        return count

    return dp(0, True, 0)
```

---

## 11. Segment Tree

```python
class SegmentTree:
    def __init__(self, nums: list):
        self.n = len(nums)
        self.tree = [0] * (4 * self.n)
        self._build(nums, 0, 0, self.n - 1)

    def _build(self, nums, node, start, end):
        if start == end:
            self.tree[node] = nums[start]
        else:
            mid = (start + end) // 2
            self._build(nums, 2*node + 1, start, mid)
            self._build(nums, 2*node + 2, mid + 1, end)
            self.tree[node] = self.tree[2*node + 1] + self.tree[2*node + 2]

    def update(self, idx: int, val: int):
        self._update(0, 0, self.n - 1, idx, val)

    def _update(self, node, start, end, idx, val):
        if start == end:
            self.tree[node] = val
        else:
            mid = (start + end) // 2
            if idx <= mid:
                self._update(2*node + 1, start, mid, idx, val)
            else:
                self._update(2*node + 2, mid + 1, end, idx, val)
            self.tree[node] = self.tree[2*node + 1] + self.tree[2*node + 2]

    def query(self, left: int, right: int) -> int:
        return self._query(0, 0, self.n - 1, left, right)

    def _query(self, node, start, end, left, right):
        if right < start or left > end:
            return 0
        if left <= start and end <= right:
            return self.tree[node]
        mid = (start + end) // 2
        return (self._query(2*node + 1, start, mid, left, right) +
                self._query(2*node + 2, mid + 1, end, left, right))
```

**Use When:**
- Range queries (sum, min, max)
- Point updates
- O(log n) operations needed

---

## 12. Fenwick Tree (Binary Indexed Tree)

```python
class FenwickTree:
    def __init__(self, n: int):
        self.n = n
        self.tree = [0] * (n + 1)

    def update(self, i: int, delta: int):
        while i <= self.n:
            self.tree[i] += delta
            i += i & (-i)

    def query(self, i: int) -> int:
        total = 0
        while i > 0:
            total += self.tree[i]
            i -= i & (-i)
        return total

    def range_query(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)
```

**Use When:**
- Prefix sum queries
- Point updates
- More memory efficient than segment tree

---

## Quick Reference: When to Use

| Problem Type | Pattern |
|--------------|---------|
| Prefix matching | Trie |
| Shortest path (positive weights) | Dijkstra |
| Shortest path (negative weights) | Bellman-Ford |
| Dependencies / Order | Topological Sort |
| Connect all nodes min cost | Prim's / Kruskal's |
| Maximum subarray | Kadane's |
| String comparison | LCS / Edit Distance |
| Game theory / Optimal play | Interval DP |
| Range queries | Segment Tree / Fenwick |
