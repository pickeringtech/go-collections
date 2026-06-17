package collections

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

// The RW facade constructors must return the read-optimised array implementation,
// not just any thread-safe type, and must preserve the seeded values.

func TestNewConcurrentRWList_BackedByRWArray(t *testing.T) {
	list := NewConcurrentRWList(1, 2, 3)
	if _, ok := list.(*lists.ConcurrentRWArray[int]); !ok {
		t.Fatalf("NewConcurrentRWList returned %T, want *lists.ConcurrentRWArray[int]", list)
	}
	if got := list.AsSlice(); len(got) != 3 || got[0] != 1 || got[2] != 3 {
		t.Errorf("AsSlice() = %v, want [1 2 3]", got)
	}
}

func TestNewConcurrentRWQueue_BackedByRWArray(t *testing.T) {
	q := NewConcurrentRWQueue(1, 2, 3)
	if _, ok := q.(*lists.ConcurrentRWArray[int]); !ok {
		t.Fatalf("NewConcurrentRWQueue returned %T, want *lists.ConcurrentRWArray[int]", q)
	}
	if front, ok := q.PeekFront(); !ok || front != 1 {
		t.Errorf("PeekFront() = (%d, %t), want (1, true)", front, ok)
	}
}

func TestNewConcurrentRWStack_BackedByRWArray(t *testing.T) {
	s := NewConcurrentRWStack(1, 2, 3)
	if _, ok := s.(*lists.ConcurrentRWArray[int]); !ok {
		t.Fatalf("NewConcurrentRWStack returned %T, want *lists.ConcurrentRWArray[int]", s)
	}
	if end, ok := s.PeekEnd(); !ok || end != 3 {
		t.Errorf("PeekEnd() = (%d, %t), want (3, true)", end, ok)
	}
}

// The builders must only choose the RW implementation when both Concurrent() and
// RW() are set; RW() alone (without Concurrent()) stays non-concurrent.

func TestListBuilder_RWPath(t *testing.T) {
	rw := NewListBuilder[int]().Concurrent().RW().Add(1, 2, 3).Build()
	if _, ok := rw.(*lists.ConcurrentRWArray[int]); !ok {
		t.Errorf("Concurrent().RW() built %T, want *lists.ConcurrentRWArray[int]", rw)
	}

	rwOnly := NewListBuilder[int]().RW().Add(1, 2, 3).Build()
	if _, ok := rwOnly.(*lists.ConcurrentRWArray[int]); ok {
		t.Errorf("RW() without Concurrent() built an RW array, want a non-concurrent list")
	}

	concurrent := NewListBuilder[int]().Concurrent().Add(1).Build()
	if _, ok := concurrent.(*lists.ConcurrentRWArray[int]); ok {
		t.Errorf("Concurrent() without RW() built an RW array, want a plain concurrent list")
	}
}

func TestQueueBuilder_RWPath(t *testing.T) {
	rw := NewQueueBuilder[int]().Concurrent().RW().Add(1, 2, 3).Build()
	if _, ok := rw.(*lists.ConcurrentRWArray[int]); !ok {
		t.Errorf("Concurrent().RW() built %T, want *lists.ConcurrentRWArray[int]", rw)
	}

	rwOnly := NewQueueBuilder[int]().RW().Add(1).Build()
	if _, ok := rwOnly.(*lists.ConcurrentRWArray[int]); ok {
		t.Errorf("RW() without Concurrent() built an RW array, want a non-concurrent queue")
	}
}

func TestStackBuilder_RWPath(t *testing.T) {
	rw := NewStackBuilder[int]().Concurrent().RW().Add(1, 2, 3).Build()
	if _, ok := rw.(*lists.ConcurrentRWArray[int]); !ok {
		t.Errorf("Concurrent().RW() built %T, want *lists.ConcurrentRWArray[int]", rw)
	}

	rwOnly := NewStackBuilder[int]().RW().Add(1).Build()
	if _, ok := rwOnly.(*lists.ConcurrentRWArray[int]); ok {
		t.Errorf("RW() without Concurrent() built an RW array, want a non-concurrent stack")
	}
}
