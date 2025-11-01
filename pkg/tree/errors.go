package tree

import (
	"errors"
)

var (
	ErrNilParent     = errors.New("parent is nil")
	ErrChildNotFound = errors.New("child not found")
)
