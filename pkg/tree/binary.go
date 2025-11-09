package tree

import (
	"cmp"
)

const (
	rootNode = iota
	leftNode
	rightNode
)

type BinaryNode[T cmp.Ordered] struct {
	hierarchy int
	level     int
	*NodeValue[T]
	left  *BinaryNode[T]
	right *BinaryNode[T]
}

func NewBinaryNode[T cmp.Ordered](level int, nodeValue NodeValue[T], left, right *BinaryNode[T]) *BinaryNode[T] {
	return &BinaryNode[T]{
		NodeValue: &nodeValue,
		left:      left,
		right:     right,
		level:     level,
	}
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
