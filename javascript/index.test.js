/**
 * Tests for the main index.js file
 * Tests the algorithm collection exports and utility functions
 */

const {
  // Individual algorithms
  anagrams,
  palindrome,
  reversestring,
  capitalize,
  maxchar,
  vowels,
  fib,
  fizzbuzz,
  reverseint,
  steps,
  pyramid,
  bubbleSort,
  selectionSort,
  mergeSort,
  merge,
  LinkedList,
  Queue,
  Stack,
  Tree,
  levelwidth,
  midpoint,
  fromlast,
  circular,
  qfroms,
  weave,
  chunk,
  matrix,
  events,
  
  // Organized collections
  algorithms,
  
  // Utility functions
  getAlgorithm,
  listAlgorithms,
  getCategory,
  listCategories,
} = require('./index');

describe('Main Index Exports', () => {
  describe('Individual Algorithm Exports', () => {
    test('exports string algorithms', () => {
      expect(typeof anagrams).toBe('function');
      expect(typeof palindrome).toBe('function');
      expect(typeof reversestring).toBe('function');
      expect(typeof capitalize).toBe('function');
      expect(typeof maxchar).toBe('function');
      expect(typeof vowels).toBe('function');
    });

    test('exports math algorithms', () => {
      expect(typeof fib).toBe('function');
      expect(typeof fizzbuzz).toBe('function');
      expect(typeof reverseint).toBe('function');
      expect(typeof steps).toBe('function');
      expect(typeof pyramid).toBe('function');
    });

    test('exports sorting algorithms', () => {
      expect(typeof bubbleSort).toBe('function');
      expect(typeof selectionSort).toBe('function');
      expect(typeof mergeSort).toBe('function');
      expect(typeof merge).toBe('function');
    });

    test('exports data structures', () => {
      expect(typeof LinkedList).toBe('function');
      expect(typeof Queue).toBe('function');
      expect(typeof Stack).toBe('function');
      expect(typeof Tree).toBe('function');
    });

    test('exports tree algorithms', () => {
      expect(typeof levelwidth).toBe('function');
    });

    test('exports linked list algorithms', () => {
      expect(typeof midpoint).toBe('function');
      expect(typeof fromlast).toBe('function');
      expect(typeof circular).toBe('function');
    });

    test('exports queue/stack algorithms', () => {
      expect(typeof qfroms).toBe('function');
      expect(typeof weave).toBe('function');
    });

    test('exports array algorithms', () => {
      expect(typeof chunk).toBe('function');
      expect(typeof matrix).toBe('function');
    });

    test('exports utility functions', () => {
      expect(typeof events).toBe('function');
    });
  });

  describe('Organized Algorithm Collections', () => {
    test('exports algorithms object with correct structure', () => {
      expect(typeof algorithms).toBe('object');
      expect(algorithms).toHaveProperty('strings');
      expect(algorithms).toHaveProperty('math');
      expect(algorithms).toHaveProperty('sorting');
      expect(algorithms).toHaveProperty('dataStructures');
      expect(algorithms).toHaveProperty('trees');
      expect(algorithms).toHaveProperty('linkedLists');
      expect(algorithms).toHaveProperty('queuesStacks');
      expect(algorithms).toHaveProperty('arrays');
      expect(algorithms).toHaveProperty('utilities');
    });

    test('string algorithms category contains correct algorithms', () => {
      const { strings } = algorithms;
      expect(strings).toHaveProperty('anagrams');
      expect(strings).toHaveProperty('palindrome');
      expect(strings).toHaveProperty('reversestring');
      expect(strings).toHaveProperty('capitalize');
      expect(strings).toHaveProperty('maxchar');
      expect(strings).toHaveProperty('vowels');
    });

    test('math algorithms category contains correct algorithms', () => {
      const { math } = algorithms;
      expect(math).toHaveProperty('fib');
      expect(math).toHaveProperty('fizzbuzz');
      expect(math).toHaveProperty('reverseint');
      expect(math).toHaveProperty('steps');
      expect(math).toHaveProperty('pyramid');
    });

    test('sorting algorithms category contains correct algorithms', () => {
      const { sorting } = algorithms;
      expect(sorting).toHaveProperty('bubbleSort');
      expect(sorting).toHaveProperty('selectionSort');
      expect(sorting).toHaveProperty('mergeSort');
      expect(sorting).toHaveProperty('merge');
    });
  });

  describe('Utility Functions', () => {
    describe('getAlgorithm', () => {
      test('returns algorithm by name', () => {
        const algorithm = getAlgorithm('anagrams');
        expect(typeof algorithm).toBe('function');
        expect(algorithm).toBe(anagrams);
      });

      test('returns algorithm from different categories', () => {
        expect(getAlgorithm('bubbleSort')).toBe(bubbleSort);
        expect(getAlgorithm('fib')).toBe(fib);
        expect(getAlgorithm('LinkedList')).toBe(LinkedList);
      });

      test('throws error for non-existent algorithm', () => {
        expect(() => getAlgorithm('nonExistentAlgorithm')).toThrow(
          "Algorithm 'nonExistentAlgorithm' not found"
        );
      });
    });

    describe('listAlgorithms', () => {
      test('returns array of algorithm names', () => {
        const algorithmList = listAlgorithms();
        expect(Array.isArray(algorithmList)).toBe(true);
        expect(algorithmList.length).toBeGreaterThan(0);
      });

      test('returns sorted list', () => {
        const algorithmList = listAlgorithms();
        const sortedList = [...algorithmList].sort();
        expect(algorithmList).toEqual(sortedList);
      });

      test('includes expected algorithms', () => {
        const algorithmList = listAlgorithms();
        expect(algorithmList).toContain('anagrams');
        expect(algorithmList).toContain('bubbleSort');
        expect(algorithmList).toContain('fib');
        expect(algorithmList).toContain('LinkedList');
      });
    });

    describe('getCategory', () => {
      test('returns category by name', () => {
        const stringsCategory = getCategory('strings');
        expect(typeof stringsCategory).toBe('object');
        expect(stringsCategory).toHaveProperty('anagrams');
      });

      test('returns different categories correctly', () => {
        expect(getCategory('math')).toHaveProperty('fib');
        expect(getCategory('sorting')).toHaveProperty('bubbleSort');
        expect(getCategory('dataStructures')).toHaveProperty('LinkedList');
      });

      test('throws error for non-existent category', () => {
        expect(() => getCategory('nonExistentCategory')).toThrow(
          "Category 'nonExistentCategory' not found"
        );
      });
    });

    describe('listCategories', () => {
      test('returns array of category names', () => {
        const categories = listCategories();
        expect(Array.isArray(categories)).toBe(true);
        expect(categories.length).toBeGreaterThan(0);
      });

      test('includes expected categories', () => {
        const categories = listCategories();
        expect(categories).toContain('strings');
        expect(categories).toContain('math');
        expect(categories).toContain('sorting');
        expect(categories).toContain('dataStructures');
        expect(categories).toContain('trees');
        expect(categories).toContain('linkedLists');
        expect(categories).toContain('queuesStacks');
        expect(categories).toContain('arrays');
        expect(categories).toContain('utilities');
      });
    });
  });

  describe('Integration Tests', () => {
    test('algorithms can be accessed through different methods', () => {
      // Direct export
      expect(anagrams('listen', 'silent')).toBe(true);
      
      // Through getAlgorithm
      const anagramsFunc = getAlgorithm('anagrams');
      expect(anagramsFunc('listen', 'silent')).toBe(true);
      
      // Through category
      const stringsCategory = getCategory('strings');
      expect(stringsCategory.anagrams('listen', 'silent')).toBe(true);
    });

    test('data structures can be instantiated', () => {
      const list = new LinkedList();
      expect(list).toBeInstanceOf(LinkedList);
      
      const queue = new Queue();
      expect(queue).toBeInstanceOf(Queue);
      
      const stack = new Stack();
      expect(stack).toBeInstanceOf(Stack);
    });
  });
});
