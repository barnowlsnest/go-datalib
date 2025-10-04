package serial

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// NSumTestSuite tests the NSum hash function
type NSumTestSuite struct {
	suite.Suite
}

func (s *NSumTestSuite) TestNSum_Consistency() {
	from := uint64(123)
	to := uint64(456)

	hash1 := NSum(from, to)
	hash2 := NSum(from, to)
	hash3 := NSum(from, to)

	assert.Equal(s.T(), hash1, hash2, "same inputs should produce same hash")
	assert.Equal(s.T(), hash1, hash3, "same inputs should produce same hash")
}

func (s *NSumTestSuite) TestNSum_Deterministic() {
	// Get actual hash values
	hash00 := NSum(0, 0)
	hash11 := NSum(1, 1)
	hash100_200 := NSum(100, 200)

	// Verify they're deterministic by computing again
	assert.Equal(s.T(), hash00, NSum(0, 0))
	assert.Equal(s.T(), hash11, NSum(1, 1))
	assert.Equal(s.T(), hash100_200, NSum(100, 200))
}

func (s *NSumTestSuite) TestNSum_OrderMatters() {
	// NSum should produce different results for different order
	hash1 := NSum(123, 456)
	hash2 := NSum(456, 123)

	assert.NotEqual(s.T(), hash1, hash2,
		"NSum(a, b) should differ from NSum(b, a)")
}

func (s *NSumTestSuite) TestNSum_ZeroValues() {
	hash00 := NSum(0, 0)
	hash01 := NSum(0, 1)
	hash10 := NSum(1, 0)

	// NSum(0, 0) produces magic constant, not zero
	assert.NotEqual(s.T(), hash00, hash01, "NSum(0, 1) should differ from NSum(0, 0)")
	assert.NotEqual(s.T(), hash00, hash10, "NSum(1, 0) should differ from NSum(0, 0)")
	assert.NotEqual(s.T(), hash01, hash10, "NSum(0, 1) should differ from NSum(1, 0)")
}

func (s *NSumTestSuite) TestNSum_MaxValues() {
	maxUint := uint64(18446744073709551615) // max uint64

	hash1 := NSum(maxUint, maxUint)
	hash2 := NSum(maxUint, 0)
	hash3 := NSum(0, maxUint)

	// All should produce valid hashes
	assert.Greater(s.T(), hash1, uint64(0))
	assert.Greater(s.T(), hash2, uint64(0))
	assert.Greater(s.T(), hash3, uint64(0))

	// All should be different
	assert.NotEqual(s.T(), hash1, hash2)
	assert.NotEqual(s.T(), hash1, hash3)
	assert.NotEqual(s.T(), hash2, hash3)
}

func (s *NSumTestSuite) TestNSum_Distribution() {
	// Test that different inputs produce well-distributed hashes
	hashes := make(map[uint64]bool)
	collisions := 0

	for i := uint64(0); i < 1000; i++ {
		hash := NSum(i, i+1)
		if hashes[hash] {
			collisions++
		}
		hashes[hash] = true
	}

	// Expect very few or no collisions
	assert.LessOrEqual(s.T(), collisions, 5,
		"should have minimal collisions in 1000 sequential inputs")
	assert.GreaterOrEqual(s.T(), len(hashes), 995,
		"should have good distribution (at least 995 unique hashes)")
}

func (s *NSumTestSuite) TestNSum_PairUniqueness() {
	// Test that different pairs produce unique hashes
	type pair struct {
		from, to uint64
	}

	pairs := []pair{
		{1, 2}, {2, 1},
		{10, 20}, {20, 10},
		{100, 200}, {200, 100},
		{1000, 2000}, {2000, 1000},
	}

	hashes := make(map[uint64]pair)
	for _, p := range pairs {
		hash := NSum(p.from, p.to)

		if existing, found := hashes[hash]; found {
			s.T().Errorf("collision: NSum(%d, %d) == NSum(%d, %d) = %d",
				p.from, p.to, existing.from, existing.to, hash)
		}
		hashes[hash] = p
	}

	assert.Equal(s.T(), len(pairs), len(hashes),
		"all pairs should produce unique hashes")
}

// HashcodeTestSuite tests the internal hashcode function
type HashcodeTestSuite struct {
	suite.Suite
}

func (s *HashcodeTestSuite) TestHashcode_Consistency() {
	base := uint64(100)
	val := uint64(200)

	hash1 := hashcode(base, val)
	hash2 := hashcode(base, val)

	assert.Equal(s.T(), hash1, hash2,
		"hashcode should be consistent for same inputs")
}

func (s *HashcodeTestSuite) TestHashcode_GoldenRatio() {
	// The magic constant should be the golden ratio
	expectedMagic := uint64(0x9e3779b97f4a7c15)
	assert.Equal(s.T(), expectedMagic, magic,
		"magic constant should be golden ratio (64-bit)")
}

func (s *HashcodeTestSuite) TestHashcode_Commutative() {
	// hashcode is NOT commutative (order matters)
	base := uint64(123)
	val := uint64(456)

	hash1 := hashcode(base, val)
	hash2 := hashcode(val, base)

	assert.NotEqual(s.T(), hash1, hash2,
		"hashcode should not be commutative")
}

func (s *HashcodeTestSuite) TestHashcode_ZeroBase() {
	val := uint64(100)

	hash := hashcode(0, val)
	expected := val + magic // when base=0: 0 ^ (val + magic + 0 + 0) = val + magic

	assert.Equal(s.T(), expected, hash,
		"hashcode with zero base should equal val + magic")
}

func (s *HashcodeTestSuite) TestHashcode_BitOperations() {
	// Verify the bit operations work as expected
	base := uint64(8) // Binary: 1000
	val := uint64(0)

	// base << 6 = 512 (1000000000 in binary)
	// base >> 2 = 2   (10 in binary)
	// hashcode = 8 ^ (0 + magic + 512 + 2) = 8 ^ (magic + 514)

	hash := hashcode(base, val)
	expected := base ^ (val + magic + (base << 6) + (base >> 2))

	assert.Equal(s.T(), expected, hash,
		"hashcode should correctly apply bit operations")
}

// EdgeCasesTestSuite tests edge cases for NSum
type EdgeCasesNSumTestSuite struct {
	suite.Suite
}

func (s *EdgeCasesNSumTestSuite) TestNSum_SequentialPairs() {
	// Test sequential pairs like (0,1), (1,2), (2,3), etc.
	hashes := make(map[uint64]bool)

	for i := uint64(0); i < 100; i++ {
		hash := NSum(i, i+1)
		assert.False(s.T(), hashes[hash],
			"sequential pairs should produce unique hashes")
		hashes[hash] = true
	}
}

func (s *EdgeCasesNSumTestSuite) TestNSum_SamePairs() {
	// Test pairs where from == to
	hashes := make(map[uint64]bool)

	for i := uint64(0); i < 100; i++ {
		hash := NSum(i, i)
		hashes[hash] = true
	}

	// All same pairs should produce unique hashes
	assert.Equal(s.T(), 100, len(hashes),
		"NSum(i, i) should produce unique hashes for different i")
}

func (s *EdgeCasesNSumTestSuite) TestNSum_LargeDifference() {
	// Test pairs with large differences
	testCases := []struct {
		from, to uint64
	}{
		{0, 1000000},
		{1, 1000001},
		{1000000, 0},
		{1000001, 1},
	}

	hashes := make(map[uint64]bool)
	for _, tc := range testCases {
		hash := NSum(tc.from, tc.to)
		hashes[hash] = true
	}

	assert.Equal(s.T(), len(testCases), len(hashes),
		"pairs with large differences should produce unique hashes")
}

func (s *EdgeCasesNSumTestSuite) TestNSum_PowersOfTwo() {
	// Test powers of two
	hashes := make(map[uint64]bool)

	for i := uint(0); i < 10; i++ {
		val := uint64(1) << i // 2^i
		hash := NSum(val, val)
		hashes[hash] = true
	}

	assert.Equal(s.T(), 10, len(hashes),
		"powers of two should produce unique hashes")
}

// Benchmark tests
func BenchmarkNSum(b *testing.B) {
	from := uint64(12345)
	to := uint64(67890)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NSum(from, to)
	}
}

func BenchmarkNSum_Sequential(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NSum(uint64(i), uint64(i+1))
	}
}

func BenchmarkNSum_DifferentPairs(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		from := uint64(i % 1000)
		to := uint64((i + 1) % 1000)
		NSum(from, to)
	}
}

func BenchmarkHashcode(b *testing.B) {
	base := uint64(12345)
	val := uint64(67890)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashcode(base, val)
	}
}

func BenchmarkNSum_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := uint64(0)
		for pb.Next() {
			_ = NSum(i, i+1)
			i++
		}
	})
}

// Use case example tests
func TestNSum_GraphEdgeHashing(t *testing.T) {
	// Example: hashing graph edges (node pairs)
	type edge struct {
		from, to uint64
	}

	edges := []edge{
		{1, 2}, {2, 3}, {3, 4}, {1, 3},
	}

	edgeHashes := make(map[uint64]edge)
	for _, e := range edges {
		hash := NSum(e.from, e.to)
		edgeHashes[hash] = e
	}

	// All edges should have unique hashes
	assert.Equal(t, len(edges), len(edgeHashes))

	// Verify we can look up edges by hash
	hash13 := NSum(1, 3)
	assert.Equal(t, edge{1, 3}, edgeHashes[hash13])
}

func TestNSum_CompositeKey(t *testing.T) {
	// Example: creating composite keys from two IDs
	userID := uint64(12345)
	sessionID := uint64(67890)

	// Create a consistent key for the user-session pair
	key1 := NSum(userID, sessionID)
	key2 := NSum(userID, sessionID)

	assert.Equal(t, key1, key2, "composite key should be consistent")

	// Different session should produce different key
	differentSessionKey := NSum(userID, sessionID+1)
	assert.NotEqual(t, key1, differentSessionKey)
}

func ExampleNSum() {
	// Hash two node IDs to create an edge identifier
	nodeA := uint64(123)
	nodeB := uint64(456)

	edgeHash := NSum(nodeA, nodeB)
	fmt.Printf("Edge hash: %d\n", edgeHash)

	// Same inputs always produce same output
	edgeHash2 := NSum(nodeA, nodeB)
	fmt.Printf("Consistent: %t\n", edgeHash == edgeHash2)

	// Order matters
	reverseHash := NSum(nodeB, nodeA)
	fmt.Printf("Order matters: %t\n", edgeHash != reverseHash)

	// Output:
	// Edge hash: 11400714819323206848
	// Consistent: true
	// Order matters: true
}

// Test suite runners
func TestNSumTestSuite(t *testing.T) {
	suite.Run(t, new(NSumTestSuite))
}

func TestHashcodeTestSuite(t *testing.T) {
	suite.Run(t, new(HashcodeTestSuite))
}

func TestEdgeCasesNSumTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesNSumTestSuite))
}
