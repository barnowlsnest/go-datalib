package utils

import (
	"math"
)

func SafeIntToUint8(n int) (uint8, error) {
	if n < 0 || n > math.MaxUint8 {
		return 0, ErrIntToUint8
	}
	return uint8(n), nil
}
