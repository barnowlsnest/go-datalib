package mtree

import (
	"errors"

	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

type (
	ChildrenSlice  = []string
	HierarchyModel = map[string]ChildrenSlice
)

const RootTag = "#root"

func Hierarchy(m HierarchyModel, maxBreadth int, nextID func() uint64) (*Node[string], error) {
	rootDef, rootDefined := m[RootTag]
	switch {
	case !rootDefined:
		return nil, ErrRootTagNotFound
	case len(rootDef) != 1:
		return nil, errors.Join(ErrHierarchyModel, errors.New("only 1 root allowed"))
	case nextID == nil:
		return nil, ErrNil
	case maxBreadth < 1:
		return nil, errors.Join(ErrHierarchyModel, errors.New("max breadth should be at least 1"))
	}

	rootNodeVal := rootDef[0]
	rootChildren, rootExists := m[rootNodeVal]
	if !rootExists {
		return nil, errors.Join(ErrHierarchyModel, errors.New("root ref not found"))
	}

	rootID := nextID()
	rootNode, errRoot := NewNode[string](rootID, maxBreadth, ValueOpt[string](rootNodeVal))
	if errRoot != nil {
		return nil, errRoot
	}
	if ok := rootNode.asRoot(); !ok {
		return nil, errors.Join(ErrHierarchyModel, errors.New("unable set root state"))
	}

	stack := list.NewStack()
	lookup := make(map[uint64]*Node[string])
	var (
		parent   *Node[string]
		children []string
	)
	parent = rootNode
	children = rootChildren
buildBranch:
	for _, childVal := range children {
		childID := nextID()
		childNode, errChild := NewNode[string](childID, maxBreadth, ValueOpt[string](childVal))
		if errChild != nil {
			return nil, errChild
		}
		if errAttach := parent.AttachChild(childNode); errAttach != nil {
			return nil, errAttach
		}

		lookup[childID] = childNode
		stack.Push(node.ID(childID))
	}

	if stack.IsEmpty() {
		return rootNode, nil
	}

	n := stack.Pop()
	if childNode := lookup[n.ID()]; childNode != nil {
		parent = childNode
		children = m[childNode.Val()]
		goto buildBranch
	} else {
		return nil, ErrNil
	}
}

func ToModel(n *Node[string]) (HierarchyModel, error) {
	if n == nil {
		return nil, ErrNil
	}
	if ok := n.IsRoot(); !ok {
		return nil, errors.Join(ErrHierarchyModel, errors.New("not root"))
	}

	m, lookup := make(HierarchyModel), make(map[uint64]*Node[string])
	m[RootTag] = ChildrenSlice{n.Val()}
	rootID := n.ID()
	lookup[rootID] = n
	queue := list.NewQueue()
	queue.Enqueue(node.ID(rootID))
	for !queue.IsEmpty() {
		next := queue.Dequeue()
		if next == nil {
			return nil, ErrHierarchyModel
		}

		nextNode := lookup[next.ID()]
		if nextNode == nil {
			return nil, ErrNil
		}

		for id, child := range nextNode.ChildrenIter() {
			queue.Enqueue(node.ID(id))
			lookup[id] = child
			if m[nextNode.Val()] == nil {
				m[nextNode.Val()] = make(ChildrenSlice, 0, nextNode.MaxBreadth())
			}
			m[nextNode.Val()] = append(m[nextNode.Val()], child.Val())
		}
	}

	return m, nil
}
