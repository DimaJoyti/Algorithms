# 🎯 Mastering Coding Interviews: The Complete Guide

> A comprehensive roadmap from preparation to offer

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Interview Process Breakdown](#interview-process-breakdown)
3. [Preparation Strategy](#preparation-strategy)
4. [Pattern Recognition System](#pattern-recognition-system)
5. [The Interview Mindset](#the-interview-mindset)
6. [Communication Framework](#communication-framework)
7. [Common Mistakes & How to Avoid Them](#common-mistakes--how-to-avoid-them)
8. [Mock Interview Protocol](#mock-interview-protocol)
9. [Timeline & Roadmap](#timeline--roadmap)
10. [Company-Specific Prep](#company-specific-prep)
11. [Behavioral Integration](#behavioral-integration)
12. [Resources](#resources)

---

## Overview

### What This Guide Covers

```
┌─────────────────────────────────────────────────────────────┐
│                    INTERVIEW SUCCESS                         │
├─────────────────────────────────────────────────────────────┤
│  Technical Skills (40%)  │  Communication (30%)              │
│  - Patterns              │  - Explaining approach            │
│  - Implementation        │  - Asking clarifying questions    │
│  - Optimization          │  - Handling hints                 │
├─────────────────────────────────────────────────────────────┤
│  Problem Solving (20%)   │  Behavioral (10%)                 │
│  - Breaking down         │  - Past experiences               │
│  - Edge cases            │  - Team fit                       │
│  - Trade-offs            │  - Culture alignment              │
└─────────────────────────────────────────────────────────────┘
```

### Key Success Metrics

| Level | Problems Solved | Patterns Known | Mock Interviews | Timeline |
|-------|-----------------|----------------|-----------------|----------|
| Beginner | 50-100 | 5-8 | 3-5 | 2-3 months |
| Intermediate | 150-250 | 12-15 | 8-12 | 1-2 months |
| Advanced | 300+ | 18+ | 15+ | 2-4 weeks |

---

## Interview Process Breakdown

### Standard Tech Interview Flow

```
Application → Screen → Phone Interview → Onsite → Offer
    │            │            │              │
    │            │            │              └── 4-6 rounds
    │            │            └── 1-2 coding problems
    │            └── 15-30 min recruiter call
    └── Resume + Referral
```

### Interview Types

| Type | Duration | Focus | Weight |
|------|----------|-------|--------|
| Phone Screen | 45-60 min | 1-2 coding problems | High |
| Technical Onsite | 45 min each | Coding + system design | Critical |
| Behavioral | 30-45 min | Past experiences | Medium |
| System Design | 45 min | Architecture (senior) | High |

### What Interviewers Evaluate

**✅ Passing Signals:**
- Clear problem understanding
- Systematic approach
- Working solution
- Edge case handling
- Time/space optimization
- Clean code
- Good communication

**❌ Failing Signals:**
- Jumping to code without understanding
- No clarification questions
- Silent coding
- Missing edge cases
- Can't explain approach
- No optimization discussion

---

## Preparation Strategy

### Phase 1: Foundation (Weeks 1-4)

**Focus:** Core data structures & basic patterns

**Daily Routine:**
```
Morning (1h):   Learn 1 new pattern/concept
Evening (1h):   Practice 2-3 problems using that pattern
Weekly:         1 mock interview + review mistakes
```

**Checklist:**
- [ ] Master arrays & strings operations
- [ ] Understand HashMap/HashSet deeply
- [ ] Practice two pointers extensively
- [ ] Learn sliding window technique
- [ ] Implement binary search from scratch
- [ ] Understand recursion fundamentals

### Phase 2: Patterns Deep Dive (Weeks 5-8)

**Focus:** Recognize and apply patterns quickly

**Pattern Mastery Framework:**
```
For each pattern:
1. Understand WHEN to use it
2. Know the TEMPLATE code
3. Practice 5-10 problems
4. Identify edge cases
5. Time yourself
```

**Key Patterns Priority:**
1. **Arrays & Hashing** → Most common, foundation for others
2. **Two Pointers** → Sorted arrays, pairs, palindromes
3. **Sliding Window** → Subarrays, substrings
4. **Binary Search** → Sorted data, O(log n) requirement
5. **Trees (DFS/BFS)** → Hierarchical data
6. **Linked Lists** → Pointer manipulation
7. **Heaps** → Top K, priority-based
8. **Graphs** → Relationships, connections
9. **Dynamic Programming** → Optimization problems
10. **Backtracking** → All combinations/permutations

### Phase 3: Fluency (Weeks 9-12)

**Focus:** Speed, accuracy, and communication

**Daily Practice:**
```
Timed Sessions:
- 20 min: Solve 1 medium problem
- 35 min: Solve 1 medium-hard problem
- 45 min: Full mock interview simulation

Review:
- Analyze failed attempts
- Study optimal solutions
- Add to pattern notes
```

### Phase 4: Polish (Weeks 13-16)

**Focus:** Company-specific prep & confidence

**Activities:**
- [ ] Complete company-specific problem lists
- [ ] Practice behavioral stories
- [ ] Do 2-3 mock interviews per week
- [ ] Review all patterns (spaced repetition)
- [ ] Prepare questions for interviewers

---

## Pattern Recognition System

### Quick Decision Tree

```
Problem Type?
│
├─── SORTED ARRAY
│    ├── Find target → Binary Search
│    ├── Find pairs/sum → Two Pointers
│    └── Rotated → Modified Binary Search
│
├─── SUBARRAY/SUBSTRING
│    ├── Fixed size → Sliding Window (fixed)
│    ├── Variable size → Sliding Window (variable)
│    └── With constraint → Hash + Sliding Window
│
├─── PAIRS/COUNTING
│    ├── Find if exists → HashSet
│    ├── Find all/count → HashMap
│    └── Sorted → Two Pointers
│
├─── TREE
│    ├── Depth/path → DFS
│    ├── Level order → BFS
│    ├── Search in BST → BST properties
│    └── Serialization → DFS with markers
│
├─── GRAPH
│    ├── Connected components → DFS/Union Find
│    ├── Shortest path (unweighted) → BFS
│    ├── Shortest path (weighted) → Dijkstra
│    ├── Dependencies → Topological Sort
│    └── Minimum spanning tree → Kruskal/Prim
│
├─── OPTIMIZATION
│    ├── Overlapping subproblems → DP
│    ├── All combinations → Backtracking
│    ├── Local optimal = global → Greedy
│    └── Priority-based → Heap
│
├─── STRING
│    ├── Prefix matching → Trie
│    ├── Pattern matching → KMP/Rabin-Karp
│    ├── Parentheses → Stack
│    └── Comparison → Two Pointers/DP
│
└─── SPECIAL
     ├── Top K elements → Heap
     ├── K smallest/largest → Quick Select
     ├── Median → Two Heaps
     └── LRU/LFU → HashMap + LinkedList
```

### Pattern Templates

#### 1. Two Pointers

```python
def two_pointers(arr):
    l, r = 0, len(arr) - 1

    while l < r:
        # Process current state
        current = process(arr[l], arr[r])

        if condition_met(current):
            return result
        elif need_more(current):
            l += 1  # or modify based on problem
        else:
            r -= 1

    return default_result
```

**Use When:**
- Array is sorted
- Finding pairs with condition
- Palindrome validation
- Container problems

---

#### 2. Sliding Window

```python
def sliding_window(arr, k):
    # Fixed window size
    window_sum = sum(arr[:k])
    max_sum = window_sum

    for i in range(k, len(arr)):
        window_sum += arr[i] - arr[i - k]
        max_sum = max(max_sum, window_sum)

    return max_sum

def sliding_window_variable(arr):
    # Variable window size
    l = 0
    result = 0

    for r in range(len(arr)):
        # Add arr[r] to window

        while invalid_condition():
            # Remove arr[l] from window
            l += 1

        # Update result

    return result
```

**Use When:**
- Contiguous subarray/substring
- "Longest/shortest subarray where..."
- "Maximum/minimum window"
- Fixed k-size window

---

#### 3. Binary Search

```python
def binary_search(arr, target):
    l, r = 0, len(arr) - 1

    while l <= r:
        mid = l + (r - l) // 2  # Avoid overflow

        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            l = mid + 1
        else:
            r = mid - 1

    return -1  # Not found

# Find first True in [F, F, F, T, T, T]
def binary_search_predicate():
    l, r = 0, n  # Note: r = n, not n-1

    while l < r:
        mid = l + (r - l) // 2

        if condition(mid):
            r = mid
        else:
            l = mid + 1

    return l
```

**Use When:**
- Sorted array
- "Find first/last occurrence"
- "Minimum value such that..."
- Rotated sorted array
- Search space reduction

---

#### 4. DFS (Trees/Graphs)

```python
# Tree DFS
def dfs_tree(node):
    if not node:
        return base_case

    left = dfs_tree(node.left)
    right = dfs_tree(node.right)

    return combine(left, right, node.val)

# Graph DFS
def dfs_graph(node, visited, graph):
    if node in visited:
        return

    visited.add(node)

    for neighbor in graph[node]:
        dfs_graph(neighbor, visited, graph)
```

**Use When:**
- Path problems
- Depth/height calculations
- Connected components
- Cycle detection
- Tree traversals

---

#### 5. BFS (Level Order)

```python
from collections import deque

def bfs_tree(root):
    if not root:
        return []

    result = []
    queue = deque([root])

    while queue:
        level = []
        level_size = len(queue)

        for _ in range(level_size):
            node = queue.popleft()
            level.append(node.val)

            if node.left:
                queue.append(node.left)
            if node.right:
                queue.append(node.right)

        result.append(level)

    return result

def bfs_graph(start, graph):
    visited = set([start])
    queue = deque([start])

    while queue:
        node = queue.popleft()

        for neighbor in graph[node]:
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(neighbor)
```

**Use When:**
- Level order traversal
- Shortest path (unweighted)
- Minimum steps to target
- "Nearest" problems
- Propagation/contagion

---

#### 6. Heap / Priority Queue

```python
import heapq

# Min heap (default)
heap = []
heapq.heappush(heap, item)
min_val = heapq.heappop(heap)

# Max heap (negate values)
max_heap = []
heapq.heappush(max_heap, -value)
max_val = -heapq.heappop(max_heap)

# K largest elements
def k_largest(arr, k):
    heap = []
    for num in arr:
        heapq.heappush(heap, num)
        if len(heap) > k:
            heapq.heappop(heap)
    return heap
```

**Use When:**
- Top K elements
- Merge K sorted arrays
- Find median (two heaps)
- Priority-based processing
- "Kth smallest/largest"

---

#### 7. Dynamic Programming

```python
# 1D DP
def dp_1d(n):
    if n <= 1:
        return n

    dp = [0] * (n + 1)
    dp[0] = base_0
    dp[1] = base_1

    for i in range(2, n + 1):
        dp[i] = recurrence(dp, i)

    return dp[n]

# Space optimized
def dp_optimized(n):
    a, b = base_0, base_1

    for _ in range(2, n + 1):
        a, b = b, f(a, b)

    return b

# 2D DP
def dp_2d(m, n):
    dp = [[0] * (n + 1) for _ in range(m + 1)]

    # Initialize base cases
    for i in range(m + 1):
        dp[i][0] = base_i
    for j in range(n + 1):
        dp[0][j] = base_j

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            dp[i][j] = recurrence(dp, i, j)

    return dp[m][n]
```

**Use When:**
- Overlapping subproblems
- Optimal substructure
- "Maximum/minimum ways"
- "Count number of ways"
- Optimization problems

---

#### 8. Backtracking

```python
def backtrack(path, options, result):
    if is_complete(path):
        result.append(path[:])  # Copy!
        return

    for i, option in enumerate(options):
        # Choose
        path.append(option)

        # Explore
        backtrack(path, options[i+1:], result)  # or options if repetition allowed

        # Unchoose (backtrack)
        path.pop()

# With pruning
def backtrack_pruned(path, options, result):
    if is_complete(path):
        result.append(path[:])
        return

    for option in options:
        if not is_valid(path, option):
            continue  # Prune invalid branches

        path.append(option)
        backtrack_pruned(path, options, result)
        path.pop()
```

**Use When:**
- All permutations
- All combinations
- All subsets
- Constraint satisfaction
- "Find all possible..."

---

#### 9. Union Find

```python
class UnionFind:
    def __init__(self, n):
        self.parent = list(range(n))
        self.rank = [0] * n

    def find(self, x):
        if self.parent[x] != x:
            self.parent[x] = self.find(self.parent[x])  # Path compression
        return self.parent[x]

    def union(self, x, y):
        px, py = self.find(x), self.find(y)
        if px == py:
            return False  # Already connected

        # Union by rank
        if self.rank[px] < self.rank[py]:
            px, py = py, px
        self.parent[py] = px
        self.rank[px] += self.rank[px] == self.rank[py]
        return True
```

**Use When:**
- Connected components
- Cycle detection in graphs
- "Are these connected?"
- Dynamic connectivity
- Kruskal's MST

---

#### 10. Monotonic Stack

```python
# Next Greater Element
def next_greater(arr):
    result = [-1] * len(arr)
    stack = []  # Indices

    for i in range(len(arr)):
        while stack and arr[stack[-1]] < arr[i]:
            result[stack.pop()] = arr[i]
        stack.append(i)

    return result

# Largest Rectangle in Histogram
def largest_rectangle(heights):
    stack = []  # Indices
    max_area = 0

    for i in range(len(heights) + 1):
        h = heights[i] if i < len(heights) else 0

        while stack and heights[stack[-1]] > h:
            height = heights[stack.pop()]
            width = i if not stack else i - stack[-1] - 1
            max_area = max(max_area, height * width)

        stack.append(i)

    return max_area
```

**Use When:**
- Next greater/smaller element
- Temperature/spread problems
- Histogram problems
- Parentheses matching
- "First element to left/right that..."

---

## The Interview Mindset

### Before the Interview

**Mental Preparation:**
```
✅ DO:
- Get 7-8 hours sleep
- Eat a light meal
- Review pattern templates (30 min)
- Prepare your environment (quiet, water, notepad)
- Have your IDE/editor ready

❌ DON'T:
- Cram new problems
- Stay up late
- Compare yourself to others
- Overthink potential questions
```

### During the Interview

**The REACT Method:**

```
R - Read the problem carefully
E - Explore examples (ask for clarification)
A - Analyze approach (discuss before coding)
C - Code the solution (explain as you go)
T - Test and optimize (edge cases, complexity)
```

### Time Management (45 min interview)

```
Minutes 0-5:   Understand problem + clarify
Minutes 5-10:  Discuss approach + get buy-in
Minutes 10-25: Implement solution
Minutes 25-35: Test + debug
Minutes 35-45: Optimize + discuss alternatives
```

---

## Communication Framework

### Clarifying Questions Template

```
1. "Let me make sure I understand the problem..."
   [Restate in your own words]

2. "Can you give me an example input/output?"
   [Work through it together]

3. "What should I return if...?"
   [Ask about edge cases]

4. "Are there constraints on...?"
   [Input size, values, time/space requirements]

5. "Is the input sorted/unique/valid?"
   [Clarify assumptions]
```

### Explaining Your Approach

```
Structure:
1. "I'm thinking of using [PATTERN] because..."
2. "The time complexity would be O(?) because..."
3. "Space complexity is O(?) because..."
4. "An alternative approach would be [X], but [Y] is better because..."
```

### While Coding

```
✅ "Let me start by handling the base case..."
✅ "Now I'll iterate through..."
✅ "I need to track [X] and [Y]..."
✅ "This handles the edge case where..."
✅ "Let me add a check for..."

❌ [Silence for 5+ minutes]
❌ "Hmm, let me think..." [long pause]
❌ [Typing without explanation]
```

### Handling Hints

```
When interviewer gives a hint:

1. Acknowledge: "That's a good point..."
2. Process: "So if I use [hint], then..."
3. Apply: "Let me modify the code to..."
4. Confirm: "Does this align with what you were thinking?"
```

---

## Common Mistakes & How to Avoid Them

### Technical Mistakes

| Mistake | Fix |
|---------|-----|
| Jumping to code | Spend 5 min understanding + planning |
| Missing edge cases | Create edge case checklist |
| Not analyzing complexity | Always state O(?) before coding |
| Ignoring constraints | Read constraints, design accordingly |
| Silent debugging | Narrate your debugging process |
| Overcomplicating | Start simple, then optimize |

### Communication Mistakes

| Mistake | Fix |
|---------|-----|
| No questions | Always ask 2-3 clarifying questions |
| Silent coding | Explain each step |
| Ignoring hints | Acknowledge and incorporate |
| Defensive attitude | Stay open to suggestions |
| Rushing | Slow down, think out loud |

### Pattern-Specific Pitfalls

**Arrays/Strings:**
- ❌ Off-by-one errors in indexing
- ✅ Use `range(len(arr))` or `enumerate()`

**Two Pointers:**
- ❌ Not handling duplicates
- ✅ Skip duplicates with `while l < r and arr[l] == arr[l-1]`

**Sliding Window:**
- ❌ Forgetting to shrink window
- ✅ Always have condition to move left pointer

**Binary Search:**
- ❌ Infinite loop with `mid` calculation
- ✅ Use `l + (r - l) // 2` and correct update

**Trees:**
- ❌ Not handling null/empty cases
- ✅ Always check `if not node: return`

**DP:**
- ❌ Wrong base cases
- ✅ Verify base cases with small examples

---

## Mock Interview Protocol

### Self-Mock Structure

```
1. Pick a problem (don't look at solution)
2. Set 45-minute timer
3. Record yourself (audio/video)
4. Follow REACT method
5. Review recording for:
   - Communication gaps
   - Time management
   - Technical accuracy
6. Study optimal solution
7. Re-implement from memory next day
```

### Peer Mock Structure

```
Roles:
- Interviewer: Has solution, gives hints if stuck
- Interviewee: Thinks out loud, asks questions

Flow:
1. 5 min: Interviewer explains problem
2. 35 min: Solve with realistic constraints
3. 5 min: Feedback exchange

Feedback Template:
✅ Strengths: [What went well]
⚠️ Areas for improvement: [What to work on]
💡 Suggestion: [Specific advice]
```

### Mock Interview Platforms

| Platform | Best For | Cost |
|----------|----------|------|
| Pramp | Peer mocks | Free |
| interviewing.io | Expert mocks | Paid |
| LeetCode Mock | Timed practice | Free |
| CodeSignal | Realistic environment | Freemium |

---

## Timeline & Roadmap

### 12-Week Plan (Beginner → Interview Ready)

#### Weeks 1-2: Foundation
```
Days 1-7:   Arrays + Hashing (15 problems)
Days 8-14:  Two Pointers + Strings (15 problems)

Milestone: 30 problems, 3 patterns
```

#### Weeks 3-4: Core Patterns
```
Days 15-21: Sliding Window + Binary Search (15 problems)
Days 22-28: Linked Lists + Stacks (15 problems)

Milestone: 60 problems, 6 patterns
```

#### Weeks 5-6: Trees & Graphs
```
Days 29-35: Trees (DFS/BFS) (15 problems)
Days 36-42: Graphs + Heaps (15 problems)

Milestone: 90 problems, 9 patterns
Mock: 2 mock interviews
```

#### Weeks 7-8: Advanced Patterns
```
Days 43-49: Backtracking + DP 1D (15 problems)
Days 50-56: DP 2D + Greedy (15 problems)

Milestone: 120 problems, 12 patterns
Mock: 4 mock interviews
```

#### Weeks 9-10: Fluency
```
Days 57-63: Mixed practice (20 problems)
Days 64-70: Timed practice (20 problems)

Milestone: 160 problems, 15 patterns
Mock: 6 mock interviews
```

#### Weeks 11-12: Polish
```
Days 71-77: Company-specific problems
Days 78-84: Final review + mocks

Milestone: 180+ problems, 15+ patterns
Mock: 8+ mock interviews
```

### 4-Week Intensive (Experienced)

```
Week 1: Patterns 1-6 (40 problems)
Week 2: Patterns 7-12 (40 problems)
Week 3: Mixed + Timed (40 problems)
Week 4: Company prep + Mocks (20 problems + 8 mocks)
```

---

## Company-Specific Prep

### FAANG+ Patterns

| Company | Focus Areas | Difficulty |
|---------|-------------|------------|
| Google | DP, Graphs, Design | Hard |
| Meta | Arrays, Trees, Two Pointers | Medium-Hard |
| Amazon | Trees, DP, Behavioral | Medium |
| Apple | DP, System Design | Medium-Hard |
| Netflix | System Design, Culture | Variable |
| Microsoft | Trees, Graphs, DP | Medium |

### Company Problem Lists

**Google:**
- DP (optimization problems)
- Graph algorithms
- System design (senior)
- Focus: Correctness + edge cases

**Meta:**
- Two pointers (many variants)
- Trees (all traversals)
- Fast implementation
- Focus: Speed + correctness

**Amazon:**
- Trees + BST
- DP (moderate difficulty)
- Behavioral (Leadership Principles)
- Focus: Working solution + LPs

**Apple:**
- DP (varied)
- Clean code
- System design
- Focus: Elegance + correctness

---

## Behavioral Integration

### STAR Method

```
S - Situation: Set the context
T - Task: What you needed to accomplish
A - Action: What YOU did (not team)
R - Result: Quantifiable outcome
```

### Common Behavioral Questions

**Technical:**
- "Tell me about a challenging bug you fixed"
- "Describe a time you had to learn something quickly"
- "How do you approach code reviews?"

**Collaboration:**
- "Tell me about a conflict with a teammate"
- "How do you handle disagreements about technical decisions?"
- "Describe working with a difficult stakeholder"

**Growth:**
- "What's a technical mistake you made?"
- "How do you stay current with technology?"
- "Tell me about feedback that changed your approach"

### Prepare 5-7 Stories

```
Story 1: Technical challenge
Story 2: Leadership/ownership
Story 3: Conflict resolution
Story 4: Failure + learning
Story 5: Collaboration
Story 6: Innovation/creativity
Story 7: Mentorship/teaching
```

---

## Resources

### Essential Books

| Book | Focus | Level |
|------|-------|-------|
| Cracking the Coding Interview | Comprehensive prep | All |
| Elements of Programming Interviews | Problem patterns | Intermediate |
| Grokking Algorithms | Algorithm basics | Beginner |
| System Design Interview | Architecture | Senior |

### Online Platforms

| Platform | Best For |
|----------|----------|
| LeetCode | Problem practice |
| NeetCode | Pattern-based learning |
| AlgoExpert | Video explanations |
| InterviewBit | Company prep |
| Codeforces | Competitive programming |

### Problem Lists

| List | Problems | Focus |
|------|----------|-------|
| Blind 75 | 75 | Essential patterns |
| NeetCode 150 | 150 | Comprehensive |
| Grind 75 | 75 | Time-based |
| LeetCode Top 100 | 100 | Most common |

### Time Complexity Reference

| Complexity | n=10 | n=100 | n=1000 | n=10^6 |
|-----------|------|-------|--------|--------|
| O(1) | ✓ | ✓ | ✓ | ✓ |
| O(log n) | ✓ | ✓ | ✓ | ✓ |
| O(n) | ✓ | ✓ | ✓ | ✓ |
| O(n log n) | ✓ | ✓ | ✓ | ~ |
| O(n²) | ✓ | ~ | ✗ | ✗ |
| O(2^n) | ~ | ✗ | ✗ | ✗ |

✓ = Fast enough | ~ = Might be okay | ✗ = Too slow

---

## Quick Reference Cards

### Before Every Problem

```
□ Understand problem (restate in own words)
□ Clarify inputs/outputs
□ Ask about constraints
□ Discuss edge cases
□ State approach + complexity
□ Get interviewer buy-in
□ Code while explaining
□ Test with examples
□ Handle edge cases
□ Discuss optimizations
```

### Complexity Quick Guide

```
HashMap operations: O(1) average
Binary search: O(log n)
Sorting: O(n log n)
Tree traversal: O(n)
BFS/DFS: O(V + E)
Heap operations: O(log n)
Two pointers: O(n)
Sliding window: O(n)
```

### Common Edge Cases

```
□ Empty input
□ Single element
□ Two elements
□ All same elements
□ Sorted (ascending/descending)
□ Negative numbers
□ Duplicates
□ Integer overflow
□ Null/undefined values
□ Maximum size input
```

---

## Final Checklist

### Week Before Interview

- [ ] Review all pattern templates
- [ ] Complete 2-3 mock interviews
- [ ] Prepare behavioral stories
- [ ] Research company + role
- [ ] Prepare questions for interviewers
- [ ] Test your setup (if remote)

### Day Before Interview

- [ ] Light review (no new problems)
- [ ] Prepare interview environment
- [ ] Get 8 hours sleep
- [ ] Eat well, stay hydrated

### Day of Interview

- [ ] Eat a light meal
- [ ] Review pattern templates (15 min)
- [ ] Arrive early / test connection
- [ ] Have water and notepad ready
- [ ] Stay calm and confident

---

## Summary

**Success = Preparation × Communication × Mindset**

```
┌─────────────────────────────────────────────────────────────┐
│                      THE WINNING FORMULA                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   PREPARATION          COMMUNICATION         MINDSET         │
│   ───────────          ─────────────         ────────        │
│   • Patterns           • Clarify             • Confident     │
│   • Practice           • Explain             • Calm          │
│   • Mocks              • Collaborate         • Curious       │
│                                                              │
│   = Technical Skill    = Interview Presence  = Performance   │
│                                                              │
│                    YOU'VE GOT THIS! 💪                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 📚 Complete Documentation Suite

This guide is part of a comprehensive interview preparation system:

| Document | Purpose | When to Use |
| --------- | -------- | ------------ |
| **[Mastery Guide](mastery-guide.md)** | Main reference | Daily study |
| **[Cheatsheet](cheatsheet.md)** | 1-page quick ref | Day of interview |
| **[System Design](system-design.md)** | Senior interviews | System design prep |
| **[Basic Patterns](../patterns/basic-patterns.md)** | Core patterns | Foundation |
| **[Advanced Patterns](../patterns/advanced-patterns.md)** | Complex patterns | Medium/Hard problems |
| **[Decision Tree](../patterns/decision-tree.md)** | Pattern selection | Quick decisions |
| **[Weekly Practice Sets](../practice/weekly-sets.md)** | Structured practice | 12-week plan |

---

## 🧠 Behavioral Question Bank

### Technical Leadership

```
Q: "Tell me about a time you had to make a difficult technical decision"

STAR Template:
S: "We had a tight deadline for a payment system migration"
T: "I needed to decide between rewriting or incremental migration"
A: "I analyzed risks, consulted team, created rollback plan, chose incremental"
R: "Zero downtime, 100% data integrity, delivered on time"

Key Points:
- Show decision framework
- Demonstrate risk assessment
- Highlight team collaboration
```

### Conflict Resolution

```
Q: "Describe a disagreement with a teammate about implementation"

Good Response Elements:
- Listened to their perspective
- Focused on data/evidence
- Found common ground
- Prioritized project goals over ego
- Reached consensus or escalated appropriately
```

### Handling Failure

```
Q: "Tell me about a time you made a mistake"

Structure:
1. Own the mistake directly
2. Explain the impact
3. Describe immediate fix
4. Share long-term prevention
5. Show what you learned

Avoid:
- Blaming others
- Minimizing impact
- No lessons learned
```

### 10 Must-Prepare Stories

| # | Theme | Question Type |
|---|-------|---------------|
| 1 | Technical Challenge | "Most difficult bug you fixed" |
| 2 | Leadership | "Led a project without authority" |
| 3 | Conflict | "Disagreed with a technical decision" |
| 4 | Failure | "Project that didn't go as planned" |
| 5 | Collaboration | "Worked with difficult team member" |
| 6 | Growth | "Learned a new technology quickly" |
| 7 | Innovation | "Improved an existing process" |
| 8 | Mentorship | "Helped a teammate grow" |
| 9 | Ownership | "Took initiative beyond your role" |
| 10 | Adaptability | "Handled changing requirements" |

---

## 🎮 Advanced DP Patterns

### State Machine DP

```python
# Best Time to Buy/Sell Stock with Cooldown
def maxProfit(prices):
    n = len(prices)
    if n < 2:
        return 0

    # States: holding, sold (cooldown), rest
    hold = -prices[0]  # Max profit while holding
    sold = 0           # Max profit just sold (cooldown)
    rest = 0           # Max profit in rest state

    for i in range(1, n):
        prev_hold, prev_sold, prev_rest = hold, sold, rest

        # Transitions
        hold = max(prev_hold, prev_rest - prices[i])  # Buy or keep
        sold = prev_hold + prices[i]                   # Sell today
        rest = max(prev_rest, prev_sold)               # Rest or from cooldown

    return max(sold, rest)
```

### Digit DP

```python
# Count numbers <= n where digit sum equals target
def countNumbers(n: str, target: int) -> int:
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

### Bitmask DP

```python
# Traveling Salesman Problem
def tsp(costs):
    n = len(costs)

    @lru_cache(None)
    def dp(mask: int, pos: int) -> int:
        if mask == (1 << n) - 1:
            return costs[pos][0]  # Return to start

        result = float('inf')
        for city in range(n):
            if not (mask & (1 << city)):
                result = min(result, costs[pos][city] + dp(mask | (1 << city), city))

        return result

    return dp(1, 0)  # Start at city 0, visited only city 0
```

---

## 🔧 Language-Specific Tips

### Python

```python
# Common shortcuts
from collections import defaultdict, Counter, deque
from functools import lru_cache
import heapq

# Counter for frequency
freq = Counter(arr)
most_common = freq.most_common(k)

# DefaultDict
graph = defaultdict(list)

# Deque for BFS
q = deque([start])
q.popleft()
q.append(neighbor)

# Heap
heapq.heappush(heap, item)
heapq.heappop(heap)

# LRU Cache for memoization
@lru_cache(None)
def dp(state):
    return result
```

### JavaScript/TypeScript

```javascript
// Map for frequency
const freq = new Map();
for (const num of arr) {
    freq.set(num, (freq.get(num) || 0) + 1);
}

// Set for uniqueness
const seen = new Set();

// Queue for BFS
const queue = [start];
while (queue.length) {
    const node = queue.shift();
    // process
}

// Sort with comparator
arr.sort((a, b) => a - b);

// Object as HashMap
const map = {};
map[key] = value;
```

### Go

```go
// Map for frequency
freq := make(map[int]int)
for _, num := range arr {
    freq[num]++
}

// Slice operations
sort.Ints(arr)
sort.Slice(arr, func(i, j int) bool {
    return arr[i] < arr[j]
})

// Heap (implement container/heap interface)
type MinHeap []int
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
```

---

## 📊 Company Interview Styles

| Company | Style | Focus | Tips |
|---------|-------|-------|------|
| **Google** | 4-5 rounds | DP, Graphs, Scale | Focus on correctness, explain approach |
| **Meta** | 2 coding + behavioral | Speed, Clean code | Practice fast implementation |
| **Amazon** | Coding + LP | Working solution | Emphasize Leadership Principles |
| **Apple** | Domain-specific | Clean code, Design | Show passion for products |
| **Microsoft** | Balanced | Communication | Collaborative problem solving |
| **Netflix** | Senior-focused | System design | Culture fit, technical depth |
| **Stripe** | Practical | Real-world problems | Think about edge cases |
| **Airbnb** | Full-stack | Product thinking | Consider user experience |

---

## ⏰ Final Week Checklist

### 7 Days Before

```
□ Review all pattern templates (1/day)
□ 2 mock interviews
□ Prepare 7 behavioral stories
□ Research company culture
□ Prepare 3-5 questions for interviewers
□ Test remote setup (camera, mic, IDE)
□ Plan interview day logistics
```

### Day Before

```
□ Light review only (no new problems)
□ Review pattern templates (30 min)
□ Review behavioral stories (30 min)
□ Prepare interview space
□ Charge devices
□ Get 8 hours sleep
```

### Day of Interview

```
□ Eat a light meal
□ Review cheat sheet (15 min max)
□ Arrive/test connection 15 min early
□ Have water, notepad, pen ready
□ Take deep breaths between rounds
```

---

## 🏆 Success Metrics

### By Preparation Level

| Level | Problems | Patterns | Mocks | Timeline |
|-------|----------|----------|-------|----------|
| Beginner | 50-100 | 5-8 | 3-5 | 2-3 months |
| Intermediate | 150-250 | 12-15 | 8-12 | 1-2 months |
| Advanced | 300+ | 18+ | 15+ | 2-4 weeks |

### Weekly Goals

```
Week 1-4:   15-20 problems/week, 1 mock
Week 5-8:   15-20 problems/week, 2 mocks
Week 9-12:  20-25 problems/week, 2-3 mocks
```

---

*Remember: Every expert was once a beginner. Consistent practice beats cramming every time.*

**Good luck with your interviews!** 🚀
