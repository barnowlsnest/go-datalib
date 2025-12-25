package tree

import (
	"fmt"
	"strings"

	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

const (
	DefaultMaxDepth   = 64
	DefaultMaxBreadth = 32
)

type (
	Segment[T comparable] struct {
		alias      string
		id         uint64
		maxDepth   int
		maxBreadth int
		cap        int
		root       *Node[T]
		levelMap   map[int][]uint64
		nodeMap    map[uint64]*Node[T]
	}

	Selector[T comparable] struct {
		Type  string
		ID    uint64
		Level int
		Limit int
	}

	VisitorFunc[T comparable] func(n *Node[T]) bool

	traverser interface {
		add(id uint64)
		next() (uint64, bool)
		isEmpty() bool
	}

	stackTraverser struct {
		stack *list.Stack
	}

	queueTraverser struct {
		queue *list.Queue
	}
)

func NewSegment[T comparable](alias string, id uint64, maxBreadth, maxDepth int) *Segment[T] {
	var (
		mAlias   string
		mDepth   int
		mBreadth int
	)
	if maxDepth <= 0 {
		mDepth = DefaultMaxDepth
	} else {
		mDepth = maxDepth
	}
	if maxBreadth <= 0 {
		mBreadth = DefaultMaxBreadth
	} else {
		mBreadth = maxBreadth
	}

	mAlias = strings.ReplaceAll(alias, " ", "")
	if mAlias == "" {
		mAlias = fmt.Sprintf("seg.%d", id)
	}

	return &Segment[T]{
		id:         id,
		alias:      mAlias,
		maxDepth:   mDepth,
		maxBreadth: mBreadth,
		cap:        mDepth * mBreadth,
		levelMap:   make(map[int][]uint64, mDepth),
		nodeMap:    make(map[uint64]*Node[T]),
	}
}

func (s *Segment[T]) Alias() string {
	return s.alias
}

func (s *Segment[T]) ID() uint64 {
	return s.id
}

func (s *Segment[T]) Capacity() int {
	return s.cap
}

func (s *Segment[T]) Height() int {
	return len(s.levelMap)
}

func (s *Segment[T]) Length() int {
	return len(s.nodeMap)
}

func (s *Segment[T]) RemainingCapacity() int {
	return s.cap - len(s.nodeMap)
}

func (s *Segment[T]) Root() (*Node[T], bool) {
	if s.root == nil {
		return nil, false
	}

	return s.root, true
}

func (s *Segment[T]) NodeByID(id uint64) (*Node[T], error) {
	n, exists := s.nodeMap[id]
	if !exists {
		return nil, ErrNodeNotFound
	}

	return n, nil
}

func (s *Segment[T]) nodesAtLevel(level int) ([]*Node[T], error) {
	nodes, existsLevel := s.levelMap[level]
	if !existsLevel {
		return nil, ErrSegmentLevelNotFound
	}
	if level == 0 {
		return []*Node[T]{s.root}, nil
	}

	parents := s.levelMap[level-1]
	levelSlice := make([]*Node[T], 0, len(parents)*s.maxBreadth)
	for _, n := range nodes {
		vNode, existsNode := s.nodeMap[n]
		if !existsNode {
			return nil, ErrSegmentDoesNotHaveNode
		}

		levelSlice = append(levelSlice, vNode)
	}

	return levelSlice, nil
}

func (t *stackTraverser) add(id uint64) {
	t.stack.Push(node.ID(id))
}

func (t *stackTraverser) next() (uint64, bool) {
	n := t.stack.Pop()
	if n == nil {
		return 0, false
	}
	return n.ID(), true
}

func (t *stackTraverser) isEmpty() bool {
	return t.stack.IsEmpty()
}

func (t *queueTraverser) add(id uint64) {
	t.queue.Enqueue(node.ID(id))
}

func (t *queueTraverser) next() (uint64, bool) {
	n := t.queue.Dequeue()
	if n == nil {
		return 0, false
	}
	return n.ID(), true
}

func (t *queueTraverser) isEmpty() bool {
	return t.queue.IsEmpty()
}

func (s *Segment[T]) traverse(t traverser, visitor VisitorFunc[T]) error {
	if s.root == nil {
		return nil
	}

	t.add(s.root.ID())

	for !t.isEmpty() {
		id, ok := t.next()
		if !ok {
			return nil
		}

		treeNode, err := s.NodeByID(id)
		if err != nil {
			return err
		}
		if !visitor(treeNode) {
			return nil
		}

		for _, child := range treeNode.ChildrenIter() {
			t.add(child.ID())
		}
	}

	return nil
}

func (s *Segment[T]) DFS(visitor VisitorFunc[T]) error {
	return s.traverse(&stackTraverser{stack: list.NewStack()}, visitor)
}

func (s *Segment[T]) BFS(visitor VisitorFunc[T]) error {
	return s.traverse(&queueTraverser{queue: list.NewQueue()}, visitor)
}

func (s *Segment[T]) ForEachNodeAtLevel(level int, visitor VisitorFunc[T]) error {
	nodes, err := s.nodesAtLevel(level)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		if !visitor(n) {
			break
		}
	}

	return nil
}
