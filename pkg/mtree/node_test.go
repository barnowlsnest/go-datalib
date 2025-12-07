package mtree

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

const testDefaultGroup = "default_test"

type NodeTestSuite struct {
	suite.Suite
	seq *serial.Serial
}

func TestNodeTestSuite(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

func (s *NodeTestSuite) SetupTest() {
	s.seq = serial.Seq()
}

func (s *NodeTestSuite) nextGroupID(name string) uint64 {
	s.T().Helper()
	return s.seq.Next(name)
}

func (s *NodeTestSuite) nextDefaultGroupID() uint64 {
	return s.nextGroupID(testDefaultGroup)
}

func (s *NodeTestSuite) TestNewNode() {
	id := s.nextDefaultGroupID()
	n, err := NewNode[int](id, 0)
	s.NotNil(n)
	s.NoError(err)
	s.Equal(0, n.Val())
	s.Equal(-1, n.Level())
	s.True(n.IsDetached())
	s.False(n.IsRoot())
	s.False(n.IsAttached())
	s.False(n.HasParent())
}

func (s *NodeTestSuite) TestNewNode_ID() {
	id := s.nextDefaultGroupID()
	n, err := NewNode[int](id, 0)
	s.NotNil(n)
	s.NoError(err)
	s.Equal(id, n.ID())
}

func (s *NodeTestSuite) TestNewNode_MaxBreadth() {
	expectedMaxBreadth := 3
	id := s.nextDefaultGroupID()
	n, err := NewNode[int](id, expectedMaxBreadth)
	s.NotNil(n)
	s.NoError(err)
	s.Equal(expectedMaxBreadth, n.MaxBreadth())
	s.Equal(expectedMaxBreadth, n.Capacity())
	s.Equal(0, n.Breadth())
}

func (s *NodeTestSuite) TestNewNode_ChildOpt() {
	// create children nodes
	childID1, childID2 := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	child1, err := NewNode[string](childID1, 0, ValueOpt("child1"))
	s.NotNil(child1)
	s.Require().NoError(err)

	child2, err := NewNode[string](childID2, 0, ValueOpt("child2"))
	s.NotNil(child2)
	s.Require().NoError(err)

	// create a parent node with children
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](
		parentID, 2,
		ValueOpt("parent"),
		ChildOpt[string](child1),
		ChildOpt[string](child2),
	)
	s.NotNil(parent)
	s.Require().NoError(err)

	// assert parentness
	s.True(child1.IsChildOf(parent))
	s.True(child1.HasParent())
	s.True(child1.IsAttached())
	s.Equal(child1.Parent(), parent)

	s.True(child2.IsChildOf(parent))
	s.True(child2.HasParent())
	s.True(child2.IsAttached())
	s.Equal(child2.Parent(), parent)

	// assert children
	s.False(parent.HasParent())
	s.Equal(0, parent.Level())
	s.Equal(1, child1.Level())
	s.Equal(1, child2.Level())

	actualChildren := make(map[uint64]*Node[string])
	for id, child := range parent.ChildrenIter() {
		actualChildren[id] = child
	}

	s.Equal(
		map[uint64]*Node[string]{
			serial.NSum(parentID, childID1): child1,
			serial.NSum(parentID, childID2): child2,
		},
		actualChildren,
	)
}

func (s *NodeTestSuite) TestNewNode_ValueOpt() {
	expectedVal := math.Pi
	id := s.nextDefaultGroupID()
	n, err := NewNode[float64](id, 1, ValueOpt[float64](expectedVal))
	s.NotNil(n)
	s.Require().NoError(err)
	s.Equal(expectedVal, n.Val())
}

func (s *NodeTestSuite) TestNewNode_LevelOpt() {
	expectedLevel := 10
	id := s.nextDefaultGroupID()
	n, err := NewNode[int](id, 1, LevelOpt[int](expectedLevel))
	s.NotNil(n)
	s.Require().NoError(err)
	s.Equal(expectedLevel, n.Level())
}

func (s *NodeTestSuite) TestNewNode_ParentOpt() {
	childID, parentID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0, ParentOpt[int](parent))
	s.NotNil(child)
	s.Require().NoError(err)
	s.True(child.IsChildOf(parent))
	s.True(child.HasParent())
	s.True(parent.HasChildren())
	s.True(parent.HasChild(child))
}

func (s *NodeTestSuite) TestNode_WithValue() {
	id := s.nextDefaultGroupID()
	n, err := NewNode[float64](id, 0)
	s.NotNil(n)
	s.Require().NoError(err)
	s.Equal(0., n.Val())
	expectedVal := math.Pi
	n.WithValue(expectedVal)
	s.Equal(expectedVal, n.Val())
}

func (s *NodeTestSuite) TestNode_SelectChildByID() {
	parentID, childID1, childID2 := s.nextDefaultGroupID(), s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 2)
	s.NotNil(parent)
	s.Require().NoError(err)

	child1, err := NewNode[int](childID1, 0, ParentOpt[int](parent))
	s.NotNil(child1)
	s.Require().NoError(err)

	child2, err := NewNode[int](childID2, 0, ParentOpt[int](parent))
	s.NotNil(child2)
	s.Require().NoError(err)

	selectedChild1, err := parent.SelectChildByID(childID1)
	s.NotNil(selectedChild1)
	s.Require().NoError(err)

	selectedChild2, err := parent.SelectChildByID(childID2)
	s.NotNil(selectedChild2)
	s.Require().NoError(err)

	s.Equal(child1, selectedChild1)
	s.Equal(child2, selectedChild2)
}

func (s *NodeTestSuite) TestNode_SelectChildrenFunc() {
	boxID := s.nextDefaultGroupID()
	box, err := NewNode[string](boxID, 5, ValueOpt[string]("boxOfFruits"))
	s.NotNil(box)
	s.Require().NoError(err)

	appleIDs := []uint64{s.nextDefaultGroupID(), s.nextDefaultGroupID()}
	orangeIDs := []uint64{s.nextDefaultGroupID(), s.nextDefaultGroupID(), s.nextDefaultGroupID()}

	expectedApples := make([]*Node[string], len(appleIDs))
	for i, id := range appleIDs {
		apple, err := NewNode[string](id, 0, ValueOpt[string]("apple"), ParentOpt[string](box))
		s.NotNil(apple)
		s.NoError(err)
		expectedApples[i] = apple
	}

	expectedOranges := make([]*Node[string], len(orangeIDs))
	for i, id := range orangeIDs {
		orange, err := NewNode[string](id, 0, ValueOpt[string]("orange"), ParentOpt[string](box))
		s.NotNil(orange)
		s.NoError(err)
		expectedOranges[i] = orange
	}

	actualOranges, err := box.SelectChildrenFunc(func(fruit *Node[string]) bool {
		return fruit.Val() == "orange"
	})
	s.NotNil(actualOranges)
	s.Require().NoError(err)
	s.ElementsMatch(expectedOranges, actualOranges)

	actualApples, err := box.SelectChildrenFunc(func(fruit *Node[string]) bool {
		return fruit.Val() == "apple"
	})
	s.NotNil(actualOranges)
	s.Require().NoError(err)
	s.ElementsMatch(expectedApples, actualApples)
}

func (s *NodeTestSuite) TestNode_Swap_Once() {
	model := HierarchyModel{
		RootTag: {"CEO"},
		"CEO":   {"CTO", "CFO"},
		"CTO":   {"PSA", "PSE", "DM"},
		"CFO":   {"SEM", "PA"},
	}

	ctoFn := func(n *Node[string]) bool {
		return n.Val() == "CTO"
	}

	dmFn := func(n *Node[string]) bool {
		return n.Val() == "DM"
	}

	cfoFn := func(n *Node[string]) bool {
		return n.Val() == "CFO"
	}

	semFn := func(n *Node[string]) bool {
		return n.Val() == "SEM"
	}

	ceo, err := Hierarchy(
		model, 10,
		func() uint64 {
			return s.nextGroupID("hierarchy1")
		},
	)
	s.NotNil(ceo)
	s.Require().NoError(err)

	cto, err := ceo.SelectChildrenFunc(ctoFn)
	s.NotNil(cto)
	s.Require().NoError(err)
	s.Require().Len(cto, 1)

	dm, err := cto[0].SelectChildrenFunc(dmFn)
	s.NotNil(dm)
	s.Require().NoError(err)
	s.Require().Len(dm, 1)

	cfo, err := ceo.SelectChildrenFunc(cfoFn)
	s.NotNil(cfo)
	s.Require().NoError(err)
	s.Require().Len(cfo, 1)

	sem, err := cfo[0].SelectChildrenFunc(semFn)
	s.NotNil(sem)
	s.Require().NoError(err)
	s.Require().Len(sem, 1)

	err = sem[0].Swap(dm[0])
	s.Require().NoError(err)

	expectedModel := HierarchyModel{
		RootTag: {"CEO"},
		"CEO":   {"CTO", "CFO"},
		"CTO":   {"PSA", "PSE", "SEM"},
		"CFO":   {"DM", "PA"},
	}

	actualModel, err := ToModel(ceo)
	s.NotNil(actualModel)
	s.Require().NoError(err)

	s.ElementsMatch(expectedModel["CTO"], actualModel["CTO"])
	s.ElementsMatch(expectedModel["CFO"], actualModel["CFO"])
}
