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

// addToLevelMap adds a node ID to the level map at the specified level.
func (s *Segment[T]) addToLevelMap(level int, id uint64) {
	if _, exists := s.levelMap[level]; !exists {
		s.levelMap[level] = make([]uint64, 0, s.maxBreadth)
	}
	s.levelMap[level] = append(s.levelMap[level], id)
}

// removeFromLevelMap removes a node ID from the level map at the specified level.
func (s *Segment[T]) removeFromLevelMap(level int, id uint64) {
	if ids, exists := s.levelMap[level]; exists {
		for i, nodeID := range ids {
			if nodeID == id {
				s.levelMap[level] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
		if len(s.levelMap[level]) == 0 {
			delete(s.levelMap, level)
		}
	}
}

// Insert adds a node to the segment. If parentID is 0 and the segment is empty,
// the node becomes the root. Otherwise, the node is attached as a child of the parent.
// This method maintains consistency between levelMap, nodeMap, and Node children relations.
func (s *Segment[T]) Insert(n *Node[T], parentID uint64) error {
	if n == nil {
		return fmt.Errorf("cannot insert: %w", ErrNil)
	}

	if _, exists := s.nodeMap[n.ID()]; exists {
		return ErrNodeAlreadyInSegment
	}

	if s.RemainingCapacity() <= 0 {
		return ErrSegmentFull
	}

	// Auto-root: if segment is empty and no parent specified, make this the root
	if s.root == nil && parentID == 0 {
		if !n.asRoot() {
			// Force root state if node was previously attached elsewhere
			n.Detach()
			n.state = root
			n.level = 0
			n.parent = nil
		}
		s.root = n
		s.nodeMap[n.ID()] = n
		s.addToLevelMap(0, n.ID())
		return nil
	}

	// If segment has nodes, we need a valid parent
	if parentID == 0 {
		return fmt.Errorf("cannot insert without parent in non-empty segment: %w", ErrParentNotInSegment)
	}

	parent, exists := s.nodeMap[parentID]
	if !exists {
		return ErrParentNotInSegment
	}

	// Check max depth constraint
	newLevel := parent.Level() + 1
	if newLevel >= s.maxDepth {
		return ErrSegmentMaxDepth
	}

	// Detach from any previous parent before inserting
	n.Detach()

	// Attach to parent (this updates Node's parent, level, state, and parent's children map)
	if err := parent.AttachChild(n); err != nil {
		return err
	}

	// Update segment maps
	s.nodeMap[n.ID()] = n
	s.addToLevelMap(n.Level(), n.ID())

	return nil
}

// RemoveCascade removes a node and all its descendants from the segment.
// This method maintains consistency between levelMap, nodeMap, and Node children relations.
func (s *Segment[T]) RemoveCascade(id uint64) error {
	n, exists := s.nodeMap[id]
	if !exists {
		return ErrNodeNotFound
	}

	// Collect all descendants using DFS
	toRemove := make([]*Node[T], 0)
	var collectDescendants func(node *Node[T])
	collectDescendants = func(node *Node[T]) {
		toRemove = append(toRemove, node)
		for _, child := range node.children {
			collectDescendants(child)
		}
	}
	collectDescendants(n)

	// Remove all collected nodes (in reverse order to handle children first)
	for i := len(toRemove) - 1; i >= 0; i-- {
		treeNode := toRemove[i]
		s.removeFromLevelMap(treeNode.Level(), treeNode.ID())
		delete(s.nodeMap, treeNode.ID())
		treeNode.Detach()
	}

	// If we removed the root, clear it
	if s.root != nil && s.root.ID() == id {
		s.root = nil
	}

	return nil
}

// RemovePromote removes a node and promotes its children to the removed node's parent.
// If the node is root and has children, returns an error (use RemoveCascade instead).
// This method maintains consistency between levelMap, nodeMap, and Node children relations.
func (s *Segment[T]) RemovePromote(id uint64) error {
	n, exists := s.nodeMap[id]
	if !exists {
		return ErrNodeNotFound
	}

	// Cannot promote children of root (they would need a new parent)
	if n.IsRoot() && n.HasChildren() {
		return ErrCannotRemoveRoot
	}

	parent := n.Parent()

	// Promote children to parent
	if parent != nil && n.HasChildren() {
		// Update levels for all descendants
		var updateLevels func(node *Node[T], levelDelta int)
		updateLevels = func(node *Node[T], levelDelta int) {
			oldLevel := node.Level()
			s.removeFromLevelMap(oldLevel, node.ID())
			node.level = oldLevel + levelDelta
			s.addToLevelMap(node.Level(), node.ID())
			for _, child := range node.children {
				updateLevels(child, levelDelta)
			}
		}

		for _, child := range n.children {
			child.Detach()
			if err := parent.AttachChild(child); err != nil {
				return err
			}
			// Update level map for child and its descendants (level decreased by 1)
			updateLevels(child, 0) // AttachChild already set correct level, just update map
		}
	}

	// Remove the node itself
	s.removeFromLevelMap(n.Level(), n.ID())
	delete(s.nodeMap, n.ID())
	n.Detach()

	// If we removed the root (which had no children), clear it
	if s.root != nil && s.root.ID() == id {
		s.root = nil
	}

	return nil
}

// Link establishes a parent-child relationship between two nodes already in the segment.
// The child is detached from its current parent (if any) and attached to the new parent.
// This method maintains consistency between levelMap, nodeMap, and Node children relations.
func (s *Segment[T]) Link(parentID, childID uint64) error {
	parent, parentExists := s.nodeMap[parentID]
	child, childExists := s.nodeMap[childID]

	if !parentExists || !childExists {
		return ErrNodesNotInSegment
	}

	// Check max depth constraint
	newLevel := parent.Level() + 1
	if newLevel >= s.maxDepth {
		return ErrSegmentMaxDepth
	}

	// Collect old levels BEFORE detaching (since Detach changes level to -1)
	type levelInfo struct {
		node     *Node[T]
		oldLevel int
	}
	var collectLevels func(node *Node[T]) []levelInfo
	collectLevels = func(node *Node[T]) []levelInfo {
		result := []levelInfo{{node: node, oldLevel: node.Level()}}
		for _, ch := range node.children {
			result = append(result, collectLevels(ch)...)
		}
		return result
	}
	oldLevels := collectLevels(child)

	// If child was root, we need to clear root
	wasRoot := child.IsRoot()

	// Detach from current parent
	child.Detach()

	// Attach to new parent
	if err := parent.AttachChild(child); err != nil {
		return err
	}

	// Update level maps: remove from old levels, add to new levels
	for _, info := range oldLevels {
		s.removeFromLevelMap(info.oldLevel, info.node.ID())
	}

	// Add to new levels (recalculate based on new position)
	var addToNewLevels func(node *Node[T])
	addToNewLevels = func(node *Node[T]) {
		s.addToLevelMap(node.Level(), node.ID())
		for _, ch := range node.children {
			addToNewLevels(ch)
		}
	}
	addToNewLevels(child)

	// If child was root, clear segment root
	if wasRoot {
		s.root = nil
	}

	return nil
}

// Unlink breaks the parent-child relationship, keeping both nodes in the segment.
// The child becomes detached (level -1, no parent) but remains in nodeMap.
// Note: The child is removed from levelMap since it no longer has a valid level.
func (s *Segment[T]) Unlink(parentID, childID uint64) error {
	parent, parentExists := s.nodeMap[parentID]
	child, childExists := s.nodeMap[childID]

	if !parentExists || !childExists {
		return ErrNodesNotInSegment
	}

	if !child.IsChildOf(parent) {
		return fmt.Errorf("child is not a child of parent: %w", ErrNodeNotFound)
	}

	// Remove child and descendants from level map (they become detached)
	var removeFromLevels func(node *Node[T])
	removeFromLevels = func(node *Node[T]) {
		s.removeFromLevelMap(node.Level(), node.ID())
		for _, ch := range node.children {
			removeFromLevels(ch)
		}
	}
	removeFromLevels(child)

	// Detach child from parent
	child.Detach()

	return nil
}

// Select returns all nodes matching the predicate function.
func (s *Segment[T]) Select(predicate func(*Node[T]) bool) []*Node[T] {
	result := make([]*Node[T], 0)
	for _, n := range s.nodeMap {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

// SelectAtLevel returns all nodes at the specified level matching the predicate.
func (s *Segment[T]) SelectAtLevel(level int, predicate func(*Node[T]) bool) ([]*Node[T], error) {
	nodes, err := s.nodesAtLevel(level)
	if err != nil {
		return nil, err
	}

	result := make([]*Node[T], 0, len(nodes))
	for _, n := range nodes {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result, nil
}

// SelectOne returns the first node matching the predicate, or error if none found.
func (s *Segment[T]) SelectOne(predicate func(*Node[T]) bool) (*Node[T], error) {
	for _, n := range s.nodeMap {
		if predicate(n) {
			return n, nil
		}
	}
	return nil, ErrNoMatch
}
