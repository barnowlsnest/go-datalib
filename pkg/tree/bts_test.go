package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BSTTestSuite struct {
	suite.Suite
}

func (s *BSTTestSuite) TestNewBST() {
	bst := NewBST[int]()
	s.Require().NotNil(bst)
	s.Require().Equal(0, bst.Size())
	s.Require().True(bst.IsEmpty())
	s.Require().Nil(bst.Root())
}

func (s *BSTTestSuite) TestInsert() {
	testCases := []struct {
		name         string
		insertValues []struct {
			id  uint64
			val int
		}
		expectedSize   int
		expectedRoot   int
		expectedHeight int
	}{
		{
			name: "insert single value",
			insertValues: []struct {
				id  uint64
				val int
			}{
				{1, 50},
			},
			expectedSize:   1,
			expectedRoot:   50,
			expectedHeight: 0,
		},
		{
			name: "insert multiple values in order",
			insertValues: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
			},
			expectedSize:   3,
			expectedRoot:   50,
			expectedHeight: 1,
		},
		{
			name: "insert values creating unbalanced tree",
			insertValues: []struct {
				id  uint64
				val int
			}{
				{1, 10},
				{2, 20},
				{3, 30},
				{4, 40},
			},
			expectedSize:   4,
			expectedRoot:   10,
			expectedHeight: 3,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()

			for _, iv := range tc.insertValues {
				bst.Insert(iv.id, iv.val)
			}

			s.Require().Equal(tc.expectedSize, bst.Size())
			s.Require().False(bst.IsEmpty())
			s.Require().NotNil(bst.Root())
			s.Require().Equal(tc.expectedRoot, bst.Root().val)
			s.Require().Equal(tc.expectedHeight, bst.Height())
		})
	}
}

func (s *BSTTestSuite) TestSearch() {
	bst := NewBST[int]()
	values := []struct {
		id  uint64
		val int
	}{
		{1, 50},
		{2, 30},
		{3, 70},
		{4, 20},
		{5, 40},
		{6, 60},
		{7, 80},
	}

	for _, v := range values {
		bst.Insert(v.id, v.val)
	}

	testCases := []struct {
		name        string
		searchValue int
		shouldFind  bool
	}{
		{
			name:        "search for root value",
			searchValue: 50,
			shouldFind:  true,
		},
		{
			name:        "search for left child",
			searchValue: 30,
			shouldFind:  true,
		},
		{
			name:        "search for right child",
			searchValue: 70,
			shouldFind:  true,
		},
		{
			name:        "search for leaf node",
			searchValue: 20,
			shouldFind:  true,
		},
		{
			name:        "search for non-existent value",
			searchValue: 100,
			shouldFind:  false,
		},
		{
			name:        "search for another non-existent value",
			searchValue: 25,
			shouldFind:  false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := bst.Search(tc.searchValue)
			if tc.shouldFind {
				s.Require().NotNil(result)
				s.Require().Equal(tc.searchValue, result.val)
			} else {
				s.Require().Nil(result)
			}
		})
	}
}

func (s *BSTTestSuite) TestContains() {
	bst := NewBST[int]()
	bst.Insert(1, 50)
	bst.Insert(2, 30)
	bst.Insert(3, 70)

	testCases := []struct {
		name     string
		value    int
		expected bool
	}{
		{
			name:     "contains existing value",
			value:    50,
			expected: true,
		},
		{
			name:     "does not contain non-existent value",
			value:    100,
			expected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := bst.Contains(tc.value)
			s.Require().Equal(tc.expected, result)
		})
	}
}

func (s *BSTTestSuite) TestMin() {
	testCases := []struct {
		name       string
		insertVals []struct {
			id  uint64
			val int
		}
		expectedMin int
		shouldFind  bool
	}{
		{
			name: "empty tree",
			insertVals: []struct {
				id  uint64
				val int
			}{},
			expectedMin: 0,
			shouldFind:  false,
		},
		{
			name: "single node",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
			},
			expectedMin: 50,
			shouldFind:  true,
		},
		{
			name: "multiple nodes",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
				{4, 20},
				{5, 40},
			},
			expectedMin: 20,
			shouldFind:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for _, v := range tc.insertVals {
				bst.Insert(v.id, v.val)
			}

			minVal, found := bst.Min()
			s.Require().Equal(tc.shouldFind, found)
			if found {
				s.Require().Equal(tc.expectedMin, minVal)
			}
		})
	}
}

func (s *BSTTestSuite) TestMax() {
	testCases := []struct {
		name       string
		insertVals []struct {
			id  uint64
			val int
		}
		expectedMax int
		shouldFind  bool
	}{
		{
			name: "empty tree",
			insertVals: []struct {
				id  uint64
				val int
			}{},
			expectedMax: 0,
			shouldFind:  false,
		},
		{
			name: "single node",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
			},
			expectedMax: 50,
			shouldFind:  true,
		},
		{
			name: "multiple nodes",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
				{4, 60},
				{5, 80},
			},
			expectedMax: 80,
			shouldFind:  true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for _, v := range tc.insertVals {
				bst.Insert(v.id, v.val)
			}

			maxVal, found := bst.Max()
			s.Require().Equal(tc.shouldFind, found)
			if found {
				s.Require().Equal(tc.expectedMax, maxVal)
			}
		})
	}
}

func (s *BSTTestSuite) TestDelete() {
	testCases := []struct {
		name       string
		insertVals []struct {
			id  uint64
			val int
		}
		deleteValue     int
		expectedDeleted bool
		expectedSize    int
		remainingValues []int
		shouldContain   map[int]bool
	}{
		{
			name: "delete from empty tree",
			insertVals: []struct {
				id  uint64
				val int
			}{},
			deleteValue:     50,
			expectedDeleted: false,
			expectedSize:    0,
		},
		{
			name: "delete leaf node",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
				{4, 20},
			},
			deleteValue:     20,
			expectedDeleted: true,
			expectedSize:    3,
			shouldContain: map[int]bool{
				50: true,
				30: true,
				70: true,
				20: false,
			},
		},
		{
			name: "delete node with one child (left)",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 20},
			},
			deleteValue:     30,
			expectedDeleted: true,
			expectedSize:    2,
			shouldContain: map[int]bool{
				50: true,
				20: true,
				30: false,
			},
		},
		{
			name: "delete node with one child (right)",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 40},
			},
			deleteValue:     30,
			expectedDeleted: true,
			expectedSize:    2,
			shouldContain: map[int]bool{
				50: true,
				40: true,
				30: false,
			},
		},
		{
			name: "delete node with two children",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
				{4, 20},
				{5, 40},
				{6, 60},
				{7, 80},
			},
			deleteValue:     50,
			expectedDeleted: true,
			expectedSize:    6,
			shouldContain: map[int]bool{
				30: true,
				70: true,
				20: true,
				40: true,
				60: true,
				80: true,
				50: false,
			},
		},
		{
			name: "delete root node with only left child",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
			},
			deleteValue:     50,
			expectedDeleted: true,
			expectedSize:    1,
			shouldContain: map[int]bool{
				30: true,
				50: false,
			},
		},
		{
			name: "delete root node with only right child",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 70},
			},
			deleteValue:     50,
			expectedDeleted: true,
			expectedSize:    1,
			shouldContain: map[int]bool{
				70: true,
				50: false,
			},
		},
		{
			name: "delete non-existent value",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
			},
			deleteValue:     100,
			expectedDeleted: false,
			expectedSize:    3,
			shouldContain: map[int]bool{
				50: true,
				30: true,
				70: true,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for _, v := range tc.insertVals {
				bst.Insert(v.id, v.val)
			}

			deleted := bst.Delete(tc.deleteValue)
			s.Require().Equal(tc.expectedDeleted, deleted)
			s.Require().Equal(tc.expectedSize, bst.Size())

			if tc.shouldContain != nil {
				for val, shouldExist := range tc.shouldContain {
					s.Require().Equal(shouldExist, bst.Contains(val),
						"Value %d should exist: %v", val, shouldExist)
				}
			}
		})
	}
}

func (s *BSTTestSuite) TestInOrderTraversal() {
	bst := NewBST[int]()
	values := []struct {
		id  uint64
		val int
	}{
		{1, 50},
		{2, 30},
		{3, 70},
		{4, 20},
		{5, 40},
		{6, 60},
		{7, 80},
	}

	for _, v := range values {
		bst.Insert(v.id, v.val)
	}

	var result []int
	bst.InOrder(func(node *BinaryNode[int]) {
		result = append(result, node.val)
	})

	expected := []int{20, 30, 40, 50, 60, 70, 80}
	s.Require().Equal(expected, result)
}

func (s *BSTTestSuite) TestPreOrderTraversal() {
	bst := NewBST[int]()
	values := []struct {
		id  uint64
		val int
	}{
		{1, 50},
		{2, 30},
		{3, 70},
		{4, 20},
		{5, 40},
		{6, 60},
		{7, 80},
	}

	for _, v := range values {
		bst.Insert(v.id, v.val)
	}

	var result []int
	bst.PreOrder(func(node *BinaryNode[int]) {
		result = append(result, node.val)
	})

	expected := []int{50, 30, 20, 40, 70, 60, 80}
	s.Require().Equal(expected, result)
}

func (s *BSTTestSuite) TestPostOrderTraversal() {
	bst := NewBST[int]()
	values := []struct {
		id  uint64
		val int
	}{
		{1, 50},
		{2, 30},
		{3, 70},
		{4, 20},
		{5, 40},
		{6, 60},
		{7, 80},
	}

	for _, v := range values {
		bst.Insert(v.id, v.val)
	}

	var result []int
	bst.PostOrder(func(node *BinaryNode[int]) {
		result = append(result, node.val)
	})

	expected := []int{20, 40, 30, 60, 80, 70, 50}
	s.Require().Equal(expected, result)
}

func (s *BSTTestSuite) TestHeight() {
	testCases := []struct {
		name       string
		insertVals []struct {
			id  uint64
			val int
		}
		expectedHeight int
	}{
		{
			name: "empty tree",
			insertVals: []struct {
				id  uint64
				val int
			}{},
			expectedHeight: -1,
		},
		{
			name: "single node",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
			},
			expectedHeight: 0,
		},
		{
			name: "balanced tree height 1",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
			},
			expectedHeight: 1,
		},
		{
			name: "balanced tree height 2",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 50},
				{2, 30},
				{3, 70},
				{4, 20},
				{5, 40},
				{6, 60},
				{7, 80},
			},
			expectedHeight: 2,
		},
		{
			name: "unbalanced tree (right-skewed)",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 10},
				{2, 20},
				{3, 30},
				{4, 40},
			},
			expectedHeight: 3,
		},
		{
			name: "unbalanced tree (left-skewed)",
			insertVals: []struct {
				id  uint64
				val int
			}{
				{1, 40},
				{2, 30},
				{3, 20},
				{4, 10},
			},
			expectedHeight: 3,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			bst := NewBST[int]()
			for _, v := range tc.insertVals {
				bst.Insert(v.id, v.val)
			}

			height := bst.Height()
			s.Require().Equal(tc.expectedHeight, height)
		})
	}
}

func (s *BSTTestSuite) TestClear() {
	bst := NewBST[int]()
	bst.Insert(1, 50)
	bst.Insert(2, 30)
	bst.Insert(3, 70)

	s.Require().Equal(3, bst.Size())
	s.Require().False(bst.IsEmpty())

	bst.Clear()

	s.Require().Equal(0, bst.Size())
	s.Require().True(bst.IsEmpty())
	s.Require().Nil(bst.Root())
}

func (s *BSTTestSuite) TestBinaryNodeHierarchy() {
	bst := NewBST[int]()
	bst.Insert(1, 50)
	bst.Insert(2, 30)
	bst.Insert(3, 70)

	root := bst.Root()
	s.Require().True(root.IsRoot())
	s.Require().False(root.IsLeft())
	s.Require().False(root.IsRight())

	left := root.Left()
	s.Require().NotNil(left)
	s.Require().False(left.IsRoot())
	s.Require().True(left.IsLeft())
	s.Require().False(left.IsRight())

	right := root.Right()
	s.Require().NotNil(right)
	s.Require().False(right.IsRoot())
	s.Require().False(right.IsLeft())
	s.Require().True(right.IsRight())
}

func (s *BSTTestSuite) TestBinaryNodeLevels() {
	bst := NewBST[int]()
	values := []struct {
		id  uint64
		val int
	}{
		{1, 50},
		{2, 30},
		{3, 70},
		{4, 20},
		{5, 40},
	}

	for _, v := range values {
		bst.Insert(v.id, v.val)
	}

	root := bst.Root()
	s.Require().Equal(0, root.Level())

	left := root.Left()
	s.Require().Equal(1, left.Level())

	right := root.Right()
	s.Require().Equal(1, right.Level())

	leftLeft := left.Left()
	s.Require().Equal(2, leftLeft.Level())

	leftRight := left.Right()
	s.Require().Equal(2, leftRight.Level())
}

func TestBST(t *testing.T) {
	suite.Run(t, new(BSTTestSuite))
}
