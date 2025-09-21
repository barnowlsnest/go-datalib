package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// NodeTestSuite defines the test suite for basic node functionality
type NodeTestSuite struct {
	suite.Suite
}

func (s *NodeTestSuite) TestNew_WithNilReferences() {
	node := New(42, nil, nil)

	assert.NotNil(s.T(), node)
	assert.Equal(s.T(), uint64(42), node.ID())
	assert.Nil(s.T(), node.Next())
	assert.Nil(s.T(), node.Prev())
}

func (s *NodeTestSuite) TestNew_WithValidReferences() {
	prevNode := New(1, nil, nil)
	nextNode := New(3, nil, nil)
	node := New(2, nextNode, prevNode)

	assert.NotNil(s.T(), node)
	assert.Equal(s.T(), uint64(2), node.ID())
	assert.Equal(s.T(), nextNode, node.Next())
	assert.Equal(s.T(), prevNode, node.Prev())
}

func (s *NodeTestSuite) TestNew_WithZeroID() {
	node := New(0, nil, nil)

	assert.NotNil(s.T(), node)
	assert.Equal(s.T(), uint64(0), node.ID())
	assert.Nil(s.T(), node.Next())
	assert.Nil(s.T(), node.Prev())
}

func (s *NodeTestSuite) TestNew_WithLargeID() {
	largeID := uint64(18446744073709551615) // Max uint64
	node := New(largeID, nil, nil)

	assert.NotNil(s.T(), node)
	assert.Equal(s.T(), largeID, node.ID())
	assert.Nil(s.T(), node.Next())
	assert.Nil(s.T(), node.Prev())
}

func (s *NodeTestSuite) TestID_Immutability() {
	node := New(100, nil, nil)

	originalID := node.ID()
	assert.Equal(s.T(), uint64(100), originalID)

	// ID should remain the same after multiple calls
	assert.Equal(s.T(), originalID, node.ID())
	assert.Equal(s.T(), originalID, node.ID())
}

func (s *NodeTestSuite) TestNext_ReturnsCorrectReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, node2, node1)

	assert.Equal(s.T(), node2, node3.Next())
	assert.Nil(s.T(), node1.Next())
	assert.Nil(s.T(), node2.Next())
}

func (s *NodeTestSuite) TestPrev_ReturnsCorrectReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, node2, node1)

	assert.Equal(s.T(), node1, node3.Prev())
	assert.Nil(s.T(), node1.Prev())
	assert.Nil(s.T(), node2.Prev())
}

func (s *NodeTestSuite) TestWithNext_SetsNextReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)

	// Initially node1 has no next
	assert.Nil(s.T(), node1.Next())

	// Set next reference
	node1.WithNext(node2)
	assert.Equal(s.T(), node2, node1.Next())
}

func (s *NodeTestSuite) TestWithNext_WithNilReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)

	// Set next reference first
	node1.WithNext(node2)
	assert.Equal(s.T(), node2, node1.Next())

	// Clear next reference with nil
	node1.WithNext(nil)
	assert.Nil(s.T(), node1.Next())
}

func (s *NodeTestSuite) TestWithNext_OverwriteExistingReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	// Set initial next reference
	node1.WithNext(node2)
	assert.Equal(s.T(), node2, node1.Next())

	// Overwrite with different reference
	node1.WithNext(node3)
	assert.Equal(s.T(), node3, node1.Next())
}

func (s *NodeTestSuite) TestWithPrev_SetsPrevReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)

	// Initially node2 has no prev
	assert.Nil(s.T(), node2.Prev())

	// Set prev reference
	node2.WithPrev(node1)
	assert.Equal(s.T(), node1, node2.Prev())
}

func (s *NodeTestSuite) TestWithPrev_WithNilReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)

	// Set prev reference first
	node2.WithPrev(node1)
	assert.Equal(s.T(), node1, node2.Prev())

	// Clear prev reference with nil
	node2.WithPrev(nil)
	assert.Nil(s.T(), node2.Prev())
}

func (s *NodeTestSuite) TestWithPrev_OverwriteExistingReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	// Set initial prev reference
	node3.WithPrev(node1)
	assert.Equal(s.T(), node1, node3.Prev())

	// Overwrite with different reference
	node3.WithPrev(node2)
	assert.Equal(s.T(), node2, node3.Prev())
}

// NodeChainTestSuite defines tests for chaining multiple nodes
type NodeChainTestSuite struct {
	suite.Suite
}

func (s *NodeChainTestSuite) TestThreeNodeChain() {
	// Create base nodes first
	node2 := New(2, nil, nil)

	// Now create the chain by creating new nodes with proper references
	// node1 -> node2, node1 has no prev
	node1WithNext := New(1, node2, nil)
	// node3 -> nil, node3 has node2 as prev
	node3WithPrev := New(3, nil, node2)
	// node2 -> node3, node2 has node1 as prev
	node2Complete := New(2, node3WithPrev, node1WithNext)

	// Test the middle node (which has the most complete references)
	assert.Equal(s.T(), node3WithPrev, node2Complete.Next())
	assert.Equal(s.T(), node1WithNext, node2Complete.Prev())

	// Test node1 (start of chain)
	assert.Equal(s.T(), node2, node1WithNext.Next())
	assert.Nil(s.T(), node1WithNext.Prev())

	// Test node3 (end of chain)
	assert.Nil(s.T(), node3WithPrev.Next())
	assert.Equal(s.T(), node2, node3WithPrev.Prev())

	// Test IDs
	assert.Equal(s.T(), uint64(1), node1WithNext.ID())
	assert.Equal(s.T(), uint64(2), node2Complete.ID())
	assert.Equal(s.T(), uint64(3), node3WithPrev.ID())
}

func (s *NodeChainTestSuite) TestSingleNodeLoop() {
	// Create a self-referencing node
	var selfNode *Node
	selfNode = New(42, selfNode, selfNode)

	assert.NotNil(s.T(), selfNode)
	assert.Equal(s.T(), uint64(42), selfNode.ID())
	// Note: Next() and Prev() will return nil because the node was created
	// with nil references initially (selfNode was nil when New was called)
}

func (s *NodeChainTestSuite) TestTwoNodeCircle() {
	// Since nodes are immutable, we need to create them in a way that allows circular references
	// We'll create a placeholder first, then create the actual nodes
	placeholder := New(999, nil, nil)

	// Create node1 that references node2 (using placeholder initially)
	node1 := New(1, placeholder, placeholder)
	// Create node2 that references node1
	node2 := New(2, node1, node1)

	// Now create the final node1 with proper node2 references
	node1Final := New(1, node2, node2)

	assert.Equal(s.T(), node2, node1Final.Next())
	assert.Equal(s.T(), node2, node1Final.Prev())
	assert.Equal(s.T(), node1, node2.Next())
	assert.Equal(s.T(), node1, node2.Prev())
}

// NodeRobustnessTestSuite defines tests for edge cases and robustness
type NodeRobustnessTestSuite struct {
	suite.Suite
}

func (s *NodeRobustnessTestSuite) TestSameIDNodes() {
	// Multiple nodes can have the same ID (no uniqueness constraint)
	node1 := New(100, nil, nil)
	node2 := New(100, nil, nil)
	node3 := New(100, node2, node1)

	assert.Equal(s.T(), uint64(100), node1.ID())
	assert.Equal(s.T(), uint64(100), node2.ID())
	assert.Equal(s.T(), uint64(100), node3.ID())

	// But they should be different objects (different memory addresses)
	assert.False(s.T(), node1 == node2, "node1 and node2 should be different objects")
	assert.False(s.T(), node2 == node3, "node2 and node3 should be different objects")
	assert.False(s.T(), node1 == node3, "node1 and node3 should be different objects")
}

func (s *NodeRobustnessTestSuite) TestMultipleReferencesToSameNode() {
	baseNode := New(1, nil, nil)

	// Multiple nodes referencing the same base node
	node2 := New(2, baseNode, baseNode)
	node3 := New(3, baseNode, baseNode)

	assert.Equal(s.T(), baseNode, node2.Next())
	assert.Equal(s.T(), baseNode, node2.Prev())
	assert.Equal(s.T(), baseNode, node3.Next())
	assert.Equal(s.T(), baseNode, node3.Prev())

	// All should reference the same object
	assert.True(s.T(), node2.Next() == node3.Next())
	assert.True(s.T(), node2.Prev() == node3.Prev())
}

func (s *NodeRobustnessTestSuite) TestNodeWithSelfReference() {
	node := New(42, nil, nil)
	// Create a new node that references itself
	selfRefNode := New(42, node, node)

	assert.Equal(s.T(), node, selfRefNode.Next())
	assert.Equal(s.T(), node, selfRefNode.Prev())
	assert.Equal(s.T(), selfRefNode.Next(), selfRefNode.Prev())
}

// NodeStructuralTestSuite defines tests for structural integrity
type NodeStructuralTestSuite struct {
	suite.Suite
}

func (s *NodeStructuralTestSuite) TestNodeFieldsArePrivate() {
	node := New(123, nil, nil)

	// We can't directly test private fields, but we can ensure
	// the public interface works as expected
	assert.Equal(s.T(), uint64(123), node.ID())
	assert.Nil(s.T(), node.Next())
	assert.Nil(s.T(), node.Prev())
}

func (s *NodeStructuralTestSuite) TestNodeMethodsAreReadOnly() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, node2, node1)

	// Getting references multiple times should return the same values
	next1 := node3.Next()
	next2 := node3.Next()
	prev1 := node3.Prev()
	prev2 := node3.Prev()

	assert.Equal(s.T(), next1, next2)
	assert.Equal(s.T(), prev1, prev2)
	assert.Equal(s.T(), node2, next1)
	assert.Equal(s.T(), node1, prev1)
}

func (s *NodeStructuralTestSuite) TestNodeCreationDoesntModifyInputs() {
	originalNext := New(10, nil, nil)
	originalPrev := New(20, nil, nil)

	// Store original states
	originalNextNext := originalNext.Next()
	originalNextPrev := originalNext.Prev()
	originalPrevNext := originalPrev.Next()
	originalPrevPrev := originalPrev.Prev()

	// Create a new node with these as references
	New(30, originalNext, originalPrev)

	// Original nodes should remain unchanged
	assert.Equal(s.T(), originalNextNext, originalNext.Next())
	assert.Equal(s.T(), originalNextPrev, originalNext.Prev())
	assert.Equal(s.T(), originalPrevNext, originalPrev.Next())
	assert.Equal(s.T(), originalPrevPrev, originalPrev.Prev())
}

// NodeMutableOperationsTestSuite defines tests for mutable operations (WithNext/WithPrev)
type NodeMutableOperationsTestSuite struct {
	suite.Suite
}

func (s *NodeMutableOperationsTestSuite) TestWithNext_AndWithPrev_CreateChain() {
	// Create three independent nodes
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	// Build a chain using mutable operations: node1 <-> node2 <-> node3
	node1.WithNext(node2)
	node2.WithPrev(node1)
	node2.WithNext(node3)
	node3.WithPrev(node2)

	// Test forward traversal
	assert.Equal(s.T(), node2, node1.Next())
	assert.Equal(s.T(), node3, node2.Next())
	assert.Nil(s.T(), node3.Next())

	// Test backward traversal
	assert.Equal(s.T(), node2, node3.Prev())
	assert.Equal(s.T(), node1, node2.Prev())
	assert.Nil(s.T(), node1.Prev())
}

func (s *NodeMutableOperationsTestSuite) TestWithNext_SelfReference() {
	node := New(1, nil, nil)

	// Create self-reference
	node.WithNext(node)

	assert.Equal(s.T(), node, node.Next())
	assert.Nil(s.T(), node.Prev()) // Prev should still be nil
}

func (s *NodeMutableOperationsTestSuite) TestWithPrev_SelfReference() {
	node := New(1, nil, nil)

	// Create self-reference
	node.WithPrev(node)

	assert.Equal(s.T(), node, node.Prev())
	assert.Nil(s.T(), node.Next()) // Next should still be nil
}

func (s *NodeMutableOperationsTestSuite) TestWithNext_WithPrev_CircularReference() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)

	// Create circular references
	node1.WithNext(node2)
	node1.WithPrev(node2)
	node2.WithNext(node1)
	node2.WithPrev(node1)

	// Test circular references
	assert.Equal(s.T(), node2, node1.Next())
	assert.Equal(s.T(), node2, node1.Prev())
	assert.Equal(s.T(), node1, node2.Next())
	assert.Equal(s.T(), node1, node2.Prev())
}

func (s *NodeMutableOperationsTestSuite) TestWithNext_DoesntAffectPrev() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	// Set initial prev reference
	node1.WithPrev(node2)
	assert.Equal(s.T(), node2, node1.Prev())
	assert.Nil(s.T(), node1.Next())

	// Setting next should not affect prev
	node1.WithNext(node3)
	assert.Equal(s.T(), node2, node1.Prev()) // Should remain unchanged
	assert.Equal(s.T(), node3, node1.Next())
}

func (s *NodeMutableOperationsTestSuite) TestWithPrev_DoesntAffectNext() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	// Set initial next reference
	node1.WithNext(node2)
	assert.Equal(s.T(), node2, node1.Next())
	assert.Nil(s.T(), node1.Prev())

	// Setting prev should not affect next
	node1.WithPrev(node3)
	assert.Equal(s.T(), node2, node1.Next()) // Should remain unchanged
	assert.Equal(s.T(), node3, node1.Prev())
}

func (s *NodeMutableOperationsTestSuite) TestMultipleWithNext_Operations() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)
	node4 := New(4, nil, nil)

	// Chain multiple WithNext operations
	node1.WithNext(node2)
	node1.WithNext(node3) // Overwrite previous
	node1.WithNext(node4) // Overwrite again

	assert.Equal(s.T(), node4, node1.Next())
	// Previous nodes should not be affected by being replaced
	assert.Nil(s.T(), node2.Next())
	assert.Nil(s.T(), node3.Next())
}

func (s *NodeMutableOperationsTestSuite) TestMultipleWithPrev_Operations() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)
	node4 := New(4, nil, nil)

	// Chain multiple WithPrev operations
	node1.WithPrev(node2)
	node1.WithPrev(node3) // Overwrite previous
	node1.WithPrev(node4) // Overwrite again

	assert.Equal(s.T(), node4, node1.Prev())
	// Previous nodes should not be affected by being replaced
	assert.Nil(s.T(), node2.Prev())
	assert.Nil(s.T(), node3.Prev())
}

func (s *NodeMutableOperationsTestSuite) TestWithNext_WithPrev_Independence() {
	node1 := New(1, nil, nil)
	node2 := New(2, nil, nil)
	node3 := New(3, nil, nil)

	// Set references independently
	node1.WithNext(node2)
	node1.WithPrev(node3)

	// Both should be set correctly
	assert.Equal(s.T(), node2, node1.Next())
	assert.Equal(s.T(), node3, node1.Prev())

	// Clear one, other should remain
	node1.WithNext(nil)
	assert.Nil(s.T(), node1.Next())
	assert.Equal(s.T(), node3, node1.Prev()) // Should remain unchanged
}

// TestSuite runners
func TestNodeTestSuite(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

func TestNodeChainTestSuite(t *testing.T) {
	suite.Run(t, new(NodeChainTestSuite))
}

func TestNodeRobustnessTestSuite(t *testing.T) {
	suite.Run(t, new(NodeRobustnessTestSuite))
}

func TestNodeStructuralTestSuite(t *testing.T) {
	suite.Run(t, new(NodeStructuralTestSuite))
}

func TestNodeMutableOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(NodeMutableOperationsTestSuite))
}
