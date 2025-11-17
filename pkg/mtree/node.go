package mtree

import (
	"errors"
	"fmt"
	"iter"

	"golang.org/x/exp/slices"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

const (
	root = iota
	attached
	detached
)

type (
	NodeOption[T comparable] func(n *Node[T])

	Node[T comparable] struct {
		id         uint64
		level      int
		maxBreadth int
		state      int
		val        T
		parent     *Node[T]
		children   map[uint64]*Node[T]
	}
)

func NewNode[T comparable](id uint64, maxBreadth int, opts ...NodeOption[T]) *Node[T] {
	n := &Node[T]{
		id:         id,
		level:      -1,
		state:      detached,
		parent:     nil,
		maxBreadth: maxBreadth,
		children:   make(map[uint64]*Node[T], maxBreadth),
	}
	for _, opt := range opts {
		opt(n)
	}

	return n
}

func NewRoot[T comparable](id uint64, maxBreadth int, val T) *Node[T] {
	rootNode := NewNode[T](id, maxBreadth, ValueOpt(val))
	_ = rootNode.AsRoot()
	return rootNode
}

func ValueOpt[T comparable](val T) NodeOption[T] {
	return func(n *Node[T]) {
		n.val = val
	}
}

func LevelOpt[T comparable](level int) NodeOption[T] {
	return func(n *Node[T]) {
		n.level = level
	}
}

func ParentOpt[T comparable](parent *Node[T]) NodeOption[T] {
	return func(n *Node[T]) {
		n.parent = parent
	}
}

func ChildOpt[T comparable](child *Node[T]) NodeOption[T] {
	return func(n *Node[T]) {
		if n == nil {
			return
		}

		child.parent = n
		n.children[serial.NSum(n.id, child.id)] = child
		child.level = n.level + 1
		child.state = attached
	}
}

func (n *Node[T]) verifyMaxBreadth(count int) error {
	if n.Capacity() < count {
		return ErrMaxBreadth
	}

	return nil
}

func (n *Node[T]) ID() uint64 {
	return n.id
}

func (n *Node[T]) Level() int {
	return n.level
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

func (n *Node[T]) Children() iter.Seq2[uint64, *Node[T]] {
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

func (n *Node[T]) DetachChildFunc(successorFn func(child *Node[T]) bool) int {
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

func (n *Node[T]) SelectChildrenFunc(successorFn func(child *Node[T]) bool) ([]*Node[T], error) {
	nodes := make([]*Node[T], 0, n.maxBreadth)
	for _, child := range n.children {
		if ok := successorFn(child); !ok {
			continue
		}

		copyChild := *child
		nodes = append(nodes, &copyChild)
	}

	switch {
	case len(nodes) > 0:
		return nodes, nil
	default:
		return nil, ErrNoMatch
	}
}

func (n *Node[T]) SelectChildByID(id uint64) (*Node[T], error) {
	relID := serial.NSum(n.id, id)
	child, exists := n.children[relID]
	if !exists {
		return nil, ErrNodeNotFound
	}

	copyChild := *child

	return &copyChild, nil
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

	targetParent := target.parent

	if target.IsRoot() {
		n.Detach()
		n.AsRoot()
	}

	if targetParent != nil {
		n.Detach()
		if err := targetParent.attach(n); err != nil {
			return err
		}
	}

	if n.IsRoot() {
		target.Detach()
		target.AsRoot()
	}

	sourceParent := n.parent
	if sourceParent != nil {
		target.Detach()
		if err := sourceParent.attach(target); err != nil {
			return err
		}
	}

	target.children, n.children = n.children, target.children

	return nil
}

func (n *Node[T]) AsRoot() bool {
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

func (n *Node[T]) IsAttached() bool {
	return n.state == attached
}

func (n *Node[T]) IsDetached() bool {
	return n.state == detached
}

func (n *Node[T]) Capacity() int {
	return n.maxBreadth - len(n.children)
}
