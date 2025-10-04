package serial

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
)

const (
	// cacheLineBytes defines the typical CPU cache line size in bytes.
	// This constant is used for memory alignment to prevent false sharing
	// between different shards in concurrent access scenarios.
	cacheLineBytes = 64

	// atomicSharding defines the number of shards used for distributing
	// atomic operations. This value should be a power of 2 for efficient
	// bit masking in the hash function.
	atomicSharding = 64
)

type (
	// alignedShard represents a single shard in the serial ID generator.
	//
	// Each shard contains an atomic counter and padding to prevent false sharing.
	// False sharing occurs when multiple CPU cores access different variables
	// that reside on the same cache line, causing unnecessary cache invalidations
	// and performance degradation.
	//
	// The padding ensures that each shard occupies its own cache line,
	// allowing truly independent atomic operations across different shards.
	alignedShard struct {
		// id is the atomic counter for this shard.
		// It's incremented atomically to generate unique sequential IDs.
		id uint64

		// _ is padding to prevent false sharing between shards.
		// The padding size ensures each shard occupies exactly one cache line.
		_ [cacheLineBytes - 8]byte // should prevent false sharing
	}

	// Serial implements a high-performance, thread-safe serial ID generator.
	//
	// This implementation uses sharding to distribute atomic operations across
	// multiple cache-line-aligned counters, reducing contention and improving
	// performance in high-concurrency scenarios.
	//
	// Key features:
	//   - Thread-safe atomic operations
	//   - Sharding to reduce contention
	//   - Cache-line alignment to prevent false sharing
	//   - Key-based distribution for consistent shard selection
	//   - High-performance suitable for concurrent systems
	//
	// The Serial generator uses FNV-1a hashing to distribute keys across shards,
	// ensuring that the same key always maps to the same shard while providing
	// good distribution across different keys.
	//
	// Thread Safety:
	// Serial is fully thread-safe. Multiple goroutines can safely call
	// Next() and Current() concurrently without external synchronization.
	//
	// Performance:
	// The sharded design significantly reduces contention compared to a single
	// atomic counter, making it suitable for high-throughput applications.
	Serial struct {
		// shards is an array of cache-line-aligned atomic counters.
		// Each shard operates independently to reduce contention.
		shards [atomicSharding]alignedShard
	}
)

// hash computes a shard index for the given key using FNV-1a hashing.
//
// This function provides consistent and well-distributed mapping of keys
// to shard indices. The same key will always map to the same shard,
// ensuring consistent behavior while distributing load across shards.
//
// Parameters:
//   - key: The string key to hash for shard selection
//
// Returns:
//   - A shard index in the range [0, atomicSharding)
//
// Panics:
//   - If the hash function fails to write the key bytes (highly unlikely)
//
// Implementation:
// Uses FNV-1a hash algorithm for good distribution properties and applies
// bit masking for efficient modulo operation (requires atomicSharding to be power of 2).
func hash(key string) uint64 {
	hash := fnv.New64()
	_, err := hash.Write([]byte(key))
	if err != nil {
		panic(err)
	}
	return hash.Sum64() & (atomicSharding - 1)
}

// Next generates and returns the next sequential ID for the given key.
//
// This method atomically increments the counter for the shard associated
// with the key and returns the new value. Each key maintains its own
// independent sequence, allowing for predictable ID generation patterns.
//
// Parameters:
//   - key: The string key that determines which shard to use
//
// Returns:
//   - The next sequential ID for the given key (starting from 1)
//
// Thread Safety:
// This method is fully thread-safe and can be called concurrently
// from multiple goroutines without synchronization.
//
// Example:
//
//	serial := &Serial{}
//	id1 := serial.Next("user")     // Returns 1
//	id2 := serial.Next("user")     // Returns 2
//	id3 := serial.Next("product")  // Returns 1 (different key)
//	id4 := serial.Next("user")     // Returns 3
func (s *Serial) Next(key string) uint64 {
	return atomic.AddUint64(&s.shards[hash(key)].id, 1)
}

// Current returns the current ID value for the given key without incrementing.
//
// This method provides read-only access to the current counter value
// for the shard associated with the key. It's useful for checking
// the current state without generating a new ID.
//
// Parameters:
//   - key: The string key that determines which shard to read
//
// Returns:
//   - The current ID value for the given key (0 if never incremented)
//
// Thread Safety:
// This method is fully thread-safe and can be called concurrently
// from multiple goroutines without synchronization.
//
// Example:
//
//	serial := &Serial{}
//	current := serial.Current("user")  // Returns 0 (not yet incremented)
//	serial.Next("user")                // Returns 1
//	current = serial.Current("user")   // Returns 1
func (s *Serial) Current(key string) uint64 {
	return atomic.LoadUint64(&s.shards[hash(key)].id)
}

var (
	// ids is the singleton instance of the Serial generator.
	// It's initialized once using sync.Once for thread-safe singleton pattern.
	ids *Serial

	// once ensures that the Serial singleton is initialized exactly once,
	// even in concurrent scenarios.
	once sync.Once
)

// Seq returns the singleton instance of the Serial ID generator.
//
// This function implements the singleton pattern using sync.Once to ensure
// that only one Serial instance exists throughout the application lifecycle.
// The singleton approach is useful for maintaining global ID sequences
// across the entire application.
//
// Returns:
//   - The singleton Serial instance, ready for use
//
// Thread Safety:
// This function is fully thread-safe. Multiple goroutines can call
// Seq() concurrently and will always receive the same instance.
//
// Example:
//
//	// Get the global serial generator
//	serial := Seq()
//	userID := serial.Next("user")
//
//	// All calls to Seq() return the same instance
//	serial2 := Seq()
//	// serial == serial2 (same instance)
func Seq() *Serial {
	once.Do(func() {
		ids = &Serial{}
	})
	return ids
}
