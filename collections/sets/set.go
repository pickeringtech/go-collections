package sets

type Hash[T comparable] map[T]struct{}

func NewHash[T comparable](values ...T) Hash[T] {
	m := make(Hash[T])
	for _, value := range values {
		m[value] = struct{}{}
	}
	return m
}
