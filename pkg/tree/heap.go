package tree

import (
	"cmp"
)

// Heap is a generic binary heap that can function as either a min-heap or max-heap.
// It provides O(log n) insertion and removal operations with O(1) peek.
//
// The heap is implemented using a slice for optimal cache locality and performance.
// Elements are stored in level-order, where for element at index i:
//   - Parent is at index (i-1)/2
//   - Left child is at index 2*i+1
//   - Right child is at index 2*i+2
type Heap[T any] struct {
	data []T
	less func(T, T) bool // Comparison function: less(a, b) returns true if a should be higher in heap than b
}

// NewHeap creates a new empty heap with the given comparison function.
// For a min-heap, pass a "less than" function: func(a, b T) bool { return a < b }
// For a max-heap, pass a "greater than" function: func(a, b T) bool { return a > b }
//
// Example min-heap:
//
//	h := heap.New(func(a, b int) bool { return a < b })
//
// Example max-heap:
//
//	h := heap.New(func(a, b int) bool { return a > b })
func NewHeap[T any](less func(T, T) bool) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: less,
	}
}

// NewMin creates a new min-heap for ordered types (types that support <, >, etc).
// The minimum element will always be at the top.
//
// Example:
//
//	h := heap.NewMin[int]()
func NewMin[T cmp.Ordered]() *Heap[T] {
	return NewHeap(func(a, b T) bool { return a < b })
}

// NewMax creates a new max-heap for ordered types (types that support <, >, etc).
// The maximum element will always be at the top.
//
// Example:
//
//	h := heap.NewMax[int]()
func NewMax[T cmp.Ordered]() *Heap[T] {
	return NewHeap(func(a, b T) bool { return a > b })
}

// NewWithCapacity creates a new heap with pre-allocated capacity.
// This can improve performance when the expected size is known in advance.
func NewWithCapacity[T any](less func(T, T) bool, capacity int) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0, capacity),
		less: less,
	}
}

// HeapFromSlice creates a heap from an existing slice using heapify.
// This is more efficient than inserting elements one by one: O(n) vs O(n log n).
//
// The input slice is copied, so modifications to the heap won't affect the original slice.
func HeapFromSlice[T any](slice []T, less func(T, T) bool) *Heap[T] {
	data := make([]T, len(slice))
	copy(data, slice)

	h := &Heap[T]{
		data: data,
		less: less,
	}

	h.heapify()
	return h
}

// Push adds a new element to the heap.
// Time complexity: O(log n)
func (h *Heap[T]) Push(value T) {
	h.data = append(h.data, value)
	h.bubbleUp(len(h.data) - 1)
}

// Pop removes and returns the top element (min or max depending on heap type).
// Returns the element and true if successful, or zero value and false if heap is empty.
// Time complexity: O(log n)
func (h *Heap[T]) Pop() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}

	root := h.data[0]
	lastIdx := len(h.data) - 1

	// Move last element to root
	h.data[0] = h.data[lastIdx]
	h.data = h.data[:lastIdx]

	// Restore heap property
	if len(h.data) > 0 {
		h.bubbleDown(0)
	}

	return root, true
}

// Peek returns the top element without removing it.
// Returns the element and true if successful, or zero value and false if heap is empty.
// Time complexity: O(1)
func (h *Heap[T]) Peek() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}
	return h.data[0], true
}

// Size returns the number of elements in the heap.
// Time complexity: O(1)
func (h *Heap[T]) Size() int {
	return len(h.data)
}

// IsEmpty returns true if the heap contains no elements.
// Time complexity: O(1)
func (h *Heap[T]) IsEmpty() bool {
	return len(h.data) == 0
}

// Clear removes all elements from the heap.
// Time complexity: O(1)
func (h *Heap[T]) Clear() {
	h.data = h.data[:0]
}

// ToSlice returns a copy of the heap's internal data as a slice.
// The slice is in heap order (level-order), not sorted order.
// Time complexity: O(n)
func (h *Heap[T]) ToSlice() []T {
	result := make([]T, len(h.data))
	copy(result, h.data)
	return result
}

// parent returns the index of the parent node.
func parent(i int) int {
	return (i - 1) / 2
}

// leftChild returns the index of the left child node.
func leftChild(i int) int {
	return 2*i + 1
}

// rightChild returns the index of the right child node.
func rightChild(i int) int {
	return 2*i + 2
}

// bubbleUp moves an element up the heap until the heap property is restored.
// Used after insertion.
func (h *Heap[T]) bubbleUp(i int) {
	for i > 0 {
		p := parent(i)
		// If current element should be higher than parent, swap
		if h.less(h.data[i], h.data[p]) {
			h.data[i], h.data[p] = h.data[p], h.data[i]
			i = p
		} else {
			break
		}
	}
}

// bubbleDown moves an element down the heap until the heap property is restored.
// Used after removal.
func (h *Heap[T]) bubbleDown(i int) {
	n := len(h.data)

	for {
		smallest := i
		left := leftChild(i)
		right := rightChild(i)

		// Find the element that should be highest among parent and children
		if left < n && h.less(h.data[left], h.data[smallest]) {
			smallest = left
		}
		if right < n && h.less(h.data[right], h.data[smallest]) {
			smallest = right
		}

		// If parent is already in correct position, we're done
		if smallest == i {
			break
		}

		// Swap and continue
		h.data[i], h.data[smallest] = h.data[smallest], h.data[i]
		i = smallest
	}
}

// heapify converts an arbitrary slice into a valid heap.
// This is more efficient than inserting elements one by one.
// Time complexity: O(n)
func (h *Heap[T]) heapify() {
	// Start from the last non-leaf node and work backwards
	// Last non-leaf node is at index (n/2 - 1)
	n := len(h.data)
	for i := n/2 - 1; i >= 0; i-- {
		h.bubbleDown(i)
	}
}
