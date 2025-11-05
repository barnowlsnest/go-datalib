package tree

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type NodeValueTestSuite struct {
	suite.Suite
}

func (s *NodeValueTestSuite) TestNodeValue_Constructor() {
	testCases := []struct {
		name    string
		id      uint64
		val     int
		constFn func() *NodeValue[int]
	}{
		{
			name: "NewNode",
			id:   math.MaxUint64,
			val:  0,
			constFn: func() *NodeValue[int] {
				return NewNode[int](math.MaxUint64)
			},
		},
		{
			name: "NewNodeValue",
			id:   math.MaxUint64,
			val:  10,
			constFn: func() *NodeValue[int] {
				return NewNodeValue[int](math.MaxUint64, 10)
			},
		},
		{
			name: "Node",
			id:   math.MaxUint64,
			val:  20,
			constFn: func() *NodeValue[int] {
				return Node[int](math.MaxUint64, 20)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			n := tc.constFn()
			s.Require().NotNil(n)
			s.Require().NotNil(n.node)
			s.Require().Equal(tc.id, n.node.ID())
			s.Require().Equal(tc.val, n.val)
			s.Require().Nil(n.node.Next())
			s.Require().Nil(n.node.Prev())
		})
	}
}

func (s *NodeValueTestSuite) TestNodeValue_Props() {
	n := Node[float64](5, 0.4321)
	p, err := n.Props()
	s.Require().NoError(err)
	s.Require().Equal(uint64(5), p.ID)
	s.Require().Equal(0.4321, p.Value)
}

func (s *NodeValueTestSuite) TestNodeValue_WithValue() {
	n := NewNode[int](uint64(100))
	s.Require().Equal(0, n.val)
	n.WithValue(10)
	s.Require().Equal(10, n.val)
}

func (s *NodeValueTestSuite) TestNodeValue_WithParent() {
	testCase := []struct {
		name      string
		parent    *NodeValue[int]
		child     *NodeValue[int]
		shouldErr bool
	}{
		{
			name:      "should set parent",
			parent:    NewNode[int](uint64(1)),
			child:     NewNode[int](uint64(2)),
			shouldErr: false,
		},
		{
			name:      "should not set nil parent",
			parent:    nil,
			child:     NewNode[int](uint64(2)),
			shouldErr: true,
		},
	}

	for _, tc := range testCase {
		s.Run(tc.name, func() {
			hash, err := tc.child.WithParent(tc.parent)
			if tc.shouldErr {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)
			parentID := tc.parent.node.ID()
			childID := tc.child.node.ID()
			s.Require().Equal(serial.NSum(childID, parentID), hash)
			s.Require().True(tc.child.IsChildOf(tc.parent))
		})
	}
}

func (s *NodeValueTestSuite) TestNodeValue_HasParent() {
	n := NewNode[int](uint64(1))
	s.Require().False(n.HasParent())
	_, err := n.WithParent(NewNode[int](uint64(2)))
	s.Require().NoError(err)
	s.Require().True(n.HasParent())
	n.UnlinkParent()
	s.Require().False(n.HasParent())
}

func (s *NodeValueTestSuite) TestNodeValue_Equal() {
	testCases := []struct {
		name        string
		n1          *NodeValue[int]
		n2          *NodeValue[int]
		shouldEqual bool
	}{
		{
			name:        "should equal when values are equal",
			n1:          NewNode[int](uint64(1)),
			n2:          NewNode[int](uint64(1)),
			shouldEqual: true,
		},
		{
			name: "should equal when both have nil underlying *node.Node",
			n1: func() *NodeValue[int] {
				n := NewNode[int](uint64(1))
				n.node = nil
				return n
			}(),
			n2: func() *NodeValue[int] {
				n := NewNode[int](uint64(1))
				n.node = nil
				return n
			}(),
			shouldEqual: true,
		},
		{
			name: "should not be equal when second have nil underlying *node.Node",
			n1:   NewNode[int](uint64(1)),
			n2: func() *NodeValue[int] {
				n := NewNode[int](uint64(1))
				n.node = nil
				return n
			}(),
			shouldEqual: false,
		},
		{
			name: "should not be equal when first have nil underlying *node.Node",
			n1: func() *NodeValue[int] {
				n := NewNode[int](uint64(1))
				n.node = nil
				return n
			}(),
			n2:          nil,
			shouldEqual: false,
		},
		{
			name:        "should be equal when both nil",
			n1:          nil,
			n2:          nil,
			shouldEqual: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.shouldEqual {
				s.Require().True(tc.shouldEqual, tc.n1.Equal(tc.n2))
			} else {
				s.Require().False(tc.shouldEqual, tc.n1.Equal(tc.n2))
			}
		})
	}
}

func TestNodeValue(t *testing.T) {
	suite.Run(t, new(NodeValueTestSuite))
}
