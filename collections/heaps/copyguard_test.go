package heaps

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Heaps(t *testing.T) {
	nocopytest.AssertLockerImpl(t)

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
			nocopytest.AssertNoCopyFirstField(t, tc.typ)
		})
	}
}
