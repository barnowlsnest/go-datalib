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

func TestLRUTestSuite(t *testing.T) {
	suite.Run(t, new(LRUTestSuite))
}
