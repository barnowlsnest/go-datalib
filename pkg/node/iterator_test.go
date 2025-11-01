package node

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ForwardIteratorTestSuite tests forward iteration functionality
type ForwardIteratorTestSuite struct {
	suite.Suite
}

func (s *ForwardIteratorTestSuite) TestForward_NilNode() {
	it := Forward(nil)

	s.Require().NotNil(it)
	s.Require().Nil(it.cur)
}

func (s *ForwardIteratorTestSuite) TestForward_SingleNode() {
	node := New(1, nil, nil)
	it := Forward(node)

	s.Require().NotNil(it)
	s.Require().Equal(node, it.cur)
}

func (s *ForwardIteratorTestSuite) TestCurr_NilNode() {
	it := Forward(nil)

	curr, err := it.Curr()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(curr)
}

func (s *ForwardIteratorTestSuite) TestCurr_ValidNode() {
	node := New(1, nil, nil)
	it := Forward(node)

	curr, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node, curr)
	s.Require().Equal(uint64(1), curr.ID())
}

func (s *ForwardIteratorTestSuite) TestCurr_DoesNotAdvance() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node1.WithNext(node2)

	it := Forward(node1)

	// Call Curr multiple times
	curr1, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node1, curr1)

	curr2, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node1, curr2)
	s.Require().Equal(curr1, curr2)
}

func (s *ForwardIteratorTestSuite) TestHasNext_NilNode() {
	it := Forward(nil)

	s.Require().False(it.HasNext())
}

func (s *ForwardIteratorTestSuite) TestHasNext_SingleNode() {
	node := New(1, nil, nil)
	it := Forward(node)

	s.Require().False(it.HasNext())
}

func (s *ForwardIteratorTestSuite) TestHasNext_WithNextNode() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node1.WithNext(node2)

	it := Forward(node1)

	s.Require().True(it.HasNext())
}

func (s *ForwardIteratorTestSuite) TestNext_NilNode() {
	it := Forward(nil)

	next, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next)
}

func (s *ForwardIteratorTestSuite) TestNext_SingleNode() {
	node := New(1, nil, nil)
	it := Forward(node)

	// First Next() returns nil node with no error (reached end)
	next, err := it.Next()
	s.Require().NoError(err)
	s.Require().Nil(next)

	// Second Next() returns ErrEOI
	next2, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next2)
}

func (s *ForwardIteratorTestSuite) TestNext_TwoNodes() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node1.WithNext(node2)

	it := Forward(node1)

	// First Next() should return node2
	next, err := it.Next()
	s.Require().NoError(err)
	s.Require().NotNil(next)
	s.Require().Equal(node2, next)
	s.Require().Equal(uint64(2), next.ID())

	// Second Next() returns nil (reached end)
	next2, err := it.Next()
	s.Require().NoError(err)
	s.Require().Nil(next2)

	// Third Next() returns ErrEOI
	next3, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next3)
}

func (s *ForwardIteratorTestSuite) TestNext_LinearChain() {
	// Create chain: 1 -> 2 -> 3 -> 4 -> 5
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = New(uint64(i+1), nil, nil)
	}
	for i := 0; i < 4; i++ {
		nodes[i].WithNext(nodes[i+1])
	}

	it := Forward(nodes[0])

	// Iterate through all nodes
	for i := 1; i < 5; i++ {
		s.Require().True(it.HasNext())
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().Equal(nodes[i], next)
		s.Require().Equal(uint64(i+1), next.ID())
	}

	// No more nodes
	s.Require().False(it.HasNext())
	// First Next() after exhaustion returns (nil, nil)
	next, err := it.Next()
	s.Require().NoError(err)
	s.Require().Nil(next)

	// Second Next() returns ErrEOI
	next2, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next2)
}

func (s *ForwardIteratorTestSuite) TestCurrAfterNext() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)
	node1.WithNext(node2)
	node2.WithNext(node3)

	it := Forward(node1)

	// Initial Curr
	curr, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node1, curr)

	// After first Next
	_, _ = it.Next()
	curr, err = it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node2, curr)

	// After second Next
	_, _ = it.Next()
	curr, err = it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node3, curr)

	// After final Next (should be EOI)
	_, _ = it.Next()
	curr, err = it.Curr()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(curr)
}

// BackwardIteratorTestSuite tests backward iteration functionality
type BackwardIteratorTestSuite struct {
	suite.Suite
}

func (s *BackwardIteratorTestSuite) TestBackward_NilNode() {
	it := Backward(nil)

	s.Require().NotNil(it)
	s.Require().Nil(it.cur)
}

func (s *BackwardIteratorTestSuite) TestBackward_SingleNode() {
	node := New(1, nil, nil)
	it := Backward(node)

	s.Require().NotNil(it)
	s.Require().Equal(node, it.cur)
}

func (s *BackwardIteratorTestSuite) TestCurr_NilNode() {
	it := Backward(nil)

	curr, err := it.Curr()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(curr)
}

func (s *BackwardIteratorTestSuite) TestCurr_ValidNode() {
	node := New(1, nil, nil)
	it := Backward(node)

	curr, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node, curr)
	s.Require().Equal(uint64(1), curr.ID())
}

func (s *BackwardIteratorTestSuite) TestCurr_DoesNotAdvance() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node2.WithPrev(node1)

	it := Backward(node2)

	// Call Curr multiple times
	curr1, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node2, curr1)

	curr2, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node2, curr2)
	s.Require().Equal(curr1, curr2)
}

func (s *BackwardIteratorTestSuite) TestHasNext_NilNode() {
	it := Backward(nil)

	s.Require().False(it.HasNext())
}

func (s *BackwardIteratorTestSuite) TestHasNext_SingleNode() {
	node := New(1, nil, nil)
	it := Backward(node)

	s.Require().False(it.HasNext())
}

func (s *BackwardIteratorTestSuite) TestHasNext_WithPrevNode() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node2.WithPrev(node1)

	it := Backward(node2)

	s.Require().True(it.HasNext())
}

func (s *BackwardIteratorTestSuite) TestNext_NilNode() {
	it := Backward(nil)

	next, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next)
}

func (s *BackwardIteratorTestSuite) TestNext_SingleNode() {
	node := New(1, nil, nil)
	it := Backward(node)

	// First Next() returns nil node with no error (reached end)
	next, err := it.Next()
	s.Require().NoError(err)
	s.Require().Nil(next)

	// Second Next() returns ErrEOI
	next2, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next2)
}

func (s *BackwardIteratorTestSuite) TestNext_TwoNodes() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node2.WithPrev(node1)

	it := Backward(node2)

	// First Next() should return node1
	next, err := it.Next()
	s.Require().NoError(err)
	s.Require().NotNil(next)
	s.Require().Equal(node1, next)
	s.Require().Equal(uint64(1), next.ID())

	// Second Next() returns nil (reached end)
	next2, err := it.Next()
	s.Require().NoError(err)
	s.Require().Nil(next2)

	// Third Next() returns ErrEOI
	next3, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next3)
}

func (s *BackwardIteratorTestSuite) TestNext_LinearChain() {
	// Create chain: 1 <- 2 <- 3 <- 4 <- 5
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = New(uint64(i+1), nil, nil)
	}
	for i := 1; i < 5; i++ {
		nodes[i].WithPrev(nodes[i-1])
	}

	it := Backward(nodes[4]) // Start from node 5

	// Iterate backward through all nodes
	for i := 3; i >= 0; i-- {
		s.Require().True(it.HasNext())
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().Equal(nodes[i], next)
		s.Require().Equal(uint64(i+1), next.ID())
	}

	// No more nodes
	s.Require().False(it.HasNext())
	// First Next() after exhaustion returns (nil, nil)
	next, err := it.Next()
	s.Require().NoError(err)
	s.Require().Nil(next)

	// Second Next() returns ErrEOI
	next2, err := it.Next()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(next2)
}

func (s *BackwardIteratorTestSuite) TestCurrAfterNext() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)
	node2.WithPrev(node1)
	node3.WithPrev(node2)

	it := Backward(node3)

	// Initial Curr
	curr, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node3, curr)

	// After first Next
	_, _ = it.Next()
	curr, err = it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node2, curr)

	// After second Next
	_, _ = it.Next()
	curr, err = it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(node1, curr)

	// After final Next (should be EOI)
	_, _ = it.Next()
	curr, err = it.Curr()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrEOI)
	s.Require().Nil(curr)
}

// IteratorDataIntegrityTestSuite tests data integrity during iteration
type IteratorDataIntegrityTestSuite struct {
	suite.Suite
}

func (s *IteratorDataIntegrityTestSuite) TestForwardIterator_PreservesNodeData() {
	// Create chain with various ID values
	node1 := New(100, nil, nil)
	node2 := New(200, nil, nil)
	node3 := New(300, nil, nil)
	node1.WithNext(node2)
	node2.WithNext(node3)

	it := Forward(node1)

	// Verify data integrity at each step
	curr, _ := it.Curr()
	s.Require().Equal(uint64(100), curr.ID())

	next1, _ := it.Next()
	s.Require().Equal(uint64(200), next1.ID())

	next2, _ := it.Next()
	s.Require().Equal(uint64(300), next2.ID())
}

func (s *IteratorDataIntegrityTestSuite) TestBackwardIterator_PreservesNodeData() {
	// Create chain with various ID values
	node1 := New(100, nil, nil)
	node2 := New(200, nil, nil)
	node3 := New(300, nil, nil)
	node2.WithPrev(node1)
	node3.WithPrev(node2)

	it := Backward(node3)

	// Verify data integrity at each step
	curr, _ := it.Curr()
	s.Require().Equal(uint64(300), curr.ID())

	next1, _ := it.Next()
	s.Require().Equal(uint64(200), next1.ID())

	next2, _ := it.Next()
	s.Require().Equal(uint64(100), next2.ID())
}

func (s *IteratorDataIntegrityTestSuite) TestForwardIterator_DoublyLinkedChain() {
	// Create doubly-linked chain: 1 <-> 2 <-> 3
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	node1.WithNext(node2)
	node2.WithPrev(node1)
	node2.WithNext(node3)
	node3.WithPrev(node2)

	it := Forward(node1)

	// Verify forward iteration
	curr, _ := it.Curr()
	s.Require().Equal(uint64(1), curr.ID())
	s.Require().NotNil(curr.Next())
	s.Require().Nil(curr.Prev())

	next1, _ := it.Next()
	s.Require().Equal(uint64(2), next1.ID())
	s.Require().NotNil(next1.Next())
	s.Require().NotNil(next1.Prev())

	next2, _ := it.Next()
	s.Require().Equal(uint64(3), next2.ID())
	s.Require().Nil(next2.Next())
	s.Require().NotNil(next2.Prev())
}

func (s *IteratorDataIntegrityTestSuite) TestBackwardIterator_DoublyLinkedChain() {
	// Create doubly-linked chain: 1 <-> 2 <-> 3
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	node1.WithNext(node2)
	node2.WithPrev(node1)
	node2.WithNext(node3)
	node3.WithPrev(node2)

	it := Backward(node3)

	// Verify backward iteration
	curr, _ := it.Curr()
	s.Require().Equal(uint64(3), curr.ID())
	s.Require().Nil(curr.Next())
	s.Require().NotNil(curr.Prev())

	next1, _ := it.Next()
	s.Require().Equal(uint64(2), next1.ID())
	s.Require().NotNil(next1.Next())
	s.Require().NotNil(next1.Prev())

	next2, _ := it.Next()
	s.Require().Equal(uint64(1), next2.ID())
	s.Require().NotNil(next2.Next())
	s.Require().Nil(next2.Prev())
}

// IteratorNilSafetyTestSuite tests that iterators handle nil nodes safely
type IteratorNilSafetyTestSuite struct {
	suite.Suite
}

func (s *IteratorNilSafetyTestSuite) TestForwardIterator_NilChainOperations() {
	it := Forward(nil)

	// All operations should return errors, not panic
	s.Require().NotPanics(func() {
		s.Require().False(it.HasNext())
	})

	s.Require().NotPanics(func() {
		curr, err := it.Curr()
		s.Require().Error(err)
		s.Require().Nil(curr)
	})

	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().Error(err)
		s.Require().Nil(next)
	})
}

func (s *IteratorNilSafetyTestSuite) TestBackwardIterator_NilChainOperations() {
	it := Backward(nil)

	// All operations should return errors, not panic
	s.Require().NotPanics(func() {
		s.Require().False(it.HasNext())
	})

	s.Require().NotPanics(func() {
		curr, err := it.Curr()
		s.Require().Error(err)
		s.Require().Nil(curr)
	})

	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().Error(err)
		s.Require().Nil(next)
	})
}

func (s *IteratorNilSafetyTestSuite) TestForwardIterator_AfterExhaustion() {
	node := New(1, nil, nil)
	it := Forward(node)

	// Exhaust the iterator
	_, _ = it.Next()

	// Multiple calls after exhaustion should not panic
	for i := 0; i < 10; i++ {
		s.Require().NotPanics(func() {
			next, err := it.Next()
			s.Require().Error(err)
			s.Require().ErrorIs(err, ErrEOI)
			s.Require().Nil(next)
		})
	}
}

func (s *IteratorNilSafetyTestSuite) TestBackwardIterator_AfterExhaustion() {
	node := New(1, nil, nil)
	it := Backward(node)

	// Exhaust the iterator
	_, _ = it.Next()

	// Multiple calls after exhaustion should not panic
	for i := 0; i < 10; i++ {
		s.Require().NotPanics(func() {
			next, err := it.Next()
			s.Require().Error(err)
			s.Require().ErrorIs(err, ErrEOI)
			s.Require().Nil(next)
		})
	}
}

func (s *IteratorNilSafetyTestSuite) TestForwardIterator_BrokenChain() {
	// Create chain with a nil gap: 1 -> 2 -> nil
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node1.WithNext(node2)

	it := Forward(node1)

	// First next should work
	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().NotNil(next)
	})

	// Next should hit nil - first returns (nil, nil)
	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().Nil(next)
	})

	// Third call returns ErrEOI
	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().Error(err)
		s.Require().ErrorIs(err, ErrEOI)
		s.Require().Nil(next)
	})
}

func (s *IteratorNilSafetyTestSuite) TestBackwardIterator_BrokenChain() {
	// Create chain with a nil gap: nil <- 1 <- 2
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node2.WithPrev(node1)

	it := Backward(node2)

	// First next should work
	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().NotNil(next)
	})

	// Next should hit nil - first returns (nil, nil)
	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().Nil(next)
	})

	// Third call returns ErrEOI
	s.Require().NotPanics(func() {
		next, err := it.Next()
		s.Require().Error(err)
		s.Require().ErrorIs(err, ErrEOI)
		s.Require().Nil(next)
	})
}

// IteratorEdgeCasesTestSuite tests edge cases and boundary conditions
type IteratorEdgeCasesTestSuite struct {
	suite.Suite
}

func (s *IteratorEdgeCasesTestSuite) TestForwardIterator_LargeChain() {
	// Create a chain of 1000 nodes
	const chainLength = 1000
	nodes := make([]*Node, chainLength)
	for i := 0; i < chainLength; i++ {
		nodes[i] = New(uint64(i+1), nil, nil)
	}
	for i := 0; i < chainLength-1; i++ {
		nodes[i].WithNext(nodes[i+1])
	}

	it := Forward(nodes[0])

	// Iterate through all nodes
	count := 0
	curr, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), curr.ID())
	count++

	for it.HasNext() {
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().NotNil(next)
		count++
		s.Require().Equal(uint64(count), next.ID())
	}

	s.Require().Equal(chainLength, count)
}

func (s *IteratorEdgeCasesTestSuite) TestBackwardIterator_LargeChain() {
	// Create a chain of 1000 nodes
	const chainLength = 1000
	nodes := make([]*Node, chainLength)
	for i := 0; i < chainLength; i++ {
		nodes[i] = New(uint64(i+1), nil, nil)
	}
	for i := 1; i < chainLength; i++ {
		nodes[i].WithPrev(nodes[i-1])
	}

	it := Backward(nodes[chainLength-1])

	// Iterate backward through all nodes
	count := chainLength
	curr, err := it.Curr()
	s.Require().NoError(err)
	s.Require().Equal(uint64(count), curr.ID())

	for it.HasNext() {
		count--
		next, err := it.Next()
		s.Require().NoError(err)
		s.Require().NotNil(next)
		s.Require().Equal(uint64(count), next.ID())
	}

	s.Require().Equal(1, count)
}

func (s *IteratorEdgeCasesTestSuite) TestForwardAndBackward_SameChain() {
	// Create doubly-linked chain: 1 <-> 2 <-> 3 <-> 4 <-> 5
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = New(uint64(i+1), nil, nil)
	}
	for i := 0; i < 4; i++ {
		nodes[i].WithNext(nodes[i+1])
		nodes[i+1].WithPrev(nodes[i])
	}

	// Forward iteration from start
	fwdIt := Forward(nodes[0])
	fwdIDs := []uint64{}
	curr, _ := fwdIt.Curr()
	fwdIDs = append(fwdIDs, curr.ID())
	for fwdIt.HasNext() {
		next, _ := fwdIt.Next()
		fwdIDs = append(fwdIDs, next.ID())
	}

	// Backward iteration from end
	bwdIt := Backward(nodes[4])
	bwdIDs := []uint64{}
	curr, _ = bwdIt.Curr()
	bwdIDs = append(bwdIDs, curr.ID())
	for bwdIt.HasNext() {
		next, _ := bwdIt.Next()
		bwdIDs = append(bwdIDs, next.ID())
	}

	// Verify forward iteration
	s.Require().Equal([]uint64{1, 2, 3, 4, 5}, fwdIDs)

	// Verify backward iteration (should be reverse)
	s.Require().Equal([]uint64{5, 4, 3, 2, 1}, bwdIDs)
}

func (s *IteratorEdgeCasesTestSuite) TestIterator_MaxUint64ID() {
	// Test with maximum uint64 value
	maxID := ^uint64(0) // Max uint64
	node := New(maxID, nil, nil)

	fwdIt := Forward(node)
	curr, err := fwdIt.Curr()
	s.Require().NoError(err)
	s.Require().Equal(maxID, curr.ID())

	bwdIt := Backward(node)
	curr, err = bwdIt.Curr()
	s.Require().NoError(err)
	s.Require().Equal(maxID, curr.ID())
}

func (s *IteratorEdgeCasesTestSuite) TestIterator_ZeroID() {
	// Test with zero ID
	node := New(0, nil, nil)

	fwdIt := Forward(node)
	curr, err := fwdIt.Curr()
	s.Require().NoError(err)
	s.Require().Equal(uint64(0), curr.ID())

	bwdIt := Backward(node)
	curr, err = bwdIt.Curr()
	s.Require().NoError(err)
	s.Require().Equal(uint64(0), curr.ID())
}

// Test suite runners
func TestForwardIteratorTestSuite(t *testing.T) {
	suite.Run(t, new(ForwardIteratorTestSuite))
}

func TestBackwardIteratorTestSuite(t *testing.T) {
	suite.Run(t, new(BackwardIteratorTestSuite))
}

func TestIteratorDataIntegrityTestSuite(t *testing.T) {
	suite.Run(t, new(IteratorDataIntegrityTestSuite))
}

func TestIteratorNilSafetyTestSuite(t *testing.T) {
	suite.Run(t, new(IteratorNilSafetyTestSuite))
}

func TestIteratorEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(IteratorEdgeCasesTestSuite))
}
