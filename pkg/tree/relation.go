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
		parent NodeProps[T]
		child  NodeProps[T]
	}

	RelationProps[T cmp.Ordered] struct {
		Hash     uint64
		ParentID uint64
		ChildID  uint64
	}
)

func NewRelation[T cmp.Ordered](parent, child NodeProps[T]) *Relation[T] {
	return &Relation[T]{
		hash:   serial.NSum(parent.ID, child.ID),
		parent: parent,
		child:  child,
	}
}

func Rel[T cmp.Ordered](parent, child NodeProps[T]) *Relation[T] {
	return NewRelation(parent, child)
}

func (r *Relation[T]) Props() RelationProps[T] {
	return RelationProps[T]{
		Hash:     r.hash,
		ParentID: r.parent.ID,
		ChildID:  r.child.ID,
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
