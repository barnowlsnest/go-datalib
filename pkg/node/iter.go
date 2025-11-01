package node

import (
	"iter"
)

func move(it Iterable) iter.Seq2[int, *Node] {
	return func(yield func(int, *Node) bool) {
		var (
			ne  *Node
			err error
		)

		// First yield the current node
		ne, err = it.Curr()
		if err != nil {
			return
		}

		var i int
		if !yield(i, ne) {
			return
		}
		i++

		// Then iterate through remaining nodes
		for it.HasNext() {
			ne, err = it.Next()
			if err != nil {
				return
			}
			if !yield(i, ne) {
				return
			}
			i++
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

func NextNodes(n *Node) iter.Seq2[int, *Node] {
	return move(Forward(n))
}

func PrevNodes(n *Node) iter.Seq2[int, *Node] {
	return move(Backward(n))
}
