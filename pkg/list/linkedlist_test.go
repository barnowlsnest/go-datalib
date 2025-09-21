package list

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/barnowlsnest/go-datalib/pkg/node"
)

func TestNewLinkedList(t *testing.T) {
	t.Run("should create empty linked list", func(t *testing.T) {
		list := New()

		assert.NotNil(t, list)
		assert.Nil(t, list.head)
		assert.Nil(t, list.tail)
		assert.Equal(t, 0, list.size)
	})
}

func TestPush(t *testing.T) {
	t.Run("should push to empty list", func(t *testing.T) {
		list := New()
		n := node.New(1, nil, nil)

		list.Push(n)

		assert.Equal(t, 1, list.size)
		assert.Equal(t, n, list.head)
		assert.Equal(t, n, list.tail)
		assert.Nil(t, n.Next())
		assert.Nil(t, n.Prev())
	})

	t.Run("should push to non-empty list", func(t *testing.T) {
		list := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)

		list.Push(node1)
		list.Push(node2)

		assert.Equal(t, 2, list.size)
		assert.Equal(t, node1, list.head)
		assert.Equal(t, node2, list.tail)
		assert.Equal(t, node2, node1.Next())
		assert.Equal(t, node1, node2.Prev())
		assert.Nil(t, node2.Next())
	})
}

func TestPop(t *testing.T) {
	t.Run("should return nil when popping from empty list", func(t *testing.T) {
		list := New()

		result := list.Pop()

		assert.Nil(t, result)
		assert.Equal(t, 0, list.size)
	})

	t.Run("should pop from list with one element", func(t *testing.T) {
		list := New()
		n := node.New(1, nil, nil)
		list.Push(n)

		result := list.Pop()

		assert.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 0, list.size)
		assert.Nil(t, list.head)
		assert.Nil(t, list.tail)
	})

	t.Run("should pop from list with multiple elements", func(t *testing.T) {
		list := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		list.Push(node1)
		list.Push(node2)

		result := list.Pop()

		assert.NotNil(t, result)
		assert.Equal(t, uint64(2), result.ID())
		assert.Equal(t, 1, list.size)
		assert.Equal(t, node1, list.head)
		assert.Equal(t, node1, list.tail)
		assert.Nil(t, node1.Next())
	})
}

func TestUnshift(t *testing.T) {
	t.Run("should unshift to empty list", func(t *testing.T) {
		list := New()
		n := node.New(1, nil, nil)

		list.Unshift(n)

		assert.Equal(t, 1, list.size)
		assert.Equal(t, n, list.head)
		assert.Equal(t, n, list.tail)
		assert.Nil(t, n.Next())
		assert.Nil(t, n.Prev())
	})

	t.Run("should unshift to non-empty list", func(t *testing.T) {
		list := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)

		list.Push(node1)
		list.Unshift(node2)

		assert.Equal(t, 2, list.size)
		assert.Equal(t, node2, list.head)
		assert.Equal(t, node1, list.tail)
		assert.Equal(t, node1, node2.Next())
		assert.Equal(t, node2, node1.Prev())
		assert.Nil(t, node1.Next())
	})
}

func TestShift(t *testing.T) {
	t.Run("should return nil when shifting from empty list", func(t *testing.T) {
		list := New()

		result := list.Shift()

		assert.Nil(t, result)
		assert.Equal(t, 0, list.size)
	})

	t.Run("should shift from list with one element", func(t *testing.T) {
		list := New()
		n := node.New(1, nil, nil)
		list.Push(n)

		result := list.Shift()

		assert.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 0, list.size)
		assert.Nil(t, list.head)
		assert.Nil(t, list.tail)
	})

	t.Run("should shift from list with multiple elements", func(t *testing.T) {
		list := New()
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		list.Push(node1)
		list.Push(node2)

		result := list.Shift()

		assert.NotNil(t, result)
		assert.Equal(t, uint64(1), result.ID())
		assert.Equal(t, 1, list.size)
		assert.Equal(t, node2, list.head)
		assert.Equal(t, node2, list.tail)
		assert.Nil(t, node2.Prev())
	})
}

func TestCombinedOperations(t *testing.T) {
	t.Run("should handle push, pop, unshift, shift in sequence", func(t *testing.T) {
		list := New()

		// Push two nodes
		node1 := node.New(1, nil, nil)
		node2 := node.New(2, nil, nil)
		list.Push(node1)
		list.Push(node2)

		assert.Equal(t, 2, list.size)

		// Pop one node
		popped := list.Pop()
		assert.Equal(t, uint64(2), popped.ID())
		assert.Equal(t, 1, list.size)

		// Unshift a new node
		node3 := node.New(3, nil, nil)
		list.Unshift(node3)
		assert.Equal(t, 2, list.size)
		assert.Equal(t, node3, list.head)
		assert.Equal(t, node1, list.tail)

		// Shift a node
		shifted := list.Shift()
		assert.Equal(t, uint64(3), shifted.ID())
		assert.Equal(t, 1, list.size)
		assert.Equal(t, node1, list.head)
		assert.Equal(t, node1, list.tail)
	})
}
