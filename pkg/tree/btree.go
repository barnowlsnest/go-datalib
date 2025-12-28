package tree

import (
	"cmp"
	"iter"
)

const (
	// DefaultMinDegree is the default minimum degree (t) for a B-tree.
	// A B-tree with minimum degree t has:
	//   - Each node (except root) has at least t-1 keys
	//   - Each node has at most 2t-1 keys
	//   - Each internal node (except root) has at least t children
	//   - Each internal node has at most 2t children
	DefaultMinDegree = 2
)

type (
	// BTreeEntry represents a key-value pair stored in the B-tree.
	BTreeEntry[K cmp.Ordered, V any] struct {
		Key   K
		Value V
	}

	// btreeNode represents an internal node in the B-tree.
	btreeNode[K cmp.Ordered, V any] struct {
		entries  []BTreeEntry[K, V]
		children []*btreeNode[K, V]
		leaf     bool
	}

	// BTree is a self-balancing tree data structure that maintains sorted data
	// and allows searches, sequential access, insertions, and deletions in
	// logarithmic time. This implementation is optimized for in-memory use
	// and is suitable for indexing message offsets in a commit log.
	BTree[K cmp.Ordered, V any] struct {
		root      *btreeNode[K, V]
		minDegree int
		size      int
	}

	// BTreeOption is a functional option for configuring a BTree during creation.
	BTreeOption[K cmp.Ordered, V any] func(t *BTree[K, V])
)

// NewBTree creates a new B-tree with the specified minimum degree.
// If minDegree < 2, DefaultMinDegree (2) is used.
//
// The minimum degree t determines the node capacity:
//   - Each node (except root) has at least t-1 keys
//   - Each node has at most 2t-1 keys
//
// Example:
//
//	tree := NewBTree[uint64, string](3)
//	tree.Insert(1, "first message")
//	tree.Insert(2, "second message")
func NewBTree[K cmp.Ordered, V any](minDegree int, opts ...BTreeOption[K, V]) *BTree[K, V] {
	if minDegree < 2 {
		minDegree = DefaultMinDegree
	}

	t := &BTree[K, V]{
		minDegree: minDegree,
		size:      0,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// newNode creates a new B-tree node.
func newNode[K cmp.Ordered, V any](minDegree int, leaf bool) *btreeNode[K, V] {
	return &btreeNode[K, V]{
		entries:  make([]BTreeEntry[K, V], 0, 2*minDegree-1),
		children: make([]*btreeNode[K, V], 0, 2*minDegree),
		leaf:     leaf,
	}
}

// Size returns the number of entries in the B-tree.
func (t *BTree[K, V]) Size() int {
	return t.size
}

// IsEmpty returns true if the B-tree contains no entries.
func (t *BTree[K, V]) IsEmpty() bool {
	return t.size == 0
}

// MinDegree returns the minimum degree of the B-tree.
func (t *BTree[K, V]) MinDegree() int {
	return t.minDegree
}

// Height returns the height of the B-tree.
// An empty tree has height 0.
func (t *BTree[K, V]) Height() int {
	if t.root == nil {
		return 0
	}

	height := 1
	node := t.root
	for !node.leaf {
		height++
		node = node.children[0]
	}

	return height
}

// Insert adds a key-value pair to the B-tree.
// If the key already exists, the value is updated.
func (t *BTree[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = newNode[K, V](t.minDegree, true)
		t.root.entries = append(t.root.entries, BTreeEntry[K, V]{Key: key, Value: value})
		t.size++
		return
	}

	// Check if key exists and update
	if t.update(t.root, key, value) {
		return
	}

	// If root is full, split it
	if len(t.root.entries) == 2*t.minDegree-1 {
		newRoot := newNode[K, V](t.minDegree, false)
		newRoot.children = append(newRoot.children, t.root)
		t.splitChild(newRoot, 0)
		t.root = newRoot
	}

	t.insertNonFull(t.root, key, value)
	t.size++
}

// update attempts to update an existing key's value.
// Returns true if key was found and updated, false otherwise.
func (t *BTree[K, V]) update(node *btreeNode[K, V], key K, value V) bool {
	i := 0
	for i < len(node.entries) && key > node.entries[i].Key {
		i++
	}

	if i < len(node.entries) && key == node.entries[i].Key {
		node.entries[i].Value = value
		return true
	}

	if node.leaf {
		return false
	}

	return t.update(node.children[i], key, value)
}

// splitChild splits the i-th child of parent when it's full.
func (t *BTree[K, V]) splitChild(parent *btreeNode[K, V], i int) {
	minDeg := t.minDegree
	fullChild := parent.children[i]
	newChild := newNode[K, V](minDeg, fullChild.leaf)

	// Move the upper half of entries to new child
	midIndex := minDeg - 1
	newChild.entries = append(newChild.entries, fullChild.entries[midIndex+1:]...)

	// Move the upper half of children to new child (if not leaf)
	if !fullChild.leaf {
		newChild.children = append(newChild.children, fullChild.children[minDeg:]...)
		fullChild.children = fullChild.children[:minDeg]
	}

	// Get the median entry to promote
	medianEntry := fullChild.entries[midIndex]
	fullChild.entries = fullChild.entries[:midIndex]

	// Insert new child into parent's children
	parent.children = append(parent.children, nil)
	copy(parent.children[i+2:], parent.children[i+1:])
	parent.children[i+1] = newChild

	// Insert median entry into parent
	parent.entries = append(parent.entries, BTreeEntry[K, V]{})
	copy(parent.entries[i+1:], parent.entries[i:])
	parent.entries[i] = medianEntry
}

// insertNonFull inserts a key-value pair into a non-full node.
func (t *BTree[K, V]) insertNonFull(node *btreeNode[K, V], key K, value V) {
	i := len(node.entries) - 1

	if node.leaf {
		// Find position and insert
		node.entries = append(node.entries, BTreeEntry[K, V]{})
		for i >= 0 && key < node.entries[i].Key {
			node.entries[i+1] = node.entries[i]
			i--
		}
		node.entries[i+1] = BTreeEntry[K, V]{Key: key, Value: value}
		return
	}

	// Find child to descend to
	for i >= 0 && key < node.entries[i].Key {
		i--
	}
	i++

	// Split child if full
	if len(node.children[i].entries) == 2*t.minDegree-1 {
		t.splitChild(node, i)
		if key > node.entries[i].Key {
			i++
		}
	}

	t.insertNonFull(node.children[i], key, value)
}

// Search finds the value associated with the given key.
// Returns the value and true if found, zero value and false otherwise.
func (t *BTree[K, V]) Search(key K) (V, bool) {
	if t.root == nil {
		var zero V
		return zero, false
	}

	return t.search(t.root, key)
}

func (t *BTree[K, V]) search(node *btreeNode[K, V], key K) (V, bool) {
	i := 0
	for i < len(node.entries) && key > node.entries[i].Key {
		i++
	}

	if i < len(node.entries) && key == node.entries[i].Key {
		return node.entries[i].Value, true
	}

	if node.leaf {
		var zero V
		return zero, false
	}

	return t.search(node.children[i], key)
}

// Contains returns true if the key exists in the B-tree.
func (t *BTree[K, V]) Contains(key K) bool {
	_, found := t.Search(key)
	return found
}

// Delete removes a key from the B-tree.
// Returns true if the key was found and deleted, false otherwise.
func (t *BTree[K, V]) Delete(key K) bool {
	if t.root == nil {
		return false
	}

	deleted := t.delete(t.root, key)
	if deleted {
		t.size--

		// If root has no entries and has a child, make that child the new root
		if len(t.root.entries) == 0 {
			if t.root.leaf {
				t.root = nil
			} else {
				t.root = t.root.children[0]
			}
		}
	}

	return deleted
}

func (t *BTree[K, V]) delete(node *btreeNode[K, V], key K) bool {
	i := 0
	for i < len(node.entries) && key > node.entries[i].Key {
		i++
	}

	// Case 1: Key is in this node
	if i < len(node.entries) && key == node.entries[i].Key {
		if node.leaf {
			// Case 1a: Node is a leaf, simply remove the key
			node.entries = append(node.entries[:i], node.entries[i+1:]...)
			return true
		}

		// Case 1b: Node is internal
		return t.deleteFromInternal(node, i)
	}

	// Key is not in this node
	if node.leaf {
		return false
	}

	// Case 3: Key might be in child[i]
	return t.deleteFromChild(node, i, key)
}

// deleteFromInternal handles deletion when key is in an internal node.
func (t *BTree[K, V]) deleteFromInternal(node *btreeNode[K, V], i int) bool {
	minDeg := t.minDegree

	// Case 2a: Left child has >= t keys
	if len(node.children[i].entries) >= minDeg {
		pred := t.getPredecessor(node.children[i])
		node.entries[i] = pred
		return t.delete(node.children[i], pred.Key)
	}

	// Case 2b: Right child has >= t keys
	if len(node.children[i+1].entries) >= minDeg {
		succ := t.getSuccessor(node.children[i+1])
		node.entries[i] = succ
		return t.delete(node.children[i+1], succ.Key)
	}

	// Case 2c: Both children have t-1 keys, merge them
	// Save the key before merging since it moves to the child
	keyToDelete := node.entries[i].Key
	t.merge(node, i)
	return t.delete(node.children[i], keyToDelete)
}

// deleteFromChild handles deletion when key might be in a child.
func (t *BTree[K, V]) deleteFromChild(node *btreeNode[K, V], i int, key K) bool {
	minDeg := t.minDegree
	child := node.children[i]

	// If child has only t-1 keys, we need to ensure it has at least t keys
	if len(child.entries) < minDeg {
		// Try to borrow from left sibling
		if i > 0 && len(node.children[i-1].entries) >= minDeg {
			t.borrowFromLeft(node, i)
		} else if i < len(node.children)-1 && len(node.children[i+1].entries) >= minDeg {
			// Try to borrow from right sibling
			t.borrowFromRight(node, i)
		} else {
			// Merge with a sibling
			if i < len(node.children)-1 {
				t.merge(node, i)
			} else {
				t.merge(node, i-1)
				i--
			}
		}
	}

	return t.delete(node.children[i], key)
}

// getPredecessor returns the predecessor (largest key in left subtree).
func (t *BTree[K, V]) getPredecessor(node *btreeNode[K, V]) BTreeEntry[K, V] {
	for !node.leaf {
		node = node.children[len(node.children)-1]
	}
	return node.entries[len(node.entries)-1]
}

// getSuccessor returns the successor (smallest key in right subtree).
func (t *BTree[K, V]) getSuccessor(node *btreeNode[K, V]) BTreeEntry[K, V] {
	for !node.leaf {
		node = node.children[0]
	}
	return node.entries[0]
}

// borrowFromLeft borrows an entry from the left sibling.
func (t *BTree[K, V]) borrowFromLeft(parent *btreeNode[K, V], i int) {
	child := parent.children[i]
	leftSibling := parent.children[i-1]

	// Move parent entry down to child
	child.entries = append([]BTreeEntry[K, V]{parent.entries[i-1]}, child.entries...)

	// Move last entry from left sibling up to parent
	parent.entries[i-1] = leftSibling.entries[len(leftSibling.entries)-1]
	leftSibling.entries = leftSibling.entries[:len(leftSibling.entries)-1]

	// Move last child from left sibling to child
	if !leftSibling.leaf {
		child.children = append([]*btreeNode[K, V]{leftSibling.children[len(leftSibling.children)-1]}, child.children...)
		leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
	}
}

// borrowFromRight borrows an entry from the right sibling.
func (t *BTree[K, V]) borrowFromRight(parent *btreeNode[K, V], i int) {
	child := parent.children[i]
	rightSibling := parent.children[i+1]

	// Move parent entry down to child
	child.entries = append(child.entries, parent.entries[i])

	// Move first entry from right sibling up to parent
	parent.entries[i] = rightSibling.entries[0]
	rightSibling.entries = rightSibling.entries[1:]

	// Move first child from right sibling to child
	if !rightSibling.leaf {
		child.children = append(child.children, rightSibling.children[0])
		rightSibling.children = rightSibling.children[1:]
	}
}

// merge merges child[i] with child[i+1].
func (t *BTree[K, V]) merge(parent *btreeNode[K, V], i int) {
	left := parent.children[i]
	right := parent.children[i+1]

	// Move parent entry down to left child
	left.entries = append(left.entries, parent.entries[i])

	// Move all entries from right child to left child
	left.entries = append(left.entries, right.entries...)

	// Move all children from right child to left child
	if !left.leaf {
		left.children = append(left.children, right.children...)
	}

	// Remove entry from parent
	parent.entries = append(parent.entries[:i], parent.entries[i+1:]...)

	// Remove right child from parent
	parent.children = append(parent.children[:i+1], parent.children[i+2:]...)
}

// Min returns the minimum key-value pair in the B-tree.
// Returns zero values and false if the tree is empty.
func (t *BTree[K, V]) Min() (key K, value V, found bool) {
	if t.root == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	node := t.root
	for !node.leaf {
		node = node.children[0]
	}

	entry := node.entries[0]
	return entry.Key, entry.Value, true
}

// Max returns the maximum key-value pair in the B-tree.
// Returns zero values and false if the tree is empty.
func (t *BTree[K, V]) Max() (key K, value V, found bool) {
	if t.root == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	node := t.root
	for !node.leaf {
		node = node.children[len(node.children)-1]
	}

	entry := node.entries[len(node.entries)-1]
	return entry.Key, entry.Value, true
}

// Range returns an iterator over all entries with keys in [from, to].
// The entries are yielded in ascending key order.
func (t *BTree[K, V]) Range(from, to K) iter.Seq[BTreeEntry[K, V]] {
	return func(yield func(BTreeEntry[K, V]) bool) {
		if t.root == nil || from > to {
			return
		}
		t.rangeTraverse(t.root, from, to, yield)
	}
}

func (t *BTree[K, V]) rangeTraverse(node *btreeNode[K, V], from, to K, yield func(BTreeEntry[K, V]) bool) bool {
	i := 0
	for i < len(node.entries) && node.entries[i].Key < from {
		i++
	}

	for i < len(node.entries) {
		// Visit left child if not a leaf
		if !node.leaf && i < len(node.children) {
			if !t.rangeTraverse(node.children[i], from, to, yield) {
				return false
			}
		}

		// Check if we've passed the upper bound
		if node.entries[i].Key > to {
			return true
		}

		// Yield the current entry
		if !yield(node.entries[i]) {
			return false
		}

		i++
	}

	// Visit rightmost child if not a leaf
	if !node.leaf && i < len(node.children) {
		return t.rangeTraverse(node.children[i], from, to, yield)
	}

	return true
}

// All returns an iterator over all entries in ascending key order.
func (t *BTree[K, V]) All() iter.Seq[BTreeEntry[K, V]] {
	return func(yield func(BTreeEntry[K, V]) bool) {
		if t.root == nil {
			return
		}
		t.inOrderTraverse(t.root, yield)
	}
}

func (t *BTree[K, V]) inOrderTraverse(node *btreeNode[K, V], yield func(BTreeEntry[K, V]) bool) bool {
	for i := 0; i < len(node.entries); i++ {
		// Visit left child if not a leaf
		if !node.leaf {
			if !t.inOrderTraverse(node.children[i], yield) {
				return false
			}
		}

		// Yield the current entry
		if !yield(node.entries[i]) {
			return false
		}
	}

	// Visit rightmost child if not a leaf
	if !node.leaf {
		return t.inOrderTraverse(node.children[len(node.children)-1], yield)
	}

	return true
}

// Clear removes all entries from the B-tree.
func (t *BTree[K, V]) Clear() {
	t.root = nil
	t.size = 0
}

// Floor returns the largest entry with a key <= the given key.
// Returns zero values and false if no such entry exists.
func (t *BTree[K, V]) Floor(key K) (floorKey K, floorValue V, found bool) {
	if t.root == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	entry, found := t.floor(t.root, key)
	if !found {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	return entry.Key, entry.Value, true
}

func (t *BTree[K, V]) floor(node *btreeNode[K, V], key K) (BTreeEntry[K, V], bool) {
	i := 0
	for i < len(node.entries) && key >= node.entries[i].Key {
		i++
	}

	// Check if we found an exact match at i-1
	if i > 0 && node.entries[i-1].Key == key {
		return node.entries[i-1], true
	}

	// If leaf, return the largest key less than target
	if node.leaf {
		if i > 0 {
			return node.entries[i-1], true
		}
		return BTreeEntry[K, V]{}, false
	}

	// Try to find in appropriate child
	result, found := t.floor(node.children[i], key)
	if found {
		return result, true
	}

	// If not found in child and we have a smaller key in this node
	if i > 0 {
		return node.entries[i-1], true
	}

	return BTreeEntry[K, V]{}, false
}

// Ceiling returns the smallest entry with a key >= the given key.
// Returns zero values and false if no such entry exists.
func (t *BTree[K, V]) Ceiling(key K) (ceilingKey K, ceilingValue V, found bool) {
	if t.root == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	entry, found := t.ceiling(t.root, key)
	if !found {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	return entry.Key, entry.Value, true
}

func (t *BTree[K, V]) ceiling(node *btreeNode[K, V], key K) (BTreeEntry[K, V], bool) {
	i := 0
	for i < len(node.entries) && key > node.entries[i].Key {
		i++
	}

	// Check if we found an exact match or a larger key
	if i < len(node.entries) && node.entries[i].Key >= key {
		// If leaf, this is our answer
		if node.leaf {
			return node.entries[i], true
		}

		// Try to find smaller ceiling in left child
		result, found := t.ceiling(node.children[i], key)
		if found {
			return result, true
		}

		return node.entries[i], true
	}

	// If leaf and no valid key found
	if node.leaf {
		return BTreeEntry[K, V]{}, false
	}

	// Try rightmost child
	return t.ceiling(node.children[i], key)
}

// Keys returns all keys in ascending order.
func (t *BTree[K, V]) Keys() []K {
	keys := make([]K, 0, t.size)
	for entry := range t.All() {
		keys = append(keys, entry.Key)
	}
	return keys
}

// Values returns all values in key-ascending order.
func (t *BTree[K, V]) Values() []V {
	values := make([]V, 0, t.size)
	for entry := range t.All() {
		values = append(values, entry.Value)
	}
	return values
}
