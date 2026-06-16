package collections

import "testing"

func TestListBuilder_Plain(t *testing.T) {
	// Non-concurrent path (the Example test covers the concurrent path).
	list := NewListBuilder[int]().Add(1, 2).Add(3).Build()
	if list.Length() != 3 {
		t.Errorf("Length() = %d, want 3", list.Length())
	}
}

func TestQueueBuilder(t *testing.T) {
	tests := []struct {
		name       string
		concurrent bool
	}{
		{"plain", false},
		{"concurrent", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewQueueBuilder[int]().RW().Add(1, 2).Add(3)
			if tt.concurrent {
				b = b.Concurrent()
			}
			q := b.Build()
			front, ok := q.PeekFront()
			if !ok || front != 1 {
				t.Errorf("PeekFront() = (%d, %t), want (1, true)", front, ok)
			}
		})
	}
}

func TestStackBuilder(t *testing.T) {
	tests := []struct {
		name       string
		concurrent bool
	}{
		{"plain", false},
		{"concurrent", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewStackBuilder[int]().RW().Add(1, 2).Add(3)
			if tt.concurrent {
				b = b.Concurrent()
			}
			s := b.Build()
			end, ok := s.PeekEnd()
			if !ok || end != 3 {
				t.Errorf("PeekEnd() = (%d, %t), want (3, true)", end, ok)
			}
		})
	}
}
