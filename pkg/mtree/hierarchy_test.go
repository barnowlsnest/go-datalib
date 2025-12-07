package mtree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type HierarchyTestSuite struct {
	suite.Suite
}

func TestHierarchyTestSuite(t *testing.T) {
	suite.Run(t, new(HierarchyTestSuite))
}

func nextID() uint64 {
	return serial.Seq().Next("hierarchyTest")
}

func (s *HierarchyTestSuite) TestHierarchy() {
	expectedModel := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"B", "C", "D"},
		"D":     ChildrenSlice{"E", "F"},
	}

	n, err := Hierarchy(expectedModel, 5, nextID)
	s.NotNil(n)
	s.NoError(err)

	actualModel, err := ToModel(n)
	s.NotNil(actualModel)
	s.NoError(err)
	s.ElementsMatch(expectedModel["A"], actualModel["A"])
	s.ElementsMatch(expectedModel["D"], actualModel["D"])
	s.ElementsMatch(expectedModel[RootTag], actualModel[RootTag])
}

// Test cyclic reference detection - direct cycle
func (s *HierarchyTestSuite) TestHierarchy_DirectCycle() {
	cyclicModel := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"B"},
		"B":     ChildrenSlice{"A"}, // B points back to A - direct cycle
	}

	done := make(chan struct{})
	var n *Node[string]
	var err error

	go func() {
		defer close(done)
		n, err = Hierarchy(cyclicModel, 5, nextID)
	}()

	select {
	case <-done:
		// Should return an error, not hang
		s.Nil(n)
		s.Error(err)
		s.Contains(err.Error(), "cycle")
	case <-time.After(2 * time.Second):
		s.Fail("Hierarchy() hung on cyclic model - infinite loop detected!")
	}
}

// Test cyclic reference detection - indirect cycle
func (s *HierarchyTestSuite) TestHierarchy_IndirectCycle() {
	cyclicModel := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"B"},
		"B":     ChildrenSlice{"C"},
		"C":     ChildrenSlice{"A"}, // C points back to A - indirect cycle
	}

	done := make(chan struct{})
	var n *Node[string]
	var err error

	go func() {
		defer close(done)
		n, err = Hierarchy(cyclicModel, 5, nextID)
	}()

	select {
	case <-done:
		// Should return an error, not hang
		s.Nil(n)
		s.Error(err)
		s.Contains(err.Error(), "cycle")
	case <-time.After(2 * time.Second):
		s.Fail("Hierarchy() hung on cyclic model - infinite loop detected!")
	}
}

// Test self-reference
func (s *HierarchyTestSuite) TestHierarchy_SelfReference() {
	selfRefModel := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"A"}, // A references itself
	}

	done := make(chan struct{})
	var n *Node[string]
	var err error

	go func() {
		defer close(done)
		n, err = Hierarchy(selfRefModel, 5, nextID)
	}()

	select {
	case <-done:
		// Should return an error, not hang
		s.Nil(n)
		s.Error(err)
		s.Contains(err.Error(), "cycle")
	case <-time.After(2 * time.Second):
		s.Fail("Hierarchy() hung on self-referencing model - infinite loop detected!")
	}
}

// Test missing root tag
func (s *HierarchyTestSuite) TestHierarchy_MissingRootTag() {
	model := HierarchyModel{
		"A": ChildrenSlice{"B", "C"},
	}

	n, err := Hierarchy(model, 5, nextID)
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrRootTagNotFound)
}

// Test multiple roots
func (s *HierarchyTestSuite) TestHierarchy_MultipleRoots() {
	model := HierarchyModel{
		RootTag: ChildrenSlice{"A", "B"}, // Multiple roots not allowed
		"A":     ChildrenSlice{"C"},
		"B":     ChildrenSlice{"D"},
	}

	n, err := Hierarchy(model, 5, nextID)
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrHierarchyModel)
}

// Test root reference not found
func (s *HierarchyTestSuite) TestHierarchy_RootRefNotFound() {
	model := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		// "A" is not defined in the model
	}

	n, err := Hierarchy(model, 5, nextID)
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrHierarchyModel)
}

// Test nil nextID function
func (s *HierarchyTestSuite) TestHierarchy_NilNextID() {
	model := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"B"},
	}

	n, err := Hierarchy(model, 5, nil)
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

// Test maxBreadth less than 1
func (s *HierarchyTestSuite) TestHierarchy_InvalidMaxBreadth() {
	model := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"B"},
	}

	n, err := Hierarchy(model, 0, nextID)
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrHierarchyModel)
}

// Test ToModel with nil node
func (s *HierarchyTestSuite) TestToModel_NilNode() {
	model, err := ToModel(nil)
	s.Nil(model)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

// Test ToModel with non-root node
func (s *HierarchyTestSuite) TestToModel_NonRootNode() {
	parentID, childID := nextID(), nextID()
	parent, err := NewNode[string](parentID, 1, ValueOpt[string]("parent"))
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[string](childID, 0, ValueOpt[string]("child"), ParentOpt[string](parent))
	s.NotNil(child)
	s.Require().NoError(err)

	model, err := ToModel(child)
	s.Nil(model)
	s.Error(err)
	s.ErrorIs(err, ErrHierarchyModel)
}
