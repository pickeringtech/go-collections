package collections

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/deques"
)

func ExampleDequeBuilder_Build() {
	d := NewDequeBuilder[int]().
		Add(1, 2, 3).
		Add(4, 5).
		Build()

	fmt.Printf("deque: %v\n", d.AsSlice())
	// Output: deque: [1 2 3 4 5]
}

func ExampleDequeBuilder_Bounded() {
	d := NewDequeBuilder[int]().
		Bounded(3, deques.OverwriteOldest).
		Add(1, 2, 3, 4).
		Build()

	fmt.Printf("deque: %v cap: %d\n", d.AsSlice(), d.Capacity())
	// Output: deque: [2 3 4] cap: 3
}

func TestDequeBuilderVariants(t *testing.T) {
	tests := []struct {
		name    string
		build   func() deques.Deque[int]
		wantTyp interface{}
	}{
		{
			name:    "plain unbounded",
			build:   func() deques.Deque[int] { return NewDequeBuilder[int]().Add(1, 2).Build() },
			wantTyp: &deques.RingBuffer[int]{},
		},
		{
			name:    "concurrent unbounded",
			build:   func() deques.Deque[int] { return NewDequeBuilder[int]().Concurrent().Add(1, 2).Build() },
			wantTyp: &deques.ConcurrentRingBuffer[int]{},
		},
		{
			name:    "concurrent rw unbounded",
			build:   func() deques.Deque[int] { return NewDequeBuilder[int]().Concurrent().RW().Add(1, 2).Build() },
			wantTyp: &deques.ConcurrentRWRingBuffer[int]{},
		},
		{
			name: "plain bounded",
			build: func() deques.Deque[int] {
				return NewDequeBuilder[int]().Bounded(5, deques.RejectWhenFull).Add(1, 2).Build()
			},
			wantTyp: &deques.RingBuffer[int]{},
		},
		{
			name: "concurrent bounded",
			build: func() deques.Deque[int] {
				return NewDequeBuilder[int]().Concurrent().Bounded(5, deques.RejectWhenFull).Add(1, 2).Build()
			},
			wantTyp: &deques.ConcurrentRingBuffer[int]{},
		},
		{
			name: "concurrent rw bounded",
			build: func() deques.Deque[int] {
				return NewDequeBuilder[int]().Concurrent().RW().Bounded(5, deques.RejectWhenFull).Add(1, 2).Build()
			},
			wantTyp: &deques.ConcurrentRWRingBuffer[int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.build()
			gotType := reflect.TypeOf(d)
			wantType := reflect.TypeOf(tt.wantTyp)
			if gotType != wantType {
				t.Errorf("Build() type = %v, want %v", gotType, wantType)
			}
			if !reflect.DeepEqual(d.AsSlice(), []int{1, 2}) {
				t.Errorf("Build() contents = %v, want [1 2]", d.AsSlice())
			}
		})
	}
}

func TestNewDequeWrappers(t *testing.T) {
	tests := []struct {
		name    string
		got     deques.Deque[int]
		wantTyp interface{}
		want    []int
	}{
		{"NewDeque", NewDeque(1, 2, 3), &deques.RingBuffer[int]{}, []int{1, 2, 3}},
		{"NewConcurrentDeque", NewConcurrentDeque(1, 2, 3), &deques.ConcurrentRingBuffer[int]{}, []int{1, 2, 3}},
		{"NewConcurrentRWDeque", NewConcurrentRWDeque(1, 2, 3), &deques.ConcurrentRWRingBuffer[int]{}, []int{1, 2, 3}},
		{"NewBoundedDeque", NewBoundedDeque(2, deques.OverwriteOldest, 1, 2, 3), &deques.RingBuffer[int]{}, []int{2, 3}},
		{"NewBoundedConcurrentDeque", NewBoundedConcurrentDeque(2, deques.OverwriteOldest, 1, 2, 3), &deques.ConcurrentRingBuffer[int]{}, []int{2, 3}},
		{"NewBoundedConcurrentRWDeque", NewBoundedConcurrentRWDeque(2, deques.RejectWhenFull, 1, 2, 3), &deques.ConcurrentRWRingBuffer[int]{}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := reflect.TypeOf(tt.got)
			wantType := reflect.TypeOf(tt.wantTyp)
			if gotType != wantType {
				t.Errorf("%s type = %v, want %v", tt.name, gotType, wantType)
			}
			if !reflect.DeepEqual(tt.got.AsSlice(), tt.want) {
				t.Errorf("%s contents = %v, want %v", tt.name, tt.got.AsSlice(), tt.want)
			}
		})
	}
}
