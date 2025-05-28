package math

import (
	"reflect"
	"testing"
)

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Negative input", -1, 0},
		{"Zero", 0, 0},
		{"One", 1, 1},
		{"Two", 2, 1},
		{"Five", 5, 5},
		{"Ten", 10, 55},
		{"Fifteen", 15, 610},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Fibonacci(tt.input)
			if result != tt.expected {
				t.Errorf("Fibonacci(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFibonacciRecursive(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 0},
		{"One", 1, 1},
		{"Five", 5, 5},
		{"Eight", 8, 21}, // Keep small for recursive version
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciRecursive(tt.input)
			if result != tt.expected {
				t.Errorf("FibonacciRecursive(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFibonacciMemoized(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 0},
		{"One", 1, 1},
		{"Ten", 10, 55},
		{"Twenty", 20, 6765},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FibonacciMemoized(tt.input)
			if result != tt.expected {
				t.Errorf("FibonacciMemoized(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Negative", -5, false},
		{"Zero", 0, false},
		{"One", 1, false},
		{"Two", 2, true},
		{"Three", 3, true},
		{"Four", 4, false},
		{"Five", 5, true},
		{"Nine", 9, false},
		{"Eleven", 11, true},
		{"Fifteen", 15, false},
		{"Seventeen", 17, true},
		{"Large prime", 97, true},
		{"Large composite", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrime(tt.input)
			if result != tt.expected {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSieveOfEratosthenes(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected []int
	}{
		{"Zero", 0, []int{}},
		{"One", 1, []int{}},
		{"Ten", 10, []int{2, 3, 5, 7}},
		{"Twenty", 20, []int{2, 3, 5, 7, 11, 13, 17, 19}},
		{"Thirty", 30, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SieveOfEratosthenes(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SieveOfEratosthenes(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Both zero", 0, 0, 0},
		{"One zero", 5, 0, 5},
		{"Same numbers", 12, 12, 12},
		{"Coprime", 17, 13, 1},
		{"Common factors", 48, 18, 6},
		{"Negative numbers", -48, 18, 6},
		{"Large numbers", 1071, 462, 21},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GCD(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"One zero", 5, 0, 0},
		{"Same numbers", 12, 12, 12},
		{"Coprime", 7, 13, 91},
		{"Common factors", 12, 18, 36},
		{"Powers of 2", 8, 12, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LCM(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("LCM(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		exponent int
		expected int
	}{
		{"Zero exponent", 5, 0, 1},
		{"Negative exponent", 5, -2, 0},
		{"Base zero", 0, 5, 0},
		{"Base one", 1, 100, 1},
		{"Small power", 2, 3, 8},
		{"Large power", 3, 4, 81},
		{"Negative base", -2, 3, 8}, // abs(-2) = 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Power(tt.base, tt.exponent)
			if result != tt.expected {
				t.Errorf("Power(%d, %d) = %d, want %d", tt.base, tt.exponent, result, tt.expected)
			}
		})
	}
}

func TestPowerMod(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		exponent int
		modulus  int
		expected int
	}{
		{"Simple case", 2, 3, 5, 3}, // 2^3 % 5 = 8 % 5 = 3
		{"Large exponent", 2, 10, 1000, 24}, // 2^10 % 1000 = 1024 % 1000 = 24
		{"Modulus 1", 5, 3, 1, 0},
		{"Base larger than modulus", 10, 2, 7, 2}, // 10^2 % 7 = 100 % 7 = 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PowerMod(tt.base, tt.exponent, tt.modulus)
			if result != tt.expected {
				t.Errorf("PowerMod(%d, %d, %d) = %d, want %d", tt.base, tt.exponent, tt.modulus, result, tt.expected)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Negative", -1, 0},
		{"Zero", 0, 1},
		{"One", 1, 1},
		{"Five", 5, 120},
		{"Seven", 7, 5040},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Factorial(tt.input)
			if result != tt.expected {
				t.Errorf("Factorial(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFactorialRecursive(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 1},
		{"One", 1, 1},
		{"Five", 5, 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FactorialRecursive(tt.input)
			if result != tt.expected {
				t.Errorf("FactorialRecursive(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCombination(t *testing.T) {
	tests := []struct {
		name     string
		n, r     int
		expected int
	}{
		{"Invalid r negative", 5, -1, 0},
		{"Invalid r > n", 5, 6, 0},
		{"r = 0", 5, 0, 1},
		{"r = n", 5, 5, 1},
		{"C(5,2)", 5, 2, 10},
		{"C(10,3)", 10, 3, 120},
		{"C(6,3)", 6, 3, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Combination(tt.n, tt.r)
			if result != tt.expected {
				t.Errorf("Combination(%d, %d) = %d, want %d", tt.n, tt.r, result, tt.expected)
			}
		})
	}
}

func TestPermutation(t *testing.T) {
	tests := []struct {
		name     string
		n, r     int
		expected int
	}{
		{"Invalid r negative", 5, -1, 0},
		{"Invalid r > n", 5, 6, 0},
		{"r = 0", 5, 0, 1},
		{"P(5,2)", 5, 2, 20},
		{"P(10,3)", 10, 3, 720},
		{"P(6,6)", 6, 6, 720},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Permutation(tt.n, tt.r)
			if result != tt.expected {
				t.Errorf("Permutation(%d, %d) = %d, want %d", tt.n, tt.r, result, tt.expected)
			}
		})
	}
}

func TestSquareRoot(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Negative", -4, 0},
		{"Zero", 0, 0},
		{"One", 1, 1},
		{"Perfect square", 16, 4},
		{"Non-perfect square", 15, 3},
		{"Large number", 100, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SquareRoot(tt.input)
			if result != tt.expected {
				t.Errorf("SquareRoot(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPerfectSquare(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Negative", -4, false},
		{"Zero", 0, true},
		{"One", 1, true},
		{"Perfect square", 16, true},
		{"Non-perfect square", 15, false},
		{"Large perfect square", 144, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPerfectSquare(tt.input)
			if result != tt.expected {
				t.Errorf("IsPerfectSquare(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDigitSum(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 0},
		{"Single digit", 7, 7},
		{"Multiple digits", 123, 6},
		{"Negative", -123, 6},
		{"Large number", 9876, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DigitSum(tt.input)
			if result != tt.expected {
				t.Errorf("DigitSum(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReverseInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 0},
		{"Single digit", 7, 7},
		{"Multiple digits", 123, 321},
		{"Negative", -123, -321},
		{"Trailing zeros", 1200, 21},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReverseInteger(tt.input)
			if result != tt.expected {
				t.Errorf("ReverseInteger(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPalindromeNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Negative", -121, false},
		{"Single digit", 7, true},
		{"Palindrome", 121, true},
		{"Not palindrome", 123, false},
		{"Large palindrome", 12321, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindromeNumber(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindromeNumber(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountDigits(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Zero", 0, 1},
		{"Single digit", 7, 1},
		{"Multiple digits", 123, 3},
		{"Negative", -123, 3},
		{"Large number", 123456, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountDigits(tt.input)
			if result != tt.expected {
				t.Errorf("CountDigits(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkFibonacciMemoized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciMemoized(20)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(97)
	}
}

func BenchmarkSieveOfEratosthenes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SieveOfEratosthenes(1000)
	}
}
