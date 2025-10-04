package list

import (
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

// LinkedList implements a doubly-linked list data structure.
//
// This implementation provides efficient insertion and removal operations
// at both ends of the list, making it suitable for implementing queues,
// stacks, and deques. The list maintains references to both head and tail
// nodes for O(1) operations at either end.
//
// Key features:
//   - O(1) insertion and removal at both ends
//   - Automatic size tracking
//   - Memory-efficient with minimal overhead
//   - Safe handling of empty list conditions
//
// The LinkedList uses Node instances as nodes, where each Node contains
// an ID and bidirectional references to adjacent nodes.
//
// Thread Safety:
// LinkedList is not thread-safe. Concurrent access requires external
// synchronization mechanisms.
//
// Memory Management:
// The list creates copies of nodes during Pop() and Shift() operations
// to prevent memory leaks and ensure clean separation from the list structure.
type LinkedList struct {
	// size tracks the current number of nodes in the list.
	// It's automatically maintained by all modification operations.
	size int

	// head points to the first node in the list, or nil if the list is empty.
	head *node.Node

	// tail points to the last node in the list, or nil if the list is empty.
	tail *node.Node
}

// New creates a new empty LinkedList.
//
// The returned list is ready for use and has zero size with nil head and tail pointers.
//
// Returns:
//   - A new empty LinkedList instance
//
// Example:
//
//	list := New()
//	list.Push(node.New(1, nil, nil))
//	fmt.Printf("List size: %d", list.size) // Output: List size: 1
func New() *LinkedList {
	return &LinkedList{
		head: nil,
		tail: nil,
	}
}

// Push adds a node to the end (tail) of the list.
//
// This operation is O(1) and increases the list size by 1. The node's
// next and prev pointers will be updated to maintain list integrity.
// If the node already has next/prev references, they will be overwritten.
//
// Parameters:
//   - n: The Node to add to the end of the list. Must not be nil.
//
// Behavior:
//   - If the list is empty, the node becomes both head and tail
//   - If the list has nodes, the node is linked after the current tail
//   - The list size is automatically incremented
//
// Example:
//
//	list := New()
//	node := node.New(1, nil, nil)
//	list.Push(node)
//	// List now contains one node at both head and tail positions
func (list *LinkedList) Push(n *node.Node) {
	defer func() {
		list.size++
	}()

	if list.tail == nil {
		list.head = n
		list.tail = n
	} else {
		list.tail.WithNext(n)
		n.WithPrev(list.tail)
		list.tail = n
	}
}

// Pop removes and returns the last node (tail) from the list.
//
// This operation is O(1) and decreases the list size by 1. The returned
// node is a copy of the original node with nil next/prev pointers to
// prevent memory leaks and maintain clean separation from the list.
//
// Returns:
//   - A copy of the removed Node, or nil if the list is empty
//
// Behavior:
//   - If the list is empty, returns nil without changing the size
//   - If the list has one node, both head and tail become nil
//   - If the list has multiple nodes, the tail is moved to the previous node
//   - The list size is automatically decremented (only if size > 0)
//
// Memory Safety:
// The returned node is a copy with cleared next/prev pointers, ensuring
// no references back to the list structure remain.
//
// Example:
//
//	list := New()
//	list.Push(node.New(1, nil, nil))
//	node := list.Pop()
//	if node != nil {
//		fmt.Printf("Popped node ID: %d", node.ID())
//	}
func (list *LinkedList) Pop() *node.Node {
	if list.tail == nil {
		return nil
	}

	pTail := list.tail
	if list.tail.Prev() != nil {
		list.tail = list.tail.Prev()
		list.tail.WithNext(nil)
	} else {
		list.head = nil
		list.tail = nil
	}

	return list.cleanAndCopyNode(pTail)
}

// Unshift adds a node to the beginning (head) of the list.
//
// This operation is O(1) and increases the list size by 1. The node's
// next and prev pointers will be updated to maintain list integrity.
// If the node already has next/prev references, they will be overwritten.
//
// Parameters:
//   - n: The Node to add to the beginning of the list. Must not be nil.
//
// Behavior:
//   - If the list is empty, the node becomes both head and tail
//   - If the list has nodes, the node is linked before the current head
//   - The list size is automatically incremented
//
// Example:
//
//	list := New()
//	node := node.New(1, nil, nil)
//	list.Unshift(node)
//	// List now contains one node at both head and tail positions
func (list *LinkedList) Unshift(n *node.Node) {
	defer func() {
		list.size++
	}()

	if list.head == nil {
		list.head = n
		list.tail = n
	} else {
		n.WithNext(list.head)
		list.head.WithPrev(n)
		list.head = n
	}
}

// Shift removes and returns the first node (head) from the list.
//
// This operation is O(1) and decreases the list size by 1. The returned
// node is a copy of the original node with nil next/prev pointers to
// prevent memory leaks and maintain clean separation from the list.
//
// Returns:
//   - A copy of the removed Node, or nil if the list is empty
//
// Behavior:
//   - If the list is empty, returns nil without changing the size
//   - If the list has one node, both head and tail become nil
//   - If the list has multiple nodes, the head is moved to the next node
//   - The list size is automatically decremented (only if size > 0)
//
// Memory Safety:
// The returned node is a copy with cleared next/prev pointers, ensuring
// no references back to the list structure remain.
//
// Example:
//
//	list := New()
//	list.Unshift(node.New(1, nil, nil))
//	node := list.Shift()
//	if node != nil {
//		fmt.Printf("Shifted node ID: %d", node.ID())
//	}
func (list *LinkedList) Shift() *node.Node {
	if list.head == nil {
		return nil
	}

	pHead := list.head

	if list.head.Next() != nil {
		list.head = list.head.Next()
		list.head.WithPrev(nil)
	} else {
		list.head = nil
		list.tail = nil
	}

	return list.cleanAndCopyNode(pHead)
}

// Size returns the current number of nodes in the list.
//
// Returns:
//   - The current number of nodes in the list
func (list *LinkedList) Size() int {
	return list.size
}

// Head returns the first node in the list.
//
// Returns:
//   - The first node in the list, or an empty Node if the list is empty
func (list *LinkedList) Head() (node.Node, bool) {
	if list.head == nil {
		return node.Node{}, false
	}

	headCopy := *list.head
	return headCopy, true
}

// Tail returns the last node in the list.
//
// Returns:
//   - The last node in the list, or an empty Node if the list is empty
func (list *LinkedList) Tail() (node.Node, bool) {
	if list.tail == nil {
		return node.Node{}, false
	}

	tailCopy := *list.tail
	return tailCopy, true
}

// cleanAndCopyNode cleans a node's references and returns a copy.
// This helper method decrements the list size and ensures the node
// is properly disconnected from the list structure.
func (list *LinkedList) cleanAndCopyNode(n *node.Node) *node.Node {
	defer func() {
		if list.size > 0 {
			list.size--
		}
	}()

	n.WithNext(nil)
	n.WithPrev(nil)
	nCopy := *n

	return &nCopy
}
