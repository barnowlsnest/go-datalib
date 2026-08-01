// Package tree provides generic tree data structures: a multi-way tree, a
// binary search tree, a B-tree, a binary heap, a Fenwick tree and a bounded
// tree segment.
//
// The package offers:
//   - Node[T comparable], a multi-way tree node with a configurable maximum
//     breadth, automatic level tracking, parent-child management and the
//     ValueOpt, LevelOpt, ParentOpt and ChildOpt construction options.
//     Children are attached, detached, moved and swapped through dedicated
//     methods, and selected with SelectChildrenFunc, SelectOneChildFunc or
//     the context-aware, concurrent SelectOneChildByEachValue.
//   - Hierarchy and ToModel, which convert between a HierarchyModel map and a
//     Node[string] tree rooted at RootTag, rejecting cycles.
//   - BinaryNode[T cmp.Ordered], a binary node embedding node.Node, and
//     BST[T cmp.Ordered] with Insert, Delete, Search, Height and the InOrder,
//     PreOrder, PostOrder and LevelOrder traversals.
//   - BTree[K cmp.Ordered, V any], a B-tree with a configurable minimum degree
//     (DefaultMinDegree) supporting Insert, Search, Delete and range iteration
//     over BTreeEntry values.
//   - Heap[T any], a binary heap driven by a caller-supplied less function,
//     with the NewMin and NewMax shorthands for ordered types.
//   - Fenwick[T], a binary indexed tree answering prefix and range sums in
//     logarithmic time via Update, Query and RangeQuery.
//   - Segment[T comparable], an aliased tree segment constrained by breadth and
//     depth (DefaultMaxDepth) that tracks capacity, links and unlinks nodes by
//     ID, removes them cascading or promoting their children, and searches them
//     with BFS, DFS or per-level selection.
//
// These types are not thread-safe on their own; the only concurrency is
// internal to SelectOneChildByEachValue. Concurrent mutation requires external
// synchronization.
package tree
