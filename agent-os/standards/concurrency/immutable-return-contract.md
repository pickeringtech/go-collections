# Concurrent Return Contract

Operating on a thread-safe collection yields a thread-safe result.

- **Immutable ops** (`Filter`, `Add`, `Remove`, `Union`...) return a **new instance of the same concurrent type**, behind the interface (`Set[T]`/`Dict[K,V]`/`List[T]`). Build the copy while holding the lock, then wrap it in the matching concurrent type.
- **`InPlace` ops** mutate the receiver and return **void** (or a status `bool`) — never the collection. The caller already holds it.

```go
// Correct: ConcurrentHashRW.Filter returns another ConcurrentHashRW
func (ch *ConcurrentHashRW[T]) Filter(fn func(T) bool) Set[T] {
	ch.lock.RLock()
	elements := make([]T, 0, len(ch.data))
	for element := range ch.data {
		elements = append(elements, element)
	}
	ch.lock.RUnlock() // released before the predicate runs — see callback-reentrancy

	result := NewConcurrentHashRW[T]()
	for _, element := range elements {
		if fn(element) {
			result.data[element] = struct{}{}
		}
	}
	return result
}
```

- Principle of least surprise: thread-safe in → thread-safe out. Never downgrade to a plain (non-thread-safe) type like `Hash[T]`.
- Returning a fresh instance keeps the result independent of the receiver's lock.
- The predicate `fn` is a user callback, so it runs **outside** the lock against a snapshot — see [[callback-reentrancy]].

> Gap: existing concurrent `Filter`/etc. currently return the plain type (e.g. `Hash[T]`). These should be migrated to return the same concurrent type.
