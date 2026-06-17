// Command leaderboard tallies a stream of scoring events per player and prints
// the standings two ways: every player's total in name order, then the top-N
// players by score.
//
// It demonstrates three facade paths the other examples leave uncovered:
//   - concurrent UpdateInPlace: each event is applied from its own goroutine to
//     a concurrent dict, so the read-modify-write that accumulates a player's
//     total is race-free. Addition is commutative, so the totals — and thus the
//     output — stay deterministic regardless of the order the goroutines run in.
//   - MapSorted: the unordered tally is projected into a Tree-backed sorted dict
//     so it can be iterated in ascending key (name) order.
//   - a heap (priority queue): the totals are heapified and drained in priority
//     order to rank the players, highest score first.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/pickeringtech/go-collections/collections"
	"github.com/pickeringtech/go-collections/collections/dicts"
)

// event is a single "player scored points" record parsed from the input.
type event struct {
	name  string
	score int
}

// standing pairs a player with their accumulated total for ranking.
type standing struct {
	name  string
	total int
}

func main() {
	raw := flag.String("events", "alice:5,bob:3,alice:2,carol:9,bob:4,alice:1,carol:1", "comma-separated name:score events")
	top := flag.Int("top", 3, "how many players to rank")
	flag.Parse()

	events, err := parse(*raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "leaderboard:", err)
		os.Exit(1)
	}

	totals := tally(events)

	fmt.Printf("events: %d scored across %d players\n", len(events), totals.Length())

	// MapSorted projects the unordered tally into a sorted-by-name view so the
	// per-player block prints in a stable, alphabetical order.
	byName := dicts.MapSorted(totals, func(name string, total int) (string, int) {
		return name, total
	})
	fmt.Println("totals (by name):")
	for name, total := range byName.All() {
		fmt.Printf("  %-8s %d\n", name, total)
	}

	fmt.Printf("top %d (by score):\n", *top)
	for i, s := range rank(totals, *top) {
		fmt.Printf("  %d. %-8s %d\n", i+1, s.name, s.total)
	}
}

// tally accumulates each event's score onto its player's running total. Every
// event is applied from its own goroutine via UpdateInPlace on a concurrent
// dict, which serialises the read-modify-write so no increment is lost.
func tally(events []event) dicts.MutableDict[string, int] {
	totals := collections.NewConcurrentDict[string, int]()
	var wg sync.WaitGroup
	for _, e := range events {
		wg.Add(1)
		go func() {
			defer wg.Done()
			totals.UpdateInPlace(e.name, func(old int, _ bool) int {
				return old + e.score
			})
		}()
	}
	wg.Wait()
	return totals
}

// rank heapifies the totals and drains the top n in priority order. The
// comparator ranks higher totals first and breaks ties by name so the order is
// fully determined.
func rank(totals dicts.Dict[string, int], n int) []standing {
	less := func(a, b standing) bool {
		if a.total != b.total {
			return a.total > b.total
		}
		return a.name < b.name
	}
	heap := collections.NewHeap[standing](less)
	totals.ForEach(func(name string, total int) {
		heap = heap.Push(standing{name: name, total: total})
	})

	ranked := heap.AsSortedSlice()
	if n > len(ranked) {
		n = len(ranked)
	}
	return ranked[:n]
}

// parse turns "alice:5, bob:3" into events, trimming spaces and dropping blanks.
func parse(s string) ([]event, error) {
	out := []event{}
	for _, field := range strings.Split(s, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		name, scoreRaw, ok := strings.Cut(field, ":")
		if !ok {
			return nil, fmt.Errorf("%q is not a name:score pair", field)
		}
		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("%q has an empty name", field)
		}
		score, err := strconv.Atoi(strings.TrimSpace(scoreRaw))
		if err != nil {
			return nil, fmt.Errorf("%q has a non-integer score", field)
		}
		out = append(out, event{name: name, score: score})
	}
	return out, nil
}
