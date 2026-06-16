// Command collection-transform runs a list of words through the free-function
// Map / FlatMap / Reduce helpers that operate over the collection interfaces.
//
// These are free functions rather than methods because a general T -> U
// transform needs a type parameter on the operation itself, which Go methods
// cannot express (golang/go#49085). Each helper returns a collection interface
// (List / Dict / Set) backed by the default concrete type, so results chain on
// into other collection helpers.
//
// It demonstrates, end-to-end:
//   - lists.Map     — word -> length (T -> U, a new element type);
//   - lists.FlatMap — word -> its runes (one element expands into many);
//   - lists.Reduce  — fold the lengths into a single total;
//   - sets.Map      — collapse the lengths into their distinct values;
//   - dicts.Map     — build a word -> length Dict, then Reduce it to a total.
package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/pickeringtech/go-collections/collections/dicts"
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/collections/sets"
)

func main() {
	raw := flag.String("words", "the,quick,brown,fox,the,lazy,dog", "comma-separated words to transform")
	flag.Parse()

	words := parse(*raw)
	list := lists.NewArray(words...)

	// Map: T -> U. Each word becomes its length, a different element type.
	lengths := lists.Map(list, func(w string) int { return len(w) })

	// FlatMap: each word expands into its runes, concatenated into one List.
	runes := lists.FlatMap(list, func(w string) lists.List[string] {
		out := []string{}
		for _, r := range w {
			out = append(out, string(r))
		}
		return lists.NewArray(out...)
	})

	// Reduce: fold the lengths into a single accumulated total.
	totalLen := lists.Reduce(list, 0, func(acc int, w string) int {
		return acc + len(w)
	})

	// sets.Map collapses the lengths into their distinct values (order is
	// unspecified over a Set, so sort for a stable display).
	distinctLengths := sets.Map(setOfLengths(lengths), func(n int) int { return n })
	distinct := distinctLengths.AsSlice()
	sort.Ints(distinct)

	// dicts.Map builds a word -> length Dict; dicts.Reduce folds it to a total.
	lengthByWord := dicts.Map(dictOfWords(words), func(_ int, w string) (string, int) {
		return w, len(w)
	})
	dictTotal := dicts.Reduce(lengthByWord, 0, func(acc int, _ string, v int) int {
		return acc + v
	})

	fmt.Printf("input:       %v\n", words)
	fmt.Printf("lengths:     %v\n", lengths.AsSlice())
	fmt.Printf("runes:       %v\n", runes.AsSlice())
	fmt.Printf("total:       %d\n", totalLen)
	fmt.Printf("distinct:    %v\n", distinct)
	fmt.Printf("dict total:  %d (over %d distinct words)\n", dictTotal, lengthByWord.Length())
}

// setOfLengths copies a List of lengths into a Set so it can be mapped with
// sets.Map.
func setOfLengths(lengths lists.List[int]) sets.Set[int] {
	return sets.NewHash(lengths.AsSlice()...)
}

// dictOfWords builds an index -> word Dict from the input, the shape dicts.Map
// consumes.
func dictOfWords(words []string) dicts.Dict[int, string] {
	pairs := make([]dicts.Pair[int, string], 0, len(words))
	for i, w := range words {
		pairs = append(pairs, dicts.Pair[int, string]{Key: i, Value: w})
	}
	return dicts.NewHash(pairs...)
}

// parse splits "the, quick ,fox" into [the quick fox], trimming spaces and
// dropping blank fields.
func parse(s string) []string {
	out := []string{}
	for _, field := range strings.Split(s, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		out = append(out, field)
	}
	return out
}
