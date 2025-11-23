package mtree

import (
	"fmt"
	"iter"
	"slices"
)

type MTree[T comparable] struct {
	maxBreadth int
	maxDepth   int
	capacity   int
	root       *Node[T]
	nodes      map[uint64]*Node[T]
	levels     map[int][]uint64
}

func New[T comparable](n *Node[T], maxBreadth, maxDepth int) (*MTree[T], error) {
	if n == nil {
		return nil, ErrInvalidRoot
	}

	if !n.IsRoot() {
		return nil, ErrInvalidRoot
	}

	capacity := maxBreadth * maxDepth
	nodes := make(map[uint64]*Node[T], capacity)
	nodes[n.id] = n
	levels := make(map[int][]uint64, maxDepth)
	levels[0] = []uint64{n.id}

	return &MTree[T]{
		maxBreadth: maxBreadth,
		maxDepth:   maxDepth,
		capacity:   capacity,
		root:       n,
		nodes:      nodes,
		levels:     levels,
	}, nil
}

func (t *MTree[T]) MaxBreadth() int {
	return t.maxBreadth
}

func (t *MTree[T]) MaxDepth() int {
	return t.maxDepth
}

func (t *MTree[T]) Capacity() int {
	return t.capacity
}

func (t *MTree[T]) Size() int {
	return len(t.nodes)
}

func (t *MTree[T]) ContainsNode(n *Node[T]) bool {
	if n == nil {
		return false
	}

	return t.ContainsID(n.ID())
}

func (t *MTree[T]) ContainsID(id uint64) bool {
	_, exists := t.nodes[id]

	return exists
}

func (t *MTree[T]) CurrDepth() int {
	var maxDepth int
	for level := range t.levels {
		if level > maxDepth {
			maxDepth = level
		}
	}

	return maxDepth
}

func (t *MTree[T]) DepthCapacity() int {
	return t.MaxDepth() - t.CurrDepth()
}

func (t *MTree[T]) Root() *Node[T] {
	return t.root
}

func (t *MTree[T]) Add(n *Node[T], upsert bool) error {
	if n == nil {
		return ErrNil
	}

	_, nodeExists := t.nodes[n.id]
	switch {
	case nodeExists && !upsert:
		return fmt.Errorf("node %d already exists: %w", n.id, ErrNodeConflict)
	default:
		t.nodes[n.id] = n
	}

	_, exists := t.levels[n.level]
	switch {
	case !exists:
		t.levels[n.level] = []uint64{n.id}
		return nil
	default:
		t.levels[n.level] = append(t.levels[n.level], n.id)
		return nil
	}
}

func (t *MTree[T]) Remove(n *Node[T]) error {
	if n == nil {
		return ErrNil
	}

	_, exists := t.nodes[n.id]
	if !exists {
		return ErrNodeNotFound
	}

	delete(t.nodes, n.id)
	if l, ok := t.levels[n.level]; ok {
		t.levels[n.level] = slices.DeleteFunc(l, func(siblingId uint64) bool {
			return siblingId == n.id
		})
	}

	return nil
}

func (t *MTree[T]) Level(level int) iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		ids, exists := t.levels[level]
		if !exists {
			return
		}

		for _, id := range ids {
			n, ok := t.nodes[id]
			if !ok {
				continue
			}
			if !yield(n) {
				return
			}
		}
	}
}
