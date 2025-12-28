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
	ErrSegmentLevelNotFound   = errors.New("segment level not found")
	ErrSegmentDoesNotHaveNode = errors.New("segment does not contain node")
	ErrSegmentFull            = errors.New("segment capacity exceeded")
	ErrSegmentMaxDepth        = errors.New("segment max depth exceeded")
	ErrNodeAlreadyInSegment   = errors.New("node already exists in segment")
	ErrParentNotInSegment     = errors.New("parent node not in segment")
	ErrCannotRemoveRoot       = errors.New("cannot remove root with children using promote strategy")
	ErrNodesNotInSegment      = errors.New("one or both nodes not in segment")
)
