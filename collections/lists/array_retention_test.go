package lists_test

import (
	"runtime"
	"testing"
	"weak"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// retained is large enough to get its own allocation (rather than share a tiny
// allocator block), so a weak pointer to it tracks exactly this object's
// reachability.
type retained struct {
	_ [128]byte
}

// removeWatched builds a two-element Array (so the survivor keeps the backing
// array alive), removes one element and returns the list, a weak pointer to the
// removed element and whether the removal returned it. Doing this in a helper
// leaves no strong reference to the element in the caller's frame.
func removeWatched(
	watch func(l *lists.Array[*retained]) *retained,
	remove func(l *lists.Array[*retained]) (*retained, bool),
) (*lists.Array[*retained], weak.Pointer[retained], bool) {
	l := lists.NewArray(&retained{}, &retained{})
	watched := weak.Make(watch(l))
	removed, ok := remove(l)
	return l, watched, ok && removed == watched.Value()
}

// TestArray_RemovalReleasesElement asserts that an element removed from either
// end of an Array is no longer reachable through the list's backing storage:
// once the caller drops it, a garbage collection must reclaim it while the list
// itself is still alive.
func TestArray_RemovalReleasesElement(t *testing.T) {
	type args struct {
		watch  func(l *lists.Array[*retained]) *retained
		remove func(l *lists.Array[*retained]) (*retained, bool)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "DequeueInPlace releases the removed front element",
			args: args{
				watch: func(l *lists.Array[*retained]) *retained {
					front, _ := l.PeekFront()
					return front
				},
				remove: func(l *lists.Array[*retained]) (*retained, bool) { return l.DequeueInPlace() },
			},
		},
		{
			name: "PopInPlace releases the removed end element",
			args: args{
				watch: func(l *lists.Array[*retained]) *retained {
					end, _ := l.PeekEnd()
					return end
				},
				remove: func(l *lists.Array[*retained]) (*retained, bool) { return l.PopInPlace() },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, watched, removedWatched := removeWatched(tt.args.watch, tt.args.remove)
			if !removedWatched {
				t.Fatalf("removal did not return the watched element")
			}

			runtime.GC()
			if watched.Value() != nil {
				t.Errorf("removed element is still reachable through the list's backing storage")
			}
			if l.Length() != 1 {
				t.Errorf("Length() = %d, want 1", l.Length())
			}
			runtime.KeepAlive(l)
		})
	}
}
