package serial

// magic is the golden ratio constant used in hash calculations.
// This value (0x9e3779b97f4a7c15) is the 64-bit representation of
// the golden ratio multiplied by 2^64, providing good hash distribution
// properties and reducing clustering in hash-based data structures.
const magic uint64 = 0x9e3779b97f4a7c15 // Golden ratio (64-bit version)

// hashcode computes a hash value by combining two uint64 values.
//
// This function implements a fast hash algorithm that combines the base
// and val parameters using XOR operations and bit shifting. The algorithm
// incorporates the golden ratio constant to improve hash distribution
// and reduce collisions.
//
// The hash function uses the following operations:
//   - Adds the golden ratio constant to val
//   - Applies bit shifting to base (left shift by 6, right shift by 2)
//   - Combines all values using XOR for the final hash
//
// Parameters:
//   - base: The base value for hash calculation
//   - val: The value to be hashed with the base
//
// Returns:
//   - A 64-bit hash value combining both input parameters
//
// Performance:
// This function is designed for speed and uses only bitwise operations,
// making it suitable for high-frequency hash calculations.
func hashcode(base, val uint64) uint64 {
	return base ^ (val + magic + (base << 6) + (base >> 2))
}

// NSum computes a hash-based sum of two uint64 values.
//
// This function provides a consistent way to combine two numeric values
// into a single hash value. It's particularly useful for creating
// composite keys or generating consistent identifiers from pairs of values.
//
// The function uses the internal hashcode algorithm, which incorporates
// the golden ratio constant for good distribution properties.
//
// Parameters:
//   - from: The first value in the sum calculation
//   - to: The second value in the sum calculation
//
// Returns:
//   - A hash value representing the combination of both input values
//
// Use Cases:
//   - Creating composite keys from two numeric identifiers
//   - Generating consistent hash values for pairs of values
//   - Building hash-based data structures with compound keys
//
// Example:
//
//	// Create a hash from two node IDs
//	nodeA := uint64(123)
//	nodeB := uint64(456)
//	edgeHash := NSum(nodeA, nodeB)
//
//	// The same inputs always produce the same hash
//	assert.Equal(t, NSum(nodeA, nodeB), NSum(nodeA, nodeB))
func NSum(from, to uint64) uint64 {
	if from > to {
		from, to = to, from
	}
	return hashcode(from, to)
}
