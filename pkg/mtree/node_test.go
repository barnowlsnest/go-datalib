package mtree

import (
	"testing"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
	"github.com/stretchr/testify/suite"
)

const defaultMaxBreadth = 3

type NodeTestSuite struct {
	suite.Suite
}

func TestNode(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

func (s *NodeTestSuite) ID() uint64 {
	s.T().Helper()
	return serial.Seq().Next("NodeTestSuite")
}

func (s *NodeTestSuite) TestNewNodeDefaults() {
	var (
		id = s.ID()
	)
	testCases := []struct {
		name string
		spec func(n *Node[int])
	}{
		{
			name: "should not has parent by default",
			spec: func(n *Node[int]) {
				s.Nil(n.Parent())
				s.False(n.HasParent())
				s.False(n.IsChildOf(nil))
			},
		},
		{
			name: "should not be root, but detached by default",
			spec: func(n *Node[int]) {
				s.False(n.IsRoot())
				s.False(n.IsAttached())
				s.True(n.IsDetached())
			},
		},
		{
			name: "should has id",
			spec: func(n *Node[int]) {
				s.Equal(id, n.ID())
			},
		},
		{
			name: "should has MaxBreadth and max capacity by default",
			spec: func(n *Node[int]) {
				s.Equal(defaultMaxBreadth, n.MaxBreadth())
				s.Equal(n.MaxBreadth(), n.Capacity())
			},
		},
		{
			name: "should not have children by default",
			spec: func(n *Node[int]) {
				s.False(n.HasChildren())
				s.Equal(0, n.Breadth())
			},
		},
		{
			name: "should have -1 level by default",
			spec: func(n *Node[int]) {
				s.Equal(-1, n.Level())
			},
		},
	}

	n := NewNode[int](DefaultShard, id, defaultMaxBreadth)
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.spec(n)
		})
	}
}

func (s *NodeTestSuite) TestNewRoot() {
	rootNode := NewRoot[int](DefaultShard, 1, 3, 100)
	s.NotNil(rootNode)
	s.True(rootNode.IsRoot())
	s.True(rootNode.AsRoot())
	s.Equal(0, rootNode.Level())
}
