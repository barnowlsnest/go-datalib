package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// NodeChildrenTestSuite tests the NodeChildren structure and its methods
type NodeChildrenTestSuite struct {
	suite.Suite
}

func (s *NodeChildrenTestSuite) TestNewNodeChildren_NilParent() {
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	children, err := NewNodeChildren[string](nil, child1, child2)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNilParent)
	s.Require().Nil(children)
}

func (s *NodeChildrenTestSuite) TestNewNodeChildren_ValidParent() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	children, err := NewNodeChildren(parent, child1, child2)

	s.Require().NoError(err)
	s.Require().NotNil(children)
	s.Require().Equal(2, children.Size())
	s.Require().Equal(parent, children.Parent())
}

func (s *NodeChildrenTestSuite) TestNewNodeChildren_NoChildren() {
	parent := NewNode(uint64(100), "parent", nil)

	children, err := NewNodeChildren[string](parent)

	s.Require().NoError(err)
	s.Require().NotNil(children)
	s.Require().Equal(0, children.Size())
}

func (s *NodeChildrenTestSuite) TestNewNodeChildren_NilChildrenFiltered() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	var nilChild *Node[string]
	child2 := NewNode(uint64(2), "child2", nil)

	children, err := NewNodeChildren(parent, child1, nilChild, child2)

	s.Require().NoError(err)
	s.Require().NotNil(children)
	s.Require().Equal(2, children.Size())
}

func (s *NodeChildrenTestSuite) TestNewNodeChildren_AllNilChildren() {
	parent := NewNode(uint64(100), "parent", nil)
	var nilChild1, nilChild2 *Node[string]

	children, err := NewNodeChildren(parent, nilChild1, nilChild2)

	s.Require().NoError(err)
	s.Require().NotNil(children)
	s.Require().Equal(0, children.Size())
}

func (s *NodeChildrenTestSuite) TestNewNodeChildren_SetsParentRelationship() {
	parent := NewNode(uint64(100), "parent", nil)
	parent.BeholdRoot() // Make it root with level 0
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	_, err := NewNodeChildren(parent, child1, child2)

	s.Require().NoError(err)
	s.Require().Equal(1, child1.Level())
	s.Require().Equal(1, child2.Level())
	s.Require().NotNil(parent.Ptr().Next())
}

func (s *NodeChildrenTestSuite) TestSize() {
	testCases := []struct {
		name         string
		numChildren  int
		expectedSize int
	}{
		{"no children", 0, 0},
		{"one child", 1, 1},
		{"three children", 3, 3},
		{"five children", 5, 5},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			parent := NewNode(uint64(100), "parent", nil)
			nodes := make([]*Node[string], tc.numChildren)
			for i := 0; i < tc.numChildren; i++ {
				nodes[i] = NewNode(uint64(i+1), "child", nil)
			}

			children, err := NewNodeChildren(parent, nodes...)

			s.Require().NoError(err)
			s.Require().Equal(tc.expectedSize, children.Size())
		})
	}
}

func (s *NodeChildrenTestSuite) TestParent() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)

	children, err := NewNodeChildren(parent, child1)

	s.Require().NoError(err)
	s.Require().Equal(parent, children.Parent())
	s.Require().Equal(uint64(100), children.Parent().ID())
	s.Require().Equal("parent", children.Parent().Value())
}

func (s *NodeChildrenTestSuite) TestChild_Found() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	children, err := NewNodeChildren(parent, child1, child2, child3)
	s.Require().NoError(err)

	testCases := []struct {
		name     string
		id       uint64
		expected *Node[string]
	}{
		{"first child", 1, child1},
		{"second child", 2, child2},
		{"third child", 3, child3},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			child, err := children.Child(tc.id)

			s.Require().NoError(err)
			s.Require().NotNil(child)
			s.Require().Equal(tc.expected, child.Node())
		})
	}
}

func (s *NodeChildrenTestSuite) TestChild_NotFound() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)

	children, err := NewNodeChildren(parent, child1)
	s.Require().NoError(err)

	child, err := children.Child(999)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrChildNotFound)
	s.Require().Nil(child)
}

func (s *NodeChildrenTestSuite) TestChildNth_Found() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	children, err := NewNodeChildren(parent, child1, child2, child3)
	s.Require().NoError(err)

	testCases := []struct {
		name     string
		index    int
		expected *Node[string]
	}{
		{"first child (index 0)", 0, child1},
		{"second child (index 1)", 1, child2},
		{"third child (index 2)", 2, child3},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			child, err := children.ChildNth(tc.index)

			s.Require().NoError(err)
			s.Require().NotNil(child)
			s.Require().Equal(tc.expected, child.Node())
		})
	}
}

func (s *NodeChildrenTestSuite) TestChildNth_NegativeIndex() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)

	children, err := NewNodeChildren(parent, child1)
	s.Require().NoError(err)

	child, err := children.ChildNth(-1)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrChildNotFound)
	s.Require().Nil(child)
}

func (s *NodeChildrenTestSuite) TestChildNth_OutOfBounds() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)

	children, err := NewNodeChildren(parent, child1)
	s.Require().NoError(err)

	child, err := children.ChildNth(10)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrChildNotFound)
	s.Require().Nil(child)
}

func (s *NodeChildrenTestSuite) TestHasChild() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	children, err := NewNodeChildren(parent, child1, child2)
	s.Require().NoError(err)

	testCases := []struct {
		name     string
		id       uint64
		expected bool
	}{
		{"child 1 exists", 1, true},
		{"child 2 exists", 2, true},
		{"child 999 does not exist", 999, false},
		{"child 0 does not exist", 0, false},
		{"parent ID does not count as child", 100, false},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := children.HasChild(tc.id)
			s.Require().Equal(tc.expected, result)
		})
	}
}

func (s *NodeChildrenTestSuite) TestRelation() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	children, err := NewNodeChildren(parent, child1, child2)
	s.Require().NoError(err)

	testCases := []struct {
		name     string
		id       uint64
		expected string
	}{
		{"parent relation", 100, ParentRel},
		{"child 1 relation", 1, ChildRel},
		{"child 2 relation", 2, ChildRel},
		{"unrelated ID", 999, UnRelated},
		{"another unrelated ID", 50, UnRelated},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			relation := children.Relation(tc.id)
			s.Require().Equal(tc.expected, relation)
		})
	}
}

func (s *NodeChildrenTestSuite) TestNodes_Iteration() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	children, err := NewNodeChildren(parent, child1, child2, child3)
	s.Require().NoError(err)

	var collected []uint64
	var childNodes []*ChildNode[string]

	for hash, child := range children.Nodes() {
		collected = append(collected, hash)
		childNodes = append(childNodes, child)
	}

	// Should collect all 3 children
	s.Require().Len(collected, 3)
	s.Require().Len(childNodes, 3)

	// Verify each child is present
	var nodeIDs []uint64
	for _, cn := range childNodes {
		nodeIDs = append(nodeIDs, cn.Node().ID())
	}
	s.Require().Contains(nodeIDs, uint64(1))
	s.Require().Contains(nodeIDs, uint64(2))
	s.Require().Contains(nodeIDs, uint64(3))
}

func (s *NodeChildrenTestSuite) TestNodes_EarlyBreak() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	children, err := NewNodeChildren(parent, child1, child2, child3)
	s.Require().NoError(err)

	count := 0
	for _, child := range children.Nodes() {
		s.Require().NotNil(child)
		count++
		if count >= 2 {
			break
		}
	}

	s.Require().Equal(2, count)
}

func (s *NodeChildrenTestSuite) TestNodes_EmptyChildren() {
	parent := NewNode(uint64(100), "parent", nil)

	children, err := NewNodeChildren[string](parent)
	s.Require().NoError(err)

	count := 0
	for range children.Nodes() {
		count++
	}

	s.Require().Equal(0, count)
}

// ChildNodeTestSuite tests the ChildNode structure and its methods
type ChildNodeTestSuite struct {
	suite.Suite
}

func (s *ChildNodeTestSuite) TestNewChild() {
	parent := NewNode(uint64(100), "parent", nil)
	child := NewNode(uint64(1), "child", nil)

	childNode := newChild(parent, child)

	s.Require().NotNil(childNode)
	s.Require().Equal(parent, childNode.Parent())
	s.Require().Equal(child, childNode.Node())
	s.Require().NotZero(childNode.Hash())
}

func (s *ChildNodeTestSuite) TestHash() {
	parent := NewNode(uint64(100), "parent", nil)
	child := NewNode(uint64(1), "child", nil)

	childNode := newChild(parent, child)

	hash := childNode.Hash()
	s.Require().NotZero(hash)

	// Hash should be consistent
	s.Require().Equal(hash, childNode.Hash())
}

func (s *ChildNodeTestSuite) TestHash_UniquePerPair() {
	parent1 := NewNode(uint64(100), "parent1", nil)
	parent2 := NewNode(uint64(200), "parent2", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	// Different parent-child pairs should have different hashes
	hash1 := newChild(parent1, child1).Hash()
	hash2 := newChild(parent1, child2).Hash()
	hash3 := newChild(parent2, child1).Hash()
	hash4 := newChild(parent2, child2).Hash()

	s.Require().NotEqual(hash1, hash2)
	s.Require().NotEqual(hash1, hash3)
	s.Require().NotEqual(hash1, hash4)
	s.Require().NotEqual(hash2, hash3)
	s.Require().NotEqual(hash2, hash4)
	s.Require().NotEqual(hash3, hash4)
}

func (s *ChildNodeTestSuite) TestParent() {
	parent := NewNode(uint64(100), "parent", nil)
	child := NewNode(uint64(1), "child", nil)

	childNode := newChild(parent, child)

	s.Require().Equal(parent, childNode.Parent())
	s.Require().Equal(uint64(100), childNode.Parent().ID())
	s.Require().Equal("parent", childNode.Parent().Value())
}

func (s *ChildNodeTestSuite) TestNode() {
	parent := NewNode(uint64(100), "parent", nil)
	child := NewNode(uint64(1), "child", nil)

	childNode := newChild(parent, child)

	s.Require().Equal(child, childNode.Node())
	s.Require().Equal(uint64(1), childNode.Node().ID())
	s.Require().Equal("child", childNode.Node().Value())
}

func (s *ChildNodeTestSuite) TestChildNode_FromNodeChildren() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	children, err := NewNodeChildren(parent, child1, child2)
	s.Require().NoError(err)

	// Retrieve child node via NodeChildren
	childNode, err := children.Child(1)
	s.Require().NoError(err)
	s.Require().NotNil(childNode)

	// Verify ChildNode methods work correctly
	s.Require().Equal(parent, childNode.Parent())
	s.Require().Equal(child1, childNode.Node())
	s.Require().NotZero(childNode.Hash())
}

// ChildrenIntegrationTestSuite tests integration scenarios
type ChildrenIntegrationTestSuite struct {
	suite.Suite
}

func (s *ChildrenIntegrationTestSuite) TestCompleteTree() {
	// Create a simple tree: root -> [child1, child2, child3]
	root := NewNode(uint64(0), "root", nil)
	root.BeholdRoot()

	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)
	child3 := NewNode(uint64(3), "child3", nil)

	children, err := NewNodeChildren(root, child1, child2, child3)
	s.Require().NoError(err)

	// Verify tree structure
	s.Require().True(root.IsRoot())
	s.Require().Equal(0, root.Level())
	s.Require().Equal(1, child1.Level())
	s.Require().Equal(1, child2.Level())
	s.Require().Equal(1, child3.Level())

	// Verify children relationships
	s.Require().Equal(3, children.Size())
	s.Require().True(children.HasChild(1))
	s.Require().True(children.HasChild(2))
	s.Require().True(children.HasChild(3))

	// Verify relations
	s.Require().Equal(ParentRel, children.Relation(0))
	s.Require().Equal(ChildRel, children.Relation(1))
	s.Require().Equal(ChildRel, children.Relation(2))
	s.Require().Equal(ChildRel, children.Relation(3))
	s.Require().Equal(UnRelated, children.Relation(999))
}

func (s *ChildrenIntegrationTestSuite) TestMultiLevelTree() {
	// Create a two-level tree
	root := NewNode(uint64(0), "root", nil)
	root.BeholdRoot()

	child1 := NewNode(uint64(1), "child1", nil)
	child2 := NewNode(uint64(2), "child2", nil)

	rootChildren, err := NewNodeChildren(root, child1, child2)
	s.Require().NoError(err)

	// Add grandchildren to child1
	grandchild1 := NewNode(uint64(11), "grandchild1", nil)
	grandchild2 := NewNode(uint64(12), "grandchild2", nil)

	child1Children, err := NewNodeChildren(child1, grandchild1, grandchild2)
	s.Require().NoError(err)

	// Verify levels
	s.Require().Equal(0, root.Level())
	s.Require().Equal(1, child1.Level())
	s.Require().Equal(1, child2.Level())
	s.Require().Equal(2, grandchild1.Level())
	s.Require().Equal(2, grandchild2.Level())

	// Verify relationships
	s.Require().Equal(2, rootChildren.Size())
	s.Require().Equal(2, child1Children.Size())
	s.Require().True(rootChildren.HasChild(1))
	s.Require().True(child1Children.HasChild(11))
	s.Require().True(child1Children.HasChild(12))
}

func (s *ChildrenIntegrationTestSuite) TestChildNode_AccessThroughMethods() {
	parent := NewNode(uint64(100), "parent", nil)
	child1 := NewNode(uint64(1), "value1", nil)
	child2 := NewNode(uint64(2), "value2", nil)
	child3 := NewNode(uint64(3), "value3", nil)

	children, err := NewNodeChildren(parent, child1, child2, child3)
	s.Require().NoError(err)

	// Test Child() method
	cn1, err := children.Child(1)
	s.Require().NoError(err)
	s.Require().Equal("value1", cn1.Node().Value())

	// Test ChildNth() method
	cn0, err := children.ChildNth(0)
	s.Require().NoError(err)
	s.Require().Equal("value1", cn0.Node().Value())

	// Test iteration
	values := []string{}
	for _, child := range children.Nodes() {
		values = append(values, child.Node().Value())
	}
	s.Require().Len(values, 3)
	s.Require().Contains(values, "value1")
	s.Require().Contains(values, "value2")
	s.Require().Contains(values, "value3")
}

func (s *ChildrenIntegrationTestSuite) TestDifferentTypes() {
	// Test with int type
	intParent := NewNode(uint64(100), 42, nil)
	intChild := NewNode(uint64(1), 10, nil)

	intChildren, err := NewNodeChildren(intParent, intChild)
	s.Require().NoError(err)
	s.Require().Equal(1, intChildren.Size())

	ic, err := intChildren.Child(1)
	s.Require().NoError(err)
	s.Require().Equal(10, ic.Node().Value())

	// Test with float64 type
	floatParent := NewNode(uint64(200), 3.14, nil)
	floatChild := NewNode(uint64(2), 2.71, nil)

	floatChildren, err := NewNodeChildren(floatParent, floatChild)
	s.Require().NoError(err)
	s.Require().Equal(1, floatChildren.Size())

	fc, err := floatChildren.Child(2)
	s.Require().NoError(err)
	s.Require().Equal(2.71, fc.Node().Value())
}

// Test suite runners
func TestNodeChildrenTestSuite(t *testing.T) {
	suite.Run(t, new(NodeChildrenTestSuite))
}

func TestChildNodeTestSuite(t *testing.T) {
	suite.Run(t, new(ChildNodeTestSuite))
}

func TestChildrenIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ChildrenIntegrationTestSuite))
}
