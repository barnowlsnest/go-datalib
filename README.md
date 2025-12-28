# go-datalib

A Go library providing fundamental data structures and high-performance utilities for building efficient applications.

## Features

### Data Structures

#### Linear Data Structures
- **LinkedList** - Doubly-linked list with O(1) operations at both ends
- **Stack** - LIFO data structure built on LinkedList
- **Queue** - FIFO data structure built on LinkedList
- **Node** - Foundation for building custom linked data structures

#### Tree Structures
- **BST (Binary Search Tree)** - Iterative BST with O(log n) average-case operations, supports multiple traversal orders
- **Heap** - Generic binary heap (min/max) with O(log n) insert/delete and O(1) peek
- **Fenwick Tree (Binary Indexed Tree)** - Efficient prefix sums and point updates in O(log n) time
- **Segment Tree** - Generic segment tree for range queries with configurable depth/breadth, DFS/BFS traversal, and level-based node organization
- **B-Tree** - Self-balancing tree with O(log n) operations, range queries, and floor/ceiling lookups
- **MTree (Multi-way Tree)** - Generic M-way tree with configurable breadth/depth, hierarchy building, and cycle detection

#### Graph Structures
- **DAG (Directed Acyclic Graph)** - Directed graph with cycle detection via Kahn's algorithm, group-based node organization

### Utilities

- **Serial** - High-performance, thread-safe ID generator with sharding and cache-line alignment
- **NSum** - Fast hash function for combining uint64 pairs using the golden ratio

## Installation

```bash
go get github.com/barnowlsnest/go-datalib
```

## Quick Start

### LinkedList

```go
import "github.com/barnowlsnest/go-datalib/pkg/list"
import "github.com/barnowlsnest/go-datalib/pkg/node"

list := list.New()
list.Push(node.New(1, nil, nil))
list.Push(node.New(2, nil, nil))

n := list.Pop() // Returns node with ID 2
```

### Stack

```go
import (
	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

s := list.NewStack()
s.Push(node.New(1, nil, nil))
s.Push(node.New(2, nil, nil))

top, ok := s.Peek()  // View top without removing
n := s.Pop()         // Remove and return top
isEmpty := s.IsEmpty()
size := s.Size()
```

### Queue

```go
import (
	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

q := list.NewQueue()
q.Enqueue(node.New(1, nil, nil))
q.Enqueue(node.New(2, nil, nil))

front, okFront := q.PeekFront()  // View front without removing
rear, okRear := q.PeekRear()     // View rear without removing
n := q.Dequeue()                 // Remove from front
isEmpty := q.IsEmpty()
size := q.Size()
```

### Serial ID Generator

```go
import "github.com/barnowlsnest/go-datalib/pkg/serial"

// Create instance or use singleton
gen := &serial.Serial{}
// or
gen := serial.Seq()  // Global singleton

// Generate sequential IDs per key
id1 := gen.Next("user")     // Returns 1
id2 := gen.Next("user")     // Returns 2
id3 := gen.Next("product")  // Returns 1 (different key)

// Read current value
current := gen.Current("user")  // Returns 2
```

### NSum Hash Function

```go
import "github.com/barnowlsnest/go-datalib/pkg/serial"

// Hash two uint64 values into a single hash
nodeA := uint64(123)
nodeB := uint64(456)
edgeHash := serial.NSum(nodeA, nodeB)

// Order is normalized: NSum(a, b) == NSum(b, a)
// Useful for creating undirected edge identifiers
```

### Binary Search Tree (BST)

```go
import (
	"fmt"
	"github.com/barnowlsnest/go-datalib/pkg/tree"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

// Create a new BST
bst := tree.NewBST[int]()

// Insert values (node, value)
bst.Insert(node.New(1, nil, nil), 50)
bst.Insert(node.New(2, nil, nil), 30)
bst.Insert(node.New(3, nil, nil), 70)
bst.Insert(node.New(4, nil, nil), 20)
bst.Insert(node.New(5, nil, nil), 40)

// Search for a value
found := bst.Search(30)  // Returns *BinaryNode[int] or nil

// Delete a value
deleted := bst.Delete(30)  // Returns true if deleted

// Traverse the tree (in-order, pre-order, post-order, level-order)
bst.InOrder(func(n *tree.BinaryNode[int]) {
	fmt.Println(n.Value())  // Prints: 20, 30, 40, 50, 70
})

// Get min/max values
min := bst.Min()  // Returns node with value 20
max := bst.Max()  // Returns node with value 70

// Check tree properties
size := bst.Size()      // Number of nodes
height := bst.Height()  // Tree height
```

### Heap (Min/Max Binary Heap)

```go
import "github.com/barnowlsnest/go-datalib/pkg/tree"

// Create a min-heap (smallest element on top)
minHeap := tree.NewMin[int]()
minHeap.Push(50)
minHeap.Push(30)
minHeap.Push(70)
minHeap.Push(20)

top, _ := minHeap.Peek()  // Returns 20 (smallest)
min, _ := minHeap.Pop()   // Removes and returns 20

// Create a max-heap (largest element on top)
maxHeap := tree.NewMax[int]()
maxHeap.Push(50)
maxHeap.Push(30)
maxHeap.Push(70)

max, _ := maxHeap.Pop()  // Removes and returns 70 (largest)

// Build heap from existing slice (O(n) heapify)
data := []int{3, 2, 1, 5, 4}
heap := tree.HeapFromSlice(data, func(a, b int) bool { return a < b })
```

### Fenwick Tree (Binary Indexed Tree)

```go
import "github.com/barnowlsnest/go-datalib/pkg/tree"

// Create from slice
data := []int{3, 2, -1, 6, 5, 4, -3, 3, 7, 2, 3}
ft := tree.FromSlice(data)

// Query prefix sum (1-indexed)
sum := ft.Query(5)  // Sum of elements from index 1 to 5

// Update element (add delta)
ft.Update(3, 5)   // Add 5 to element at index 3
ft.Update(3, -2)  // Subtract 2 from element at index 3

// Range sum query
rangeSum := ft.RangeQuery(2, 7)  // Sum from index 2 to 7

// Create empty tree
ft2 := tree.NewFenwick[int](100)  // Size 100
```

### Segment Tree

```go
import (
	"fmt"
	"github.com/barnowlsnest/go-datalib/pkg/tree"
)

// Create a new segment tree with alias, ID, max breadth, and max depth
seg := tree.NewSegment[string]("users", 1, 10, 5)

// Create and insert root node
root, err := tree.NewNode[string](1, 10, tree.ValueOpt("root"))
if err != nil {
	panic(err)
}
seg.Insert(root, 0)  // parentID 0 for root

// Insert child nodes
child1, err := tree.NewNode[string](2, 10, tree.ValueOpt("child1"))
if err != nil {
	panic(err)
}
seg.Insert(child1, root.ID())

child2, err := tree.NewNode[string](3, 10, tree.ValueOpt("child2"))
if err != nil {
	panic(err)
}
seg.Insert(child2, root.ID())

// Query segment properties
height := seg.Height()              // Current tree height
length := seg.Length()              // Number of nodes
capacity := seg.Capacity()          // Max total nodes
remaining := seg.RemainingCapacity()

// Get nodes
rootNode, _ := seg.Root()
nodeByID, _ := seg.NodeByID(2)

// Traversal
seg.DFS(func(n *tree.Node[string]) bool {
	fmt.Println(n.Val())
	return true  // continue traversal
})

seg.BFS(func(n *tree.Node[string]) bool {
	fmt.Println(n.Val())
	return true
})

// Level-based operations
seg.ForEachNodeAtLevel(1, func(n *tree.Node[string]) bool {
	return true
})

// Selection
nodes := seg.Select(func(n *tree.Node[string]) bool {
	return n.Val() == "child1"
})

// Node management
seg.Link(root.ID(), child1.ID())        // Link existing nodes
seg.Unlink(root.ID(), child1.ID())      // Break relationship
seg.RemoveCascade(child1.ID())          // Remove node and descendants
seg.RemovePromote(child1.ID())          // Remove and promote children
```

### B-Tree

```go
import (
	"fmt"
	"github.com/barnowlsnest/go-datalib/pkg/tree"
)

// Create a B-tree with minimum degree 3
bt := tree.NewBTree[uint64, string](3)

// Insert key-value pairs
bt.Insert(100, "first")
bt.Insert(200, "second")
bt.Insert(50, "third")
bt.Insert(150, "fourth")

// Search for a value
value, found := bt.Search(100)  // Returns "first", true

// Check if key exists
exists := bt.Contains(100)  // true

// Delete a key
deleted := bt.Delete(100)  // true

// Get min/max entries
minKey, minVal, _ := bt.Min()  // 50, "third"
maxKey, maxVal, _ := bt.Max()  // 200, "second"

// Floor/Ceiling queries
floorKey, floorVal, _ := bt.Floor(125)    // Largest key <= 125
ceilKey, ceilVal, _ := bt.Ceiling(125)    // Smallest key >= 125

// Range iteration
for entry := range bt.Range(50, 200) {
	fmt.Printf("%d: %s\n", entry.Key, entry.Value)
}

// Iterate all entries in sorted order
for entry := range bt.All() {
	fmt.Printf("%d: %s\n", entry.Key, entry.Value)
}

// Get all keys/values
keys := bt.Keys()      // []uint64
values := bt.Values()  // []string

// Tree properties
size := bt.Size()
height := bt.Height()
isEmpty := bt.IsEmpty()
bt.Clear()
```

### Directed Acyclic Graph (DAG)

```go
import "github.com/barnowlsnest/go-datalib/pkg/dag"

// Create a new DAG
g := dag.New()

// Create groups first
g.AddGroup("build")
g.AddGroup("test")

// Define nodes with group membership
task1 := dag.GroupNode{ID: 1, Group: "build"}
task2 := dag.GroupNode{ID: 2, Group: "build"}
task3 := dag.GroupNode{ID: 3, Group: "test"}

// Add nodes to their groups
g.AddNode(task1)
g.AddNode(task2)
g.AddNode(task3)

// Add edges (dependencies) - edge IDs are auto-generated using NSum
g.AddEdge(task1, task2)  // task1 -> task2
g.AddEdge(task2, task3)  // task2 -> task3

// Check if graph is acyclic (async via channel)
isAcyclic := <-g.IsAcyclic()  // Returns true if no cycles

// Query relationships
hasEdge := g.HasEdge(task1, task2)           // true
predecessors, _ := g.GetBackRefsOf(task3)    // Returns [task2]

// Iterate over neighbors
g.ForEachNeighbour(task1, func(edge dag.AdjacencyEdge, err error) {
	// edge.From, edge.To, edge.Edge (ID)
})

// Get all nodes in a group
buildNodes, _ := g.GetNodes("build")  // Returns [task1, task2]

// List all groups
groups := g.ListGroups()  // Returns ["build", "test"]
```

### Multi-way Tree (MTree)

```go
import (
	"context"
	"github.com/barnowlsnest/go-datalib/pkg/tree"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

// Create nodes with maximum breadth (max children)
root, _ := tree.NewNode[string](1, 5, tree.ValueOpt("CEO"))

// Attach children using ParentOpt during creation
engineering, _ := tree.NewNode[string](2, 3,
	tree.ValueOpt("Engineering"),
	tree.ParentOpt(root),
)

sales, _ := tree.NewNode[string](3, 3,
	tree.ValueOpt("Sales"),
	tree.ParentOpt(root),
)

// Or attach children manually
backend, _ := tree.NewNode[string](4, 0, tree.ValueOpt("Backend"))
engineering.AttachChild(backend)

// Query operations
children, _ := root.SelectChildrenFunc(func(n *tree.Node[string]) bool {
	return n.Val() == "Engineering"
})

// Select single child by predicate
eng, _ := root.SelectOneChildFunc(func(n *tree.Node[string]) bool {
	return n.Val() == "Engineering"
})

// Concurrent child selection by multiple values
ctx := context.Background()
selected, _ := root.SelectOneChildByEachValue(ctx, "Engineering", "Sales")

// Build from model with cycle detection
model := tree.HierarchyModel{
	tree.RootTag: {"Company"},
	"Company":     {"Engineering", "Sales"},
	"Engineering": {"Frontend", "Backend"},
}
idGen := func() uint64 { return serial.Seq().Next("company") }
rootNode, err := tree.Hierarchy(model, 10, idGen)
if err != nil {
	// Handles cycles: A→B→A will return error
	panic(err)
}

// Convert back to model
modelCopy, _ := tree.ToModel(rootNode)

// Node operations
backend.Detach()                    // Detach from parent
backend.Move(sales)                 // Move to new parent
engineering.MoveChildren(sales)     // Move all children to new parent
root.DetachChild(sales)             // Detach specific child
root.DetachChildFunc(func(n *tree.Node[string]) bool {
	return n.Val() == "Sales"       // Detach children matching predicate
})
```

## Performance

### Serial Benchmarks (Apple M1 Max)

- **Next() same key**: 6.9 ns/op (0 allocs)
- **Next() parallel same key**: 74.5 ns/op
- **Next() parallel different keys**: 22.8 ns/op
- **Current() parallel**: 0.9 ns/op
- **NSum()**: ~1-2 ns/op

The sharded design with cache-line alignment significantly reduces contention in high-concurrency scenarios.

## Architecture

### Memory Safety

All data structures return copies of nodes during Pop/Shift/Dequeue operations with cleared references to prevent memory leaks.

### Thread Safety

- **Serial package**: Fully thread-safe using atomic operations
- **Linear data structures** (LinkedList, Stack, Queue): Require external synchronization for concurrent access
- **Tree structures** (BST, Heap, Fenwick, MTree): Require external synchronization for concurrent access
- **Graph structures** (DAG): Require external synchronization for concurrent access
- **MTree.SelectOneChildByEachValue**: Context-aware concurrent child selection with proper goroutine synchronization

### Complexity Summary

| Data Structure | Insert                   | Delete                   | Search/Query             | Space  |
|----------------|--------------------------|--------------------------|--------------------------|--------|
| LinkedList     | O(1)                     | O(1)                     | O(n)                     | O(n)   |
| Stack          | O(1)                     | O(1)                     | O(1) peek                | O(n)   |
| Queue          | O(1)                     | O(1)                     | O(1) peek                | O(n)   |
| BST            | O(log n) avg, O(n) worst | O(log n) avg, O(n) worst | O(log n) avg, O(n) worst | O(n)   |
| Heap           | O(log n)                 | O(log n)                 | O(1) peek                | O(n)   |
| Fenwick Tree   | O(log n)                 | N/A                      | O(log n)                 | O(n)   |
| Segment Tree   | O(1)                     | O(1)                     | O(n) traversal           | O(n)   |
| B-Tree         | O(log n)                 | O(log n)                 | O(log n)                 | O(n)   |
| MTree          | O(1) attach              | O(1) detach              | O(n) traversal           | O(n)   |
| DAG            | O(1)                     | O(1)                     | O(V+E) cycle detection   | O(V+E) |

### Performance Optimizations

- **Cache-line alignment**: Serial package uses 64-byte alignment to prevent false sharing
- **Sharding**: 64 shards distribute atomic operations for reduced contention
- **Zero allocations**: All operations avoid heap allocations where possible

## Testing

This project uses [Task](https://taskfile.dev) for build automation. All test, build, and benchmark commands should be run via Task:

```bash
# Run all tests with coverage
task go-test

# Check coverage with threshold validation
task go-coverage

# Run all benchmarks
task go-bench

# Run serial package benchmarks only
task go-bench-serial
```

## Development

### Available Tasks

```bash
# Code Quality
task go-fmt        # Format code
task go-vet        # Run go vet
task go-lint       # Run linter with auto-fix
task sanity        # Run all checks: format, vet, lint, test, coverage

# Build
task go-build      # Build all packages
task build         # Full build pipeline (sanity + build)

# Testing & Benchmarks
task go-test       # Run tests with coverage
task go-coverage   # Check coverage meets 80% threshold
task go-bench      # Run all benchmarks with memory stats
task go-bench-serial  # Run serial package benchmarks

# Maintenance
task go-update     # Update dependencies (go mod tidy)
```

### Quick Commands

```bash
# Most common workflow
task sanity        # Format, vet, lint, test, coverage check
task build         # Full build pipeline
```

## License

See [LICENSE](LICENSE) file for details.
