package lists_test

import (
	"github.com/pickeringtech/go-collections/collections/lists"
	"testing"
)

func TestLinked_Basic(t *testing.T) {
	l := lists.NewLinked(1, 2, 3)
	
	if l.Length() != 3 {
		t.Errorf("Expected length 3, got %d", l.Length())
	}
	
	value := l.Get(1, -1)
	if value != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", value)
	}
}

func TestDoublyLinked_Basic(t *testing.T) {
	dl := lists.NewDoublyLinked(1, 2, 3)
	
	if dl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", dl.Length())
	}
	
	value := dl.Get(1, -1)
	if value != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", value)
	}
}

func TestConcurrentLinked_Basic(t *testing.T) {
	cl := lists.NewConcurrentLinked(1, 2, 3)
	
	if cl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", cl.Length())
	}
	
	value := cl.Get(1, -1)
	if value != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", value)
	}
}

func TestConcurrentDoublyLinked_Basic(t *testing.T) {
	cdl := lists.NewConcurrentDoublyLinked(1, 2, 3)
	
	if cdl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", cdl.Length())
	}
	
	value := cdl.Get(1, -1)
	if value != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", value)
	}
}

func TestConcurrentRWLinked_Basic(t *testing.T) {
	crwl := lists.NewConcurrentRWLinked(1, 2, 3)
	
	if crwl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", crwl.Length())
	}
	
	value := crwl.Get(1, -1)
	if value != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", value)
	}
}

func TestConcurrentRWDoublyLinked_Basic(t *testing.T) {
	crwdl := lists.NewConcurrentRWDoublyLinked(1, 2, 3)
	
	if crwdl.Length() != 3 {
		t.Errorf("Expected length 3, got %d", crwdl.Length())
	}
	
	value := crwdl.Get(1, -1)
	if value != 2 {
		t.Errorf("Expected value 2 at index 1, got %d", value)
	}
}
