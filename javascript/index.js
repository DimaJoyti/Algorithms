/**
 * Algorithms Collection - JavaScript Implementation
 *
 * This file exports all available algorithms organized by categories.
 * Each algorithm includes implementation, tests, and documentation.
 *
 * @author Algorithms Collection Contributors
 * @license MIT
 */

// String Algorithms
const anagrams = require('./anagrams');
const palindrome = require('./palindrome');
const reversestring = require('./reversestring');
const capitalize = require('./capitalize');
const maxchar = require('./maxchar');
const vowels = require('./vowels');

// Mathematical Algorithms
const fib = require('./fib');
const fizzbuzz = require('./fizzbuzz');
const reverseint = require('./reverseint');
const steps = require('./steps');
const pyramid = require('./pyramid');

// Sorting Algorithms
const { bubbleSort, selectionSort, mergeSort, merge } = require('./sorting');
const { quickSort, quickSortRandomized, quickSortPure } = require('./quicksort');
const { heapSort, buildMaxHeap, extractMax, insertIntoHeap } = require('./heapsort');

// Search Algorithms
const { binarySearch, binarySearchRecursive, searchRotatedArray } = require('./binarysearch');

// Graph Algorithms
const { dijkstra, shortestPath, allPairsShortestPath } = require('./dijkstra');
const { astar, manhattanDistance, euclideanDistance } = require('./astar');

// Machine Learning Algorithms
const { LinearRegression, KMeans, KNearestNeighbors } = require('./machinelearning');

// Data Structures
const { LinkedList, Node: LinkedListNode } = require('./linkedlist');
const Queue = require('./queue');
const Stack = require('./stack');
const { Tree, Node: TreeNode } = require('./tree');

// Tree Algorithms
const { Node: BSTNode } = require('./bst');
const { Node: ValidateNode } = require('./validate');
const levelwidth = require('./levelwidth');

// Linked List Algorithms
const midpoint = require('./midpoint');
const fromlast = require('./fromlast');
const circular = require('./circular');

// Queue/Stack Algorithms
const qfroms = require('./qfroms');
const weave = require('./weave');

// Array Algorithms
const chunk = require('./chunk');
const matrix = require('./matrix');

// Event System
const events = require('./events');

/**
 * String Algorithms Category
 */
const stringAlgorithms = {
  anagrams,
  palindrome,
  reversestring,
  capitalize,
  maxchar,
  vowels,
};

/**
 * Mathematical Algorithms Category
 */
const mathAlgorithms = {
  fib,
  fizzbuzz,
  reverseint,
  steps,
  pyramid,
};

/**
 * Sorting Algorithms Category
 */
const sortingAlgorithms = {
  bubbleSort,
  selectionSort,
  mergeSort,
  merge,
  quickSort,
  quickSortRandomized,
  quickSortPure,
  heapSort,
  buildMaxHeap,
  extractMax,
  insertIntoHeap,
};

/**
 * Search Algorithms Category
 */
const searchAlgorithms = {
  binarySearch,
  binarySearchRecursive,
  searchRotatedArray,
};

/**
 * Graph Algorithms Category
 */
const graphAlgorithms = {
  dijkstra,
  shortestPath,
  allPairsShortestPath,
  astar,
  manhattanDistance,
  euclideanDistance,
};

/**
 * Machine Learning Algorithms Category
 */
const machineLearningAlgorithms = {
  LinearRegression,
  KMeans,
  KNearestNeighbors,
};

/**
 * Data Structures Category
 */
const dataStructures = {
  LinkedList,
  LinkedListNode,
  Queue,
  Stack,
  Tree,
  TreeNode,
  BSTNode,
  ValidateNode,
};

/**
 * Tree Algorithms Category
 */
const treeAlgorithms = {
  levelwidth,
  // Add more tree algorithms here
};

/**
 * Linked List Algorithms Category
 */
const linkedListAlgorithms = {
  midpoint,
  fromlast,
  circular,
};

/**
 * Queue/Stack Algorithms Category
 */
const queueStackAlgorithms = {
  qfroms,
  weave,
};

/**
 * Array Algorithms Category
 */
const arrayAlgorithms = {
  chunk,
  matrix,
};

/**
 * Utility Functions Category
 */
const utilities = {
  events,
};

/**
 * All algorithms organized by category
 */
const algorithms = {
  strings: stringAlgorithms,
  math: mathAlgorithms,
  sorting: sortingAlgorithms,
  search: searchAlgorithms,
  graphs: graphAlgorithms,
  machineLearning: machineLearningAlgorithms,
  dataStructures,
  trees: treeAlgorithms,
  linkedLists: linkedListAlgorithms,
  queuesStacks: queueStackAlgorithms,
  arrays: arrayAlgorithms,
  utilities,
};

/**
 * Get algorithm by name
 * @param {string} name - Algorithm name
 * @returns {Function|Object} Algorithm implementation
 */
function getAlgorithm(name) {
  for (const category of Object.values(algorithms)) {
    if (category[name]) {
      return category[name];
    }
  }
  throw new Error(`Algorithm '${name}' not found`);
}

/**
 * List all available algorithms
 * @returns {Array<string>} Array of algorithm names
 */
function listAlgorithms() {
  const allAlgorithms = [];
  for (const category of Object.values(algorithms)) {
    allAlgorithms.push(...Object.keys(category));
  }
  return allAlgorithms.sort();
}

/**
 * Get algorithms by category
 * @param {string} category - Category name
 * @returns {Object} Algorithms in the category
 */
function getCategory(category) {
  if (!algorithms[category]) {
    throw new Error(`Category '${category}' not found`);
  }
  return algorithms[category];
}

/**
 * List all categories
 * @returns {Array<string>} Array of category names
 */
function listCategories() {
  return Object.keys(algorithms);
}

// Export everything
module.exports = {
  // Individual algorithms
  ...stringAlgorithms,
  ...mathAlgorithms,
  ...sortingAlgorithms,
  ...searchAlgorithms,
  ...graphAlgorithms,
  ...machineLearningAlgorithms,
  ...dataStructures,
  ...treeAlgorithms,
  ...linkedListAlgorithms,
  ...queueStackAlgorithms,
  ...arrayAlgorithms,
  ...utilities,

  // Organized by category
  algorithms,

  // Utility functions
  getAlgorithm,
  listAlgorithms,
  getCategory,
  listCategories,
};
