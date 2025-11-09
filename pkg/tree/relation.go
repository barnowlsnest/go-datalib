package tree

import (
	"cmp"
	"fmt"

	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

const (
	ParentRelation = "parent"
	ChildRelation  = "child"
	Unrelated      = "unrelated"
)

type (
	Relation[T cmp.Ordered] struct {
		hash   uint64
		parent *NodeProps[T]
		child  *NodeProps[T]
	}

	RelationProps[T cmp.Ordered] struct {
		Hash     uint64
		ParentID uint64
		ChildID  uint64
	}
)

func NewRelation[T cmp.Ordered](parent, child *NodeValue[T]) (*Relation[T], error) {
	if parent == nil {
		return nil, ErrParentNil
	}

	if child == nil {
		return nil, ErrNil
	}

	var (
		p, c NodeProps[T]
		err  error
	)

	p, err = parent.Props()
	if err != nil {
		return nil, err
	}

	c, err = child.Props()
	if err != nil {
		return nil, err
	}

	return &Relation[T]{
		hash:   serial.NSum(p.ID, c.ID),
		parent: &p,
		child:  &c,
	}, nil
}

func Rel[T cmp.Ordered](parent, child *NodeValue[T]) (*Relation[T], error) {
	return NewRelation(parent, child)
}

func (r *Relation[T]) Props() RelationProps[T] {
	var parentID, childID uint64
	if r.parent != nil {
		parentID = r.parent.ID
	}

	if r.child != nil {
		childID = r.child.ID
	}

	return RelationProps[T]{
		Hash:     r.hash,
		ParentID: parentID,
		ChildID:  childID,
	}
}

func (r *Relation[T]) Type(id uint64) string {
	switch {
	case r.parent.ID == id:
		return ParentRelation
	case r.child.ID == id:
		return ChildRelation
	default:
		return Unrelated
	}
}

func (r *Relation[T]) Equal(other *Relation[T]) bool {
	if r == nil {
		return other == nil
	}

	if other == nil {
		return false
	}

	return r.hash == other.hash
}

func (r *Relation[T]) String() string {
	if r == nil {
		return "nil"
	}

	return fmt.Sprintf("%d:%d,%d", r.hash, r.parent.ID, r.child.ID)
}
