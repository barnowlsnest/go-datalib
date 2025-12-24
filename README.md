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
import "github.com/barnowlsnest/go-datalib/pkg/stack"
import "github.com/barnowlsnest/go-datalib/pkg/node"

s := stack.New()
s.Push(node.New(1, nil, nil))
s.Push(node.New(2, nil, nil))

top, ok := s.Peek()  // View top without removing
n := s.Pop()         // Remove and return top
```

### Queue

```go
import "github.com/barnowlsnest/go-datalib/pkg/queue"
import "github.com/barnowlsnest/go-datalib/pkg/node"

q := queue.New()
q.Enqueue(node.New(1, nil, nil))
q.Enqueue(node.New(2, nil, nil))

front, ok := q.PeekFront()  // View front
n := q.Dequeue()            // Remove from front
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

// Hash two uint64 values
nodeA := uint64(123)
nodeB := uint64(456)
edgeHash := serial.NSum(nodeA, nodeB)

// Order matters: NSum(a, b) != NSum(b, a)
```

### Binary Search Tree (BST)

```go
import (
	"github.com/barnowlsnest/go-datalib/pkg/tree"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

// Create a new BST
bst := tree.NewBST[int]()

// Insert values
bst.Insert(node.New(1, nil, nil), 50)
bst.Insert(node.New(2, nil, nil), 30)
bst.Insert(node.New(3, nil, nil), 70)
bst.Insert(node.New(4, nil, nil), 20)
bst.Insert(node.New(5, nil, nil), 40)

// Search for a value
found, exists := bst.Search(30)  // Returns node and true

// Delete a value
deleted := bst.Delete(30)  // Returns true if deleted

// Traverse the tree (in-order, pre-order, post-order, level-order)
nodes := bst.InOrder()  // Returns [20, 30, 40, 50, 70]
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

### Directed Acyclic Graph (DAG)

```go
import (
	"github.com/barnowlsnest/go-datalib/pkg/dag"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

// Create a new DAG
g := dag.New()

// Add nodes to groups
task1 := dag.GroupNode{Group: "build", ID: 1}
task2 := dag.GroupNode{Group: "build", ID: 2}
task3 := dag.GroupNode{Group: "test", ID: 3}

g.AddNode(task1)
g.AddNode(task2)
g.AddNode(task3)

// Add edges (dependencies)
edgeID := serial.NSum(task1.ID, task2.ID)
g.AddEdge(task1.ID, task2.ID, edgeID)  // task1 -> task2
g.AddEdge(task2.ID, task3.ID, serial.NSum(task2.ID, task3.ID))  // task2 -> task3

// Check if graph is acyclic (detect cycles)
isDAG, err := g.IsAcyclic()
if !isDAG {
	// Handle cycle
}

// Get topological sort order
order, err := g.TopologicalSort()  // Returns valid execution order

// Query relationships
hasEdge := g.HasEdge(task1.ID, task2.ID)  // true
predecessors := g.GetBackRefs(task3.ID)   // Returns [task2]
```

### Multi-way Tree (MTree)

```go
import (
	mtree "github.com/barnowlsnest/go-datalib/pkg/tree"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

// Create nodes with maximum breadth (max children)
root, _ := mtree.NewNode[string](1, 5, mtree.ValueOpt("CEO"))
root.asRoot() // Mark as root node

engineering, _ := mtree.NewNode[string](2, 3,
	mtree.ValueOpt("Engineering"),
	mtree.ParentOpt(root),
)

sales, _ := mtree.NewNode[string](3, 3,
	mtree.ValueOpt("Sales"),
	mtree.ParentOpt(root),
)

// Query operations
children, _ := root.SelectChildrenFunc(func(n *mtree.Node[string]) bool {
	return n.Val() == "Engineering"
})

// Build from model with cycle detection
model := mtree.HierarchyModel{
	mtree.RootTag: {"Company"},
	"Company":     {"Engineering", "Sales"},
	"Engineering": {"Frontend", "Backend"},
}
idGen := func() uint64 { return serial.Seq().Next("company") }
tree, err := mtree.Hierarchy(model, 10, idGen)
if err != nil {
	// Handles cycles: A→B→A will return error
	panic(err)
}

// Convert back to model
modelCopy, _ := mtree.ToModel(tree)

// Use Container for bounded trees with level tracking
root2, _ := mtree.NewNode[string](serial.Seq().Next("org"), 5, mtree.ValueOpt("root"))
container, _ := mtree.NewContainer(root2, 5, 10)  // Max 5 children, 10 levels deep

child1, _ := mtree.NewNode[string](serial.Seq().Next("org"), 5, mtree.ValueOpt("child1"))
container.Insert(child1, false)

// Iterate nodes at specific level
for node := range container.NodesIter(1) {
	fmt.Println(node.Val())
}

// Check container status
fmt.Printf("Size: %d, Depth: %d, Capacity: %d\n",
	container.Size(), container.Depth(), container.Capacity())
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
| Stack          | O(1)                     | O(1)                     |  O(1) peek               | O(n)   |
| Queue          | O(1)                     | O(1)                     | O(1) peek                | O(n)   |
| BST            | O(log n) avg, O(n) worst | O(log n) avg, O(n) worst | O(log n) avg, O(n) worst | O(n)   |
| Heap           | O(log n)                 | O(log n)                 | O(1) peek                | O(n)   |
| Fenwick Tree   | O(log n)                 | N/A                      | O(log n)                 | O(n)   |
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
