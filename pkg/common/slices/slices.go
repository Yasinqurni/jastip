package slices

// Contains check if a slice contains a value
func Contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

// Max find highest value in slices
func Max[T float32 | float64 | int | int64](slice []T) T {
	max := slice[0]
	for _, v := range slice {
		if v > max {
			max = v
		}
	}

	return max
}

// Min find lowest value in slices
func Min[T float32 | float64 | int | int64](slice []T) T {
	min := slice[0]
	for _, v := range slice {
		if v < min {
			min = v
		}
	}

	return min
}
