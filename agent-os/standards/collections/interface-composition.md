# Capability Interface Composition

Collection contracts are built from small, single-responsibility role interfaces composed into one aggregate interface. Never declare a fat interface directly.

Role interfaces (one capability each):

```go
type Indexable[T comparable] interface { Contains(element T) bool; Length() int; IsEmpty() bool }
type Iterable[T comparable]  interface { ForEach(fn func(element T)) }
type Filterable[T comparable] interface { Filter(fn func(element T) bool) Set[T] }
type Searchable[T comparable] interface { Find(...); AllMatch(...); AnyMatch(...) }
type Convertible[T comparable] interface { AsSlice() []T; AsMap() map[T]struct{} }
// ... Insertable, Removable, SetOperations
```

Aggregate interface = composition of roles:

```go
type Set[T comparable] interface {
	Indexable[T]
	Iterable[T]
	Filterable[T]
	Searchable[T]
	Convertible[T]
	Insertable[T]
	Removable[T]
	SetOperations[T]
}
```

- Common roles across collections: `Indexable`, `Iterable`, `Filterable`, `Searchable`, `Convertible`, `Insertable`, `Removable`.
- Every interface and method carries a doc comment, including the `(value, true)` / `(zero, false)` return convention for lookups.
- See [[mutable-immutable-hierarchy]] for how mutable variants compose, and [[interface-guards]] for conformance assertions.
