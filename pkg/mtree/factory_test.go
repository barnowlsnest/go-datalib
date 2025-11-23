package mtree

import (
	"testing"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
	"github.com/stretchr/testify/suite"
)

const testMaxBreadth = 3

type NewNodeFactoryTestSuite struct {
	suite.Suite
}

func TestNewNodeFactory(t *testing.T) {
	suite.Run(t, new(NewNodeFactoryTestSuite))
}

func (s *NewNodeFactoryTestSuite) TestNode() {
	factory := NewFactory[int]("node", testMaxBreadth, serial.Seq())
	n1 := factory.Node(ValueOpt[int](10))
	s.NotNil(n1)
	n2 := factory.Node(ValueOpt[int](20))
	s.NotNil(n2)
	n3 := factory.Node(ValueOpt[int](30), ChildOpt[int](n1), ChildOpt[int](n2))
	s.NotNil(n3)
	s.True(n1.IsChildOf(n3))
	s.True(n2.IsChildOf(n3))
	s.False(n3.IsChildOf(n1))
	s.False(n3.IsChildOf(n2))
	s.True(n3.HasChildren())
}

func (s *NewNodeFactoryTestSuite) TestRoot() {
	factory := NewFactory[int]("root", testMaxBreadth, serial.Seq())
	rootNode := factory.Root(50)
	s.Nil(rootNode)
	s.Equal(50, rootNode.Val())
	s.True(rootNode.IsRoot())
}

func (s *NewNodeFactoryTestSuite) TestRootWithChildren() {
	factory := NewFactory[int]("rootWithChildren", testMaxBreadth, serial.Seq())
	n1 := factory.Node(ValueOpt[int](1))
	n2 := factory.Node(ValueOpt[int](2))
	n3 := factory.Node(ValueOpt[int](3))
	r, err := factory.RootWithChildren(100, n1, n2, n3)
	s.NoError(err)
	s.NotNil(r)
	s.True(r.HasChildren())
	s.True(n1.IsChildOf(r))
	s.True(n2.IsChildOf(r))
	s.True(n3.IsChildOf(r))
	s.Equal(3, r.Breadth())
	s.Equal(0, r.Capacity())
}
