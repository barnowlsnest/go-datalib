package heap

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// MinHeapTestSuite tests min-heap functionality
type MinHeapTestSuite struct {
	suite.Suite
}

func (s *MinHeapTestSuite) TestNewMin_EmptyHeap() {
	h := NewMin[int]()

	s.Require().NotNil(h)
	s.Require().True(h.IsEmpty())
	s.Require().Equal(0, h.Size())
}

func (s *MinHeapTestSuite) TestPush_SingleElement() {
	h := NewMin[int]()

	h.Push(5)

	s.Require().False(h.IsEmpty())
	s.Require().Equal(1, h.Size())

	val, ok := h.Peek()
	s.Require().True(ok)
	s.Require().Equal(5, val)
}

func (s *MinHeapTestSuite) TestPush_MultipleElements_MaintainsMinAtTop() {
	h := NewMin[int]()

	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(1)
	h.Push(9)

	s.Require().Equal(5, h.Size())

	val, ok := h.Peek()
	s.Require().True(ok)
	s.Require().Equal(1, val, "Min element should be at top")
}

func (s *MinHeapTestSuite) TestPop_RemovesMinElement() {
	h := NewMin[int]()
	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(1)

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal(1, val, "Should pop minimum element")
	s.Require().Equal(3, h.Size())

	val, ok = h.Peek()
	s.Require().True(ok)
	s.Require().Equal(3, val, "Next minimum should be at top")
}

func (s *MinHeapTestSuite) TestPop_EmptyHeap() {
	h := NewMin[int]()

	val, ok := h.Pop()
	s.Require().False(ok)
	s.Require().Equal(0, val)
}

func (s *MinHeapTestSuite) TestPop_AllElements_InSortedOrder() {
	h := NewMin[int]()
	values := []int{5, 3, 7, 1, 9, 2, 8, 4, 6}

	for _, v := range values {
		h.Push(v)
	}

	var result []int
	for !h.IsEmpty() {
		val, ok := h.Pop()
		s.Require().True(ok)
		result = append(result, val)
	}

	// Should be in ascending order (min-heap)
	s.Require().Equal([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, result)
}

func (s *MinHeapTestSuite) TestPeek_DoesNotRemoveElement() {
	h := NewMin[int]()
	h.Push(5)
	h.Push(3)

	val1, ok1 := h.Peek()
	val2, ok2 := h.Peek()

	s.Require().True(ok1)
	s.Require().True(ok2)
	s.Require().Equal(3, val1)
	s.Require().Equal(3, val2)
	s.Require().Equal(2, h.Size(), "Size should not change after Peek")
}

func (s *MinHeapTestSuite) TestPeek_EmptyHeap() {
	h := NewMin[int]()

	val, ok := h.Peek()
	s.Require().False(ok)
	s.Require().Equal(0, val)
}

// MaxHeapTestSuite tests max-heap functionality
type MaxHeapTestSuite struct {
	suite.Suite
}

func (s *MaxHeapTestSuite) TestNewMax_EmptyHeap() {
	h := NewMax[int]()

	s.Require().NotNil(h)
	s.Require().True(h.IsEmpty())
	s.Require().Equal(0, h.Size())
}

func (s *MaxHeapTestSuite) TestPush_MaintainsMaxAtTop() {
	h := NewMax[int]()

	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(1)
	h.Push(9)

	s.Require().Equal(5, h.Size())

	val, ok := h.Peek()
	s.Require().True(ok)
	s.Require().Equal(9, val, "Max element should be at top")
}

func (s *MaxHeapTestSuite) TestPop_RemovesMaxElement() {
	h := NewMax[int]()
	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(9)

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal(9, val, "Should pop maximum element")

	val, ok = h.Peek()
	s.Require().True(ok)
	s.Require().Equal(7, val, "Next maximum should be at top")
}

func (s *MaxHeapTestSuite) TestPop_AllElements_InReverseSortedOrder() {
	h := NewMax[int]()
	values := []int{5, 3, 7, 1, 9, 2, 8, 4, 6}

	for _, v := range values {
		h.Push(v)
	}

	var result []int
	for !h.IsEmpty() {
		val, ok := h.Pop()
		s.Require().True(ok)
		result = append(result, val)
	}

	// Should be in descending order (max-heap)
	s.Require().Equal([]int{9, 8, 7, 6, 5, 4, 3, 2, 1}, result)
}

// CustomComparisonTestSuite tests heaps with custom comparison functions
type CustomComparisonTestSuite struct {
	suite.Suite
}

type Person struct {
	Name string
	Age  int
}

func (s *CustomComparisonTestSuite) TestCustomComparison_ByAge() {
	// Min-heap by age
	h := New(func(a, b Person) bool {
		return a.Age < b.Age
	})

	h.Push(Person{"Alice", 30})
	h.Push(Person{"Bob", 25})
	h.Push(Person{"Charlie", 35})
	h.Push(Person{"David", 20})

	youngest, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal("David", youngest.Name)
	s.Require().Equal(20, youngest.Age)
}

func (s *CustomComparisonTestSuite) TestCustomComparison_ByName() {
	// Min-heap by name (alphabetical)
	h := New(func(a, b Person) bool {
		return a.Name < b.Name
	})

	h.Push(Person{"Charlie", 35})
	h.Push(Person{"Alice", 30})
	h.Push(Person{"David", 20})
	h.Push(Person{"Bob", 25})

	first, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal("Alice", first.Name)
}

func (s *CustomComparisonTestSuite) TestCustomComparison_ReverseOrder() {
	// Max-heap for strings (reverse alphabetical)
	h := New(func(a, b string) bool {
		return a > b
	})

	h.Push("apple")
	h.Push("banana")
	h.Push("cherry")
	h.Push("date")

	top, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal("date", top)
}

// FromSliceTestSuite tests heap creation from existing slices
type FromSliceTestSuite struct {
	suite.Suite
}

func (s *FromSliceTestSuite) TestFromSlice_CreatesValidMinHeap() {
	input := []int{5, 3, 7, 1, 9, 2, 8}

	h := FromSlice(input, func(a, b int) bool { return a < b })

	s.Require().Equal(7, h.Size())

	// Should pop in sorted order
	var result []int
	for !h.IsEmpty() {
		val, _ := h.Pop()
		result = append(result, val)
	}

	s.Require().Equal([]int{1, 2, 3, 5, 7, 8, 9}, result)
}

func (s *FromSliceTestSuite) TestFromSlice_DoesNotModifyOriginal() {
	input := []int{5, 3, 7, 1, 9}
	original := make([]int, len(input))
	copy(original, input)

	h := FromSlice(input, func(a, b int) bool { return a < b })

	// Modify heap
	h.Pop()
	h.Push(100)

	// Original should be unchanged
	s.Require().Equal(original, input)
}

func (s *FromSliceTestSuite) TestFromSlice_EmptySlice() {
	h := FromSlice([]int{}, func(a, b int) bool { return a < b })

	s.Require().True(h.IsEmpty())
	s.Require().Equal(0, h.Size())
}

func (s *FromSliceTestSuite) TestFromSlice_SingleElement() {
	h := FromSlice([]int{42}, func(a, b int) bool { return a < b })

	s.Require().Equal(1, h.Size())

	val, ok := h.Peek()
	s.Require().True(ok)
	s.Require().Equal(42, val)
}

// HeapOperationsTestSuite tests various heap operations
type HeapOperationsTestSuite struct {
	suite.Suite
}

func (s *HeapOperationsTestSuite) TestClear_RemovesAllElements() {
	h := NewMin[int]()
	h.Push(1)
	h.Push(2)
	h.Push(3)

	s.Require().Equal(3, h.Size())

	h.Clear()

	s.Require().True(h.IsEmpty())
	s.Require().Equal(0, h.Size())

	_, ok := h.Peek()
	s.Require().False(ok)
}

func (s *HeapOperationsTestSuite) TestToSlice_ReturnsHeapData() {
	h := NewMin[int]()
	h.Push(3)
	h.Push(1)
	h.Push(4)

	slice := h.ToSlice()

	s.Require().Len(slice, 3)
	s.Require().Contains(slice, 1)
	s.Require().Contains(slice, 3)
	s.Require().Contains(slice, 4)
}

func (s *HeapOperationsTestSuite) TestToSlice_DoesNotAffectHeap() {
	h := NewMin[int]()
	h.Push(1)
	h.Push(2)

	slice := h.ToSlice()
	slice[0] = 999

	val, _ := h.Peek()
	s.Require().NotEqual(999, val, "Modifying slice should not affect heap")
}

func (s *HeapOperationsTestSuite) TestNewWithCapacity_PreallocatesSpace() {
	h := NewWithCapacity(func(a, b int) bool { return a < b }, 100)

	s.Require().True(h.IsEmpty())
	s.Require().Equal(0, h.Size())

	// Should be able to add 100 elements without reallocation
	for i := 0; i < 100; i++ {
		h.Push(i)
	}

	s.Require().Equal(100, h.Size())
}

// EdgeCasesTestSuite tests edge cases and boundary conditions
type EdgeCasesTestSuite struct {
	suite.Suite
}

func (s *EdgeCasesTestSuite) TestSingleElement_PushPopPeek() {
	h := NewMin[int]()

	h.Push(42)
	s.Require().Equal(1, h.Size())

	val, ok := h.Peek()
	s.Require().True(ok)
	s.Require().Equal(42, val)
	s.Require().Equal(1, h.Size(), "Peek should not change size")

	val, ok = h.Pop()
	s.Require().True(ok)
	s.Require().Equal(42, val)
	s.Require().Equal(0, h.Size())
	s.Require().True(h.IsEmpty())
}

func (s *EdgeCasesTestSuite) TestDuplicateValues() {
	h := NewMin[int]()

	h.Push(5)
	h.Push(5)
	h.Push(5)

	s.Require().Equal(3, h.Size())

	for i := 0; i < 3; i++ {
		val, ok := h.Pop()
		s.Require().True(ok)
		s.Require().Equal(5, val)
	}

	s.Require().True(h.IsEmpty())
}

func (s *EdgeCasesTestSuite) TestAlreadySortedInput_MinHeap() {
	h := NewMin[int]()

	// Insert in ascending order
	for i := 1; i <= 10; i++ {
		h.Push(i)
	}

	// Should still pop in order
	for i := 1; i <= 10; i++ {
		val, ok := h.Pop()
		s.Require().True(ok)
		s.Require().Equal(i, val)
	}
}

func (s *EdgeCasesTestSuite) TestAlreadySortedInput_MaxHeap() {
	h := NewMax[int]()

	// Insert in descending order
	for i := 10; i >= 1; i-- {
		h.Push(i)
	}

	// Should still pop in descending order
	for i := 10; i >= 1; i-- {
		val, ok := h.Pop()
		s.Require().True(ok)
		s.Require().Equal(i, val)
	}
}

func (s *EdgeCasesTestSuite) TestZeroValue() {
	h := NewMin[int]()

	h.Push(0)
	h.Push(1)
	h.Push(-1)

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal(-1, val)

	val, ok = h.Pop()
	s.Require().True(ok)
	s.Require().Equal(0, val)
}

func (s *EdgeCasesTestSuite) TestNegativeNumbers() {
	h := NewMin[int]()

	h.Push(-5)
	h.Push(-3)
	h.Push(-7)
	h.Push(-1)

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal(-7, val)
}

func (s *EdgeCasesTestSuite) TestLargeHeap() {
	h := NewMin[int]()

	// Insert 1000 elements
	for i := 1000; i > 0; i-- {
		h.Push(i)
	}

	s.Require().Equal(1000, h.Size())

	// Should pop in order
	for i := 1; i <= 1000; i++ {
		val, ok := h.Pop()
		s.Require().True(ok)
		s.Require().Equal(i, val)
	}

	s.Require().True(h.IsEmpty())
}

// TypesTestSuite tests heap with different types
type TypesTestSuite struct {
	suite.Suite
}

func (s *TypesTestSuite) TestFloat64() {
	h := NewMin[float64]()

	h.Push(3.14)
	h.Push(2.71)
	h.Push(1.41)

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().InDelta(1.41, val, 0.001)
}

func (s *TypesTestSuite) TestString() {
	h := NewMin[string]()

	h.Push("charlie")
	h.Push("alice")
	h.Push("bob")

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal("alice", val)
}

func (s *TypesTestSuite) TestUint64() {
	h := NewMax[uint64]()

	h.Push(100)
	h.Push(200)
	h.Push(50)

	val, ok := h.Pop()
	s.Require().True(ok)
	s.Require().Equal(uint64(200), val)
}

// Test suite runners
func TestMinHeapTestSuite(t *testing.T) {
	suite.Run(t, new(MinHeapTestSuite))
}

func TestMaxHeapTestSuite(t *testing.T) {
	suite.Run(t, new(MaxHeapTestSuite))
}

func TestCustomComparisonTestSuite(t *testing.T) {
	suite.Run(t, new(CustomComparisonTestSuite))
}

func TestFromSliceTestSuite(t *testing.T) {
	suite.Run(t, new(FromSliceTestSuite))
}

func TestHeapOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(HeapOperationsTestSuite))
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}
