package node

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// IterSeqTestSuite tests Go 1.23+ iterator functionality
type IterSeqTestSuite struct {
	suite.Suite
}

func (s *IterSeqTestSuite) TestNextNodes_NilNode() {
	var collected []int

	for id, node := range NextNodes(nil) {
		collected = append(collected, id)
		s.Require().Nil(node)
	}

	// Should not iterate over nil
	s.Require().Empty(collected)
}

func (s *IterSeqTestSuite) TestNextNodes_SingleNode() {
	node := New(1, nil, nil)
	var collected []int
	var nodes []*Node

	for id, n := range NextNodes(node) {
		collected = append(collected, id)
		nodes = append(nodes, n)
	}

	// Single node should yield itself (this was the bug - it was being skipped)
	s.Require().Equal([]int{0}, collected, "Should yield index 0")
	s.Require().Len(nodes, 1, "Should yield exactly one node")
	s.Require().Equal(node, nodes[0], "Should yield the starting node")
	s.Require().Equal(uint64(1), nodes[0].ID(), "CreateNode should have correct ID")
}

func (s *IterSeqTestSuite) TestNextNodes_LinearChain() {
	// Create chain: 1 -> 2 -> 3 -> 4 -> 5
	nodeList := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodeList[i] = New(uint64(i+1), nil, nil)
	}
	for i := 0; i < 4; i++ {
		nodeList[i].WithNext(nodeList[i+1])
	}

	var collected []int
	var nodes []*Node

	for id, node := range NextNodes(nodeList[0]) {
		collected = append(collected, id)
		nodes = append(nodes, node)
	}

	// Should collect ALL nodes: 1, 2, 3, 4, 5 (starting from node 1)
	s.Require().Equal([]int{0, 1, 2, 3, 4}, collected, "Should have indices 0-4")
	s.Require().Len(nodes, 5, "Should yield all 5 nodes")
	s.Require().Equal(nodeList[0], nodes[0], "First should be node 1")
	s.Require().Equal(nodeList[1], nodes[1], "Second should be node 2")
	s.Require().Equal(nodeList[2], nodes[2], "Third should be node 3")
	s.Require().Equal(nodeList[3], nodes[3], "Fourth should be node 4")
	s.Require().Equal(nodeList[4], nodes[4], "Fifth should be node 5")
}

func (s *IterSeqTestSuite) TestNextNodes_EarlyBreak() {
	// Create chain: 1 -> 2 -> 3 -> 4 -> 5
	nodeList := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodeList[i] = New(uint64(i+1), nil, nil)
	}
	for i := 0; i < 4; i++ {
		nodeList[i].WithNext(nodeList[i+1])
	}

	var collected []int
	count := 0

	for id, node := range NextNodes(nodeList[0]) {
		collected = append(collected, id)
		s.Require().NotNil(node)
		count++
		if count >= 2 {
			break
		}
	}

	// Should only collect 2 nodes due to early break
	s.Require().Len(collected, 2)
}

func (s *IterSeqTestSuite) TestPrevNodes_NilNode() {
	var collected []int

	for id, node := range PrevNodes(nil) {
		collected = append(collected, id)
		s.Require().Nil(node)
	}

	// Should not iterate over nil
	s.Require().Empty(collected)
}

func (s *IterSeqTestSuite) TestPrevNodes_SingleNode() {
	node := New(1, nil, nil)
	var collected []int
	var nodes []*Node

	for id, n := range PrevNodes(node) {
		collected = append(collected, id)
		nodes = append(nodes, n)
	}

	// Single node should yield itself (this was the bug - it was being skipped)
	s.Require().Equal([]int{0}, collected, "Should yield index 0")
	s.Require().Len(nodes, 1, "Should yield exactly one node")
	s.Require().Equal(node, nodes[0], "Should yield the starting node")
	s.Require().Equal(uint64(1), nodes[0].ID(), "CreateNode should have correct ID")
}

func (s *IterSeqTestSuite) TestPrevNodes_LinearChain() {
	// Create chain: 1 <- 2 <- 3 <- 4 <- 5
	nodeList := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodeList[i] = New(uint64(i+1), nil, nil)
	}
	for i := 1; i < 5; i++ {
		nodeList[i].WithPrev(nodeList[i-1])
	}

	var collected []int
	var nodes []*Node

	for id, node := range PrevNodes(nodeList[4]) {
		collected = append(collected, id)
		nodes = append(nodes, node)
	}

	// Should collect ALL nodes: 5, 4, 3, 2, 1 (starting from node 5, going backward)
	s.Require().Equal([]int{0, 1, 2, 3, 4}, collected, "Should have indices 0-4")
	s.Require().Len(nodes, 5, "Should yield all 5 nodes")
	s.Require().Equal(nodeList[4], nodes[0], "First should be node 5")
	s.Require().Equal(nodeList[3], nodes[1], "Second should be node 4")
	s.Require().Equal(nodeList[2], nodes[2], "Third should be node 3")
	s.Require().Equal(nodeList[1], nodes[3], "Fourth should be node 2")
	s.Require().Equal(nodeList[0], nodes[4], "Fifth should be node 1")
}

func (s *IterSeqTestSuite) TestPrevNodes_EarlyBreak() {
	// Create chain: 1 <- 2 <- 3 <- 4 <- 5
	nodeList := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodeList[i] = New(uint64(i+1), nil, nil)
	}
	for i := 1; i < 5; i++ {
		nodeList[i].WithPrev(nodeList[i-1])
	}

	var collected []int
	count := 0

	for id, node := range PrevNodes(nodeList[4]) {
		collected = append(collected, id)
		s.Require().NotNil(node)
		count++
		if count >= 2 {
			break
		}
	}

	// Should only collect 2 nodes due to early break
	s.Require().Len(collected, 2)
}

func (s *IterSeqTestSuite) TestNextNodes_DoublyLinkedChain() {
	// Create doubly-linked chain: 1 <-> 2 <-> 3
	node1 := New(100, nil, nil)
	node2 := New(200, nil, nil)
	node3 := New(300, nil, nil)

	node1.WithNext(node2)
	node2.WithPrev(node1)
	node2.WithNext(node3)
	node3.WithPrev(node2)

	var nodeIDs []uint64

	for _, node := range NextNodes(node1) {
		nodeIDs = append(nodeIDs, node.ID())
	}

	// Should include starting node and all following nodes
	s.Require().Equal([]uint64{100, 200, 300}, nodeIDs, "Should include all nodes starting from node1")
}

func (s *IterSeqTestSuite) TestPrevNodes_DoublyLinkedChain() {
	// Create doubly-linked chain: 1 <-> 2 <-> 3
	node1 := New(100, nil, nil)
	node2 := New(200, nil, nil)
	node3 := New(300, nil, nil)

	node1.WithNext(node2)
	node2.WithPrev(node1)
	node2.WithNext(node3)
	node3.WithPrev(node2)

	var nodeIDs []uint64

	for _, node := range PrevNodes(node3) {
		nodeIDs = append(nodeIDs, node.ID())
	}

	// Should include starting node and all previous nodes
	s.Require().Equal([]uint64{300, 200, 100}, nodeIDs, "Should include all nodes starting from node3 going backward")
}

// NodeConstructorTestSuite tests node constructor functions
type NodeConstructorTestSuite struct {
	suite.Suite
}

func (s *NodeConstructorTestSuite) TestID_CreatesNode() {
	node := ID(42)

	s.Require().NotNil(node)
	s.Require().Equal(uint64(42), node.ID())
	s.Require().Nil(node.Next())
	s.Require().Nil(node.Prev())
}

func (s *NodeConstructorTestSuite) TestID_ZeroID() {
	node := ID(0)

	s.Require().NotNil(node)
	s.Require().Equal(uint64(0), node.ID())
}

func (s *NodeConstructorTestSuite) TestID_MaxUint64() {
	maxID := ^uint64(0)
	node := ID(maxID)

	s.Require().NotNil(node)
	s.Require().Equal(maxID, node.ID())
}

func (s *NodeConstructorTestSuite) TestID_MultipleNodes() {
	node1 := ID(1)
	node2 := ID(2)
	node3 := ID(3)

	// Verify each node has correct ID
	s.Require().Equal(uint64(1), node1.ID())
	s.Require().Equal(uint64(2), node2.ID())
	s.Require().Equal(uint64(3), node3.ID())

	// Link them together
	node1.WithNext(node2)
	node2.WithPrev(node1)
	node2.WithNext(node3)
	node3.WithPrev(node2)

	// Verify links work correctly
	s.Require().Equal(node2, node1.Next())
	s.Require().Equal(node1, node2.Prev())
	s.Require().Equal(node3, node2.Next())
	s.Require().Equal(node2, node3.Prev())
}

// Test suite runners
func TestIterSeqTestSuite(t *testing.T) {
	suite.Run(t, new(IterSeqTestSuite))
}

func TestNodeConstructorTestSuite(t *testing.T) {
	suite.Run(t, new(NodeConstructorTestSuite))
}
