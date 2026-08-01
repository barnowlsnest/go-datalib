// Package serial provides a thread-safe, sharded generator of monotonic
// serial identifiers.
//
// A Serial keeps 64 counters, each padded to its own CPU cache line so that
// concurrent increments of different keys do not contend through false
// sharing. Keys are mapped onto shards with FNV-1a hashing, so a given key
// always resolves to the same counter while distinct keys spread across the
// shard array. Counters advance with atomic operations; no locks are held on
// the hot path.
//
// The package offers:
//   - Seq to construct a generator, with Next to advance a key's counter and
//     Current to read it without advancing
//   - NSum to fold two uint64 values into a single hash, useful for composite
//     keys built from a pair of identifiers
//
// Serial is fully thread-safe: Next and Current may be called concurrently
// from any number of goroutines without external synchronization.
package serial
