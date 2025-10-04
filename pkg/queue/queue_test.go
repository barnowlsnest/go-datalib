package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/barnowlsnest/go-datalib/pkg/node"
)

func TestNewQueue(t *testing.T) {
	t.Run("should create empty queue", func(t *testing.T) {
		q := New()

		assert.NotNil(t, q)
		assert.Equal(t, 0, q.Size())
		assert.True(t, q.IsEmpty())
	})
}

func TestEnqueue(t *testing.T) {
	t.Run("should enqueue to empty queue", func(t *testing.T) {
		q := New()
		n := node.New(1, nil, nil)

		q.Enqueue(n)

		assert.Equal(t, 1, q.Size())
		assert.False(t, q.IsEmpty())
	})

	t.Run("should enqueue multiple elements", func(t *testing.T) {
		q := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		q.Enqueue(node1)
		q.Enqueue(node2)
		q.Enqueue(node3)

		assert.Equal(t, 3, q.Size())
		assert.False(t, q.IsEmpty())
	})
}

func TestDequeue(t *testing.T) {
	t.Run("should return nil when dequeuing from empty queue", func(t *testing.T) {
		q := New()

		result := q.Dequeue()

		assert.Nil(t, result)
		assert.Equal(t, 0, q.Size())
		assert.True(t, q.IsEmpty())
	})

	t.Run("should dequeue from queue with one element", func(t *testing.T) {
		q := New()
		n := node.New(1, nil, nil)
		q.Enqueue(n)

		result := q.Dequeue()

		assert.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 0, q.Size())
		assert.True(t, q.IsEmpty())
	})

	t.Run("should dequeue in FIFO order", func(t *testing.T) {
		q := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		q.Enqueue(node1)
		q.Enqueue(node2)
		q.Enqueue(node3)

		// Dequeue should return in same order: 1, 2, 3
		result1 := q.Dequeue()
		assert.Equal(t, uint64(1), result1.ID())
		assert.Equal(t, 2, q.Size())

		result2 := q.Dequeue()
		assert.Equal(t, uint64(2), result2.ID())
		assert.Equal(t, 1, q.Size())

		result3 := q.Dequeue()
		assert.Equal(t, uint64(3), result3.ID())
		assert.Equal(t, 0, q.Size())
		assert.True(t, q.IsEmpty())
	})
}

func TestPeekFront(t *testing.T) {
	t.Run("should return empty node when peeking empty queue", func(t *testing.T) {
		q := New()

		result, ok := q.PeekFront()

		assert.False(t, ok)
		assert.Equal(t, node.Node{}, result)
		assert.Equal(t, 0, q.Size())
	})

	t.Run("should peek front element without removing it", func(t *testing.T) {
		q := New()
		n := node.New(1, nil, nil)
		q.Enqueue(n)

		result, ok := q.PeekFront()

		assert.True(t, ok)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 1, q.Size()) // Size unchanged
	})

	t.Run("should peek correct front element after multiple enqueues", func(t *testing.T) {
		q := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		q.Enqueue(node1)
		q.Enqueue(node2)
		q.Enqueue(node3)

		result, ok := q.PeekFront()

		assert.True(t, ok)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 3, q.Size()) // Size unchanged
	})

	t.Run("should peek correct front element after dequeue", func(t *testing.T) {
		q := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)

		q.Enqueue(node1)
		q.Enqueue(node2)
		q.Dequeue()

		result, ok := q.PeekFront()

		assert.True(t, ok)
		assert.Equal(t, uint64(2), result.ID())
		assert.Equal(t, 1, q.Size())
	})
}

func TestPeekRear(t *testing.T) {
	t.Run("should return empty node when peeking empty queue", func(t *testing.T) {
		q := New()

		result, ok := q.PeekRear()

		assert.False(t, ok)
		assert.Equal(t, node.Node{}, result)
		assert.Equal(t, 0, q.Size())
	})

	t.Run("should peek rear element without removing it", func(t *testing.T) {
		q := New()
		n := node.New(1, nil, nil)
		q.Enqueue(n)

		result, ok := q.PeekRear()

		assert.True(t, ok)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 1, q.Size()) // Size unchanged
	})

	t.Run("should peek correct rear element after multiple enqueues", func(t *testing.T) {
		q := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		node3 := node.New(3, nil, nil)

		q.Enqueue(node1)
		q.Enqueue(node2)
		q.Enqueue(node3)

		result, ok := q.PeekRear()

		assert.True(t, ok)
		assert.Equal(t, uint64(3), result.ID())
		assert.Equal(t, 3, q.Size()) // Size unchanged
	})

	t.Run("should peek correct rear element after enqueue", func(t *testing.T) {
		q := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)

		q.Enqueue(node1)
		q.Enqueue(node2)

		result, ok := q.PeekRear()

		assert.True(t, ok)
		assert.Equal(t, uint64(2), result.ID())
		assert.Equal(t, 2, q.Size())
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("should return true for new queue", func(t *testing.T) {
		q := New()

		assert.True(t, q.IsEmpty())
	})

	t.Run("should return false after enqueue", func(t *testing.T) {
		q := New()
		q.Enqueue(node.New(1, nil, nil))

		assert.False(t, q.IsEmpty())
	})

	t.Run("should return true after enqueue and dequeue", func(t *testing.T) {
		q := New()
		q.Enqueue(node.New(1, nil, nil))
		q.Dequeue()

		assert.True(t, q.IsEmpty())
	})
}

func TestQueueCombinedOperations(t *testing.T) {
	t.Run("should handle enqueue, dequeue, peek in sequence", func(t *testing.T) {
		q := New()

		// Enqueue elements
		q.Enqueue(node.New(1, nil, nil))
		q.Enqueue(node.New(2, nil, nil))
		assert.Equal(t, 2, q.Size())

		// Peek front and rear
		front, ok := q.PeekFront()
		assert.True(t, ok)
		assert.Equal(t, uint64(1), front.ID())

		rear, ok := q.PeekRear()
		assert.True(t, ok)
		assert.Equal(t, uint64(2), rear.ID())
		assert.Equal(t, 2, q.Size())

		// Dequeue one
		dequeued := q.Dequeue()
		assert.Equal(t, uint64(1), dequeued.ID())
		assert.Equal(t, 1, q.Size())

		// Enqueue another
		q.Enqueue(node.New(3, nil, nil))
		assert.Equal(t, 2, q.Size())

		// Peek should show new state
		front, ok = q.PeekFront()
		assert.True(t, ok)
		assert.Equal(t, uint64(2), front.ID())

		rear, ok = q.PeekRear()
		assert.True(t, ok)
		assert.Equal(t, uint64(3), rear.ID())

		// Dequeue remaining
		q.Dequeue()
		q.Dequeue()
		assert.True(t, q.IsEmpty())
	})

	t.Run("should handle multiple enqueue and dequeue cycles", func(t *testing.T) {
		q := New()

		// First cycle
		q.Enqueue(node.New(1, nil, nil))
		q.Enqueue(node.New(2, nil, nil))
		q.Dequeue()
		q.Dequeue()
		assert.True(t, q.IsEmpty())

		// Second cycle
		q.Enqueue(node.New(3, nil, nil))
		q.Enqueue(node.New(4, nil, nil))
		result := q.Dequeue()
		assert.Equal(t, uint64(3), result.ID())
		assert.Equal(t, 1, q.Size())
	})

	t.Run("should maintain FIFO with interleaved operations", func(t *testing.T) {
		q := New()

		q.Enqueue(node.New(1, nil, nil))
		q.Enqueue(node.New(2, nil, nil))

		result1 := q.Dequeue()
		assert.Equal(t, uint64(1), result1.ID())

		q.Enqueue(node.New(3, nil, nil))

		result2 := q.Dequeue()
		assert.Equal(t, uint64(2), result2.ID())

		result3 := q.Dequeue()
		assert.Equal(t, uint64(3), result3.ID())

		assert.True(t, q.IsEmpty())
	})
}