package tree

import (
	"golang.org/x/exp/constraints"
)

// Fenwick is a data structure that efficiently
// supports prefix sum queries and point updates in O(log n) time.
//
// The tree uses 1-based indexing internally, where index relationships are determined
// by bit manipulation rather than explicit pointers:
//   - Parent of index i: i - (i & -i)
//   - Next index: i + (i & -i)
//
// Common use cases:
//   - Cumulative frequency tables
//   - Range sum queries with updates
//   - Counting inversions
//   - 2D range sum queries (with 2D Fenwick Fenwick)
type Fenwick[T constraints.Unsigned | constraints.Integer | constraints.Float] struct {
	tree []T
	n    int
}

// NewFenwick creates a new Fenwick with the given size.
// The tree is initialized with all zeros.
//
// Example:
//
//	ft := fenwick.New[int](10)
func NewFenwick[T constraints.Integer | constraints.Float](size int) *Fenwick[T] {
	if size < 0 {
		size = 0
	}
	return &Fenwick[T]{
		tree: make([]T, size+1), // index 0 is unused, indices 1..n are used
		n:    size,
	}
}

// FromSlice creates a Fenwick from an existing slice.
// The input slice is treated as 0-indexed, but internally the tree uses 1-based indexing.
// This is more efficient than creating an empty tree and updating each element individually.
// Time complexity: O(n log n)
//
// Example:
//
//	data := []int{3, 2, -1, 6, 5, 4, -3, 3, 7, 2, 3}
//	ft := fenwick.FromSlice(data)
func FromSlice[T constraints.Integer | constraints.Float](data []T) *Fenwick[T] {
	n := len(data)
	tree := &Fenwick[T]{
		tree: make([]T, n+1),
		n:    n,
	}

	// Build tree by updating each element
	for i := 0; i < n; i++ {
		tree.Update(i+1, data[i]) // Convert to 1-indexed
	}

	return tree
}

// Size returns the size of the Fenwick.
// Time complexity: O(1)
func (t *Fenwick[T]) Size() int {
	return t.n
}

// Update adds delta to the element at the given 1-based index.
// The update propagates to all relevant ranges in O(log n) time.
//
// Time complexity: O(log n)
//
// Example:
//
//	ft.Update(3, 5)  // Add 5 to index 3
//	ft.Update(3, -2) // Subtract 2 from index 3
func (t *Fenwick[T]) Update(index int, delta T) {
	if index <= 0 || index > t.n {
		return // Out of bounds, silently ignore
	}

	// Propagate update upward through the tree
	for index <= t.n {
		t.tree[index] += delta
		index += index & -index // Move to parent (add lowest set bit)
	}
}

// Query returns the prefix sum from index 1 to the given 1-based index (inclusive).
// Time complexity: O(log n)
//
// Example:
//
//	sum := ft.Query(5) // Sum of elements from index 1 to 5
func (t *Fenwick[T]) Query(index int) T {
	if index <= 0 {
		var zero T
		return zero
	}
	if index > t.n {
		index = t.n
	}

	var sum T
	// Traverse downward through the tree
	for index > 0 {
		sum += t.tree[index]
		index -= index & -index // Move to next range (remove lowest set bit)
	}

	return sum
}

// RangeQuery returns the sum of elements in the range [left, right] (1-based, inclusive).
// Time complexity: O(log n)
//
// Example:
//
//	sum := ft.RangeQuery(3, 7) // Sum of elements from index 3 to 7
func (t *Fenwick[T]) RangeQuery(left, right int) T {
	if left > right || left <= 0 || right > t.n {
		var zero T
		return zero
	}

	if left == 1 {
		return t.Query(right)
	}

	return t.Query(right) - t.Query(left-1)
}

// Set sets the element at the given 1-based index to the specified value.
// This is implemented as: Update(index, newValue - currentValue)
// Time complexity: O(log n)
//
// Example:
//
//	ft.Set(3, 10) // Set index 3 to value 10
func (t *Fenwick[T]) Set(index int, value T) {
	if index <= 0 || index > t.n {
		return
	}

	// Calculate current value and update with the difference
	currentValue := t.Get(index)
	t.Update(index, value-currentValue)
}

// Get returns the value at the given 1-based index.
// This is implemented as: Query(index) - Query(index-1)
// Time complexity: O(log n)
//
// Example:
//
//	val := ft.Get(3) // Get value at index 3
func (t *Fenwick[T]) Get(index int) T {
	if index <= 0 || index > t.n {
		var zero T
		return zero
	}

	if index == 1 {
		return t.Query(1)
	}

	return t.Query(index) - t.Query(index-1)
}

// Clear resets all elements in the Fenwick to zero.
// Time complexity: O(n)
func (t *Fenwick[T]) Clear() {
	for i := range t.tree {
		t.tree[i] = 0
	}
}

// ToSlice returns a 0-indexed slice containing all values in the Fenwick.
// The returned slice is a copy, so modifications won't affect the tree.
// Time complexity: O(n log n)
//
// Example:
//
//	data := ft.ToSlice() // Returns []T with 0-indexed values
func (t *Fenwick[T]) ToSlice() []T {
	if t.n == 0 {
		return []T{}
	}

	result := make([]T, t.n)
	for i := 1; i <= t.n; i++ {
		result[i-1] = t.Get(i) // Convert from 1-indexed to 0-indexed
	}

	return result
}
