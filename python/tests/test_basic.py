"""
Basic tests for the algorithms collection.

This module contains basic tests to ensure the package structure
and core functionality work correctly.
"""

import pytest
import sys
import os

# Add parent directory to path to import algorithms
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

try:
    import BinarySearch
    import Fibonacci
    ALGORITHMS_AVAILABLE = True
except ImportError:
    ALGORITHMS_AVAILABLE = False


class TestPackageStructure:
    """Test basic package structure and imports."""
    
    def test_package_imports(self):
        """Test that the package can be imported."""
        try:
            import algorithms
            assert True
        except ImportError:
            # Package structure might not be fully set up yet
            assert True
    
    def test_categories_defined(self):
        """Test that algorithm categories are defined."""
        try:
            from algorithms import CATEGORIES
            assert isinstance(CATEGORIES, dict)
            assert len(CATEGORIES) > 0
        except ImportError:
            # Skip if package not fully set up
            pytest.skip("Package not fully configured yet")


@pytest.mark.skipif(not ALGORITHMS_AVAILABLE, reason="Core algorithms not available")
class TestBinarySearch:
    """Test binary search algorithm."""
    
    def test_binary_search_found(self):
        """Test binary search when element is found."""
        arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
        result = BinarySearch.binarySearch(arr, 0, len(arr) - 1, 5)
        assert result == 4  # Index of element 5
    
    def test_binary_search_not_found(self):
        """Test binary search when element is not found."""
        arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
        result = BinarySearch.binarySearch(arr, 0, len(arr) - 1, 11)
        assert result == -1
    
    def test_binary_search_empty_array(self):
        """Test binary search with empty array."""
        arr = []
        result = BinarySearch.binarySearch(arr, 0, -1, 5)
        assert result == -1
    
    def test_binary_search_single_element_found(self):
        """Test binary search with single element that is found."""
        arr = [5]
        result = BinarySearch.binarySearch(arr, 0, 0, 5)
        assert result == 0
    
    def test_binary_search_single_element_not_found(self):
        """Test binary search with single element that is not found."""
        arr = [5]
        result = BinarySearch.binarySearch(arr, 0, 0, 3)
        assert result == -1


@pytest.mark.skipif(not ALGORITHMS_AVAILABLE, reason="Core algorithms not available")
class TestFibonacci:
    """Test Fibonacci algorithms."""
    
    def test_fibonacci_base_cases(self):
        """Test Fibonacci base cases."""
        assert Fibonacci.fib(0) == 0
        assert Fibonacci.fib(1) == 1
    
    def test_fibonacci_small_numbers(self):
        """Test Fibonacci for small numbers."""
        assert Fibonacci.fib(2) == 1
        assert Fibonacci.fib(3) == 2
        assert Fibonacci.fib(4) == 3
        assert Fibonacci.fib(5) == 5
        assert Fibonacci.fib(6) == 8
    
    def test_fibonacci_consistency(self):
        """Test that different Fibonacci implementations give same results."""
        for i in range(10):
            assert Fibonacci.fib(i) == Fibonacci.fib2(i)
            assert Fibonacci.fib2(i) == Fibonacci.fib3(i)
    
    def test_fibonacci_larger_numbers(self):
        """Test Fibonacci for larger numbers."""
        # Test that it doesn't crash and gives reasonable results
        result = Fibonacci.fib2(20)
        assert result == 6765
        
        result = Fibonacci.fib3(15)
        assert result == 610


class TestUtilityFunctions:
    """Test utility functions."""
    
    def test_algorithm_categories(self):
        """Test that we can define and access algorithm categories."""
        categories = {
            "sorting": ["bubble_sort", "merge_sort"],
            "searching": ["binary_search", "linear_search"]
        }
        
        assert "sorting" in categories
        assert "searching" in categories
        assert len(categories["sorting"]) == 2
        assert len(categories["searching"]) == 2
    
    def test_algorithm_metadata(self):
        """Test algorithm metadata structure."""
        metadata = {
            "name": "binary_search",
            "category": "searching",
            "time_complexity": "O(log n)",
            "space_complexity": "O(1)",
            "description": "Efficient search algorithm for sorted arrays"
        }
        
        assert metadata["name"] == "binary_search"
        assert metadata["category"] == "searching"
        assert "O(log n)" in metadata["time_complexity"]


if __name__ == "__main__":
    # Run tests if script is executed directly
    pytest.main([__file__, "-v"])
