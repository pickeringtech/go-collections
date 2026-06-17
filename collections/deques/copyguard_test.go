package deques

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy/nocopytest"
)

func TestCopyGuard_Deques(t *testing.T) {
	nocopytest.AssertLockerImpl(t)

	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{
			name: "ConcurrentRingBuffer",
			typ:  reflect.TypeOf(ConcurrentRingBuffer[int]{}),
		},
		{
			name: "ConcurrentRWRingBuffer",
			typ:  reflect.TypeOf(ConcurrentRWRingBuffer[int]{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nocopytest.AssertNoCopyFirstField(t, tc.typ)
		})
	}
}
