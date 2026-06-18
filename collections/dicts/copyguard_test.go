package dicts

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Dicts(t *testing.T) {
	impl := nocopytest.NoCopyImplementsLocker()
	if !impl {
		t.Error("*nocopy.NoCopy must implement sync.Locker")
	}

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentHash",
			typ:  reflect.TypeOf(ConcurrentHash[string, int]{}),
		},
		{
			name: "ConcurrentHashRW",
			typ:  reflect.TypeOf(ConcurrentHashRW[string, int]{}),
		},
		{
			name: "ConcurrentTree",
			typ:  reflect.TypeOf(ConcurrentTree[string, int]{}),
		},
		{
			name: "ConcurrentTreeRW",
			typ:  reflect.TypeOf(ConcurrentTreeRW[string, int]{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok := nocopytest.HasNoCopyFirstField(tc.typ)
			if !ok {
				t.Errorf("%s: first field is not nocopy.NoCopy", tc.typ)
			}
		})
	}
}
