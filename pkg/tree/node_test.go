package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// NodeConstructorTestSuite tests the NewNode constructor
type NodeConstructorTestSuite struct {
	suite.Suite
}

func (s *NodeConstructorTestSuite) TestNewNode_WithoutChildren() {
	node := NewNode(uint64(42), "test-value", nil)

	s.Require().NotNil(node)
	s.Require().Equal(uint64(42), node.ID())
	s.Require().Equal("test-value", node.Value())
	s.Require().Equal(-1, node.Level())
	s.Require().Nil(node.Children())
	s.Require().NotNil(node.Ptr())
}

func (s *NodeConstructorTestSuite) TestNewNode_WithChildren() {
	parent := NewNode(uint64(1), "parent", nil)
	child := NewNode(uint64(2), "child", nil)
	children, err := NewNodeChildren(parent, child)
	s.Require().NoError(err)

	node := NewNode(uint64(100), "test", children)

	s.Require().NotNil(node)
	s.Require().Equal(children, node.Children())
}

func (s *NodeConstructorTestSuite) TestNewNode_DifferentTypes() {
	testCases := []struct {
		name  string
		id    uint64
		value interface{}
	}{
		{"string type", 1, "hello"},
		{"int type", 2, 42},
		{"float64 type", 3, 3.14},
		{"bool type", 4, true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			switch v := tc.value.(type) {
			case string:
				node := NewNode(tc.id, v, nil)
				s.Require().Equal(tc.id, node.ID())
				s.Require().Equal(v, node.Value())
			case int:
				node := NewNode(tc.id, v, nil)
				s.Require().Equal(tc.id, node.ID())
				s.Require().Equal(v, node.Value())
			case float64:
				node := NewNode(tc.id, v, nil)
				s.Require().Equal(tc.id, node.ID())
				s.Require().Equal(v, node.Value())
			case bool:
				node := NewNode(tc.id, v, nil)
				s.Require().Equal(tc.id, node.ID())
				s.Require().Equal(v, node.Value())
			}
		})
	}
}

func (s *NodeConstructorTestSuite) TestNewNode_ZeroAndMaxID() {
	// Test with zero ID
	zeroNode := NewNode(uint64(0), "zero", nil)
	s.Require().Equal(uint64(0), zeroNode.ID())

	// Test with max uint64 ID
	maxID := ^uint64(0)
	maxNode := NewNode(maxID, "max", nil)
	s.Require().Equal(maxID, maxNode.ID())
}

// NodeAccessorTestSuite tests accessor methods
type NodeAccessorTestSuite struct {
	suite.Suite
}

func (s *NodeAccessorTestSuite) TestValue() {
	testCases := []struct {
		name  string
		value string
	}{
		{"empty string", ""},
		{"simple string", "test"},
		{"complex string", "hello world 123!@#"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			node := NewNode(uint64(1), tc.value, nil)
			s.Require().Equal(tc.value, node.Value())
		})
	}
}

func (s *NodeAccessorTestSuite) TestLevel() {
	node := NewNode(uint64(1), "test", nil)

	// Initial level should be -1
	s.Require().Equal(-1, node.Level())
}

func (s *NodeAccessorTestSuite) TestID() {
	testCases := []struct {
		name string
		id   uint64
	}{
		{"zero ID", 0},
		{"small ID", 42},
		{"large ID", 1000000},
		{"max ID", ^uint64(0)},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			node := NewNode(tc.id, "test", nil)
			s.Require().Equal(tc.id, node.ID())
		})
	}
}

func (s *NodeAccessorTestSuite) TestPtr() {
	node := NewNode(uint64(42), "test", nil)

	ptr := node.Ptr()
	s.Require().NotNil(ptr)
	s.Require().Equal(uint64(42), ptr.ID())
}

func (s *NodeAccessorTestSuite) TestChildren_Nil() {
	node := NewNode(uint64(1), "test", nil)
	s.Require().Nil(node.Children())
}

func (s *NodeAccessorTestSuite) TestChildren_NonNil() {
	parent := NewNode(uint64(1), "parent", nil)
	child := NewNode(uint64(2), "child", nil)
	children, err := NewNodeChildren(parent, child)
	s.Require().NoError(err)

	node := NewNode(uint64(100), "test", children)
	s.Require().NotNil(node.Children())
	s.Require().Equal(children, node.Children())
}

// NodeRootTestSuite tests root-related functionality
type NodeRootTestSuite struct {
	suite.Suite
}

func (s *NodeRootTestSuite) TestIsRoot_NewNode() {
	node := NewNode(uint64(1), "test", nil)

	// New node is not a root (level is -1, not 0)
	s.Require().False(node.IsRoot())
}

func (s *NodeRootTestSuite) TestIsRoot_AfterBeholdRoot() {
	node := NewNode(uint64(1), "test", nil)
	node.BeholdRoot()

	// After BeholdRoot, it should be a root
	s.Require().True(node.IsRoot())
	s.Require().Equal(0, node.Level())
	s.Require().Nil(node.Ptr().Prev())
}

func (s *NodeRootTestSuite) TestIsRoot_WithPrevNode() {
	node1 := NewNode(uint64(1), "node1", nil)
	node2 := NewNode(uint64(2), "node2", nil)

	node1.BeholdRoot()
	node2.WithParent(node1)

	// node1 should be root
	s.Require().True(node1.IsRoot())
	s.Require().Equal(0, node1.Level())

	// node2 should not be root
	s.Require().False(node2.IsRoot())
	s.Require().Equal(1, node2.Level())
}

func (s *NodeRootTestSuite) TestBeholdRoot() {
	node := NewNode(uint64(1), "test", nil)

	// Initially not a root
	s.Require().Equal(-1, node.Level())

	node.BeholdRoot()

	// After BeholdRoot
	s.Require().Equal(0, node.Level())
	s.Require().Nil(node.Ptr().Prev())
	s.Require().True(node.IsRoot())
}

func (s *NodeRootTestSuite) TestBeholdRoot_MultipleTimes() {
	node := NewNode(uint64(1), "test", nil)

	// Call BeholdRoot multiple times
	node.BeholdRoot()
	s.Require().True(node.IsRoot())

	node.BeholdRoot()
	s.Require().True(node.IsRoot())
	s.Require().Equal(0, node.Level())
}

func (s *NodeRootTestSuite) TestBeholdRoot_AfterHavingParent() {
	parent := NewNode(uint64(1), "parent", nil)
	child := NewNode(uint64(2), "child", nil)

	parent.BeholdRoot()
	child.WithParent(parent)

	s.Require().Equal(1, child.Level())
	s.Require().NotNil(child.Ptr().Prev())

	// Make child a root
	child.BeholdRoot()

	s.Require().True(child.IsRoot())
	s.Require().Equal(0, child.Level())
	s.Require().Nil(child.Ptr().Prev())
}

// NodeParentTestSuite tests parent relationship functionality
type NodeParentTestSuite struct {
	suite.Suite
}

func (s *NodeParentTestSuite) TestWithParent_ValidParent() {
	parent := NewNode(uint64(1), "parent", nil)
	child := NewNode(uint64(2), "child", nil)

	parent.BeholdRoot()
	child.WithParent(parent)

	s.Require().Equal(1, child.Level())
	s.Require().NotNil(child.Ptr().Prev())
	s.Require().Equal(parent.Ptr(), child.Ptr().Prev())
}

func (s *NodeParentTestSuite) TestWithParent_NilParent() {
	node := NewNode(uint64(1), "test", nil)

	// WithParent with nil should not panic
	s.Require().NotPanics(func() {
		node.WithParent(nil)
	})

	// Node should remain unchanged
	s.Require().Equal(-1, node.Level())
}

func (s *NodeParentTestSuite) TestWithParent_LevelPropagation() {
	// Create a chain: root -> child1 -> child2 -> child3
	root := NewNode(uint64(0), "root", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	root.BeholdRoot()
	child1.WithParent(root)
	child2.WithParent(child1)
	child3.WithParent(child2)

	s.Require().Equal(0, root.Level())
	s.Require().Equal(1, child1.Level())
	s.Require().Equal(2, child2.Level())
	s.Require().Equal(3, child3.Level())
}

func (s *NodeParentTestSuite) TestWithParent_ChangingParent() {
	parent1 := NewNode(uint64(1), "parent1", nil)
	parent2 := NewNode(uint64(2), "parent2", nil)
	child := NewNode(uint64(3), "child", nil)

	parent1.BeholdRoot()
	parent2.BeholdRoot()

	// Initially set parent1
	child.WithParent(parent1)
	s.Require().Equal(1, child.Level())
	s.Require().Equal(parent1.Ptr(), child.Ptr().Prev())

	// Change to parent2
	child.WithParent(parent2)
	s.Require().Equal(1, child.Level())
	s.Require().Equal(parent2.Ptr(), child.Ptr().Prev())
}

func (s *NodeParentTestSuite) TestWithParent_DifferentParentLevels() {
	testCases := []struct {
		name          string
		parentLevel   int
		expectedLevel int
	}{
		{"parent at level 0", 0, 1},
		{"parent at level 1", 1, 2},
		{"parent at level 5", 5, 6},
		{"parent at level 10", 10, 11},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Create a chain to get parent to desired level
			nodes := make([]*Node[string], tc.parentLevel+2)
			for i := range nodes {
				nodes[i] = NewNode(uint64(i), "node", nil)
			}

			nodes[0].BeholdRoot()
			for i := 1; i <= tc.parentLevel; i++ {
				nodes[i].WithParent(nodes[i-1])
			}

			parent := nodes[tc.parentLevel]
			child := nodes[tc.parentLevel+1]
			child.WithParent(parent)

			s.Require().Equal(tc.parentLevel, parent.Level())
			s.Require().Equal(tc.expectedLevel, child.Level())
		})
	}
}

// NodeIntegrationTestSuite tests integration scenarios
type NodeIntegrationTestSuite struct {
	suite.Suite
}

func (s *NodeIntegrationTestSuite) TestSimpleTree() {
	// Create: root -> [child1, child2]
	root := NewNode(uint64(0), "root", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	root.BeholdRoot()
	child1.WithParent(root)
	child2.WithParent(root)

	// Verify root
	s.Require().True(root.IsRoot())
	s.Require().Equal(0, root.Level())
	s.Require().Equal("root", root.Value())

	// Verify children
	s.Require().False(child1.IsRoot())
	s.Require().Equal(1, child1.Level())
	s.Require().Equal("child1", child1.Value())

	s.Require().False(child2.IsRoot())
	s.Require().Equal(1, child2.Level())
	s.Require().Equal("child2", child2.Value())
}

func (s *NodeIntegrationTestSuite) TestDeepTree() {
	// Create a deep tree with 10 levels
	const depth = 10
	nodes := make([]*Node[int], depth)

	for i := 0; i < depth; i++ {
		nodes[i] = NewNode(uint64(i), i*100, nil)
	}

	nodes[0].BeholdRoot()
	for i := 1; i < depth; i++ {
		nodes[i].WithParent(nodes[i-1])
	}

	// Verify all levels
	for i := 0; i < depth; i++ {
		s.Require().Equal(i, nodes[i].Level())
		s.Require().Equal(i*100, nodes[i].Value())
		s.Require().Equal(uint64(i), nodes[i].ID())

		if i == 0 {
			s.Require().True(nodes[i].IsRoot())
		} else {
			s.Require().False(nodes[i].IsRoot())
		}
	}
}

func (s *NodeIntegrationTestSuite) TestTreeWithNodeChildren() {
	// Create root with NodeChildren
	root := NewNode(uint64(0), "root", nil)
	root.BeholdRoot()

	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	children, err := NewNodeChildren(root, child1, child2, child3)
	s.Require().NoError(err)

	// Create a node with these children
	nodeWithChildren := NewNode(uint64(100), "parent", children)

	s.Require().NotNil(nodeWithChildren.Children())
	s.Require().Equal(3, nodeWithChildren.Children().Size())
	s.Require().Equal("parent", nodeWithChildren.Value())
}

func (s *NodeIntegrationTestSuite) TestMultiBranchTree() {
	// Create a tree with multiple branches:
	//       root
	//      /    \
	//  branch1  branch2
	//   /  \      /  \
	// l1a l1b   l2a l2b

	root := NewNode(uint64(0), "root", nil)
	branch1 := NewNode(uint64(1), "branch1", nil)
	branch2 := NewNode(uint64(2), "branch2", nil)
	leaf1a := NewNode(uint64(11), "leaf1a", nil)
	leaf1b := NewNode(uint64(12), "leaf1b", nil)
	leaf2a := NewNode(uint64(21), "leaf2a", nil)
	leaf2b := NewNode(uint64(22), "leaf2b", nil)

	root.BeholdRoot()
	branch1.WithParent(root)
	branch2.WithParent(root)
	leaf1a.WithParent(branch1)
	leaf1b.WithParent(branch1)
	leaf2a.WithParent(branch2)
	leaf2b.WithParent(branch2)

	// Verify structure
	s.Require().Equal(0, root.Level())
	s.Require().Equal(1, branch1.Level())
	s.Require().Equal(1, branch2.Level())
	s.Require().Equal(2, leaf1a.Level())
	s.Require().Equal(2, leaf1b.Level())
	s.Require().Equal(2, leaf2a.Level())
	s.Require().Equal(2, leaf2b.Level())

	// Verify only root is root
	s.Require().True(root.IsRoot())
	s.Require().False(branch1.IsRoot())
	s.Require().False(branch2.IsRoot())
	s.Require().False(leaf1a.IsRoot())
	s.Require().False(leaf1b.IsRoot())
	s.Require().False(leaf2a.IsRoot())
	s.Require().False(leaf2b.IsRoot())
}

func (s *NodeIntegrationTestSuite) TestPtrConsistency() {
	parent := NewNode(uint64(1), "parent", nil)
	child := NewNode(uint64(2), "child", nil)

	parent.BeholdRoot()
	child.WithParent(parent)

	// Verify Ptr() returns consistent pointer
	ptr1 := parent.Ptr()
	ptr2 := parent.Ptr()
	s.Require().Equal(ptr1, ptr2)
	s.Require().Equal(uint64(1), ptr1.ID())
	s.Require().Equal(uint64(1), ptr2.ID())

	// Verify parent-child ptr relationship
	s.Require().Equal(parent.Ptr(), child.Ptr().Prev())
}

func (s *NodeIntegrationTestSuite) TestComplexValueTypes() {
	// Test with struct as value
	type CustomStruct struct {
		Name string
		Age  int
	}

	value := CustomStruct{Name: "Alice", Age: 30}
	node := NewNode(uint64(1), value, nil)

	s.Require().Equal("Alice", node.Value().Name)
	s.Require().Equal(30, node.Value().Age)
}

// Test suite runners
func TestNodeConstructorTestSuite(t *testing.T) {
	suite.Run(t, new(NodeConstructorTestSuite))
}

func TestNodeAccessorTestSuite(t *testing.T) {
	suite.Run(t, new(NodeAccessorTestSuite))
}

func TestNodeRootTestSuite(t *testing.T) {
	suite.Run(t, new(NodeRootTestSuite))
}

func TestNodeParentTestSuite(t *testing.T) {
	suite.Run(t, new(NodeParentTestSuite))
}

func TestNodeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(NodeIntegrationTestSuite))
}
