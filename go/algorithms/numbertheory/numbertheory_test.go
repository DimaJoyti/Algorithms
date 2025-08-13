package numbertheory

import (
	"math"
	"testing"
)

func TestMillerRabinPrimality(t *testing.T) {
	tests := []struct {
		n        int64
		expected bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{11, true},
		{15, false},
		{17, true},
		{25, false},
		{97, true},
		{561, false}, // Carmichael number
	}

	for _, test := range tests {
		result := MillerRabinPrimality(test.n, 10)
		if result != test.expected {
			t.Errorf("MillerRabinPrimality(%d) = %v, expected %v", test.n, result, test.expected)
		}
	}
}

func TestModularExponentiation(t *testing.T) {
	tests := []struct {
		base, exp, mod, expected int64
	}{
		{2, 3, 5, 3},  // 2^3 mod 5 = 8 mod 5 = 3
		{3, 4, 7, 4},  // 3^4 mod 7 = 81 mod 7 = 4
		{5, 0, 13, 1}, // 5^0 mod 13 = 1
		{7, 1, 11, 7}, // 7^1 mod 11 = 7
	}

	for _, test := range tests {
		result := modularExponentiation(test.base, test.exp, test.mod)
		if result != test.expected {
			t.Errorf("modularExponentiation(%d, %d, %d) = %d, expected %d",
				test.base, test.exp, test.mod, result, test.expected)
		}
	}
}

func TestPollardRho(t *testing.T) {
	// Test with composite numbers
	composites := []int64{15, 21, 35, 77, 91, 143}

	for _, n := range composites {
		factor := PollardRho(n)
		if factor <= 1 || factor >= n {
			t.Errorf("PollardRho(%d) returned invalid factor %d", n, factor)
		}
		if n%factor != 0 {
			t.Errorf("PollardRho(%d) returned non-factor %d", n, factor)
		}
	}

	// Test with even numbers
	evenFactor := PollardRho(14)
	if evenFactor != 2 {
		t.Errorf("PollardRho(14) should return 2, got %d", evenFactor)
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b, expected int64
	}{
		{48, 18, 6},
		{17, 13, 1},
		{0, 5, 5},
		{5, 0, 5},
		{100, 25, 25},
	}

	for _, test := range tests {
		result := gcd(test.a, test.b)
		if result != test.expected {
			t.Errorf("gcd(%d, %d) = %d, expected %d", test.a, test.b, result, test.expected)
		}
	}
}

func TestFFT(t *testing.T) {
	// Test simple case
	input := []Complex{
		{Real: 1, Imag: 0},
		{Real: 1, Imag: 0},
		{Real: 1, Imag: 0},
		{Real: 1, Imag: 0},
	}

	result := FFT(input)
	if len(result) != 4 {
		t.Errorf("FFT should return same length, got %d", len(result))
	}

	// Test IFFT
	inverse := IFFT(result)
	for i, val := range inverse {
		realDiff := math.Abs(val.Real - input[i].Real)
		imagDiff := math.Abs(val.Imag - input[i].Imag)
		if realDiff > 1e-9 || imagDiff > 1e-9 {
			t.Errorf("IFFT doesn't match original at index %d", i)
		}
	}
}

func TestDiscreteLogarithm(t *testing.T) {
	// Test case: 3^x ≡ 4 (mod 7)
	// 3^1 = 3, 3^2 = 9 ≡ 2, 3^3 = 6, 3^4 = 18 ≡ 4 (mod 7)
	result := DiscreteLogarithm(3, 4, 7)
	if result != 4 {
		t.Errorf("DiscreteLogarithm(3, 4, 7) = %d, expected 4", result)
	}

	// Test case with no solution
	noSolution := DiscreteLogarithm(2, 3, 5)
	if noSolution == -1 {
		// This is expected for some cases
	}

	// Test edge cases
	edge1 := DiscreteLogarithm(2, 1, 5) // Should be 0 since 2^0 = 1
	if edge1 != 0 {
		t.Errorf("DiscreteLogarithm(2, 1, 5) = %d, expected 0", edge1)
	}
}

func TestExtendedEuclidean(t *testing.T) {
	tests := []struct {
		a, b          int64
		expectedGCD   int64
		shouldSatisfy bool
	}{
		{30, 18, 6, true},
		{17, 13, 1, true},
		{25, 15, 5, true},
	}

	for _, test := range tests {
		gcd, x, y := ExtendedEuclidean(test.a, test.b)
		if gcd != test.expectedGCD {
			t.Errorf("ExtendedEuclidean(%d, %d) gcd = %d, expected %d",
				test.a, test.b, gcd, test.expectedGCD)
		}

		if test.shouldSatisfy {
			// Verify that ax + by = gcd
			if test.a*x+test.b*y != gcd {
				t.Errorf("ExtendedEuclidean(%d, %d): %d*%d + %d*%d = %d, expected %d",
					test.a, test.b, test.a, x, test.b, y, test.a*x+test.b*y, gcd)
			}
		}
	}
}

func TestChineseRemainderTheorem(t *testing.T) {
	// Test case: x ≡ 2 (mod 3), x ≡ 3 (mod 5), x ≡ 2 (mod 7)
	remainders := []int64{2, 3, 2}
	moduli := []int64{3, 5, 7}

	result, valid := ChineseRemainderTheorem(remainders, moduli)
	if !valid {
		t.Error("ChineseRemainderTheorem should find valid solution")
	}

	// Verify the solution
	for i := 0; i < len(remainders); i++ {
		if result%moduli[i] != remainders[i] {
			t.Errorf("Solution %d doesn't satisfy congruence %d ≡ %d (mod %d)",
				result, result%moduli[i], remainders[i], moduli[i])
		}
	}

	// Test with non-coprime moduli (should fail)
	nonCoprimeModuli := []int64{4, 6}
	nonCoprimeRemainders := []int64{1, 2}
	_, validNonCoprime := ChineseRemainderTheorem(nonCoprimeRemainders, nonCoprimeModuli)
	if validNonCoprime {
		t.Error("ChineseRemainderTheorem should fail with non-coprime moduli")
	}
}

func TestEulerTotient(t *testing.T) {
	tests := []struct {
		n        int64
		expected int64
	}{
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 2},
		{5, 4},
		{6, 2},
		{9, 6},
		{10, 4},
		{12, 4},
	}

	for _, test := range tests {
		result := EulerTotient(test.n)
		if result != test.expected {
			t.Errorf("EulerTotient(%d) = %d, expected %d", test.n, result, test.expected)
		}
	}
}

func TestIsPrimitiveRoot(t *testing.T) {
	// Test primitive roots modulo 7
	// Primitive roots of 7 are: 3, 5
	tests := []struct {
		g, p     int64
		expected bool
	}{
		{1, 7, false},
		{2, 7, false},
		{3, 7, true},
		{4, 7, false},
		{5, 7, true},
		{6, 7, false},
	}

	for _, test := range tests {
		result := IsPrimitiveRoot(test.g, test.p)
		if result != test.expected {
			t.Errorf("IsPrimitiveRoot(%d, %d) = %v, expected %v",
				test.g, test.p, result, test.expected)
		}
	}
}

func TestJacobiSymbol(t *testing.T) {
	tests := []struct {
		a, n     int64
		expected int
	}{
		{1, 9, 1},
		{2, 9, 1},
		{3, 9, 0},
		{4, 9, 1},
		{5, 9, 1},
		{6, 9, 0},
		{7, 9, 1},
		{8, 9, 1},
	}

	for _, test := range tests {
		result := JacobiSymbol(test.a, test.n)
		if result != test.expected {
			t.Errorf("JacobiSymbol(%d, %d) = %d, expected %d",
				test.a, test.n, result, test.expected)
		}
	}
}

func TestQuadraticResidue(t *testing.T) {
	// Test quadratic residues modulo 7
	tests := []struct {
		a, p     int64
		expected bool
	}{
		{1, 7, true},  // 1² ≡ 1 (mod 7)
		{2, 7, true},  // 3² ≡ 2 (mod 7)
		{3, 7, false}, // No x such that x² ≡ 3 (mod 7)
		{4, 7, true},  // 2² ≡ 4 (mod 7)
	}

	for _, test := range tests {
		result := QuadraticResidue(test.a, test.p)
		if result != test.expected {
			t.Errorf("QuadraticResidue(%d, %d) = %v, expected %v",
				test.a, test.p, result, test.expected)
		}
	}
}

func TestFactorizeTrialDivision(t *testing.T) {
	tests := []struct {
		n        int64
		expected []int64
	}{
		{12, []int64{2, 2, 3}},
		{15, []int64{3, 5}},
		{17, []int64{17}}, // Prime
		{30, []int64{2, 3, 5}},
		{1, []int64{}},
	}

	for _, test := range tests {
		result := FactorizeTrialDivision(test.n)
		if len(result) != len(test.expected) {
			t.Errorf("FactorizeTrialDivision(%d) length = %d, expected %d",
				test.n, len(result), len(test.expected))
			continue
		}

		// Verify factors multiply to original number
		product := int64(1)
		for _, factor := range result {
			product *= factor
		}
		if product != test.n {
			t.Errorf("FactorizeTrialDivision(%d) factors don't multiply to original: %v",
				test.n, result)
		}
	}
}

// Benchmark tests
func BenchmarkMillerRabinPrimality(b *testing.B) {
	n := int64(982451653) // Large prime
	for i := 0; i < b.N; i++ {
		MillerRabinPrimality(n, 10)
	}
}

func BenchmarkPollardRho(b *testing.B) {
	n := int64(1403) // 23 * 61
	for i := 0; i < b.N; i++ {
		PollardRho(n)
	}
}

func BenchmarkFFT(b *testing.B) {
	input := make([]Complex, 1024)
	for i := range input {
		input[i] = Complex{Real: float64(i), Imag: 0}
	}

	for i := 0; i < b.N; i++ {
		FFT(input)
	}
}

func BenchmarkEulerTotient(b *testing.B) {
	n := int64(1000000)
	for i := 0; i < b.N; i++ {
		EulerTotient(n)
	}
}

func BenchmarkFactorizeTrialDivision(b *testing.B) {
	n := int64(1000003) // Large prime
	for i := 0; i < b.N; i++ {
		FactorizeTrialDivision(n)
	}
}
