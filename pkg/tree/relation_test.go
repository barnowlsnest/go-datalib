package tree

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RelationTestSuite struct {
	suite.Suite
}

func (s *RelationTestSuite) TestNewRelation() {
	n1, n2 := NewNode[int](1), NewNode[int](2)

	var (
		r   *Relation[int]
		err error
	)
	r, err = NewRelation[int](n1, n2)
	s.Require().NoError(err)
	s.Require().NotNil(r)
	h1 := r.hash

	r, err = NewRelation[int](n2, n1)
	s.Require().NoError(err)
	s.Require().NotNil(r)
	h2 := r.hash

	s.Require().Equal(h1, h2)
}

func (s *RelationTestSuite) TestRel() {
	n1, n2 := NewNode[int](1), NewNode[int](2)
	r, err := Rel[int](n1, n2)
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func (s *RelationTestSuite) TestProps() {
	n1, n2 := Node[int](1, 3), Node[int](2, 5)
	r, err := Rel[int](n1, n2)
	s.Require().NoError(err)
	p := r.Props()

	nProps1, err := n1.Props()
	s.Require().NoError(err)

	nProps2, err := n2.Props()
	s.Require().NoError(err)

	s.Require().Equal(r.hash, p.Hash)
	s.Require().Equal(nProps1.ID, p.ParentID)
	s.Require().Equal(nProps2.ID, p.ChildID)
}

func (s *RelationTestSuite) TestType() {
	n1, n2 := Node[int](1, 10), Node[int](2, 20)
	r, err := Rel[int](n1, n2)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		id           uint64
		expectedType string
	}{
		{
			name:         "parent node returns ParentRelation",
			id:           1,
			expectedType: ParentRelation,
		},
		{
			name:         "child node returns ChildRelation",
			id:           2,
			expectedType: ChildRelation,
		},
		{
			name:         "unrelated node returns Unrelated",
			id:           999,
			expectedType: Unrelated,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := r.Type(tc.id)
			s.Require().Equal(tc.expectedType, result)
		})
	}
}

func (s *RelationTestSuite) TestEqual() {
	n1, n2 := Node[int](1, 10), Node[int](2, 20)
	n3, n4 := Node[int](3, 30), Node[int](4, 40)

	testCases := []struct {
		name          string
		r1            *Relation[int]
		r2            *Relation[int]
		shouldBeEqual bool
	}{
		{
			name: "equal relations with same parent and child",
			r1: func() *Relation[int] {
				r, _ := Rel[int](n1, n2)
				return r
			}(),
			r2: func() *Relation[int] {
				r, _ := Rel[int](n1, n2)
				return r
			}(),
			shouldBeEqual: true,
		},
		{
			name: "equal relations with reversed order (same hash)",
			r1: func() *Relation[int] {
				r, _ := Rel[int](n1, n2)
				return r
			}(),
			r2: func() *Relation[int] {
				r, _ := Rel[int](n2, n1)
				return r
			}(),
			shouldBeEqual: true,
		},
		{
			name: "not equal relations with different nodes",
			r1: func() *Relation[int] {
				r, _ := Rel[int](n1, n2)
				return r
			}(),
			r2: func() *Relation[int] {
				r, _ := Rel[int](n3, n4)
				return r
			}(),
			shouldBeEqual: false,
		},
		{
			name:          "both relations are nil",
			r1:            nil,
			r2:            nil,
			shouldBeEqual: true,
		},
		{
			name: "first relation is nil, second is not",
			r1:   nil,
			r2: func() *Relation[int] {
				r, _ := Rel[int](n1, n2)
				return r
			}(),
			shouldBeEqual: false,
		},
		{
			name: "first relation is not nil, second is nil",
			r1: func() *Relation[int] {
				r, _ := Rel[int](n1, n2)
				return r
			}(),
			r2:            nil,
			shouldBeEqual: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			result := tc.r1.Equal(tc.r2)
			s.Require().Equal(tc.shouldBeEqual, result)
		})
	}
}

func TestRelation(t *testing.T) {
	suite.Run(t, new(RelationTestSuite))
}
