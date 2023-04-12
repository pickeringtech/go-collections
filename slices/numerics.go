package slices

// Sum adds up each element of the input slice, returning the total result
func Sum[T byte | int | float32 | float64](input []T) T {
	var result T
	for _, element := range input {
		result += element
	}
	return result
}

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
