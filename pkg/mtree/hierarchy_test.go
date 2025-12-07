package mtree

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type HierarchyTestSuite struct {
	suite.Suite
}

func TestHierarchyTestSuite(t *testing.T) {
	suite.Run(t, new(HierarchyTestSuite))
}

func nextID() uint64 {
	return serial.Seq().Next("hierarchyTest")
}

func (s *HierarchyTestSuite) TestHierarchy() {
	expectedModel := HierarchyModel{
		RootTag: ChildrenSlice{"A"},
		"A":     ChildrenSlice{"B", "C", "D"},
		"D":     ChildrenSlice{"E", "F"},
	}

	n, err := Hierarchy(expectedModel, 5, nextID)
	s.NotNil(n)
	s.NoError(err)

	actualModel, err := ToModel(n)
	s.NotNil(actualModel)
	s.NoError(err)
	s.ElementsMatch(expectedModel["A"], actualModel["A"])
	s.ElementsMatch(expectedModel["D"], actualModel["D"])
	s.ElementsMatch(expectedModel[RootTag], actualModel[RootTag])
}
