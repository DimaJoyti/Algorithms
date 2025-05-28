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
