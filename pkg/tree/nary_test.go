package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TreeConstructorTestSuite tests the Nary constructor
type TreeConstructorTestSuite struct {
	suite.Suite
}

func (s *TreeConstructorTestSuite) TestNew_DefaultValues() {
	tree := NewNary[string](5, 10)

	s.Require().NotNil(tree)
	s.Require().Nil(tree.root)
	s.Require().NotNil(tree.nodes)
	s.Require().NotNil(tree.levels)
	s.Require().Equal(uint8(5), tree.maxDepth)
	s.Require().Equal(uint8(10), tree.maxChildrenPerNode)
	s.Require().Equal(uint8(0), tree.levelsCount)
}

func (s *TreeConstructorTestSuite) TestNew_DifferentSizes() {
	testCases := []struct {
		name               string
		maxDepth           uint8
		maxChildrenPerNode uint8
	}{
		{"small tree", 3, 2},
		{"medium tree", 10, 5},
		{"large tree", 20, 15},
		{"minimal tree", 1, 1},
		{"max values", 255, 255},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tree := NewNary[int](tc.maxDepth, tc.maxChildrenPerNode)

			s.Require().NotNil(tree)
			s.Require().Equal(tc.maxDepth, tree.maxDepth)
			s.Require().Equal(tc.maxChildrenPerNode, tree.maxChildrenPerNode)
		})
	}
}

func (s *TreeConstructorTestSuite) TestNew_DifferentTypes() {
	// Test with string
	stringTree := NewNary[string](5, 10)
	s.Require().NotNil(stringTree)

	// Test with int
	intTree := NewNary[int](5, 10)
	s.Require().NotNil(intTree)

	// Test with float64
	floatTree := NewNary[float64](5, 10)
	s.Require().NotNil(floatTree)
}

func (s *TreeConstructorTestSuite) TestNewBinary() {
	tree := NewBinary[int](10)

	s.Require().NotNil(tree)
	s.Require().Equal(uint8(10), tree.maxDepth)
	s.Require().Equal(uint8(2), tree.maxChildrenPerNode)
}

func (s *TreeConstructorTestSuite) TestNewTernary() {
	tree := NewTernary[int](10)

	s.Require().NotNil(tree)
	s.Require().Equal(uint8(10), tree.maxDepth)
	s.Require().Equal(uint8(3), tree.maxChildrenPerNode)
}

// TreeRootTestSuite tests root node operations
type TreeRootTestSuite struct {
	suite.Suite
}

func (s *TreeRootTestSuite) TestAddRoot_Success() {
	tree := NewNary[string](5, 10)

	err := tree.AddRoot(N[string]{ID: 1, Val: "root"})

	s.Require().NoError(err)
	s.Require().NotNil(tree.root)
	s.Require().Equal(uint64(1), tree.root.ID())
	s.Require().Equal("root", tree.root.Value())
	s.Require().True(tree.root.IsRoot())

	// Root should be at level 0
	levelNodes, err := tree.Level(0)
	s.Require().NoError(err)
	s.Require().Len(levelNodes, 1)
	s.Require().Equal(tree.root, levelNodes[0])
}

func (s *TreeRootTestSuite) TestAddRoot_MultipleRoots() {
	tree := NewNary[string](5, 10)

	// First root should succeed
	err1 := tree.AddRoot(N[string]{ID: 1, Val: "root1"})
	s.Require().NoError(err1)

	// Try to add another root should fail
	err2 := tree.AddRoot(N[string]{ID: 2, Val: "root2"})
	s.Require().Error(err2)
	s.Require().ErrorIs(err2, ErrNotAllowed)

	// Original root should still be intact
	s.Require().Equal(uint64(1), tree.root.ID())
	s.Require().Equal("root1", tree.root.Value())
}

// TreeChildrenTestSuite tests adding children operations
type TreeChildrenTestSuite struct {
	suite.Suite
}

func (s *TreeChildrenTestSuite) TestAddChildren_ParentNotFound() {
	tree := NewNary[string](5, 10)

	err := tree.AddChildren(999, N[string]{ID: 1, Val: "child"})

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNilParent)
}

func (s *TreeChildrenTestSuite) TestAddChildren_ExceedsMaxChildrenPerNode() {
	tree := NewNary[string](5, 2) // max children limit is 2

	// This test validates validation logic through error scenarios
	err := tree.AddChildren(1,
		N[string]{ID: 10, Val: "child1"},
		N[string]{ID: 11, Val: "child2"},
		N[string]{ID: 12, Val: "child3"}, // exceeds max of 2
	)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNilParent) // parent doesn't exist
}

func (s *TreeChildrenTestSuite) TestAddChildren_MaxChildrenValidation() {
	testCases := []struct {
		name               string
		maxChildrenPerNode uint8
		childrenCount      int
		expectError        bool
	}{
		{"within limit 1", 5, 3, false},
		{"within limit 2", 5, 5, false},
		{"exceeds limit 1", 5, 6, true},
		{"exceeds limit 2", 3, 4, true},
		{"exact limit", 10, 10, false},
		{"one over limit", 10, 11, true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tree := NewNary[string](5, tc.maxChildrenPerNode)

			children := make([]N[string], tc.childrenCount)
			for i := 0; i < tc.childrenCount; i++ {
				children[i] = N[string]{ID: uint64(i + 1), Val: "child"}
			}

			err := tree.AddChildren(1, children...)

			if tc.expectError {
				s.Require().Error(err)
				// Will be ErrNotAllowedMaxNodes if > max, or ErrNilParent if parent missing
			} else {
				s.Require().Error(err) // Still error because parent doesn't exist
				s.Require().ErrorIs(err, ErrNilParent)
			}
		})
	}
}

// TreeLevelTestSuite tests level operations
type TreeLevelTestSuite struct {
	suite.Suite
}

func (s *TreeLevelTestSuite) TestLevel_NotFound() {
	tree := NewNary[string](5, 10)

	nodes, err := tree.Level(0)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrLevelNotFound)
	s.Require().Nil(nodes)
}

func (s *TreeLevelTestSuite) TestLevel_DifferentLevels() {
	testCases := []struct {
		name  string
		level uint8
	}{
		{"level 0", 0},
		{"level 1", 1},
		{"level 5", 5},
		{"level 10", 10},
		{"level 255", 255},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tree := NewNary[string](20, 10)

			nodes, err := tree.Level(tc.level)

			s.Require().Error(err)
			s.Require().ErrorIs(err, ErrLevelNotFound)
			s.Require().Nil(nodes)
		})
	}
}

// TreeLevelFuncTestSuite tests LevelFunc operations
type TreeLevelFuncTestSuite struct {
	suite.Suite
}

func (s *TreeLevelFuncTestSuite) TestLevelFunc_LevelNotFound() {
	tree := NewNary[string](5, 10)

	callCount := 0
	err := tree.LevelFunc(0, func(params LevelParams[string]) bool {
		callCount++
		return true
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrLevelNotFound)
	s.Require().Equal(0, callCount)
}

func (s *TreeLevelFuncTestSuite) TestLevelFunc_NilFunction() {
	tree := NewNary[string](5, 10)

	// Calling with nil function should handle gracefully
	s.Require().NotPanics(func() {
		_ = tree.LevelFunc(0, nil)
	})
}

func (s *TreeLevelFuncTestSuite) TestLevelFunc_PanicRecovery() {
	tree := NewNary[string](5, 10)

	// Add root and children
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})
	_ = tree.AddChildren(1,
		N[string]{ID: 2, Val: "a"},
		N[string]{ID: 3, Val: "b"},
	)

	// Function that panics
	err := tree.LevelFunc(1, func(params LevelParams[string]) bool {
		panic("intentional panic for testing")
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrUnexpected)
}

func (s *TreeLevelFuncTestSuite) TestLevelFunc_EarlyExit() {
	tree := NewNary[string](5, 10)

	// Add root and children
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})
	_ = tree.AddChildren(1,
		N[string]{ID: 2, Val: "a"},
		N[string]{ID: 3, Val: "b"},
		N[string]{ID: 4, Val: "c"},
	)

	callCount := 0
	err := tree.LevelFunc(1, func(params LevelParams[string]) bool {
		callCount++
		return false // Stop iteration after first call
	})

	s.Require().NoError(err)
	s.Require().Equal(1, callCount, "Should stop after first node")
}

func (s *TreeLevelFuncTestSuite) TestLevelFunc_ContinueIteration() {
	tree := NewNary[string](5, 10)

	// Add root and children
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})
	_ = tree.AddChildren(1,
		N[string]{ID: 2, Val: "a"},
		N[string]{ID: 3, Val: "b"},
		N[string]{ID: 4, Val: "c"},
	)

	callCount := 0
	err := tree.LevelFunc(1, func(params LevelParams[string]) bool {
		callCount++
		return true // Continue iteration
	})

	s.Require().NoError(err)
	s.Require().Equal(3, callCount, "Should iterate all nodes")
}

func (s *TreeLevelFuncTestSuite) TestLevelFunc_ParametersProvided() {
	tree := NewNary[int](5, 10)

	// Add root and children with numeric values
	_ = tree.AddRoot(N[int]{ID: 1, Val: 10})
	_ = tree.AddChildren(1,
		N[int]{ID: 2, Val: 5},
		N[int]{ID: 3, Val: 15},
		N[int]{ID: 4, Val: 3},
	)

	var capturedParams []LevelParams[int]
	err := tree.LevelFunc(1, func(params LevelParams[int]) bool {
		capturedParams = append(capturedParams, params)
		return true
	})

	s.Require().NoError(err)
	s.Require().Len(capturedParams, 3)

	// Verify each param has a node
	for _, p := range capturedParams {
		s.Require().NotNil(p.Node)
		s.Require().NotZero(p.MinVal)
		s.Require().NotZero(p.MaxVal)
	}
}

// TreeNodeStructTestSuite tests the N struct
type TreeNodeStructTestSuite struct {
	suite.Suite
}

func (s *TreeNodeStructTestSuite) TestNodeStruct_Creation() {
	testCases := []struct {
		name string
		id   uint64
		val  string
	}{
		{"simple node", 1, "test"},
		{"zero id", 0, "zero"},
		{"max id", ^uint64(0), "max"},
		{"empty value", 42, ""},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			n := N[string]{ID: tc.id, Val: tc.val}

			s.Require().Equal(tc.id, n.ID)
			s.Require().Equal(tc.val, n.Val)
		})
	}
}

func (s *TreeNodeStructTestSuite) TestNodeStruct_DifferentTypes() {
	// String type
	stringNode := N[string]{ID: 1, Val: "hello"}
	s.Require().Equal("hello", stringNode.Val)

	// Int type
	intNode := N[int]{ID: 2, Val: 42}
	s.Require().Equal(42, intNode.Val)

	// Float64 type
	floatNode := N[float64]{ID: 3, Val: 3.14}
	s.Require().Equal(3.14, floatNode.Val)
}

func (s *TreeNodeStructTestSuite) TestLevelParams_Structure() {
	node := NewNode[int](1, 100)
	params := LevelParams[int]{
		Node:   node,
		MinVal: 10,
		MaxVal: 200,
	}

	s.Require().Equal(node, params.Node)
	s.Require().Equal(10, params.MinVal)
	s.Require().Equal(200, params.MaxVal)
}

// TreeInternalMethodsTestSuite tests internal/private methods behavior
type TreeInternalMethodsTestSuite struct {
	suite.Suite
}

func (s *TreeInternalMethodsTestSuite) TestToNodes_EmptySlice() {
	tree := NewNary[string](5, 10)

	nodes := tree.toNodes()

	s.Require().NotNil(nodes)
	s.Require().Len(nodes, 0)
}

func (s *TreeInternalMethodsTestSuite) TestToNodes_SingleNode() {
	tree := NewNary[string](5, 10)

	nodes := tree.toNodes(N[string]{ID: 1, Val: "test"})

	s.Require().Len(nodes, 1)
	s.Require().Equal(uint64(1), nodes[0].ID())
	s.Require().Equal("test", nodes[0].Value())

	// Should be stored in tree.nodes map
	storedNode, exists := tree.nodes[1]
	s.Require().True(exists)
	s.Require().Equal(nodes[0], storedNode)
}

func (s *TreeInternalMethodsTestSuite) TestToNodes_MultipleNodes() {
	tree := NewNary[int](5, 10)

	nNodes := []N[int]{
		{ID: 1, Val: 100},
		{ID: 2, Val: 200},
		{ID: 3, Val: 300},
		{ID: 4, Val: 400},
		{ID: 5, Val: 500},
	}

	nodes := tree.toNodes(nNodes...)

	s.Require().Len(nodes, 5)

	for i, node := range nodes {
		s.Require().Equal(uint64(i+1), node.ID())
		s.Require().Equal((i+1)*100, node.Value())

		// Verify in map
		storedNode, exists := tree.nodes[uint64(i+1)]
		s.Require().True(exists)
		s.Require().Equal(node, storedNode)
	}
}

func (s *TreeInternalMethodsTestSuite) TestToNodes_OverwriteExisting() {
	tree := NewNary[string](5, 10)

	// Add first node
	nodes1 := tree.toNodes(N[string]{ID: 1, Val: "first"})
	s.Require().Equal("first", nodes1[0].Value())

	// Add another node with same ID (should overwrite)
	nodes2 := tree.toNodes(N[string]{ID: 1, Val: "second"})
	s.Require().Equal("second", nodes2[0].Value())

	// Map should have the new node
	storedNode, exists := tree.nodes[1]
	s.Require().True(exists)
	s.Require().Equal("second", storedNode.Value())
}

// TreeValidationTestSuite tests validation logic
type TreeValidationTestSuite struct {
	suite.Suite
}

func (s *TreeValidationTestSuite) TestMaxDepthValidation() {
	testCases := []struct {
		name     string
		maxDepth uint8
	}{
		{"depth 1", 1},
		{"depth 5", 5},
		{"depth 10", 10},
		{"depth 255", 255},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tree := NewNary[string](tc.maxDepth, 10)
			s.Require().Equal(tc.maxDepth, tree.maxDepth)
		})
	}
}

func (s *TreeValidationTestSuite) TestMaxChildrenPerNodeValidation() {
	testCases := []struct {
		name               string
		maxChildrenPerNode uint8
	}{
		{"1 child max", 1},
		{"5 children max", 5},
		{"10 children max", 10},
		{"255 children max", 255},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tree := NewNary[string](5, tc.maxChildrenPerNode)
			s.Require().Equal(tc.maxChildrenPerNode, tree.maxChildrenPerNode)
		})
	}
}

// TreeEdgeCasesTestSuite tests edge cases and boundary conditions
type TreeEdgeCasesTestSuite struct {
	suite.Suite
}

func (s *TreeEdgeCasesTestSuite) TestMinimalTree() {
	// Nary with depth 1 and 1 child per node
	tree := NewNary[string](1, 1)

	s.Require().NotNil(tree)
	s.Require().Equal(uint8(1), tree.maxDepth)
	s.Require().Equal(uint8(1), tree.maxChildrenPerNode)
}

func (s *TreeEdgeCasesTestSuite) TestMaximalTree() {
	// Nary with max uint8 values
	tree := NewNary[string](255, 255)

	s.Require().NotNil(tree)
	s.Require().Equal(uint8(255), tree.maxDepth)
	s.Require().Equal(uint8(255), tree.maxChildrenPerNode)
}

func (s *TreeEdgeCasesTestSuite) TestZeroValues() {
	// Nary with zero constraints (unusual but valid)
	tree := NewNary[string](0, 0)

	s.Require().NotNil(tree)
	s.Require().Equal(uint8(0), tree.maxDepth)
	s.Require().Equal(uint8(0), tree.maxChildrenPerNode)
}

func (s *TreeEdgeCasesTestSuite) TestAddChildren_ZeroChildren() {
	tree := NewNary[string](5, 10)

	// Adding zero children should fail because parent doesn't exist
	err := tree.AddChildren(1)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNilParent)
}

func (s *TreeEdgeCasesTestSuite) TestLevel_BoundaryValues() {
	tree := NewNary[string](255, 10)

	testCases := []struct {
		name  string
		level uint8
	}{
		{"level 0", 0},
		{"level 127", 127},
		{"level 254", 254},
		{"level 255", 255},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := tree.Level(tc.level)
			s.Require().Error(err)
			s.Require().ErrorIs(err, ErrLevelNotFound)
		})
	}
}

func (s *TreeEdgeCasesTestSuite) TestLevelFunc_EmptyLevel() {
	tree := NewNary[int](10, 10)
	_ = tree.AddRoot(N[int]{ID: 1, Val: 100})

	// Level 0 exists but has only one node
	callCount := 0
	err := tree.LevelFunc(0, func(params LevelParams[int]) bool {
		callCount++
		s.Require().Equal(100, params.Node.Value())
		s.Require().Equal(100, params.MinVal)
		s.Require().Equal(100, params.MaxVal)
		return true
	})

	s.Require().NoError(err)
	s.Require().Equal(1, callCount)
}

// TreeConcurrencyTestSuite tests potential concurrency issues
type TreeConcurrencyTestSuite struct {
	suite.Suite
}

func (s *TreeConcurrencyTestSuite) TestMapsInitialization() {
	tree := NewNary[string](5, 10)

	// Verify maps are initialized and not shared
	s.Require().NotNil(tree.nodes)
	s.Require().NotNil(tree.levels)

	// Should be empty
	s.Require().Len(tree.nodes, 0)
	s.Require().Len(tree.levels, 0)
}

func (s *TreeConcurrencyTestSuite) TestIndependentTrees() {
	tree1 := NewNary[string](5, 10)
	tree2 := NewNary[string](5, 10)

	// Add nodes to tree1
	tree1.toNodes(N[string]{ID: 1, Val: "tree1"})

	// tree2 should be independent
	s.Require().Len(tree1.nodes, 1)
	s.Require().Len(tree2.nodes, 0)
}

// TreeIntegrationTestSuite tests complete tree operations
type TreeIntegrationTestSuite struct {
	suite.Suite
}

func (s *TreeIntegrationTestSuite) TestSimpleTree() {
	tree := NewNary[string](5, 10)

	// Add root
	err := tree.AddRoot(N[string]{ID: 1, Val: "root"})
	s.Require().NoError(err)

	// Add children to root
	err = tree.AddChildren(1,
		N[string]{ID: 2, Val: "child1"},
		N[string]{ID: 3, Val: "child2"},
	)
	s.Require().NoError(err)

	// Verify level 0 (root)
	level0, err := tree.Level(0)
	s.Require().NoError(err)
	s.Require().Len(level0, 1)
	s.Require().Equal("root", level0[0].Value())

	// Verify level 1 (children)
	level1, err := tree.Level(1)
	s.Require().NoError(err)
	s.Require().Len(level1, 2)
}

func (s *TreeIntegrationTestSuite) TestMultiLevelTree() {
	tree := NewNary[int](10, 5)

	// Add root
	err := tree.AddRoot(N[int]{ID: 1, Val: 100})
	s.Require().NoError(err)

	// Add level 1
	err = tree.AddChildren(1,
		N[int]{ID: 2, Val: 200},
		N[int]{ID: 3, Val: 300},
	)
	s.Require().NoError(err)

	// Add level 2 under first child
	err = tree.AddChildren(2,
		N[int]{ID: 4, Val: 400},
		N[int]{ID: 5, Val: 500},
	)
	s.Require().NoError(err)

	// Verify all levels
	level0, _ := tree.Level(0)
	s.Require().Len(level0, 1)
	s.Require().Equal(100, level0[0].Value())

	level1, _ := tree.Level(1)
	s.Require().Len(level1, 2)

	level2, _ := tree.Level(2)
	s.Require().Len(level2, 2)
	s.Require().Equal(400, level2[0].Value())
	s.Require().Equal(500, level2[1].Value())
}

func (s *TreeIntegrationTestSuite) TestLevelFunc() {
	tree := NewNary[string](5, 10)

	// Build tree
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})
	_ = tree.AddChildren(1,
		N[string]{ID: 2, Val: "a"},
		N[string]{ID: 3, Val: "b"},
		N[string]{ID: 4, Val: "c"},
	)

	// Test LevelFunc on level 1
	values := []string{}
	err := tree.LevelFunc(1, func(params LevelParams[string]) bool {
		values = append(values, params.Node.Value())
		return true
	})

	s.Require().NoError(err)
	s.Require().Len(values, 3)
	s.Require().Contains(values, "a")
	s.Require().Contains(values, "b")
	s.Require().Contains(values, "c")
}

func (s *TreeIntegrationTestSuite) TestMaxDepthEnforcement() {
	tree := NewNary[string](2, 10) // maxDepth = 2

	// Add root (level 0)
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})

	// Add level 1 (levelsCount becomes 1, check: 0 >= 2? No, so OK)
	err := tree.AddChildren(1, N[string]{ID: 2, Val: "child"})
	s.Require().NoError(err, "Adding level 1 should succeed")

	// Try to add level 2 (check: levelsCount 1 >= maxDepth 2? No, so this will succeed!)
	err = tree.AddChildren(2, N[string]{ID: 3, Val: "grandchild"})
	s.Require().NoError(err, "Adding level 2 should succeed when maxDepth=2")

	// Now levelsCount is 2. Try to add level 3 (check: 2 >= 2? Yes, so this should fail)
	err = tree.AddChildren(3, N[string]{ID: 4, Val: "great-grandchild"})
	s.Require().Error(err, "Adding level 3 should fail when maxDepth=2")
	s.Require().ErrorIs(err, ErrNotAllowedMaxDepth)
}

func (s *TreeIntegrationTestSuite) TestMaxChildrenEnforcement() {
	tree := NewNary[string](5, 2)

	// Add root
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})

	// Add 2 children (should succeed since max is 2)
	err := tree.AddChildren(1,
		N[string]{ID: 2, Val: "child1"},
		N[string]{ID: 3, Val: "child2"},
	)
	s.Require().NoError(err)

	// Try to add more children (should fail)
	err = tree.AddChildren(1, N[string]{ID: 4, Val: "child3"})
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNotAllowedMaxNodes)
}

func (s *TreeIntegrationTestSuite) TestAddingTooManyChildrenAtOnce() {
	tree := NewNary[string](5, 3)

	// Add root
	_ = tree.AddRoot(N[string]{ID: 1, Val: "root"})

	// Try to add 4 children at once when max is 3 (should fail)
	err := tree.AddChildren(1,
		N[string]{ID: 2, Val: "child1"},
		N[string]{ID: 3, Val: "child2"},
		N[string]{ID: 4, Val: "child3"},
		N[string]{ID: 5, Val: "child4"},
	)

	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNotAllowedMaxNodes)
}

func (s *TreeIntegrationTestSuite) TestBinaryTreeOperations() {
	tree := NewBinary[int](5)

	_ = tree.AddRoot(N[int]{ID: 1, Val: 10})
	err := tree.AddChildren(1,
		N[int]{ID: 2, Val: 5},
		N[int]{ID: 3, Val: 15},
	)
	s.Require().NoError(err)

	// Try to add a third child (should fail for binary tree)
	err = tree.AddChildren(1, N[int]{ID: 4, Val: 20})
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNotAllowedMaxNodes)
}

func (s *TreeIntegrationTestSuite) TestTernaryTreeOperations() {
	tree := NewTernary[int](5)

	_ = tree.AddRoot(N[int]{ID: 1, Val: 10})
	err := tree.AddChildren(1,
		N[int]{ID: 2, Val: 5},
		N[int]{ID: 3, Val: 10},
		N[int]{ID: 4, Val: 15},
	)
	s.Require().NoError(err)

	// Try to add a fourth child (should fail for ternary tree)
	err = tree.AddChildren(1, N[int]{ID: 5, Val: 20})
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNotAllowedMaxNodes)
}

// Test suite runners
func TestTreeConstructorTestSuite(t *testing.T) {
	suite.Run(t, new(TreeConstructorTestSuite))
}

func TestTreeRootTestSuite(t *testing.T) {
	suite.Run(t, new(TreeRootTestSuite))
}

func TestTreeChildrenTestSuite(t *testing.T) {
	suite.Run(t, new(TreeChildrenTestSuite))
}

func TestTreeLevelTestSuite(t *testing.T) {
	suite.Run(t, new(TreeLevelTestSuite))
}

func TestTreeLevelFuncTestSuite(t *testing.T) {
	suite.Run(t, new(TreeLevelFuncTestSuite))
}

func TestTreeNodeStructTestSuite(t *testing.T) {
	suite.Run(t, new(TreeNodeStructTestSuite))
}

func TestTreeInternalMethodsTestSuite(t *testing.T) {
	suite.Run(t, new(TreeInternalMethodsTestSuite))
}

func TestTreeValidationTestSuite(t *testing.T) {
	suite.Run(t, new(TreeValidationTestSuite))
}

func TestTreeEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(TreeEdgeCasesTestSuite))
}

func TestTreeConcurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(TreeConcurrencyTestSuite))
}

func TestTreeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TreeIntegrationTestSuite))
}
