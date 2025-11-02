package utils

import "cmp"

// MinMax returns the minimum and maximum values from a slice of ordered values.
// Returns an error if the slice is empty.
func MinMax[T cmp.Ordered](nums []T) (min, max T, err error) {
	if len(nums) == 0 {
		return min, max, ErrEmptySlice
	}

	min = nums[0]
	max = nums[0]

	for _, num := range nums[1:] {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	return min, max, nil
}