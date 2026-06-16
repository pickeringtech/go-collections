# Mutable / Immutable Dual Hierarchy

Every collection exposes two interfaces: an immutable base and a mutable extension.

- **Base interface** (`Dict`, `Set`, `List`) — methods return a **new** collection; the receiver is never modified.
- **`Mutable*` interface** — embeds the base, then adds the in-place (`Mutable*`) capability roles.

```go
type Set[T comparable] interface {
	Indexable[T]; Iterable[T]; Filterable[T]; /* ... */ Insertable[T]; Removable[T]
}

type MutableSet[T comparable] interface {
	Set[T]                  // everything the immutable set can do
	MutableFilterable[T]    // + FilterInPlace
	MutableInsertable[T]    // + AddInPlace, AddManyInPlace
	MutableRemovable[T]     // + RemoveInPlace, Clear
}
```

## Every mutating capability ships in both forms

For each role that changes contents, provide a pair:

| Form | Interface | Method | Returns |
|------|-----------|--------|---------|
| Immutable | `Insertable` | `Add(e) Set[T]` | new collection |
| In-place | `MutableInsertable` | `AddInPlace(e)` | nothing (mutates) |

- The returns-new version lives on the base; the `InPlace` version lives on the matching `Mutable*` role. See [[inplace-suffix]].
- **Exception:** pure read/query operations (`Contains`, `IsSubsetOf`, `AllMatch`) don't mutate and have no in-place form — they belong only on the base.
