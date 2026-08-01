// Package list provides a doubly-linked list and the stack and queue
// abstractions built on top of it.
//
// LinkedList links nodes from the node package and keeps references to both
// ends, so insertion and removal at the head or the tail are O(1). Removed
// nodes are returned as detached copies so the caller cannot reach back into
// the internal chain. Size is tracked incrementally and reported in O(1).
//
// The package offers:
//   - LinkedList with Push, Pop, Unshift and Shift, their PushID, PopID,
//     UnshiftID and ShiftID convenience forms, and the IterNext and IterPrev
//     range-over-func traversals
//   - Stack, a LIFO wrapper exposing Push, Pop, Peek and IsEmpty
//   - Queue, a FIFO wrapper exposing Enqueue, Dequeue, PeekFront, PeekRear
//     and IsEmpty
//
// None of these types are thread-safe; concurrent use requires external
// synchronization.
package list
