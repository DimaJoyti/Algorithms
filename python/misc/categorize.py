#!/usr/bin/env python3
import os
import shutil

CATEGORIES = {
    'arrays': [
        'array', 'sort', 'rotate', 'duplicate', 'missing', 'max', 'min',
        'subset', 'sequence', 'subarray', 'kth', 'majority', 'product',
        'merge', 'partition', 'find', 'search', 'two-sum', 'three-sum',
        'four-sum', 'container', 'trapping', 'jump', 'spiral', 'diagonal',
        'transverse', 'sum', 'water', 'game', 'flip', 'climbing', '买入',
        'stock', 'candy', 'wiggle', 'odd', 'even', 'score', 'average',
        'h-index', 'pascals', 'triangle', 'fizz', 'buzz', 'integer',
        'continuous', 'consecutive', 'shortest', 'longest', 'subsequence',
        'submatrix', 'wave', 'anti-diagonal', 'arranging', 'building',
        'daily', 'descent', 'elevator', 'final', 'grid', 'image', 'increasing'
    ],
    'strings': [
        'string', 'word', 'char', 'text', 'sentence', 'paragraph',
        'anagram', 'palindrome', 'reverse', 'encode', 'decode', 'compress',
        'pattern', 'matching', 'regex', 'longest', 'shortest', 'unique',
        'count', 'replace', 'convert', 'transform', 'zigzag', 'implement',
        'strstr', 'atoi', 'atoi', 'roman', 'valid', 'rabbit', 'judge',
        'license', 'break', 'compare', 'concatenate', 'read', 'write',
        'keyboard', 'keyboard', 'keyboard'
    ],
    'linked_lists': [
        'linked', 'list', 'node', 'cycle', 'reverse', 'add', 'delete',
        'merge', 'sort', 'rotate', 'copy', 'random', 'flatten'
    ],
    'trees': [
        'tree', 'binary', 'bst', 'trie', 'n-ary', 'root', 'leaf', 'node',
        'traversal', 'inorder', 'preorder', 'postorder', 'level', 'order',
        'width', 'depth', 'height', 'balance', 'serialize', 'deserialize',
        'construct', 'insert', 'delete', 'search', 'find', 'lowest',
        'common', 'ancestor', 'path', 'sum', 'max', 'min', 'average',
        'mirror', 'symmetric', 'identical', 'complete', 'perfect',
        'full', 'subtree', 'prune', 'cameras', 'boundary', 'leaf', 'buddy',
        'all-nodes', 'distance', 'closest', 'valid', 'recover', 'trim'
    ],
    'graphs': [
        'graph', 'network', 'node', 'edge', 'vertex', 'bfs', 'dfs',
        'traverse', 'path', 'route', 'shortest', 'longest', 'cycle',
        'connected', 'component', 'island', 'bridge', 'articulation',
        'topological', 'order', 'dependency', 'course', 'schedule',
        'alien', 'evaluate', 'division', 'clone', 'connect', 'population',
        'number-of', 'provinces', 'accounts', 'roads', 'flights', 'cheapest',
        'min-cost', 'max-area', 'regions', 'cut', 'flood', 'surrounded'
    ],
    'dynamic_programming': [
        'dp', 'dynamic', 'climbing', 'house', 'robber', 'coin', 'change',
        'knapsack', 'partition', 'edit', 'distance', 'word', 'break',
        'decode', 'ways', 'combination', 'subsequence', 'subarray',
        'string', 'number', 'integer', 'triangle', 'matrix', 'stair',
        'fibonacci', 'gold', 'mine', 'egg', 'drop', 'burglar', 'stock',
        'buy', 'sell', 'transaction', 'fee', 'cooldown', 'III', 'IV',
        'v', 'profit', 'max-profit', 'dungeon', 'cherry', 'pickup',
        'predict', 'winner', 'guess', 'word', 'stone', 'game', 'predict',
        'interleaving', 'scramble', 'strange', 'regular', 'matching',
        'parenthesis', 'brackets', 'arrange', 'balloon', 'burst', 'word-break'
    ],
    'math': [
        'math', 'add', 'multiply', 'divide', 'subtract', 'power', 'sqrt',
        'prime', 'factorial', 'fibonacci', 'gcd', 'lcm', 'mod', 'modulo',
        'digit', 'number', 'integer', 'fraction', 'decimal', 'roman',
        'convert', 'base', 'binary', 'hex', 'octal', 'random', 'round',
        'angle', 'degree', 'radian', 'sin', 'cos', 'tan', 'circle', 'point',
        'coordinate', 'distance', 'area', 'volume', 'perimeter', 'surface',
        'complex', 'imaginary', 'rational', 'ugly', 'happy', 'perfect',
        'armstrong', 'self', 'dividing', 'multiplying', 'adding', 'subtracting'
    ],
    'bit_manipulation': [
        'bit', 'binary', 'xor', 'shift', 'mask', 'flip', 'set', 'clear',
        'test', 'count', 'ones', 'zeros', 'complement', 'range', 'single',
        'missing', 'duplicate', 'power', 'divide'
    ],
    'backtracking': [
        'backtrack', 'permutation', 'combination', 'subset', 'partition',
        'solve', 'n-queens', 'sudoku', 'word', 'search', 'pattern', 'generate',
        'parenthesis', 'bracket', 'brace', 'telephone', 'letter', 'combination-sum',
        'restore', 'abbreviation', 'abbreviate'
    ],
    'stack': [
        'stack', 'push', 'pop', 'peek', 'top', 'min', 'max', 'valid',
        'parentheses', 'brackets', 'tags', 'html', 'markdown', 'score',
        'evaluate', 'reverse', 'polish', 'notation', 'rpn', 'flatten'
    ],
    'queue': [
        'queue', 'deque', 'circular', 'first', 'last', 'moving', 'window',
        'sliding', 'average', 'max', 'min'
    ],
    'heap': [
        'heap', 'priority', 'queue', 'kth', 'largest', 'smallest', 'top',
        'frequency', 'median', 'sliding', 'window', 'moving'
    ],
    'design': [
        'design', 'implement', 'class', 'data', 'structure', 'twitter',
        'facebook', 'instagram', 'linkedin', 'tweet', 'post', 'follow',
        'friend', 'timeline', 'feed', 'cache', 'lru', 'lfu', 'redis',
        'database', 'storage', 'file', 'system', 'parking', 'tic', 'tac',
        'toe', 'phone', 'directory', 'authentication', 'manager', 'token',
        'logger', 'logger', 'counter', 'hit', 'counter', 'moving', 'average',
        'min-stack', 'max-stack', 'queue-with', 'stack-with', 'hashmap',
        'hashset', 'linked-list', 'array', 'string', 'serialize', 'deserialize',
        'binary-tree', 'expression', 'trie', 'prefix', 'tree', 'word',
        'search-autocomplete', 'ocomplete', 'movie', 'rental', 'log',
        'excel', 'sum', 'formula'
    ],
    'greedy': [
        'greedy', 'assign', 'candy', 'task', 'schedule', 'interval',
        'meeting', 'jump', 'gas', 'station', 'minimum', 'maximum', 'best',
        'optimal', 'profit', 'strategy', 'two', 'keys', 'keyboard'
    ],
    'sliding_window': [
        'sliding', 'window', 'substring', 'subarray', 'minimum', 'maximum',
        'average', 'sum', 'longest', 'shortest', 'contains', 'duplicate',
        'unique', 'atmost', 'exactly', 'k-different', 'fruit', 'substrings',
        'min-number', 'max-number', 'subarrays', 'replace'
    ],
    'two_pointers': [
        'two', 'pointers', 'two-pointer', 'sorted', 'squares', 'pairs',
        'triplets', 'quadruplets', 'sum', 'container', 'water', 'trap',
        'palindrome', 'linked-list', 'cycle', 'meeting', 'point', 'contact',
        'remove', 'element', 'duplicates', 'sorted', 'color', 'sort'
    ],
    'hash_table': [
        'hash', 'table', 'map', 'dictionary', 'counter', 'frequency',
        'unique', 'duplicate', 'group', 'anagram', 'isomorphic', 'word',
        'pattern', 'cloned', 'copy', 'cache', 'LRU', 'LFU', 'bulls', 'cows',
        'keyboard'
    ],
    'recursion': [
        'recursion', 'recursive', 'base', 'case', 'memo', 'memoization',
        'tower', 'hanoi', 'josephus', 'permutation', 'combination', 'climb'
    ],
    'trie': [
        'trie', 'prefix', 'tree', 'autocomplete', 'spell', 'checker',
        'suggestion', 'word', 'search', 'insert', 'delete', 'startsWith'
    ],
    'misc': [
        'calculator', 'switcher', 'senate', 'bulls', 'cows', 'bombs', 'brick',
        'wall', 'asteroid', 'collision', 'capacity', 'ship', 'packages',
        'fleet', 'pooling', 'allocation', 'booking', 'couples', 'cracking',
        'defuse', 'destination', 'detect', 'capital', 'dice', 'diet',
        'distribute', 'candies', 'people', 'cat', 'mouse', 'broken', 'can-i-win',
        'place-flowers', 'cinema', 'seat', 'campus', 'bikes', 'puzzle',
        'before-after', 'battleship', 'board', 'check', 'straight', 'line',
        'exist', 'day-of-week', 'year', 'money', 'bank', 'fill', 'emulate',
        'cpu', 'magic', 'index', 'triple', 'one-away', 'away', 'swaps'
    ]
}

def categorize(filename):
    name = filename.lower().replace('.py', '').replace('-', ' ').replace('_', ' ')
    
    # Special cases first
    if 'linked' in name and 'list' in name:
        return 'linked_lists'
    if 'binary search tree' in name or 'binary-tree' in name or 'bst' in name:
        return 'trees'
    if 'dynamic programming' in name or 'dp' in name:
        return 'dynamic_programming'
    
    for category, keywords in CATEGORIES.items():
        for keyword in keywords:
            if keyword in name:
                return category
    return None

def main():
    # Get all .py files not in categories yet
    files = [f for f in os.listdir('.') if f.endswith('.py') and not os.path.isdir(f)]
    
    categorized = 0
    uncategorized = []
    
    for filename in files:
        # Skip __init__.py and special files
        if filename in ['__init__.py', 'requirements.txt', 'setup.py']:
            continue
            
        category = categorize(filename)
        
        if category:
            dest = f'{category}/{filename}'
            try:
                shutil.move(filename, dest)
                categorized += 1
            except Exception as e:
                print(f"Error moving {filename}: {e}")
        else:
            uncategorized.append(filename)
    
    print(f"Categorized: {categorized}")
    print(f"Uncategorized: {len(uncategorized)}")
    if uncategorized:
        print("\nUncategorized files:")
        for f in sorted(uncategorized)[:50]:
            print(f"  {f}")

if __name__ == '__main__':
    main()
