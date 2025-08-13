// Package strings implements various string algorithms
package strings

import (
	"sort"
	"strings"
	"unicode"
)

// ReverseString reverses a string
// Time Complexity: O(n), Space Complexity: O(n)
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome
// Time Complexity: O(n), Space Complexity: O(1)
func IsPalindrome(s string) bool {
	// Convert to lowercase and remove non-alphanumeric characters
	cleaned := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned += strings.ToLower(string(r))
		}
	}

	left, right := 0, len(cleaned)-1
	for left < right {
		if cleaned[left] != cleaned[right] {
			return false
		}
		left++
		right--
	}
	return true
}

// AreAnagrams checks if two strings are anagrams
// Time Complexity: O(n log n), Space Complexity: O(n)
func AreAnagrams(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	// Convert to lowercase and sort
	r1 := []rune(strings.ToLower(s1))
	r2 := []rune(strings.ToLower(s2))

	sort.Slice(r1, func(i, j int) bool { return r1[i] < r1[j] })
	sort.Slice(r2, func(i, j int) bool { return r2[i] < r2[j] })

	return string(r1) == string(r2)
}

// AreAnagramsOptimized checks if two strings are anagrams using character counting
// Time Complexity: O(n), Space Complexity: O(1) for ASCII
func AreAnagramsOptimized(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	charCount := make(map[rune]int)

	// Count characters in first string
	for _, r := range strings.ToLower(s1) {
		charCount[r]++
	}

	// Subtract characters from second string
	for _, r := range strings.ToLower(s2) {
		charCount[r]--
		if charCount[r] < 0 {
			return false
		}
	}

	// Check if all counts are zero
	for _, count := range charCount {
		if count != 0 {
			return false
		}
	}

	return true
}

// CountVowels counts the number of vowels in a string
// Time Complexity: O(n), Space Complexity: O(1)
func CountVowels(s string) int {
	vowels := "aeiouAEIOU"
	count := 0

	for _, r := range s {
		if strings.ContainsRune(vowels, r) {
			count++
		}
	}

	return count
}

// MaxCharacter finds the most frequently occurring character
// In case of ties, returns the first occurring character
// Time Complexity: O(n), Space Complexity: O(k) where k is unique characters
func MaxCharacter(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}

	charCount := make(map[rune]int)
	firstOccurrence := make(map[rune]int)

	// Count characters and track first occurrence
	for i, r := range s {
		charCount[r]++
		if _, exists := firstOccurrence[r]; !exists {
			firstOccurrence[r] = i
		}
	}

	// Find maximum, preferring first occurring character in case of ties
	var maxChar rune
	maxCount := 0
	earliestIndex := len(s)

	for char, count := range charCount {
		if count > maxCount || (count == maxCount && firstOccurrence[char] < earliestIndex) {
			maxCount = count
			maxChar = char
			earliestIndex = firstOccurrence[char]
		}
	}

	return maxChar, maxCount
}

// Capitalize capitalizes the first letter of each word
// Time Complexity: O(n), Space Complexity: O(n)
func Capitalize(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// LongestCommonSubstring finds the longest common substring between two strings
// Time Complexity: O(m*n), Space Complexity: O(m*n)
func LongestCommonSubstring(s1, s2 string) string {
	m, n := len(s1), len(s2)
	if m == 0 || n == 0 {
		return ""
	}

	// Create DP table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	maxLength := 0
	endPos := 0

	// Fill DP table
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLength {
					maxLength = dp[i][j]
					endPos = i
				}
			}
		}
	}

	if maxLength == 0 {
		return ""
	}

	return s1[endPos-maxLength : endPos]
}

// LongestCommonSubsequence finds the length of longest common subsequence
// Time Complexity: O(m*n), Space Complexity: O(m*n)
func LongestCommonSubsequence(s1, s2 string) int {
	m, n := len(s1), len(s2)

	// Create DP table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Fill DP table
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

// EditDistance calculates the minimum edit distance (Levenshtein distance)
// Time Complexity: O(m*n), Space Complexity: O(m*n)
func EditDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)

	// Create DP table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize base cases
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Fill DP table
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}

	return dp[m][n]
}

// KMPSearch implements Knuth-Morris-Pratt string matching algorithm
// Time Complexity: O(n + m), Space Complexity: O(m)
func KMPSearch(text, pattern string) []int {
	if len(pattern) == 0 {
		return []int{}
	}

	// Build failure function
	lps := buildLPS(pattern)

	var matches []int
	i, j := 0, 0

	for i < len(text) {
		if text[i] == pattern[j] {
			i++
			j++
		}

		if j == len(pattern) {
			matches = append(matches, i-j)
			j = lps[j-1]
		} else if i < len(text) && text[i] != pattern[j] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}

	return matches
}

// buildLPS builds the Longest Proper Prefix which is also Suffix array
func buildLPS(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)
	length := 0
	i := 1

	for i < m {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	return lps
}

// RabinKarpSearch implements Rabin-Karp string matching algorithm
// Time Complexity: O(n + m) average, O(nm) worst, Space Complexity: O(1)
func RabinKarpSearch(text, pattern string) []int {
	const prime = 101

	n, m := len(text), len(pattern)
	if m > n {
		return []int{}
	}

	var matches []int

	// Calculate hash values
	patternHash := 0
	textHash := 0
	h := 1

	// Calculate h = pow(256, m-1) % prime
	for i := 0; i < m-1; i++ {
		h = (h * 256) % prime
	}

	// Calculate initial hash values
	for i := 0; i < m; i++ {
		patternHash = (256*patternHash + int(pattern[i])) % prime
		textHash = (256*textHash + int(text[i])) % prime
	}

	// Slide the pattern over text
	for i := 0; i <= n-m; i++ {
		// Check if hash values match
		if patternHash == textHash {
			// Check characters one by one
			match := true
			for j := 0; j < m; j++ {
				if text[i+j] != pattern[j] {
					match = false
					break
				}
			}
			if match {
				matches = append(matches, i)
			}
		}

		// Calculate hash for next window
		if i < n-m {
			textHash = (256*(textHash-int(text[i])*h) + int(text[i+m])) % prime
			if textHash < 0 {
				textHash += prime
			}
		}
	}

	return matches
}

// IsSubsequence checks if s is a subsequence of t
// Time Complexity: O(n), Space Complexity: O(1)
func IsSubsequence(s, t string) bool {
	i := 0
	for j := 0; j < len(t) && i < len(s); j++ {
		if s[i] == t[j] {
			i++
		}
	}
	return i == len(s)
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// SuffixArray constructs suffix array for a given string
// Time Complexity: O(n log n), Space Complexity: O(n)
func SuffixArray(s string) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}

	// Create array of suffixes with their indices
	type suffix struct {
		index int
		str   string
	}

	suffixes := make([]suffix, n)
	for i := 0; i < n; i++ {
		suffixes[i] = suffix{index: i, str: s[i:]}
	}

	// Sort suffixes lexicographically
	sort.Slice(suffixes, func(i, j int) bool {
		return suffixes[i].str < suffixes[j].str
	})

	// Extract indices
	result := make([]int, n)
	for i, suf := range suffixes {
		result[i] = suf.index
	}

	return result
}

// LongestRepeatedSubstring finds the longest substring that appears at least twice
// Time Complexity: O(n²), Space Complexity: O(n)
func LongestRepeatedSubstring(s string) string {
	n := len(s)
	if n <= 1 {
		return ""
	}

	maxLength := 0
	result := ""

	// Check all possible substrings
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			substring := s[i:j]
			if len(substring) > maxLength {
				// Check if this substring appears again
				if strings.Contains(s[j:], substring) {
					maxLength = len(substring)
					result = substring
				}
			}
		}
	}

	return result
}

// AhoCorasickNode represents a node in the Aho-Corasick automaton
type AhoCorasickNode struct {
	children map[rune]*AhoCorasickNode
	failure  *AhoCorasickNode
	output   []string
}

// AhoCorasick implements the Aho-Corasick algorithm for multiple pattern matching
type AhoCorasick struct {
	root *AhoCorasickNode
}

// NewAhoCorasick creates a new Aho-Corasick automaton
func NewAhoCorasick() *AhoCorasick {
	return &AhoCorasick{
		root: &AhoCorasickNode{
			children: make(map[rune]*AhoCorasickNode),
			output:   make([]string, 0),
		},
	}
}

// AddPattern adds a pattern to the automaton
func (ac *AhoCorasick) AddPattern(pattern string) {
	current := ac.root
	for _, char := range pattern {
		if _, exists := current.children[char]; !exists {
			current.children[char] = &AhoCorasickNode{
				children: make(map[rune]*AhoCorasickNode),
				output:   make([]string, 0),
			}
		}
		current = current.children[char]
	}
	current.output = append(current.output, pattern)
}

// BuildFailureFunction builds the failure function for the automaton
func (ac *AhoCorasick) BuildFailureFunction() {
	queue := make([]*AhoCorasickNode, 0)

	// Initialize failure function for depth 1 nodes
	for _, child := range ac.root.children {
		child.failure = ac.root
		queue = append(queue, child)
	}

	// Build failure function using BFS
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for char, child := range current.children {
			queue = append(queue, child)

			// Find failure link
			failure := current.failure
			for failure != nil && failure.children[char] == nil {
				failure = failure.failure
			}

			if failure != nil {
				child.failure = failure.children[char]
			} else {
				child.failure = ac.root
			}

			// Add output from failure node
			if child.failure != nil {
				child.output = append(child.output, child.failure.output...)
			}
		}
	}
}

// Search finds all occurrences of patterns in the text
// Time Complexity: O(n + m + z), where n is text length, m is total pattern length, z is number of matches
func (ac *AhoCorasick) Search(text string) map[string][]int {
	ac.BuildFailureFunction()

	result := make(map[string][]int)
	current := ac.root

	for i, char := range text {
		// Follow failure links until we find a match or reach root
		for current != nil && current.children[char] == nil {
			current = current.failure
		}

		if current == nil {
			current = ac.root
			continue
		}

		current = current.children[char]

		// Check for pattern matches
		for _, pattern := range current.output {
			if result[pattern] == nil {
				result[pattern] = make([]int, 0)
			}
			result[pattern] = append(result[pattern], i-len(pattern)+1)
		}
	}

	return result
}

// PalindromicTree (Eertree) for finding all palindromic substrings
type PalindromicTreeNode struct {
	edges  map[rune]*PalindromicTreeNode
	link   *PalindromicTreeNode
	length int
	count  int
}

// PalindromicTree implements Eertree data structure
type PalindromicTree struct {
	nodes []PalindromicTreeNode
	s     string
	ptr   int
}

// NewPalindromicTree creates a new palindromic tree
func NewPalindromicTree() *PalindromicTree {
	pt := &PalindromicTree{
		nodes: make([]PalindromicTreeNode, 2),
		s:     "",
		ptr:   1,
	}

	// Initialize odd and even length roots
	pt.nodes[0] = PalindromicTreeNode{
		edges:  make(map[rune]*PalindromicTreeNode),
		length: -1,
	}
	pt.nodes[1] = PalindromicTreeNode{
		edges:  make(map[rune]*PalindromicTreeNode),
		length: 0,
		link:   &pt.nodes[0],
	}

	return pt
}

// AddChar adds a character to the palindromic tree
func (pt *PalindromicTree) AddChar(c rune) {
	pt.s += string(c)
	pos := len(pt.s) - 1

	// Find the longest palindrome ending at current position
	cur := pt.ptr
	for {
		curLen := pt.nodes[cur].length
		if pos-curLen >= 1 && rune(pt.s[pos-curLen-1]) == c {
			break
		}
		cur = pt.getLink(cur)
	}

	// Check if palindrome already exists
	if edge, exists := pt.nodes[cur].edges[c]; exists {
		pt.ptr = pt.getNodeIndex(edge)
		return
	}

	// Create new node
	newNodeIndex := len(pt.nodes)
	pt.nodes = append(pt.nodes, PalindromicTreeNode{
		edges:  make(map[rune]*PalindromicTreeNode),
		length: pt.nodes[cur].length + 2,
		count:  0,
	})

	pt.nodes[cur].edges[c] = &pt.nodes[newNodeIndex]

	// Find suffix link
	if pt.nodes[newNodeIndex].length == 1 {
		pt.nodes[newNodeIndex].link = &pt.nodes[1]
	} else {
		cur = pt.getLink(cur)
		for {
			curLen := pt.nodes[cur].length
			if pos-curLen >= 1 && rune(pt.s[pos-curLen-1]) == c {
				pt.nodes[newNodeIndex].link = pt.nodes[cur].edges[c]
				break
			}
			cur = pt.getLink(cur)
		}
	}

	pt.ptr = newNodeIndex
}

func (pt *PalindromicTree) getLink(nodeIndex int) int {
	if pt.nodes[nodeIndex].link == nil {
		return 0
	}
	return pt.getNodeIndex(pt.nodes[nodeIndex].link)
}

func (pt *PalindromicTree) getNodeIndex(node *PalindromicTreeNode) int {
	for i := range pt.nodes {
		if &pt.nodes[i] == node {
			return i
		}
	}
	return 0
}

// CountPalindromes counts all palindromic substrings in a string
func CountPalindromes(s string) int {
	pt := NewPalindromicTree()
	for _, c := range s {
		pt.AddChar(c)
	}

	count := 0
	for i := 2; i < len(pt.nodes); i++ {
		count++
	}
	return count
}

// AllPalindromes returns all unique palindromic substrings
func AllPalindromes(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	palindromes := make(map[string]bool)
	n := len(s)

	// Check all possible centers for palindromes
	for center := 0; center < 2*n-1; center++ {
		left := center / 2
		right := left + center%2

		// Expand around center
		for left >= 0 && right < n && s[left] == s[right] {
			palindromes[s[left:right+1]] = true
			left--
			right++
		}
	}

	result := make([]string, 0, len(palindromes))
	for palindrome := range palindromes {
		result = append(result, palindrome)
	}

	sort.Strings(result)
	return result
}

// LongestCommonPrefix finds the longest common prefix of an array of strings
// Time Complexity: O(S), where S is sum of all characters, Space Complexity: O(1)
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	// Find minimum length
	minLen := len(strs[0])
	for _, str := range strs[1:] {
		if len(str) < minLen {
			minLen = len(str)
		}
	}

	// Check character by character
	for i := 0; i < minLen; i++ {
		char := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if strs[j][i] != char {
				return strs[0][:i]
			}
		}
	}

	return strs[0][:minLen]
}

// StringRotation checks if s2 is a rotation of s1
// Time Complexity: O(n), Space Complexity: O(n)
func StringRotation(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	if len(s1) == 0 {
		return true
	}

	// s2 is a rotation of s1 if s2 is a substring of s1+s1
	return strings.Contains(s1+s1, s2)
}
