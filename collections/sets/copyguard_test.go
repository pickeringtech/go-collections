package sets

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Sets(t *testing.T) {
	nocopytest.AssertLockerImpl(t)

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentHash",
			typ:  reflect.TypeOf(ConcurrentHash[string]{}),
		},
		{
			name: "ConcurrentHashRW",
			typ:  reflect.TypeOf(ConcurrentHashRW[string]{}),
		},
		{
			name: "ConcurrentTreeSet",
			typ:  reflect.TypeOf(ConcurrentTreeSet[string]{}),
		},
		{
			name: "ConcurrentTreeSetRW",
			typ:  reflect.TypeOf(ConcurrentTreeSetRW[string]{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nocopytest.AssertNoCopyFirstField(t, tc.typ)
		})
	}
}
