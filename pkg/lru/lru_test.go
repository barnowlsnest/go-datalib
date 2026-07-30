package lru

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// LRUTestSuite defines tests for the LRU cache.
type LRUTestSuite struct {
	suite.Suite
}

func (s *LRUTestSuite) TestNew_DefaultsCapacityWhenNonPositive() {
	s.Require().Equal(defaultLRUCapacity, New[int](0).capacity)
	s.Require().Equal(defaultLRUCapacity, New[int](-5).capacity)
	s.Require().Equal(8, New[int](8).capacity)
}

func (s *LRUTestSuite) TestGet_MissReturnsError() {
	cache := New[int](4)

	value, err := cache.Get(42)

	s.Require().Nil(value)
	s.Require().ErrorIs(err, ErrCacheMiss)
}

func (s *LRUTestSuite) TestPutThenGet_ReturnsValue() {
	cache := New[int](4)

	cache.Put(1, new(7))
	got, err := cache.Get(1)

	s.Require().NoError(err)
	s.Require().Equal(7, *got)
	s.Require().Equal(1, cache.Len())
}

func (s *LRUTestSuite) TestPut_ExistingKeyUpdatesValue() {
	cache := New[int](4)

	cache.Put(1, new(1))
	cache.Put(1, new(2))
	got, err := cache.Get(1)

	s.Require().NoError(err)
	s.Require().Equal(2, *got)
	s.Require().Equal(1, cache.Len())
}

func (s *LRUTestSuite) TestPut_EvictsLeastRecentlyUsed() {
	cache := New[int](2)

	cache.Put(1, new(1))
	cache.Put(2, new(2))
	cache.Put(3, new(3)) // exceeds capacity -> evicts key 1 (LRU)

	s.Require().Equal(2, cache.Len())

	_, err := cache.Get(1)
	s.Require().ErrorIs(err, ErrCacheMiss)

	got2, err := cache.Get(2)
	s.Require().NoError(err)
	s.Require().Equal(2, *got2)

	got3, err := cache.Get(3)
	s.Require().NoError(err)
	s.Require().Equal(3, *got3)
}

func (s *LRUTestSuite) TestGet_RefreshesRecency() {
	cache := New[int](2)

	cache.Put(1, new(1))
	cache.Put(2, new(2))

	// Access key 1, making key 2 the least recently used.
	_, err := cache.Get(1)
	s.Require().NoError(err)

	cache.Put(3, new(3)) // evicts key 2, not key 1

	_, err = cache.Get(2)
	s.Require().ErrorIs(err, ErrCacheMiss)

	got1, err := cache.Get(1)
	s.Require().NoError(err)
	s.Require().Equal(1, *got1)
}

func (s *LRUTestSuite) TestPut_ExistingKeyDoesNotEvict() {
	cache := New[int](2)

	cache.Put(1, new(1))
	cache.Put(2, new(2))
	cache.Put(1, new(11)) // update existing -> no eviction

	s.Require().Equal(2, cache.Len())

	got1, err := cache.Get(1)
	s.Require().NoError(err)
	s.Require().Equal(11, *got1)

	got2, err := cache.Get(2)
	s.Require().NoError(err)
	s.Require().Equal(2, *got2)
}

func (s *LRUTestSuite) TestDelete_RemovesExistingKey() {
	cache := New[int](4)
	cache.Put(1, new(1))
	cache.Put(2, new(2))

	cache.Delete(1)

	s.Require().Equal(1, cache.Len())

	_, err := cache.Get(1)
	s.Require().ErrorIs(err, ErrCacheMiss)

	got2, err := cache.Get(2)
	s.Require().NoError(err)
	s.Require().Equal(2, *got2)
}

func (s *LRUTestSuite) TestDelete_MissingKeyIsNoOp() {
	cache := New[int](4)
	cache.Put(1, new(1))

	cache.Delete(99)

	s.Require().Equal(1, cache.Len())

	got1, err := cache.Get(1)
	s.Require().NoError(err)
	s.Require().Equal(1, *got1)
}

func (s *LRUTestSuite) TestDelete_NoKeysIsNoOp() {
	cache := New[int](4)
	cache.Put(1, new(1))

	cache.Delete()

	s.Require().Equal(1, cache.Len())

	got1, err := cache.Get(1)
	s.Require().NoError(err)
	s.Require().Equal(1, *got1)
}

func (s *LRUTestSuite) TestDelete_MultipleKeys() {
	cache := New[int](4)
	cache.Put(1, new(1))
	cache.Put(2, new(2))
	cache.Put(3, new(3))

	// Mixes present (1, 3) and absent (99) keys.
	cache.Delete(1, 99, 3)

	s.Require().Equal(1, cache.Len())

	_, err := cache.Get(1)
	s.Require().ErrorIs(err, ErrCacheMiss)

	_, err = cache.Get(3)
	s.Require().ErrorIs(err, ErrCacheMiss)

	got2, err := cache.Get(2)
	s.Require().NoError(err)
	s.Require().Equal(2, *got2)
}

func (s *LRUTestSuite) TestDelete_OnlyEntry() {
	cache := New[int](4)
	cache.Put(1, new(1))

	cache.Delete(1)

	s.Require().Equal(0, cache.Len())

	_, err := cache.Get(1)
	s.Require().ErrorIs(err, ErrCacheMiss)
}

func (s *LRUTestSuite) TestDelete_LeastRecentlyUsedEntry() {
	cache := New[int](3)
	cache.Put(1, new(1)) // becomes LRU (tail)
	cache.Put(2, new(2))
	cache.Put(3, new(3))

	// Delete the current LRU entry (the list tail).
	cache.Delete(1)

	s.Require().Equal(2, cache.Len())

	_, err := cache.Get(1)
	s.Require().ErrorIs(err, ErrCacheMiss)

	got2, err := cache.Get(2)
	s.Require().NoError(err)
	s.Require().Equal(2, *got2)

	got3, err := cache.Get(3)
	s.Require().NoError(err)
	s.Require().Equal(3, *got3)
}

func (s *LRUTestSuite) TestDelete_PreservesRecencyOfSurvivors() {
	cache := New[int](2)
	cache.Put(1, new(1))
	cache.Put(2, new(2)) // order MRU->LRU: 2, 1

	// Delete the MRU entry; key 1 must remain the only (and LRU) entry.
	cache.Delete(2)
	s.Require().Equal(1, cache.Len())

	// Insert two more so the cache evicts. Key 1 is LRU and should go first.
	cache.Put(3, new(3)) // order: 3, 1
	cache.Put(4, new(4)) // exceeds capacity -> evicts key 1

	_, err := cache.Get(1)
	s.Require().ErrorIs(err, ErrCacheMiss)

	got3, err := cache.Get(3)
	s.Require().NoError(err)
	s.Require().Equal(3, *got3)

	got4, err := cache.Get(4)
	s.Require().NoError(err)
	s.Require().Equal(4, *got4)
}

func (s *LRUTestSuite) TestDelete_FreesCapacityForReinsertion() {
	cache := New[int](2)
	cache.Put(1, new(1))
	cache.Put(2, new(2))

	cache.Delete(1)

	// Re-inserting must not evict key 2, since deletion freed a slot.
	cache.Put(3, new(3))

	s.Require().Equal(2, cache.Len())

	got2, err := cache.Get(2)
	s.Require().NoError(err)
	s.Require().Equal(2, *got2)

	got3, err := cache.Get(3)
	s.Require().NoError(err)
	s.Require().Equal(3, *got3)
}

func (s *LRUTestSuite) TestClear_RemovesAllEntries() {
	cache := New[int](4)
	cache.Put(1, new(1))
	cache.Put(2, new(2))
	cache.Put(3, new(3))

	cache.Clear()

	s.Require().Equal(0, cache.Len())

	for _, key := range []uint64{1, 2, 3} {
		_, err := cache.Get(key)
		s.Require().ErrorIs(err, ErrCacheMiss)
	}
}

func (s *LRUTestSuite) TestClear_EmptyCacheIsNoOp() {
	cache := New[int](4)

	cache.Clear()

	s.Require().Equal(0, cache.Len())
}

func (s *LRUTestSuite) TestClear_PreservesCapacityForReuse() {
	cache := New[int](2)
	cache.Put(1, new(1))
	cache.Put(2, new(2))

	cache.Clear()

	// The cache is reusable at full capacity and its recency list is reset:
	// filling it exactly must not evict anything.
	cache.Put(3, new(3))
	cache.Put(4, new(4))

	s.Require().Equal(2, cache.Len())

	got3, err := cache.Get(3)
	s.Require().NoError(err)
	s.Require().Equal(3, *got3)

	got4, err := cache.Get(4)
	s.Require().NoError(err)
	s.Require().Equal(4, *got4)
}

func TestLRUTestSuite(t *testing.T) {
	suite.Run(t, new(LRUTestSuite))
}
