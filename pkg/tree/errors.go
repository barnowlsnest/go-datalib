package tree

import (
	"errors"
)

var (
	ErrNilParent          = errors.New("parent is nil")
	ErrChildNotFound      = errors.New("child not found")
	ErrNotAllowed         = errors.New("not allowed")
	ErrNotAllowedMaxNodes = errors.New("not allowed: add nodes more then maxChildrenPerNode")
	ErrNotAllowedMaxDepth = errors.New("not allowed: maxDepth reached")
	ErrLevelNotFound      = errors.New("level not found")
	ErrNodeNotFound       = errors.New("node not found")
	ErrUnexpected         = errors.New("unexpected err")
)
