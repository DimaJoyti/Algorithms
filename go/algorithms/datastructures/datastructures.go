// Package datastructures implements various data structures
package datastructures

import (
	"errors"
	"fmt"
)

// Stack implementation using slice
type Stack struct {
	items []interface{}
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{items: make([]interface{}, 0)}
}

// Push adds an item to the top of the stack
func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack
func (s *Stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("stack is empty")
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, nil
}

// Peek returns the top item without removing it
func (s *Stack) Peek() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("stack is empty")
	}

	return s.items[len(s.items)-1], nil
}

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s *Stack) Size() int {
	return len(s.items)
}

// Queue implementation using slice
type Queue struct {
	items []interface{}
}

// NewQueue creates a new queue
func NewQueue() *Queue {
	return &Queue{items: make([]interface{}, 0)}
}

// Enqueue adds an item to the rear of the queue
func (q *Queue) Enqueue(item interface{}) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue
func (q *Queue) Dequeue() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Front returns the front item without removing it
func (q *Queue) Front() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	return q.items[0], nil
}

// IsEmpty checks if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue
func (q *Queue) Size() int {
	return len(q.items)
}

// LinkedListNode represents a node in a linked list
type LinkedListNode struct {
	Data interface{}
	Next *LinkedListNode
}

// LinkedList implementation
type LinkedList struct {
	Head *LinkedListNode
	size int
}

// NewLinkedList creates a new linked list
func NewLinkedList() *LinkedList {
	return &LinkedList{Head: nil, size: 0}
}

// Insert adds a new node at the beginning
func (ll *LinkedList) Insert(data interface{}) {
	newNode := &LinkedListNode{Data: data, Next: ll.Head}
	ll.Head = newNode
	ll.size++
}

// Append adds a new node at the end
func (ll *LinkedList) Append(data interface{}) {
	newNode := &LinkedListNode{Data: data, Next: nil}

	if ll.Head == nil {
		ll.Head = newNode
	} else {
		current := ll.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
	ll.size++
}

// Delete removes the first occurrence of data
func (ll *LinkedList) Delete(data interface{}) bool {
	if ll.Head == nil {
		return false
	}

	if ll.Head.Data == data {
		ll.Head = ll.Head.Next
		ll.size--
		return true
	}

	current := ll.Head
	for current.Next != nil {
		if current.Next.Data == data {
			current.Next = current.Next.Next
			ll.size--
			return true
		}
		current = current.Next
	}

	return false
}

// Find searches for data in the list
func (ll *LinkedList) Find(data interface{}) *LinkedListNode {
	current := ll.Head
	for current != nil {
		if current.Data == data {
			return current
		}
		current = current.Next
	}
	return nil
}

// Size returns the number of nodes
func (ll *LinkedList) Size() int {
	return ll.size
}

// IsEmpty checks if the list is empty
func (ll *LinkedList) IsEmpty() bool {
	return ll.Head == nil
}

// ToSlice converts the linked list to a slice
func (ll *LinkedList) ToSlice() []interface{} {
	var result []interface{}
	current := ll.Head
	for current != nil {
		result = append(result, current.Data)
		current = current.Next
	}
	return result
}

// BinaryTreeNode represents a node in a binary tree
type BinaryTreeNode struct {
	Data  interface{}
	Left  *BinaryTreeNode
	Right *BinaryTreeNode
}

// BinaryTree implementation
type BinaryTree struct {
	Root *BinaryTreeNode
}

// NewBinaryTree creates a new binary tree
func NewBinaryTree() *BinaryTree {
	return &BinaryTree{Root: nil}
}

// Insert adds a new node (for BST)
func (bt *BinaryTree) Insert(data interface{}) {
	bt.Root = bt.insertNode(bt.Root, data)
}

func (bt *BinaryTree) insertNode(node *BinaryTreeNode, data interface{}) *BinaryTreeNode {
	if node == nil {
		return &BinaryTreeNode{Data: data, Left: nil, Right: nil}
	}

	// For simplicity, assume data is int for comparison
	if dataInt, ok := data.(int); ok {
		if nodeInt, ok := node.Data.(int); ok {
			if dataInt < nodeInt {
				node.Left = bt.insertNode(node.Left, data)
			} else {
				node.Right = bt.insertNode(node.Right, data)
			}
		}
	}

	return node
}

// Search finds a node with given data
func (bt *BinaryTree) Search(data interface{}) *BinaryTreeNode {
	return bt.searchNode(bt.Root, data)
}

func (bt *BinaryTree) searchNode(node *BinaryTreeNode, data interface{}) *BinaryTreeNode {
	if node == nil {
		return nil
	}

	if node.Data == data {
		return node
	}

	// For BST search
	if dataInt, ok := data.(int); ok {
		if nodeInt, ok := node.Data.(int); ok {
			if dataInt < nodeInt {
				return bt.searchNode(node.Left, data)
			} else {
				return bt.searchNode(node.Right, data)
			}
		}
	}

	// For general tree, search both subtrees
	left := bt.searchNode(node.Left, data)
	if left != nil {
		return left
	}
	return bt.searchNode(node.Right, data)
}

// InorderTraversal returns nodes in inorder
func (bt *BinaryTree) InorderTraversal() []interface{} {
	var result []interface{}
	bt.inorder(bt.Root, &result)
	return result
}

func (bt *BinaryTree) inorder(node *BinaryTreeNode, result *[]interface{}) {
	if node != nil {
		bt.inorder(node.Left, result)
		*result = append(*result, node.Data)
		bt.inorder(node.Right, result)
	}
}

// PreorderTraversal returns nodes in preorder
func (bt *BinaryTree) PreorderTraversal() []interface{} {
	var result []interface{}
	bt.preorder(bt.Root, &result)
	return result
}

func (bt *BinaryTree) preorder(node *BinaryTreeNode, result *[]interface{}) {
	if node != nil {
		*result = append(*result, node.Data)
		bt.preorder(node.Left, result)
		bt.preorder(node.Right, result)
	}
}

// PostorderTraversal returns nodes in postorder
func (bt *BinaryTree) PostorderTraversal() []interface{} {
	var result []interface{}
	bt.postorder(bt.Root, &result)
	return result
}

func (bt *BinaryTree) postorder(node *BinaryTreeNode, result *[]interface{}) {
	if node != nil {
		bt.postorder(node.Left, result)
		bt.postorder(node.Right, result)
		*result = append(*result, node.Data)
	}
}

// Height calculates the height of the tree
func (bt *BinaryTree) Height() int {
	return bt.height(bt.Root)
}

func (bt *BinaryTree) height(node *BinaryTreeNode) int {
	if node == nil {
		return 0
	}

	leftHeight := bt.height(node.Left)
	rightHeight := bt.height(node.Right)

	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// MinHeap implementation
type MinHeap struct {
	items []int
}

// NewMinHeap creates a new min heap
func NewMinHeap() *MinHeap {
	return &MinHeap{items: make([]int, 0)}
}

// Insert adds a new item to the heap
func (h *MinHeap) Insert(item int) {
	h.items = append(h.items, item)
	h.heapifyUp(len(h.items) - 1)
}

// ExtractMin removes and returns the minimum item
func (h *MinHeap) ExtractMin() (int, error) {
	if len(h.items) == 0 {
		return 0, errors.New("heap is empty")
	}

	min := h.items[0]
	lastIndex := len(h.items) - 1
	h.items[0] = h.items[lastIndex]
	h.items = h.items[:lastIndex]

	if len(h.items) > 0 {
		h.heapifyDown(0)
	}

	return min, nil
}

// Peek returns the minimum item without removing it
func (h *MinHeap) Peek() (int, error) {
	if len(h.items) == 0 {
		return 0, errors.New("heap is empty")
	}
	return h.items[0], nil
}

// Size returns the number of items in the heap
func (h *MinHeap) Size() int {
	return len(h.items)
}

// IsEmpty checks if the heap is empty
func (h *MinHeap) IsEmpty() bool {
	return len(h.items) == 0
}

func (h *MinHeap) heapifyUp(index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if h.items[index] >= h.items[parentIndex] {
			break
		}
		h.items[index], h.items[parentIndex] = h.items[parentIndex], h.items[index]
		index = parentIndex
	}
}

func (h *MinHeap) heapifyDown(index int) {
	for {
		leftChild := 2*index + 1
		rightChild := 2*index + 2
		smallest := index

		if leftChild < len(h.items) && h.items[leftChild] < h.items[smallest] {
			smallest = leftChild
		}

		if rightChild < len(h.items) && h.items[rightChild] < h.items[smallest] {
			smallest = rightChild
		}

		if smallest == index {
			break
		}

		h.items[index], h.items[smallest] = h.items[smallest], h.items[index]
		index = smallest
	}
}

// String returns a string representation of the heap
func (h *MinHeap) String() string {
	return fmt.Sprintf("MinHeap%v", h.items)
}

// SplayTreeNode represents a node in a splay tree
type SplayTreeNode struct {
	Key    int
	Left   *SplayTreeNode
	Right  *SplayTreeNode
	Parent *SplayTreeNode
}

// SplayTree implementation with self-balancing property
type SplayTree struct {
	Root *SplayTreeNode
}

// NewSplayTree creates a new splay tree
func NewSplayTree() *SplayTree {
	return &SplayTree{Root: nil}
}

// Insert adds a new key to the splay tree
// Time Complexity: O(log n) amortized, Space Complexity: O(1)
func (st *SplayTree) Insert(key int) {
	if st.Root == nil {
		st.Root = &SplayTreeNode{Key: key}
		return
	}

	st.splay(key)

	if st.Root.Key == key {
		return // Key already exists
	}

	newNode := &SplayTreeNode{Key: key}

	if key < st.Root.Key {
		newNode.Left = st.Root.Left
		newNode.Right = st.Root
		if st.Root.Left != nil {
			st.Root.Left.Parent = newNode
		}
		st.Root.Left = nil
		st.Root.Parent = newNode
	} else {
		newNode.Right = st.Root.Right
		newNode.Left = st.Root
		if st.Root.Right != nil {
			st.Root.Right.Parent = newNode
		}
		st.Root.Right = nil
		st.Root.Parent = newNode
	}

	st.Root = newNode
}

// Search finds a key in the splay tree
// Time Complexity: O(log n) amortized, Space Complexity: O(1)
func (st *SplayTree) Search(key int) bool {
	if st.Root == nil {
		return false
	}

	st.splay(key)
	return st.Root.Key == key
}

// Delete removes a key from the splay tree
// Time Complexity: O(log n) amortized, Space Complexity: O(1)
func (st *SplayTree) Delete(key int) bool {
	if st.Root == nil {
		return false
	}

	st.splay(key)

	if st.Root.Key != key {
		return false // Key not found
	}

	if st.Root.Left == nil {
		st.Root = st.Root.Right
		if st.Root != nil {
			st.Root.Parent = nil
		}
	} else if st.Root.Right == nil {
		st.Root = st.Root.Left
		if st.Root != nil {
			st.Root.Parent = nil
		}
	} else {
		// Node has both children
		leftSubtree := st.Root.Left
		rightSubtree := st.Root.Right

		st.Root = leftSubtree
		st.Root.Parent = nil

		// Find maximum in left subtree and splay it
		for st.Root.Right != nil {
			st.Root = st.Root.Right
		}
		st.splay(st.Root.Key)

		// Attach right subtree
		st.Root.Right = rightSubtree
		rightSubtree.Parent = st.Root
	}

	return true
}

// splay performs the splay operation (simplified version)
func (st *SplayTree) splay(key int) {
	if st.Root == nil {
		return
	}

	// Find the node or the closest node
	current := st.Root
	var parent *SplayTreeNode

	for current != nil && current.Key != key {
		parent = current
		if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}

	// If key not found, splay the last accessed node
	if current == nil {
		current = parent
	}

	// Perform bottom-up splaying
	for current != nil && current != st.Root {
		if current.Parent == st.Root {
			// Zig step
			if current == current.Parent.Left {
				st.rotateRight(current.Parent)
			} else {
				st.rotateLeft(current.Parent)
			}
		} else {
			// Zig-Zig or Zig-Zag step
			parent := current.Parent
			grandparent := parent.Parent

			if current == parent.Left && parent == grandparent.Left {
				// Zig-Zig (Left-Left)
				st.rotateRight(grandparent)
				st.rotateRight(parent)
			} else if current == parent.Right && parent == grandparent.Right {
				// Zig-Zig (Right-Right)
				st.rotateLeft(grandparent)
				st.rotateLeft(parent)
			} else if current == parent.Left && parent == grandparent.Right {
				// Zig-Zag (Left-Right)
				st.rotateRight(parent)
				st.rotateLeft(grandparent)
			} else {
				// Zig-Zag (Right-Left)
				st.rotateLeft(parent)
				st.rotateRight(grandparent)
			}
		}
	}
}

// rotateLeft performs left rotation
func (st *SplayTree) rotateLeft(node *SplayTreeNode) {
	rightChild := node.Right
	node.Right = rightChild.Left
	if rightChild.Left != nil {
		rightChild.Left.Parent = node
	}
	rightChild.Parent = node.Parent
	if node.Parent == nil {
		st.Root = rightChild
	} else if node == node.Parent.Left {
		node.Parent.Left = rightChild
	} else {
		node.Parent.Right = rightChild
	}
	rightChild.Left = node
	node.Parent = rightChild
}

// rotateRight performs right rotation
func (st *SplayTree) rotateRight(node *SplayTreeNode) {
	leftChild := node.Left
	node.Left = leftChild.Right
	if leftChild.Right != nil {
		leftChild.Right.Parent = node
	}
	leftChild.Parent = node.Parent
	if node.Parent == nil {
		st.Root = leftChild
	} else if node == node.Parent.Right {
		node.Parent.Right = leftChild
	} else {
		node.Parent.Left = leftChild
	}
	leftChild.Right = node
	node.Parent = leftChild
}

// SkipListNode represents a node in a skip list
type SkipListNode struct {
	Key     int
	Value   interface{}
	Forward []*SkipListNode
}

// SkipList implementation with probabilistic balancing
type SkipList struct {
	Header   *SkipListNode
	Level    int
	MaxLevel int
}

// NewSkipList creates a new skip list
func NewSkipList(maxLevel int) *SkipList {
	header := &SkipListNode{
		Key:     -1,
		Forward: make([]*SkipListNode, maxLevel+1),
	}
	return &SkipList{
		Header:   header,
		Level:    0,
		MaxLevel: maxLevel,
	}
}

// randomLevel generates a random level for new nodes
func (sl *SkipList) randomLevel() int {
	level := 0
	for level < sl.MaxLevel && (level == 0 || (level < sl.MaxLevel && (level&1) == 0)) {
		level++
	}
	return level
}

// Insert adds a key-value pair to the skip list
// Time Complexity: O(log n) expected, Space Complexity: O(log n)
func (sl *SkipList) Insert(key int, value interface{}) {
	update := make([]*SkipListNode, sl.MaxLevel+1)
	current := sl.Header

	// Find position to insert
	for i := sl.Level; i >= 0; i-- {
		for current.Forward[i] != nil && current.Forward[i].Key < key {
			current = current.Forward[i]
		}
		update[i] = current
	}

	current = current.Forward[0]

	// If key already exists, update value
	if current != nil && current.Key == key {
		current.Value = value
		return
	}

	// Generate random level for new node
	newLevel := sl.randomLevel()
	if newLevel > sl.Level {
		for i := sl.Level + 1; i <= newLevel; i++ {
			update[i] = sl.Header
		}
		sl.Level = newLevel
	}

	// Create new node
	newNode := &SkipListNode{
		Key:     key,
		Value:   value,
		Forward: make([]*SkipListNode, newLevel+1),
	}

	// Update forward pointers
	for i := 0; i <= newLevel; i++ {
		newNode.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = newNode
	}
}

// Search finds a value by key in the skip list
// Time Complexity: O(log n) expected, Space Complexity: O(1)
func (sl *SkipList) Search(key int) (interface{}, bool) {
	current := sl.Header

	for i := sl.Level; i >= 0; i-- {
		for current.Forward[i] != nil && current.Forward[i].Key < key {
			current = current.Forward[i]
		}
	}

	current = current.Forward[0]

	if current != nil && current.Key == key {
		return current.Value, true
	}

	return nil, false
}

// Delete removes a key from the skip list
// Time Complexity: O(log n) expected, Space Complexity: O(log n)
func (sl *SkipList) Delete(key int) bool {
	update := make([]*SkipListNode, sl.MaxLevel+1)
	current := sl.Header

	// Find position to delete
	for i := sl.Level; i >= 0; i-- {
		for current.Forward[i] != nil && current.Forward[i].Key < key {
			current = current.Forward[i]
		}
		update[i] = current
	}

	current = current.Forward[0]

	if current == nil || current.Key != key {
		return false // Key not found
	}

	// Update forward pointers
	for i := 0; i <= sl.Level; i++ {
		if update[i].Forward[i] != current {
			break
		}
		update[i].Forward[i] = current.Forward[i]
	}

	// Update level
	for sl.Level > 0 && sl.Header.Forward[sl.Level] == nil {
		sl.Level--
	}

	return true
}

// DisjointSet (Union-Find) with path compression and union by rank
type DisjointSet struct {
	parent []int
	rank   []int
	size   int
}

// NewDisjointSet creates a new disjoint set with n elements
func NewDisjointSet(n int) *DisjointSet {
	parent := make([]int, n)
	rank := make([]int, n)

	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 0
	}

	return &DisjointSet{
		parent: parent,
		rank:   rank,
		size:   n,
	}
}

// Find finds the representative of the set containing x with path compression
// Time Complexity: O(α(n)) amortized, Space Complexity: O(1)
func (ds *DisjointSet) Find(x int) int {
	if x < 0 || x >= ds.size {
		return -1
	}

	if ds.parent[x] != x {
		ds.parent[x] = ds.Find(ds.parent[x]) // Path compression
	}

	return ds.parent[x]
}

// Union unites two sets containing x and y using union by rank
// Time Complexity: O(α(n)) amortized, Space Complexity: O(1)
func (ds *DisjointSet) Union(x, y int) bool {
	rootX := ds.Find(x)
	rootY := ds.Find(y)

	if rootX == -1 || rootY == -1 || rootX == rootY {
		return false
	}

	// Union by rank
	if ds.rank[rootX] < ds.rank[rootY] {
		ds.parent[rootX] = rootY
	} else if ds.rank[rootX] > ds.rank[rootY] {
		ds.parent[rootY] = rootX
	} else {
		ds.parent[rootY] = rootX
		ds.rank[rootX]++
	}

	return true
}

// Connected checks if two elements are in the same set
// Time Complexity: O(α(n)) amortized, Space Complexity: O(1)
func (ds *DisjointSet) Connected(x, y int) bool {
	return ds.Find(x) == ds.Find(y) && ds.Find(x) != -1
}

// CountSets returns the number of disjoint sets
// Time Complexity: O(n), Space Complexity: O(1)
func (ds *DisjointSet) CountSets() int {
	roots := make(map[int]bool)
	for i := 0; i < ds.size; i++ {
		root := ds.Find(i)
		if root != -1 {
			roots[root] = true
		}
	}
	return len(roots)
}
