package list

import (
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

// Stack implements a LIFO (Last In, First Out) data structure.
//
// This implementation uses a LinkedList internally to provide O(1) push and pop
// operations. The stack follows the classic stack semantics where elements are
// added and removed from the same end (top of the stack).
//
// Key features:
//   - O(1) push and pop operations
//   - O(1) peek operation to view top element without removal
//   - Automatic size tracking
//   - Safe handling of empty stack conditions
//
// Thread Safety:
// Stack is not thread-safe. Concurrent access requires external
// synchronization mechanisms.
//
// Memory Management:
// The stack creates copies of nodes during Pop() operations to prevent
// memory leaks and ensure clean separation from the internal structure.
type Stack struct {
	// list is the internal LinkedList used to store stack elements.
	// Elements are pushed and popped from the tail (top of stack).
	list *LinkedList
}

// NewStack creates a new empty Stack.
//
// The returned stack is ready for use and has zero size.
//
// Returns:
//   - A new empty Stack instance
//
// Example:
//
//	s := New()
//	s.Push(node.New(1, nil, nil))
//	fmt.Printf("Stack size: %d", s.Size()) // Output: Stack size: 1
func NewStack() *Stack {
	return &Stack{
		list: New(),
	}
}

// Push adds an element to the top of the stack.
//
// This operation is O(1) and increases the stack size by 1.
//
// Parameters:
//   - n: The Node to add to the top of the stack. Must not be nil.
//
// Example:
//
//	s := New()
//	s.Push(node.New(1, nil, nil))
//	s.Push(node.New(2, nil, nil))
//	// Stack now has 2 at the top, 1 at the bottom
func (s *Stack) Push(n *node.Node) {
	s.list.Push(n)
}

// Pop removes and returns the element from the top of the stack.
//
// This operation is O(1) and decreases the stack size by 1. The returned
// node is a copy with nil next/prev pointers.
//
// Returns:
//   - A copy of the removed Node, or nil if the stack is empty
//
// Example:
//
//	s := New()
//	s.Push(node.New(1, nil, nil))
//	s.Push(node.New(2, nil, nil))
//	n := s.Pop() // Returns node with ID 2
//	if n != nil {
//		fmt.Printf("Popped: %d", n.ID())
//	}
func (s *Stack) Pop() *node.Node {
	return s.list.Pop()
}

// Peek returns the element at the top of the stack without removing it.
//
// This operation is O(1) and does not modify the stack.
//
// Returns:
//   - A copy of the top Node and true, or an empty Node and false if the stack is empty
//
// Example:
//
//	s := New()
//	s.Push(node.New(1, nil, nil))
//	n, ok := s.Peek()
//	if ok {
//		fmt.Printf("Top: %d", n.ID()) // Output: Top: 1
//	}
//	fmt.Printf("Size: %d", s.Size()) // Output: Size: 1 (unchanged)
func (s *Stack) Peek() (node.Node, bool) {
	return s.list.Tail()
}

// Size returns the current number of elements in the stack.
//
// Returns:
//   - The current number of elements in the stack
func (s *Stack) Size() int {
	return s.list.Size()
}

// IsEmpty returns true if the stack contains no elements.
//
// Returns:
//   - true if the stack is empty, false otherwise
func (s *Stack) IsEmpty() bool {
	return s.list.Size() == 0
}
