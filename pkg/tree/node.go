// Package mtree provides a generic multi-way tree (M-Tree) implementation with support
// for hierarchical data structures, node operations, and cycle detection.
//
// The package offers:
//   - Generic Node[T] type for building tree structures with any comparable type
//   - Configurable maximum breadth (max children per node)
//   - Parent-child relationship management with automatic level tracking
//   - Tree traversal and node selection operations
//   - Concurrent-safe node selection with context support
//   - Hierarchy model builder with cycle detection
//   - Node swapping and movement operations
package tree

import (
	"context"
	"errors"
	"fmt"
	"iter"

	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

const (
	root = iota
	attached
	detached
)

type (
	// NodeOption is a functional option for configuring a Node during creation.
	NodeOption[T comparable] func(n *Node[T]) error

	// Node represents a node in a multi-way tree with generic value type T.
	// Each node can have multiple children (up to maxBreadth), a single parent,
	// and maintains its level in the tree hierarchy.
	//
	// Node states:
	//   - root: The top-level node with no parent (level 0)
	//   - attached: A child node connected to a parent
	//   - detached: An orphaned node not connected to any parent (level -1)
	Node[T comparable] struct {
		id         uint64
		level      int
		maxBreadth int
		state      int
		val        T
		parent     *Node[T]
		children   map[uint64]*Node[T]
	}

	// NodeSuccessorFunc is a predicate function for filtering/selecting child nodes.
	NodeSuccessorFunc[T comparable] func(child *Node[T]) bool
)

// NewNode creates a new tree node with the given ID and maximum breadth.
// The maxBreadth parameter limits the number of children this node can have.
//
// Optional configuration can be applied using NodeOption functions:
//   - ValueOpt: Set the node's value
//   - LevelOpt: Set the node's level (typically for root nodes)
//   - ParentOpt: Attach this node as a child of a parent
//   - ChildOpt: Attach existing nodes as children
//
// Returns an error if:
//   - Any option returns an error
//   - MaxBreadth would be exceeded when attaching children
//
// Example:
//
//	root, err := NewNode[string](1, 5, ValueOpt("root"), LevelOpt(0))
//	child, err := NewNode[string](2, 3, ValueOpt("child"), ParentOpt(root))
func NewNode[T comparable](id uint64, maxBreadth int, opts ...NodeOption[T]) (*Node[T], error) {
	n := &Node[T]{
		id:         id,
		level:      -1,
		state:      detached,
		parent:     nil,
		maxBreadth: maxBreadth,
		children:   make(map[uint64]*Node[T], maxBreadth),
	}
	for _, opt := range opts {
		if err := opt(n); err != nil {
			return nil, err
		}
	}

	return n, nil
}

func ValueOpt[T comparable](val T) NodeOption[T] {
	return func(n *Node[T]) error {
		n.WithValue(val)

		return nil
	}
}

func LevelOpt[T comparable](level int) NodeOption[T] {
	return func(n *Node[T]) error {
		n.level = level

		return nil
	}
}

func ParentOpt[T comparable](parent *Node[T]) NodeOption[T] {
	return func(n *Node[T]) error {
		if parent == nil {
			return ErrNil
		}

		if err := parent.AttachChild(n); err != nil {
			return err
		}

		return nil
	}
}

func ChildOpt[T comparable](child *Node[T]) NodeOption[T] {
	return func(n *Node[T]) error {
		if child == nil {
			return ErrNil
		}

		if n.level < 0 {
			n.level = 0
		}

		if err := n.AttachChild(child); err != nil {
			return err
		}

		return nil
	}
}

func (n *Node[T]) verifyMaxBreadth(count int) error {
	if n.Capacity() < count {
		return ErrMaxBreadth
	}

	return nil
}

func (n *Node[T]) asRoot() bool {
	if n.IsRoot() {
		return true
	}

	switch {
	case n.parent != nil:
		return false
	case n.level >= 0:
		return false
	}

	n.state = root
	n.level = 0
	n.parent = nil

	return true
}

func (n *Node[T]) ID() uint64 {
	return n.id
}

func (n *Node[T]) Level() int {
	return n.level
}

// setLevel sets the node's level. Package-private method for use by Segment
// when updating levels after structural changes.
func (n *Node[T]) setLevel(level int) {
	n.level = level
}

func (n *Node[T]) Val() T {
	return n.val
}

func (n *Node[T]) WithValue(val T) {
	n.val = val
}

func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

func (n *Node[T]) HasParent() bool {
	return n.parent != nil
}

func (n *Node[T]) IsRoot() bool {
	return n.state == root
}

func (n *Node[T]) HasChild(child *Node[T]) bool {
	if child == nil {
		return false
	}

	if len(n.children) == 0 {
		return false
	}

	_, ok := n.children[serial.NSum(n.id, child.id)]

	return ok
}

func (n *Node[T]) MaxBreadth() int {
	return n.maxBreadth
}

func (n *Node[T]) IsChildOf(parentNode *Node[T]) bool {
	switch {
	case parentNode == nil:
		return false
	case n.parent != parentNode:
		return false
	}

	_, exists := parentNode.children[serial.NSum(parentNode.id, n.id)]

	return exists
}

func (n *Node[T]) HasChildren() bool {
	return len(n.children) > 0
}

func (n *Node[T]) Breadth() int {
	return len(n.children)
}

func (n *Node[T]) attach(child *Node[T]) error {
	switch {
	case n == nil:
		return fmt.Errorf("not valid parent: %w", ErrNil)
	case child == nil:
		return fmt.Errorf("not valid child: %w", ErrNil)
	}

	relID := serial.NSum(n.id, child.id)
	n.children[relID] = child
	child.parent = n
	child.level = n.level + 1
	child.state = attached

	return nil
}

func (n *Node[T]) AttachChild(childNode *Node[T]) error {
	if err := n.verifyMaxBreadth(1); err != nil {
		return err
	}

	return n.attach(childNode)
}

func (n *Node[T]) AttachMany(children ...*Node[T]) error {
	var err error
	clean := slices.DeleteFunc(children, func(n *Node[T]) bool { return n == nil })
	if err := n.verifyMaxBreadth(len(clean)); err != nil {
		return err
	}

	errCollector := make([]error, 0, len(clean))
	for _, child := range clean {
		if err = n.attach(child); err != nil {
			errCollector = append(errCollector, err)
		}
	}

	if len(errCollector) > 0 {
		collectedErrors := errors.Join(errCollector...)
		return fmt.Errorf("inconsistent attach: %w", collectedErrors)
	}

	return nil
}

func (n *Node[T]) ChildrenIter() iter.Seq2[uint64, *Node[T]] {
	return func(yield func(uint64, *Node[T]) bool) {
		for id, child := range n.children {
			if !yield(id, child) {
				return
			}
		}
	}
}

func (n *Node[T]) DetachChild(child *Node[T]) error {
	if child == nil {
		return fmt.Errorf("nil child node:%w", ErrNil)
	}

	id := serial.NSum(n.id, child.id)
	childNode, exists := n.children[id]
	if !exists {
		return ErrNodeNotFound
	}

	childNode.Detach()

	return nil
}

func (n *Node[T]) DetachChildFunc(successorFn NodeSuccessorFunc[T]) int {
	if successorFn == nil {
		return 0
	}

	var count int
	for _, child := range n.children {
		if ok := successorFn(child); !ok {
			continue
		}

		child.Detach()
		count++
	}

	return count
}

func (n *Node[T]) SelectChildrenFunc(successorFn NodeSuccessorFunc[T]) ([]*Node[T], error) {
	nodes := make([]*Node[T], 0, n.maxBreadth)
	for _, child := range n.children {
		if ok := successorFn(child); !ok {
			continue
		}

		nodes = append(nodes, child)
	}

	switch {
	case len(nodes) > 0:
		return nodes, nil
	default:
		return nil, ErrNoMatch
	}
}

func (n *Node[T]) SelectOneChildFunc(successorFn NodeSuccessorFunc[T]) (*Node[T], error) {
	for _, child := range n.children {
		if successorFn(child) {
			return child, nil
		}
	}

	return nil, ErrNoMatch
}

func (n *Node[T]) SelectOneChildByEachValue(ctx context.Context, values ...T) (map[T]*Node[T], error) {
	dedup := make(map[T]struct{}, len(values))
	for _, val := range values {
		dedup[val] = struct{}{}
	}

	if len(dedup) == 0 {
		return make(map[T]*Node[T]), nil
	}

	eg := errgroup.Group{}
	childCh := make(chan *Node[T], len(dedup))
	errCh := make(chan error, 1)

	for val := range dedup {
		val := val // Required: capture loop variable for goroutine closure
		eg.Go(func() error {
			child, err := n.SelectOneChildFunc(func(n *Node[T]) bool {
				return n.Val() == val
			})
			if err != nil {
				return err
			}

			childCh <- child

			return nil
		})
	}

	// Wait for all goroutines in a separate goroutine
	go func() {
		if err := eg.Wait(); err != nil {
			errCh <- err
			close(errCh)
		} else {
			close(childCh)
			close(errCh)
		}
	}()

	res := make(map[T]*Node[T])
	expectedCount := len(dedup)
	receivedCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err, ok := <-errCh:
			if ok && err != nil {
				return nil, err
			}
			// errCh closed without error - continue draining childCh
		case child, ok := <-childCh:
			if !ok {
				// childCh closed - all results received
				return res, nil
			}
			res[child.Val()] = child
			receivedCount++
			if receivedCount == expectedCount {
				// All expected children received
				return res, nil
			}
		}
	}
}

func (n *Node[T]) SelectChildByID(id uint64) (*Node[T], error) {
	relID := serial.NSum(n.id, id)
	child, exists := n.children[relID]
	if !exists {
		return nil, ErrNodeNotFound
	}

	return child, nil
}

func (n *Node[T]) Detach() {
	p := n.parent
	if p == nil {
		return
	}

	n.parent = nil
	delete(p.children, serial.NSum(p.id, n.id))
	n.state = detached
	n.level = -1
}

func (n *Node[T]) MoveChildren(newParent *Node[T]) error {
	if newParent == nil {
		return fmt.Errorf("nil parent node:%w", ErrNil)
	}

	if err := newParent.verifyMaxBreadth(n.Breadth()); err != nil {
		return err
	}

	errCollector := make([]error, 0, len(n.children))
	for _, child := range n.children {
		child.Detach()
		if err := newParent.attach(child); err != nil {
			errCollector = append(errCollector, err)
		}
	}

	if len(errCollector) > 0 {
		collectedErrors := errors.Join(errCollector...)
		return fmt.Errorf("inconsistent children move: %w", collectedErrors)
	}

	return nil
}

func (n *Node[T]) Move(newParent *Node[T]) error {
	if newParent == nil {
		return fmt.Errorf("nil parent node:%w", ErrNil)
	}

	if err := newParent.verifyMaxBreadth(1); err != nil {
		return err
	}

	n.Detach()
	return newParent.attach(n)
}

func (n *Node[T]) Swap(target *Node[T]) error {
	if target == nil {
		return fmt.Errorf("nil target node: %w", ErrNil)
	}

	parent := n.parent
	targetParent := target.parent

	n.Detach()
	target.Detach()

	if target.IsRoot() {
		n.asRoot()
	}
	if n.IsRoot() {
		target.asRoot()
	}

	if targetParent != nil {
		if err := targetParent.attach(n); err != nil {
			return err
		}
	}

	if parent != nil {
		target.Detach()
		if err := parent.attach(target); err != nil {
			return err
		}
	}

	target.children, n.children = n.children, target.children

	return nil
}

func (n *Node[T]) IsAttached() bool {
	return n.state == attached
}

func (n *Node[T]) IsDetached() bool {
	return n.state == detached
}

func (n *Node[T]) Capacity() int {
	return n.MaxBreadth() - n.Breadth()
}
