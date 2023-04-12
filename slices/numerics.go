package slices

// Sum adds up each element of the input slice, returning the total result
func Sum[T byte | int | float32 | float64](input []T) T {
	var result T
	for _, element := range input {
		result += element
	}
	return result
}
