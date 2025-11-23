package mtree

import (
	"fmt"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

const DefaultShard = "default"

type Factory[T comparable] struct {
	maxDepth   int
	maxBreadth int
	serialID   *serial.Serial
}

func (factory *Factory[T]) id(shard string) uint64 {
	return factory.serialID.Next(shard)
}

func NewFactory[T comparable](maxBreadth, maxDepth int, id *serial.Serial) *Factory[T] {
	return &Factory[T]{
		serialID:   id,
		maxBreadth: maxBreadth,
		maxDepth:   maxDepth,
	}
}

func (factory *Factory[T]) Node(shard string, opts ...NodeOption[T]) *Node[T] {
	return NewNode[T](shard, factory.id(shard), factory.maxBreadth, opts...)
}

func (factory *Factory[T]) Root(shard string, value T) *Node[T] {
	return NewRoot[T](shard, factory.id(shard), factory.maxBreadth, value)
}

func (factory *Factory[T]) RootWithChildren(shard string, value T, children ...*Node[T]) (*Node[T], error) {
	r := factory.Root(shard, value)
	if err := r.AttachMany(children...); err != nil {
		return nil, fmt.Errorf("factory err: %w", err)
	}

	return r, nil
}

func (factory *Factory[T]) Tree(root *Node[T]) (*MTree[T], error) {
	return New[T](root, factory.maxBreadth, factory.maxDepth)
}
