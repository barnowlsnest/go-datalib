package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BSTTestSuite is the main test suite for BST operations
type BSTTestSuite struct {
	suite.Suite
	bst *BST[int]
}

func (s *BSTTestSuite) SetupTest() {
	s.bst = NewBST[int]()
}

// buildTree is a helper to build a tree from values
func (s *BSTTestSuite) buildTree(values []int) {
	for i, v := range values {
		s.bst.Insert(NewNodeValue(uint64(i+1), v))
	}
}

// collectValuesInt is a helper to collect values during traversal
func collectValuesInt(traversalFunc func(func(*BinaryNode[int]))) []int {
	var values []int
	traversalFunc(func(node *BinaryNode[int]) {
		props, _ := node.Props()
		values = append(values, props.Value)
	})
	return values
}

func TestBSTTestSuite(t *testing.T) {
	suite.Run(t, new(BSTTestSuite))
}

// Test basic operations
func (s *BSTTestSuite) TestNewBST() {
	testCases := []struct {
		name     string
		checkFn  func() bool
		expected bool
	}{
		{"is not nil", func() bool { return s.bst != nil }, true},
		{"is empty", func() bool { return s.bst.IsEmpty() }, true},
		{"size is zero", func() bool { return s.bst.Size() == 0 }, true},
		{"root is nil", func() bool { return s.bst.Root() == nil }, true},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			assert.Equal(s.T(), tc.expected, tc.checkFn())
		})
	}
}

func (s *BSTTestSuite) TestInsert() {
	testCases := []struct {
		name          string
		insertValues  []int
		expectedSize  int
		expectedEmpty bool
		checkRoot     bool
	}{
		{
			name:          "single node",
			insertValues:  []int{50},
			expectedSize:  1,
			expectedEmpty: false,
			checkRoot:     true,
		},
		{
			name:          "multiple nodes",
			insertValues:  []int{50, 30, 70, 20, 40},
			expectedSize:  5,
			expectedEmpty: false,
			checkRoot:     false,
		},
		{
			name:          "with duplicate",
			insertValues:  []int{50, 30, 50},
			expectedSize:  2,
			expectedEmpty: false,
			checkRoot:     false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for i, v := range tc.insertValues {
				bst.Insert(NewNodeValue(uint64(i+1), v))
			}
			assert.Equal(s.T(), tc.expectedSize, bst.Size())
			assert.Equal(s.T(), tc.expectedEmpty, bst.IsEmpty())
			if tc.checkRoot {
				assert.NotNil(s.T(), bst.Root())
				assert.True(s.T(), bst.Root().IsRoot())
			}
		})
	}
}

func (s *BSTTestSuite) TestInsertEdgeCases() {
	testCases := []struct {
		name         string
		value        *NodeValue[int]
		shouldInsert bool
	}{
		{"nil value", nil, false},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			inserted := bst.Insert(tc.value)
			assert.Equal(s.T(), tc.shouldInsert, inserted)
		})
	}
}

func (s *BSTTestSuite) TestSearch() {
	treeValues := []int{50, 30, 70, 20, 40, 80}
	s.buildTree(treeValues)

	testCases := []struct {
		name        string
		searchValue int
		shouldFind  bool
		expectValue int
	}{
		{"find root", 50, true, 50},
		{"find leaf left", 20, true, 20},
		{"find leaf right", 80, true, 80},
		{"find internal", 30, true, 30},
		{"not found", 999, false, 0},
		{"not found negative", -10, false, 0},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			node := s.bst.Search(tc.searchValue)
			if tc.shouldFind {
				assert.NotNil(s.T(), node)
				props, _ := node.Props()
				assert.Equal(s.T(), tc.expectValue, props.Value)
			} else {
				assert.Nil(s.T(), node)
			}
		})
	}
}

func (s *BSTTestSuite) TestSearchEmptyTree() {
	node := s.bst.Search(50)
	assert.Nil(s.T(), node)
}

func (s *BSTTestSuite) TestDelete() {
	testCases := []struct {
		name          string
		treeValues    []int
		deleteValue   int
		shouldDelete  bool
		expectedSize  int
		verifyAbsent  []int
		verifyPresent []int
	}{
		{
			name:          "from empty tree",
			treeValues:    []int{},
			deleteValue:   50,
			shouldDelete:  false,
			expectedSize:  0,
			verifyAbsent:  []int{50},
			verifyPresent: []int{},
		},
		{
			name:          "non-existing value",
			treeValues:    []int{50},
			deleteValue:   999,
			shouldDelete:  false,
			expectedSize:  1,
			verifyAbsent:  []int{999},
			verifyPresent: []int{50},
		},
		{
			name:          "leaf node",
			treeValues:    []int{50, 30, 70},
			deleteValue:   30,
			shouldDelete:  true,
			expectedSize:  2,
			verifyAbsent:  []int{30},
			verifyPresent: []int{50, 70},
		},
		{
			name:          "node with left child only",
			treeValues:    []int{50, 30, 20},
			deleteValue:   30,
			shouldDelete:  true,
			expectedSize:  2,
			verifyAbsent:  []int{30},
			verifyPresent: []int{50, 20},
		},
		{
			name:          "node with right child only",
			treeValues:    []int{50, 30, 40},
			deleteValue:   30,
			shouldDelete:  true,
			expectedSize:  2,
			verifyAbsent:  []int{30},
			verifyPresent: []int{50, 40},
		},
		{
			name:          "node with two children",
			treeValues:    []int{50, 30, 70, 20, 40, 60, 80},
			deleteValue:   50,
			shouldDelete:  true,
			expectedSize:  6,
			verifyAbsent:  []int{50},
			verifyPresent: []int{30, 70, 20, 40, 60, 80},
		},
		{
			name:          "root with no children",
			treeValues:    []int{50},
			deleteValue:   50,
			shouldDelete:  true,
			expectedSize:  0,
			verifyAbsent:  []int{50},
			verifyPresent: []int{},
		},
		{
			name:          "root with one child",
			treeValues:    []int{50, 30},
			deleteValue:   50,
			shouldDelete:  true,
			expectedSize:  1,
			verifyAbsent:  []int{50},
			verifyPresent: []int{30},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for i, v := range tc.treeValues {
				bst.Insert(NewNodeValue(uint64(i+1), v))
			}

			deleted := bst.Delete(tc.deleteValue)
			assert.Equal(s.T(), tc.shouldDelete, deleted)
			assert.Equal(s.T(), tc.expectedSize, bst.Size())

			for _, v := range tc.verifyAbsent {
				assert.Nil(s.T(), bst.Search(v), "value %d should be absent", v)
			}

			for _, v := range tc.verifyPresent {
				assert.NotNil(s.T(), bst.Search(v), "value %d should be present", v)
			}
		})
	}
}

func (s *BSTTestSuite) TestDeleteAllNodes() {
	values := []int{50, 30, 70}
	s.buildTree(values)

	for _, v := range values {
		s.bst.Delete(v)
	}

	assert.True(s.T(), s.bst.IsEmpty())
	assert.Equal(s.T(), 0, s.bst.Size())
	assert.Nil(s.T(), s.bst.Root())
}

func (s *BSTTestSuite) TestTraversals() {
	treeValues := []int{50, 30, 70, 20, 40, 80}
	s.buildTree(treeValues)

	testCases := []struct {
		name     string
		traverse func(func(*BinaryNode[int]))
		expected []int
	}{
		{
			name:     "InOrder",
			traverse: s.bst.InOrder,
			expected: []int{20, 30, 40, 50, 70, 80},
		},
		{
			name:     "PreOrder",
			traverse: s.bst.PreOrder,
			expected: []int{50, 30, 20, 40, 70, 80},
		},
		{
			name:     "PostOrder",
			traverse: s.bst.PostOrder,
			expected: []int{20, 40, 30, 80, 70, 50},
		},
		{
			name:     "LevelOrder",
			traverse: s.bst.LevelOrder,
			expected: []int{50, 30, 70, 20, 40, 80},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			values := collectValuesInt(tc.traverse)
			assert.Equal(s.T(), tc.expected, values)
		})
	}
}

func (s *BSTTestSuite) TestTraversalEdgeCases() {
	testCases := []struct {
		name        string
		setup       func() *BST[int]
		traversalFn func(bst *BST[int], visit func(*BinaryNode[int]))
	}{
		{
			name:        "InOrder on empty tree",
			setup:       func() *BST[int] { return NewBST[int]() },
			traversalFn: func(bst *BST[int], visit func(*BinaryNode[int])) { bst.InOrder(visit) },
		},
		{
			name:        "PreOrder on empty tree",
			setup:       func() *BST[int] { return NewBST[int]() },
			traversalFn: func(bst *BST[int], visit func(*BinaryNode[int])) { bst.PreOrder(visit) },
		},
		{
			name:        "PostOrder on empty tree",
			setup:       func() *BST[int] { return NewBST[int]() },
			traversalFn: func(bst *BST[int], visit func(*BinaryNode[int])) { bst.PostOrder(visit) },
		},
		{
			name:        "LevelOrder on empty tree",
			setup:       func() *BST[int] { return NewBST[int]() },
			traversalFn: func(bst *BST[int], visit func(*BinaryNode[int])) { bst.LevelOrder(visit) },
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := tc.setup()
			var values []int
			tc.traversalFn(bst, func(node *BinaryNode[int]) {
				props, _ := node.Props()
				values = append(values, props.Value)
			})
			assert.Empty(s.T(), values)
		})
	}
}

func (s *BSTTestSuite) TestTraversalWithNilVisitor() {
	s.buildTree([]int{50, 30, 70})

	// Should not panic
	s.bst.InOrder(nil)
	s.bst.PreOrder(nil)
	s.bst.PostOrder(nil)
	s.bst.LevelOrder(nil)

	assert.Equal(s.T(), 3, s.bst.Size())
}

func (s *BSTTestSuite) TestMinMaxHeight() {
	testCases := []struct {
		name           string
		treeValues     []int
		expectedMin    *int
		expectedMax    *int
		expectedHeight int
	}{
		{
			name:           "empty tree",
			treeValues:     []int{},
			expectedMin:    nil,
			expectedMax:    nil,
			expectedHeight: -1,
		},
		{
			name:           "single node",
			treeValues:     []int{50},
			expectedMin:    intPtr(50),
			expectedMax:    intPtr(50),
			expectedHeight: 0,
		},
		{
			name:           "balanced tree",
			treeValues:     []int{50, 30, 70, 20, 40, 80},
			expectedMin:    intPtr(20),
			expectedMax:    intPtr(80),
			expectedHeight: 2,
		},
		{
			name:           "right-skewed tree",
			treeValues:     []int{10, 20, 30, 40},
			expectedMin:    intPtr(10),
			expectedMax:    intPtr(40),
			expectedHeight: 3,
		},
		{
			name:           "left-skewed tree",
			treeValues:     []int{40, 30, 20, 10},
			expectedMin:    intPtr(10),
			expectedMax:    intPtr(40),
			expectedHeight: 3,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for i, v := range tc.treeValues {
				bst.Insert(NewNodeValue(uint64(i+1), v))
			}

			// Test Min
			minNode := bst.Min()
			if tc.expectedMin == nil {
				assert.Nil(s.T(), minNode)
			} else {
				assert.NotNil(s.T(), minNode)
				props, _ := minNode.Props()
				assert.Equal(s.T(), *tc.expectedMin, props.Value)
			}

			// Test Max
			maxNode := bst.Max()
			if tc.expectedMax == nil {
				assert.Nil(s.T(), maxNode)
			} else {
				assert.NotNil(s.T(), maxNode)
				props, _ := maxNode.Props()
				assert.Equal(s.T(), *tc.expectedMax, props.Value)
			}

			// Test Height
			height := bst.Height()
			assert.Equal(s.T(), tc.expectedHeight, height)
		})
	}
}

func (s *BSTTestSuite) TestComplexScenarios() {
	testCases := []struct {
		name       string
		operations func(bst *BST[int])
		verify     func(t *testing.T, bst *BST[int])
	}{
		{
			name: "mixed insert and delete",
			operations: func(bst *BST[int]) {
				bst.Insert(NewNodeValue(1, 50))
				bst.Insert(NewNodeValue(2, 30))
				bst.Insert(NewNodeValue(3, 70))
				bst.Delete(30)
				bst.Insert(NewNodeValue(4, 40))
				bst.Delete(70)
				bst.Insert(NewNodeValue(5, 60))
			},
			verify: func(t *testing.T, bst *BST[int]) {
				assert.Equal(t, 3, bst.Size())
				assert.NotNil(t, bst.Search(50))
				assert.NotNil(t, bst.Search(40))
				assert.NotNil(t, bst.Search(60))
				assert.Nil(t, bst.Search(30))
				assert.Nil(t, bst.Search(70))
			},
		},
		{
			name: "BST property maintained after delete",
			operations: func(bst *BST[int]) {
				values := []int{50, 30, 70, 20, 40, 60, 80}
				for i, v := range values {
					bst.Insert(NewNodeValue(uint64(i+1), v))
				}
				bst.Delete(50) // Delete node with two children
			},
			verify: func(t *testing.T, bst *BST[int]) {
				// InOrder should produce sorted sequence
				var values []int
				bst.InOrder(func(node *BinaryNode[int]) {
					props, _ := node.Props()
					values = append(values, props.Value)
				})
				for i := 1; i < len(values); i++ {
					assert.Less(t, values[i-1], values[i], "BST property violated")
				}
			},
		},
		{
			name: "large number of operations",
			operations: func(bst *BST[int]) {
				// Insert 100 nodes
				for i := 0; i < 100; i++ {
					bst.Insert(NewNodeValue(uint64(i), i))
				}
				// Delete every other node
				for i := 0; i < 100; i += 2 {
					bst.Delete(i)
				}
			},
			verify: func(t *testing.T, bst *BST[int]) {
				assert.Equal(t, 50, bst.Size())
				// Verify odd numbers exist
				for i := 1; i < 100; i += 2 {
					assert.NotNil(t, bst.Search(i))
				}
				// Verify even numbers don't exist
				for i := 0; i < 100; i += 2 {
					assert.Nil(t, bst.Search(i))
				}
			},
		},
		{
			name: "delete in reverse insertion order",
			operations: func(bst *BST[int]) {
				values := []int{50, 30, 70, 20, 40, 60, 80}
				for i, v := range values {
					bst.Insert(NewNodeValue(uint64(i+1), v))
				}
				for i := len(values) - 1; i >= 0; i-- {
					bst.Delete(values[i])
				}
			},
			verify: func(t *testing.T, bst *BST[int]) {
				assert.True(t, bst.IsEmpty())
				assert.Equal(t, 0, bst.Size())
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			tc.operations(bst)
			tc.verify(s.T(), bst)
		})
	}
}

func (s *BSTTestSuite) TestDifferentTypes() {
	s.Run("string type", func() {
		bst := NewBST[string]()
		words := []string{"dog", "cat", "elephant", "ant"}
		for i, w := range words {
			bst.Insert(NewNodeValue(uint64(i+1), w))
		}

		assert.Equal(s.T(), 4, bst.Size())

		// InOrder should give sorted strings
		var values []string
		bst.InOrder(func(node *BinaryNode[string]) {
			props, _ := node.Props()
			values = append(values, props.Value)
		})

		expected := []string{"ant", "cat", "dog", "elephant"}
		assert.Equal(s.T(), expected, values)
	})

	s.Run("float64 type", func() {
		bst := NewBST[float64]()
		nums := []float64{3.14, 2.71, 1.41, 1.73}
		for i, n := range nums {
			bst.Insert(NewNodeValue(uint64(i+1), n))
		}

		assert.Equal(s.T(), 4, bst.Size())

		minNode := bst.Min()
		assert.NotNil(s.T(), minNode)
		props, _ := minNode.Props()
		assert.Equal(s.T(), 1.41, props.Value)
	})
}

func (s *BSTTestSuite) TestAllTraversalsVisitAllNodes() {
	treeValues := []int{50, 30, 70, 20, 40, 60, 80}
	s.buildTree(treeValues)

	traversals := map[string]func(func(*BinaryNode[int])){
		"InOrder":    s.bst.InOrder,
		"PreOrder":   s.bst.PreOrder,
		"PostOrder":  s.bst.PostOrder,
		"LevelOrder": s.bst.LevelOrder,
	}

	for name, traverseFn := range traversals {
		s.Run(name, func() {
			count := 0
			traverseFn(func(node *BinaryNode[int]) {
				count++
			})
			assert.Equal(s.T(), len(treeValues), count)
		})
	}
}

// Helper function to create int pointer
func intPtr(v int) *int {
	return &v
}
