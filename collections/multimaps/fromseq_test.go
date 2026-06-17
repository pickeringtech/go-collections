package multimaps_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/multimaps"
)

func TestListMultimapFromSeq2(t *testing.T) {
	source := multimaps.NewListMultimap(
		multimaps.Entry[string, int]{Key: "a", Value: 1},
		multimaps.Entry[string, int]{Key: "a", Value: 2},
		multimaps.Entry[string, int]{Key: "b", Value: 3},
	)
	got := multimaps.ListMultimapFromSeq2(source.All())

	if !reflect.DeepEqual(got.Get("a"), []int{1, 2}) {
		t.Errorf(`Get("a") = %v, want [1 2]`, got.Get("a"))
	}
	if got.Length() != 3 {
		t.Errorf("Length() = %d, want 3", got.Length())
	}
}

func TestSetMultimapFromSeq2_CollapsesDuplicates(t *testing.T) {
	source := multimaps.NewListMultimap(
		multimaps.Entry[string, int]{Key: "a", Value: 1},
		multimaps.Entry[string, int]{Key: "a", Value: 1},
		multimaps.Entry[string, int]{Key: "a", Value: 2},
	)
	got := multimaps.SetMultimapFromSeq2(source.All())

	values := got.Get("a")
	sort.Ints(values)
	if !reflect.DeepEqual(values, []int{1, 2}) {
		t.Errorf(`Get("a") = %v, want [1 2] (duplicates collapsed)`, values)
	}
}

func TestListMultimapFromSeq2_Empty(t *testing.T) {
	got := multimaps.ListMultimapFromSeq2(multimaps.NewListMultimap[string, int]().All())
	if !got.IsEmpty() {
		t.Errorf("FromSeq2 over empty sequence should be empty")
	}
}

func ExampleListMultimapFromSeq2() {
	source := multimaps.NewListMultimap(
		multimaps.Entry[string, int]{Key: "fruit", Value: 1},
		multimaps.Entry[string, int]{Key: "fruit", Value: 2},
	)
	m := multimaps.ListMultimapFromSeq2(source.All())
	fmt.Println(m.Get("fruit"))
	// Output: [1 2]
}
