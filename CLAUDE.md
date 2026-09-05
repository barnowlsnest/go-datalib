# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go data structures library (`go-datalib` v5) providing fundamental data structures including linked lists, stacks, queues, trees (BST, B-tree, heap, Fenwick, segment), directed acyclic graphs, an LRU cache, and a thread-safe serial ID generator. The module uses modern Go features including generics and range-over-func iterators.

**Module:** `github.com/barnowlsnest/go-datalib/v5` (Go 1.26.5, see `go.mod`)

## Commands

### Testing
- Run all tests: `go test -v ./...`
- Test specific package: `go test -v ./pkg/node`, `go test -v ./pkg/list`, `go test -v ./pkg/tree`, etc.
- Run tests with coverage: `go test -cover ./...`
- Run benchmarks: `go test -bench=. -benchmem ./...`
- Run serial benchmarks: `go test -bench=. -benchmem ./pkg/serial`

### Building
- Build the module: `go build ./...`
- Check module dependencies: `go mod tidy`

### Development
- Format code: `go fmt ./...`
- Vet code: `go vet ./...`
- Lint code: `golangci-lint run --fix`
- Clean module cache: `go clean -modcache`

## Architecture

### Core Components

The library is structured with a layered architecture across six packages:

1. **Node Package** (`pkg/node/`): Provides the fundamental `Node` type for doubly-linked list structures
   - Immutable ID field (uint64) with getter access
   - Mutable next/prev pointers via `WithNext()` and `WithPrev()` methods
   - Iterator support: `ForwardIterator`, `BackwardIterator`, and `Iterable` interface
   - Range-over-func support via `NextNodes()` and `PrevNodes()` returning `iter.Seq2`
   - Custom errors: `ErrEOI` (end of iteration), `ErrNil`

2. **List Package** (`pkg/list/`): Implements `LinkedList`, `Stack`, and `Queue`
   - **LinkedList**: O(1) operations at both ends — `Push()`, `Pop()`, `Unshift()`, `Shift()`
     - ID convenience methods: `PushID()`, `PopID()`, `UnshiftID()`, `ShiftID()`
     - Iteration: `IterNext()`, `IterPrev()` for range-based traversal
     - Automatic size tracking; memory-safe node copying during removal
   - **Stack**: LIFO wrapper around LinkedList — `Push()`, `Pop()`, `Peek()`, `IsEmpty()`
   - **Queue**: FIFO wrapper around LinkedList — `Enqueue()`, `Dequeue()`, `PeekFront()`, `PeekRear()`, `IsEmpty()`

3. **Tree Package** (`pkg/tree/`): Multiple tree data structures using generics
   - **Node[T comparable]**: Generic multi-way tree node with children management
     - `AttachChild()`, `AttachMany()`, `DetachChild()`, `MoveChildren()`, `Swap()`
     - Selection: `SelectChildrenFunc()`, `SelectOneChildFunc()`, `SelectOneChildByEachValue()` (concurrent with context)
     - Options pattern: `ValueOpt()`, `LevelOpt()`, `ParentOpt()`, `ChildOpt()`
   - **BinaryNode[T cmp.Ordered]**: Binary tree node embedding `*node.Node`
   - **BST[T cmp.Ordered]**: Binary Search Tree with `Insert()`, `Delete()`, `Search()`, traversals (InOrder, PreOrder, PostOrder, LevelOrder)
   - **Heap[T any]**: Generic binary heap with custom comparators — `NewMin()`, `NewMax()`, `Push()`, `Pop()`, `Peek()`
   - **BTree[K cmp.Ordered, V any]**: B-tree with configurable min degree — `Insert()`, `Search()`, `Delete()`, range iteration
   - **Fenwick[T]**: Binary Indexed Tree for prefix sums — `Update()`, `Query()`, `RangeQuery()`
   - **Segment[T comparable]**: Tree segment with capacity/depth constraints, aliases, level mapping
   - **Hierarchy**: Build trees from hierarchy models with cycle detection (uses `RootTag = "#root"`)

4. **DAG Package** (`pkg/dag/`): Directed Acyclic Graph with group support
   - `Graph` type with group-based node organization
   - Types: `GroupNode`, `AdjacencyEdge`, `BackRefEdge`
   - Uses `uuid.UUID` for IDs; type aliases for `NodeID`, `EdgeID`, `GroupName`
   - Cycle detection and validation

5. **LRU Package** (`pkg/lru/`): Generic, thread-safe Least-Recently-Used cache
   - `LRU[T any]` keyed by `uint64`, backed by `pkg/list.LinkedList` (recency order) plus a map for O(1) lookup
   - `New(capacity int)` falls back to `defaultLRUCapacity` (1024) for non-positive capacity
   - `Put()`, `Get()` (returns `ErrCacheMiss` on miss), eviction of the least-recently-used entry at capacity
   - Synchronized via an internal `sync.Mutex`

6. **Serial Package** (`pkg/serial/`): Thread-safe sharded ID generator
   - `Serial` type with 64 cache-line-aligned shards to prevent false sharing
   - `Next(key string) uint64` and `Current(key string) uint64`
   - FNV-1a hashing for key distribution
   - `NSum(from, to uint64) uint64` for order-independent hash combining

### Key Dependencies

- `github.com/google/uuid` — UUID generation (DAG package)
- `github.com/stretchr/testify` — Testing framework
- `golang.org/x/exp` — Iterators, constraints
- `golang.org/x/sync` — Errgroup for concurrent operations

### Design Patterns

- **Generics**: Extensive use throughout tree and heap implementations
- **Options pattern**: Tree node creation via functional options
- **Clean separation of concerns**: Node handles structure, higher-level types handle operations
- **Memory safety**: Node copies are returned during Pop/Shift to prevent memory leaks
- **Defensive programming**: Nil checks and safe empty collection handling throughout
- **Encapsulation**: Private fields with controlled public access via methods
- **Iterator compatibility**: Supports Go 1.22+ range-over-func via `iter.Seq2`
- **Concurrency**: Sharded atomics in Serial; context-aware concurrent selection in Tree

### Testing Strategy

The project uses comprehensive test suites with testify:
- **Node package**: Multiple test suites covering basic functionality, chaining, robustness, structural integrity, mutable operations, and iterators
- **List package**: Functional tests for LinkedList, Stack, and Queue including edge cases
- **Tree package**: Tests for BST, B-tree, heap, Fenwick tree, segment, hierarchy, and generic nodes
- **DAG package**: Graph construction, edge validation, and error handling tests
- **LRU package**: Capacity defaulting, hit/miss behavior, eviction order, and concurrency
- **Serial package**: Concurrency tests and benchmarks for sharded ID generation
- **Coverage**: All packages maintain thorough test coverage; 80% minimum enforced

## Development Notes

- All struct fields are private; use provided getter/setter methods
- Node IDs are uint64 and immutable after creation
- LinkedList, Stack, and Queue are not thread-safe; external synchronization required for concurrent access
- Serial package is thread-safe by design (sharded atomics); LRU is thread-safe via an internal mutex
- Test files use testify/assert and testify/suite for organized test structure

## Sanity Checks

**CRITICAL**: You MUST ALWAYS run the sanity task after ANY code change before considering the task complete:

```bash
task sanity
```

This single command runs all necessary checks in the correct order:
1. Updates dependencies (`go mod tidy`)
2. Formats code (`go fmt ./...`)
3. Vets code (`go vet ./...`)
4. Runs linter (`golangci-lint run --fix`)
5. Runs all tests (`go test -cover ./...`)
6. Verifies code coverage meets 80% threshold

**NEVER use individual go commands for sanity checks - ALWAYS use `task sanity`.**

These checks ensure code quality, prevent regressions, and maintain project standards. Never skip sanity checks.