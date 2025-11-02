package tree

import (
	"cmp"
	"fmt"

	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
	"github.com/barnowlsnest/go-datalib/pkg/utils"
)

type (
	N[T cmp.Ordered] struct {
		ID  uint64
		Val T
	}

	LevelParams[T cmp.Ordered] struct {
		Node   *Node[T]
		MinVal T
		MaxVal T
	}

	LevelFn[T cmp.Ordered] func(params LevelParams[T]) bool

	Nary[T cmp.Ordered] struct {
		root               *Node[T]
		nodes              map[uint64]*Node[T]
		levels             map[uint8]*list.LinkedList
		levelsCount        uint8
		maxDepth           uint8
		maxChildrenPerNode uint8
	}
)

func NewNary[T cmp.Ordered](maxDepth, maxNodesPerNode uint8) *Nary[T] {
	return &Nary[T]{
		levels:             make(map[uint8]*list.LinkedList, maxDepth),
		nodes:              make(map[uint64]*Node[T], maxDepth*maxNodesPerNode),
		maxDepth:           maxDepth,
		maxChildrenPerNode: maxNodesPerNode,
	}
}

func NewBinary[T cmp.Ordered](maxDepth uint8) *Nary[T] {
	return NewNary[T](maxDepth, 2)
}

func NewTernary[T cmp.Ordered](maxDepth uint8) *Nary[T] {
	return NewNary[T](maxDepth, 3)
}

func (t *Nary[T]) toNodes(nNodes ...N[T]) []*Node[T] {
	nodes := make([]*Node[T], 0, len(nNodes))
	for _, n := range nNodes {
		tNode := NewNode[T](n.ID, n.Val)
		t.nodes[n.ID] = tNode
		nodes = append(nodes, tNode)
	}

	return nodes
}

func (t *Nary[T]) AddRoot(n N[T]) error {
	if t.root != nil {
		return ErrNotAllowed
	}

	nodes := t.toNodes(n)
	t.root = nodes[0]
	t.root.BeholdRoot()
	rootLevel := list.New()
	rootLevel.Push(node.ID(n.ID))
	t.levels[0] = rootLevel

	return nil
}

func (t *Nary[T]) AddChildren(pID uint64, nNodes ...N[T]) error {
	parent, exists := t.nodes[pID]
	if !exists {
		return ErrNilParent
	}

	if t.levelsCount >= t.maxDepth {
		return ErrNotAllowedMaxDepth
	}

	nNodesLen := len(nNodes)

	if nNodesLen > int(t.maxChildrenPerNode) {
		return ErrNotAllowedMaxNodes
	}

	pChildren := parent.Children()
	if pChildren != nil && pChildren.Size()+nNodesLen > int(t.maxChildrenPerNode) {
		return ErrNotAllowedMaxNodes
	}

	nodes := t.toNodes(nNodes...)
	children, err := NewNodeChildren[T](parent, nodes...)
	if err != nil {
		return err
	}

	parent.WithChildren(children)
	nextLevel := parent.Level() + 1
	uint8Level, errConv := utils.SafeIntToUint8(nextLevel)
	if errConv != nil {
		return fmt.Errorf("err level: %w", errConv)
	}

	t.levelsCount = uint8Level
	levelNodes := list.New()
	for _, n := range nodes {
		levelNodes.Push(node.ID(n.ID()))
	}

	t.levels[uint8Level] = levelNodes

	return nil
}

func (t *Nary[T]) Level(level uint8) ([]*Node[T], error) {
	levelNodes, exists := t.levels[level]
	if !exists {
		return nil, ErrLevelNotFound
	}

	nodes := make([]*Node[T], 0, levelNodes.Size())
	for _, n := range levelNodes.IterNext() {
		treeNode, treeNodeExists := t.nodes[n.ID()]
		if !treeNodeExists {
			return nil, ErrUnexpected
		}

		nodes = append(nodes, treeNode)
	}

	return nodes, nil
}

func (t *Nary[T]) LevelFunc(level uint8, fn LevelFn[T]) error {
	nodes, err := t.Level(level)
	if err != nil {
		return err
	}

	values := make([]T, 0, len(nodes))
	for _, n := range nodes {
		values = append(values, n.Value())
	}

	minVal, maxVal, err := utils.MinMax(values)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		var shouldCont bool
		err := func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = ErrUnexpected
					return
				}
			}()
			shouldCont = fn(LevelParams[T]{
				Node:   n,
				MinVal: minVal,
				MaxVal: maxVal,
			})
			return nil
		}()

		if err != nil {
			return err
		}

		if !shouldCont {
			break
		}
	}

	return nil
}
