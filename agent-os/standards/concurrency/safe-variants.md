# Thread-Safe Collection Variants

Each collection offers two thread-safe variants alongside the plain type:

- **`Concurrent<X>`** — backed by `*sync.Mutex`. Every method locks. Use when reads and writes are balanced.
- **`ConcurrentRW<X>`** — backed by `*sync.RWMutex`. Reads use `RLock`, writes use `Lock`. Use when reads dominate.

```go
type ConcurrentHash[T comparable]   struct { data map[T]struct{}; lock *sync.Mutex }
type ConcurrentHashRW[T comparable] struct { data map[T]struct{}; lock *sync.RWMutex }
```

- Constructor `New<Type>(values ...T)` allocates the lock pointer and seeds data.
- Both variants implement the same `Set`/`Dict`/`List` + `Mutable*` interfaces as the plain type — assert with [[interface-guards]].
- The plain type (`Hash`, `Array`, `Linked`) stays lock-free; thread safety is opt-in by choosing a variant.
