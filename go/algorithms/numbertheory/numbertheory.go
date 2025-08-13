// Package numbertheory implements various number theory algorithms
package numbertheory

import (
	"crypto/rand"
	"math"
	"math/big"
	mathrand "math/rand"
	"time"
)

// MillerRabinPrimality performs Miller-Rabin primality test
// Time Complexity: O(k log³ n), Space Complexity: O(1)
func MillerRabinPrimality(n int64, k int) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Write n-1 as d * 2^r
	d := n - 1
	r := 0
	for d%2 == 0 {
		d /= 2
		r++
	}

	// Perform k rounds of testing
	mathrand.Seed(time.Now().UnixNano())
	for i := 0; i < k; i++ {
		a := mathrand.Int63n(n-3) + 2 // Random number in [2, n-2]
		x := modularExponentiation(a, d, n)

		if x == 1 || x == n-1 {
			continue
		}

		composite := true
		for j := 0; j < r-1; j++ {
			x = modularExponentiation(x, 2, n)
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

// modularExponentiation computes (base^exp) % mod efficiently
func modularExponentiation(base, exp, mod int64) int64 {
	result := int64(1)
	base = base % mod

	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % mod
		}
		exp = exp >> 1
		base = (base * base) % mod
	}

	return result
}

// PollardRho finds a non-trivial factor of n using Pollard's rho algorithm
// Time Complexity: O(n^(1/4)), Space Complexity: O(1)
func PollardRho(n int64) int64 {
	if n%2 == 0 {
		return 2
	}

	// Choose random starting values
	mathrand.Seed(time.Now().UnixNano())
	x := mathrand.Int63n(n-2) + 2
	y := x
	c := mathrand.Int63n(n-1) + 1
	d := int64(1)

	// Function f(x) = (x² + c) mod n
	f := func(x int64) int64 {
		return (x*x + c) % n
	}

	for d == 1 {
		x = f(x)
		y = f(f(y))
		d = gcd(abs(x-y), n)
	}

	if d == n {
		return PollardRho(n) // Try again with different parameters
	}

	return d
}

// gcd computes the greatest common divisor
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// abs returns the absolute value
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Complex represents a complex number for FFT
type Complex struct {
	Real, Imag float64
}

// Add adds two complex numbers
func (c Complex) Add(other Complex) Complex {
	return Complex{Real: c.Real + other.Real, Imag: c.Imag + other.Imag}
}

// Sub subtracts two complex numbers
func (c Complex) Sub(other Complex) Complex {
	return Complex{Real: c.Real - other.Real, Imag: c.Imag - other.Imag}
}

// Mul multiplies two complex numbers
func (c Complex) Mul(other Complex) Complex {
	return Complex{
		Real: c.Real*other.Real - c.Imag*other.Imag,
		Imag: c.Real*other.Imag + c.Imag*other.Real,
	}
}

// FFT performs Fast Fourier Transform
// Time Complexity: O(n log n), Space Complexity: O(n)
func FFT(a []Complex) []Complex {
	n := len(a)
	if n <= 1 {
		return a
	}

	// Ensure n is a power of 2
	if n&(n-1) != 0 {
		// Pad with zeros to next power of 2
		nextPow2 := 1
		for nextPow2 < n {
			nextPow2 <<= 1
		}
		padded := make([]Complex, nextPow2)
		copy(padded, a)
		return FFT(padded)
	}

	// Divide
	even := make([]Complex, n/2)
	odd := make([]Complex, n/2)

	for i := 0; i < n/2; i++ {
		even[i] = a[2*i]
		odd[i] = a[2*i+1]
	}

	// Conquer
	evenFFT := FFT(even)
	oddFFT := FFT(odd)

	// Combine
	result := make([]Complex, n)
	for i := 0; i < n/2; i++ {
		angle := -2 * math.Pi * float64(i) / float64(n)
		w := Complex{Real: math.Cos(angle), Imag: math.Sin(angle)}
		t := w.Mul(oddFFT[i])

		result[i] = evenFFT[i].Add(t)
		result[i+n/2] = evenFFT[i].Sub(t)
	}

	return result
}

// IFFT performs Inverse Fast Fourier Transform
// Time Complexity: O(n log n), Space Complexity: O(n)
func IFFT(a []Complex) []Complex {
	n := len(a)
	if n <= 1 {
		return a
	}

	// Conjugate the complex numbers
	for i := range a {
		a[i].Imag = -a[i].Imag
	}

	// Apply FFT
	result := FFT(a)

	// Conjugate again and scale
	for i := range result {
		result[i].Imag = -result[i].Imag
		result[i].Real /= float64(n)
		result[i].Imag /= float64(n)
	}

	return result
}

// DiscreteLogarithm solves g^x ≡ h (mod p) using baby-step giant-step algorithm
// Time Complexity: O(√p), Space Complexity: O(√p)
func DiscreteLogarithm(g, h, p int64) int64 {
	if p <= 1 {
		return -1
	}

	m := int64(math.Ceil(math.Sqrt(float64(p))))

	// Baby steps: compute g^j mod p for j = 0, 1, ..., m-1
	babySteps := make(map[int64]int64)
	gamma := int64(1)

	for j := int64(0); j < m; j++ {
		if gamma == h {
			return j
		}
		babySteps[gamma] = j
		gamma = (gamma * g) % p
	}

	// Giant steps: compute h * (g^(-m))^i mod p for i = 0, 1, ..., m-1
	gInvM := modularInverse(modularExponentiation(g, m, p), p)
	y := h

	for i := int64(0); i < m; i++ {
		if j, exists := babySteps[y]; exists {
			result := i*m + j
			if result < p {
				return result
			}
		}
		y = (y * gInvM) % p
	}

	return -1 // No solution found
}

// modularInverse computes the modular inverse of a modulo m using extended Euclidean algorithm
func modularInverse(a, m int64) int64 {
	if gcd(a, m) != 1 {
		return -1 // Inverse doesn't exist
	}

	// Extended Euclidean Algorithm
	m0 := m
	x0, x1 := int64(0), int64(1)

	if m == 1 {
		return 0
	}

	for a > 1 {
		q := a / m
		m, a = a%m, m
		x0, x1 = x1-q*x0, x0
	}

	if x1 < 0 {
		x1 += m0
	}

	return x1
}

// ExtendedEuclidean computes gcd(a, b) and finds x, y such that ax + by = gcd(a, b)
// Time Complexity: O(log min(a, b)), Space Complexity: O(1)
func ExtendedEuclidean(a, b int64) (int64, int64, int64) {
	if a == 0 {
		return b, 0, 1
	}

	gcd, x1, y1 := ExtendedEuclidean(b%a, a)
	x := y1 - (b/a)*x1
	y := x1

	return gcd, x, y
}

// ChineseRemainderTheorem solves system of congruences using CRT
// Time Complexity: O(n²), Space Complexity: O(1)
func ChineseRemainderTheorem(remainders, moduli []int64) (int64, bool) {
	if len(remainders) != len(moduli) {
		return 0, false
	}

	n := len(remainders)
	if n == 0 {
		return 0, true
	}

	// Check if moduli are pairwise coprime
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if gcd(moduli[i], moduli[j]) != 1 {
				return 0, false
			}
		}
	}

	// Compute product of all moduli
	M := int64(1)
	for _, mod := range moduli {
		M *= mod
	}

	result := int64(0)
	for i := 0; i < n; i++ {
		Mi := M / moduli[i]
		yi := modularInverse(Mi, moduli[i])
		if yi == -1 {
			return 0, false
		}
		result = (result + remainders[i]*Mi*yi) % M
	}

	if result < 0 {
		result += M
	}

	return result, true
}

// EulerTotient computes Euler's totient function φ(n)
// Time Complexity: O(√n), Space Complexity: O(1)
func EulerTotient(n int64) int64 {
	if n <= 1 {
		return n
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
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			result -= result / i
			for n%i == 0 {
				n /= i
			}
		}
	}

	// If n is still greater than 1, it's a prime
	if n > 1 {
		result -= result / n
	}

	return result
}

// IsPrimitivRoot checks if g is a primitive root modulo p
// Time Complexity: O(√φ(p) * log p), Space Complexity: O(√φ(p))
func IsPrimitiveRoot(g, p int64) bool {
	if gcd(g, p) != 1 {
		return false
	}

	phi := EulerTotient(p)

	// Find all prime factors of phi
	factors := make([]int64, 0)
	temp := phi

	for i := int64(2); i*i <= temp; i++ {
		if temp%i == 0 {
			factors = append(factors, i)
			for temp%i == 0 {
				temp /= i
			}
		}
	}
	if temp > 1 {
		factors = append(factors, temp)
	}

	// Check if g^(phi/factor) ≢ 1 (mod p) for all prime factors
	for _, factor := range factors {
		if modularExponentiation(g, phi/factor, p) == 1 {
			return false
		}
	}

	return true
}

// JacobiSymbol computes the Jacobi symbol (a/n)
// Time Complexity: O(log n), Space Complexity: O(1)
func JacobiSymbol(a, n int64) int {
	if n <= 0 || n%2 == 0 {
		return 0
	}

	a = a % n
	result := 1

	for a != 0 {
		for a%2 == 0 {
			a /= 2
			if n%8 == 3 || n%8 == 5 {
				result = -result
			}
		}

		a, n = n, a
		if a%4 == 3 && n%4 == 3 {
			result = -result
		}
		a = a % n
	}

	if n == 1 {
		return result
	}
	return 0
}

// QuadraticResidue checks if a is a quadratic residue modulo p using Legendre symbol
// Time Complexity: O(log p), Space Complexity: O(1)
func QuadraticResidue(a, p int64) bool {
	if p == 2 {
		return a%2 == 0
	}

	// Use Legendre symbol: a^((p-1)/2) mod p
	legendre := modularExponentiation(a, (p-1)/2, p)
	return legendre == 1
}

// CarmichaelFunction computes the Carmichael function λ(n)
// Time Complexity: O(√n), Space Complexity: O(log n)
func CarmichaelFunction(n int64) int64 {
	if n <= 1 {
		return 1
	}

	// Factor n into prime powers
	factors := make(map[int64]int64)
	temp := n

	for i := int64(2); i*i <= temp; i++ {
		for temp%i == 0 {
			factors[i]++
			temp /= i
		}
	}
	if temp > 1 {
		factors[temp] = 1
	}

	result := int64(1)

	for prime, power := range factors {
		var lambda int64
		if prime == 2 && power >= 3 {
			lambda = int64(math.Pow(2, float64(power-2)))
		} else {
			lambda = int64(math.Pow(float64(prime), float64(power-1))) * (prime - 1)
		}
		result = lcm(result, lambda)
	}

	return result
}

// lcm computes the least common multiple
func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

// PrimitiveRoots finds all primitive roots modulo p
// Time Complexity: O(p * √φ(p) * log p), Space Complexity: O(φ(p))
func PrimitiveRoots(p int64) []int64 {
	if p <= 1 {
		return []int64{}
	}

	roots := make([]int64, 0)
	for g := int64(1); g < p; g++ {
		if IsPrimitiveRoot(g, p) {
			roots = append(roots, g)
		}
	}

	return roots
}

// BigIntMillerRabin performs Miller-Rabin test on big integers
func BigIntMillerRabin(n *big.Int, k int) bool {
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// Write n-1 as d * 2^r
	nMinus1 := new(big.Int).Sub(n, big.NewInt(1))
	d := new(big.Int).Set(nMinus1)
	r := 0

	for new(big.Int).Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		d.Div(d, big.NewInt(2))
		r++
	}

	// Perform k rounds of testing
	for i := 0; i < k; i++ {
		// Generate random a in [2, n-2]
		maxVal := new(big.Int).Sub(n, big.NewInt(3))
		a, _ := rand.Int(rand.Reader, maxVal)
		a.Add(a, big.NewInt(2))

		x := new(big.Int).Exp(a, d, n)

		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(nMinus1) == 0 {
			continue
		}

		composite := true
		for j := 0; j < r-1; j++ {
			x.Exp(x, big.NewInt(2), n)
			if x.Cmp(nMinus1) == 0 {
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

// GeneratePrime generates a random prime number of specified bit length
// Time Complexity: O(k * b³), Space Complexity: O(b) where b is bit length
func GeneratePrime(bitLength int) *big.Int {
	for {
		// Generate random odd number of specified bit length
		maxVal := new(big.Int).Lsh(big.NewInt(1), uint(bitLength))
		candidate, _ := rand.Int(rand.Reader, maxVal)
		candidate.SetBit(candidate, 0, 1)           // Make it odd
		candidate.SetBit(candidate, bitLength-1, 1) // Ensure it has the right bit length

		if BigIntMillerRabin(candidate, 20) {
			return candidate
		}
	}
}

// FactorizeTrialDivision performs trial division factorization
// Time Complexity: O(√n), Space Complexity: O(log n)
func FactorizeTrialDivision(n int64) []int64 {
	if n <= 1 {
		return []int64{}
	}

	factors := make([]int64, 0)

	// Check for factor 2
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}

	// Check for odd factors
	for i := int64(3); i*i <= n; i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}

	// If n is still greater than 1, it's a prime
	if n > 1 {
		factors = append(factors, n)
	}

	return factors
}
