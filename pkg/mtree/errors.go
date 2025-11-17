package mtree

import (
	"errors"
)

var (
	ErrNodeNotFound = errors.New("node not found err")
	ErrNoMatch      = errors.New("no node match err")
	ErrNil          = errors.New("nil err")
	ErrInvalidRoot  = errors.New("invalid root err")
	ErrMaxBreadth   = errors.New("max breadth err")
)
