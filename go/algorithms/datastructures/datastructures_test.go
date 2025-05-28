package datastructures

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()

	// Test empty stack
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}

	if stack.Size() != 0 {
		t.Errorf("Empty stack size should be 0, got %d", stack.Size())
	}

	// Test push
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.IsEmpty() {
		t.Error("Stack should not be empty after push")
	}

	if stack.Size() != 3 {
		t.Errorf("Stack size should be 3, got %d", stack.Size())
	}

	// Test peek
	top, err := stack.Peek()
	if err != nil {
		t.Errorf("Peek should not return error: %v", err)
	}
	if top != 3 {
		t.Errorf("Top element should be 3, got %v", top)
	}

	// Test pop
	item, err := stack.Pop()
	if err != nil {
		t.Errorf("Pop should not return error: %v", err)
	}
	if item != 3 {
		t.Errorf("Popped item should be 3, got %v", item)
	}

	if stack.Size() != 2 {
		t.Errorf("Stack size should be 2 after pop, got %d", stack.Size())
	}

	// Test pop empty stack
	stack.Pop() // Remove 2
	stack.Pop() // Remove 1

	_, err = stack.Pop()
	if err == nil {
		t.Error("Pop on empty stack should return error")
	}

	_, err = stack.Peek()
	if err == nil {
		t.Error("Peek on empty stack should return error")
	}
}

func TestQueue(t *testing.T) {
	queue := NewQueue()

	// Test empty queue
	if !queue.IsEmpty() {
		t.Error("New queue should be empty")
	}

	if queue.Size() != 0 {
		t.Errorf("Empty queue size should be 0, got %d", queue.Size())
	}

	// Test enqueue
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	if queue.IsEmpty() {
		t.Error("Queue should not be empty after enqueue")
	}

	if queue.Size() != 3 {
		t.Errorf("Queue size should be 3, got %d", queue.Size())
	}

	// Test front
	front, err := queue.Front()
	if err != nil {
		t.Errorf("Front should not return error: %v", err)
	}
	if front != 1 {
		t.Errorf("Front element should be 1, got %v", front)
	}

	// Test dequeue
	item, err := queue.Dequeue()
	if err != nil {
		t.Errorf("Dequeue should not return error: %v", err)
	}
	if item != 1 {
		t.Errorf("Dequeued item should be 1, got %v", item)
	}

	if queue.Size() != 2 {
		t.Errorf("Queue size should be 2 after dequeue, got %d", queue.Size())
	}

	// Test dequeue empty queue
	queue.Dequeue() // Remove 2
	queue.Dequeue() // Remove 3

	_, err = queue.Dequeue()
	if err == nil {
		t.Error("Dequeue on empty queue should return error")
	}

	_, err = queue.Front()
	if err == nil {
		t.Error("Front on empty queue should return error")
	}
}

func TestLinkedList(t *testing.T) {
	ll := NewLinkedList()

	// Test empty list
	if !ll.IsEmpty() {
		t.Error("New linked list should be empty")
	}

	if ll.Size() != 0 {
		t.Errorf("Empty list size should be 0, got %d", ll.Size())
	}

	// Test insert
	ll.Insert(1)
	ll.Insert(2)
	ll.Insert(3)

	if ll.IsEmpty() {
		t.Error("List should not be empty after insert")
	}

	if ll.Size() != 3 {
		t.Errorf("List size should be 3, got %d", ll.Size())
	}

	// Test find
	node := ll.Find(2)
	if node == nil {
		t.Error("Should find node with data 2")
	}
	if node.Data != 2 {
		t.Errorf("Found node should have data 2, got %v", node.Data)
	}

	// Test append
	ll.Append(4)
	if ll.Size() != 4 {
		t.Errorf("List size should be 4 after append, got %d", ll.Size())
	}

	// Test to slice
	slice := ll.ToSlice()
	expected := []interface{}{3, 2, 1, 4} // Insert adds to beginning
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("ToSlice() = %v, want %v", slice, expected)
	}

	// Test delete
	deleted := ll.Delete(2)
	if !deleted {
		t.Error("Should successfully delete existing element")
	}

	if ll.Size() != 3 {
		t.Errorf("List size should be 3 after delete, got %d", ll.Size())
	}

	// Test delete non-existent
	deleted = ll.Delete(10)
	if deleted {
		t.Error("Should not delete non-existent element")
	}
}

func TestBinaryTree(t *testing.T) {
	bt := NewBinaryTree()

	// Test insert
	bt.Insert(5)
	bt.Insert(3)
	bt.Insert(7)
	bt.Insert(1)
	bt.Insert(9)

	// Test search
	node := bt.Search(3)
	if node == nil {
		t.Error("Should find node with data 3")
	}
	if node.Data != 3 {
		t.Errorf("Found node should have data 3, got %v", node.Data)
	}

	// Test search non-existent
	node = bt.Search(10)
	if node != nil {
		t.Error("Should not find non-existent node")
	}

	// Test traversals
	inorder := bt.InorderTraversal()
	expectedInorder := []interface{}{1, 3, 5, 7, 9}
	if !reflect.DeepEqual(inorder, expectedInorder) {
		t.Errorf("Inorder traversal = %v, want %v", inorder, expectedInorder)
	}

	preorder := bt.PreorderTraversal()
	expectedPreorder := []interface{}{5, 3, 1, 7, 9}
	if !reflect.DeepEqual(preorder, expectedPreorder) {
		t.Errorf("Preorder traversal = %v, want %v", preorder, expectedPreorder)
	}

	postorder := bt.PostorderTraversal()
	expectedPostorder := []interface{}{1, 3, 9, 7, 5}
	if !reflect.DeepEqual(postorder, expectedPostorder) {
		t.Errorf("Postorder traversal = %v, want %v", postorder, expectedPostorder)
	}

	// Test height
	height := bt.Height()
	if height != 3 {
		t.Errorf("Tree height should be 3, got %d", height)
	}
}

func TestMinHeap(t *testing.T) {
	heap := NewMinHeap()

	// Test empty heap
	if !heap.IsEmpty() {
		t.Error("New heap should be empty")
	}

	if heap.Size() != 0 {
		t.Errorf("Empty heap size should be 0, got %d", heap.Size())
	}

	// Test insert
	heap.Insert(5)
	heap.Insert(3)
	heap.Insert(8)
	heap.Insert(1)
	heap.Insert(9)

	if heap.IsEmpty() {
		t.Error("Heap should not be empty after insert")
	}

	if heap.Size() != 5 {
		t.Errorf("Heap size should be 5, got %d", heap.Size())
	}

	// Test peek
	min, err := heap.Peek()
	if err != nil {
		t.Errorf("Peek should not return error: %v", err)
	}
	if min != 1 {
		t.Errorf("Min element should be 1, got %d", min)
	}

	// Test extract min
	extracted, err := heap.ExtractMin()
	if err != nil {
		t.Errorf("ExtractMin should not return error: %v", err)
	}
	if extracted != 1 {
		t.Errorf("Extracted min should be 1, got %d", extracted)
	}

	if heap.Size() != 4 {
		t.Errorf("Heap size should be 4 after extract, got %d", heap.Size())
	}

	// Test that next min is correct
	min, err = heap.Peek()
	if err != nil {
		t.Errorf("Peek should not return error: %v", err)
	}
	if min != 3 {
		t.Errorf("Next min element should be 3, got %d", min)
	}

	// Test extract all elements
	extracted, _ = heap.ExtractMin() // 3
	extracted, _ = heap.ExtractMin() // 5
	extracted, _ = heap.ExtractMin() // 8
	extracted, _ = heap.ExtractMin() // 9

	if extracted != 9 {
		t.Errorf("Last extracted should be 9, got %d", extracted)
	}

	// Test extract from empty heap
	_, err = heap.ExtractMin()
	if err == nil {
		t.Error("ExtractMin on empty heap should return error")
	}

	_, err = heap.Peek()
	if err == nil {
		t.Error("Peek on empty heap should return error")
	}
}

func TestMinHeapProperty(t *testing.T) {
	heap := NewMinHeap()
	
	// Insert random values
	values := []int{15, 10, 20, 8, 25, 5, 12}
	for _, v := range values {
		heap.Insert(v)
	}

	// Extract all values and verify they come out in sorted order
	var extracted []int
	for !heap.IsEmpty() {
		min, _ := heap.ExtractMin()
		extracted = append(extracted, min)
	}

	// Verify sorted order
	for i := 1; i < len(extracted); i++ {
		if extracted[i] < extracted[i-1] {
			t.Errorf("Heap property violated: %v is not sorted", extracted)
		}
	}
}

// Benchmark tests
func BenchmarkStackPush(b *testing.B) {
	stack := NewStack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkQueueEnqueue(b *testing.B) {
	queue := NewQueue()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
	}
}

func BenchmarkLinkedListInsert(b *testing.B) {
	ll := NewLinkedList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ll.Insert(i)
	}
}

func BenchmarkMinHeapInsert(b *testing.B) {
	heap := NewMinHeap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i)
	}
}
