package mtree

import (
	"context"
	"math"
	"testing"
	"time"

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

func (s *NodeTestSuite) TestNode_Swap_TargetRoot() {
	model := HierarchyModel{
		RootTag: {"play"},
		"play":  {"intro", "ad", "game_scene"},
		"intro": {"clip1", "clip2"},
	}

	r, err := Hierarchy(
		model, 3,
		func() uint64 {
			return s.nextGroupID("hierarchy2")
		},
	)
	s.NotNil(r)
	s.Require().NoError(err)

	intro := func(n *Node[string]) bool {
		return n.Val() == "intro"
	}

	introNode, err := r.SelectOneChildFunc(intro)
	s.NotNil(introNode)
	s.Require().NoError(err)

	err = r.Swap(introNode)
	s.Require().NoError(err)
	r = introNode

	actualModel, err := ToModel(r)
	s.NotNil(actualModel)
	s.Require().NoError(err)

	expectedModel := HierarchyModel{
		RootTag: {"intro"},
		"intro": {"play", "ad", "game_scene"},
		"play":  {"clip1", "clip2"},
	}

	s.ElementsMatch(expectedModel[RootTag], actualModel[RootTag])
	s.ElementsMatch(expectedModel["intro"], actualModel["intro"])
	s.ElementsMatch(expectedModel["play"], actualModel["play"])
}

// Test error handling in NewNode options
func (s *NodeTestSuite) TestNewNode_ParentOpt_Nil() {
	id := s.nextDefaultGroupID()
	n, err := NewNode[int](id, 0, ParentOpt[int](nil))
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

func (s *NodeTestSuite) TestNewNode_ChildOpt_Nil() {
	id := s.nextDefaultGroupID()
	n, err := NewNode[int](id, 0, ChildOpt[int](nil))
	s.Nil(n)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

func (s *NodeTestSuite) TestNewNode_ParentOpt_MaxBreadthExceeded() {
	parentID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 0)
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0, ParentOpt[int](parent))
	s.Nil(child)
	s.Error(err)
	s.ErrorIs(err, ErrMaxBreadth)
}

func (s *NodeTestSuite) TestNewNode_ChildOpt_MaxBreadthExceeded() {
	childID1, childID2, parentID := s.nextDefaultGroupID(), s.nextDefaultGroupID(), s.nextDefaultGroupID()
	child1, err := NewNode[int](childID1, 0)
	s.NotNil(child1)
	s.Require().NoError(err)

	child2, err := NewNode[int](childID2, 0)
	s.NotNil(child2)
	s.Require().NoError(err)

	parent, err := NewNode[int](parentID, 1, ChildOpt[int](child1), ChildOpt[int](child2))
	s.Nil(parent)
	s.Error(err)
	s.ErrorIs(err, ErrMaxBreadth)
}

// Test AttachMany function
func (s *NodeTestSuite) TestNode_AttachMany() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	childIDs := []uint64{
		s.nextDefaultGroupID(),
		s.nextDefaultGroupID(),
		s.nextDefaultGroupID(),
	}

	children := make([]*Node[string], len(childIDs))
	for i, id := range childIDs {
		child, err := NewNode[string](id, 0, ValueOpt[string]("child"))
		s.NotNil(child)
		s.Require().NoError(err)
		children[i] = child
	}

	err = parent.AttachMany(children...)
	s.NoError(err)
	s.Equal(3, parent.Breadth())

	for _, child := range children {
		s.True(parent.HasChild(child))
		s.True(child.IsChildOf(parent))
	}
}

func (s *NodeTestSuite) TestNode_AttachMany_WithNilChildren() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	child1, err := NewNode[string](s.nextDefaultGroupID(), 0)
	s.NotNil(child1)
	s.Require().NoError(err)

	child2, err := NewNode[string](s.nextDefaultGroupID(), 0)
	s.NotNil(child2)
	s.Require().NoError(err)

	// Include nil children - they should be filtered out
	err = parent.AttachMany(child1, nil, child2, nil)
	s.NoError(err)
	s.Equal(2, parent.Breadth())
}

func (s *NodeTestSuite) TestNode_AttachMany_MaxBreadthExceeded() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 2)
	s.NotNil(parent)
	s.Require().NoError(err)

	children := make([]*Node[int], 3)
	for i := range children {
		children[i], err = NewNode[int](s.nextDefaultGroupID(), 0)
		s.NotNil(children[i])
		s.Require().NoError(err)
	}

	err = parent.AttachMany(children...)
	s.Error(err)
	s.ErrorIs(err, ErrMaxBreadth)
}

// Test DetachChild function
func (s *NodeTestSuite) TestNode_DetachChild() {
	parentID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 2)
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0, ParentOpt[int](parent))
	s.NotNil(child)
	s.Require().NoError(err)

	s.True(parent.HasChild(child))
	s.Equal(1, parent.Breadth())

	err = parent.DetachChild(child)
	s.NoError(err)
	s.False(parent.HasChild(child))
	s.Equal(0, parent.Breadth())
	s.True(child.IsDetached())
	s.False(child.HasParent())
}

func (s *NodeTestSuite) TestNode_DetachChild_Nil() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	err = parent.DetachChild(nil)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

func (s *NodeTestSuite) TestNode_DetachChild_NotFound() {
	parentID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0)
	s.NotNil(child)
	s.Require().NoError(err)

	err = parent.DetachChild(child)
	s.Error(err)
	s.ErrorIs(err, ErrNodeNotFound)
}

// Test DetachChildFunc
func (s *NodeTestSuite) TestNode_DetachChildFunc() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 5, ValueOpt[string]("parent"))
	s.NotNil(parent)
	s.Require().NoError(err)

	appleIDs := []uint64{s.nextDefaultGroupID(), s.nextDefaultGroupID()}
	orangeIDs := []uint64{s.nextDefaultGroupID(), s.nextDefaultGroupID()}

	for _, id := range appleIDs {
		apple, err := NewNode[string](id, 0, ValueOpt[string]("apple"), ParentOpt[string](parent))
		s.NotNil(apple)
		s.NoError(err)
	}

	for _, id := range orangeIDs {
		orange, err := NewNode[string](id, 0, ValueOpt[string]("orange"), ParentOpt[string](parent))
		s.NotNil(orange)
		s.NoError(err)
	}

	s.Equal(4, parent.Breadth())

	count := parent.DetachChildFunc(func(n *Node[string]) bool {
		return n.Val() == "apple"
	})

	s.Equal(2, count)
	s.Equal(2, parent.Breadth())
}

func (s *NodeTestSuite) TestNode_DetachChildFunc_NilFunc() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	count := parent.DetachChildFunc(nil)
	s.Equal(0, count)
}

// Test Move function
func (s *NodeTestSuite) TestNode_Move() {
	parent1ID, parent2ID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent1, err := NewNode[string](parent1ID, 2, ValueOpt[string]("parent1"))
	s.NotNil(parent1)
	s.Require().NoError(err)

	parent2, err := NewNode[string](parent2ID, 2, ValueOpt[string]("parent2"))
	s.NotNil(parent2)
	s.Require().NoError(err)

	child, err := NewNode[string](childID, 0, ValueOpt[string]("child"), ParentOpt[string](parent1))
	s.NotNil(child)
	s.Require().NoError(err)

	s.True(parent1.HasChild(child))
	s.False(parent2.HasChild(child))

	err = child.Move(parent2)
	s.NoError(err)

	s.False(parent1.HasChild(child))
	s.True(parent2.HasChild(child))
	s.Equal(parent2, child.Parent())
}

func (s *NodeTestSuite) TestNode_Move_NilParent() {
	childID := s.nextDefaultGroupID()
	child, err := NewNode[int](childID, 0)
	s.NotNil(child)
	s.Require().NoError(err)

	err = child.Move(nil)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

func (s *NodeTestSuite) TestNode_Move_MaxBreadthExceeded() {
	parent1ID, parent2ID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent1, err := NewNode[int](parent1ID, 1)
	s.NotNil(parent1)
	s.Require().NoError(err)

	parent2, err := NewNode[int](parent2ID, 0)
	s.NotNil(parent2)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0, ParentOpt[int](parent1))
	s.NotNil(child)
	s.Require().NoError(err)

	err = child.Move(parent2)
	s.Error(err)
	s.ErrorIs(err, ErrMaxBreadth)
}

// Test MoveChildren function
func (s *NodeTestSuite) TestNode_MoveChildren() {
	parent1ID, parent2ID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent1, err := NewNode[string](parent1ID, 3, ValueOpt[string]("parent1"))
	s.NotNil(parent1)
	s.Require().NoError(err)

	parent2, err := NewNode[string](parent2ID, 3, ValueOpt[string]("parent2"))
	s.NotNil(parent2)
	s.Require().NoError(err)

	childIDs := []uint64{s.nextDefaultGroupID(), s.nextDefaultGroupID()}
	for _, id := range childIDs {
		child, err := NewNode[string](id, 0, ValueOpt[string]("child"), ParentOpt[string](parent1))
		s.NotNil(child)
		s.NoError(err)
	}

	s.Equal(2, parent1.Breadth())
	s.Equal(0, parent2.Breadth())

	err = parent1.MoveChildren(parent2)
	s.NoError(err)

	s.Equal(0, parent1.Breadth())
	s.Equal(2, parent2.Breadth())
}

func (s *NodeTestSuite) TestNode_MoveChildren_NilParent() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	err = parent.MoveChildren(nil)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

func (s *NodeTestSuite) TestNode_MoveChildren_MaxBreadthExceeded() {
	parent1ID, parent2ID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent1, err := NewNode[int](parent1ID, 3)
	s.NotNil(parent1)
	s.Require().NoError(err)

	parent2, err := NewNode[int](parent2ID, 1)
	s.NotNil(parent2)
	s.Require().NoError(err)

	// Add 2 children to parent1
	for i := 0; i < 2; i++ {
		child, err := NewNode[int](s.nextDefaultGroupID(), 0, ParentOpt[int](parent1))
		s.NotNil(child)
		s.NoError(err)
	}

	err = parent1.MoveChildren(parent2)
	s.Error(err)
	s.ErrorIs(err, ErrMaxBreadth)
}

// Test SelectOneChildFunc no match
func (s *NodeTestSuite) TestNode_SelectOneChildFunc_NoMatch() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 2, ValueOpt[string]("parent"))
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string]("apple"), ParentOpt[string](parent))
	s.NotNil(child)
	s.NoError(err)

	result, err := parent.SelectOneChildFunc(func(n *Node[string]) bool {
		return n.Val() == "orange"
	})
	s.Nil(result)
	s.Error(err)
	s.ErrorIs(err, ErrNoMatch)
}

// Test SelectChildrenFunc no match
func (s *NodeTestSuite) TestNode_SelectChildrenFunc_NoMatch() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 2, ValueOpt[string]("parent"))
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string]("apple"), ParentOpt[string](parent))
	s.NotNil(child)
	s.NoError(err)

	results, err := parent.SelectChildrenFunc(func(n *Node[string]) bool {
		return n.Val() == "orange"
	})
	s.Nil(results)
	s.Error(err)
	s.ErrorIs(err, ErrNoMatch)
}

// Test SelectChildByID not found
func (s *NodeTestSuite) TestNode_SelectChildByID_NotFound() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	nonExistentID := s.nextDefaultGroupID()
	result, err := parent.SelectChildByID(nonExistentID)
	s.Nil(result)
	s.Error(err)
	s.ErrorIs(err, ErrNodeNotFound)
}

// Test Swap with nil target
func (s *NodeTestSuite) TestNode_Swap_NilTarget() {
	nodeID := s.nextDefaultGroupID()
	node, err := NewNode[int](nodeID, 1)
	s.NotNil(node)
	s.Require().NoError(err)

	err = node.Swap(nil)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

// Test HasChild with nil
func (s *NodeTestSuite) TestNode_HasChild_Nil() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	s.False(parent.HasChild(nil))
}

// Test IsChildOf with nil parent
func (s *NodeTestSuite) TestNode_IsChildOf_NilParent() {
	childID := s.nextDefaultGroupID()
	child, err := NewNode[int](childID, 0)
	s.NotNil(child)
	s.Require().NoError(err)

	s.False(child.IsChildOf(nil))
}

// Test Detach when already detached
func (s *NodeTestSuite) TestNode_Detach_AlreadyDetached() {
	nodeID := s.nextDefaultGroupID()
	node, err := NewNode[int](nodeID, 0)
	s.NotNil(node)
	s.Require().NoError(err)

	s.True(node.IsDetached())
	node.Detach() // Should not panic
	s.True(node.IsDetached())
}

// Test Capacity calculation
func (s *NodeTestSuite) TestNode_Capacity() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	s.Equal(5, parent.Capacity())

	for i := 0; i < 3; i++ {
		child, err := NewNode[int](s.nextDefaultGroupID(), 0, ParentOpt[int](parent))
		s.NotNil(child)
		s.NoError(err)
	}

	s.Equal(2, parent.Capacity())
}

// Test AttachChild error when max breadth reached
func (s *NodeTestSuite) TestNode_AttachChild_MaxBreadthReached() {
	parentID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 0)
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0)
	s.NotNil(child)
	s.Require().NoError(err)

	err = parent.AttachChild(child)
	s.Error(err)
	s.ErrorIs(err, ErrMaxBreadth)
}

// Test SelectOneChildByEachValue function
func (s *NodeTestSuite) TestNode_SelectOneChildByEachValue() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	// Create children with different values
	values := []string{"apple", "orange", "banana"}
	expectedChildren := make(map[string]*Node[string])

	for _, val := range values {
		child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string](val), ParentOpt[string](parent))
		s.NotNil(child)
		s.NoError(err)
		expectedChildren[val] = child
	}

	ctx := context.Background()
	result, err := parent.SelectOneChildByEachValue(ctx, "apple", "orange", "banana")
	s.NotNil(result)
	s.NoError(err)
	s.Len(result, 3)

	for val, expectedChild := range expectedChildren {
		actualChild, exists := result[val]
		s.True(exists, "Expected value %s not found in result", val)
		s.Equal(expectedChild, actualChild)
	}
}

// Test SelectOneChildByEachValue with duplicate values (should deduplicate)
func (s *NodeTestSuite) TestNode_SelectOneChildByEachValue_Duplicates() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 3)
	s.NotNil(parent)
	s.Require().NoError(err)

	// Create children
	values := []string{"apple", "orange"}
	for _, val := range values {
		child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string](val), ParentOpt[string](parent))
		s.NotNil(child)
		s.NoError(err)
	}

	ctx := context.Background()
	// Request with duplicates - should deduplicate
	result, err := parent.SelectOneChildByEachValue(ctx, "apple", "apple", "orange", "orange")
	s.NotNil(result)
	s.NoError(err)
	s.Len(result, 2) // Only 2 unique values
}

// Test SelectOneChildByEachValue with context cancellation
func (s *NodeTestSuite) TestNode_SelectOneChildByEachValue_ContextCanceled() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	// Create children
	for _, val := range []string{"apple", "orange", "banana"} {
		child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string](val), ParentOpt[string](parent))
		s.NotNil(child)
		s.NoError(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := parent.SelectOneChildByEachValue(ctx, "apple", "orange", "banana")
	s.Nil(result)
	s.Error(err)
	s.ErrorIs(err, context.Canceled)
}

// Test SelectOneChildByEachValue with timeout
func (s *NodeTestSuite) TestNode_SelectOneChildByEachValue_Timeout() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	// Create children
	for _, val := range []string{"apple", "orange"} {
		child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string](val), ParentOpt[string](parent))
		s.NotNil(child)
		s.NoError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure timeout occurs

	result, err := parent.SelectOneChildByEachValue(ctx, "apple", "orange")
	// Might succeed if fast enough, or fail with timeout
	if err != nil {
		s.ErrorIs(err, context.DeadlineExceeded)
	}
	_ = result
}

// Test SelectOneChildByEachValue with missing value
func (s *NodeTestSuite) TestNode_SelectOneChildByEachValue_NoMatch() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[string](parentID, 3)
	s.NotNil(parent)
	s.Require().NoError(err)

	// Create children
	child, err := NewNode[string](s.nextDefaultGroupID(), 0, ValueOpt[string]("apple"), ParentOpt[string](parent))
	s.NotNil(child)
	s.NoError(err)

	ctx := context.Background()
	result, err := parent.SelectOneChildByEachValue(ctx, "apple", "orange")
	s.Nil(result)
	s.Error(err)
}

// Test HasChildren when empty
func (s *NodeTestSuite) TestNode_HasChildren_Empty() {
	parentID := s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 5)
	s.NotNil(parent)
	s.Require().NoError(err)

	s.False(parent.HasChildren())
}

// Test IsChildOf when child belongs to different parent
func (s *NodeTestSuite) TestNode_IsChildOf_WrongParent() {
	parent1ID, parent2ID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent1, err := NewNode[int](parent1ID, 1)
	s.NotNil(parent1)
	s.Require().NoError(err)

	parent2, err := NewNode[int](parent2ID, 1)
	s.NotNil(parent2)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0, ParentOpt[int](parent1))
	s.NotNil(child)
	s.Require().NoError(err)

	s.True(child.IsChildOf(parent1))
	s.False(child.IsChildOf(parent2))
}

// Test HasChild with empty children map
func (s *NodeTestSuite) TestNode_HasChild_EmptyChildren() {
	parentID, childID := s.nextDefaultGroupID(), s.nextDefaultGroupID()
	parent, err := NewNode[int](parentID, 1)
	s.NotNil(parent)
	s.Require().NoError(err)

	child, err := NewNode[int](childID, 0)
	s.NotNil(child)
	s.Require().NoError(err)

	s.False(parent.HasChild(child))
}
