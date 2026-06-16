package dicts_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/pickeringtech/go-collections/collections/dicts"
)

// --- Hash coverage gaps ---------------------------------------------------

func TestHash_ForEachKey(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	var keys []string
	h.ForEachKey(func(key string) {
		keys = append(keys, key)
	})
	sort.Strings(keys)

	if !reflect.DeepEqual(keys, []string{"a", "b"}) {
		t.Errorf("ForEachKey keys = %v, want [a b]", keys)
	}
}

func TestHash_ForEachValue(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	sum := 0
	h.ForEachValue(func(value int) {
		sum += value
	})

	if sum != 3 {
		t.Errorf("ForEachValue sum = %d, want 3", sum)
	}
}

func TestHash_FilterInPlace(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)

	h.FilterInPlace(func(key string, value int) bool {
		return value%2 == 1
	})

	if h.Length() != 2 {
		t.Errorf("FilterInPlace length = %d, want 2", h.Length())
	}
	if h.Contains("b") {
		t.Error("FilterInPlace should have removed b")
	}
}

func TestHash_Find(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	k, v, ok := h.Find(func(key string, value int) bool {
		return value == 2
	})
	if !ok || k != "b" || v != 2 {
		t.Errorf("Find = %q, %d, %v, want b, 2, true", k, v, ok)
	}

	k, v, ok = h.Find(func(key string, value int) bool {
		return value == 99
	})
	if ok || k != "" || v != 0 {
		t.Errorf("Find(missing) = %q, %d, %v, want \"\", 0, false", k, v, ok)
	}
}

func TestHash_FindKey(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	k, ok := h.FindKey(func(key string) bool {
		return key == "a"
	})
	if !ok || k != "a" {
		t.Errorf("FindKey = %q, %v, want a, true", k, ok)
	}

	k, ok = h.FindKey(func(key string) bool {
		return key == "zzz"
	})
	if ok || k != "" {
		t.Errorf("FindKey(missing) = %q, %v, want \"\", false", k, ok)
	}
}

func TestHash_FindValue(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	v, ok := h.FindValue(func(value int) bool {
		return value == 2
	})
	if !ok || v != 2 {
		t.Errorf("FindValue = %d, %v, want 2, true", v, ok)
	}

	v, ok = h.FindValue(func(value int) bool {
		return value == 99
	})
	if ok || v != 0 {
		t.Errorf("FindValue(missing) = %d, %v, want 0, false", v, ok)
	}
}

func TestHash_ContainsValue(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	if !h.ContainsValue(1) {
		t.Error("ContainsValue(1) = false, want true")
	}
	if h.ContainsValue(99) {
		t.Error("ContainsValue(99) = true, want false")
	}
}

func TestHash_AsMap(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	m := h.AsMap()
	if len(m) != 2 || m["a"] != 1 || m["b"] != 2 {
		t.Errorf("AsMap() = %v, want a=1 b=2", m)
	}

	// Returned map is an independent copy.
	m["a"] = 100
	if v, _ := h.Get("a", -1); v != 1 {
		t.Error("AsMap returned map is not independent of the dict")
	}
}

func TestHash_PutMany(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
	)

	result := h.PutMany(
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)
	if result.Length() != 3 {
		t.Errorf("PutMany result length = %d, want 3", result.Length())
	}
	if h.Length() != 1 {
		t.Errorf("original length = %d, want 1 (PutMany must not mutate)", h.Length())
	}
}

func TestHash_PutManyInPlace(t *testing.T) {
	h := dicts.NewHash[string, int]()

	h.PutManyInPlace(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)
	if h.Length() != 2 {
		t.Errorf("Length after PutManyInPlace = %d, want 2", h.Length())
	}
}

func TestHash_RemoveMany(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)

	result := h.RemoveMany("a", "b")
	if result.Length() != 1 {
		t.Errorf("RemoveMany result length = %d, want 1", result.Length())
	}
	if !result.Contains("c") {
		t.Error("RemoveMany result should retain c")
	}
	if h.Length() != 3 {
		t.Errorf("original length = %d, want 3 (RemoveMany must not mutate)", h.Length())
	}
}

func TestHash_RemoveInPlace(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
	)

	v, ok := h.RemoveInPlace("a")
	if !ok || v != 1 {
		t.Errorf("RemoveInPlace(a) = %d, %v, want 1, true", v, ok)
	}

	v, ok = h.RemoveInPlace("missing")
	if ok || v != 0 {
		t.Errorf("RemoveInPlace(missing) = %d, %v, want 0, false", v, ok)
	}
}

func TestHash_RemoveManyInPlace(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
		dicts.Pair[string, int]{Key: "c", Value: 3},
	)

	h.RemoveManyInPlace("a", "c")
	if h.Length() != 1 {
		t.Errorf("Length after RemoveManyInPlace = %d, want 1", h.Length())
	}
	if !h.Contains("b") {
		t.Error("RemoveManyInPlace should retain b")
	}
}

func TestHash_Clear(t *testing.T) {
	h := dicts.NewHash(
		dicts.Pair[string, int]{Key: "a", Value: 1},
		dicts.Pair[string, int]{Key: "b", Value: 2},
	)

	h.Clear()
	if !h.IsEmpty() {
		t.Error("Clear should leave the dict empty")
	}
}

// --- Tree coverage gaps ---------------------------------------------------

func TestTree_Put(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	result := tree.Put(1, "one")
	if result.Length() != 2 {
		t.Errorf("Put result length = %d, want 2", result.Length())
	}
	if tree.Length() != 1 {
		t.Errorf("original length = %d, want 1 (Put must not mutate)", tree.Length())
	}
}

func TestTree_PutMany(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	result := tree.PutMany(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 3, Value: "three"},
	)
	if result.Length() != 3 {
		t.Errorf("PutMany result length = %d, want 3", result.Length())
	}
	if tree.Length() != 1 {
		t.Errorf("original length = %d, want 1 (PutMany must not mutate)", tree.Length())
	}
}

func TestTree_PutManyInPlace(t *testing.T) {
	tree := dicts.NewTree[int, string]()

	tree.PutManyInPlace(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)
	if tree.Length() != 2 {
		t.Errorf("Length after PutManyInPlace = %d, want 2", tree.Length())
	}
}

func TestTree_PutInPlace_UpdateExisting(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
	)

	tree.PutInPlace(1, "ONE")
	v, ok := tree.Get(1, "")
	if !ok || v != "ONE" {
		t.Errorf("after update Get(1) = %q, %v, want ONE, true", v, ok)
	}
	if tree.Length() != 1 {
		t.Errorf("Length = %d, want 1 (update must not grow size)", tree.Length())
	}
}

func TestTree_ForEachKey(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	var keys []int
	tree.ForEachKey(func(key int) {
		keys = append(keys, key)
	})

	if !reflect.DeepEqual(keys, []int{1, 2, 3}) {
		t.Errorf("ForEachKey = %v, want [1 2 3]", keys)
	}
}

func TestTree_ForEachValue(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	var values []string
	tree.ForEachValue(func(value string) {
		values = append(values, value)
	})

	if !reflect.DeepEqual(values, []string{"one", "two", "three"}) {
		t.Errorf("ForEachValue = %v, want [one two three]", values)
	}
}

func TestTree_FilterInPlace(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
		dicts.Pair[int, string]{Key: 3, Value: "three"},
		dicts.Pair[int, string]{Key: 4, Value: "four"},
	)

	tree.FilterInPlace(func(key int, value string) bool {
		return key%2 == 0
	})

	if !reflect.DeepEqual(tree.Keys(), []int{2, 4}) {
		t.Errorf("FilterInPlace keys = %v, want [2 4]", tree.Keys())
	}
}

func TestTree_FindKey(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	k, ok := tree.FindKey(func(key int) bool {
		return key > 1
	})
	if !ok || k != 2 {
		t.Errorf("FindKey = %d, %v, want 2, true", k, ok)
	}

	k, ok = tree.FindKey(func(key int) bool {
		return key > 100
	})
	if ok || k != 0 {
		t.Errorf("FindKey(missing) = %d, %v, want 0, false", k, ok)
	}
}

func TestTree_FindValue(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	v, ok := tree.FindValue(func(value string) bool {
		return value == "two"
	})
	if !ok || v != "two" {
		t.Errorf("FindValue = %q, %v, want two, true", v, ok)
	}

	v, ok = tree.FindValue(func(value string) bool {
		return value == "nope"
	})
	if ok || v != "" {
		t.Errorf("FindValue(missing) = %q, %v, want \"\", false", v, ok)
	}
}

func TestTree_ContainsValue(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	if !tree.ContainsValue("one") {
		t.Error("ContainsValue(one) = false, want true")
	}
	if tree.ContainsValue("nope") {
		t.Error("ContainsValue(nope) = true, want false")
	}
}

func TestTree_Values(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 2, Value: "two"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
	)

	if !reflect.DeepEqual(tree.Values(), []string{"one", "two"}) {
		t.Errorf("Values() = %v, want [one two]", tree.Values())
	}
}

func TestTree_Items(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 2, Value: "two"},
		dicts.Pair[int, string]{Key: 1, Value: "one"},
	)

	items := tree.Items()
	want := []dicts.Pair[int, string]{
		{Key: 1, Value: "one"},
		{Key: 2, Value: "two"},
	}
	if !reflect.DeepEqual(items, want) {
		t.Errorf("Items() = %v, want %v", items, want)
	}
}

func TestTree_AsMap(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
	)

	m := tree.AsMap()
	if len(m) != 2 || m[1] != "one" || m[2] != "two" {
		t.Errorf("AsMap() = %v, want 1=one 2=two", m)
	}
}

func TestTree_RemoveMany(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
		dicts.Pair[int, string]{Key: 3, Value: "three"},
	)

	result := tree.RemoveMany(1, 3)
	if !reflect.DeepEqual(result.Keys(), []int{2}) {
		t.Errorf("RemoveMany keys = %v, want [2]", result.Keys())
	}
	if tree.Length() != 3 {
		t.Errorf("original length = %d, want 3 (RemoveMany must not mutate)", tree.Length())
	}
}

func TestTree_RemoveManyInPlace(t *testing.T) {
	tree := dicts.NewTree(
		dicts.Pair[int, string]{Key: 1, Value: "one"},
		dicts.Pair[int, string]{Key: 2, Value: "two"},
		dicts.Pair[int, string]{Key: 3, Value: "three"},
	)

	tree.RemoveManyInPlace(1, 2)
	if !reflect.DeepEqual(tree.Keys(), []int{3}) {
		t.Errorf("RemoveManyInPlace keys = %v, want [3]", tree.Keys())
	}
}

// TestTree_RemoveInPlace_AllCases exercises every branch of removeNode:
// leaf, single-left-child, single-right-child, and two-children nodes, plus
// removing a non-existent key from a populated tree.
func TestTree_RemoveInPlace_AllCases(t *testing.T) {
	build := func() *dicts.Tree[int, string] {
		// Insertion order yields this shape:
		//        50
		//       /  \
		//     30    70
		//    /  \     \
		//  20    40    80
		//  /
		// 10
		return dicts.NewTree(
			dicts.Pair[int, string]{Key: 50, Value: "50"},
			dicts.Pair[int, string]{Key: 30, Value: "30"},
			dicts.Pair[int, string]{Key: 70, Value: "70"},
			dicts.Pair[int, string]{Key: 20, Value: "20"},
			dicts.Pair[int, string]{Key: 40, Value: "40"},
			dicts.Pair[int, string]{Key: 80, Value: "80"},
			dicts.Pair[int, string]{Key: 10, Value: "10"},
		)
	}

	t.Run("leaf node", func(t *testing.T) {
		tree := build()
		v, ok := tree.RemoveInPlace(40)
		if !ok || v != "40" {
			t.Errorf("RemoveInPlace(40) = %q, %v, want 40, true", v, ok)
		}
		if tree.Contains(40) {
			t.Error("40 should be gone")
		}
		// Tree remains a valid BST (sorted order preserved).
		if !sortedKeys(tree) {
			t.Errorf("keys not sorted after leaf removal: %v", tree.Keys())
		}
	})

	t.Run("node with only right child", func(t *testing.T) {
		tree := build()
		// 70 has only a right child (80).
		v, ok := tree.RemoveInPlace(70)
		if !ok || v != "70" {
			t.Errorf("RemoveInPlace(70) = %q, %v, want 70, true", v, ok)
		}
		if tree.Contains(70) || !tree.Contains(80) {
			t.Error("70 should be removed and 80 retained")
		}
		if !sortedKeys(tree) {
			t.Errorf("keys not sorted after right-only removal: %v", tree.Keys())
		}
	})

	t.Run("node with only left child", func(t *testing.T) {
		tree := build()
		// 20 has only a left child (10).
		v, ok := tree.RemoveInPlace(20)
		if !ok || v != "20" {
			t.Errorf("RemoveInPlace(20) = %q, %v, want 20, true", v, ok)
		}
		if tree.Contains(20) || !tree.Contains(10) {
			t.Error("20 should be removed and 10 retained")
		}
		if !sortedKeys(tree) {
			t.Errorf("keys not sorted after left-only removal: %v", tree.Keys())
		}
	})

	t.Run("node with two children", func(t *testing.T) {
		tree := build()
		// 30 has two children (20 and 40); removal triggers findMin successor.
		v, ok := tree.RemoveInPlace(30)
		if !ok || v != "30" {
			t.Errorf("RemoveInPlace(30) = %q, %v, want 30, true", v, ok)
		}
		if tree.Contains(30) {
			t.Error("30 should be removed")
		}
		if !sortedKeys(tree) {
			t.Errorf("keys not sorted after two-child removal: %v", tree.Keys())
		}
		// All other keys intact.
		for _, k := range []int{10, 20, 40, 50, 70, 80} {
			if !tree.Contains(k) {
				t.Errorf("key %d should still be present", k)
			}
		}
	})

	t.Run("root node", func(t *testing.T) {
		tree := build()
		v, ok := tree.RemoveInPlace(50)
		if !ok || v != "50" {
			t.Errorf("RemoveInPlace(50) = %q, %v, want 50, true", v, ok)
		}
		if tree.Contains(50) {
			t.Error("root 50 should be removed")
		}
		if !sortedKeys(tree) {
			t.Errorf("keys not sorted after root removal: %v", tree.Keys())
		}
	})

	t.Run("missing key in populated tree", func(t *testing.T) {
		tree := build()
		before := tree.Length()
		v, ok := tree.RemoveInPlace(999)
		if ok || v != "" {
			t.Errorf("RemoveInPlace(999) = %q, %v, want \"\", false", v, ok)
		}
		if tree.Length() != before {
			t.Errorf("length changed after removing missing key: %d, want %d", tree.Length(), before)
		}
	})
}

// sortedKeys reports whether the tree's keys come back in ascending order.
func sortedKeys(tree *dicts.Tree[int, string]) bool {
	keys := tree.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i-1] >= keys[i] {
			return false
		}
	}
	return true
}
