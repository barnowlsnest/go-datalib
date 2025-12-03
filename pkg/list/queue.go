package list

import (
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

// Queue implements a FIFO (First In, First Out) data structure.
//
// This implementation uses a LinkedList internally to provide O(1) enqueue and
// dequeue operations. The queue follows classic queue semantics where elements
// are added at the rear (tail) and removed from the front (head).
//
// Key features:
//   - O(1) enqueue and dequeue operations
//   - O(1) peek operations for both front and rear elements
//   - Automatic size tracking
//   - Safe handling of empty queue conditions
//
// Thread Safety:
// Queue is not thread-safe. Concurrent access requires external
// synchronization mechanisms.
//
// Memory Management:
// The queue creates copies of nodes during Dequeue() operations to prevent
// memory leaks and ensure clean separation from the internal structure.
type Queue struct {
	// list is the internal LinkedList used to store queue elements.
	// Elements are enqueued at the tail and dequeued from the head.
	list *LinkedList
}

// NewQueue creates a new empty Queue.
//
// The returned queue is ready for use and has zero size.
//
// Returns:
//   - A new empty Queue instance
//
// Example:
//
//	q := New()
//	q.Enqueue(node.New(1, nil, nil))
//	fmt.Printf("Queue size: %d", q.Size()) // Output: Queue size: 1
func NewQueue() *Queue {
	return &Queue{
		list: New(),
	}
}

// Enqueue adds an element to the rear of the queue.
//
// This operation is O(1) and increases the queue size by 1.
//
// Parameters:
//   - n: The CreateNode to add to the rear of the queue. Must not be nil.
//
// Example:
//
//	q := New()
//	q.Enqueue(node.New(1, nil, nil))
//	q.Enqueue(node.New(2, nil, nil))
//	// Queue now has 1 at the front, 2 at the rear
func (q *Queue) Enqueue(n *node.Node) {
	q.list.Push(n)
}

// Dequeue removes and returns the element from the front of the queue.
//
// This operation is O(1) and decreases the queue size by 1. The returned
// node is a copy with nil next/prev pointers.
//
// Returns:
//   - A copy of the removed CreateNode, or nil if the queue is empty
//
// Example:
//
//	q := New()
//	q.Enqueue(node.New(1, nil, nil))
//	q.Enqueue(node.New(2, nil, nil))
//	n := q.Dequeue() // Returns node with ID 1
//	if n != nil {
//		fmt.Printf("Dequeued: %d", n.ID())
//	}
func (q *Queue) Dequeue() *node.Node {
	return q.list.Shift()
}

// PeekFront returns the element at the front of the queue without removing it.
//
// This operation is O(1) and does not modify the queue.
//
// Returns:
//   - A copy of the front CreateNode and true, or an empty CreateNode and false if the queue is empty
//
// Example:
//
//	q := New()
//	q.Enqueue(node.New(1, nil, nil))
//	n, ok := q.PeekFront()
//	if ok {
//		fmt.Printf("Front: %d", n.ID()) // Output: Front: 1
//	}
//	fmt.Printf("Size: %d", q.Size()) // Output: Size: 1 (unchanged)
func (q *Queue) PeekFront() (node.Node, bool) {
	return q.list.Head()
}

// PeekRear returns the element at the rear of the queue without removing it.
//
// This operation is O(1) and does not modify the queue.
//
// Returns:
//   - A copy of the rear CreateNode and true, or an empty CreateNode and false if the queue is empty
//
// Example:
//
//	q := New()
//	q.Enqueue(node.New(1, nil, nil))
//	q.Enqueue(node.New(2, nil, nil))
//	n, ok := q.PeekRear()
//	if ok {
//		fmt.Printf("Rear: %d", n.ID()) // Output: Rear: 2
//	}
func (q *Queue) PeekRear() (node.Node, bool) {
	return q.list.Tail()
}

// Size returns the current number of elements in the queue.
//
// Returns:
//   - The current number of elements in the queue
func (q *Queue) Size() int {
	return q.list.Size()
}

// IsEmpty returns true if the queue contains no elements.
//
// Returns:
//   - true if the queue is empty, false otherwise
func (q *Queue) IsEmpty() bool {
	return q.list.Size() == 0
}
