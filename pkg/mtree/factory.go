package mtree

import (
	"fmt"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type NodeFactory[T comparable] struct {
	shard      string
	maxBreadth int
	serialID   *serial.Serial
}

func (factory *NodeFactory[T]) id() uint64 {
	return factory.serialID.Next(factory.shard)
}

func NewNodeFactory[T comparable](shard string, maxBreadth int, id *serial.Serial) *NodeFactory[T] {
	return &NodeFactory[T]{
		shard:      shard,
		serialID:   id,
		maxBreadth: maxBreadth,
	}
}

func (factory *NodeFactory[T]) Node(opts ...NodeOption[T]) *Node[T] {
	return NewNode[T](factory.id(), factory.maxBreadth, opts...)
}

func (factory *NodeFactory[T]) Root(value T) *Node[T] {
	return NewRoot[T](factory.id(), factory.maxBreadth, value)
}

func (factory *NodeFactory[T]) RootWithChildren(value T, children ...*Node[T]) (*Node[T], error) {
	r := factory.Root(value)
	if err := r.AttachMany(children...); err != nil {
		return nil, fmt.Errorf("factory err: %w", err)
	}

	return r, nil
}
