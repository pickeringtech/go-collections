package collections_test

import (
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections"
	"github.com/pickeringtech/go-collections/collections/lists"
)

// queueContents reads a Queue front to back through its immutable Dequeue, so
// the queue under test is left untouched.
func queueContents(q lists.Queue[int]) []int {
	contents := []int{}
	for {
		value, ok, rest := q.Dequeue()
		if !ok {
			return contents
		}
		contents = append(contents, value)
		q = rest
	}
}

// stackContents reads a Stack bottom to top through its immutable Pop, so the
// stack under test is left untouched.
func stackContents(s lists.Stack[int]) []int {
	contents := []int{}
	for {
		value, ok, rest := s.Pop()
		if !ok {
			return contents
		}
		contents = append([]int{value}, contents...)
		s = rest
	}
}

// branchedContents is what TestBuilders_BranchesAreIndependent reads back from
// a base builder and two builders branched from it.
type branchedContents struct {
	base, x, y []int
}

// TestBuilders_BranchesAreIndependent branches two builders from one base and
// asserts each builds only its own values. The base is grown by two Add calls
// so its backing slice has spare capacity, which is what let sibling branches
// overwrite each other's values when Add appended in place.
func TestBuilders_BranchesAreIndependent(t *testing.T) {
	tests := []struct {
		name  string
		build func() branchedContents
	}{
		{
			name: "ListBuilder",
			build: func() branchedContents {
				base := collections.NewListBuilder[int]().Add(1, 2).Add(3)
				x, y := base.Add(4), base.Add(5)
				return branchedContents{base.Build().AsSlice(), x.Build().AsSlice(), y.Build().AsSlice()}
			},
		},
		{
			name: "QueueBuilder",
			build: func() branchedContents {
				base := collections.NewQueueBuilder[int]().Add(1, 2).Add(3)
				x, y := base.Add(4), base.Add(5)
				return branchedContents{queueContents(base.Build()), queueContents(x.Build()), queueContents(y.Build())}
			},
		},
		{
			name: "StackBuilder",
			build: func() branchedContents {
				base := collections.NewStackBuilder[int]().Add(1, 2).Add(3)
				x, y := base.Add(4), base.Add(5)
				return branchedContents{stackContents(base.Build()), stackContents(x.Build()), stackContents(y.Build())}
			},
		},
		{
			name: "DequeBuilder",
			build: func() branchedContents {
				base := collections.NewDequeBuilder[int]().Add(1, 2).Add(3)
				x, y := base.Add(4), base.Add(5)
				return branchedContents{base.Build().AsSlice(), x.Build().AsSlice(), y.Build().AsSlice()}
			},
		},
	}
	want := branchedContents{base: []int{1, 2, 3}, x: []int{1, 2, 3, 4}, y: []int{1, 2, 3, 5}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.build()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("built (base, x, y) = %v, want %v", got, want)
			}
		})
	}
}
