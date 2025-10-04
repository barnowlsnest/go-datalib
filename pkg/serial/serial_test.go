package serial

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AlignmentTestSuite verifies memory alignment and struct size
type AlignmentTestSuite struct {
	suite.Suite
}

func (s *AlignmentTestSuite) TestAlignedShardSize() {
	// Verify that alignedShard is exactly one cache line (64 bytes)
	size := unsafe.Sizeof(alignedShard{})
	assert.Equal(s.T(), uintptr(cacheLineBytes), size,
		"alignedShard should be exactly %d bytes to occupy one cache line", cacheLineBytes)
}

func (s *AlignmentTestSuite) TestAlignedShardAlignment() {
	// Verify proper alignment of the id field
	var shard alignedShard
	idOffset := unsafe.Offsetof(shard.id)
	assert.Equal(s.T(), uintptr(0), idOffset,
		"id field should be at offset 0 for proper alignment")
}

func (s *AlignmentTestSuite) TestSerialShardArrayAlignment() {
	// Verify that shards in the array are properly separated
	serial := &Serial{}

	// Check distance between consecutive shards
	shard0Addr := uintptr(unsafe.Pointer(&serial.shards[0]))
	shard1Addr := uintptr(unsafe.Pointer(&serial.shards[1]))

	distance := shard1Addr - shard0Addr
	assert.Equal(s.T(), uintptr(cacheLineBytes), distance,
		"consecutive shards should be %d bytes apart", cacheLineBytes)
}

func (s *AlignmentTestSuite) TestSerialStructSize() {
	// Verify total size of Serial struct
	size := unsafe.Sizeof(Serial{})
	expectedSize := uintptr(atomicSharding * cacheLineBytes)
	assert.Equal(s.T(), expectedSize, size,
		"Serial struct should be %d bytes (%d shards * %d bytes)",
		expectedSize, atomicSharding, cacheLineBytes)
}

// BasicFunctionalityTestSuite tests core functionality
type BasicFunctionalityTestSuite struct {
	suite.Suite
}

func (s *BasicFunctionalityTestSuite) TestNext_FirstCall() {
	serial := &Serial{}

	id := serial.Next("test")
	assert.Equal(s.T(), uint64(1), id, "first Next() should return 1")
}

func (s *BasicFunctionalityTestSuite) TestNext_Sequential() {
	serial := &Serial{}

	id1 := serial.Next("test")
	id2 := serial.Next("test")
	id3 := serial.Next("test")

	assert.Equal(s.T(), uint64(1), id1)
	assert.Equal(s.T(), uint64(2), id2)
	assert.Equal(s.T(), uint64(3), id3)
}

func (s *BasicFunctionalityTestSuite) TestNext_DifferentKeys() {
	serial := &Serial{}

	id1 := serial.Next("user")
	id2 := serial.Next("product")
	id3 := serial.Next("order")

	// Different keys may or may not share shards, but all should start from 1
	// if they're using different shards
	assert.Greater(s.T(), id1, uint64(0))
	assert.Greater(s.T(), id2, uint64(0))
	assert.Greater(s.T(), id3, uint64(0))
}

func (s *BasicFunctionalityTestSuite) TestNext_SameKeySameShard() {
	serial := &Serial{}

	key := "test-key"
	id1 := serial.Next(key)
	id2 := serial.Next(key)

	// Same key should always use same shard
	assert.Equal(s.T(), id1+1, id2, "same key should increment sequentially")
}

func (s *BasicFunctionalityTestSuite) TestCurrent_InitialValue() {
	serial := &Serial{}

	current := serial.Current("test")
	assert.Equal(s.T(), uint64(0), current, "initial Current() should return 0")
}

func (s *BasicFunctionalityTestSuite) TestCurrent_AfterNext() {
	serial := &Serial{}

	serial.Next("test")
	serial.Next("test")
	current := serial.Current("test")

	assert.Equal(s.T(), uint64(2), current, "Current() should return last Next() value")
}

func (s *BasicFunctionalityTestSuite) TestCurrent_DoesNotIncrement() {
	serial := &Serial{}

	serial.Next("test") // id = 1
	current1 := serial.Current("test")
	current2 := serial.Current("test")
	current3 := serial.Current("test")

	assert.Equal(s.T(), uint64(1), current1)
	assert.Equal(s.T(), uint64(1), current2)
	assert.Equal(s.T(), uint64(1), current3)
}

// HashingTestSuite tests the hash function
type HashingTestSuite struct {
	suite.Suite
}

func (s *HashingTestSuite) TestHash_Consistency() {
	key := "test-key"
	hash1 := hash(key)
	hash2 := hash(key)
	hash3 := hash(key)

	assert.Equal(s.T(), hash1, hash2, "same key should produce same hash")
	assert.Equal(s.T(), hash1, hash3, "same key should produce same hash")
}

func (s *HashingTestSuite) TestHash_Range() {
	keys := []string{"user", "product", "order", "invoice", "payment"}

	for _, key := range keys {
		h := hash(key)
		assert.Less(s.T(), h, uint64(atomicSharding),
			"hash(%s) should be less than %d", key, atomicSharding)
	}
}

func (s *HashingTestSuite) TestHash_Distribution() {
	// Test that different keys produce different hashes (mostly)
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	hashes := make(map[uint64]bool)

	for _, key := range keys {
		h := hash(key)
		hashes[h] = true
	}

	// We expect good distribution, so most should be unique
	// (allowing for some collisions due to limited shard count)
	assert.GreaterOrEqual(s.T(), len(hashes), 5,
		"hash should distribute keys reasonably well")
}

// ConcurrencyTestSuite tests thread safety
type ConcurrencyTestSuite struct {
	suite.Suite
}

func (s *ConcurrencyTestSuite) TestNext_Concurrent_SameKey() {
	serial := &Serial{}
	key := "test"
	iterations := 1000
	goroutines := 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				serial.Next(key)
			}
		}()
	}

	wg.Wait()

	final := serial.Current(key)
	expected := uint64(goroutines * iterations)
	assert.Equal(s.T(), expected, final,
		"concurrent Next() calls should produce correct total count")
}

func (s *ConcurrencyTestSuite) TestNext_Concurrent_DifferentKeys() {
	serial := &Serial{}
	iterations := 1000
	goroutines := 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		key := fmt.Sprintf("key-%d", i)
		go func(k string) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				serial.Next(k)
			}
		}(key)
	}

	wg.Wait()

	// Each key should have exactly iterations count
	for i := 0; i < goroutines; i++ {
		key := fmt.Sprintf("key-%d", i)
		count := serial.Current(key)
		assert.Equal(s.T(), uint64(iterations), count,
			"key %s should have count %d", key, iterations)
	}
}

func (s *ConcurrencyTestSuite) TestNext_HighContention() {
	serial := &Serial{}
	key := "high-contention"
	iterations := 10000
	goroutines := 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				serial.Next(key)
			}
		}()
	}

	wg.Wait()

	final := serial.Current(key)
	expected := uint64(goroutines * iterations)
	assert.Equal(s.T(), expected, final,
		"high contention should still produce correct count")
}

func (s *ConcurrencyTestSuite) TestCurrent_ConcurrentReads() {
	serial := &Serial{}
	key := "read-test"

	// Set initial value
	serial.Next(key)
	serial.Next(key)
	serial.Next(key)

	iterations := 1000
	goroutines := 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Multiple concurrent reads
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				current := serial.Current(key)
				assert.Equal(s.T(), uint64(3), current)
			}
		}()
	}

	wg.Wait()
}

func (s *ConcurrencyTestSuite) TestMixedReadWrite() {
	serial := &Serial{}
	key := "mixed"
	iterations := 1000
	readers := 5
	writers := 5

	var wg sync.WaitGroup
	wg.Add(readers + writers)

	// Start writers
	for i := 0; i < writers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				serial.Next(key)
			}
		}()
	}

	// Start readers
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				serial.Current(key)
			}
		}()
	}

	wg.Wait()

	final := serial.Current(key)
	expected := uint64(writers * iterations)
	assert.Equal(s.T(), expected, final)
}

// SingletonTestSuite tests the Seq singleton
type SingletonTestSuite struct {
	suite.Suite
}

func (s *SingletonTestSuite) TestSeq_ReturnsSameInstance() {
	// Reset singleton for testing
	once = sync.Once{}
	ids = nil

	serial1 := Seq()
	serial2 := Seq()
	serial3 := Seq()

	// All should be the same instance
	assert.Same(s.T(), serial1, serial2)
	assert.Same(s.T(), serial1, serial3)
}

func (s *SingletonTestSuite) TestSeq_ConcurrentAccess() {
	// Reset singleton for testing
	once = sync.Once{}
	ids = nil

	goroutines := 100
	instances := make([]*Serial, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			instances[idx] = Seq()
		}(i)
	}

	wg.Wait()

	// All should be the same instance
	first := instances[0]
	for i := 1; i < goroutines; i++ {
		assert.Same(s.T(), first, instances[i],
			"all Seq() calls should return same instance")
	}
}

func (s *SingletonTestSuite) TestSeq_SharedState() {
	// Reset singleton for testing
	once = sync.Once{}
	ids = nil

	serial1 := Seq()
	serial1.Next("test")
	serial1.Next("test")

	serial2 := Seq()
	current := serial2.Current("test")

	assert.Equal(s.T(), uint64(2), current,
		"singleton should share state across calls")
}

// EdgeCasesTestSuite tests edge cases
type EdgeCasesTestSuite struct {
	suite.Suite
}

func (s *EdgeCasesTestSuite) TestEmptyKey() {
	serial := &Serial{}

	id1 := serial.Next("")
	id2 := serial.Next("")

	assert.Equal(s.T(), uint64(1), id1)
	assert.Equal(s.T(), uint64(2), id2)
}

func (s *EdgeCasesTestSuite) TestVeryLongKey() {
	serial := &Serial{}
	longKey := string(make([]byte, 10000))

	id := serial.Next(longKey)
	assert.Equal(s.T(), uint64(1), id)
}

func (s *EdgeCasesTestSuite) TestUnicodeKeys() {
	serial := &Serial{}

	id1 := serial.Next("用户")
	id2 := serial.Next("用户")
	id3 := serial.Next("продукт")

	assert.Equal(s.T(), uint64(1), id1)
	assert.Equal(s.T(), uint64(2), id2)
	assert.Equal(s.T(), uint64(1), id3)
}

func (s *EdgeCasesTestSuite) TestManyKeys() {
	serial := &Serial{}
	numKeys := 1000

	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		serial.Next(key)
	}

	// Verify all keys have their counts
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		count := serial.Current(key)
		assert.Greater(s.T(), count, uint64(0))
	}
}

// Benchmark tests
func BenchmarkSerial_Next_SameKey(b *testing.B) {
	serial := &Serial{}
	key := "benchmark"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serial.Next(key)
	}
}

func BenchmarkSerial_Next_DifferentKeys(b *testing.B) {
	serial := &Serial{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%100)
		serial.Next(key)
	}
}

func BenchmarkSerial_Current(b *testing.B) {
	serial := &Serial{}
	key := "benchmark"
	serial.Next(key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serial.Current(key)
	}
}

func BenchmarkSerial_Parallel_Next_SameKey(b *testing.B) {
	serial := &Serial{}
	key := "parallel"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			serial.Next(key)
		}
	})
}

func BenchmarkSerial_Parallel_Next_DifferentKeys(b *testing.B) {
	serial := &Serial{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var localCounter uint64
		for pb.Next() {
			key := fmt.Sprintf("key-%d", atomic.AddUint64(&localCounter, 1)%100)
			serial.Next(key)
		}
	})
}

func BenchmarkSerial_Parallel_Current(b *testing.B) {
	serial := &Serial{}
	key := "parallel-read"
	serial.Next(key)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			serial.Current(key)
		}
	})
}

func BenchmarkSerial_Parallel_Mixed(b *testing.B) {
	serial := &Serial{}
	key := "mixed"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				serial.Current(key)
			} else {
				serial.Next(key)
			}
			i++
		}
	})
}

func BenchmarkHash(b *testing.B) {
	key := "benchmark-key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash(key)
	}
}

// Comparison benchmark: single atomic counter (no sharding)
func BenchmarkAtomic_NoSharding_Parallel(b *testing.B) {
	var counter uint64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddUint64(&counter, 1)
		}
	})
}

// Test suite runners
func TestAlignmentTestSuite(t *testing.T) {
	suite.Run(t, new(AlignmentTestSuite))
}

func TestBasicFunctionalityTestSuite(t *testing.T) {
	suite.Run(t, new(BasicFunctionalityTestSuite))
}

func TestHashingTestSuite(t *testing.T) {
	suite.Run(t, new(HashingTestSuite))
}

func TestConcurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrencyTestSuite))
}

func TestSingletonTestSuite(t *testing.T) {
	suite.Run(t, new(SingletonTestSuite))
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}
