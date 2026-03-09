# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go data structures library (`go-datalib`) implementing fundamental data structures with a focus on doubly-linked lists. The project follows standard Go module conventions and provides clean APIs for common data structure operations.

## Commands

### Testing
- Run all tests: `go test -v ./...`
- Test specific package: `go test -v ./pkg/node` or `go test -v ./pkg/list`
- Run tests with coverage: `go test -cover ./...`

### Building
- Build the module: `go build ./...`
- Check module dependencies: `go mod tidy`

### Development
- Format code: `go fmt ./...`
- Vet code: `go vet ./...`
- Clean module cache: `go clean -modcache`

## Architecture

### Core Components

The library is structured with a layered architecture:

1. **Node Package** (`pkg/node/`): Provides the fundamental `Node` type for doubly-linked list structures
   - Immutable ID field with getter access
   - Mutable next/prev pointers via `WithNext()` and `WithPrev()` methods
   - Thread-safety must be handled by containing structures

2. **List Package** (`pkg/list/`): Implements `LinkedList` using the Node package
   - O(1) operations at both ends: `Push()`, `Pop()`, `Unshift()`, `Shift()`
   - Automatic size tracking
   - Memory-safe node copying during removal operations
   - Returns value copies through `Head()` and `Tail()` methods

### Design Patterns

- **Clean separation of concerns**: Node handles structure, LinkedList handles operations
- **Memory safety**: Node copies are returned during Pop/Shift to prevent memory leaks
- **Defensive programming**: Nil checks and safe empty list handling throughout
- **Encapsulation**: Private fields with controlled public access via methods

### Testing Strategy

The project uses comprehensive test suites with testify:
- **Node package**: Multiple test suites covering basic functionality, chaining, robustness, structural integrity, and mutable operations
- **List package**: Functional tests for all CRUD operations and combined operation scenarios
- **Coverage**: Both packages have thorough test coverage including edge cases

## Development Notes

- All struct fields are private; use provided getter/setter methods
- Node IDs are uint64 and immutable after creation
- LinkedList is not thread-safe; external synchronization required for concurrent access
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