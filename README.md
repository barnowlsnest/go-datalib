# go-datalib

A Go library providing fundamental data structures and high-performance utilities for building efficient applications.

## Features

### Data Structures

- **LinkedList** - Doubly-linked list with O(1) operations at both ends
- **Stack** - LIFO data structure built on LinkedList
- **Queue** - FIFO data structure built on LinkedList
- **Node** - Foundation for building custom linked data structures
- **MTree (Multi-way Tree)** - Generic tree structure with configurable max breadth, parent-child operations, and cycle detection

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

### Multi-way Tree (MTree)

```go
import (
	"github.com/barnowlsnest/go-datalib/pkg/mtree"
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
- **Data structures** (LinkedList, Stack, Queue, MTree): Require external synchronization for concurrent access
- **MTree.SelectOneChildByEachValue**: Context-aware concurrent child selection with proper goroutine synchronization

### Performance Optimizations

- **Cache-line alignment**: Serial package uses 64-byte alignment to prevent false sharing
- **Sharding**: 64 shards distribute atomic operations for reduced contention
- **Zero allocations**: All operations avoid heap allocations where possible

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem ./pkg/serial
```

## Development

This project uses [Task](https://taskfile.dev) for build automation:

```bash
task go-test       # Run tests
task go-bench      # Run benchmarks
task go-lint       # Run linter
task sanity        # Format, vet, lint, test
task build         # Full build pipeline
```

## License

See [LICENSE](LICENSE) file for details.