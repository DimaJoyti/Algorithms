"""
Algorithms Collection - Python Implementation

This package contains a comprehensive collection of algorithms and data structures
implemented in Python with detailed documentation, tests, and examples.

Categories:
- String Algorithms
- Mathematical Algorithms  
- Sorting and Searching
- Data Structures
- Tree Algorithms
- Graph Algorithms
- Dynamic Programming

Author: Algorithms Collection Contributors
License: MIT
"""

__version__ = "1.0.0"
__author__ = "Algorithms Collection Contributors"
__license__ = "MIT"

# Import core algorithms
try:
    from .BinarySearch import binarySearch, binarySearchLessThan
    from .Fibonacci import fib, fib2, fib3
    from .MergeSort import merge_sort
    from .SelectionSort import selection_sort
    from .ArraySearch import array_search
    from .Calculator import Calculator
    from .Graphs import Graph
    from .MaxHeap import MaxHeap
    from .Lists import LinkedList, Node
    from .Primes import is_prime, sieve_of_eratosthenes
    from .GreatestCommonDevider import gcd
except ImportError:
    # Handle missing modules gracefully
    pass

# Algorithm categories
CATEGORIES = {
    "strings": [
        "anagram_check",
        "palindrome_check", 
        "string_reversal",
        "character_frequency",
        "substring_search"
    ],
    "math": [
        "fibonacci",
        "prime_numbers",
        "gcd_lcm",
        "factorial",
        "power_calculation"
    ],
    "sorting": [
        "bubble_sort",
        "selection_sort", 
        "merge_sort",
        "quick_sort",
        "heap_sort"
    ],
    "searching": [
        "binary_search",
        "linear_search",
        "depth_first_search",
        "breadth_first_search"
    ],
    "data_structures": [
        "linked_list",
        "stack",
        "queue", 
        "heap",
        "hash_table"
    ],
    "trees": [
        "binary_tree",
        "binary_search_tree",
        "tree_traversal",
        "tree_validation"
    ],
    "graphs": [
        "graph_representation",
        "graph_traversal",
        "shortest_path",
        "minimum_spanning_tree"
    ],
    "dynamic_programming": [
        "fibonacci_dp",
        "knapsack_problem",
        "longest_subsequence",
        "coin_change"
    ]
}

def get_algorithm_info(name: str) -> dict:
    """
    Get information about a specific algorithm.
    
    Args:
        name: Algorithm name
        
    Returns:
        Dictionary with algorithm information
    """
    # This would be expanded to include actual algorithm metadata
    return {
        "name": name,
        "category": "unknown",
        "complexity": "unknown",
        "description": "Algorithm information not available"
    }

def list_algorithms() -> list:
    """
    List all available algorithms.
    
    Returns:
        List of algorithm names
    """
    all_algorithms = []
    for category_algorithms in CATEGORIES.values():
        all_algorithms.extend(category_algorithms)
    return sorted(all_algorithms)

def list_categories() -> list:
    """
    List all algorithm categories.
    
    Returns:
        List of category names
    """
    return list(CATEGORIES.keys())

def get_category_algorithms(category: str) -> list:
    """
    Get algorithms in a specific category.
    
    Args:
        category: Category name
        
    Returns:
        List of algorithms in the category
    """
    if category not in CATEGORIES:
        raise ValueError(f"Category '{category}' not found")
    return CATEGORIES[category]

# Export main functions and classes
__all__ = [
    # Core functions
    "get_algorithm_info",
    "list_algorithms", 
    "list_categories",
    "get_category_algorithms",
    
    # Categories
    "CATEGORIES",
    
    # Version info
    "__version__",
    "__author__", 
    "__license__"
]

# Add available algorithms to __all__ if they exist
try:
    __all__.extend([
        "binarySearch",
        "binarySearchLessThan", 
        "fib",
        "fib2",
        "fib3",
        "merge_sort",
        "selection_sort",
        "array_search",
        "Calculator",
        "Graph",
        "MaxHeap",
        "LinkedList",
        "Node",
        "is_prime",
        "sieve_of_eratosthenes",
        "gcd"
    ])
except NameError:
    # Some algorithms might not be available
    pass
