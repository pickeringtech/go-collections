package collections

type List[T any] []T

func NewList[T any]() List[T] {
	return List[T]{}
}

// AtOrDefault retrieves the element at the given index or returns the specified default value.
func (l List[T]) AtOrDefault(idx int, defaultValue T) T {
	if idx < 0 || idx >= l.Size() {
		return defaultValue
	}
	return l[idx]
}

func (l List[T]) Size() int {
	return len(l)
}
