package collections

import (
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/collections/multimaps"
)

// listConstructor names a List-returning constructor so the same assertions can
// run against every backing implementation.
type listConstructor struct {
	name string
	make func(...int) interface {
		AsSlice() []int
		Length() int
	}
}

func TestListConstructors(t *testing.T) {
	constructors := []listConstructor{
		{"NewList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewList(v...)
		}},
		{"NewConcurrentList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewConcurrentList(v...)
		}},
		{"NewLinkedList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewLinkedList(v...)
		}},
		{"NewConcurrentLinkedList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewConcurrentLinkedList(v...)
		}},
		{"NewConcurrentRWLinkedList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewConcurrentRWLinkedList(v...)
		}},
		{"NewDoublyLinkedList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewDoublyLinkedList(v...)
		}},
		{"NewConcurrentDoublyLinkedList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewConcurrentDoublyLinkedList(v...)
		}},
		{"NewConcurrentRWDoublyLinkedList", func(v ...int) interface {
			AsSlice() []int
			Length() int
		} {
			return NewConcurrentRWDoublyLinkedList(v...)
		}},
	}

	for _, c := range constructors {
		t.Run(c.name, func(t *testing.T) {
			list := c.make(1, 2, 3)
			if got := list.Length(); got != 3 {
				t.Errorf("Length() = %d, want 3", got)
			}
			got := list.AsSlice()
			want := []int{1, 2, 3}
			if len(got) != len(want) {
				t.Fatalf("AsSlice() = %v, want %v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("AsSlice()[%d] = %d, want %d", i, got[i], want[i])
				}
			}
		})
	}
}

func TestQueueConstructors(t *testing.T) {
	queues := map[string]func(...int) interface {
		PeekFront() (int, bool)
	}{
		"NewQueue":           func(v ...int) interface{ PeekFront() (int, bool) } { return NewQueue(v...) },
		"NewConcurrentQueue": func(v ...int) interface{ PeekFront() (int, bool) } { return NewConcurrentQueue(v...) },
	}
	for name, make := range queues {
		t.Run(name, func(t *testing.T) {
			q := make(10, 20, 30)
			front, ok := q.PeekFront()
			if !ok || front != 10 {
				t.Errorf("PeekFront() = (%d, %t), want (10, true)", front, ok)
			}
		})
	}
}

func TestStackConstructors(t *testing.T) {
	stacks := map[string]func(...int) interface {
		PeekEnd() (int, bool)
	}{
		"NewStack":           func(v ...int) interface{ PeekEnd() (int, bool) } { return NewStack(v...) },
		"NewConcurrentStack": func(v ...int) interface{ PeekEnd() (int, bool) } { return NewConcurrentStack(v...) },
	}
	for name, make := range stacks {
		t.Run(name, func(t *testing.T) {
			s := make(10, 20, 30)
			end, ok := s.PeekEnd()
			if !ok || end != 30 {
				t.Errorf("PeekEnd() = (%d, %t), want (30, true)", end, ok)
			}
		})
	}
}

func TestDictConstructors(t *testing.T) {
	dictFns := map[string]func(...dicts.Pair[string, int]) dicts.Dict[string, int]{
		"NewDict":             NewDict[string, int],
		"NewConcurrentDict":   NewConcurrentDict[string, int],
		"NewConcurrentRWDict": NewConcurrentRWDict[string, int],
	}
	for name, make := range dictFns {
		t.Run(name, func(t *testing.T) {
			d := make(
				dicts.Pair[string, int]{Key: "a", Value: 1},
				dicts.Pair[string, int]{Key: "b", Value: 2},
			)
			if d.Length() != 2 {
				t.Errorf("Length() = %d, want 2", d.Length())
			}
			value, ok := d.Get("a", -1)
			if !ok || value != 1 {
				t.Errorf("Get(a) = (%d, %t), want (1, true)", value, ok)
			}
			_, ok = d.Get("missing", -1)
			if ok {
				t.Errorf("Get(missing) ok = true, want false")
			}
		})
	}
}

func TestMultimapConstructors(t *testing.T) {
	multimapFns := map[string]func(...multimaps.Entry[string, int]) multimaps.Multimap[string, int]{
		"NewListMultimap":             NewListMultimap[string, int],
		"NewConcurrentListMultimap":   NewConcurrentListMultimap[string, int],
		"NewConcurrentRWListMultimap": NewConcurrentRWListMultimap[string, int],
		"NewSetMultimap":              NewSetMultimap[string, int],
		"NewConcurrentSetMultimap":    NewConcurrentSetMultimap[string, int],
		"NewConcurrentRWSetMultimap":  NewConcurrentRWSetMultimap[string, int],
	}
	for name, make := range multimapFns {
		t.Run(name, func(t *testing.T) {
			m := make(
				multimaps.Entry[string, int]{Key: "a", Value: 1},
				multimaps.Entry[string, int]{Key: "a", Value: 2},
				multimaps.Entry[string, int]{Key: "b", Value: 3},
			)
			if m.KeyCount() != 2 {
				t.Errorf("KeyCount() = %d, want 2", m.KeyCount())
			}
			if !m.ContainsEntry("a", 1) {
				t.Errorf("ContainsEntry(a, 1) = false, want true")
			}
			if m.ContainsEntry("a", 99) {
				t.Errorf("ContainsEntry(a, 99) = true, want false")
			}
		})
	}
}

func TestSetConstructors(t *testing.T) {
	setFns := map[string]func(...int) interface {
		Contains(int) bool
		Length() int
	}{
		"NewSet": func(v ...int) interface {
			Contains(int) bool
			Length() int
		} {
			return NewSet(v...)
		},
		"NewConcurrentSet": func(v ...int) interface {
			Contains(int) bool
			Length() int
		} {
			return NewConcurrentSet(v...)
		},
		"NewConcurrentRWSet": func(v ...int) interface {
			Contains(int) bool
			Length() int
		} {
			return NewConcurrentRWSet(v...)
		},
	}
	for name, make := range setFns {
		t.Run(name, func(t *testing.T) {
			s := make(1, 2, 3)
			if s.Length() != 3 {
				t.Errorf("Length() = %d, want 3", s.Length())
			}
			if !s.Contains(2) {
				t.Errorf("Contains(2) = false, want true")
			}
			if s.Contains(99) {
				t.Errorf("Contains(99) = true, want false")
			}
		})
	}
}
