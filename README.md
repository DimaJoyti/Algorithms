# 🚀 Algorithms & Data Structures Collection

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![JavaScript](https://img.shields.io/badge/JavaScript-ES6+-yellow.svg)](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
[![Python](https://img.shields.io/badge/Python-3.7+-blue.svg)](https://www.python.org/)
[![Tests](https://img.shields.io/badge/Tests-Jest%20%7C%20PyTest-green.svg)](https://jestjs.io/)

A comprehensive collection of algorithms and data structures implemented in JavaScript and Python with detailed explanations, tests, and usage examples.

## 📋 Table of Contents

- [Project Overview](#project-overview)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Algorithm Categories](#algorithm-categories)
- [Installation](#installation)
- [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Project Overview

This repository contains:

- **32+ JavaScript algorithms** with full test coverage
- **1400+ Python algorithms** organized by category (LeetCode solutions)
- Detailed complexity analysis (Big O notation)
- Interactive examples and visualizations
- Comparative analysis of different approaches

### 🎓 Who is this for?

- Students learning algorithms and data structures
- Developers preparing for technical interviews
- Educators and mentors
- Anyone looking to deepen their Computer Science knowledge

## ⚡ Quick Start

### JavaScript

```bash
# Clone the repository
git clone https://github.com/your-username/algorithms.git
cd algorithms/javascript

# Install dependencies
npm install

# Run tests
npm test

# Run specific test
npm test anagrams
```

### Python

```bash
# Navigate to Python directory
cd python

# Install dependencies
pip install -r requirements.txt

# Run tests
pytest

# Run specific file
python BinarySearch.py
```

## 📁 Project Structure

```
algorithms/
├── javascript/                 # JavaScript implementations
│   ├── anagrams/             # Anagram checking
│   ├── sorting/              # Sorting algorithms
│   ├── tree/                 # Trees and traversals
│   ├── linkedlist/           # Linked lists
│   └── ...                   # Other algorithms
├── python/                    # Python implementations (organized by category)
│   ├── arrays/               # 880 array problems
│   ├── strings/              # 227 string problems
│   ├── trees/                # 191 tree problems
│   ├── dynamic_programming/  # 142 DP problems
│   ├── graphs/               # 41 graph problems
│   ├── linked_lists/         # 69 linked list problems
│   ├── math/                 # 40 math problems
│   ├── misc/                 # 133 uncategorized problems
│   └── ...                   # Other categories
├── docs/                      # Documentation
├── tests/                     # Additional tests
└── README.md                  # This file
```

## 🗂️ Algorithm Categories

### 🔤 String Algorithms

- **Anagrams** - Check if two strings are anagrams
- **Palindrome** - Check if string is palindrome
- **Reverse String** - Reverse a string
- **Capitalize** - Capitalize words in string
- **Max Character** - Find most frequent character

### 🔢 Mathematical Algorithms

- **Fibonacci** - Fibonacci numbers (recursion + memoization)
- **FizzBuzz** - Classic FizzBuzz problem
- **Reverse Integer** - Reverse digits of integer
- **Steps/Pyramid** - Print step/pyramid patterns

### 📊 Sorting and Searching

- **Bubble Sort** - Bubble sort O(n²)
- **Selection Sort** - Selection sort O(n²)
- **Merge Sort** - Merge sort O(n log n)
- **Binary Search** - Binary search O(log n)

### 🔗 Data Structures

- **Linked Lists** - Linked list implementations
- **Stacks** - Stack data structure
- **Queues** - Queue data structure
- **Trees** - Tree structures (BFS, DFS)
- **Binary Search Trees** - BST implementations

### 🌐 Graphs and Trees

- **Tree Traversal** - Tree traversal algorithms
- **Level Width** - Calculate tree level widths
- **Validate BST** - BST validation
- **Graph Algorithms** - Graph algorithms

### 🧮 Dynamic Programming

- **Fibonacci DP** - Fibonacci with memoization
- **Climbing Stairs** - Staircase climbing problem
- **Coin Change** - Coin change problem
- **Longest Subsequence** - Longest common subsequence

## 🛠️ Installation

### System Requirements

- **Node.js** 14+ for JavaScript
- **Python** 3.7+ for Python
- **Git** for cloning repository

### Step by Step

1. **Clone Repository**

```bash
git clone https://github.com/your-username/algorithms.git
cd algorithms
```

2. **JavaScript Setup**

```bash
cd javascript
npm install
npm test
```

3. **Python Setup**

```bash
cd ../python
pip install -r requirements.txt
pytest
```

## 🧪 Running Tests

### JavaScript Tests

```bash
# All tests
npm test

# Specific algorithm
npm test anagrams
npm test sorting
npm test tree

# With verbose output
npm test -- --verbose
```

### Python Tests

```bash
# All tests
pytest

# Specific file
pytest test_binary_search.py

# With coverage
pytest --cov=.
```

## 📚 Usage Examples

### JavaScript Example

```javascript
const anagrams = require('./anagrams');
const { bubbleSort } = require('./sorting');

// Check anagrams
console.log(anagrams('listen', 'silent')); // true

// Sort array
const arr = [64, 34, 25, 12, 22, 11, 90];
console.log(bubbleSort(arr)); // [11, 12, 22, 25, 34, 64, 90]
```

### Python Example

```python
from BinarySearch import binarySearch
from Fibonacci import fib

# Binary search
arr = [1, 3, 5, 7, 9, 11]
result = binarySearch(arr, 0, len(arr)-1, 7)
print(f"Element found at position: {result}")

# Fibonacci numbers
print(f"10th Fibonacci number: {fib(10)}")
```

## 🤝 Contributing

We welcome contributions from the community! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

### How to Help

1. 🐛 **Found a bug?** Create an issue
2. 💡 **Have an idea?** Discuss in Discussions
3. 🔧 **Want to add algorithm?** Create a Pull Request
4. 📖 **Improve documentation?** Always welcome!

### Quick Start for Contributors

```bash
# Fork the repository
git clone https://github.com/your-username/algorithms.git
cd algorithms

# Create new branch
git checkout -b feature/new-algorithm

# Make changes and add tests
# ...

# Commit changes
git commit -m "Add: new sorting algorithm"

# Push changes
git push origin feature/new-algorithm

# Create Pull Request
```

## 📊 Project Statistics

- **JavaScript**: 32 algorithms with tests
- **Python**: 1400+ algorithms organized in 20 categories
- **Test Coverage**: 95%+
- **Supported Languages**: JavaScript, Python
- **Planned**: Go, Rust, Java

## 🔗 Useful Links

- [LeetCode](https://leetcode.com/) - Practice platform
- [HackerRank](https://www.hackerrank.com/) - Coding challenges
- [GeeksforGeeks](https://www.geeksforgeeks.org/) - Algorithms and data structures
- [Visualgo](https://visualgo.net/) - Algorithm visualizations
- [Big O Cheat Sheet](https://www.bigocheatsheet.com/) - Complexity reference

## 📄 License

This project is licensed under the [MIT License](LICENSE) - see the LICENSE file for details.

## 🙏 Acknowledgments

- Thanks to all contributors for their contributions
- Inspiration from the algorithms community
- Special thanks to Computer Science educators

---

**⭐ If this project was helpful, please give it a star!**

Created with ❤️ for the developer community
