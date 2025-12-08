package mtree

import (
	"fmt"
	"iter"
	"slices"
)

// Container provides a bounded multi-way tree structure with level-based organization.
// It manages a tree with configurable maximum breadth (children per node) and depth (tree height),
// tracking nodes by ID and organizing them by level for efficient level-based traversal.
//
// The container enforces structural constraints:
//   - Maximum breadth: limits children per node
//   - Maximum depth: limits tree height
//   - Capacity: total maximum nodes (maxBreadth * maxDepth)
//   - Single root node at level 0
//
// Example:
//
//	root, _ := mtree.NewNode[string](1, 5, mtree.ValueOpt("root"))
//	container, _ := mtree.NewContainer(root, 5, 10)  // 5 children max, 10 levels deep
//	child, _ := mtree.NewNode[string](2, 5, mtree.ValueOpt("child"))
//	container.Insert(child, false)
type Container[T comparable] struct {
	maxBreadth int
	maxDepth   int
	capacity   int
	root       *Node[T]
	nodes      map[uint64]*Node[T]
	levels     map[int][]uint64
}

// NewContainer creates a new tree container with the specified root node and constraints.
// The root node will be marked as the tree root (level 0) if not already set.
//
// Parameters:
//   - root: The root node of the tree (must not be nil)
//   - maxBreadth: Maximum number of children per node
//   - maxDepth: Maximum tree depth (number of levels)
//
// Returns an error if:
//   - root is nil (ErrNil)
//   - root cannot be set as root node (ErrNodeConflict)
//
// Example:
//
//	root, _ := mtree.NewNode[int](1, 5, mtree.ValueOpt(100))
//	container, err := mtree.NewContainer(root, 5, 10)
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

// MaxBreadth returns the maximum number of children allowed per node.
func (c *Container[T]) MaxBreadth() int {
	return c.maxBreadth
}

// MaxDepth returns the maximum depth (number of levels) allowed in the tree.
func (c *Container[T]) MaxDepth() int {
	return c.maxDepth
}

// Capacity returns the maximum number of nodes the container can hold (maxBreadth * maxDepth).
func (c *Container[T]) Capacity() int {
	return c.capacity
}

// Size returns the current number of nodes stored in the container.
func (c *Container[T]) Size() int {
	return len(c.nodes)
}

// Contains checks whether the given node is present in the container.
// Returns false if n is nil.
func (c *Container[T]) Contains(n *Node[T]) bool {
	if n == nil {
		return false
	}

	return c.contains(n.ID())
}

// Depth returns the current maximum depth of the tree (highest level with nodes).
// For an empty tree, returns 0.
func (c *Container[T]) Depth() int {
	var d int
	for level := range c.levels {
		if level > d {
			d = level
		}
	}

	return d
}

// Insert adds a node to the container and associates it with the container.
// If upsert is true, replaces an existing node with the same ID; otherwise returns an error.
//
// Parameters:
//   - n: The node to insert (must not be nil)
//   - upsert: If true, allows updating existing nodes; if false, returns error on duplicate
//
// Returns an error if:
//   - n is nil (ErrNil)
//   - Node already exists and upsert is false (ErrNodeConflict)
//
// Example:
//
//	node, _ := mtree.NewNode[string](2, 5, mtree.ValueOpt("child"))
//	err := container.Insert(node, false)  // Insert new node
//	err = container.Insert(node, true)    // Update existing node
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

// Delete removes a node from the container and disassociates it.
// The node is removed from both the node map and its level tracking.
//
// Parameters:
//   - n: The node to delete (must not be nil)
//
// Returns an error if:
//   - n is nil (ErrNil)
//   - Node is not found in the container (ErrNodeNotFound)
//
// Example:
//
//	err := container.Delete(node)
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

// NodesIter returns an iterator over all nodes at the specified level.
// The iterator yields nodes in the order they were inserted at that level.
// If the level doesn't exist or has no nodes, the iterator returns immediately.
//
// Parameters:
//   - level: The tree level to iterate (0 = root level)
//
// Example:
//
//	for node := range container.NodesIter(1) {
//	    fmt.Println(node.Val())
//	}
func (c *Container[T]) NodesIter(level int) iter.Seq[*Node[T]] {
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
