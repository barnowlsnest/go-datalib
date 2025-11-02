package utils

import (
	"errors"
)

var (
	ErrIntToUint8 = errors.New("too big value for uint8")
	ErrEmptySlice = errors.New("slice is empty")
)
