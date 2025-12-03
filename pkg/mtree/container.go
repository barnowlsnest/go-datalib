package mtree

import (
	"fmt"
	"iter"
	"slices"
)

type Container[T comparable] struct {
	maxBreadth int
	maxDepth   int
	capacity   int
	root       *Node[T]
	nodes      map[uint64]*Node[T]
	levels     map[int][]uint64
}

func NewContainer[T comparable](root *Node[T], maxBreadth, maxDepth int) (*Container[T], error) {
	if root == nil {
		return nil, ErrNil
	}
	if !root.asRoot() {
		return nil, ErrNodeConflict
	}

	capacity := maxBreadth * maxDepth
	t := &Container[T]{
		root:       root,
		maxBreadth: maxBreadth,
		maxDepth:   maxDepth,
		capacity:   capacity,
		nodes:      make(map[uint64]*Node[T], capacity),
		levels:     make(map[int][]uint64, maxDepth),
	}

	return t, nil
}

func (c *Container[T]) contains(id uint64) bool {
	_, exists := c.nodes[id]

	return exists
}

func (c *Container[T]) MaxBreadth() int {
	return c.maxBreadth
}

func (c *Container[T]) MaxDepth() int {
	return c.maxDepth
}

func (c *Container[T]) Capacity() int {
	return c.capacity
}

func (c *Container[T]) Size() int {
	return len(c.nodes)
}

func (c *Container[T]) Contains(n *Node[T]) bool {
	if n == nil {
		return false
	}

	return c.contains(n.ID())
}

func (c *Container[T]) CurrentDepth() int {
	var d int
	for level := range c.levels {
		if level > d {
			d = level
		}
	}

	return d
}

func (c *Container[T]) Insert(n *Node[T], upsert bool) error {
	if n == nil {
		return ErrNil
	}

	_, nodeExists := c.nodes[n.id]
	switch {
	case nodeExists && !upsert:
		return fmt.Errorf("node %d already exists: %w", n.id, ErrNodeConflict)
	default:
		c.nodes[n.id] = n
		n.associate(c)
	}

	_, levelExists := c.levels[n.level]
	switch {
	case !levelExists:
		c.levels[n.level] = []uint64{n.id}
	default:
		c.levels[n.level] = append(c.levels[n.level], n.id)
	}

	return nil
}

func (c *Container[T]) Delete(n *Node[T]) error {
	if n == nil {
		return ErrNil
	}

	_, exists := c.nodes[n.id]
	if !exists {
		return ErrNodeNotFound
	}

	n.associate(nil)
	delete(c.nodes, n.id)

	if l, ok := c.levels[n.level]; ok {
		c.levels[n.level] = slices.DeleteFunc(l, func(siblingId uint64) bool {
			return siblingId == n.id
		})
	}

	return nil
}

func (c *Container[T]) Nodes(level int) iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		ids, exists := c.levels[level]
		if !exists {
			return
		}

		for _, id := range ids {
			n, ok := c.nodes[id]
			if !ok {
				continue
			}
			if !yield(n) {
				return
			}
		}
	}
}
