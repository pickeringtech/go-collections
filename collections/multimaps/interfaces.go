package multimaps

import "iter"

// Indexable provides basic key-to-many-values access operations for multimaps.
type Indexable[K comparable, V any] interface {
	// Get returns a copy of all values currently bound to the given key.
	// If the key has no values, an empty (non-nil) slice is returned. The
	// returned slice is independent of the multimap; mutating it does not
	// affect the receiver.
	Get(key K) []V

	// ContainsKey reports whether the given key has at least one value bound
	// to it.
	ContainsKey(key K) bool

	// ContainsEntry reports whether the given key is bound to the given value.
	// Value equality follows the backing: set-backed multimaps require V to be
	// comparable and use native equality; list-backed multimaps accept any V and
	// compare with reflect.DeepEqual.
	ContainsEntry(key K, value V) bool

	// Length returns the total number of entries (key-value associations) in
	// the multimap. A key bound to three values contributes three to the
	// length. Use KeyCount for the number of distinct keys.
	Length() int

	// KeyCount returns the number of distinct keys in the multimap. Use
	// Length for the total number of key-value associations.
	KeyCount() int

	// IsEmpty returns true if the multimap contains no entries.
	IsEmpty() bool
}

// Iterable provides iteration capabilities for multimaps.
type Iterable[K comparable, V any] interface {
	// ForEach executes the given function once for every entry (key-value
	// association), visiting each key once per value bound to it.
	ForEach(fn func(key K, value V))

	// ForEachKey executes the given function once per distinct key, passing a
	// copy of all values bound to that key.
	ForEachKey(fn func(key K, values []V))

	// All returns an iterator over every entry (key-value association),
	// suitable for use with range-over-func.
	All() iter.Seq2[K, V]

	// KeysSeq returns an iterator over the distinct keys, suitable for use
	// with range-over-func.
	KeysSeq() iter.Seq[K]
}

// Filterable provides filtering capabilities for multimaps.
type Filterable[K comparable, V any] interface {
	// Filter returns a new multimap containing only the entries that satisfy
	// the given predicate. Keys left with no values are dropped. The original
	// is not modified.
	Filter(fn func(key K, value V) bool) Multimap[K, V]
}

// MutableFilterable provides in-place filtering capabilities.
type MutableFilterable[K comparable, V any] interface {
	// FilterInPlace removes every entry that does not satisfy the given
	// predicate, modifying the multimap in place. Keys left with no values are
	// dropped.
	FilterInPlace(fn func(key K, value V) bool)
}

// Searchable provides search capabilities for multimaps.
//
// AllMatch, AnyMatch, NoneMatch and Find form the search core shared across the
// lists, dicts and sets families, here operating over individual entries.
type Searchable[K comparable, V any] interface {
	// AllMatch returns true if every entry satisfies the given predicate. It is
	// vacuously true for an empty multimap.
	AllMatch(fn func(key K, value V) bool) bool

	// AnyMatch returns true if at least one entry satisfies the given
	// predicate. It is false for an empty multimap.
	AnyMatch(fn func(key K, value V) bool) bool

	// NoneMatch returns true if no entry satisfies the given predicate. It is
	// vacuously true for an empty multimap.
	NoneMatch(fn func(key K, value V) bool) bool

	// Find returns the first entry that satisfies the given predicate. Returns
	// the key, value, and true if found; zero values and false otherwise.
	// Iteration order over keys is unspecified.
	Find(fn func(key K, value V) bool) (K, V, bool)
}

// Convertible provides conversion capabilities for multimaps.
type Convertible[K comparable, V any] interface {
	// Keys returns a slice containing each distinct key once. Order is
	// unspecified.
	Keys() []K

	// Values returns a slice containing every value across all keys. A value
	// bound to multiple keys (or bound multiple times to one list-backed key)
	// appears once per binding. Order is unspecified.
	Values() []V

	// Entries returns a slice containing every entry (key-value association).
	// Order is unspecified.
	Entries() []Entry[K, V]

	// AsMap returns the multimap as a native Go map from each key to a copy of
	// its values. The returned map is independent of the multimap.
	AsMap() map[K][]V
}

// Insertable provides insertion capabilities for multimaps.
type Insertable[K comparable, V any] interface {
	// Put returns a new multimap with the given value bound to the given key,
	// in addition to any existing values. The original is not modified.
	Put(key K, value V) Multimap[K, V]

	// PutAll returns a new multimap with all the given values bound to the
	// given key, in addition to any existing values. The original is not
	// modified.
	PutAll(key K, values ...V) Multimap[K, V]
}

// MutableInsertable provides in-place insertion capabilities.
type MutableInsertable[K comparable, V any] interface {
	// PutInPlace binds the given value to the given key, in addition to any
	// existing values, modifying the multimap in place.
	PutInPlace(key K, value V)

	// PutAllInPlace binds all the given values to the given key, in addition to
	// any existing values, modifying the multimap in place.
	PutAllInPlace(key K, values ...V)
}

// Removable provides removal capabilities for multimaps.
type Removable[K comparable, V any] interface {
	// Remove returns a new multimap with a single binding of the given value to
	// the given key removed. For list-backed multimaps the first occurrence is
	// removed; set-backed multimaps hold at most one. The original is not
	// modified.
	Remove(key K, value V) Multimap[K, V]

	// RemoveAll returns a new multimap with the given key and all of its values
	// removed. The original is not modified.
	RemoveAll(key K) Multimap[K, V]
}

// MutableRemovable provides in-place removal capabilities.
type MutableRemovable[K comparable, V any] interface {
	// RemoveInPlace removes a single binding of the given value to the given
	// key, modifying the multimap in place. Returns true if a binding was
	// present and removed; false otherwise.
	RemoveInPlace(key K, value V) bool

	// RemoveAllInPlace removes the given key and all of its values, modifying
	// the multimap in place. Returns the removed values and true if the key was
	// present; an empty slice and false otherwise.
	RemoveAllInPlace(key K) ([]V, bool)

	// Clear removes all entries from the multimap.
	Clear()
}

// Multimap represents an immutable multimap interface: one key maps to many
// values, with operations that return a new multimap rather than modifying the
// receiver.
type Multimap[K comparable, V any] interface {
	Indexable[K, V]
	Iterable[K, V]
	Filterable[K, V]
	Searchable[K, V]
	Convertible[K, V]
	Insertable[K, V]
	Removable[K, V]
}

// MutableMultimap represents a mutable multimap interface, extending Multimap
// with in-place operations that modify the receiver.
type MutableMultimap[K comparable, V any] interface {
	Multimap[K, V]
	MutableFilterable[K, V]
	MutableInsertable[K, V]
	MutableRemovable[K, V]
}
