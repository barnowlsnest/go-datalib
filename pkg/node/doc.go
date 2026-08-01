// Package node provides the fundamental Node type for doubly-linked structures.
//
// A Node carries an immutable uint64 identifier and mutable references to its
// adjacent nodes. All fields are private; next and prev are set through the
// chainable WithNext and WithPrev methods, which makes Node the building block
// for the higher-level containers in this module (linked lists, stacks, queues
// and binary trees).
//
// The package offers:
//   - Node construction via New and the ID shorthand
//   - Chainable link mutation with WithNext and WithPrev
//   - Explicit traversal through the Iterable interface and its ForwardIterator
//     and BackwardIterator implementations, signaling exhaustion with ErrEOI
//   - Range-over-func traversal via NextNodes and PrevNodes, which return
//     iter.Seq2 sequences of position and node
//
// Nodes are not thread-safe; concurrent mutation requires external
// synchronization.
package node
