package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BTreeTestSuite struct {
	suite.Suite
}

func TestBTreeTestSuite(t *testing.T) {
	suite.Run(t, new(BTreeTestSuite))
}

// ============================================================================
// Constructor Tests
// ============================================================================

func (s *BTreeTestSuite) TestNewBTree_DefaultMinDegree() {
	tree := NewBTree[int, string](0)

	s.Equal(DefaultMinDegree, tree.MinDegree())
	s.Equal(0, tree.Size())
	s.True(tree.IsEmpty())
}

func (s *BTreeTestSuite) TestNewBTree_CustomMinDegree() {
	tree := NewBTree[int, string](5)

	s.Equal(5, tree.MinDegree())
}

func (s *BTreeTestSuite) TestNewBTree_MinDegreeOneBecomesDefault() {
	tree := NewBTree[int, string](1)

	s.Equal(DefaultMinDegree, tree.MinDegree())
}

// ============================================================================
// Insert and Search Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Insert_Single() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")

	s.Equal(1, tree.Size())
	val, found := tree.Search(1)
	s.True(found)
	s.Equal("one", val)
}

func (s *BTreeTestSuite) TestBTree_Insert_Multiple_Ascending() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	s.Equal(10, tree.Size())
	for i := 1; i <= 10; i++ {
		s.True(tree.Contains(i))
	}
}

func (s *BTreeTestSuite) TestBTree_Insert_Multiple_Descending() {
	tree := NewBTree[int, string](2)

	for i := 10; i >= 1; i-- {
		tree.Insert(i, "value")
	}

	s.Equal(10, tree.Size())
	for i := 1; i <= 10; i++ {
		s.True(tree.Contains(i))
	}
}

func (s *BTreeTestSuite) TestBTree_Insert_UpdateExisting() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "original")
	tree.Insert(1, "updated")

	s.Equal(1, tree.Size())
	val, found := tree.Search(1)
	s.True(found)
	s.Equal("updated", val)
}

func (s *BTreeTestSuite) TestBTree_Insert_CausesSplit() {
	tree := NewBTree[int, string](2) // max 3 keys per node

	// Insert 4 keys to force a split
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")
	tree.Insert(4, "four")

	s.Equal(4, tree.Size())
	for i := 1; i <= 4; i++ {
		s.True(tree.Contains(i))
	}
}

func (s *BTreeTestSuite) TestBTree_Insert_LargeDataset() {
	tree := NewBTree[int, string](3)

	for i := 1; i <= 1000; i++ {
		tree.Insert(i, "value")
	}

	s.Equal(1000, tree.Size())
	for i := 1; i <= 1000; i++ {
		s.True(tree.Contains(i), "Key %d should exist", i)
	}
}

func (s *BTreeTestSuite) TestBTree_Search_NotFound() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(3, "three")

	val, found := tree.Search(2)
	s.False(found)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Search_EmptyTree() {
	tree := NewBTree[int, string](2)

	val, found := tree.Search(1)
	s.False(found)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Contains() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(2, "two")

	s.True(tree.Contains(1))
	s.True(tree.Contains(2))
	s.False(tree.Contains(3))
}

// ============================================================================
// Delete Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Delete_Leaf() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")

	deleted := tree.Delete(3)

	s.True(deleted)
	s.Equal(2, tree.Size())
	s.False(tree.Contains(3))
	s.True(tree.Contains(1))
	s.True(tree.Contains(2))
}

func (s *BTreeTestSuite) TestBTree_Delete_NotFound() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")

	deleted := tree.Delete(999)

	s.False(deleted)
	s.Equal(1, tree.Size())
}

func (s *BTreeTestSuite) TestBTree_Delete_EmptyTree() {
	tree := NewBTree[int, string](2)

	deleted := tree.Delete(1)

	s.False(deleted)
	s.Equal(0, tree.Size())
}

func (s *BTreeTestSuite) TestBTree_Delete_SingleElement() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	deleted := tree.Delete(1)

	s.True(deleted)
	s.True(tree.IsEmpty())
	s.Nil(tree.root)
}

func (s *BTreeTestSuite) TestBTree_Delete_AllElements() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	for i := 1; i <= 10; i++ {
		s.True(tree.Delete(i))
	}

	s.True(tree.IsEmpty())
}

func (s *BTreeTestSuite) TestBTree_Delete_FromInternal() {
	tree := NewBTree[int, string](2)

	// Build tree that has internal nodes
	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	// Delete from middle (likely internal)
	s.True(tree.Delete(5))
	s.False(tree.Contains(5))
	s.Equal(9, tree.Size())

	// Verify other keys still exist
	for i := 1; i <= 10; i++ {
		if i != 5 {
			s.True(tree.Contains(i), "Key %d should still exist", i)
		}
	}
}

func (s *BTreeTestSuite) TestBTree_Delete_CausesMerge() {
	tree := NewBTree[int, string](2)

	// Build and then delete to cause merges
	for i := 1; i <= 7; i++ {
		tree.Insert(i, "value")
	}

	// Delete several keys to force merges
	tree.Delete(1)
	tree.Delete(2)
	tree.Delete(3)

	s.Equal(4, tree.Size())
	for i := 4; i <= 7; i++ {
		s.True(tree.Contains(i))
	}
}

func (s *BTreeTestSuite) TestBTree_Delete_LargeDataset() {
	tree := NewBTree[int, string](3)

	for i := 1; i <= 100; i++ {
		tree.Insert(i, "value")
	}

	// Delete even numbers
	for i := 2; i <= 100; i += 2 {
		s.True(tree.Delete(i))
	}

	s.Equal(50, tree.Size())

	// Verify odd numbers still exist
	for i := 1; i <= 100; i += 2 {
		s.True(tree.Contains(i))
	}

	// Verify even numbers are gone
	for i := 2; i <= 100; i += 2 {
		s.False(tree.Contains(i))
	}
}

// ============================================================================
// Min/Max Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Min_Empty() {
	tree := NewBTree[int, string](2)

	key, val, found := tree.Min()
	s.False(found)
	s.Equal(0, key)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Min_Single() {
	tree := NewBTree[int, string](2)

	tree.Insert(5, "five")

	key, val, found := tree.Min()
	s.True(found)
	s.Equal(5, key)
	s.Equal("five", val)
}

func (s *BTreeTestSuite) TestBTree_Min_Multiple() {
	tree := NewBTree[int, string](2)

	tree.Insert(5, "five")
	tree.Insert(1, "one")
	tree.Insert(10, "ten")

	key, val, found := tree.Min()
	s.True(found)
	s.Equal(1, key)
	s.Equal("one", val)
}

func (s *BTreeTestSuite) TestBTree_Max_Empty() {
	tree := NewBTree[int, string](2)

	key, val, found := tree.Max()
	s.False(found)
	s.Equal(0, key)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Max_Single() {
	tree := NewBTree[int, string](2)

	tree.Insert(5, "five")

	key, val, found := tree.Max()
	s.True(found)
	s.Equal(5, key)
	s.Equal("five", val)
}

func (s *BTreeTestSuite) TestBTree_Max_Multiple() {
	tree := NewBTree[int, string](2)

	tree.Insert(5, "five")
	tree.Insert(1, "one")
	tree.Insert(10, "ten")

	key, val, found := tree.Max()
	s.True(found)
	s.Equal(10, key)
	s.Equal("ten", val)
}

// ============================================================================
// Floor/Ceiling Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Floor_Empty() {
	tree := NewBTree[int, string](2)

	key, val, found := tree.Floor(5)
	s.False(found)
	s.Equal(0, key)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Floor_ExactMatch() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(5, "five")
	tree.Insert(10, "ten")

	key, val, found := tree.Floor(5)
	s.True(found)
	s.Equal(5, key)
	s.Equal("five", val)
}

func (s *BTreeTestSuite) TestBTree_Floor_LessThan() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(5, "five")
	tree.Insert(10, "ten")

	key, val, found := tree.Floor(7)
	s.True(found)
	s.Equal(5, key)
	s.Equal("five", val)
}

func (s *BTreeTestSuite) TestBTree_Floor_NoFloor() {
	tree := NewBTree[int, string](2)

	tree.Insert(5, "five")
	tree.Insert(10, "ten")

	key, val, found := tree.Floor(3)
	s.False(found)
	s.Equal(0, key)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Ceiling_Empty() {
	tree := NewBTree[int, string](2)

	key, val, found := tree.Ceiling(5)
	s.False(found)
	s.Equal(0, key)
	s.Equal("", val)
}

func (s *BTreeTestSuite) TestBTree_Ceiling_ExactMatch() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(5, "five")
	tree.Insert(10, "ten")

	key, val, found := tree.Ceiling(5)
	s.True(found)
	s.Equal(5, key)
	s.Equal("five", val)
}

func (s *BTreeTestSuite) TestBTree_Ceiling_GreaterThan() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(5, "five")
	tree.Insert(10, "ten")

	key, val, found := tree.Ceiling(3)
	s.True(found)
	s.Equal(5, key)
	s.Equal("five", val)
}

func (s *BTreeTestSuite) TestBTree_Ceiling_NoCeiling() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(5, "five")

	key, val, found := tree.Ceiling(10)
	s.False(found)
	s.Equal(0, key)
	s.Equal("", val)
}

// ============================================================================
// Range Query Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Range_Empty() {
	tree := NewBTree[int, string](2)

	var results []BTreeEntry[int, string]
	for entry := range tree.Range(1, 10) {
		results = append(results, entry)
	}

	s.Empty(results)
}

func (s *BTreeTestSuite) TestBTree_Range_Full() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	var keys []int
	for entry := range tree.Range(1, 10) {
		keys = append(keys, entry.Key)
	}

	s.Equal([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, keys)
}

func (s *BTreeTestSuite) TestBTree_Range_Partial() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	var keys []int
	for entry := range tree.Range(3, 7) {
		keys = append(keys, entry.Key)
	}

	s.Equal([]int{3, 4, 5, 6, 7}, keys)
}

func (s *BTreeTestSuite) TestBTree_Range_NoMatch() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")
	tree.Insert(10, "ten")

	var results []BTreeEntry[int, string]
	for entry := range tree.Range(5, 7) {
		results = append(results, entry)
	}

	s.Empty(results)
}

func (s *BTreeTestSuite) TestBTree_Range_InvalidBounds() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	var results []BTreeEntry[int, string]
	for entry := range tree.Range(10, 1) { // from > to
		results = append(results, entry)
	}

	s.Empty(results)
}

func (s *BTreeTestSuite) TestBTree_Range_EarlyBreak() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	var keys []int
	for entry := range tree.Range(1, 10) {
		keys = append(keys, entry.Key)
		if entry.Key == 5 {
			break
		}
	}

	s.Equal([]int{1, 2, 3, 4, 5}, keys)
}

// ============================================================================
// All Iterator Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_All_Empty() {
	tree := NewBTree[int, string](2)

	var results []BTreeEntry[int, string]
	for entry := range tree.All() {
		results = append(results, entry)
	}

	s.Empty(results)
}

func (s *BTreeTestSuite) TestBTree_All_InOrder() {
	tree := NewBTree[int, string](2)

	// Insert in random order
	tree.Insert(5, "five")
	tree.Insert(1, "one")
	tree.Insert(3, "three")
	tree.Insert(7, "seven")
	tree.Insert(2, "two")

	var keys []int
	for entry := range tree.All() {
		keys = append(keys, entry.Key)
	}

	s.Equal([]int{1, 2, 3, 5, 7}, keys)
}

func (s *BTreeTestSuite) TestBTree_All_EarlyBreak() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	count := 0
	for range tree.All() {
		count++
		if count == 3 {
			break
		}
	}

	s.Equal(3, count)
}

// ============================================================================
// Keys/Values Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Keys_Empty() {
	tree := NewBTree[int, string](2)

	keys := tree.Keys()
	s.Empty(keys)
}

func (s *BTreeTestSuite) TestBTree_Keys_InOrder() {
	tree := NewBTree[int, string](2)

	tree.Insert(5, "five")
	tree.Insert(1, "one")
	tree.Insert(3, "three")

	keys := tree.Keys()
	s.Equal([]int{1, 3, 5}, keys)
}

func (s *BTreeTestSuite) TestBTree_Values_Empty() {
	tree := NewBTree[int, string](2)

	values := tree.Values()
	s.Empty(values)
}

func (s *BTreeTestSuite) TestBTree_Values_InKeyOrder() {
	tree := NewBTree[int, string](2)

	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(2, "two")

	values := tree.Values()
	s.Equal([]string{"one", "two", "three"}, values)
}

// ============================================================================
// Height Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Height_Empty() {
	tree := NewBTree[int, string](2)

	s.Equal(0, tree.Height())
}

func (s *BTreeTestSuite) TestBTree_Height_Single() {
	tree := NewBTree[int, string](2)

	tree.Insert(1, "one")

	s.Equal(1, tree.Height())
}

func (s *BTreeTestSuite) TestBTree_Height_AfterSplit() {
	tree := NewBTree[int, string](2) // max 3 keys per node

	// Insert enough to cause splits
	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	s.Greater(tree.Height(), 1)
}

// ============================================================================
// Clear Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_Clear() {
	tree := NewBTree[int, string](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, "value")
	}

	tree.Clear()

	s.True(tree.IsEmpty())
	s.Equal(0, tree.Size())
	s.Equal(0, tree.Height())
}

// ============================================================================
// Type Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_StringKeys() {
	tree := NewBTree[string, int](2)

	tree.Insert("apple", 1)
	tree.Insert("banana", 2)
	tree.Insert("cherry", 3)

	val, found := tree.Search("banana")
	s.True(found)
	s.Equal(2, val)

	keys := tree.Keys()
	s.Equal([]string{"apple", "banana", "cherry"}, keys)
}

func (s *BTreeTestSuite) TestBTree_Uint64Keys() {
	tree := NewBTree[uint64, string](2)

	tree.Insert(100, "hundred")
	tree.Insert(50, "fifty")
	tree.Insert(200, "two hundred")

	val, found := tree.Search(100)
	s.True(found)
	s.Equal("hundred", val)

	minKey, _, _ := tree.Min()
	s.Equal(uint64(50), minKey)
}

func (s *BTreeTestSuite) TestBTree_Float64Keys() {
	tree := NewBTree[float64, string](2)

	tree.Insert(1.5, "one point five")
	tree.Insert(0.5, "half")
	tree.Insert(2.5, "two point five")

	val, found := tree.Search(1.5)
	s.True(found)
	s.Equal("one point five", val)

	keys := tree.Keys()
	s.Equal([]float64{0.5, 1.5, 2.5}, keys)
}

// ============================================================================
// Stress Tests
// ============================================================================

func (s *BTreeTestSuite) TestBTree_MixedOperations() {
	tree := NewBTree[int, string](3)

	// Insert
	for i := 1; i <= 50; i++ {
		tree.Insert(i, "value")
	}
	s.Equal(50, tree.Size())

	// Delete some
	for i := 10; i <= 20; i++ {
		tree.Delete(i)
	}
	s.Equal(39, tree.Size())

	// Insert more
	for i := 100; i <= 110; i++ {
		tree.Insert(i, "value")
	}
	s.Equal(50, tree.Size())

	// Verify structure
	for i := 1; i <= 9; i++ {
		s.True(tree.Contains(i))
	}
	for i := 10; i <= 20; i++ {
		s.False(tree.Contains(i))
	}
	for i := 21; i <= 50; i++ {
		s.True(tree.Contains(i))
	}
	for i := 100; i <= 110; i++ {
		s.True(tree.Contains(i))
	}
}

func (s *BTreeTestSuite) TestBTree_HighMinDegree() {
	tree := NewBTree[int, string](10) // Wide nodes

	for i := 1; i <= 1000; i++ {
		tree.Insert(i, "value")
	}

	s.Equal(1000, tree.Size())

	// Height should be relatively low with high min degree
	s.LessOrEqual(tree.Height(), 4)
}

// ============================================================================
// Message Queue Specific Tests (Use Case)
// ============================================================================

func (s *BTreeTestSuite) TestBTree_MessageQueueUseCase() {
	// Simulate message queue offset index
	type MessageMeta struct {
		Offset    uint64
		Timestamp int64
		Size      int32
	}

	tree := NewBTree[uint64, MessageMeta](4)

	// Insert message metadata by offset
	for i := uint64(0); i < 1000; i++ {
		tree.Insert(i, MessageMeta{
			Offset:    i,
			Timestamp: int64(1000000 + i),
			Size:      100,
		})
	}

	s.Equal(1000, tree.Size())

	// Find message by offset
	meta, found := tree.Search(500)
	s.True(found)
	s.Equal(uint64(500), meta.Offset)

	// Range query for messages
	count := 0
	for entry := range tree.Range(100, 200) {
		s.GreaterOrEqual(entry.Key, uint64(100))
		s.LessOrEqual(entry.Key, uint64(200))
		count++
	}
	s.Equal(101, count)

	// Find first message >= offset (ceiling)
	key, _, found := tree.Ceiling(555)
	s.True(found)
	s.Equal(uint64(555), key)

	// Find last message <= offset (floor)
	key, _, found = tree.Floor(555)
	s.True(found)
	s.Equal(uint64(555), key)
}
