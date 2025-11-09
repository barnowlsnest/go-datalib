package tree

import (
	"cmp"
)

// BST represents a Binary Search Tree
type BST[T cmp.Ordered] struct {
	root *BinaryNode[T]
	size int
}

// NewBST creates a new empty Binary Search Tree
func NewBST[T cmp.Ordered]() *BST[T] {
	return &BST[T]{
		root: nil,
		size: 0,
	}
}

// Size returns the number of nodes in the tree
func (bst *BST[T]) Size() int {
	return bst.size
}

// IsEmpty returns true if the tree has no nodes
func (bst *BST[T]) IsEmpty() bool {
	return bst.size == 0
}

// Root returns the root node of the tree
func (bst *BST[T]) Root() *BinaryNode[T] {
	return bst.root
}

// Insert adds a new value to the binary search tree
func (bst *BST[T]) Insert(id uint64, value T) {
	newNode := NewBinaryNode(0, *NewNodeValue(id, value), nil, nil)

	if bst.root == nil {
		newNode.AsRoot()
		bst.root = newNode
		bst.size++
		return
	}

	bst.insertNode(bst.root, newNode, 1)
	bst.size++
}

func (bst *BST[T]) insertNode(current, newNode *BinaryNode[T], level int) {
	newNode.WithLevel(level)

	if newNode.val < current.val {
		if current.left == nil {
			newNode.AsLeft()
			current.WithLeft(newNode)
		} else {
			bst.insertNode(current.left, newNode, level+1)
		}
	} else {
		if current.right == nil {
			newNode.AsRight()
			current.WithRight(newNode)
		} else {
			bst.insertNode(current.right, newNode, level+1)
		}
	}
}

// Search finds a node with the given value
func (bst *BST[T]) Search(value T) *BinaryNode[T] {
	return bst.searchNode(bst.root, value)
}

func (bst *BST[T]) searchNode(node *BinaryNode[T], value T) *BinaryNode[T] {
	if node == nil {
		return nil
	}

	if value == node.val {
		return node
	}

	if value < node.val {
		return bst.searchNode(node.left, value)
	}

	return bst.searchNode(node.right, value)
}

// Contains checks if a value exists in the tree
func (bst *BST[T]) Contains(value T) bool {
	return bst.Search(value) != nil
}

// Min returns the minimum value in the tree
func (bst *BST[T]) Min() (T, bool) {
	if bst.root == nil {
		var zero T
		return zero, false
	}

	node := bst.minNode(bst.root)
	return node.val, true
}

func (bst *BST[T]) minNode(node *BinaryNode[T]) *BinaryNode[T] {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

// Max returns the maximum value in the tree
func (bst *BST[T]) Max() (T, bool) {
	if bst.root == nil {
		var zero T
		return zero, false
	}

	node := bst.maxNode(bst.root)
	return node.val, true
}

func (bst *BST[T]) maxNode(node *BinaryNode[T]) *BinaryNode[T] {
	current := node
	for current.right != nil {
		current = current.right
	}
	return current
}

// Delete removes a node with the given value from the tree
func (bst *BST[T]) Delete(value T) bool {
	if bst.root == nil {
		return false
	}

	var deleted bool
	bst.root, deleted = bst.deleteNode(bst.root, value)
	if deleted {
		bst.size--
	}
	return deleted
}

func (bst *BST[T]) deleteNode(node *BinaryNode[T], value T) (*BinaryNode[T], bool) {
	if node == nil {
		return nil, false
	}

	var deleted bool

	if value < node.val {
		node.left, deleted = bst.deleteNode(node.left, value)
		return node, deleted
	}

	if value > node.val {
		node.right, deleted = bst.deleteNode(node.right, value)
		return node, deleted
	}

	// Node to delete found
	deleted = true

	// Case 1: Node with no children (leaf node)
	if node.left == nil && node.right == nil {
		return nil, deleted
	}

	// Case 2: Node with one child
	if node.left == nil {
		node.right.hierarchy = node.hierarchy
		return node.right, deleted
	}

	if node.right == nil {
		node.left.hierarchy = node.hierarchy
		return node.left, deleted
	}

	// Case 3: Node with two children
	// Find the inorder successor (minimum in right subtree)
	successor := bst.minNode(node.right)

	// Copy successor's value to current node
	node.val = successor.val
	node.WithValue(successor.val)

	// Delete the successor
	node.right, _ = bst.deleteNode(node.right, successor.val)

	return node, deleted
}

// InOrder performs an in-order traversal (left, root, right)
func (bst *BST[T]) InOrder(visit func(*BinaryNode[T])) {
	bst.inOrder(bst.root, visit)
}

func (bst *BST[T]) inOrder(node *BinaryNode[T], visit func(*BinaryNode[T])) {
	if node == nil {
		return
	}

	bst.inOrder(node.left, visit)
	visit(node)
	bst.inOrder(node.right, visit)
}

// PreOrder performs a pre-order traversal (root, left, right)
func (bst *BST[T]) PreOrder(visit func(*BinaryNode[T])) {
	bst.preOrder(bst.root, visit)
}

func (bst *BST[T]) preOrder(node *BinaryNode[T], visit func(*BinaryNode[T])) {
	if node == nil {
		return
	}

	visit(node)
	bst.preOrder(node.left, visit)
	bst.preOrder(node.right, visit)
}

// PostOrder performs a post-order traversal (left, right, root)
func (bst *BST[T]) PostOrder(visit func(*BinaryNode[T])) {
	bst.postOrder(bst.root, visit)
}

func (bst *BST[T]) postOrder(node *BinaryNode[T], visit func(*BinaryNode[T])) {
	if node == nil {
		return
	}

	bst.postOrder(node.left, visit)
	bst.postOrder(node.right, visit)
	visit(node)
}

// Height returns the height of the tree (longest path from root to leaf)
func (bst *BST[T]) Height() int {
	return bst.height(bst.root)
}

func (bst *BST[T]) height(node *BinaryNode[T]) int {
	if node == nil {
		return -1
	}

	leftHeight := bst.height(node.left)
	rightHeight := bst.height(node.right)

	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// Clear removes all nodes from the tree
func (bst *BST[T]) Clear() {
	bst.root = nil
	bst.size = 0
}
