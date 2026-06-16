package dicts

import (
	"iter"

	"github.com/pickeringtech/go-collections/constraints"
)

// Indexable provides basic key-value access operations for dictionaries.
type Indexable[K comparable, V any] interface {
	// Get retrieves the value associated with the given key.
	// If the key is not found, returns the default value and false.
	Get(key K, defaultValue V) (V, bool)

	// Contains checks if the given key exists in the dictionary.
	Contains(key K) bool

	// Length returns the number of key-value pairs in the dictionary.
	Length() int

	// IsEmpty returns true if the dictionary contains no key-value pairs.
	IsEmpty() bool
}

// Iterable provides iteration capabilities for dictionaries.
type Iterable[K comparable, V any] interface {
	// ForEach executes the given function for each key-value pair.
	ForEach(fn func(key K, value V))

	// ForEachKey executes the given function for each key.
	ForEachKey(fn func(key K))

	// ForEachValue executes the given function for each value.
	ForEachValue(fn func(value V))

	// All returns an iterator over key/value pairs, suitable for use with
	// range-over-func. KeysSeq and ValuesSeq are the iterator counterparts of
	// the slice-returning Keys and Values; All is named to avoid colliding with
	// them.
	All() iter.Seq2[K, V]

	// KeysSeq returns an iterator over the keys, suitable for use with
	// range-over-func. It is the iterator counterpart of the slice-returning
	// Keys.
	KeysSeq() iter.Seq[K]

	// ValuesSeq returns an iterator over the values, suitable for use with
	// range-over-func. It is the iterator counterpart of the slice-returning
	// Values.
	ValuesSeq() iter.Seq[V]
}

// Filterable provides filtering capabilities for dictionaries.
type Filterable[K comparable, V any] interface {
	// Filter returns a new dictionary containing only the key-value pairs
	// that satisfy the given predicate function.
	Filter(fn func(key K, value V) bool) Dict[K, V]
}

// MutableFilterable provides in-place filtering capabilities.
type MutableFilterable[K comparable, V any] interface {
	// FilterInPlace removes all key-value pairs that do not satisfy
	// the given predicate function, modifying the dictionary in place.
	FilterInPlace(fn func(key K, value V) bool)
}

// Searchable provides search capabilities for dictionaries.
//
// AllMatch, AnyMatch, NoneMatch and Find form the search core shared across the
// lists, dicts and sets families. The remaining methods (FindKey, FindValue and
// ContainsValue) are deliberate dict-specific extensions that reflect the
// key-value shape of a dictionary.
type Searchable[K comparable, V any] interface {
	// AllMatch returns true if every key-value pair satisfies the given
	// predicate. It is vacuously true for an empty dictionary.
	AllMatch(fn func(key K, value V) bool) bool

	// AnyMatch returns true if at least one key-value pair satisfies the given
	// predicate. It is false for an empty dictionary.
	AnyMatch(fn func(key K, value V) bool) bool

	// NoneMatch returns true if no key-value pair satisfies the given predicate.
	// It is vacuously true for an empty dictionary.
	NoneMatch(fn func(key K, value V) bool) bool

	// Find returns the first key-value pair that satisfies the given predicate.
	// Returns the key, value, and true if found; zero values and false otherwise.
	Find(fn func(key K, value V) bool) (K, V, bool)

	// FindKey returns the first key that satisfies the given predicate.
	// Returns the key and true if found; zero value and false otherwise.
	FindKey(fn func(key K) bool) (K, bool)

	// FindValue returns the first value that satisfies the given predicate.
	// Returns the value and true if found; zero value and false otherwise.
	FindValue(fn func(value V) bool) (V, bool)

	// ContainsValue checks if the given value exists in the dictionary.
	ContainsValue(value V) bool
}

// Convertible provides conversion capabilities for dictionaries.
type Convertible[K comparable, V any] interface {
	// Keys returns a slice containing all keys in the dictionary.
	Keys() []K

	// Values returns a slice containing all values in the dictionary.
	Values() []V

	// Items returns a slice containing all key-value pairs as Pair structs.
	Items() []Pair[K, V]

	// AsMap returns the dictionary as a native Go map.
	AsMap() map[K]V
}

// Insertable provides insertion capabilities for dictionaries.
type Insertable[K comparable, V any] interface {
	// Put creates a new dictionary with the given key-value pair added or updated.
	// Returns the new dictionary without modifying the original.
	Put(key K, value V) Dict[K, V]

	// PutMany creates a new dictionary with all given key-value pairs added or updated.
	// Returns the new dictionary without modifying the original.
	PutMany(pairs ...Pair[K, V]) Dict[K, V]
}

// MutableInsertable provides in-place insertion capabilities.
type MutableInsertable[K comparable, V any] interface {
	// PutInPlace adds or updates the given key-value pair in the dictionary.
	PutInPlace(key K, value V)

	// PutManyInPlace adds or updates all given key-value pairs in the dictionary.
	PutManyInPlace(pairs ...Pair[K, V])
}

// Removable provides removal capabilities for dictionaries.
type Removable[K comparable, V any] interface {
	// Remove creates a new dictionary with the given key removed.
	// Returns the new dictionary without modifying the original.
	Remove(key K) Dict[K, V]

	// RemoveMany creates a new dictionary with all given keys removed.
	// Returns the new dictionary without modifying the original.
	RemoveMany(keys ...K) Dict[K, V]
}

// MutableRemovable provides in-place removal capabilities.
type MutableRemovable[K comparable, V any] interface {
	// RemoveInPlace removes the given key from the dictionary.
	// Returns the removed value and true if the key existed; zero value and false otherwise.
	RemoveInPlace(key K) (V, bool)

	// RemoveManyInPlace removes all given keys from the dictionary.
	RemoveManyInPlace(keys ...K)

	// Clear removes all key-value pairs from the dictionary.
	Clear()
}

// Dict represents an immutable dictionary interface that provides
// comprehensive key-value operations without modifying the original dictionary.
type Dict[K comparable, V any] interface {
	Indexable[K, V]
	Iterable[K, V]
	Filterable[K, V]
	Searchable[K, V]
	Convertible[K, V]
	Insertable[K, V]
	Removable[K, V]
}

// MutableDict represents a mutable dictionary interface that provides
// comprehensive key-value operations with the ability to modify the dictionary in place.
type MutableDict[K comparable, V any] interface {
	Dict[K, V]
	MutableFilterable[K, V]
	MutableInsertable[K, V]
	MutableRemovable[K, V]
}

// Ordered provides sorted-order navigation and iteration for dictionaries whose
// keys are kept in sorted order (e.g. Tree). These are pure read-only queries, so
// — unlike the mutating roles — there is no in-place twin.
type Ordered[K constraints.Ordered, V any] interface {
	// Min returns the entry with the smallest key.
	// Returns the key, value, and true if the dictionary is non-empty;
	// zero values and false otherwise.
	Min() (K, V, bool)

	// Max returns the entry with the largest key.
	// Returns the key, value, and true if the dictionary is non-empty;
	// zero values and false otherwise.
	Max() (K, V, bool)

	// Floor returns the entry with the largest key less than or equal to the
	// given key. Returns the key, value, and true if such an entry exists;
	// zero values and false otherwise.
	Floor(key K) (K, V, bool)

	// Ceiling returns the entry with the smallest key greater than or equal to
	// the given key. Returns the key, value, and true if such an entry exists;
	// zero values and false otherwise.
	Ceiling(key K) (K, V, bool)

	// Range returns all entries whose key is within the inclusive range
	// [lo, hi], in ascending key order. Returns a non-nil (possibly empty) slice.
	Range(lo, hi K) []Pair[K, V]

	// All returns an iterator over all entries in ascending key order.
	All() iter.Seq2[K, V]

	// Backward returns an iterator over all entries in descending key order.
	Backward() iter.Seq2[K, V]

	// RangeAll returns an iterator over the entries whose key is within the
	// inclusive range [lo, hi], in ascending key order.
	RangeAll(lo, hi K) iter.Seq2[K, V]
}

// SortedDict represents an immutable dictionary whose keys are maintained in
// sorted order, adding ordered navigation and iteration on top of the standard
// Dict contract.
type SortedDict[K constraints.Ordered, V any] interface {
	Dict[K, V]
	Ordered[K, V]
}

// MutableSortedDict represents a mutable dictionary whose keys are maintained in
// sorted order, combining the SortedDict contract with in-place mutation.
type MutableSortedDict[K constraints.Ordered, V any] interface {
	SortedDict[K, V]
	MutableFilterable[K, V]
	MutableInsertable[K, V]
	MutableRemovable[K, V]
}
