# JavaScript Algorithms by Category

This directory contains JavaScript algorithms organized by categories for better navigation and understanding.

## 📁 Directory Structure

```
categories/
├── strings/              # String manipulation algorithms
├── math/                 # Mathematical algorithms
├── sorting/              # Sorting algorithms
├── data-structures/      # Data structure implementations
├── trees/                # Tree-related algorithms
└── dynamic-programming/  # Dynamic programming solutions
```

## 🔤 String Algorithms

Algorithms for string manipulation and analysis:

- **Anagrams** - Check if two strings are anagrams
- **Palindrome** - Check if a string is a palindrome
- **Reverse String** - Reverse a string
- **Capitalize** - Capitalize words in a string
- **Max Character** - Find most frequent character
- **Vowels** - Count vowels in a string

**Location**: `strings/`

## 🔢 Mathematical Algorithms

Mathematical computations and number theory:

- **Fibonacci** - Generate Fibonacci numbers
- **FizzBuzz** - Classic FizzBuzz problem
- **Reverse Integer** - Reverse digits of an integer
- **Steps** - Print step pattern
- **Pyramid** - Print pyramid pattern

**Location**: `math/`

## 📊 Sorting Algorithms

Various sorting algorithm implementations:

- **Bubble Sort** - O(n²) comparison-based sorting
- **Selection Sort** - O(n²) selection-based sorting  
- **Merge Sort** - O(n log n) divide-and-conquer sorting
- **Quick Sort** - O(n log n) average case sorting

**Location**: `sorting/`

## 🔗 Data Structures

Fundamental data structure implementations:

- **Linked List** - Singly linked list with operations
- **Stack** - LIFO data structure
- **Queue** - FIFO data structure
- **Hash Table** - Key-value storage
- **Heap** - Priority queue implementation

**Location**: `data-structures/`

## 🌳 Tree Algorithms

Tree data structures and algorithms:

- **Binary Tree** - Basic binary tree implementation
- **Binary Search Tree** - BST with search operations
- **Tree Traversal** - DFS and BFS traversals
- **Tree Validation** - Validate BST properties
- **Level Width** - Calculate width at each level

**Location**: `trees/`

## 🧮 Dynamic Programming

Optimization problems using dynamic programming:

- **Fibonacci DP** - Fibonacci with memoization
- **Coin Change** - Minimum coins for amount
- **Longest Subsequence** - Find longest common subsequence
- **Knapsack** - 0/1 knapsack problem

**Location**: `dynamic-programming/`

## 📝 Algorithm Template

When adding new algorithms, follow this structure:

```javascript
// --- Directions
// Clear description of what the algorithm does
// --- Examples
//   algorithmName(input) --> expected output
// --- Complexity
//   Time: O(?)
//   Space: O(?)

function algorithmName(input) {
  // Implementation here
}

module.exports = algorithmName;
```

## 🧪 Testing Template

Each algorithm should have comprehensive tests:

```javascript
const algorithmName = require('./index');

test('algorithm function exists', () => {
  expect(typeof algorithmName).toEqual('function');
});

test('handles basic case', () => {
  expect(algorithmName(input)).toEqual(expectedOutput);
});

test('handles edge cases', () => {
  expect(algorithmName(edgeInput)).toEqual(edgeOutput);
});
```

## 🚀 Usage Examples

```javascript
// Import specific algorithm
const { bubbleSort } = require('../sorting');

// Import entire category
const stringAlgorithms = require('../strings');

// Use algorithm
const sortedArray = bubbleSort([3, 1, 4, 1, 5]);
const isAnagram = stringAlgorithms.anagrams('listen', 'silent');
```

## 📚 Learning Resources

- [Big O Cheat Sheet](https://www.bigocheatsheet.com/)
- [Visualgo](https://visualgo.net/) - Algorithm visualizations
- [LeetCode](https://leetcode.com/) - Practice problems
- [GeeksforGeeks](https://www.geeksforgeeks.org/) - Algorithm explanations

## 🤝 Contributing

When adding new algorithms:

1. Choose the appropriate category
2. Follow the naming conventions
3. Include comprehensive tests
4. Add documentation and examples
5. Update the category README

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for detailed guidelines.
