package node

import (
	"errors"
)

var (
	// ErrEOI indicates the end of iteration has been reached.
	//
	// This error is returned by iterator methods when attempting to
	// advance past the last available node or access the current node
	// after iteration has completed.
	ErrEOI = errors.New("end of iterator")

	// ErrNil indicates an operation was attempted on a nil node.
	//
	// This error is returned by methods that require a valid node
	// but received nil instead, such as when popping from an empty list.
	ErrNil = errors.New("node is nil")
)
