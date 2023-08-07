package dicts

type Hash[K comparable, V any] map[K]V

func NewHash[K comparable, V any](entries ...Pair[K, V]) Hash[K, V] {
	m := make(Hash[K, V])
	for _, entry := range entries {
		m[entry.Key] = entry.Value
	}
	return m
}
