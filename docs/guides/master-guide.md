# Mastering Coding Interviews

> The Complete Engineer's Playbook

---

## Quick Navigation

| Chapter | Topic | Key Takeaway |
| ------- | ----- | ------------ |
| 01 | [Interview Landscape](#01-interview-landscape) | Pipeline, scoring, company differences |
| 02 | [12-Week Roadmap](#02-12-week-study-roadmap) | Phased plan, daily structure |
| 03 | [Data Structures](#03-data-structures-deep-dive) | Complexity tables, when to use |
| 04 | [Algorithm Patterns](#04-the-14-universal-algorithm-patterns) | Pattern recognition map |
| 05 | [Problem-Solving](#05-the-problem-solving-framework) | UMPIRE method, scripts |
| 06 | [Code Quality](#06-coding-best-practices) | Naming, structure, language tips |
| 07 | [System Design](#07-system-design-interviews) | RESHADED framework |
| 08 | [Behavioral](#08-behavioral-interviews) | STAR+, Amazon LPs |
| 09 | [Live Tactics](#09-live-interview-tactics) | Stuck playbook, red flags |
| 10 | [Negotiation](#10-offer-evaluation--negotiation) | Scripts, total comp |

---

## 01. Interview Landscape

### The Typical Interview Pipeline

| Stage | Duration | Format | What Actually Matters |
| ----- | -------- | ------ | --------------------- |
| Recruiter Screen | 15–30 min | Phone | Interest, availability, comp alignment |
| Online Assessment | 60–90 min | Automated | 2–3 problems, no communication |
| Technical Phone | 45–60 min | Live coding | Solution + narrating thinking |
| Onsite | 4–6 hours | Multiple rounds | DSA + system design + behavioral |
| Team Matching | 30–45 min | Informal | Culture fit, domain interest |
| Offer | 1–2 weeks | Email/phone | Total comp, leveling, start date |

### FAANG vs Startups vs Mid-Tier

| Dimension | FAANG | Mid-Tier | Startups |
| --------- | ----- | -------- | -------- |
| Algorithm Difficulty | Hard LC | Medium LC | Often take-home |
| System Design | Mandatory L4+ | Sometimes senior | Rarely formal |
| Behavioral Depth | Heavy (LPs) | Moderate | Light, culture fit |
| Rounds | 5–7 | 3–5 | 1–3 |
| Bar | Standardized, high | Varies by team | Highly variable |

### The Scoring Rubric

| Dimension | What It Means | How to Signal |
| --------- | ------------- | ------------- |
| Problem Solving | Break down new problems | State brute force, then improve |
| Communication | Explain reasoning clearly | Never go silent; narrate decisions |
| Technical Accuracy | Correct + handles edge cases | Test with 2+ examples |
| Code Quality | Clean, readable, maintainable | Meaningful names, small functions |
| Efficiency | Know complexity, can optimize | State time/space unprompted |
| Collaboration | Pleasant, receptive to hints | Embrace hints immediately |

> **The Silent Killer**: Most rejections happen because of coding silently. You can have the correct answer and still fail if the interviewer can't follow your thinking.

---

## 02. 12-Week Study Roadmap

### The Core Principle

> Deeply understanding **150 problems** is worth more than shallowly seeing 500. After solving, always ask: *What is the pattern? What other problems use this?*

### Phased Plan

| Phase | Weeks | Focus | Goal |
| ----- | ----- | ----- | ---- |
| Foundation | 1–2 | Big-O, arrays, strings, hash maps | Easy in <20 min |
| Core Structures | 3–4 | Linked lists, stacks, trees, heaps | Fluent BFS/DFS |
| Algorithm Patterns | 5–6 | Two pointers, sliding window, binary search | Recognize pattern instantly |
| Advanced Algorithms | 7–8 | DP, graphs, Union-Find, Dijkstra | 70% of Medium DPs |
| System Design | 8–9 | Scalability, databases, caching | Lead 45-min discussion |
| Mock & Polish | 10–11 | Timed mocks, behavioral, weak areas | Perform under pressure |
| Company Targeting | 12 | Company-specific banks, calibration | Know top 20 patterns |

### Daily Structure (1.5 hours)

```text
30 min → Solve 1 new problem (timed, no hints 15 min)
20 min → Review solution + optimal answer
20 min → Re-implement from scratch
20 min → Concept reading or video
```

### Weekend Structure (3–4 hours)

```text
Saturday: Full 60-min timed mock (Pramp/Interviewing.io)
Sunday:   System design deep dive + behavioral stories + week review
```

### Spaced Repetition Hack

```text
Easy (solved <10 min)   → revisit in 1 week
Medium (needed hints)   → revisit in 3 days
Hard (couldn't solve)   → revisit tomorrow

This cuts required problems by 40%
```

---

## 03. Data Structures Deep Dive

### Master Complexity Reference

| Structure | Access | Search | Insert | Delete | Space |
| --------- | ------ | ------ | ------ | ------ | ----- |
| Array | O(1) | O(n) | O(n) | O(n) | O(n) |
| Dynamic Array | O(1) | O(n) | O(1)* | O(n) | O(n) |
| Hash Map | O(1)* | O(1)* | O(1)* | O(1)* | O(n) |
| Linked List | O(n) | O(n) | O(1) | O(1) | O(n) |
| Stack | O(n) | O(n) | O(1) | O(1) | O(n) |
| Queue | O(n) | O(n) | O(1) | O(1) | O(n) |
| Binary Heap | — | O(n) | O(log n) | O(log n) | O(n) |
| BST (balanced) | O(log n) | O(log n) | O(log n) | O(log n) | O(n) |
| Trie | O(m) | O(m) | O(m) | O(m) | O(n·m) |

*Average case; amortized for dynamic array

### Tree Traversals

| Traversal | Order | Use Case |
| --------- | ----- | -------- |
| Inorder | L→Root→R | Sorted output (BST) |
| Preorder | Root→L→R | Copy tree, serialize |
| Postorder | L→R→Root | Delete, evaluate |
| BFS | Level-by-level | Shortest path |

### Graph Algorithms

| Algorithm | Use Case | Time |
| --------- | -------- | ---- |
| BFS | Unweighted shortest path | O(V+E) |
| DFS | Cycle detection, components | O(V+E) |
| Dijkstra | Weighted shortest path (+) | O((V+E) log V) |
| Bellman-Ford | Negative weights | O(VE) |
| Union-Find | Connected components | ~O(1) |
| Topological Sort | Dependencies | O(V+E) |

---

## 04. The 14 Universal Algorithm Patterns

### Pattern Recognition Map

| If problem mentions... | Think: Pattern | Key Data Structure |
| --------------------- | -------------- | ------------------ |
| Sorted array + target | Two Pointers | Array (two indices) |
| Subarray/substring constraint | Sliding Window | Array + two pointers |
| Sorted + find position | Binary Search | Array or search space |
| Top K elements | Heap | Min/max-heap |
| Prefix matching | Trie | Trie nodes |
| Parentheses / undo | Stack | Stack |
| Level-by-level tree | BFS | Queue + visited |
| Path exists / all paths | DFS/Backtracking | Recursion + state |
| Overlapping intervals | Interval Merge | Sorted intervals |
| Cycle detection | Union-Find | DSU array |
| Dependencies | Topological Sort | DAG + in-degree |
| Optimal substructure | DP | DP array/table |
| Weighted shortest path | Dijkstra | Min-heap + dist array |
| Bit manipulation | Bit Tricks | Integer operations |

### DP Pattern Playbook

| Pattern | Classic Problem | Recurrence |
| ------- | --------------- | ---------- |
| 1D Linear | Climbing Stairs | `dp[i] = dp[i-1] + dp[i-2]` |
| 1D + Choice | House Robber | `dp[i] = max(dp[i-1], dp[i-2]+nums[i])` |
| 2D Grid | Unique Paths | `dp[i][j] = dp[i-1][j] + dp[i][j-1]` |
| 0/1 Knapsack | Partition Equal Subset | `dp[j] = dp[j] or dp[j-nums[i]]` |
| Unbounded | Coin Change | `dp[i] = min(dp[i], dp[i-coin]+1)` |
| LCS | Longest Common Subseq | `dp[i][j] = dp[i-1][j-1]+1 if match` |
| LIS | Longest Increasing Subseq | `dp[i] = max(dp[j]+1) for j<i if arr[j]<arr[i]` |

> **DP Mental Model**: "If I knew the answer to every smaller version of this problem, how would I compute this one?"

---

## 05. The Problem-Solving Framework

### UMPIRE Method

| Step | Action | Time |
| ---- | ------ | ---- |
| **U**nderstand | Restate problem, ask clarifying questions | 2–3 min |
| **M**atch | Identify category, recall patterns | 1–2 min |
| **P**lan | State approach in plain English, get buy-in | 2–3 min |
| **I**mplement | Code cleanly while narrating | 15–25 min |
| **R**eview | Trace with example, find bugs | 3–5 min |
| **E**valuate | State complexity, discuss trade-offs | 2–3 min |

### Clarifying Questions Cheat Sheet

```text
INPUT:
- What is the size range of input? (n = 10? 10^5? 10^9?)
- Can values be negative? Zero? Duplicates?
- Is input sorted?

OUTPUT:
- Return indices, values, or count?
- Multiple valid answers: any one or all?
- What if no valid answer exists?

EDGES:
- Empty input? Single element?
- Memory/time constraints?
```

### Live Interview Scripts

```text
AFTER READING:
"Let me make sure I understand. We're given [restate].
 We need to return [output]. Is that right?"

BEFORE CODING:
"My plan is to use [approach] which gives O([time])
 time and O([space]) space. Does that sound good?"

WHILE CODING:
"I'm initializing a hash map here to track [what].
 I'll iterate through once, checking [condition]..."

AFTER FINISHING:
"Let me trace through with [example]...
 Time complexity is [X] because [reason].
 Space complexity is [Y]."
```

---

## 06. Coding Best Practices

### Naming Conventions

| Context | Bad | Good |
| ------- | --- | ---- |
| Loop index | i, j | left, right, idx |
| Accumulated value | s, t | currentSum, total |
| Hash map | d, m | charCount, memo |
| Boolean flag | b, f | isValid, hasVisited |
| Result | r, ans | maxProfit, minCost |

### Code Structure Signals

```text
✅ Write helper functions for repeated logic
✅ Separate concerns: input, solve, output
✅ Avoid magic numbers: ALPHABET_SIZE = 26
✅ Return early for edge cases
✅ Comment intent, not mechanics
```

### Language Tips

**Python:**

```python
from collections import defaultdict, Counter, deque
from functools import lru_cache
import heapq

# Use defaultdict instead of key checks
freq = defaultdict(int)

# Use deque for O(1) left-pop
q = deque([start])

# heapq is min-heap; negate for max-heap
```

**JavaScript:**

```javascript
// Map/Set for O(1) guaranteed
const freq = new Map();

// Sort with comparator for numbers
arr.sort((a, b) => a - b);
```

---

## 07. System Design Interviews

### RESHADED Framework

| Letter | Step | Key Questions |
| ------ | ---- | ------------- |
| **R** | Requirements | Functional/non-functional? SLAs? |
| **E** | Estimation | QPS, storage/year, bandwidth? |
| **S** | Storage Schema | SQL vs NoSQL? Indexes? |
| **H** | High-Level Design | Core components? Draw diagram. |
| **A** | API Design | REST/GraphQL? Endpoints? |
| **D** | Deep Dive | Bottlenecks? Sharding? Caching? |
| **E** | Evaluate | SPoFs? 10x scale? |
| **D** | Distinctive | Real-time? Geo-distribution? |

### Core Concepts

**Caching:**

| Strategy | Description | Trade-off |
| -------- | ----------- | --------- |
| Write-through | Write cache + DB | Strong consistency, slower |
| Write-back | Write cache only | Fast, risk of data loss |
| Cache-aside | App manages cache | Most flexible |

**Database Scaling:**

| Technique | Description | Use When |
| --------- | ----------- | -------- |
| Read replicas | Copy DB for reads | Read-heavy |
| Sharding | Partition by key | Write-heavy, large data |
| Vertical scaling | More RAM/CPU | Simple, hits ceiling |

### 8 Essential System Designs

1. URL Shortener (Base62, consistent hashing)
2. Twitter Feed (Fan-out on write vs pull)
3. Rate Limiter (Token bucket, Redis)
4. Chat System (WebSockets, message IDs)
5. Search Autocomplete (Trie + frequency)
6. Ride Sharing (Quadtree/S2 cells)
7. Distributed Cache (Consistent hashing)
8. File Storage (Chunk large files)

---

## 08. Behavioral Interviews

### STAR+ Framework

| Element | What to Include | Common Mistake |
| ------- | --------------- | -------------- |
| **S**ituation | Context: company, team, timeline | Over-explaining |
| **T**ask | YOUR specific ownership | Saying "we" |
| **A**ction | 3–5 key things YOU did | Being vague |
| **R**esult | Quantified outcome | "It went well" |
| **+** Learning | What you'd do differently | Skipping |
| **+** Impact | Broader effect | Underselling |

### 8-Story Bank

| Archetype | Questions It Answers |
| --------- | -------------------- |
| Technical deepest | "Most difficult problem" |
| Conflict | "Disagree and commit" |
| Failure | "Biggest mistake" |
| Leadership | "Influenced without authority" |
| Under pressure | "Tight deadline, incident" |
| Mentorship | "Developed a team member" |
| Data-driven | "Changed based on data" |
| Cross-functional | "Worked with non-engineers" |

### Amazon Leadership Principles (16)

| Principle | Core Expectation |
| --------- | ---------------- |
| Customer Obsession | Start with customer |
| Ownership | Act as company owner |
| Invent and Simplify | Innovate, simplify |
| Are Right, A Lot | Strong judgment |
| Learn and Be Curious | Always exploring |
| Hire and Develop | Raise the bar |
| Highest Standards | Never settle |
| Think Big | Bold, long-term |
| Bias for Action | Speed matters |
| Frugality | More with less |
| Earn Trust | Self-critical |
| Dive Deep | Stay close to details |
| Have Backbone | Disagree respectfully |
| Deliver Results | Focus on key inputs |
| Strive to be Earth's Best | Highest quality |
| Success and Scale | Responsible impact |

---

## 09. Live Interview Tactics

### First 5 Minutes

| Minute | Action | Why |
| ------ | ------ | --- |
| 0–1 | Read silently | First impressions |
| 1–2 | Restate problem | Confirm understanding |
| 2–3 | Ask clarifying questions | Prevent wasted time |
| 3–4 | State approach (brute force first) | Invite course correction |
| 4–5 | Get buy-in before coding | Prevent wrong direction |

### Stuck Playbook

```text
1. State where stuck: "I'm trying to avoid nested loop..."
2. Brute force fallback: "Let me start O(n²), then optimize"
3. Concrete example: Draw small example
4. Simplify: "What if only 2 elements?"
5. Ask gracefully: "Would a hint be appropriate?"
```

> **Never let 60 seconds pass in silence. The interviewer judges how you handle being stuck.**

### Red Flags

| Behavior | Signal | Fix |
| -------- | ------ | --- |
| Coding immediately | Impulsive | Spend 3–5 min planning |
| Silent >2 min | Poor communication | Narrate everything |
| "I don't know" + stop | Gives up | "Let me think step by step" |
| No testing | Low quality bar | Always trace examples |
| Arguing hints | Difficult | Embrace: "That's a great point" |

---

## 10. Offer Evaluation & Negotiation

### Total Compensation Components

| Component | What to Check | Negotiable? |
| --------- | ------------- | ----------- |
| Base Salary | Benchmark Levels.fyi | Yes, always |
| Equity (RSUs) | 4-yr vest, 1-yr cliff | Yes - grant size |
| Signing Bonus | Clawback clause? | Yes - highly |
| Annual Bonus | Guaranteed? | Sometimes |
| Refreshes | Annual grants? | Ask policy |
| Level | Affects future trajectory | Yes - advocate |

### Negotiation Scripts

**Deflect salary anchor:**

```text
"I'd rather focus on finding the right fit first.
 What's the band for this role?"
```

**Leverage competing offer:**

```text
"I have an offer from [Company] at [range].
 I'm more excited about [Target] because [reason].
 Is there flexibility to bring closer to [target]?"
```

**Multi-dimensional counter:**

```text
"The base is lower than I targeted—hoping for [X].
 If no flexibility on base, could we look at
 larger equity grant or higher signing bonus?"
```

### Negotiation Principles

```text
✅ Never accept on the spot
✅ Get everything in writing
✅ Know your BATNA (best alternative)
✅ Negotiate level aggressively ($50k+ impact)
✅ Frame collaboratively: "I want to make this work"
✅ Silence is powerful after stating counter
```

---

## Quick Reference

### Complexity Guide

| Complexity | Max n (1 sec) | Common Cause |
| ---------- | ------------- | ------------ |
| O(1) | Any | Hash lookup |
| O(log n) | ~10^10 | Binary search |
| O(n) | ~10^8 | Single loop |
| O(n log n) | ~10^6 | Sort + linear |
| O(n²) | ~10^4 | Nested loops |
| O(2^n) | ~20–25 | Subsets |
| O(n!) | ~12 | Permutations |

### Edge Case Checklist

```text
□ Empty input
□ Single element
□ Two elements
□ All same values
□ Sorted / reverse sorted
□ Negative numbers
□ Duplicates
□ Integer overflow
□ Null values
□ Maximum size
```

---

## Resources

| Category | Resource | Why |
| -------- | -------- | --- |
| Practice | LeetCode (NeetCode order) | Best quality |
| Video | NeetCode YouTube | Visual explanations |
| Book | Cracking the Coding Interview | Patterns + behavioral |
| System Design | System Design Primer | Free, comprehensive |
| Mocks | Pramp.com | Free peer-to-peer |
| Salary | Levels.fyi | Real comp data |

---

> **Success = Preparation × Communication × Mindset**
>
> Every expert was once a beginner. Consistent practice beats cramming.

**You've got this!**
