package slices

// Sum adds up each element of the input slice, returning the total result.  Empty or nil input results in zero.
func Sum[T byte | int | float32 | float64](input []T) T {
	var result T
	for _, element := range input {
		result += element
	}
	return result
}

// Avg calculates the average of the input, returning the result.  Empty or nil input results in zero.
func Avg[T byte | int | float32 | float64](input []T) float64 {
	var total T
	for _, element := range input {
		total += element
	}
	if total == 0 {
		return 0
	}
	return float64(total) / float64(len(input))
}

// Max finds the maximum value in the input, returning the result.  Empty or nil input results in zero.
func Max[T byte | int | float32 | float64](input []T) T {
	var result T
	for _, element := range input {
		if element > result {
			result = element
		}
	}
	return result
}
