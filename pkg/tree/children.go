package tree

import (
	"iter"
	"slices"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

const (
	ParentRel = "parent"
	ChildRel  = "child"
	UnRelated = "unrelated"
)

type (
	NodeChildren[T comparable] struct {
		parent   *Node[T]
		children map[uint64]*ChildNode[T]
		index    map[int]*ChildNode[T]
	}

	ChildNode[T comparable] struct {
		hash   uint64
		parent *Node[T]
		node   *Node[T]
	}
)

func NewNodeChildren[T comparable](parent *Node[T], nodes ...*Node[T]) (*NodeChildren[T], error) {
	if parent == nil {
		return nil, ErrNilParent
	}

	notNilNodes := slices.DeleteFunc(nodes, func(n *Node[T]) bool {
		return n == nil
	})

	m := map[uint64]*ChildNode[T]{}
	inx := map[int]*ChildNode[T]{}
	for index, node := range notNilNodes {
		if index == 0 {
			parent.ptr.WithNext(node.ptr)
		}

		node.WithParent(parent)
		child := newChild(parent, node)
		inx[index] = child
		m[node.ID()] = child
	}

	return &NodeChildren[T]{parent, m, inx}, nil
}

func (nc *NodeChildren[T]) Size() int {
	return len(nc.children)
}

func (nc *NodeChildren[T]) Parent() *Node[T] {
	return nc.parent
}

func (nc *NodeChildren[T]) Child(id uint64) (*ChildNode[T], error) {
	child, exists := nc.children[id]
	if !exists {
		return nil, ErrChildNotFound
	}

	return child, nil
}

func (nc *NodeChildren[T]) ChildNth(n int) (*ChildNode[T], error) {
	if n < 0 {
		return nil, ErrChildNotFound
	}

	if n >= len(nc.children) {
		return nil, ErrChildNotFound
	}

	child, exists := nc.index[n]
	if !exists {
		return nil, ErrChildNotFound
	}

	return child, nil
}

func (nc *NodeChildren[T]) HasChild(id uint64) bool {
	_, exists := nc.children[id]
	return exists
}

func (nc *NodeChildren[T]) Relation(id uint64) string {
	switch {
	case nc.parent.ID() == id:
		return ParentRel
	case nc.HasChild(id):
		return ChildRel
	default:
		return UnRelated
	}
}

func (nc *NodeChildren[T]) Nodes() iter.Seq2[uint64, *ChildNode[T]] {
	return func(yield func(uint64, *ChildNode[T]) bool) {
		for _, child := range nc.children {
			if !yield(child.hash, child) {
				return
			}
		}
	}
}

func newChild[T comparable](parent, node *Node[T]) *ChildNode[T] {
	return &ChildNode[T]{serial.NSum(parent.ID(), node.ID()), parent, node}
}

func (child *ChildNode[T]) Hash() uint64 {
	return child.hash
}

func (child *ChildNode[T]) Parent() *Node[T] {
	return child.parent
}

func (child *ChildNode[T]) Node() *Node[T] {
	return child.node
}
