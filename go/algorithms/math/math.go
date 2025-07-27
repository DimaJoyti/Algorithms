// Package math implements various mathematical algorithms
package math

import (
	"errors"
)

// Fibonacci calculates the nth Fibonacci number
// Time Complexity: O(n), Space Complexity: O(1)
func Fibonacci(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// FibonacciRecursive calculates Fibonacci recursively (inefficient)
// Time Complexity: O(2^n), Space Complexity: O(n)
func FibonacciRecursive(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

// FibonacciMemoized calculates Fibonacci with memoization
// Time Complexity: O(n), Space Complexity: O(n)
func FibonacciMemoized(n int) int {
	memo := make(map[int]int)
	return fibMemo(n, memo)
}

func fibMemo(n int, memo map[int]int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return n
	}

	if val, exists := memo[n]; exists {
		return val
	}

	memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
	return memo[n]
}

// IsPrime checks if a number is prime
// Time Complexity: O(√n), Space Complexity: O(1)
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// SieveOfEratosthenes finds all primes up to n
// Time Complexity: O(n log log n), Space Complexity: O(n)
func SieveOfEratosthenes(n int) []int {
	if n < 2 {
		return []int{}
	}

	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	var primes []int
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

// GCD calculates Greatest Common Divisor using Euclidean algorithm
// Time Complexity: O(log min(a,b)), Space Complexity: O(1)
func GCD(a, b int) int {
	a, b = abs(a), abs(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM calculates Least Common Multiple
// Time Complexity: O(log min(a,b)), Space Complexity: O(1)
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return abs(a*b) / GCD(a, b)
}

// Power calculates base^exponent efficiently
// Time Complexity: O(log n), Space Complexity: O(1)
func Power(base, exponent int) int {
	if exponent == 0 {
		return 1
	}
	if exponent < 0 {
		return 0 // For integer division
	}

	result := 1
	base = abs(base)

	for exponent > 0 {
		if exponent%2 == 1 {
			result *= base
		}
		base *= base
		exponent /= 2
	}

	return result
}

// PowerMod calculates (base^exponent) % modulus efficiently
// Time Complexity: O(log n), Space Complexity: O(1)
func PowerMod(base, exponent, modulus int) int {
	if modulus == 1 {
		return 0
	}

	result := 1
	base = base % modulus

	for exponent > 0 {
		if exponent%2 == 1 {
			result = (result * base) % modulus
		}
		exponent = exponent >> 1
		base = (base * base) % modulus
	}

	return result
}

// Factorial calculates n!
// Time Complexity: O(n), Space Complexity: O(1)
func Factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return 1
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// FactorialRecursive calculates n! recursively
// Time Complexity: O(n), Space Complexity: O(n)
func FactorialRecursive(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return 1
	}
	return n * FactorialRecursive(n-1)
}

// Combination calculates C(n, r) = n! / (r! * (n-r)!)
// Time Complexity: O(min(r, n-r)), Space Complexity: O(1)
func Combination(n, r int) int {
	if r < 0 || r > n {
		return 0
	}
	if r == 0 || r == n {
		return 1
	}

	// Use symmetry: C(n, r) = C(n, n-r)
	if r > n-r {
		r = n - r
	}

	result := 1
	for i := 0; i < r; i++ {
		result = result * (n - i) / (i + 1)
	}

	return result
}

// Permutation calculates P(n, r) = n! / (n-r)!
// Time Complexity: O(r), Space Complexity: O(1)
func Permutation(n, r int) int {
	if r < 0 || r > n {
		return 0
	}
	if r == 0 {
		return 1
	}

	result := 1
	for i := 0; i < r; i++ {
		result *= (n - i)
	}

	return result
}

// SquareRoot calculates integer square root using Newton's method
// Time Complexity: O(log n), Space Complexity: O(1)
func SquareRoot(n int) int {
	if n < 0 {
		return 0
	}
	if n < 2 {
		return n
	}

	x := n
	for {
		y := (x + n/x) / 2
		if y >= x {
			return x
		}
		x = y
	}
}

// IsPerfectSquare checks if n is a perfect square
// Time Complexity: O(log n), Space Complexity: O(1)
func IsPerfectSquare(n int) bool {
	if n < 0 {
		return false
	}
	root := SquareRoot(n)
	return root*root == n
}

// DigitSum calculates sum of digits
// Time Complexity: O(log n), Space Complexity: O(1)
func DigitSum(n int) int {
	n = abs(n)
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

// ReverseInteger reverses the digits of an integer
// Time Complexity: O(log n), Space Complexity: O(1)
func ReverseInteger(n int) int {
	sign := 1
	if n < 0 {
		sign = -1
		n = -n
	}

	result := 0
	for n > 0 {
		result = result*10 + n%10
		n /= 10
	}

	return result * sign
}

// IsPalindrome checks if a number is palindrome
// Time Complexity: O(log n), Space Complexity: O(1)
func IsPalindromeNumber(n int) bool {
	if n < 0 {
		return false
	}
	return n == ReverseInteger(n)
}

// CountDigits counts number of digits
// Time Complexity: O(log n), Space Complexity: O(1)
func CountDigits(n int) int {
	if n == 0 {
		return 1
	}

	n = abs(n)
	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}

// ExtendedGCD calculates GCD and coefficients x, y such that ax + by = gcd(a, b)
// Time Complexity: O(log min(a,b)), Space Complexity: O(1)
func ExtendedGCD(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}

	gcd1, x1, y1 := ExtendedGCD(b, a%b)
	x = y1
	y = x1 - (a/b)*y1

	return gcd1, x, y
}

// ModularInverse calculates modular multiplicative inverse
// Time Complexity: O(log m), Space Complexity: O(1)
func ModularInverse(a, m int) (int, error) {
	gcd, x, _ := ExtendedGCD(a, m)
	if gcd != 1 {
		return 0, errors.New("modular inverse does not exist")
	}

	// Make x positive
	result := (x%m + m) % m
	return result, nil
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FastExponentiation computes base^exp % mod efficiently
// Time Complexity: O(log exp), Space Complexity: O(1)
func FastExponentiation(base, exp, mod int) int {
	if mod == 1 {
		return 0
	}

	result := 1
	base = base % mod

	for exp > 0 {
		// If exp is odd, multiply base with result
		if exp%2 == 1 {
			result = (result * base) % mod
		}

		// Now exp must be even
		exp = exp >> 1 // exp = exp / 2
		base = (base * base) % mod
	}

	return result
}

// FastExponentiationNonMod computes base^exp without modulo
// Time Complexity: O(log exp), Space Complexity: O(1)
func FastExponentiationNonMod(base, exp int) int {
	if exp == 0 {
		return 1
	}
	if exp == 1 {
		return base
	}

	result := 1
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp >>= 1
	}

	return result
}

// MillerRabinPrimality performs Miller-Rabin primality test
// Time Complexity: O(k log³ n), Space Complexity: O(1) where k is number of rounds
func MillerRabinPrimality(n int, k int) bool {
	if n <= 1 || n == 4 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Write n-1 as d * 2^r by factoring powers of 2 from n-1
	d := n - 1
	r := 0
	for d%2 == 0 {
		d /= 2
		r++
	}

	// Witness loop
	for i := 0; i < k; i++ {
		a := 2 + (i*7)%(n-3) // Simple deterministic sequence for testing
		if a >= n {
			a = 2
		}

		x := FastExponentiation(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}

		composite := true
		for j := 0; j < r-1; j++ {
			x = FastExponentiation(x, 2, n)
			if x == n-1 {
				composite = false
				break
			}
		}

		if composite {
			return false
		}
	}

	return true
}

// IsPrimeMR checks if a number is prime using Miller-Rabin test
// Time Complexity: O(log³ n), Space Complexity: O(1)
func IsPrimeMR(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	// Use Miller-Rabin with 5 rounds for good accuracy
	return MillerRabinPrimality(n, 5)
}

// EulerTotient computes Euler's totient function φ(n)
// Time Complexity: O(√n), Space Complexity: O(1)
func EulerTotient(n int) int {
	if n <= 1 {
		return 0
	}

	result := n

	// Check for factor 2
	if n%2 == 0 {
		result -= result / 2
		for n%2 == 0 {
			n /= 2
		}
	}

	// Check for odd factors
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			result -= result / i
			for n%i == 0 {
				n /= i
			}
		}
	}

	// If n is still greater than 1, then it's a prime
	if n > 1 {
		result -= result / n
	}

	return result
}
