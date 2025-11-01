package node

import (
	"iter"
)

func move(it Iterable) iter.Seq2[uint64, *Node] {
	return func(yield func(uint64, *Node) bool) {
		var (
			ne  *Node
			err error
		)
		b, err := it.Curr()
		if err != nil {
			return
		}

		i := b.ID()
		for it.HasNext() {
			ne, err = it.Next()
			if err != nil {
				return
			}
			if !yield(i, ne) {
				return
			}
		}
	}
}

// Iterable defines the interface for traversing a doubly-linked list of nodes.
//
// Implementations include ForwardIterator (traverses via Next() pointers) and
// BackwardIterator (traverses via Prev() pointers).
//
// The iterator maintains state and will return ErrEOI when the end is reached.
// After reaching the end, Last() can be used to retrieve the final valid node.
type Iterable interface {
	// Next advances the iterator and returns the next node.
	// Returns ErrEOI if the end of iteration is reached.
	Next() (*Node, error)

	// HasNext returns true if there are more nodes to iterate over.
	HasNext() bool

	// Curr returns the current node without advancing the iterator.
	// Returns ErrEOI if the iterator has exhausted all nodes.
	Curr() (*Node, error)
}

// NextNodes creates a Go 1.23+ iterator that traverses nodes in forward direction.
//
// This function returns an iter.Seq2 that can be used with range loops in Go 1.23+.
// The iterator yields pairs of (node_id, node) as it traverses from the starting
// node through all subsequent nodes via Next() pointers.
//
// Parameters:
//   - n: The starting node for forward traversal
//
// Returns:
//   - An iter.Seq2[uint64, *Node] that yields (node_id, node) pairs
//
// Example:
//
//	for id, node := range NextNodes(startNode) {
//	    fmt.Printf("Node %d\n", id)
//	}
func NextNodes(n *Node) iter.Seq2[uint64, *Node] {
	return move(Forward(n))
}

// PrevNodes creates a Go 1.23+ iterator that traverses nodes in backward direction.
//
// This function returns an iter.Seq2 that can be used with range loops in Go 1.23+.
// The iterator yields pairs of (node_id, node) as it traverses from the starting
// node through all previous nodes via Prev() pointers.
//
// Parameters:
//   - n: The starting node for backward traversal
//
// Returns:
//   - An iter.Seq2[uint64, *Node] that yields (node_id, node) pairs
//
// Example:
//
//	for id, node := range PrevNodes(endNode) {
//	    fmt.Printf("Node %d\n", id)
//	}
func PrevNodes(n *Node) iter.Seq2[uint64, *Node] {
	return move(Backward(n))
}
