package trees

import (
	"reflect"
	"testing"
)

func TestAVLTree(t *testing.T) {
	avl := NewAVLTree()
	
	// Test insertion
	values := []int{10, 20, 30, 40, 50, 25}
	for _, val := range values {
		avl.Insert(val)
	}
	
	// Test search
	if !avl.Search(25) {
		t.Error("Expected to find 25 in AVL tree")
	}
	
	if avl.Search(100) {
		t.Error("Expected not to find 100 in AVL tree")
	}
	
	// Test inorder traversal (should be sorted)
	inorder := avl.InorderTraversal()
	expected := []int{10, 20, 25, 30, 40, 50}
	if !reflect.DeepEqual(inorder, expected) {
		t.Errorf("Expected inorder %v, got %v", expected, inorder)
	}
}

func TestTrie(t *testing.T) {
	trie := NewTrie()
	
	// Test insertion and search
	words := []string{"cat", "car", "card", "care", "careful", "cars", "carry"}
	for _, word := range words {
		trie.Insert(word)
	}
	
	// Test exact word search
	if !trie.Search("car") {
		t.Error("Expected to find 'car' in trie")
	}
	
	if trie.Search("ca") {
		t.Error("Expected not to find 'ca' as complete word in trie")
	}
	
	// Test prefix search
	if !trie.StartsWith("car") {
		t.Error("Expected to find prefix 'car' in trie")
	}
	
	if trie.StartsWith("dog") {
		t.Error("Expected not to find prefix 'dog' in trie")
	}
	
	// Test getting words with prefix
	carWords := trie.GetWordsWithPrefix("car")
	expectedCarWords := []string{"car", "card", "care", "careful", "cars", "carry"}
	if len(carWords) != len(expectedCarWords) {
		t.Errorf("Expected %d words with prefix 'car', got %d", len(expectedCarWords), len(carWords))
	}
}

func TestSegmentTree(t *testing.T) {
	arr := []int{1, 3, 5, 7, 9, 11}
	st := NewSegmentTree(arr)
	
	// Test range sum query
	sum, err := st.Query(1, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := 3 + 5 + 7 // indices 1, 2, 3
	if sum != expected {
		t.Errorf("Expected sum %d, got %d", expected, sum)
	}
	
	// Test update
	err = st.Update(1, 10) // Change arr[1] from 3 to 10
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Test query after update
	sum, err = st.Query(1, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = 10 + 5 + 7 // indices 1, 2, 3 after update
	if sum != expected {
		t.Errorf("Expected sum %d after update, got %d", expected, sum)
	}
	
	// Test invalid range
	_, err = st.Query(-1, 3)
	if err == nil {
		t.Error("Expected error for invalid range")
	}
	
	// Test invalid index for update
	err = st.Update(-1, 5)
	if err == nil {
		t.Error("Expected error for invalid index")
	}
}

func TestFenwickTree(t *testing.T) {
	arr := []int{1, 3, 5, 7, 9, 11}
	ft := NewFenwickTree(arr)
	
	// Test prefix sum
	sum, err := ft.PrefixSum(3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := 1 + 3 + 5 + 7 // indices 0, 1, 2, 3
	if sum != expected {
		t.Errorf("Expected prefix sum %d, got %d", expected, sum)
	}
	
	// Test range sum
	sum, err = ft.RangeSum(1, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = 3 + 5 + 7 // indices 1, 2, 3
	if sum != expected {
		t.Errorf("Expected range sum %d, got %d", expected, sum)
	}
	
	// Test update
	err = ft.Update(1, 7) // Add 7 to arr[1], making it 3+7=10
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Test range sum after update
	sum, err = ft.RangeSum(1, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = 10 + 5 + 7 // indices 1, 2, 3 after update
	if sum != expected {
		t.Errorf("Expected range sum %d after update, got %d", expected, sum)
	}
	
	// Test set operation
	err = ft.Set(1, 20) // Set arr[1] to 20
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	sum, err = ft.RangeSum(1, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if sum != 20 {
		t.Errorf("Expected value 20 after set, got %d", sum)
	}
	
	// Test invalid operations
	_, err = ft.PrefixSum(-1)
	if err == nil {
		t.Error("Expected error for invalid index")
	}
	
	_, err = ft.RangeSum(5, 3)
	if err == nil {
		t.Error("Expected error for invalid range")
	}
}

func TestEmptyFenwickTree(t *testing.T) {
	ft := NewEmptyFenwickTree(5)
	
	// Test initial state
	sum, err := ft.PrefixSum(2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if sum != 0 {
		t.Errorf("Expected sum 0 for empty tree, got %d", sum)
	}
	
	// Test update
	err = ft.Update(2, 10)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	sum, err = ft.PrefixSum(2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if sum != 10 {
		t.Errorf("Expected sum 10 after update, got %d", sum)
	}
}

// Benchmark tests
func BenchmarkAVLTreeInsert(b *testing.B) {
	avl := NewAVLTree()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		avl.Insert(i)
	}
}

func BenchmarkAVLTreeSearch(b *testing.B) {
	avl := NewAVLTree()
	for i := 0; i < 1000; i++ {
		avl.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		avl.Search(i % 1000)
	}
}

func BenchmarkTrieInsert(b *testing.B) {
	trie := NewTrie()
	words := []string{"apple", "application", "apply", "appreciate", "approach"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Insert(words[i%len(words)])
	}
}

func BenchmarkTrieSearch(b *testing.B) {
	trie := NewTrie()
	words := []string{"apple", "application", "apply", "appreciate", "approach"}
	for _, word := range words {
		trie.Insert(word)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Search(words[i%len(words)])
	}
}

func BenchmarkSegmentTreeQuery(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = i + 1
	}
	st := NewSegmentTree(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.Query(0, 999)
	}
}

func BenchmarkFenwickTreeQuery(b *testing.B) {
	arr := make([]int, 1000)
	for i := range arr {
		arr[i] = i + 1
	}
	ft := NewFenwickTree(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ft.RangeSum(0, 999)
	}
}

func BenchmarkFenwickTreeUpdate(b *testing.B) {
	ft := NewEmptyFenwickTree(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ft.Update(i%1000, 1)
	}
}
