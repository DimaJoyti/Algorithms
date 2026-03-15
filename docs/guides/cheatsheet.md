# Coding Interview Cheat Sheet (1 Page)

> Quick reference for interview day

---

## Before You Code (5 min)

```
1. RESTATE problem in your own words
2. ASK: Input size? Sorted? Duplicates? Edge cases?
3. STATE approach + complexity BEFORE coding
4. GET buy-in from interviewer
```

---

## Pattern Decision Tree

```
SORTED ARRAY?
├─ Find target → Binary Search
├─ Find pairs → Two Pointers
└─ Rotated → Modified BS

SUBARRAY/SUBSTRING?
├─ Fixed size → Sliding Window (fixed)
└─ Variable → Sliding Window + Hash

PAIRS/COUNTING?
├─ Exists? → HashSet
└─ Count all → HashMap

TREE?
├─ Depth/Path → DFS
├─ Level order → BFS
└─ BST search → BST properties

GRAPH?
├─ Connected? → DFS/Union Find
├─ Shortest (unweighted) → BFS
├─ Shortest (weighted) → Dijkstra
└─ Dependencies → Topological Sort

OPTIMIZATION?
├─ Overlapping subproblems → DP
├─ All combinations → Backtracking
└─ Priority-based → Heap

TOP K? → Heap | MEDIAN? → Two Heaps | PREFIX? → Trie
```

---

## Code Templates

### Two Pointers
```python
l, r = 0, len(arr) - 1
while l < r:
    if condition: return result
    elif need_more: l += 1
    else: r -= 1
```

### Sliding Window
```python
l = 0
for r in range(len(arr)):
    # add arr[r]
    while invalid():
        # remove arr[l]
        l += 1
    # update result
```

### Binary Search
```python
l, r = 0, len(arr) - 1
while l <= r:
    mid = l + (r - l) // 2
    if arr[mid] == target: return mid
    elif arr[mid] < target: l = mid + 1
    else: r = mid - 1
```

### DFS (Tree)
```python
def dfs(node):
    if not node: return base
    return combine(dfs(node.left), dfs(node.right))
```

### BFS (Level)
```python
q = deque([root])
while q:
    for _ in range(len(q)):
        node = q.popleft()
        # process, add children
```

### Backtracking
```python
def backtrack(path, options):
    if complete: result.append(path[:]); return
    for opt in options:
        path.append(opt)
        backtrack(path, remaining)
        path.pop()
```

---

## Complexity Quick Reference

| Pattern | Time | Space |
|---------|------|-------|
| Two Pointers | O(n) | O(1) |
| Sliding Window | O(n) | O(k) |
| Binary Search | O(log n) | O(1) |
| DFS/BFS | O(V+E) | O(V) |
| Heap | O(n log k) | O(k) |
| Sorting | O(n log n) | O(n) |
| DP | O(n×m) | O(n×m) |

---

## Edge Cases Checklist

```
□ Empty input []
□ Single element [1]
□ Two elements [1,2]
□ All same [1,1,1]
□ Sorted / Reverse sorted
□ Negative numbers
□ Duplicates
□ Integer overflow
□ Null/None values
```

---

## Communication Phrases

| Situation | Say |
|-----------|-----|
| Clarify | "Let me make sure I understand..." |
| Approach | "I'm thinking [pattern] because..." |
| Complexity | "Time is O(n) because..., space is O(1)" |
| Stuck | "Let me think about this differently..." |
| Hint received | "That's helpful, let me incorporate that..." |
| Done | "Let me trace through with an example..." |

---

## 45-Min Time Budget

```
0-5 min:   Understand + Clarify
5-10 min:  Approach + Complexity
10-25 min: Code + Explain
25-35 min: Test + Debug
35-45 min: Optimize + Discuss
```

---

## Common Mistakes → Fixes

| Mistake | Fix |
|---------|-----|
| Jumping to code | Plan first, get buy-in |
| Silent coding | Narrate every step |
| Missing edge cases | Test empty, single, duplicates |
| Wrong complexity | Count operations, not lines |
| Infinite loops | Check loop conditions |

---

## Behavioral STAR

```
S - Situation: Context in 1 sentence
T - Task: What YOU needed to do
A - Action: What YOU did (3-4 steps)
R - Result: Quantifiable outcome
```

**Prepare 5 stories:** Technical challenge, Leadership, Conflict, Failure, Collaboration

---

## Day of Interview

```
✅ Sleep 8 hours
✅ Eat light meal
✅ Review templates (15 min)
✅ Water + notepad ready
✅ Arrive/test connection early
```

---

*You've prepared. You're ready. You've got this!*
