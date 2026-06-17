package collections

import (
	"github.com/pickeringtech/go-collections/collections/dicts"
)

// Pair represents a key-value pair, re-exported from the dicts package so the
// facade examples (e.g. collections.NewDict(collections.Pair[K, V]{...})) work
// with a single import "github.com/pickeringtech/go-collections/collections".
type Pair[K comparable, V any] = dicts.Pair[K, V]
