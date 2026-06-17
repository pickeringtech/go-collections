package nocopytest

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/internal/nocopy"
)

// withNoCopy has nocopy.NoCopy as its first field, so HasNoCopyFirstField must return true.
type withNoCopy struct {
	_ nocopy.NoCopy
	x int
}

// wrongFirst has a non-NoCopy first field, so HasNoCopyFirstField must return false.
type wrongFirst struct {
	x int
}

func TestHasNoCopyFirstField(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want bool
	}{
		{
			name: "NoCopy first field",
			typ:  reflect.TypeOf(withNoCopy{}),
			want: true,
		},
		{
			name: "wrong first field",
			typ:  reflect.TypeOf(wrongFirst{}),
			want: false,
		},
		{
			name: "empty struct",
			typ:  reflect.TypeOf(struct{}{}),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasNoCopyFirstField(tc.typ)
			if got != tc.want {
				t.Errorf("HasNoCopyFirstField(%s) = %v, want %v", tc.typ, got, tc.want)
			}
		})
	}
}

func TestNoCopyImplementsLocker(t *testing.T) {
	impl := NoCopyImplementsLocker()
	if !impl {
		t.Error("NoCopyImplementsLocker() = false, want true")
	}
}
