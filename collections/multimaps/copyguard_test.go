package multimaps

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Multimaps(t *testing.T) {
	impl := nocopytest.NoCopyImplementsLocker()
	if !impl {
		t.Error("*nocopy.NoCopy must implement sync.Locker")
	}

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentListMultimap",
			typ:  reflect.TypeOf(ConcurrentListMultimap[string, int]{}),
		},
		{
			name: "ConcurrentRWListMultimap",
			typ:  reflect.TypeOf(ConcurrentRWListMultimap[string, int]{}),
		},
		{
			name: "ConcurrentSetMultimap",
			typ:  reflect.TypeOf(ConcurrentSetMultimap[string, string]{}),
		},
		{
			name: "ConcurrentRWSetMultimap",
			typ:  reflect.TypeOf(ConcurrentRWSetMultimap[string, string]{}),
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
