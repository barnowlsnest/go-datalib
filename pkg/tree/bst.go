package tree

import (
	"cmp"

	"github.com/barnowlsnest/go-datalib/pkg/node"
	"github.com/barnowlsnest/go-datalib/pkg/queue"
	"github.com/barnowlsnest/go-datalib/pkg/stack"
)

// BST (Binary Search Tree) is a production-ready, iterative implementation
// that maintains the BST property: for any node, all values in the left subtree
// are less than the node's value, and all values in the right subtree are greater.
//
// Key features:
//   - O(log n) average-case operations (search, insert, delete) for balanced trees
//   - O(n) worst-case for unbalanced trees (degenerates to linked list)
//   - All operations are iterative (no recursion) for better stack safety
//   - Uses stack for depth-first traversals (InOrder, PreOrder, PostOrder)
//   - Uses queue for breadth-first traversal (LevelOrder)
//   - Automatic size tracking
//
// Thread Safety:
// BST is not thread-safe. Concurrent access requires external synchronization.
//
// Memory Management:
// The BST maintains references to BinaryNode structures. When deleting nodes,
// the tree restructures to maintain BST properties.
type BST[T cmp.Ordered] struct {
	root *BinaryNode[T]
	size int
}

// New creates a new empty Binary Search Tree.
//
// Returns:
//   - A new empty BST instance ready for use
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
func New[T cmp.Ordered]() *BST[T] {
	return &BST[T]{
		root: nil,
		size: 0,
	}
}

// Insert adds a new value to the binary search tree while maintaining BST properties.
// This is an iterative implementation with O(log n) average time complexity.
//
// Parameters:
//   - value: The NodeValue to insert into the tree
//
// Returns:
//   - true if the value was inserted successfully
//   - false if the value already exists (duplicates are not allowed)
//
// Example:
//
//	bst := New[int]()
//	inserted := bst.Insert(NewNodeValue(1, 50)) // returns true
//	inserted = bst.Insert(NewNodeValue(2, 50))  // returns false (duplicate value)
func (bst *BST[T]) Insert(value *NodeValue[T]) bool {
	if value == nil {
		return false
	}

	props, err := value.Props()
	if err != nil {
		return false
	}

	newNode := NewBinaryNode(0, value, nil, nil)

	// Empty tree case
	if bst.root == nil {
		newNode.AsRoot()
		bst.root = newNode
		bst.size++
		return true
	}

	// Iterative search for insertion point
	current := bst.root
	level := 0

	for {
		currentProps, err := current.Props()
		if err != nil {
			return false
		}

		// Duplicate check
		if props.Value == currentProps.Value {
			return false
		}

		level++

		if props.Value < currentProps.Value {
			// Go left
			if !current.HasLeft() {
				newNode.AsLeft()
				newNode.WithLevel(level)
				current.WithLeft(newNode)
				bst.size++
				return true
			}
			current = current.Left()
		} else {
			// Go right
			if !current.HasRight() {
				newNode.AsRight()
				newNode.WithLevel(level)
				current.WithRight(newNode)
				bst.size++
				return true
			}
			current = current.Right()
		}
	}
}

// Search finds a value in the binary search tree using iterative binary search.
// This operation has O(log n) average time complexity.
//
// Parameters:
//   - value: The value to search for
//
// Returns:
//   - The BinaryNode containing the value if found, nil otherwise
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	node := bst.Search(50) // returns the node with value 50
//	node = bst.Search(99)  // returns nil (not found)
func (bst *BST[T]) Search(value T) *BinaryNode[T] {
	current := bst.root

	for current != nil {
		props, err := current.Props()
		if err != nil {
			return nil
		}

		if value == props.Value {
			return current
		}

		if value < props.Value {
			current = current.Left()
		} else {
			current = current.Right()
		}
	}

	return nil
}

// Delete removes a value from the binary search tree while maintaining BST properties.
// This is an iterative implementation that handles three cases:
//  1. Node with no children (leaf): simply remove
//  2. Node with one child: replace node with its child
//  3. Node with two children: replace with inorder successor (leftmost node in right subtree)
//
// Parameters:
//   - value: The value to delete from the tree
//
// Returns:
//   - true if the value was found and deleted
//   - false if the value was not found in the tree
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	deleted := bst.Delete(50) // returns true
//	deleted = bst.Delete(99)  // returns false (not found)
func (bst *BST[T]) Delete(value T) bool {
	if bst.root == nil {
		return false
	}

	// Find the node and its parent
	parent, current, isLeftChild := bst.findNodeWithParent(value)

	// Value not found
	if current == nil {
		return false
	}

	// Determine node type and handle deletion
	switch {
	case !current.HasLeft() && !current.HasRight():
		// Case 1: Leaf node (no children)
		bst.deleteLeafNode(parent, current, isLeftChild)
	case !current.HasLeft() || !current.HasRight():
		// Case 2: Node with one child
		bst.deleteNodeWithOneChild(parent, current, isLeftChild)
	default:
		// Case 3: Node with two children
		bst.deleteNodeWithTwoChildren(current)
	}

	bst.size--
	return true
}

// findNodeWithParent locates a node by value and returns its parent and position.
func (bst *BST[T]) findNodeWithParent(value T) (parent, current *BinaryNode[T], isLeftChild bool) {
	parent = nil
	current = bst.root
	isLeftChild = false

	for current != nil {
		props, err := current.Props()
		if err != nil {
			return nil, nil, false
		}

		if value == props.Value {
			return parent, current, isLeftChild
		}

		parent = current
		if value < props.Value {
			current = current.Left()
			isLeftChild = true
		} else {
			current = current.Right()
			isLeftChild = false
		}
	}

	return nil, nil, false
}

// deleteLeafNode removes a leaf node (node with no children).
func (bst *BST[T]) deleteLeafNode(parent, current *BinaryNode[T], isLeftChild bool) {
	if current == bst.root {
		bst.root = nil
		return
	}

	if isLeftChild {
		parent.WithLeft(nil)
	} else {
		parent.WithRight(nil)
	}
}

// deleteNodeWithOneChild removes a node with exactly one child.
func (bst *BST[T]) deleteNodeWithOneChild(parent, current *BinaryNode[T], isLeftChild bool) {
	// Determine which child exists
	var child *BinaryNode[T]
	if current.HasLeft() {
		child = current.Left()
	} else {
		child = current.Right()
	}

	// Replace current with its child
	if current == bst.root {
		bst.root = child
		bst.root.AsRoot()
		return
	}

	if isLeftChild {
		parent.WithLeft(child)
		child.AsLeft()
	} else {
		parent.WithRight(child)
		child.AsRight()
	}
}

// deleteNodeWithTwoChildren removes a node with two children using inorder successor.
func (bst *BST[T]) deleteNodeWithTwoChildren(current *BinaryNode[T]) {
	// Find inorder successor (leftmost node in right subtree)
	successor := bst.findMin(current.Right())

	// Get successor's value
	successorProps, err := successor.Props()
	if err != nil {
		return
	}

	// Delete successor (it has at most one child - right child)
	bst.Delete(successorProps.Value)

	// Replace current node's value with successor's value
	current.WithValue(successorProps.Value)

	// Compensate for the recursive delete that decremented size
	bst.size++
}

// findMin finds the node with minimum value in a subtree (iterative).
// Helper function used during deletion.
func (bst *BST[T]) findMin(n *BinaryNode[T]) *BinaryNode[T] {
	current := n
	for current.HasLeft() {
		current = current.Left()
	}
	return current
}

// Min returns the node with the minimum value in the tree.
// Time complexity: O(h) where h is the height of the tree.
//
// Returns:
//   - The BinaryNode with minimum value, or nil if tree is empty
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	min := bst.Min() // returns node with value 30
func (bst *BST[T]) Min() *BinaryNode[T] {
	if bst.root == nil {
		return nil
	}
	return bst.findMin(bst.root)
}

// Max returns the node with the maximum value in the tree.
// Time complexity: O(h) where h is the height of the tree.
//
// Returns:
//   - The BinaryNode with maximum value, or nil if tree is empty
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	max := bst.Max() // returns node with value 70
func (bst *BST[T]) Max() *BinaryNode[T] {
	if bst.root == nil {
		return nil
	}

	current := bst.root
	for current.HasRight() {
		current = current.Right()
	}
	return current
}

// InOrder performs an iterative in-order traversal (Left-Root-Right) using a stack.
// For a BST, this produces values in sorted ascending order.
// Time complexity: O(n), Space complexity: O(h) where h is tree height.
//
// Parameters:
//   - visit: Function to call for each node during traversal
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	bst.InOrder(func(node *BinaryNode[int]) {
//		props, _ := node.NodeValue.Props()
//		fmt.Println(props.Value) // Prints: 30, 50, 70
//	})
func (bst *BST[T]) InOrder(visit func(*BinaryNode[T])) {
	if bst.root == nil || visit == nil {
		return
	}

	s := stack.New()
	nodeMap := make(map[uint64]*BinaryNode[T])

	current := bst.root

	// Push all left nodes
	for current != nil {
		bst.addToStack(s, current, nodeMap)
		current = current.Left()
	}

	for !s.IsEmpty() {
		n := s.Pop()
		if n == nil {
			break
		}

		current = nodeMap[n.ID()]
		visit(current)

		// Process right subtree
		if current.HasRight() {
			current = current.Right()
			for current != nil {
				bst.addToStack(s, current, nodeMap)
				current = current.Left()
			}
		}
	}
}

// PreOrder performs an iterative pre-order traversal (Root-Left-Right) using a stack.
// Time complexity: O(n), Space complexity: O(h) where h is tree height.
//
// Parameters:
//   - visit: Function to call for each node during traversal
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	bst.PreOrder(func(node *BinaryNode[int]) {
//		props, _ := node.NodeValue.Props()
//		fmt.Println(props.Value) // Prints: 50, 30, 70
//	})
func (bst *BST[T]) PreOrder(visit func(*BinaryNode[T])) {
	if bst.root == nil || visit == nil {
		return
	}

	s := stack.New()
	nodeMap := make(map[uint64]*BinaryNode[T])

	bst.addToStack(s, bst.root, nodeMap)

	bst.traverseWithStack(s, nodeMap, visit, func(current *BinaryNode[T]) {
		// Push right first (so left is processed first)
		if current.HasRight() {
			bst.addToStack(s, current.Right(), nodeMap)
		}
		if current.HasLeft() {
			bst.addToStack(s, current.Left(), nodeMap)
		}
	})
}

// PostOrder performs an iterative post-order traversal (Left-Right-Root) using two stacks.
// Time complexity: O(n), Space complexity: O(h) where h is tree height.
//
// Parameters:
//   - visit: Function to call for each node during traversal
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	bst.PostOrder(func(node *BinaryNode[int]) {
//		props, _ := node.NodeValue.Props()
//		fmt.Println(props.Value) // Prints: 30, 70, 50
//	})
func (bst *BST[T]) PostOrder(visit func(*BinaryNode[T])) {
	if bst.root == nil || visit == nil {
		return
	}

	s1 := stack.New()
	s2 := stack.New()
	nodeMap := make(map[uint64]*BinaryNode[T])

	bst.addToStack(s1, bst.root, nodeMap)

	// Fill s2 with nodes in reverse post-order
	for !s1.IsEmpty() {
		n := s1.Pop()
		if n == nil {
			break
		}

		current := nodeMap[n.ID()]
		bst.addToStack(s2, current, nodeMap)

		// Push left first, then right
		if current.HasLeft() {
			bst.addToStack(s1, current.Left(), nodeMap)
		}
		if current.HasRight() {
			bst.addToStack(s1, current.Right(), nodeMap)
		}
	}

	// Process nodes in post-order
	for !s2.IsEmpty() {
		n := s2.Pop()
		if n == nil {
			break
		}
		visit(nodeMap[n.ID()])
	}
}

// LevelOrder performs an iterative level-order (breadth-first) traversal using a queue.
// Visits nodes level by level from left to right.
// Time complexity: O(n), Space complexity: O(w) where w is maximum width.
//
// Parameters:
//   - visit: Function to call for each node during traversal
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	bst.LevelOrder(func(node *BinaryNode[int]) {
//		props, _ := node.NodeValue.Props()
//		fmt.Println(props.Value) // Prints: 50, 30, 70
//	})
func (bst *BST[T]) LevelOrder(visit func(*BinaryNode[T])) {
	if bst.root == nil || visit == nil {
		return
	}

	q := queue.New()
	nodeMap := make(map[uint64]*BinaryNode[T])

	bst.addToQueue(q, bst.root, nodeMap)

	for !q.IsEmpty() {
		n := q.Dequeue()
		if n == nil {
			break
		}

		current := nodeMap[n.ID()]
		visit(current)

		if current.HasLeft() {
			bst.addToQueue(q, current.Left(), nodeMap)
		}
		if current.HasRight() {
			bst.addToQueue(q, current.Right(), nodeMap)
		}
	}
}

// Height returns the height of the tree (longest path from root to leaf).
// An empty tree has height -1, a tree with only root has height 0.
// This is an iterative level-order approach.
// Time complexity: O(n)
//
// Returns:
//   - The height of the tree
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	bst.Insert(NewNodeValue(3, 70))
//	height := bst.Height() // returns 1
func (bst *BST[T]) Height() int {
	if bst.root == nil {
		return -1
	}

	q := queue.New()
	nodeMap := make(map[uint64]*BinaryNode[T])

	bst.addToQueue(q, bst.root, nodeMap)
	height := -1

	for !q.IsEmpty() {
		levelSize := q.Size()
		height++

		for i := 0; i < levelSize; i++ {
			n := q.Dequeue()
			if n == nil {
				continue
			}

			current := nodeMap[n.ID()]

			if current.HasLeft() {
				bst.addToQueue(q, current.Left(), nodeMap)
			}
			if current.HasRight() {
				bst.addToQueue(q, current.Right(), nodeMap)
			}
		}
	}

	return height
}

// Size returns the number of nodes in the tree.
// Time complexity: O(1)
//
// Returns:
//   - The number of nodes in the tree
//
// Example:
//
//	bst := New[int]()
//	bst.Insert(NewNodeValue(1, 50))
//	bst.Insert(NewNodeValue(2, 30))
//	size := bst.Size() // returns 2
func (bst *BST[T]) Size() int {
	return bst.size
}

// IsEmpty returns true if the tree contains no nodes.
//
// Returns:
//   - true if tree is empty, false otherwise
func (bst *BST[T]) IsEmpty() bool {
	return bst.size == 0
}

// Root returns the root node of the tree.
//
// Returns:
//   - The root BinaryNode, or nil if tree is empty
func (bst *BST[T]) Root() *BinaryNode[T] {
	return bst.root
}

// traverseWithStack is a generic stack-based traversal using the strategy pattern.
// It encapsulates the common iteration logic while allowing different child addition strategies.
func (bst *BST[T]) traverseWithStack(
	s *stack.Stack,
	nodeMap map[uint64]*BinaryNode[T],
	visit func(*BinaryNode[T]),
	addChildren func(*BinaryNode[T]),
) {
	for !s.IsEmpty() {
		n := s.Pop()
		if n == nil {
			break
		}

		current := nodeMap[n.ID()]
		visit(current)
		addChildren(current)
	}
}

// addToStack is a helper function to add a BinaryNode to a stack.
// It maps the node ID to the actual BinaryNode for later retrieval.
func (bst *BST[T]) addToStack(s *stack.Stack, bn *BinaryNode[T], nodeMap map[uint64]*BinaryNode[T]) {
	if bn == nil {
		return
	}
	props, err := bn.Props()
	if err != nil {
		return
	}
	nodeMap[props.ID] = bn
	s.Push(node.ID(props.ID))
}

// addToQueue is a helper function to add a BinaryNode to a queue.
// It maps the node ID to the actual BinaryNode for later retrieval.
func (bst *BST[T]) addToQueue(q *queue.Queue, bn *BinaryNode[T], nodeMap map[uint64]*BinaryNode[T]) {
	if bn == nil {
		return
	}
	props, err := bn.Props()
	if err != nil {
		return
	}
	nodeMap[props.ID] = bn
	q.Enqueue(node.ID(props.ID))
}
