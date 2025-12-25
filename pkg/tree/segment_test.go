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

	root, err := NewNode[string](s.nextID(), 5, ValueOpt("root"), LevelOpt[string](0))
	s.Require().NoError(err)

	child1, err := NewNode[string](s.nextID(), 5, ValueOpt("child1"), ParentOpt(root))
	s.Require().NoError(err)

	child2, err := NewNode[string](s.nextID(), 5, ValueOpt("child2"), ParentOpt(root))
	s.Require().NoError(err)

	grandchild, err := NewNode[string](s.nextID(), 5, ValueOpt("grandchild"), ParentOpt(child1))
	s.Require().NoError(err)

	seg.root = root
	seg.nodeMap[root.ID()] = root
	seg.nodeMap[child1.ID()] = child1
	seg.nodeMap[child2.ID()] = child2
	seg.nodeMap[grandchild.ID()] = grandchild

	seg.levelMap[0] = []uint64{root.ID()}
	seg.levelMap[1] = []uint64{child1.ID(), child2.ID()}
	seg.levelMap[2] = []uint64{grandchild.ID()}

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
