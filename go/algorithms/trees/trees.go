// Package trees implements various advanced tree data structures
package trees

import (
	"errors"
)

// AVLNode represents a node in an AVL tree
type AVLNode struct {
	Value  int
	Height int
	Left   *AVLNode
	Right  *AVLNode
}

// AVLTree represents a self-balancing binary search tree
type AVLTree struct {
	Root *AVLNode
}

// NewAVLTree creates a new AVL tree
func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

// getHeight returns the height of a node
func (avl *AVLTree) getHeight(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

// getBalance returns the balance factor of a node
func (avl *AVLTree) getBalance(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return avl.getHeight(node.Left) - avl.getHeight(node.Right)
}

// updateHeight updates the height of a node
func (avl *AVLTree) updateHeight(node *AVLNode) {
	if node != nil {
		leftHeight := avl.getHeight(node.Left)
		rightHeight := avl.getHeight(node.Right)
		if leftHeight > rightHeight {
			node.Height = leftHeight + 1
		} else {
			node.Height = rightHeight + 1
		}
	}
}

// rightRotate performs a right rotation
func (avl *AVLTree) rightRotate(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	// Perform rotation
	x.Right = y
	y.Left = T2

	// Update heights
	avl.updateHeight(y)
	avl.updateHeight(x)

	return x
}

// leftRotate performs a left rotation
func (avl *AVLTree) leftRotate(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	// Perform rotation
	y.Left = x
	x.Right = T2

	// Update heights
	avl.updateHeight(x)
	avl.updateHeight(y)

	return y
}

// Insert inserts a value into the AVL tree
// Time Complexity: O(log n), Space Complexity: O(log n)
func (avl *AVLTree) Insert(value int) {
	avl.Root = avl.insertHelper(avl.Root, value)
}

func (avl *AVLTree) insertHelper(node *AVLNode, value int) *AVLNode {
	// Standard BST insertion
	if node == nil {
		return &AVLNode{Value: value, Height: 1}
	}

	if value < node.Value {
		node.Left = avl.insertHelper(node.Left, value)
	} else if value > node.Value {
		node.Right = avl.insertHelper(node.Right, value)
	} else {
		// Equal values not allowed
		return node
	}

	// Update height
	avl.updateHeight(node)

	// Get balance factor
	balance := avl.getBalance(node)

	// Left Left Case
	if balance > 1 && value < node.Left.Value {
		return avl.rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && value > node.Right.Value {
		return avl.leftRotate(node)
	}

	// Left Right Case
	if balance > 1 && value > node.Left.Value {
		node.Left = avl.leftRotate(node.Left)
		return avl.rightRotate(node)
	}

	// Right Left Case
	if balance < -1 && value < node.Right.Value {
		node.Right = avl.rightRotate(node.Right)
		return avl.leftRotate(node)
	}

	return node
}

// Search searches for a value in the AVL tree
// Time Complexity: O(log n), Space Complexity: O(log n)
func (avl *AVLTree) Search(value int) bool {
	return avl.searchHelper(avl.Root, value)
}

func (avl *AVLTree) searchHelper(node *AVLNode, value int) bool {
	if node == nil {
		return false
	}
	if value == node.Value {
		return true
	}
	if value < node.Value {
		return avl.searchHelper(node.Left, value)
	}
	return avl.searchHelper(node.Right, value)
}

// InorderTraversal returns inorder traversal of the tree
func (avl *AVLTree) InorderTraversal() []int {
	result := make([]int, 0)
	avl.inorderHelper(avl.Root, &result)
	return result
}

func (avl *AVLTree) inorderHelper(node *AVLNode, result *[]int) {
	if node != nil {
		avl.inorderHelper(node.Left, result)
		*result = append(*result, node.Value)
		avl.inorderHelper(node.Right, result)
	}
}

// TrieNode represents a node in a Trie
type TrieNode struct {
	Children map[rune]*TrieNode
	IsEnd    bool
}

// Trie represents a prefix tree
type Trie struct {
	Root *TrieNode
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
	return &Trie{
		Root: &TrieNode{
			Children: make(map[rune]*TrieNode),
			IsEnd:    false,
		},
	}
}

// Insert inserts a word into the trie
// Time Complexity: O(m), Space Complexity: O(m) where m is word length
func (t *Trie) Insert(word string) {
	current := t.Root
	for _, char := range word {
		if _, exists := current.Children[char]; !exists {
			current.Children[char] = &TrieNode{
				Children: make(map[rune]*TrieNode),
				IsEnd:    false,
			}
		}
		current = current.Children[char]
	}
	current.IsEnd = true
}

// Search searches for a word in the trie
// Time Complexity: O(m), Space Complexity: O(1)
func (t *Trie) Search(word string) bool {
	current := t.Root
	for _, char := range word {
		if _, exists := current.Children[char]; !exists {
			return false
		}
		current = current.Children[char]
	}
	return current.IsEnd
}

// StartsWith checks if any word in the trie starts with the given prefix
// Time Complexity: O(m), Space Complexity: O(1)
func (t *Trie) StartsWith(prefix string) bool {
	current := t.Root
	for _, char := range prefix {
		if _, exists := current.Children[char]; !exists {
			return false
		}
		current = current.Children[char]
	}
	return true
}

// GetWordsWithPrefix returns all words with the given prefix
// Time Complexity: O(p + n), Space Complexity: O(n) where p is prefix length, n is number of words
func (t *Trie) GetWordsWithPrefix(prefix string) []string {
	current := t.Root

	// Navigate to the prefix
	for _, char := range prefix {
		if _, exists := current.Children[char]; !exists {
			return []string{}
		}
		current = current.Children[char]
	}

	// Collect all words from this point
	words := make([]string, 0)
	t.collectWords(current, prefix, &words)
	return words
}

func (t *Trie) collectWords(node *TrieNode, currentWord string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, currentWord)
	}

	for char, child := range node.Children {
		t.collectWords(child, currentWord+string(char), words)
	}
}

// SegmentTree represents a segment tree for range queries
type SegmentTree struct {
	tree []int
	n    int
}

// NewSegmentTree creates a new segment tree from an array
// Time Complexity: O(n), Space Complexity: O(n)
func NewSegmentTree(arr []int) *SegmentTree {
	n := len(arr)
	tree := make([]int, 4*n) // Allocate enough space
	st := &SegmentTree{tree: tree, n: n}
	if n > 0 {
		st.build(arr, 0, 0, n-1)
	}
	return st
}

func (st *SegmentTree) build(arr []int, node, start, end int) {
	if start == end {
		st.tree[node] = arr[start]
	} else {
		mid := (start + end) / 2
		st.build(arr, 2*node+1, start, mid)
		st.build(arr, 2*node+2, mid+1, end)
		st.tree[node] = st.tree[2*node+1] + st.tree[2*node+2]
	}
}

// Query returns the sum of elements in range [l, r]
// Time Complexity: O(log n), Space Complexity: O(log n)
func (st *SegmentTree) Query(l, r int) (int, error) {
	if l < 0 || r >= st.n || l > r {
		return 0, errors.New("invalid range")
	}
	return st.queryHelper(0, 0, st.n-1, l, r), nil
}

func (st *SegmentTree) queryHelper(node, start, end, l, r int) int {
	if r < start || end < l {
		return 0 // Outside range
	}
	if l <= start && end <= r {
		return st.tree[node] // Completely inside range
	}

	mid := (start + end) / 2
	leftSum := st.queryHelper(2*node+1, start, mid, l, r)
	rightSum := st.queryHelper(2*node+2, mid+1, end, l, r)
	return leftSum + rightSum
}

// Update updates the value at index idx
// Time Complexity: O(log n), Space Complexity: O(log n)
func (st *SegmentTree) Update(idx, val int) error {
	if idx < 0 || idx >= st.n {
		return errors.New("index out of bounds")
	}
	st.updateHelper(0, 0, st.n-1, idx, val)
	return nil
}

func (st *SegmentTree) updateHelper(node, start, end, idx, val int) {
	if start == end {
		st.tree[node] = val
	} else {
		mid := (start + end) / 2
		if idx <= mid {
			st.updateHelper(2*node+1, start, mid, idx, val)
		} else {
			st.updateHelper(2*node+2, mid+1, end, idx, val)
		}
		st.tree[node] = st.tree[2*node+1] + st.tree[2*node+2]
	}
}

// FenwickTree (Binary Indexed Tree) for efficient prefix sum queries
type FenwickTree struct {
	tree []int
	n    int
}

// NewFenwickTree creates a new Fenwick tree
// Time Complexity: O(n log n), Space Complexity: O(n)
func NewFenwickTree(arr []int) *FenwickTree {
	n := len(arr)
	tree := make([]int, n+1) // 1-indexed
	ft := &FenwickTree{tree: tree, n: n}

	for i := 0; i < n; i++ {
		ft.Update(i, arr[i])
	}

	return ft
}

// NewEmptyFenwickTree creates an empty Fenwick tree of given size
func NewEmptyFenwickTree(size int) *FenwickTree {
	return &FenwickTree{
		tree: make([]int, size+1),
		n:    size,
	}
}

// Update adds delta to the element at index idx
// Time Complexity: O(log n), Space Complexity: O(1)
func (ft *FenwickTree) Update(idx, delta int) error {
	if idx < 0 || idx >= ft.n {
		return errors.New("index out of bounds")
	}

	idx++ // Convert to 1-indexed
	for idx <= ft.n {
		ft.tree[idx] += delta
		idx += idx & (-idx) // Add last set bit
	}
	return nil
}

// PrefixSum returns the sum of elements from index 0 to idx (inclusive)
// Time Complexity: O(log n), Space Complexity: O(1)
func (ft *FenwickTree) PrefixSum(idx int) (int, error) {
	if idx < 0 || idx >= ft.n {
		return 0, errors.New("index out of bounds")
	}

	idx++ // Convert to 1-indexed
	sum := 0
	for idx > 0 {
		sum += ft.tree[idx]
		idx -= idx & (-idx) // Remove last set bit
	}
	return sum, nil
}

// RangeSum returns the sum of elements in range [left, right] (inclusive)
// Time Complexity: O(log n), Space Complexity: O(1)
func (ft *FenwickTree) RangeSum(left, right int) (int, error) {
	if left < 0 || right >= ft.n || left > right {
		return 0, errors.New("invalid range")
	}

	rightSum, err := ft.PrefixSum(right)
	if err != nil {
		return 0, err
	}

	if left == 0 {
		return rightSum, nil
	}

	leftSum, err := ft.PrefixSum(left - 1)
	if err != nil {
		return 0, err
	}

	return rightSum - leftSum, nil
}

// Set sets the value at index idx to val
// Time Complexity: O(log n), Space Complexity: O(1)
func (ft *FenwickTree) Set(idx, val int) error {
	if idx < 0 || idx >= ft.n {
		return errors.New("index out of bounds")
	}

	// Get current value
	currentVal, err := ft.RangeSum(idx, idx)
	if err != nil {
		return err
	}

	// Update with the difference
	return ft.Update(idx, val-currentVal)
}
