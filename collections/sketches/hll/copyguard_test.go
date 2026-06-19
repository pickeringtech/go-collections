package hll

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_HLL(t *testing.T) {
	if !nocopytest.NoCopyImplementsLocker() {
		t.Error("*nocopy.NoCopy must implement sync.Locker")
	}

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentSketch",
			typ:  reflect.TypeOf(ConcurrentSketch[int]{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !nocopytest.HasNoCopyFirstField(tc.typ) {
				t.Errorf("%s: first field is not nocopy.NoCopy", tc.typ)
			}
		})
	}
}
