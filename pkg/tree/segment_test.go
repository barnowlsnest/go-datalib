package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type SegmentTestSuite struct {
	suite.Suite
	seq *serial.Serial
}

func TestSegmentTestSuite(t *testing.T) {
	suite.Run(t, new(SegmentTestSuite))
}

func (s *SegmentTestSuite) SetupTest() {
	s.seq = serial.Seq()
}

func (s *SegmentTestSuite) nextID() uint64 {
	return s.seq.Next("segment_test")
}

// buildTestSegment creates a segment with a tree structure for testing.
// Tree structure:
//
//	     root (level 0)
//	    /    \
//	child1  child2 (level 1)
//	  |
//	grandchild (level 2)
func (s *SegmentTestSuite) buildTestSegment() (seg *Segment[string], nodes map[string]*Node[string]) {
	seg = NewSegment[string]("test", s.nextID(), 5, 5)
	nodes = make(map[string]*Node[string])

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child1, err := NewNode[string](s.nextID(), 5, ValueOpt("child1"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child1, root.ID()))

	child2, err := NewNode[string](s.nextID(), 5, ValueOpt("child2"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child2, root.ID()))

	grandchild, err := NewNode[string](s.nextID(), 5, ValueOpt("grandchild"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(grandchild, child1.ID()))

	nodes["root"] = root
	nodes["child1"] = child1
	nodes["child2"] = child2
	nodes["grandchild"] = grandchild

	return seg, nodes
}

func (s *SegmentTestSuite) TestNewSegment() {
	seg := NewSegment[int]("test_segment", 1, 10, 20)

	s.Equal("test_segment", seg.Alias())
	s.Equal(uint64(1), seg.ID())
	s.Equal(200, seg.Capacity())
	s.Equal(0, seg.Height())
	s.Equal(0, seg.Length())
}

func (s *SegmentTestSuite) TestNewSegment_DefaultValues() {
	seg := NewSegment[int]("", 1, 0, 0)

	s.Equal("seg.1", seg.Alias())
	s.Equal(DefaultMaxDepth, seg.maxDepth)
	s.Equal(DefaultMaxBreadth, seg.maxBreadth)
	s.Equal(DefaultMaxDepth*DefaultMaxBreadth, seg.Capacity())
}

func (s *SegmentTestSuite) TestNewSegment_AliasWithSpaces() {
	seg := NewSegment[int]("test segment name", 1, 5, 5)

	s.Equal("testsegmentname", seg.Alias())
}

func (s *SegmentTestSuite) TestSegment_Alias() {
	seg := NewSegment[string]("myalias", s.nextID(), 5, 5)

	s.Equal("myalias", seg.Alias())
}

func (s *SegmentTestSuite) TestSegment_ID() {
	id := s.nextID()
	seg := NewSegment[string]("test", id, 5, 5)

	s.Equal(id, seg.ID())
}

func (s *SegmentTestSuite) TestSegment_Capacity() {
	seg := NewSegment[string]("test", s.nextID(), 4, 8)

	s.Equal(32, seg.Capacity())
}

func (s *SegmentTestSuite) TestSegment_Height() {
	seg, _ := s.buildTestSegment()

	s.Equal(3, seg.Height())
}

func (s *SegmentTestSuite) TestSegment_Height_Empty() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	s.Equal(0, seg.Height())
}

func (s *SegmentTestSuite) TestSegment_Length() {
	seg, _ := s.buildTestSegment()

	s.Equal(4, seg.Length())
}

func (s *SegmentTestSuite) TestSegment_Length_Empty() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	s.Equal(0, seg.Length())
}

func (s *SegmentTestSuite) TestSegment_RemainingCapacity() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"), LevelOpt[string](0))
	s.Require().NoError(err)

	seg.root = root
	seg.nodeMap[root.ID()] = root

	s.Equal(24, seg.RemainingCapacity())
}

func (s *SegmentTestSuite) TestSegment_Root() {
	seg, nodes := s.buildTestSegment()

	root, ok := seg.Root()
	s.True(ok)
	s.Equal(nodes["root"], root)
}

func (s *SegmentTestSuite) TestSegment_Root_Empty() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, ok := seg.Root()
	s.False(ok)
	s.Nil(root)
}

func (s *SegmentTestSuite) TestSegment_NodeByID() {
	seg, nodes := s.buildTestSegment()

	node, err := seg.NodeByID(nodes["child1"].ID())
	s.NoError(err)
	s.Equal(nodes["child1"], node)
}

func (s *SegmentTestSuite) TestSegment_NodeByID_NotFound() {
	seg, _ := s.buildTestSegment()

	node, err := seg.NodeByID(9999)
	s.Error(err)
	s.ErrorIs(err, ErrNodeNotFound)
	s.Nil(node)
}

func (s *SegmentTestSuite) TestSegment_DFS_Empty() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	visited := make([]string, 0)

	err := seg.DFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Empty(visited)
}

func (s *SegmentTestSuite) TestSegment_DFS_VisitsAllNodes() {
	seg, nodes := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.DFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Len(visited, 4)
	s.Contains(visited, "root")
	s.Contains(visited, "child1")
	s.Contains(visited, "child2")
	s.Contains(visited, "grandchild")
	_ = nodes
}

func (s *SegmentTestSuite) TestSegment_DFS_RootFirst() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.DFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Equal("root", visited[0])
}

func (s *SegmentTestSuite) TestSegment_DFS_EarlyStop() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.DFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return n.Val() != "child1"
	})

	s.NoError(err)
	s.Contains(visited, "root")
}

func (s *SegmentTestSuite) TestSegment_BFS_Empty() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	visited := make([]string, 0)

	err := seg.BFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Empty(visited)
}

func (s *SegmentTestSuite) TestSegment_BFS_VisitsAllNodes() {
	seg, nodes := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.BFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Len(visited, 4)
	s.Contains(visited, "root")
	s.Contains(visited, "child1")
	s.Contains(visited, "child2")
	s.Contains(visited, "grandchild")
	_ = nodes
}

func (s *SegmentTestSuite) TestSegment_BFS_LevelOrder() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.BFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Equal("root", visited[0])
	s.Equal("grandchild", visited[len(visited)-1])
}

func (s *SegmentTestSuite) TestSegment_BFS_EarlyStop() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.BFS(func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return n.Val() != "child1"
	})

	s.NoError(err)
	s.Contains(visited, "root")
}

func (s *SegmentTestSuite) TestSegment_ForEachNodeAtLevel_Level0() {
	seg, nodes := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.ForEachNodeAtLevel(0, func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Len(visited, 1)
	s.Equal("root", visited[0])
	_ = nodes
}

func (s *SegmentTestSuite) TestSegment_ForEachNodeAtLevel_Level1() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.ForEachNodeAtLevel(1, func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Len(visited, 2)
	s.Contains(visited, "child1")
	s.Contains(visited, "child2")
}

func (s *SegmentTestSuite) TestSegment_ForEachNodeAtLevel_Level2() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.ForEachNodeAtLevel(2, func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return true
	})

	s.NoError(err)
	s.Len(visited, 1)
	s.Equal("grandchild", visited[0])
}

func (s *SegmentTestSuite) TestSegment_ForEachNodeAtLevel_EarlyStop() {
	seg, _ := s.buildTestSegment()
	visited := make([]string, 0)

	err := seg.ForEachNodeAtLevel(1, func(n *Node[string]) bool {
		visited = append(visited, n.Val())
		return false
	})

	s.NoError(err)
	s.Len(visited, 1)
}

func (s *SegmentTestSuite) TestSegment_ForEachNodeAtLevel_NotFound() {
	seg, _ := s.buildTestSegment()

	err := seg.ForEachNodeAtLevel(10, func(n *Node[string]) bool {
		return true
	})

	s.Error(err)
	s.ErrorIs(err, ErrSegmentLevelNotFound)
}

func (s *SegmentTestSuite) TestSegment_nodesAtLevel_NodeNotInMap() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"), LevelOpt[string](0))
	s.Require().NoError(err)

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"), ParentOpt(root))
	s.Require().NoError(err)

	seg.root = root
	seg.nodeMap[root.ID()] = root
	seg.levelMap[0] = []uint64{root.ID()}
	seg.levelMap[1] = []uint64{child.ID()}

	nodes, err := seg.nodesAtLevel(1)
	s.Error(err)
	s.ErrorIs(err, ErrSegmentDoesNotHaveNode)
	s.Nil(nodes)
}

// ============================================================================
// Insert Tests
// ============================================================================

func (s *SegmentTestSuite) TestSegment_Insert_RootAutomatic() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)

	err = seg.Insert(root, 0)
	s.NoError(err)

	s.Equal(1, seg.Length())
	s.Equal(1, seg.Height())
	gotRoot, ok := seg.Root()
	s.True(ok)
	s.Equal(root, gotRoot)
	s.Equal(0, root.Level())
	s.True(root.IsRoot())
}

func (s *SegmentTestSuite) TestSegment_Insert_Child() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"))
	s.Require().NoError(err)

	err = seg.Insert(child, root.ID())
	s.NoError(err)

	s.Equal(2, seg.Length())
	s.Equal(2, seg.Height())
	s.Equal(1, child.Level())
	s.True(child.IsChildOf(root))
}

func (s *SegmentTestSuite) TestSegment_Insert_NilNode() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	err := seg.Insert(nil, 0)
	s.Error(err)
	s.ErrorIs(err, ErrNil)
}

func (s *SegmentTestSuite) TestSegment_Insert_AlreadyInSegment() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	err = seg.Insert(root, 0)
	s.Error(err)
	s.ErrorIs(err, ErrNodeAlreadyInSegment)
}

func (s *SegmentTestSuite) TestSegment_Insert_NoParentInNonEmptySegment() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"))
	s.Require().NoError(err)

	err = seg.Insert(child, 0)
	s.Error(err)
	s.ErrorIs(err, ErrParentNotInSegment)
}

func (s *SegmentTestSuite) TestSegment_Insert_ParentNotInSegment() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"))
	s.Require().NoError(err)

	err = seg.Insert(child, 99999)
	s.Error(err)
	s.ErrorIs(err, ErrParentNotInSegment)
}

func (s *SegmentTestSuite) TestSegment_Insert_MaxDepthExceeded() {
	seg := NewSegment[string]("test", s.nextID(), 5, 2) // max depth of 2

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child, root.ID()))

	grandchild, err := NewNode[string](s.nextID(), 5, ValueOpt("grandchild"))
	s.Require().NoError(err)

	err = seg.Insert(grandchild, child.ID())
	s.Error(err)
	s.ErrorIs(err, ErrSegmentMaxDepth)
}

func (s *SegmentTestSuite) TestSegment_Insert_CapacityExceeded() {
	seg := NewSegment[string]("test", s.nextID(), 1, 1) // capacity of 1

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"))
	s.Require().NoError(err)

	err = seg.Insert(child, root.ID())
	s.Error(err)
	s.ErrorIs(err, ErrSegmentFull)
}

// ============================================================================
// RemoveCascade Tests
// ============================================================================

func (s *SegmentTestSuite) TestSegment_RemoveCascade_SingleNode() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	err = seg.RemoveCascade(root.ID())
	s.NoError(err)

	s.Equal(0, seg.Length())
	s.Equal(0, seg.Height())
	gotRoot, ok := seg.Root()
	s.False(ok)
	s.Nil(gotRoot)
}

func (s *SegmentTestSuite) TestSegment_RemoveCascade_WithDescendants() {
	seg, nodes := s.buildTestSegment()

	err := seg.RemoveCascade(nodes["child1"].ID())
	s.NoError(err)

	// child1 and grandchild should be removed
	s.Equal(2, seg.Length()) // root and child2 remain
	_, err = seg.NodeByID(nodes["child1"].ID())
	s.ErrorIs(err, ErrNodeNotFound)
	_, err = seg.NodeByID(nodes["grandchild"].ID())
	s.ErrorIs(err, ErrNodeNotFound)

	// root and child2 should still be there
	_, err = seg.NodeByID(nodes["root"].ID())
	s.NoError(err)
	_, err = seg.NodeByID(nodes["child2"].ID())
	s.NoError(err)
}

func (s *SegmentTestSuite) TestSegment_RemoveCascade_Root() {
	seg, nodes := s.buildTestSegment()

	err := seg.RemoveCascade(nodes["root"].ID())
	s.NoError(err)

	s.Equal(0, seg.Length())
	s.Equal(0, seg.Height())
	gotRoot, ok := seg.Root()
	s.False(ok)
	s.Nil(gotRoot)
}

func (s *SegmentTestSuite) TestSegment_RemoveCascade_NotFound() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	err := seg.RemoveCascade(99999)
	s.Error(err)
	s.ErrorIs(err, ErrNodeNotFound)
}

// ============================================================================
// RemovePromote Tests
// ============================================================================

func (s *SegmentTestSuite) TestSegment_RemovePromote_LeafNode() {
	seg, nodes := s.buildTestSegment()

	err := seg.RemovePromote(nodes["grandchild"].ID())
	s.NoError(err)

	s.Equal(3, seg.Length())
	_, err = seg.NodeByID(nodes["grandchild"].ID())
	s.ErrorIs(err, ErrNodeNotFound)
}

func (s *SegmentTestSuite) TestSegment_RemovePromote_MiddleNode() {
	seg, nodes := s.buildTestSegment()

	err := seg.RemovePromote(nodes["child1"].ID())
	s.NoError(err)

	s.Equal(3, seg.Length())
	_, err = seg.NodeByID(nodes["child1"].ID())
	s.ErrorIs(err, ErrNodeNotFound)

	// grandchild should be promoted to root's children
	grandchild, err := seg.NodeByID(nodes["grandchild"].ID())
	s.NoError(err)
	s.True(grandchild.IsChildOf(nodes["root"]))
	s.Equal(1, grandchild.Level())
}

func (s *SegmentTestSuite) TestSegment_RemovePromote_RootWithChildren() {
	seg, nodes := s.buildTestSegment()

	err := seg.RemovePromote(nodes["root"].ID())
	s.Error(err)
	s.ErrorIs(err, ErrCannotRemoveRoot)
}

func (s *SegmentTestSuite) TestSegment_RemovePromote_RootWithoutChildren() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)
	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	err = seg.RemovePromote(root.ID())
	s.NoError(err)

	s.Equal(0, seg.Length())
	gotRoot, ok := seg.Root()
	s.False(ok)
	s.Nil(gotRoot)
}

func (s *SegmentTestSuite) TestSegment_RemovePromote_NotFound() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	err := seg.RemovePromote(99999)
	s.Error(err)
	s.ErrorIs(err, ErrNodeNotFound)
}

// ============================================================================
// Link Tests
// ============================================================================

func (s *SegmentTestSuite) TestSegment_Link_Basic() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child1, err := NewNode[string](s.nextID(), 5, ValueOpt("child1"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child1, root.ID()))

	child2, err := NewNode[string](s.nextID(), 5, ValueOpt("child2"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child2, root.ID()))

	// Link child2 as child of child1
	err = seg.Link(child1.ID(), child2.ID())
	s.NoError(err)

	s.True(child2.IsChildOf(child1))
	s.Equal(2, child2.Level())
}

func (s *SegmentTestSuite) TestSegment_Link_NodeNotInSegment() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	err = seg.Link(root.ID(), 99999)
	s.Error(err)
	s.ErrorIs(err, ErrNodesNotInSegment)
}

func (s *SegmentTestSuite) TestSegment_Link_MaxDepthExceeded() {
	seg := NewSegment[string]("test", s.nextID(), 5, 2)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child, err := NewNode[string](s.nextID(), 5, ValueOpt("child"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child, root.ID()))

	err = seg.Link(child.ID(), root.ID())
	s.Error(err)
	s.ErrorIs(err, ErrSegmentMaxDepth)
}

// ============================================================================
// Unlink Tests
// ============================================================================

func (s *SegmentTestSuite) TestSegment_Unlink_Basic() {
	seg, nodes := s.buildTestSegment()

	err := seg.Unlink(nodes["child1"].ID(), nodes["grandchild"].ID())
	s.NoError(err)

	// grandchild should be detached but still in nodeMap
	grandchild, err := seg.NodeByID(nodes["grandchild"].ID())
	s.NoError(err)
	s.True(grandchild.IsDetached())
	s.Equal(-1, grandchild.Level())
	s.False(grandchild.HasParent())

	// child1 should no longer have grandchild as child
	s.False(nodes["child1"].HasChild(grandchild))
}

func (s *SegmentTestSuite) TestSegment_Unlink_NodeNotInSegment() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	err = seg.Unlink(root.ID(), 99999)
	s.Error(err)
	s.ErrorIs(err, ErrNodesNotInSegment)
}

func (s *SegmentTestSuite) TestSegment_Unlink_NotChildOf() {
	seg, nodes := s.buildTestSegment()

	err := seg.Unlink(nodes["child2"].ID(), nodes["grandchild"].ID())
	s.Error(err)
	s.ErrorIs(err, ErrNodeNotFound)
}

// ============================================================================
// Select Tests
// ============================================================================

func (s *SegmentTestSuite) TestSegment_Select_All() {
	seg, _ := s.buildTestSegment()

	nodes := seg.Select(func(n *Node[string]) bool {
		return true
	})

	s.Len(nodes, 4)
}

func (s *SegmentTestSuite) TestSegment_Select_ByValue() {
	seg, _ := s.buildTestSegment()

	nodes := seg.Select(func(n *Node[string]) bool {
		return n.Val() == "child1"
	})

	s.Len(nodes, 1)
	s.Equal("child1", nodes[0].Val())
}

func (s *SegmentTestSuite) TestSegment_Select_NoMatch() {
	seg, _ := s.buildTestSegment()

	nodes := seg.Select(func(n *Node[string]) bool {
		return n.Val() == "nonexistent"
	})

	s.Len(nodes, 0)
}

func (s *SegmentTestSuite) TestSegment_SelectAtLevel() {
	seg, _ := s.buildTestSegment()

	nodes, err := seg.SelectAtLevel(1, func(n *Node[string]) bool {
		return true
	})

	s.NoError(err)
	s.Len(nodes, 2)
}

func (s *SegmentTestSuite) TestSegment_SelectAtLevel_WithPredicate() {
	seg, _ := s.buildTestSegment()

	nodes, err := seg.SelectAtLevel(1, func(n *Node[string]) bool {
		return n.Val() == "child1"
	})

	s.NoError(err)
	s.Len(nodes, 1)
	s.Equal("child1", nodes[0].Val())
}

func (s *SegmentTestSuite) TestSegment_SelectAtLevel_InvalidLevel() {
	seg, _ := s.buildTestSegment()

	nodes, err := seg.SelectAtLevel(99, func(n *Node[string]) bool {
		return true
	})

	s.Error(err)
	s.ErrorIs(err, ErrSegmentLevelNotFound)
	s.Nil(nodes)
}

func (s *SegmentTestSuite) TestSegment_SelectOne() {
	seg, _ := s.buildTestSegment()

	node, err := seg.SelectOne(func(n *Node[string]) bool {
		return n.Val() == "grandchild"
	})

	s.NoError(err)
	s.Equal("grandchild", node.Val())
}

func (s *SegmentTestSuite) TestSegment_SelectOne_NoMatch() {
	seg, _ := s.buildTestSegment()

	node, err := seg.SelectOne(func(n *Node[string]) bool {
		return n.Val() == "nonexistent"
	})

	s.Error(err)
	s.ErrorIs(err, ErrNoMatch)
	s.Nil(node)
}

// ============================================================================
// Integration Tests - Consistency Verification
// ============================================================================

func (s *SegmentTestSuite) TestSegment_Insert_MapsConsistency() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child1, err := NewNode[string](s.nextID(), 5, ValueOpt("child1"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child1, root.ID()))

	child2, err := NewNode[string](s.nextID(), 5, ValueOpt("child2"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child2, root.ID()))

	grandchild, err := NewNode[string](s.nextID(), 5, ValueOpt("grandchild"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(grandchild, child1.ID()))

	// Verify nodeMap
	s.Equal(4, len(seg.nodeMap))

	// Verify levelMap
	s.Equal(3, len(seg.levelMap))
	s.Len(seg.levelMap[0], 1)
	s.Len(seg.levelMap[1], 2)
	s.Len(seg.levelMap[2], 1)

	// Verify Node relations
	s.True(root.IsRoot())
	s.True(child1.IsChildOf(root))
	s.True(child2.IsChildOf(root))
	s.True(grandchild.IsChildOf(child1))
}

func (s *SegmentTestSuite) TestSegment_RemoveCascade_MapsConsistency() {
	seg, nodes := s.buildTestSegment()

	err := seg.RemoveCascade(nodes["child1"].ID())
	s.NoError(err)

	// Verify nodeMap
	s.Equal(2, len(seg.nodeMap))
	_, exists := seg.nodeMap[nodes["child1"].ID()]
	s.False(exists)
	_, exists = seg.nodeMap[nodes["grandchild"].ID()]
	s.False(exists)

	// Verify levelMap
	s.Len(seg.levelMap[0], 1)
	s.Len(seg.levelMap[1], 1) // only child2 remains
	_, exists = seg.levelMap[2]
	s.False(exists) // level 2 should be deleted

	// Verify Node relations
	s.False(nodes["root"].HasChild(nodes["child1"]))
}

func (s *SegmentTestSuite) TestSegment_Link_MapsConsistency() {
	seg := NewSegment[string]("test", s.nextID(), 5, 5)

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(root, 0))

	child1, err := NewNode[string](s.nextID(), 5, ValueOpt("child1"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child1, root.ID()))

	child2, err := NewNode[string](s.nextID(), 5, ValueOpt("child2"))
	s.Require().NoError(err)
	s.Require().NoError(seg.Insert(child2, root.ID()))

	// Move child2 under child1
	err = seg.Link(child1.ID(), child2.ID())
	s.NoError(err)

	// Verify levelMap updated correctly
	s.Len(seg.levelMap[1], 1) // only child1 at level 1
	s.Len(seg.levelMap[2], 1) // child2 moved to level 2

	// Verify Node relations
	s.True(child2.IsChildOf(child1))
	s.False(child2.IsChildOf(root))
	s.Equal(2, child2.Level())
}
