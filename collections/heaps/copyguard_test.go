package heaps

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Heaps(t *testing.T) {
	impl := nocopytest.NoCopyImplementsLocker()
	if !impl {
		t.Error("*nocopy.NoCopy must implement sync.Locker")
	}

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentBinary",
			typ:  reflect.TypeOf(ConcurrentBinary[int]{}),
		},
		{
			name: "ConcurrentRWBinary",
			typ:  reflect.TypeOf(ConcurrentRWBinary[int]{}),
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
