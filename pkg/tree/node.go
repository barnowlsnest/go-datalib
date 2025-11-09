package tree

import (
	"cmp"

	"github.com/barnowlsnest/go-datalib/pkg/node"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type (
	NodeValue[T cmp.Ordered] struct {
		node *node.Node
		val  T
	}

	NodeProps[T cmp.Ordered] struct {
		ID    uint64
		Value T
	}
)

func newNodeValue[T cmp.Ordered](n *node.Node, v T) *NodeValue[T] {
	return &NodeValue[T]{n, v}
}

func NewNode[T cmp.Ordered](id uint64) *NodeValue[T] {
	var val T
	return newNodeValue(node.ID(id), val)
}

func NewNodeValue[T cmp.Ordered](id uint64, value T) *NodeValue[T] {
	return newNodeValue(node.ID(id), value)
}

func Node[T cmp.Ordered](id uint64, value T) *NodeValue[T] {
	return NewNodeValue(id, value)
}

func (n *NodeValue[T]) Props() (NodeProps[T], error) {
	if n.node == nil {
		return NodeProps[T]{}, ErrNil
	}

	return NodeProps[T]{n.node.ID(), n.val}, nil
}

func (n *NodeValue[T]) WithValue(newVal T) {
	n.val = newVal
}

func (n *NodeValue[T]) WithParent(parent *NodeValue[T]) (uint64, error) {
	if n.node == nil {
		return 0, ErrNil
	}

	if parent == nil {
		return 0, ErrParentNil
	}

	n.node.WithPrev(parent.node)

	return serial.NSum(parent.node.ID(), n.node.ID()), nil
}

func (n *NodeValue[T]) HasParent() bool {
	if n.node == nil {
		return false
	}

	return n.node.Prev() != nil
}

func (n *NodeValue[T]) IsChildOf(parent *NodeValue[T]) bool {
	if n.node == nil {
		return false
	}

	if parent == nil {
		return false
	}

	if parent.node == nil {
		return false
	}

	return n.node.Prev() == parent.node
}

func (n *NodeValue[T]) UnlinkParent() {
	if n.node == nil {
		return
	}

	n.node.WithPrev(nil)
}

func (n *NodeValue[T]) Equal(other *NodeValue[T]) bool {
	if n == nil {
		return other == nil
	}

	if other == nil {
		return false
	}

	if n.val != other.val {
		return false
	}

	switch {
	case n.val != other.val:
		return false
	case n.node != nil && other.node == nil:
		return false
	case n.node == nil && other.node != nil:
		return false
	case n.node == nil && other.node == nil:
		return true
	default:
		return n.node.ID() == other.node.ID()
	}
}
