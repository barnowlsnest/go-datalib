package tree

import (
	"errors"
)

var (
	ErrNil       = errors.New("nil err")
	ErrParentNil = errors.New("parent nil err")
)
