package lists_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

func TestList_All(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			list := ctor.make(10, 20, 30)

			var gotIdx []int
			var gotVal []int
			for i, v := range list.All() {
				gotIdx = append(gotIdx, i)
				gotVal = append(gotVal, v)
			}

			if want := []int{0, 1, 2}; !reflect.DeepEqual(gotIdx, want) {
				t.Errorf("All() indices = %v, want %v", gotIdx, want)
			}
			if want := []int{10, 20, 30}; !reflect.DeepEqual(gotVal, want) {
				t.Errorf("All() values = %v, want %v", gotVal, want)
			}
		})
	}
}

func TestList_All_Empty(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			count := 0
			for range ctor.make().All() {
				count++
			}
			if count != 0 {
				t.Errorf("All() over empty list yielded %d times, want 0", count)
			}
		})
	}
}

func TestList_Values(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for v := range ctor.make(1, 2, 3).Values() {
				got = append(got, v)
			}
			if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
				t.Errorf("Values() = %v, want %v", got, want)
			}
		})
	}
}

func TestList_Backward(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var gotIdx []int
			var gotVal []int
			for i, v := range ctor.make(10, 20, 30).Backward() {
				gotIdx = append(gotIdx, i)
				gotVal = append(gotVal, v)
			}
			if want := []int{2, 1, 0}; !reflect.DeepEqual(gotIdx, want) {
				t.Errorf("Backward() indices = %v, want %v", gotIdx, want)
			}
			if want := []int{30, 20, 10}; !reflect.DeepEqual(gotVal, want) {
				t.Errorf("Backward() values = %v, want %v", gotVal, want)
			}
		})
	}
}

func TestList_All_EarlyBreak(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for _, v := range ctor.make(1, 2, 3, 4, 5).All() {
				got = append(got, v)
				if v == 2 {
					break
				}
			}
			if want := []int{1, 2}; !reflect.DeepEqual(got, want) {
				t.Errorf("All() with early break = %v, want %v", got, want)
			}
		})
	}
}

func TestList_Values_EarlyBreak(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for v := range ctor.make(1, 2, 3, 4, 5).Values() {
				got = append(got, v)
				if v == 3 {
					break
				}
			}
			if want := []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
				t.Errorf("Values() with early break = %v, want %v", got, want)
			}
		})
	}
}

func TestList_Backward_EarlyBreak(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			var got []int
			for _, v := range ctor.make(1, 2, 3, 4, 5).Backward() {
				got = append(got, v)
				if v == 4 {
					break
				}
			}
			if want := []int{5, 4}; !reflect.DeepEqual(got, want) {
				t.Errorf("Backward() with early break = %v, want %v", got, want)
			}
		})
	}
}

func TestFromSeq(t *testing.T) {
	for _, ctor := range allMutableListConstructors() {
		t.Run(ctor.name, func(t *testing.T) {
			source := ctor.make(7, 8, 9)
			got := lists.FromSeq(source.Values())
			if want := []int{7, 8, 9}; !reflect.DeepEqual(got.AsSlice(), want) {
				t.Errorf("FromSeq round-trip = %v, want %v", got.AsSlice(), want)
			}
		})
	}
}

func TestFromSeq_Empty(t *testing.T) {
	got := lists.FromSeq(lists.NewArray[int]().Values())
	if !got.IsEmpty() {
		t.Errorf("FromSeq over empty sequence should be empty, got %v", got.AsSlice())
	}
}

func ExampleArray_All() {
	list := lists.NewArray("a", "b", "c")
	for i, v := range list.All() {
		fmt.Printf("%d=%s ", i, v)
	}
	// Output: 0=a 1=b 2=c
}

func ExampleArray_Backward() {
	list := lists.NewArray("a", "b", "c")
	for i, v := range list.Backward() {
		fmt.Printf("%d=%s ", i, v)
	}
	// Output: 2=c 1=b 0=a
}

func ExampleFromSeq() {
	source := lists.NewArray(1, 2, 3)
	list := lists.FromSeq(source.Values())
	fmt.Println(list.AsSlice())
	// Output: [1 2 3]
}
