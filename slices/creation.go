package slices

type GeneratorFunc[T any] func(index int) T

func Generate[T any](amount int, fn GeneratorFunc[T]) []T {
	var results []T
	for i := 0; i < amount; i++ {
		results = append(results, fn(i))
	}
	return results
}
