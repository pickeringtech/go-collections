// Package nocopytest provides shared assertions for the white-box copy-guard
// tests of the concurrent collection types. It lives in its own package so the
// per-package copyguard_test.go files can reuse a single implementation rather
// than each restating the same reflection boilerplate.
package nocopytest

import (
	"reflect"
	"sync"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// noCopyType is the reflect.Type for nocopy.NoCopy used in field assertions.
var noCopyType = reflect.TypeOf(nocopy.NoCopy{})

// lockerType is the reflect.Type for sync.Locker used in interface assertions.
var lockerType = reflect.TypeOf((*sync.Locker)(nil)).Elem()

// AssertLockerImpl checks that *nocopy.NoCopy implements sync.Locker, the
// interface go vet's copylocks analyser inspects.
func AssertLockerImpl(t *testing.T) {
	t.Helper()
	ptrType := reflect.TypeOf((*nocopy.NoCopy)(nil))
	if !ptrType.Implements(lockerType) {
		t.Error("*nocopy.NoCopy does not implement sync.Locker")
	}
}

// AssertNoCopyFirstField checks that typ embeds nocopy.NoCopy as its first
// field, which is where the copy guard must sit for go vet to flag value-copies
// of the enclosing concurrent type.
func AssertNoCopyFirstField(t *testing.T, typ reflect.Type) {
	t.Helper()
	if typ.NumField() == 0 {
		t.Errorf("%s: has no fields, expected nocopy.NoCopy first", typ.Name())
		return
	}
	if first := typ.Field(0); first.Type != noCopyType {
		t.Errorf("%s: first field is %s, expected nocopy.NoCopy", typ.Name(), first.Type)
	}
}
