package list

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/barnowlsnest/go-datalib/pkg/node"
)

func TestNewStack(t *testing.T) {
	t.Run("should create empty stack", func(t *testing.T) {
		s := NewStack()

		assert.NotNil(t, s)
		assert.Equal(t, 0, s.Size())
		assert.True(t, s.IsEmpty())
	})
}

func TestPush(t *testing.T) {
	t.Run("should push to empty stack", func(t *testing.T) {
		s := NewStack()
		n := node.New(1, nil, nil)

		s.Push(n)

		assert.Equal(t, 1, s.Size())
		assert.False(t, s.IsEmpty())
	})

	t.Run("should push multiple elements", func(t *testing.T) {
		s := NewStack()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		s.Push(node1)
		s.Push(node2)
		s.Push(node3)

		assert.Equal(t, 3, s.Size())
		assert.False(t, s.IsEmpty())
	})
}

func TestPop(t *testing.T) {
	t.Run("should return nil when popping from empty stack", func(t *testing.T) {
		s := NewStack()

		result := s.Pop()

		assert.Nil(t, result)
		assert.Equal(t, 0, s.Size())
		assert.True(t, s.IsEmpty())
	})

	t.Run("should pop from stack with one element", func(t *testing.T) {
		s := NewStack()
		n := node.New(1, nil, nil)
		s.Push(n)

		result := s.Pop()

		assert.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 0, s.Size())
		assert.True(t, s.IsEmpty())
	})

	t.Run("should pop in LIFO order", func(t *testing.T) {
		s := NewStack()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		s.Push(node1)
		s.Push(node2)
		s.Push(node3)

		// Pop should return in reverse order: 3, 2, 1
		result1 := s.Pop()
		assert.Equal(t, uint64(3), result1.ID())
		assert.Equal(t, 2, s.Size())

		result2 := s.Pop()
		assert.Equal(t, uint64(2), result2.ID())
		assert.Equal(t, 1, s.Size())

		result3 := s.Pop()
		assert.Equal(t, uint64(1), result3.ID())
		assert.Equal(t, 0, s.Size())
		assert.True(t, s.IsEmpty())
	})
}

func TestPeek(t *testing.T) {
	t.Run("should return empty node when peeking empty stack", func(t *testing.T) {
		s := NewStack()

		result, ok := s.Peek()

		assert.False(t, ok)
		assert.Equal(t, node.Node{}, result)
		assert.Equal(t, 0, s.Size())
	})

	t.Run("should peek top element without removing it", func(t *testing.T) {
		s := NewStack()
		n := node.New(1, nil, nil)
		s.Push(n)

		result, ok := s.Peek()

		assert.True(t, ok)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 1, s.Size()) // Size unchanged
	})

	t.Run("should peek correct element after multiple pushes", func(t *testing.T) {
		s := NewStack()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		s.Push(node1)
		s.Push(node2)
		s.Push(node3)

		result, ok := s.Peek()

		assert.True(t, ok)
		assert.Equal(t, uint64(3), result.ID())
		assert.Equal(t, 3, s.Size()) // Size unchanged
	})

	t.Run("should peek correct element after pop", func(t *testing.T) {
		s := NewStack()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)

		s.Push(node1)
		s.Push(node2)
		s.Pop()

		result, ok := s.Peek()

		assert.True(t, ok)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 1, s.Size())
	})
}

func TestStack_IsEmpty(t *testing.T) {
	t.Run("should return true for new stack", func(t *testing.T) {
		s := NewStack()

		assert.True(t, s.IsEmpty())
	})

	t.Run("should return false after push", func(t *testing.T) {
		s := NewStack()
		s.Push(node.New(1, nil, nil))

		assert.False(t, s.IsEmpty())
	})

	t.Run("should return true after push and pop", func(t *testing.T) {
		s := NewStack()
		s.Push(node.New(1, nil, nil))
		s.Pop()

		assert.True(t, s.IsEmpty())
	})
}

func TestStackCombinedOperations(t *testing.T) {
	t.Run("should handle push, pop, peek in sequence", func(t *testing.T) {
		s := NewStack()

		// Push elements
		s.Push(node.New(1, nil, nil))
		s.Push(node.New(2, nil, nil))
		assert.Equal(t, 2, s.Size())

		// Peek
		top, ok := s.Peek()
		assert.True(t, ok)
		assert.Equal(t, uint64(2), top.ID())
		assert.Equal(t, 2, s.Size())

		// Pop one
		popped := s.Pop()
		assert.Equal(t, uint64(2), popped.ID())
		assert.Equal(t, 1, s.Size())

		// Push another
		s.Push(node.New(3, nil, nil))
		assert.Equal(t, 2, s.Size())

		// Peek should show new top
		top, ok = s.Peek()
		assert.True(t, ok)
		assert.Equal(t, uint64(3), top.ID())

		// Pop remaining
		s.Pop()
		s.Pop()
		assert.True(t, s.IsEmpty())
	})

	t.Run("should handle multiple push and pop cycles", func(t *testing.T) {
		s := NewStack()

		// First cycle
		s.Push(node.New(1, nil, nil))
		s.Push(node.New(2, nil, nil))
		s.Pop()
		s.Pop()
		assert.True(t, s.IsEmpty())

		// Second cycle
		s.Push(node.New(3, nil, nil))
		s.Push(node.New(4, nil, nil))
		result := s.Pop()
		assert.Equal(t, uint64(4), result.ID())
		assert.Equal(t, 1, s.Size())
	})
}
