// Command set-algebra builds two sets from the command line and prints the
// classic relationships between them: union, intersection, both differences,
// the subset test and the disjointness test.
//
// It demonstrates the sets package via the collections facade: collections.NewSet
// to build, then the algebraic operations from the Set interface. Output is
// sorted so it stays deterministic regardless of internal set ordering.
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/pickeringtech/go-collections/collections"
	"github.com/pickeringtech/go-collections/collections/sets"
	"github.com/pickeringtech/go-collections/slices"
)

func main() {
	aRaw := flag.String("a", "", "comma-separated members of set A")
	bRaw := flag.String("b", "", "comma-separated members of set B")
	flag.Parse()

	a := collections.NewSet(split(*aRaw)...)
	b := collections.NewSet(split(*bRaw)...)

	fmt.Printf("A = %s\n", show(a))
	fmt.Printf("B = %s\n", show(b))
	fmt.Printf("A union B        = %s\n", show(a.Union(b)))
	fmt.Printf("A intersect B    = %s\n", show(a.Intersection(b)))
	fmt.Printf("A minus B        = %s\n", show(a.Difference(b)))
	fmt.Printf("B minus A        = %s\n", show(b.Difference(a)))
	fmt.Printf("A subset of B    = %t\n", a.IsSubsetOf(b))
	fmt.Printf("A disjoint with B = %t\n", a.IsDisjoint(b))
}

// split turns "a, b ,c" into ["a" "b" "c"], trimming spaces and dropping blanks.
func split(s string) []string {
	cleaned := slices.Map(strings.Split(s, ","), strings.TrimSpace)
	return slices.Filter(cleaned, func(p string) bool { return p != "" })
}

// show renders a set as a sorted, brace-wrapped list for deterministic output.
func show(s sets.Set[string]) string {
	members := slices.SortOrderedAsc(s.AsSlice())
	if len(members) == 0 {
		return "{}"
	}
	return "{" + strings.Join(members, ", ") + "}"
}
