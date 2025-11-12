package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConstructorTestSuite tests Fenwick Tree creation
type ConstructorTestSuite struct {
	suite.Suite
}

func (s *ConstructorTestSuite) TestNew_ValidSize() {
	ft := NewFenwick[int](10)

	s.Require().NotNil(ft)
	s.Require().Equal(10, ft.Size())

	// All values should be zero initially
	for i := 1; i <= 10; i++ {
		val := ft.Get(i)
		s.Require().Equal(0, val)
	}
}

func (s *ConstructorTestSuite) TestNew_ZeroSize() {
	ft := NewFenwick[int](0)

	s.Require().NotNil(ft)
	s.Require().Equal(0, ft.Size())
}

func (s *ConstructorTestSuite) TestNew_NegativeSize() {
	ft := NewFenwick[int](-5)

	s.Require().NotNil(ft)
	s.Require().Equal(0, ft.Size())
}

func (s *ConstructorTestSuite) TestNew_DifferentTypes() {
	testCases := []struct {
		name string
		test func()
	}{
		{
			name: "int",
			test: func() {
				ft := NewFenwick[int](5)
				s.Require().Equal(5, ft.Size())
			},
		},
		{
			name: "int64",
			test: func() {
				ft := NewFenwick[int64](5)
				s.Require().Equal(5, ft.Size())
			},
		},
		{
			name: "float64",
			test: func() {
				ft := NewFenwick[float64](5)
				s.Require().Equal(5, ft.Size())
			},
		},
		{
			name: "uint",
			test: func() {
				ft := NewFenwick[uint](5)
				s.Require().Equal(5, ft.Size())
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, tc.test)
	}
}

func (s *ConstructorTestSuite) TestFromSlice_ValidData() {
	data := []int{3, 2, -1, 6, 5, 4, -3, 3, 7, 2, 3}
	ft := FromSlice(data)

	s.Require().Equal(11, ft.Size())

	// Verify each element is correctly set
	for i := 0; i < len(data); i++ {
		val := ft.Get(i + 1) // Convert to 1-indexed
		s.Require().Equal(data[i], val)
	}
}

func (s *ConstructorTestSuite) TestFromSlice_EmptySlice() {
	ft := FromSlice([]int{})

	s.Require().Equal(0, ft.Size())
}

func (s *ConstructorTestSuite) TestFromSlice_SingleElement() {
	ft := FromSlice([]int{42})

	s.Require().Equal(1, ft.Size())
	s.Require().Equal(42, ft.Get(1))
}

// UpdateTestSuite tests Update operations
type UpdateTestSuite struct {
	suite.Suite
}

func (s *UpdateTestSuite) TestUpdate_SingleElement() {
	ft := NewFenwick[int](5)

	ft.Update(3, 10)

	s.Require().Equal(10, ft.Get(3))
	s.Require().Equal(0, ft.Get(1))
	s.Require().Equal(0, ft.Get(2))
	s.Require().Equal(0, ft.Get(4))
}

func (s *UpdateTestSuite) TestUpdate_MultipleElements() {
	ft := NewFenwick[int](5)

	ft.Update(1, 5)
	ft.Update(2, 3)
	ft.Update(3, 7)
	ft.Update(4, 2)
	ft.Update(5, 8)

	s.Require().Equal(5, ft.Get(1))
	s.Require().Equal(3, ft.Get(2))
	s.Require().Equal(7, ft.Get(3))
	s.Require().Equal(2, ft.Get(4))
	s.Require().Equal(8, ft.Get(5))
}

func (s *UpdateTestSuite) TestUpdate_SameIndexMultipleTimes() {
	ft := NewFenwick[int](5)

	ft.Update(3, 10)
	ft.Update(3, 5)
	ft.Update(3, -3)

	// Should be cumulative: 10 + 5 - 3 = 12
	s.Require().Equal(12, ft.Get(3))
}

func (s *UpdateTestSuite) TestUpdate_NegativeValues() {
	ft := NewFenwick[int](5)

	ft.Update(2, 10)
	ft.Update(2, -5)

	s.Require().Equal(5, ft.Get(2))
}

func (s *UpdateTestSuite) TestUpdate_OutOfBounds() {
	ft := NewFenwick[int](5)

	// These should be silently ignored
	ft.Update(0, 10)
	ft.Update(-1, 10)
	ft.Update(6, 10)
	ft.Update(100, 10)

	// All values should still be zero
	for i := 1; i <= 5; i++ {
		s.Require().Equal(0, ft.Get(i))
	}
}

func (s *UpdateTestSuite) TestUpdate_Float64() {
	ft := NewFenwick[float64](3)

	ft.Update(1, 3.14)
	ft.Update(2, 2.71)
	ft.Update(3, 1.41)

	s.Require().InDelta(3.14, ft.Get(1), 0.001)
	s.Require().InDelta(2.71, ft.Get(2), 0.001)
	s.Require().InDelta(1.41, ft.Get(3), 0.001)
}

// QueryTestSuite tests Query operations
type QueryTestSuite struct {
	suite.Suite
}

func (s *QueryTestSuite) TestQuery_PrefixSum() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	testCases := []struct {
		index       int
		expectedSum int
	}{
		{1, 1},
		{2, 3},
		{3, 6},
		{4, 10},
		{5, 15},
	}

	for _, tc := range testCases {
		sum := ft.Query(tc.index)
		s.Require().Equal(tc.expectedSum, sum, "Query(%d) should return %d", tc.index, tc.expectedSum)
	}
}

func (s *QueryTestSuite) TestQuery_AfterUpdates() {
	ft := NewFenwick[int](5)

	ft.Update(1, 3)
	ft.Update(3, 5)
	ft.Update(5, 7)

	s.Require().Equal(3, ft.Query(1))  // 3
	s.Require().Equal(3, ft.Query(2))  // 3 + 0
	s.Require().Equal(8, ft.Query(3))  // 3 + 0 + 5
	s.Require().Equal(8, ft.Query(4))  // 3 + 0 + 5 + 0
	s.Require().Equal(15, ft.Query(5)) // 3 + 0 + 5 + 0 + 7
}

func (s *QueryTestSuite) TestQuery_NegativeValues() {
	ft := FromSlice([]int{5, -3, 7, -2, 4})

	s.Require().Equal(5, ft.Query(1))
	s.Require().Equal(2, ft.Query(2))
	s.Require().Equal(9, ft.Query(3))
	s.Require().Equal(7, ft.Query(4))
	s.Require().Equal(11, ft.Query(5))
}

func (s *QueryTestSuite) TestQuery_OutOfBounds() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	s.Require().Equal(0, ft.Query(0))
	s.Require().Equal(0, ft.Query(-1))
	s.Require().Equal(15, ft.Query(6))   // Clamped to size
	s.Require().Equal(15, ft.Query(100)) // Clamped to size
}

// RangeQueryTestSuite tests RangeQuery operations
type RangeQueryTestSuite struct {
	suite.Suite
}

func (s *RangeQueryTestSuite) TestRangeQuery_ValidRanges() {
	ft := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	testCases := []struct {
		left        int
		right       int
		expectedSum int
	}{
		{1, 1, 1},
		{1, 5, 15},
		{3, 7, 25},
		{5, 10, 45},
		{1, 10, 55},
		{8, 10, 27},
	}

	for _, tc := range testCases {
		sum := ft.RangeQuery(tc.left, tc.right)
		s.Require().Equal(tc.expectedSum, sum, "RangeQuery(%d, %d) should return %d", tc.left, tc.right, tc.expectedSum)
	}
}

func (s *RangeQueryTestSuite) TestRangeQuery_SingleElement() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	s.Require().Equal(1, ft.RangeQuery(1, 1))
	s.Require().Equal(3, ft.RangeQuery(3, 3))
	s.Require().Equal(5, ft.RangeQuery(5, 5))
}

func (s *RangeQueryTestSuite) TestRangeQuery_WithNegatives() {
	ft := FromSlice([]int{5, -3, 7, -2, 4, -1, 6})

	s.Require().Equal(2, ft.RangeQuery(1, 2))
	s.Require().Equal(2, ft.RangeQuery(2, 4))
	s.Require().Equal(14, ft.RangeQuery(3, 7))
}

func (s *RangeQueryTestSuite) TestRangeQuery_InvalidRanges() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	testCases := []struct {
		left  int
		right int
		name  string
	}{
		{5, 3, "left > right"},
		{0, 3, "left <= 0"},
		{-1, 3, "left < 0"},
		{1, 6, "right > size"},
		{1, 100, "right >> size"},
	}

	for _, tc := range testCases {
		sum := ft.RangeQuery(tc.left, tc.right)
		s.Require().Equal(0, sum, "%s should return 0", tc.name)
	}
}

// SetAndGetTestSuite tests Set and Get operations
type SetAndGetTestSuite struct {
	suite.Suite
}

func (s *SetAndGetTestSuite) TestSet_NewValue() {
	ft := NewFenwick[int](5)

	ft.Set(3, 42)

	s.Require().Equal(42, ft.Get(3))
}

func (s *SetAndGetTestSuite) TestSet_OverwriteValue() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	s.Require().Equal(3, ft.Get(3))

	ft.Set(3, 10)

	s.Require().Equal(10, ft.Get(3))

	// Verify prefix sum is updated correctly
	s.Require().Equal(22, ft.Query(5))
}

func (s *SetAndGetTestSuite) TestSet_MultipleOverwrites() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	ft.Set(2, 10)
	ft.Set(4, 20)
	ft.Set(2, 15)

	s.Require().Equal(15, ft.Get(2))
	s.Require().Equal(20, ft.Get(4))
}

func (s *SetAndGetTestSuite) TestSet_OutOfBounds() {
	ft := NewFenwick[int](5)

	// These should be silently ignored
	ft.Set(0, 10)
	ft.Set(-1, 10)
	ft.Set(6, 10)

	// All values should still be zero
	for i := 1; i <= 5; i++ {
		s.Require().Equal(0, ft.Get(i))
	}
}

func (s *SetAndGetTestSuite) TestGet_AfterMultipleUpdates() {
	ft := NewFenwick[int](3)

	ft.Update(2, 5)
	ft.Update(2, 3)
	ft.Update(2, 2)

	// Should be cumulative: 5 + 3 + 2 = 10
	s.Require().Equal(10, ft.Get(2))
}

func (s *SetAndGetTestSuite) TestGet_OutOfBounds() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	s.Require().Equal(0, ft.Get(0))
	s.Require().Equal(0, ft.Get(-1))
	s.Require().Equal(0, ft.Get(6))
	s.Require().Equal(0, ft.Get(100))
}

// UtilityTestSuite tests utility methods
type UtilityTestSuite struct {
	suite.Suite
}

func (s *UtilityTestSuite) TestClear() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	s.Require().Equal(15, ft.Query(5))

	ft.Clear()

	// All values should be zero
	for i := 1; i <= 5; i++ {
		s.Require().Equal(0, ft.Get(i))
	}
	s.Require().Equal(0, ft.Query(5))
}

func (s *UtilityTestSuite) TestToSlice() {
	original := []int{3, 2, -1, 6, 5, 4, -3, 3}
	ft := FromSlice(original)

	result := ft.ToSlice()

	s.Require().Equal(original, result)
}

func (s *UtilityTestSuite) TestToSlice_AfterUpdates() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	ft.Update(2, 10)
	ft.Set(4, 20)

	result := ft.ToSlice()

	s.Require().Equal([]int{1, 12, 3, 20, 5}, result)
}

func (s *UtilityTestSuite) TestToSlice_EmptyTree() {
	ft := NewFenwick[int](0)

	result := ft.ToSlice()

	s.Require().Empty(result)
}

func (s *UtilityTestSuite) TestToSlice_DoesNotAffectTree() {
	ft := FromSlice([]int{1, 2, 3})

	slice := ft.ToSlice()
	slice[0] = 999

	s.Require().Equal(1, ft.Get(1), "Modifying slice should not affect tree")
}

func (s *UtilityTestSuite) TestSize() {
	testCases := []struct {
		size int
	}{
		{0},
		{1},
		{10},
		{100},
		{1000},
	}

	for _, tc := range testCases {
		ft := NewFenwick[int](tc.size)
		s.Require().Equal(tc.size, ft.Size())
	}
}

// EdgeCasesTestSuite tests edge cases
type EdgeCasesTestSuite struct {
	suite.Suite
}

func (s *EdgeCasesTestSuite) TestSingleElementTree() {
	ft := NewFenwick[int](1)

	ft.Update(1, 42)

	s.Require().Equal(42, ft.Get(1))
	s.Require().Equal(42, ft.Query(1))
	s.Require().Equal(42, ft.RangeQuery(1, 1))
}

func (s *EdgeCasesTestSuite) TestAllZeros() {
	ft := NewFenwick[int](10)

	for i := 1; i <= 10; i++ {
		s.Require().Equal(0, ft.Get(i))
	}

	s.Require().Equal(0, ft.Query(10))
	s.Require().Equal(0, ft.RangeQuery(1, 10))
}

func (s *EdgeCasesTestSuite) TestAllSameValue() {
	ft := NewFenwick[int](5)

	for i := 1; i <= 5; i++ {
		ft.Update(i, 7)
	}

	for i := 1; i <= 5; i++ {
		s.Require().Equal(7, ft.Get(i))
	}

	s.Require().Equal(35, ft.Query(5))
}

func (s *EdgeCasesTestSuite) TestAlternatingValues() {
	ft := NewFenwick[int](10)

	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			ft.Update(i, 1)
		} else {
			ft.Update(i, -1)
		}
	}

	// Sum should be 0: (-1 + 1) * 5 = 0
	s.Require().Equal(0, ft.Query(10))
}

func (s *EdgeCasesTestSuite) TestLargeTree() {
	ft := NewFenwick[int](10000)

	// Update every 100th element
	for i := 100; i <= 10000; i += 100 {
		ft.Update(i, i)
	}

	// Verify a few queries
	s.Require().Equal(100, ft.Get(100))
	s.Require().Equal(500, ft.Get(500))
	s.Require().Equal(10000, ft.Get(10000))

	// Verify range query
	sum := ft.RangeQuery(1, 1000)
	expected := 100 + 200 + 300 + 400 + 500 + 600 + 700 + 800 + 900 + 1000
	s.Require().Equal(expected, sum)
}

func (s *EdgeCasesTestSuite) TestPowerOfTwo() {
	// Test with powers of 2 (important for bit manipulation)
	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		ft := NewFenwick[int](size)
		ft.Update(size, size)
		s.Require().Equal(size, ft.Get(size))
	}
}

// ComplexOperationsTestSuite tests complex operation sequences
type ComplexOperationsTestSuite struct {
	suite.Suite
}

func (s *ComplexOperationsTestSuite) TestMixedOperations() {
	ft := NewFenwick[int](10)

	// Initialize with some values
	for i := 1; i <= 10; i++ {
		ft.Update(i, i)
	}

	// Mix of operations
	ft.Set(5, 100)
	ft.Update(3, 10)
	ft.Update(7, -5)

	s.Require().Equal(13, ft.Get(3)) // 3 + 10
	s.Require().Equal(100, ft.Get(5))
	s.Require().Equal(2, ft.Get(7)) // 7 - 5

	// Verify range queries
	sum := ft.RangeQuery(3, 7)
	expected := 13 + 4 + 100 + 6 + 2
	s.Require().Equal(expected, sum)
}

func (s *ComplexOperationsTestSuite) TestFrequentUpdates() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	for i := 0; i < 100; i++ {
		ft.Update(3, 1)
	}

	s.Require().Equal(103, ft.Get(3))   // 3 + 100
	s.Require().Equal(115, ft.Query(5)) // 1 + 2 + 103 + 4 + 5
}

func (s *ComplexOperationsTestSuite) TestClearAndReuse() {
	ft := FromSlice([]int{1, 2, 3, 4, 5})

	s.Require().Equal(15, ft.Query(5))

	ft.Clear()

	for i := 1; i <= 5; i++ {
		ft.Update(i, i*2)
	}

	s.Require().Equal(30, ft.Query(5)) // 2 + 4 + 6 + 8 + 10
}

func (s *ComplexOperationsTestSuite) TestSequentialRangeQueries() {
	ft := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Query all consecutive pairs
	for i := 1; i < 10; i++ {
		sum := ft.RangeQuery(i, i+1)
		expected := i + (i + 1)
		s.Require().Equal(expected, sum)
	}
}

// TypesTestSuite tests different numeric types
type TypesHeapTestSuite struct {
	suite.Suite
}

func (s *TypesHeapTestSuite) TestInt32() {
	ft := NewFenwick[int32](5)

	ft.Update(1, 100)
	ft.Update(3, 200)
	ft.Update(5, 300)

	s.Require().Equal(int32(100), ft.Get(1))
	s.Require().Equal(int32(200), ft.Get(3))
	s.Require().Equal(int32(600), ft.Query(5))
}

func (s *TypesHeapTestSuite) TestInt64() {
	ft := NewFenwick[int64](3)

	ft.Update(1, 1000000000000)
	ft.Update(2, 2000000000000)
	ft.Update(3, 3000000000000)

	s.Require().Equal(int64(6000000000000), ft.Query(3))
}

func (s *TypesHeapTestSuite) TestFloat32() {
	ft := NewFenwick[float32](3)

	ft.Update(1, 3.14)
	ft.Update(2, 2.71)
	ft.Update(3, 1.41)

	s.Require().InDelta(float32(3.14), ft.Get(1), 0.001)
	s.Require().InDelta(float32(7.26), ft.Query(3), 0.01)
}

func (s *TypesHeapTestSuite) TestUint() {
	ft := NewFenwick[uint](5)

	for i := 1; i <= 5; i++ {
		ft.Update(i, uint(i*10))
	}

	s.Require().Equal(uint(150), ft.Query(5))
}

// Test suite runners
func TestConstructorTestSuite(t *testing.T) {
	suite.Run(t, new(ConstructorTestSuite))
}

func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}

func TestRangeQueryTestSuite(t *testing.T) {
	suite.Run(t, new(RangeQueryTestSuite))
}

func TestSetAndGetTestSuite(t *testing.T) {
	suite.Run(t, new(SetAndGetTestSuite))
}

func TestUtilityTestSuite(t *testing.T) {
	suite.Run(t, new(UtilityTestSuite))
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}

func TestComplexOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(ComplexOperationsTestSuite))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesHeapTestSuite))
}
