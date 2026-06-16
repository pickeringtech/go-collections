# InPlace Suffix Convention

A method's name signals whether it mutates the receiver.

- **No suffix** → returns a new collection, receiver untouched.
- **`InPlace` suffix** → mutates the receiver, returns nothing (or just a status bool).

```go
// Immutable — original unchanged, new dict returned
Put(key K, value V) Dict[K, V]
Remove(key K) Dict[K, V]
Filter(fn func(K, V) bool) Dict[K, V]

// In-place — mutates receiver
PutInPlace(key K, value V)
RemoveInPlace(key K) (V, bool)
FilterInPlace(fn func(K, V) bool)
```

- The `...Many` plural keeps the suffix last: `PutManyInPlace`, `RemoveManyInPlace`.
- `Clear()` is inherently mutating and lives on the `Mutable*` role (no immutable twin).
- Pair every mutating immutable method with its `InPlace` sibling — see [[mutable-immutable-hierarchy]].
