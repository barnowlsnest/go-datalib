package tree

import (
	"errors"
)

var (
	ErrNil                    = errors.New("nil err")
	ErrNodeNotFound           = errors.New("node not found err")
	ErrNoMatch                = errors.New("no node match err")
	ErrMaxBreadth             = errors.New("max breadth err")
	ErrRootTagNotFound        = errors.New("err root tag not found")
	ErrHierarchyModel         = errors.New("invalid hierarchy model")
	ErrSegmentLevelNotFound   = errors.New("err segment level not found")
	ErrNotRoot                = errors.New("err invalid root")
	ErrSegmentAlreadyHasNode  = errors.New("err segment already has node")
	ErrSegmentDoesNotHaveNode = errors.New("err segment node not found")
)
