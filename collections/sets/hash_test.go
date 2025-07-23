package sets_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/sets"
	"reflect"
	"sort"
	"testing"
)

func ExampleHash_Contains() {
	s := sets.NewHash("apple", "banana", "cherry")

	fmt.Printf("Contains 'apple': %v\n", s.Contains("apple"))
	fmt.Printf("Contains 'grape': %v\n", s.Contains("grape"))

	// Output:
	// Contains 'apple': true
	// Contains 'grape': false
}

func ExampleHash_Add() {
	s := sets.NewHash("apple", "banana")

	newS := s.Add("cherry")
	
	fmt.Printf("Original length: %d\n", s.Length())
	fmt.Printf("New length: %d\n", newS.Length())
	fmt.Printf("New contains 'cherry': %v\n", newS.Contains("cherry"))

	// Output:
	// Original length: 2
	// New length: 3
	// New contains 'cherry': true
}

func ExampleHash_Union() {
	s1 := sets.NewHash("apple", "banana")
	s2 := sets.NewHash("cherry", "date")

	union := s1.Union(s2)
	
	fmt.Printf("Union length: %d\n", union.Length())
	fmt.Printf("Union contains 'apple': %v\n", union.Contains("apple"))
	fmt.Printf("Union contains 'cherry': %v\n", union.Contains("cherry"))

	// Output:
	// Union length: 4
	// Union contains 'apple': true
	// Union contains 'cherry': true
}

func ExampleHash_Intersection() {
	s1 := sets.NewHash("apple", "banana", "cherry")
	s2 := sets.NewHash("banana", "cherry", "date")

	intersection := s1.Intersection(s2)
	
	fmt.Printf("Intersection length: %d\n", intersection.Length())
	fmt.Printf("Intersection contains 'banana': %v\n", intersection.Contains("banana"))
	fmt.Printf("Intersection contains 'apple': %v\n", intersection.Contains("apple"))

	// Output:
	// Intersection length: 2
	// Intersection contains 'banana': true
	// Intersection contains 'apple': false
}

func ExampleHash_Difference() {
	s1 := sets.NewHash("apple", "banana", "cherry")
	s2 := sets.NewHash("banana", "date")

	difference := s1.Difference(s2)
	
	fmt.Printf("Difference length: %d\n", difference.Length())
	fmt.Printf("Difference contains 'apple': %v\n", difference.Contains("apple"))
	fmt.Printf("Difference contains 'banana': %v\n", difference.Contains("banana"))

	// Output:
	// Difference length: 2
	// Difference contains 'apple': true
	// Difference contains 'banana': false
}

func ExampleHash_Filter() {
	s := sets.NewHash(1, 2, 3, 4, 5, 6)

	evens := s.Filter(func(n int) bool {
		return n%2 == 0
	})
	
	fmt.Printf("Original length: %d\n", s.Length())
	fmt.Printf("Evens length: %d\n", evens.Length())

	// Output:
	// Original length: 6
	// Evens length: 3
}

func ExampleNewHash() {
	// Create an empty set
	empty := sets.NewHash[string]()
	fmt.Printf("Empty set length: %d\n", empty.Length())

	// Create a set with initial data
	fruits := sets.NewHash("apple", "banana", "cherry")
	fmt.Printf("Fruits set length: %d\n", fruits.Length())

	// Output:
	// Empty set length: 0
	// Fruits set length: 3
}

func TestHash_Contains(t *testing.T) {
	s := sets.NewHash("apple", "banana", "cherry")

	if !s.Contains("apple") {
		t.Error("Contains() should return true for existing element")
	}
	if s.Contains("grape") {
		t.Error("Contains() should return false for non-existing element")
	}
}

func TestHash_Length(t *testing.T) {
	tests := []struct {
		name string
		set  sets.Hash[string]
		want int
	}{
		{
			name: "empty set has length 0",
			set:  sets.NewHash[string](),
			want: 0,
		},
		{
			name: "set with elements has correct length",
			set:  sets.NewHash("apple", "banana", "cherry"),
			want: 3,
		},
		{
			name: "set with duplicate elements has correct length",
			set:  sets.NewHash("apple", "apple", "banana"),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash_IsEmpty(t *testing.T) {
	empty := sets.NewHash[string]()
	notEmpty := sets.NewHash("apple")

	if !empty.IsEmpty() {
		t.Error("IsEmpty() should return true for empty set")
	}
	if notEmpty.IsEmpty() {
		t.Error("IsEmpty() should return false for non-empty set")
	}
}

func TestHash_Add(t *testing.T) {
	original := sets.NewHash("apple", "banana")
	newSet := original.Add("cherry")

	// Original should be unchanged
	if original.Length() != 2 {
		t.Errorf("Add() modified original set, length = %v, want 2", original.Length())
	}

	// New set should have the new element
	if newSet.Length() != 3 {
		t.Errorf("Add() new set length = %v, want 3", newSet.Length())
	}
	if !newSet.Contains("cherry") {
		t.Error("Add() new set should contain the new element")
	}

	// Adding existing element should not change length
	sameSet := original.Add("apple")
	if sameSet.Length() != 2 {
		t.Errorf("Add() existing element length = %v, want 2", sameSet.Length())
	}
}

func TestHash_AddInPlace(t *testing.T) {
	s := sets.NewHash("apple", "banana")

	s.AddInPlace("cherry")
	if s.Length() != 3 {
		t.Errorf("AddInPlace() length = %v, want 3", s.Length())
	}
	if !s.Contains("cherry") {
		t.Error("AddInPlace() should add the element")
	}

	// Adding existing element should not change length
	s.AddInPlace("apple")
	if s.Length() != 3 {
		t.Errorf("AddInPlace() existing element length = %v, want 3", s.Length())
	}
}

func TestHash_Remove(t *testing.T) {
	original := sets.NewHash("apple", "banana", "cherry")
	newSet := original.Remove("banana")

	// Original should be unchanged
	if original.Length() != 3 {
		t.Errorf("Remove() modified original set, length = %v, want 3", original.Length())
	}
	if !original.Contains("banana") {
		t.Error("Remove() modified original set, should still contain 'banana'")
	}

	// New set should not have the removed element
	if newSet.Length() != 2 {
		t.Errorf("Remove() new set length = %v, want 2", newSet.Length())
	}
	if newSet.Contains("banana") {
		t.Error("Remove() new set should not contain the removed element")
	}
	if !newSet.Contains("apple") || !newSet.Contains("cherry") {
		t.Error("Remove() new set should still contain other elements")
	}
}

func TestHash_RemoveInPlace(t *testing.T) {
	s := sets.NewHash("apple", "banana", "cherry")

	// Remove existing element
	removed := s.RemoveInPlace("banana")
	if !removed {
		t.Error("RemoveInPlace() should return true for existing element")
	}
	if s.Length() != 2 {
		t.Errorf("RemoveInPlace() length = %v, want 2", s.Length())
	}
	if s.Contains("banana") {
		t.Error("RemoveInPlace() should remove the element")
	}

	// Remove non-existing element
	removed = s.RemoveInPlace("grape")
	if removed {
		t.Error("RemoveInPlace() should return false for non-existing element")
	}
}

func TestHash_Union(t *testing.T) {
	s1 := sets.NewHash("apple", "banana")
	s2 := sets.NewHash("cherry", "date")

	union := s1.Union(s2)

	if union.Length() != 4 {
		t.Errorf("Union() length = %v, want 4", union.Length())
	}

	expected := []string{"apple", "banana", "cherry", "date"}
	for _, element := range expected {
		if !union.Contains(element) {
			t.Errorf("Union() should contain %s", element)
		}
	}
}

func TestHash_Intersection(t *testing.T) {
	s1 := sets.NewHash("apple", "banana", "cherry")
	s2 := sets.NewHash("banana", "cherry", "date")

	intersection := s1.Intersection(s2)

	if intersection.Length() != 2 {
		t.Errorf("Intersection() length = %v, want 2", intersection.Length())
	}
	if !intersection.Contains("banana") || !intersection.Contains("cherry") {
		t.Error("Intersection() should contain common elements")
	}
	if intersection.Contains("apple") || intersection.Contains("date") {
		t.Error("Intersection() should not contain non-common elements")
	}
}

func TestHash_AsSlice(t *testing.T) {
	s := sets.NewHash("apple", "banana", "cherry")
	slice := s.AsSlice()

	if len(slice) != 3 {
		t.Errorf("AsSlice() length = %v, want 3", len(slice))
	}

	// Sort for consistent comparison
	sort.Strings(slice)
	expected := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("AsSlice() = %v, want %v", slice, expected)
	}
}

func TestHash_SetOperations(t *testing.T) {
	s1 := sets.NewHash(1, 2, 3)
	s2 := sets.NewHash(2, 3, 4)
	s3 := sets.NewHash(1, 2)

	// IsSubsetOf
	if !s3.IsSubsetOf(s1) {
		t.Error("IsSubsetOf() should return true for subset")
	}
	if s1.IsSubsetOf(s3) {
		t.Error("IsSubsetOf() should return false for non-subset")
	}

	// IsSupersetOf
	if !s1.IsSupersetOf(s3) {
		t.Error("IsSupersetOf() should return true for superset")
	}
	if s3.IsSupersetOf(s1) {
		t.Error("IsSupersetOf() should return false for non-superset")
	}

	// IsDisjoint
	s4 := sets.NewHash(5, 6, 7)
	if !s1.IsDisjoint(s4) {
		t.Error("IsDisjoint() should return true for disjoint sets")
	}
	if s1.IsDisjoint(s2) {
		t.Error("IsDisjoint() should return false for overlapping sets")
	}

	// Equals
	s5 := sets.NewHash(1, 2, 3)
	if !s1.Equals(s5) {
		t.Error("Equals() should return true for equal sets")
	}
	if s1.Equals(s2) {
		t.Error("Equals() should return false for different sets")
	}
}
