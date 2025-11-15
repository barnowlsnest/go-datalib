package tree

import (
	"cmp"

	"github.com/barnowlsnest/go-datalib/pkg/node"
)

const (
	rootNode = iota
	leftNode
	rightNode
)

type (
	BinaryNodeOption[T cmp.Ordered] func(bn *BinaryNode[T])

	BinaryNode[T cmp.Ordered] struct {
		val       T
		hierarchy int
		level     int
		*node.Node
		left  *BinaryNode[T]
		right *BinaryNode[T]
	}
)

func WithValue[T cmp.Ordered](val T) BinaryNodeOption[T] {
	return func(bn *BinaryNode[T]) {
		bn.val = val
	}
}

func WithLevel[T cmp.Ordered](level int) BinaryNodeOption[T] {
	return func(bn *BinaryNode[T]) {
		bn.level = level
	}
}

func WithLeft[T cmp.Ordered](left *BinaryNode[T]) BinaryNodeOption[T] {
	return func(bn *BinaryNode[T]) {
		bn.left = left
	}
}

func WithRight[T cmp.Ordered](right *BinaryNode[T]) BinaryNodeOption[T] {
	return func(bn *BinaryNode[T]) {
		bn.right = right
	}
}

func NewBinaryNode[T cmp.Ordered](n *node.Node, opts ...BinaryNodeOption[T]) *BinaryNode[T] {
	bn := &BinaryNode[T]{Node: n}

	for _, opt := range opts {
		opt(bn)
	}

	return bn
}

func (bn *BinaryNode[T]) WithValue(val T) {
	bn.val = val
}

func (bn *BinaryNode[T]) Value() T {
	return bn.val
}

func (bn *BinaryNode[T]) WithLeft(left *BinaryNode[T]) {
	bn.left = left
}

func (bn *BinaryNode[T]) Left() *BinaryNode[T] {
	return bn.left
}

func (bn *BinaryNode[T]) WithRight(right *BinaryNode[T]) {
	bn.right = right
}

func (bn *BinaryNode[T]) Right() *BinaryNode[T] {
	return bn.right
}

func (bn *BinaryNode[T]) WithLevel(level int) {
	bn.level = level
}

func (bn *BinaryNode[T]) Level() int {
	return bn.level
}

func (bn *BinaryNode[T]) HasLeft() bool {
	return bn.left != nil
}

func (bn *BinaryNode[T]) HasRight() bool {
	return bn.right != nil
}

func (bn *BinaryNode[T]) HashChildren() bool {
	return bn.left != nil || bn.right != nil
}

func (bn *BinaryNode[T]) AsRoot() {
	bn.hierarchy = rootNode
}

func (bn *BinaryNode[T]) AsLeft() {
	bn.hierarchy = leftNode
}

func (bn *BinaryNode[T]) AsRight() {
	bn.hierarchy = rightNode
}

func (bn *BinaryNode[T]) IsRoot() bool {
	return bn.hierarchy == rootNode
}

func (bn *BinaryNode[T]) IsLeft() bool {
	return bn.hierarchy == leftNode
}

func (bn *BinaryNode[T]) IsRight() bool {
	return bn.hierarchy == rightNode
}
