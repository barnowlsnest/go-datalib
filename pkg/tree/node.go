package tree

import (
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

type Node[T comparable] struct {
	*node.Node
	value    T
	level    int
	children *NodeChildren[T]
}

func NewNode[T comparable](id uint64, value T) *Node[T] {
	return NewNodeWithChildren(id, value, nil)
}

func NewNodeWithChildren[T comparable](id uint64, value T, nodes *NodeChildren[T]) *Node[T] {
	return &Node[T]{
		Node:     node.ID(id),
		value:    value,
		children: nodes,
		level:    -1,
	}
}

func (n *Node[T]) Value() T {
	return n.value
}

func (n *Node[T]) Level() int {
	return n.level
}

func (n *Node[T]) IsRoot() bool {
	return n.level == 0 && n.Prev() == nil
}

func (n *Node[T]) BeholdRoot() {
	n.WithPrev(nil)
	n.level = 0
}

func (n *Node[T]) WithParent(parent *Node[T]) {
	if parent == nil {
		return
	}

	n.WithPrev(parent.Node)
	n.level = parent.level + 1
}

func (n *Node[T]) WithChildren(children *NodeChildren[T]) {
	n.children = children
}

func (n *Node[T]) Children() *NodeChildren[T] {
	return n.children
}
