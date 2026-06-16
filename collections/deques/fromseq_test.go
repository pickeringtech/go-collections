package deques_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/deques"
)

func TestFromSeq(t *testing.T) {
	source := deques.NewRingBuffer(1, 2, 3)
	got := deques.FromSeq(source.Values())
	if want := []int{1, 2, 3}; !reflect.DeepEqual(got.AsSlice(), want) {
		t.Errorf("FromSeq round-trip = %v, want %v", got.AsSlice(), want)
	}
}

func TestFromSeq_Empty(t *testing.T) {
	got := deques.FromSeq(deques.NewRingBuffer[int]().Values())
	if !got.IsEmpty() {
		t.Errorf("FromSeq over empty sequence should be empty, got %v", got.AsSlice())
	}
}

func ExampleFromSeq() {
	source := deques.NewRingBuffer(1, 2, 3)
	d := deques.FromSeq(source.Values())
	fmt.Println(d.AsSlice())
	// Output: [1 2 3]
}
