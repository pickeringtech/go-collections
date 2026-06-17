// Package nocopytest provides shared predicates for the white-box copy-guard
// tests of the concurrent collection types. It lives in its own package so the
// per-package copyguard_test.go files can reuse a single implementation rather
// than each restating the same reflection boilerplate.
package nocopytest

import (
	"reflect"
	"sync"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// noCopyType is the reflect.Type for nocopy.NoCopy used in field assertions.
var noCopyType = reflect.TypeOf(nocopy.NoCopy{})

// lockerType is the reflect.Type for sync.Locker used in interface assertions.
var lockerType = reflect.TypeOf((*sync.Locker)(nil)).Elem()

// HasNoCopyFirstField reports whether typ embeds nocopy.NoCopy as its first
// field, which is where the copy guard must sit for go vet to flag value-copies
// of the enclosing concurrent type. Returns false if typ has no fields.
func HasNoCopyFirstField(typ reflect.Type) bool {
	if typ.NumField() == 0 {
		return false
	}
	first := typ.Field(0)
	return first.Type == noCopyType
}

// NoCopyImplementsLocker reports whether *nocopy.NoCopy implements sync.Locker,
// the interface inspected by go vet's copylocks analyser.
func NoCopyImplementsLocker() bool {
	ptrType := reflect.TypeOf((*nocopy.NoCopy)(nil))
	return ptrType.Implements(lockerType)
}
