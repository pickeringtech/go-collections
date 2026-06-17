package lists

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Lists(t *testing.T) {
	nocopytest.AssertLockerImpl(t)

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentArray",
			typ:  reflect.TypeOf(ConcurrentArray[int]{}),
		},
		{
			name: "ConcurrentRWArray",
			typ:  reflect.TypeOf(ConcurrentRWArray[int]{}),
		},
		{
			name: "ConcurrentLinked",
			typ:  reflect.TypeOf(ConcurrentLinked[int]{}),
		},
		{
			name: "ConcurrentRWLinked",
			typ:  reflect.TypeOf(ConcurrentRWLinked[int]{}),
		},
		{
			name: "ConcurrentDoublyLinked",
			typ:  reflect.TypeOf(ConcurrentDoublyLinked[int]{}),
		},
		{
			name: "ConcurrentRWDoublyLinked",
			typ:  reflect.TypeOf(ConcurrentRWDoublyLinked[int]{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nocopytest.AssertNoCopyFirstField(t, tc.typ)
		})
	}
}
