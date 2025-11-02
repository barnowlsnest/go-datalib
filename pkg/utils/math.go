package utils

import "cmp"

// MinMax returns the minimum and maximum values from a slice of ordered values.
// Returns an error if the slice is empty.
func MinMax[T cmp.Ordered](nums []T) (minVal, maxVal T, err error) {
	if len(nums) == 0 {
		return minVal, maxVal, ErrEmptySlice
	}

	minVal = nums[0]
	maxVal = nums[0]

	for _, num := range nums[1:] {
		if num < minVal {
			minVal = num
		}
		if num > maxVal {
			maxVal = num
		}
	}

	return minVal, maxVal, nil
}
