package tree

import (
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

type Node[T comparable] struct {
	value    T
	level    int
	ptr      *node.Node
	children *NodeChildren[T]
}

func NewNode[T comparable](id uint64, value T, nodes *NodeChildren[T]) *Node[T] {
	return &Node[T]{
		value:    value,
		level:    -1,
		ptr:      node.ID(id),
		children: nodes,
	}
}

func (n *Node[T]) Value() T {
	return n.value
}

func (n *Node[T]) Level() int {
	return n.level
}

func (n *Node[T]) ID() uint64 {
	return n.ptr.ID()
}

func (n *Node[T]) IsRoot() bool {
	return n.level == 0 && n.ptr.Prev() == nil
}

func (n *Node[T]) BeholdRoot() {
	n.ptr.WithPrev(nil)
	n.level = 0
}

func (n *Node[T]) WithParent(parent *Node[T]) {
	if parent == nil {
		return
	}

	n.ptr.WithPrev(parent.ptr)
	n.level = parent.level + 1
}

func (n *Node[T]) Ptr() *node.Node {
	return n.ptr
}

func (n *Node[T]) Children() *NodeChildren[T] {
	return n.children
}
