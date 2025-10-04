package dag

import "errors"

// Error definitions for the entity package.
// These errors are returned by various operations to indicate specific
// failure conditions that callers can handle appropriately.
var (
	// ErrGroupNotFound is returned when attempting to access a group
	// that doesn't exist in the graph structure.
	ErrGroupNotFound = errors.New("group not found")

	// ErrGroupAlreadyExists is returned when attempting to create a group
	// that already exists in the graph structure.
	ErrGroupAlreadyExists = errors.New("group already exists")

	// ErrNodeNotFound is returned when attempting to access a node
	// that doesn't exist in the specified group or graph.
	ErrNodeNotFound = errors.New("node not found")

	// ErrInvalidEdge is returned when attempting to create or manipulate
	// an edge with invalid parameters (e.g., self-loops, duplicate edges).
	ErrInvalidEdge = errors.New("invalid edge")

	// ErrInvalidAdjacency is returned when adjacency operations fail
	// due to structural constraints or invalid node relationships.
	ErrInvalidAdjacency = errors.New("invalid adjacency")

	// ErrInvalidBackRef is returned when back-reference operations
	// encounter inconsistent or invalid reference states.
	ErrInvalidBackRef = errors.New("invalid backref")

	// ErrRecoverFromPanic is returned when a panic is recovered during
	// operation execution, allowing graceful error handling.
	ErrRecoverFromPanic = errors.New("recover from panic")
)
