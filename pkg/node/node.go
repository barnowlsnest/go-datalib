package node

// Node represents a node in a doubly-linked list structure.
//
// Each Node contains a unique identifier and maintains bidirectional
// references to adjacent nodes in the list. This structure serves as
// the fundamental building block for implementing various data structures
// like queues, stacks, and general-purpose linked lists.
//
// The Node is designed to be lightweight and memory-efficient, storing
// only the essential data needed for list operations. All fields are
// private to ensure data integrity and prevent external manipulation
// of the list structure.
//
// Thread Safety:
// Node itself is not thread-safe. Concurrent access to Node instances
// should be synchronized by the containing data structure (e.g., LinkedList).
type Node struct {
	// id is the unique identifier for this node.
	// It's immutable after creation and used to identify the node's data.
	id uint64

	// next points to the next node in the list, or nil if this is the last node.
	// This field is managed by the containing list structure.
	next *Node

	// prev points to the previous node in the list, or nil if this is the first node.
	// This field is managed by the containing list structure.
	prev *Node
}

// New creates a new Node with the specified ID and link references.
//
// This constructor allows creating a Node with predefined next and previous
// references, which is useful when inserting nodes into existing list structures
// or when building lists programmatically.
//
// Parameters:
//   - id: The unique identifier for this node
//   - next: Pointer to the next node in the list, or nil if none
//   - prev: Pointer to the previous node in the list, or nil if none
//
// Returns:
//   - A new Node instance with the specified configuration
//
// Example:
//
//	// Create a standalone node
//	node := New(1, nil, nil)
//
//	// Create a node linked between two existing nodes
//	middle := New(2, nextNode, prevNode)
func New(id uint64, next, prev *Node) *Node {
	return &Node{
		id:   id,
		next: next,
		prev: prev,
	}
}

// ID returns the unique identifier of this node.
//
// The ID is immutable after node creation and serves as the primary
// way to identify and reference the data associated with this node.
//
// Returns:
//   - The ID assigned to this node during creation
func (node *Node) ID() uint64 {
	return node.id
}

// Next returns the next node in the list.
//
// This method provides read-only access to the next node reference.
// The returned pointer may be nil if this node is the last in the list.
//
// Returns:
//   - Pointer to the next Node, or nil if this is the last node
func (node *Node) Next() *Node {
	return node.next
}

// Prev returns the previous node in the list.
//
// This method provides read-only access to the previous node reference.
// The returned pointer may be nil if this node is the first in the list.
//
// Returns:
//   - Pointer to the previous Node, or nil if this is the first node
func (node *Node) Prev() *Node {
	return node.prev
}

// WithPrev sets the previous node reference.
//
// This method allows modifying the previous node reference after creation,
// enabling dynamic list construction and manipulation. Setting n to nil
// will clear the previous node reference.
//
// Parameters:
//   - n: Pointer to the node to set as previous, or nil to clear
func (node *Node) WithPrev(n *Node) {
	node.prev = n
}

// WithNext sets the next node reference.
//
// This method allows modifying the next node reference after creation,
// enabling dynamic list construction and manipulation. Setting n to nil
// will clear the next node reference.
//
// Parameters:
//   - n: Pointer to the node to set as next, or nil to clear
func (node *Node) WithNext(n *Node) {
	node.next = n
}
