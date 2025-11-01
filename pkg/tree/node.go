package tree

import (
	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

type Node[T comparable] struct {
	value    T
	level    int
	pNode    *node.Node
	children *list.LinkedList
}

func WithChildren[T comparable](id uint64, value T, children *list.LinkedList) *Node[T] {
	return &Node[T]{
		value:    value,
		level:    -1,
		pNode:    node.New(id, nil, nil),
		children: children,
	}
}

func New[T comparable](id uint64, value T) *Node[T] {
	return &Node[T]{
		value:    value,
		level:    -1,
		pNode:    node.New(id, nil, nil),
		children: list.New(),
	}
}

func (n *Node[T]) Value() T {
	return n.value
}

func (n *Node[T]) Level() int {
	return n.level
}

func (n *Node[T]) ID() uint64 {
	return n.pNode.ID()
}

func (n *Node[T]) IsRoot() bool {
	return n.level == 0 && n.pNode.Prev() == nil
}

func (n *Node[T]) BeholdRoot() {
	n.pNode.WithPrev(nil)
	n.level = 0
}

func (n *Node[T]) withParent(parent *Node[T]) error {
	if parent == nil {
		return ErrParentNil
	}
	
	n.pNode.WithPrev(parent.pNode)
	n.level = parent.level + 1
	
	return nil
}

func (n *Node[T]) PushHead(parent *Node[T]) error {
	if err := n.withParent(parent); err != nil {
		return err
	}
	
	parent.children.Unshift(n.pNode)
	
	return nil
}

func (n *Node[T]) PushTail(parent *Node[T]) error {
	if err := n.withParent(parent); err != nil {
		return err
	}
	
	parent.children.Push(n.pNode)
	
	return nil
}
