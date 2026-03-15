# Weekly Practice Problem Sets

> Structured 12-week practice plan with specific problems

---

## Week 1: Arrays & Hashing

### Focus
HashMap, HashSet, frequency counting, grouping

### Problems (15)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Two Sum | 🟢 | HashMap | #1 |
| 2 | Contains Duplicate | 🟢 | HashSet | #217 |
| 3 | Valid Anagram | 🟢 | HashMap/Sort | #242 |
| 4 | Missing Number | 🟢 | Math/XOR | #268 |
| 5 | Single Number | 🟢 | XOR | #136 |
| 6 | Happy Number | 🟢 | HashSet | #202 |
| 7 | Group Anagrams | 🟡 | HashMap | #49 |
| 8 | Top K Frequent | 🟡 | Bucket Sort | #347 |
| 9 | Product Except Self | 🟡 | Prefix/Suffix | #238 |
| 10 | Encode/Decode Strings | 🟡 | String Manip | #271 |
| 11 | Longest Consecutive | 🟡 | HashSet | #128 |
| 12 | 4Sum II | 🟡 | HashMap | #454 |
| 13 | Subarray Sum Equals K | 🟡 | Prefix Sum | #560 |
| 14 | Continuous Subarray Sum | 🟡 | Prefix Sum | #523 |
| 15 | Count Nice Pairs | 🟡 | HashMap | #1814 |

### Daily Schedule
```
Day 1:  Problems 1-3 (45 min)
Day 2:  Problems 4-6 (45 min)
Day 3:  Problems 7-8 (45 min)
Day 4:  Problems 9-10 (45 min)
Day 5:  Problems 11-12 (45 min)
Day 6:  Problems 13-15 (45 min)
Day 7:  Review + 1 mock interview (90 min)
```

### Key Concepts
- HashMap: O(1) lookup, store computed values
- HashSet: Check existence, remove duplicates
- Prefix sum: Subarray sum = prefix[j] - prefix[i-1]

---

## Week 2: Two Pointers & Strings

### Focus
Sorted arrays, pairs, palindromes, string manipulation

### Problems (15)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Valid Palindrome | 🟢 | Two Pointers | #125 |
| 2 | Reverse String | 🟢 | Two Pointers | #344 |
| 3 | Squares of Sorted Array | 🟢 | Two Pointers | #977 |
| 4 | Move Zeroes | 🟢 | Two Pointers | #283 |
| 5 | Two Sum II | 🟢 | Two Pointers | #167 |
| 6 | 3Sum | 🟡 | Two Pointers | #15 |
| 7 | 3Sum Closest | 🟡 | Two Pointers | #16 |
| 8 | Container With Most Water | 🟡 | Two Pointers | #11 |
| 9 | Trapping Rain Water | 🔴 | Two Pointers | #42 |
| 10 | Longest Palindromic Substring | 🟡 | Expand Around Center | #5 |
| 11 | Palindromic Substrings | 🟡 | Expand Around Center | #647 |
| 12 | Valid Parentheses | 🟢 | Stack | #20 |
| 13 | Min Stack | 🟡 | Stack | #155 |
| 14 | Backspace String Compare | 🟢 | Two Pointers | #844 |
| 15 | Sort Colors | 🟡 | Dutch Flag | #75 |

### Key Concepts
- Two pointers: One from start, one from end
- Sorted array → always consider two pointers
- Skip duplicates in sorted arrays

---

## Week 3: Sliding Window

### Focus
Contiguous subarrays, substrings, fixed/variable windows

### Problems (12)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Maximum Average Subarray I | 🟢 | Fixed Window | #643 |
| 2 | Best Time to Buy/Sell Stock | 🟢 | Single Pass | #121 |
| 3 | Longest Substring Without Repeating | 🟡 | Variable Window | #3 |
| 4 | Longest Repeating Character Replacement | 🟡 | Variable Window | #424 |
| 5 | Permutation in String | 🟡 | Fixed Window | #567 |
| 6 | Find All Anagrams | 🟡 | Fixed Window | #438 |
| 7 | Minimum Window Substring | 🔴 | Variable Window | #76 |
| 8 | Sliding Window Maximum | 🔴 | Deque | #239 |
| 9 | Subarray Product Less Than K | 🟡 | Variable Window | #713 |
| 10 | Max Consecutive Ones III | 🟡 | Variable Window | #1004 |
| 11 | Longest Subarray with Sum K | 🟡 | Variable Window | - |
| 12 | Minimum Size Subarray Sum | 🟡 | Variable Window | #209 |

### Window Template
```python
l = 0
for r in range(len(arr)):
    # add arr[r] to window
    while invalid_condition():
        # remove arr[l] from window
        l += 1
    # update result
```

---

## Week 4: Binary Search

### Focus
Sorted data, search space reduction, finding boundaries

### Problems (12)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Binary Search | 🟢 | Standard | #704 |
| 2 | Search Insert Position | 🟢 | Standard | #35 |
| 3 | First Bad Version | 🟢 | Find First | #278 |
| 4 | Sqrt(x) | 🟢 | Search Space | #69 |
| 5 | Guess Number Higher/Lower | 🟢 | Standard | #374 |
| 6 | Search in Rotated Sorted | 🟡 | Modified BS | #33 |
| 7 | Find Minimum in Rotated Sorted | 🟡 | Modified BS | #153 |
| 8 | Search a 2D Matrix | 🟡 | 2D BS | #74 |
| 9 | Find Peak Element | 🟡 | Modified BS | #162 |
| 10 | Time Based Key-Value Store | 🟡 | Binary Search | #981 |
| 11 | Median of Two Sorted Arrays | 🔴 | Binary Search | #4 |
| 12 | Find Kth Smallest Pair Distance | 🔴 | Binary Search | #719 |

### BS Template (Find First True)
```python
l, r = 0, n  # Note: r = n, not n-1
while l < r:
    mid = l + (r - l) // 2
    if condition(mid):
        r = mid
    else:
        l = mid + 1
return l
```

---

## Week 5: Linked Lists

### Focus
Pointer manipulation, reversing, merging

### Problems (12)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Reverse Linked List | 🟢 | Iterative | #206 |
| 2 | Merge Two Sorted Lists | 🟢 | Merge | #21 |
| 3 | Linked List Cycle | 🟢 | Fast/Slow | #141 |
| 4 | Middle of Linked List | 🟢 | Fast/Slow | #876 |
| 5 | Remove Duplicates | 🟢 | Iteration | #83 |
| 6 | Remove Nth Node From End | 🟡 | Two Pointers | #19 |
| 7 | Reorder List | 🟡 | Fast/Slow + Reverse | #143 |
| 8 | Add Two Numbers | 🟡 | Iteration | #2 |
| 9 | Copy List with Random Pointer | 🟡 | HashMap | #138 |
| 10 | Find Duplicate Number | 🟡 | Floyd's Cycle | #287 |
| 11 | Merge K Sorted Lists | 🔴 | Heap/Divide Conquer | #23 |
| 12 | Reverse Nodes in K Group | 🔴 | Iteration | #25 |

### Key Techniques
- Fast/slow pointers for cycle detection
- Dummy node for edge cases
- Reverse: prev, curr, next pattern

---

## Week 6: Trees - Basics

### Focus
DFS, BFS, tree traversals

### Problems (15)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Maximum Depth | 🟢 | DFS | #104 |
| 2 | Invert Binary Tree | 🟢 | DFS | #226 |
| 3 | Same Tree | 🟢 | DFS | #100 |
| 4 | Symmetric Tree | 🟢 | DFS | #101 |
| 5 | Path Sum | 🟢 | DFS | #112 |
| 6 | Diameter of Binary Tree | 🟢 | DFS | #543 |
| 7 | Balanced Binary Tree | 🟢 | DFS | #110 |
| 8 | Subtree of Another Tree | 🟢 | DFS | #572 |
| 9 | Binary Tree Level Order | 🟡 | BFS | #102 |
| 10 | Average of Levels | 🟢 | BFS | #637 |
| 11 | Binary Tree Right Side View | 🟡 | BFS | #199 |
| 12 | Minimum Depth | 🟢 | BFS | #111 |
| 13 | Cousins in Binary Tree | 🟢 | BFS | #993 |
| 14 | N-ary Tree Level Order | 🟡 | BFS | #429 |
| 15 | Binary Tree Paths | 🟢 | DFS | #257 |

### Traversal Templates
```python
# DFS
def dfs(node):
    if not node: return
    # preorder: process here
    dfs(node.left)
    # inorder: process here
    dfs(node.right)
    # postorder: process here

# BFS
q = deque([root])
while q:
    for _ in range(len(q)):
        node = q.popleft()
        # process
```

---

## Week 7: Trees - Advanced & BST

### Focus
BST properties, LCA, serialization

### Problems (12)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Validate BST | 🟡 | DFS | #98 |
| 2 | Search in BST | 🟢 | BST | #700 |
| 3 | Insert into BST | 🟡 | BST | #701 |
| 4 | Delete Node in BST | 🟡 | BST | #450 |
| 5 | Kth Smallest Element | 🟡 | Inorder | #230 |
| 6 | Lowest Common Ancestor | 🟡 | DFS | #236 |
| 7 | LCA of BST | 🟢 | BST | #235 |
| 8 | Construct BST from Preorder | 🟡 | DFS | #1008 |
| 9 | Serialize/Deserialize BST | 🟡 | DFS | #449 |
| 10 | Convert Sorted Array to BST | 🟢 | Recursion | #108 |
| 11 | Binary Tree Maximum Path Sum | 🔴 | DFS | #124 |
| 12 | Count Good Nodes | 🟡 | DFS | #1448 |

### BST Property
- Left subtree < node < right subtree
- Inorder traversal gives sorted order

---

## Week 8: Heaps & Priority Queues

### Focus
Top K, merging, priority-based processing

### Problems (10)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Kth Largest Element | 🟡 | Min Heap | #215 |
| 2 | Top K Frequent | 🟡 | Bucket/Heap | #347 |
| 3 | Last Stone Weight | 🟢 | Max Heap | #1046 |
| 4 | K Closest Points to Origin | 🟡 | Heap | #973 |
| 5 | Merge K Sorted Lists | 🔴 | Min Heap | #23 |
| 6 | Find K Pairs with Smallest Sums | 🟡 | Min Heap | #373 |
| 7 | Kth Largest in Stream | 🟢 | Min Heap | #703 |
| 8 | IPO | 🔴 | Two Heaps | #502 |
| 9 | Find Median from Data Stream | 🔴 | Two Heaps | #295 |
| 10 | Sliding Window Median | 🔴 | Two Heaps | #480 |

### Heap Patterns
```python
# K largest: min heap of size K
heap = []
for num in arr:
    heapq.heappush(heap, num)
    if len(heap) > k:
        heapq.heappop(heap)

# Median: two heaps
# max_heap (left) | min_heap (right)
```

---

## Week 9: Graphs - Basics

### Focus
DFS, BFS, Union Find, connected components

### Problems (14)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Number of Islands | 🟡 | DFS/BFS | #200 |
| 2 | Max Area of Island | 🟡 | DFS | #695 |
| 3 | Flood Fill | 🟢 | DFS/BFS | #733 |
| 4 | Clone Graph | 🟡 | DFS/BFS | #133 |
| 5 | Pacific Atlantic Water Flow | 🟡 | DFS | #417 |
| 6 | Surrounded Regions | 🟡 | DFS/BFS | #130 |
| 7 | Rotting Oranges | 🟡 | BFS | #994 |
| 8 | Walls and Gates | 🟡 | BFS | #286 |
| 9 | Course Schedule | 🟡 | Topological | #207 |
| 10 | Course Schedule II | 🟡 | Topological | #210 |
| 11 | Number of Connected Components | 🟡 | Union Find | #323 |
| 12 | Graph Valid Tree | 🟡 | Union Find | #261 |
| 13 | Redundant Connection | 🟡 | Union Find | #684 |
| 14 | Keys and Rooms | 🟢 | DFS/BFS | #841 |

### Graph Representation
```python
# Adjacency list
graph = defaultdict(list)
for u, v in edges:
    graph[u].append(v)
    graph[v].append(u)  # if undirected
```

---

## Week 10: Graphs - Advanced

### Focus
Shortest path, MST, advanced traversals

### Problems (10)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Network Delay Time | 🟡 | Dijkstra | #743 |
| 2 | Cheapest Flights Within K Stops | 🟡 | Bellman-Ford | #787 |
| 3 | Path With Maximum Probability | 🟡 | Dijkstra | #1514 |
| 4 | Min Cost to Connect All Points | 🟡 | Prim's/Kruskal | #1584 |
| 5 | Word Ladder | 🔴 | BFS | #127 |
| 6 | Alien Dictionary | 🔴 | Topological | #269 |
| 7 | Reconstruct Itinerary | 🟡 | Euler Path | #332 |
| 8 | Minimum Height Trees | 🟡 | Topological | #310 |
| 9 | Accounts Merge | 🟡 | Union Find | #721 |
| 10 | Optimize Water Distribution | 🔴 | MST | #1168 |

---

## Week 11: Dynamic Programming - 1D

### Focus
Fibonacci-style, linear recurrence

### Problems (15)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Climbing Stairs | 🟢 | DP | #70 |
| 2 | Min Cost Climbing Stairs | 🟢 | DP | #746 |
| 3 | House Robber | 🟡 | DP | #198 |
| 4 | House Robber II | 🟡 | DP | #213 |
| 5 | Paint House | 🟡 | DP | #256 |
| 6 | Paint Fence | 🟡 | DP | #276 |
| 7 | Decode Ways | 🟡 | DP | #91 |
| 8 | Coin Change | 🟡 | DP | #322 |
| 9 | Minimum Coins | 🟡 | DP | #322 |
| 10 | Word Break | 🟡 | DP | #139 |
| 11 | Longest Increasing Subsequence | 🟡 | DP + BS | #300 |
| 12 | Number of Longest Increasing Subseq | 🟡 | DP | #673 |
| 13 | Partition Equal Subset Sum | 🟡 | DP | #416 |
| 14 | Maximum Subarray | 🟡 | Kadane's | #53 |
| 15 | Maximum Product Subarray | 🟡 | Kadane's | #152 |

### DP Template
```python
dp = [0] * (n + 1)
dp[0] = base_case
dp[1] = base_case

for i in range(2, n + 1):
    dp[i] = recurrence(dp, i)

return dp[n]
```

---

## Week 12: Dynamic Programming - 2D & Advanced

### Focus
2D grids, strings, interval DP

### Problems (15)

| # | Problem | Difficulty | Pattern | LeetCode |
|---|---------|------------|---------|----------|
| 1 | Unique Paths | 🟡 | 2D DP | #62 |
| 2 | Unique Paths II | 🟡 | 2D DP | #63 |
| 3 | Minimum Path Sum | 🟡 | 2D DP | #64 |
| 4 | Triangle | 🟡 | 2D DP | #120 |
| 5 | Longest Common Subsequence | 🟡 | 2D DP | #1143 |
| 6 | Edit Distance | 🟡 | 2D DP | #72 |
| 7 | Distinct Subsequences | 🔴 | 2D DP | #115 |
| 8 | Interleaving String | 🟡 | 2D DP | #97 |
| 9 | Target Sum | 🟡 | 2D DP | #494 |
| 10 | Palindromic Substrings | 🟡 | DP | #647 |
| 11 | Longest Palindromic Substring | 🟡 | DP | #5 |
| 12 | Burst Balloons | 🔴 | Interval DP | #312 |
| 13 | Regular Expression Matching | 🔴 | 2D DP | #10 |
| 14 | Wildcard Matching | 🔴 | 2D DP | #44 |
| 15 | Stone Game | 🟡 | Interval DP | #877 |

---

## Review Week: Mixed Practice

### Timed Practice Sessions

| Day | Focus | Problems | Time |
|-----|-------|----------|------|
| 1 | Arrays + Hashing | 5 | 75 min |
| 2 | Two Pointers + Sliding Window | 5 | 75 min |
| 3 | Trees + BST | 5 | 75 min |
| 4 | Graphs | 5 | 75 min |
| 5 | DP | 5 | 75 min |
| 6 | Mixed + Hard | 5 | 100 min |
| 7 | Mock Interviews | 3 | 135 min |

### Mock Interview Schedule
```
Mock 1: 45 min - Arrays/Strings focus
Mock 2: 45 min - Trees/Graphs focus
Mock 3: 45 min - DP + Mixed
```

---

## Progress Tracking

### Weekly Checklist

```
Week 1:  [ ] 15/15 problems  [ ] 1 mock interview
Week 2:  [ ] 15/15 problems  [ ] 1 mock interview
Week 3:  [ ] 12/12 problems  [ ] 1 mock interview
Week 4:  [ ] 12/12 problems  [ ] 1 mock interview
Week 5:  [ ] 12/12 problems  [ ] 1 mock interview
Week 6:  [ ] 15/15 problems  [ ] 1 mock interview
Week 7:  [ ] 12/12 problems  [ ] 1 mock interview
Week 8:  [ ] 10/10 problems  [ ] 1 mock interview
Week 9:  [ ] 14/14 problems  [ ] 1 mock interview
Week 10: [ ] 10/10 problems  [ ] 1 mock interview
Week 11: [ ] 15/15 problems  [ ] 1 mock interview
Week 12: [ ] 15/15 problems  [ ] 1 mock interview
```

### Total: 162 problems + 12 mock interviews

---

## Problem Sources

| Platform | Best For | Link |
|----------|----------|------|
| LeetCode | All problems | leetcode.com |
| NeetCode | Pattern learning | neetcode.io |
| Pramp | Mock interviews | pramp.com |
| InterviewBit | Company prep | interviewbit.com |

---

## Tips for Each Week

1. **Don't look at solutions for at least 30 minutes**
2. **Time yourself on each problem**
3. **Review solutions even if you solved it**
4. **Re-solve hard problems after 3 days**
5. **Focus on patterns, not memorization**
