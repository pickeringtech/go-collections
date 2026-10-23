# Thread-Safe Collection Variants

Each collection offers two thread-safe variants alongside the plain type:

- **`Concurrent<X>`** — guarded by a `sync.Mutex`. Every method locks. Use when reads and writes are balanced.
- **`ConcurrentRW<X>`** — guarded by a `sync.RWMutex`. Reads use `RLock`, writes use `Lock`. Use when reads dominate.

```go
type ConcurrentHash[T comparable] struct {
	_    nocopy.NoCopy // go vet copylocks guard (collections/internal/nocopy)
	data map[T]struct{}
	lock sync.Mutex // a value, never a pointer
}

type ConcurrentHashRW[T comparable] struct {
	_    nocopy.NoCopy
	data map[T]struct{}
	lock sync.RWMutex
}
```

- The mutex is held **by value**, so the lock of even a bare `&ConcurrentHash{}` is safe to take (no nil-pointer panic on first use, #124). The `_ nocopy.NoCopy` blank field makes `go vet` report any copy of the struct after first use, since a copy would carry an independent lock over shared backing data (#152). All methods use pointer receivers.
- Constructor `New<Type>(values ...T)` seeds the data; it does not need to initialise the lock. Construct via `New<Type>` — the zero value's backing data is usually nil — and document the type's actual zero-value behaviour on the type.
- Both variants implement the same `Set`/`Dict`/`List` + `Mutable*` interfaces as the plain type — assert with [[interface-guards]].
- The plain type (`Hash`, `Array`, `Linked`) stays lock-free; thread safety is opt-in by choosing a variant.
- A variant whose `FilterInPlace` must tell "nothing wrote while the predicate ran" apart from "a write raced it" also carries a `version uint64` write counter — see [[callback-reentrancy]].
