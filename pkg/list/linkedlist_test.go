package list

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/barnowlsnest/go-datalib/pkg/node"
)

// LinkedListBasicTestSuite defines tests for basic linked list functionality
type LinkedListBasicTestSuite struct {
	suite.Suite
}

func (s *LinkedListBasicTestSuite) TestNewLinkedList_ShouldCreateEmptyList() {
	list := New()

	s.Require().NotNil(list)
	s.Require().Nil(list.head)
	s.Require().Nil(list.tail)
	s.Require().Equal(0, list.size)
}

func (s *LinkedListBasicTestSuite) TestPush_ToEmptyList() {
	list := New()
	n := node.New(1, nil, nil)

	list.Push(n)

	s.Require().Equal(1, list.size)
	s.Require().Equal(n, list.head)
	s.Require().Equal(n, list.tail)
	s.Require().Nil(n.Next())
	s.Require().Nil(n.Prev())
}

func (s *LinkedListBasicTestSuite) TestPush_ToNonEmptyList() {
	list := New()
	node1 := node.New(1, nil, nil)
	node2 := node.New(2, nil, nil)

	list.Push(node1)
	list.Push(node2)

	s.Require().Equal(2, list.size)
	s.Require().Equal(node1, list.head)
	s.Require().Equal(node2, list.tail)
	s.Require().Equal(node2, node1.Next())
	s.Require().Equal(node1, node2.Prev())
	s.Require().Nil(node2.Next())
}

func (s *LinkedListBasicTestSuite) TestPop_FromEmptyList() {
	list := New()

	result := list.Pop()

	s.Require().Nil(result)
	s.Require().Equal(0, list.size)
}

func (s *LinkedListBasicTestSuite) TestPop_FromListWithOneElement() {
	list := New()
	n := node.New(1, nil, nil)
	list.Push(n)

	result := list.Pop()

	s.Require().NotNil(result)
	s.Require().Equal(uint64(1), result.ID())
	s.Require().Equal(0, list.size)
	s.Require().Nil(list.head)
	s.Require().Nil(list.tail)
}

func (s *LinkedListBasicTestSuite) TestPop_FromListWithMultipleElements() {
	list := New()
	node1 := node.New(1, nil, nil)
	node2 := node.New(2, nil, nil)
	list.Push(node1)
	list.Push(node2)

	result := list.Pop()

	s.Require().NotNil(result)
	s.Require().Equal(uint64(2), result.ID())
	s.Require().Equal(1, list.size)
	s.Require().Equal(node1, list.head)
	s.Require().Equal(node1, list.tail)
	s.Require().Nil(node1.Next())
}

func (s *LinkedListBasicTestSuite) TestUnshift_ToEmptyList() {
	list := New()
	n := node.New(1, nil, nil)

	list.Unshift(n)

	s.Require().Equal(1, list.size)
	s.Require().Equal(n, list.head)
	s.Require().Equal(n, list.tail)
	s.Require().Nil(n.Next())
	s.Require().Nil(n.Prev())
}

func (s *LinkedListBasicTestSuite) TestUnshift_ToNonEmptyList() {
	list := New()
	node1 := node.New(1, nil, nil)
	node2 := node.New(2, nil, nil)

	list.Push(node1)
	list.Unshift(node2)

	s.Require().Equal(2, list.size)
	s.Require().Equal(node2, list.head)
	s.Require().Equal(node1, list.tail)
	s.Require().Equal(node1, node2.Next())
	s.Require().Equal(node2, node1.Prev())
	s.Require().Nil(node1.Next())
}

func (s *LinkedListBasicTestSuite) TestShift_FromEmptyList() {
	list := New()

	result := list.Shift()

	s.Require().Nil(result)
	s.Require().Equal(0, list.size)
}

func (s *LinkedListBasicTestSuite) TestShift_FromListWithOneElement() {
	list := New()
	n := node.New(1, nil, nil)
	list.Push(n)

	result := list.Shift()

	s.Require().NotNil(result)
	s.Require().Equal(uint64(1), result.ID())
	s.Require().Equal(0, list.size)
	s.Require().Nil(list.head)
	s.Require().Nil(list.tail)
}

func (s *LinkedListBasicTestSuite) TestShift_FromListWithMultipleElements() {
	list := New()
	node1 := node.New(1, nil, nil)
	node2 := node.New(2, nil, nil)
	list.Push(node1)
	list.Push(node2)

	result := list.Shift()

	s.Require().NotNil(result)
	s.Require().Equal(uint64(1), result.ID())
	s.Require().Equal(1, list.size)
	s.Require().Equal(node2, list.head)
	s.Require().Equal(node2, list.tail)
	s.Require().Nil(node2.Prev())
}

// LinkedListCombinedOperationsTestSuite defines tests for combined operations
type LinkedListCombinedOperationsTestSuite struct {
	suite.Suite
}

func (s *LinkedListCombinedOperationsTestSuite) TestPushPopUnshiftShiftInSequence() {
	list := New()

	// Push two nodes
	node1 := node.New(1, nil, nil)
	node2 := node.New(2, nil, nil)
	list.Push(node1)
	list.Push(node2)

	s.Require().Equal(2, list.size)

	// Pop one node
	popped := list.Pop()
	s.Require().Equal(uint64(2), popped.ID())
	s.Require().Equal(1, list.size)

	// Unshift a new node
	node3 := node.New(3, nil, nil)
	list.Unshift(node3)
	s.Require().Equal(2, list.size)
	s.Require().Equal(node3, list.head)
	s.Require().Equal(node1, list.tail)

	// Shift a node
	shifted := list.Shift()
	s.Require().Equal(uint64(3), shifted.ID())
	s.Require().Equal(1, list.size)
	s.Require().Equal(node1, list.head)
	s.Require().Equal(node1, list.tail)
}

// LinkedListIDBasedTestSuite defines tests for ID-based operations
type LinkedListIDBasedTestSuite struct {
	suite.Suite
}

func (s *LinkedListIDBasedTestSuite) TestPushID_ToEmptyList() {
	list := New()

	list.PushID(42)

	s.Require().Equal(1, list.size)
	id, err := list.HeadID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(42), id)
}

func (s *LinkedListIDBasedTestSuite) TestPushID_MultipleNodes() {
	list := New()

	list.PushID(1)
	list.PushID(2)
	list.PushID(3)

	s.Require().Equal(3, list.size)

	headID, err := list.HeadID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), headID)

	tailID, err := list.TailID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(3), tailID)
}

func (s *LinkedListIDBasedTestSuite) TestPopID_FromEmptyList() {
	list := New()

	id, err := list.PopID()

	s.Require().Error(err)
	s.Require().Equal(node.ErrNil, err)
	s.Require().Equal(uint64(0), id)
}

func (s *LinkedListIDBasedTestSuite) TestPopID_FromList() {
	list := New()
	list.PushID(1)
	list.PushID(2)
	list.PushID(3)

	id, err := list.PopID()

	s.Require().NoError(err)
	s.Require().Equal(uint64(3), id)
	s.Require().Equal(2, list.size)
}

func (s *LinkedListIDBasedTestSuite) TestUnshiftID_ToEmptyList() {
	list := New()

	list.UnshiftID(42)

	s.Require().Equal(1, list.size)
	id, err := list.HeadID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(42), id)
}

func (s *LinkedListIDBasedTestSuite) TestUnshiftID_MultipleNodes() {
	list := New()

	list.UnshiftID(1)
	list.UnshiftID(2)
	list.UnshiftID(3)

	s.Require().Equal(3, list.size)

	headID, err := list.HeadID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(3), headID)

	tailID, err := list.TailID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), tailID)
}

func (s *LinkedListIDBasedTestSuite) TestShiftID_FromEmptyList() {
	list := New()

	id, err := list.ShiftID()

	s.Require().Error(err)
	s.Require().Equal(node.ErrNil, err)
	s.Require().Equal(uint64(0), id)
}

func (s *LinkedListIDBasedTestSuite) TestShiftID_FromList() {
	list := New()
	list.PushID(1)
	list.PushID(2)
	list.PushID(3)

	id, err := list.ShiftID()

	s.Require().NoError(err)
	s.Require().Equal(uint64(1), id)
	s.Require().Equal(2, list.size)
}

func (s *LinkedListIDBasedTestSuite) TestHeadID_EmptyList() {
	list := New()

	id, err := list.HeadID()

	s.Require().Error(err)
	s.Require().Equal(node.ErrNil, err)
	s.Require().Equal(uint64(0), id)
}

func (s *LinkedListIDBasedTestSuite) TestHeadID_NonEmptyList() {
	list := New()
	list.PushID(1)
	list.PushID(2)

	id, err := list.HeadID()

	s.Require().NoError(err)
	s.Require().Equal(uint64(1), id)
}

func (s *LinkedListIDBasedTestSuite) TestTailID_EmptyList() {
	list := New()

	id, err := list.TailID()

	s.Require().Error(err)
	s.Require().Equal(node.ErrNil, err)
	s.Require().Equal(uint64(0), id)
}

func (s *LinkedListIDBasedTestSuite) TestTailID_NonEmptyList() {
	list := New()
	list.PushID(1)
	list.PushID(2)

	id, err := list.TailID()

	s.Require().NoError(err)
	s.Require().Equal(uint64(2), id)
}

// LinkedListIDWorkflowTestSuite defines tests for complete ID-based workflows
type LinkedListIDWorkflowTestSuite struct {
	suite.Suite
}

func (s *LinkedListIDWorkflowTestSuite) TestIDBasedMethodsExclusively() {
	list := New()

	// Add items using PushID and UnshiftID
	list.PushID(2)
	list.PushID(3)
	list.UnshiftID(1)
	list.PushID(4)

	// Verify size
	s.Require().Equal(4, list.Size())

	// Verify head and tail
	headID, err := list.HeadID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), headID)

	tailID, err := list.TailID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(4), tailID)

	// Remove from both ends
	shiftedID, err := list.ShiftID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), shiftedID)

	poppedID, err := list.PopID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(4), poppedID)

	s.Require().Equal(2, list.Size())

	// Verify remaining elements
	headID, err = list.HeadID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(2), headID)

	tailID, err = list.TailID()
	s.Require().NoError(err)
	s.Require().Equal(uint64(3), tailID)
}

// TestSuite runners
func TestLinkedListBasicTestSuite(t *testing.T) {
	suite.Run(t, new(LinkedListBasicTestSuite))
}

func TestLinkedListCombinedOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(LinkedListCombinedOperationsTestSuite))
}

func TestLinkedListIDBasedTestSuite(t *testing.T) {
	suite.Run(t, new(LinkedListIDBasedTestSuite))
}

func TestLinkedListIDWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(LinkedListIDWorkflowTestSuite))
}
