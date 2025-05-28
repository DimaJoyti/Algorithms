package strings

import (
	"reflect"
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single character", "a", "a"},
		{"Simple string", "hello", "olleh"},
		{"String with spaces", "hello world", "dlrow olleh"},
		{"Unicode string", "café", "éfac"},
		{"Palindrome", "racecar", "racecar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseString(tt.input)
			if result != tt.expected {
				t.Errorf("ReverseString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Single character", "a", true},
		{"Simple palindrome", "racecar", true},
		{"Not palindrome", "hello", false},
		{"Palindrome with spaces", "A man a plan a canal Panama", true},
		{"Palindrome with punctuation", "race a car", false},
		{"Mixed case palindrome", "RaceCar", true},
		{"Numbers and letters", "A1B2b1a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAreAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected bool
	}{
		{"Empty strings", "", "", true},
		{"Simple anagrams", "listen", "silent", true},
		{"Not anagrams", "hello", "world", false},
		{"Different lengths", "abc", "abcd", false},
		{"Case insensitive", "Listen", "Silent", true},
		{"Same string", "test", "test", true},
		{"Single characters", "a", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AreAnagrams(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("AreAnagrams(%q, %q) = %v, want %v", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestAreAnagramsOptimized(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected bool
	}{
		{"Empty strings", "", "", true},
		{"Simple anagrams", "listen", "silent", true},
		{"Not anagrams", "hello", "world", false},
		{"Different lengths", "abc", "abcd", false},
		{"Case insensitive", "Listen", "Silent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AreAnagramsOptimized(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("AreAnagramsOptimized(%q, %q) = %v, want %v", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestCountVowels(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"No vowels", "bcdfg", 0},
		{"All vowels", "aeiou", 5},
		{"Mixed case", "Hello World", 3},
		{"With numbers", "a1e2i3o4u5", 5},
		{"Unicode", "café", 1}, // Only 'a' is a vowel, 'é' is not in our vowels string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountVowels(tt.input)
			if result != tt.expected {
				t.Errorf("CountVowels(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaxCharacter(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedChar  rune
		expectedCount int
	}{
		{"Empty string", "", 0, 0},
		{"Single character", "a", 'a', 1},
		{"Multiple same", "aaa", 'a', 3},
		{"Mixed characters", "abccba", 'a', 2}, // Both 'a' and 'c' have count 2, but 'a' appears first
		{"With spaces", "hello world", 'l', 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, count := MaxCharacter(tt.input)
			if char != tt.expectedChar || count != tt.expectedCount {
				t.Errorf("MaxCharacter(%q) = (%c, %d), want (%c, %d)",
					tt.input, char, count, tt.expectedChar, tt.expectedCount)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "Hello"},
		{"Multiple words", "hello world", "Hello World"},
		{"Already capitalized", "Hello World", "Hello World"},
		{"Mixed case", "hELLo WoRLd", "Hello World"},
		{"Extra spaces", "  hello   world  ", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Capitalize(tt.input)
			if result != tt.expected {
				t.Errorf("Capitalize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLongestCommonSubstring(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected string
	}{
		{"Empty strings", "", "", ""},
		{"One empty", "abc", "", ""},
		{"No common", "abc", "def", ""},
		{"Common substring", "GeeksforGeeks", "GeeksQuiz", "Geeks"},
		{"Identical strings", "test", "test", "test"},
		{"Overlapping", "abcdxyz", "xyzabcd", "abcd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LongestCommonSubstring(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("LongestCommonSubstring(%q, %q) = %q, want %q", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"Empty strings", "", "", 0},
		{"One empty", "abc", "", 0},
		{"No common", "abc", "def", 0},
		{"Common subsequence", "ABCDGH", "AEDFHR", 3},
		{"Identical strings", "test", "test", 4},
		{"Single character", "a", "a", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LongestCommonSubsequence(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want %d", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"Empty strings", "", "", 0},
		{"One empty", "abc", "", 3},
		{"Identical strings", "test", "test", 0},
		{"Single substitution", "cat", "bat", 1},
		{"Multiple operations", "kitten", "sitting", 3},
		{"Complete different", "abc", "def", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EditDistance(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("EditDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestKMPSearch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected []int
	}{
		{"Empty pattern", "hello", "", []int{}},
		{"Pattern not found", "hello", "world", []int{}},
		{"Single match", "hello", "ell", []int{1}},
		{"Multiple matches", "ababcababa", "aba", []int{0, 5, 7}},
		{"Pattern equals text", "test", "test", []int{0}},
		{"Overlapping matches", "aaaa", "aa", []int{0, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KMPSearch(tt.text, tt.pattern)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("KMPSearch(%q, %q) = %v, want %v", tt.text, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestRabinKarpSearch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected []int
	}{
		{"Pattern not found", "hello", "world", []int{}},
		{"Single match", "hello", "ell", []int{1}},
		{"Multiple matches", "ababcababa", "aba", []int{0, 5, 7}},
		{"Pattern equals text", "test", "test", []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RabinKarpSearch(tt.text, tt.pattern)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RabinKarpSearch(%q, %q) = %v, want %v", tt.text, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestIsSubsequence(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		t        string
		expected bool
	}{
		{"Empty subsequence", "", "abc", true},
		{"Empty target", "a", "", false},
		{"Valid subsequence", "ace", "abcde", true},
		{"Invalid subsequence", "aec", "abcde", false},
		{"Identical strings", "abc", "abc", true},
		{"Single character", "a", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSubsequence(tt.s, tt.t)
			if result != tt.expected {
				t.Errorf("IsSubsequence(%q, %q) = %v, want %v", tt.s, tt.t, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkReverseString(b *testing.B) {
	s := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReverseString(s)
	}
}

func BenchmarkAreAnagrams(b *testing.B) {
	s1, s2 := "listen", "silent"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AreAnagrams(s1, s2)
	}
}

func BenchmarkAreAnagramsOptimized(b *testing.B) {
	s1, s2 := "listen", "silent"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AreAnagramsOptimized(s1, s2)
	}
}
