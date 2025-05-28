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
// Time Complexity: O(n), Space Complexity: O(k) where k is unique characters
func MaxCharacter(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}
	
	charCount := make(map[rune]int)
	
	// Count characters
	for _, r := range s {
		charCount[r]++
	}
	
	// Find maximum
	var maxChar rune
	maxCount := 0
	
	for char, count := range charCount {
		if count > maxCount {
			maxCount = count
			maxChar = char
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
