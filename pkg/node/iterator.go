package node

type (
	// probe is an internal function type used to retrieve the next node during iteration.
	probe func() *Node

	// baseIter is the internal state shared by ForwardIterator and BackwardIterator.
	// It tracks the current position and the last visited node in the iteration.
	baseIterator struct {
		cur *Node // Current position in the iteration
	}

	BackwardIterator struct {
		baseIterator
	}

	ForwardIterator struct {
		baseIterator
	}
)

func (it *baseIterator) hasNext() bool {
	if it.cur == nil {
		return false
	}

	return it.cur.Next() != nil
}

func (it *baseIterator) hasPrev() bool {
	if it.cur == nil {
		return false
	}

	return it.cur.Prev() != nil
}

func (it *baseIterator) take(fn probe) (*Node, error) {
	if it.cur == nil {
		return nil, ErrEOI
	}

	it.cur = fn()
	return it.cur, nil
}

func (it *baseIterator) nextForward() (*Node, error) {
	n := it.cur
	return it.take(n.Next)
}

func (it *baseIterator) nextBackward() (*Node, error) {
	n := it.cur
	return it.take(n.Prev)
}

// Curr returns the current node without advancing the iterator.
//
// Returns:
//   - The current node, or ErrEOI if iteration has completed
func (it *baseIterator) Curr() (*Node, error) {
	if it.cur == nil {
		return nil, ErrEOI
	}

	return it.cur, nil
}

func Forward(n *Node) *ForwardIterator {
	return &ForwardIterator{baseIterator{n}}
}

func Backward(n *Node) *BackwardIterator {
	return &BackwardIterator{baseIterator{n}}
}

// Next advances the iterator backward and returns the previous node.
//
// Note: Despite the name "Next", this method moves backward through the list
// following Prev() pointers to maintain consistency with the Iterator interface.
//
// Returns:
//   - The previous node in the backward direction, or nil with ErrEOI if no more nodes exist
func (b *BackwardIterator) Next() (*Node, error) {
	return b.nextBackward()
}

// HasNext returns true if there is a previous node available in backward direction.
//
// Returns:
//   - true if calling Next() will return a node, false otherwise
func (b *BackwardIterator) HasNext() bool {
	return b.hasPrev()
}

func (f *ForwardIterator) Next() (*Node, error) {
	return f.nextForward()
}

func (f *ForwardIterator) HasNext() bool {
	return f.hasNext()
}
