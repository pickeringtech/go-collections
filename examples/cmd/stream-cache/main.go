// Command stream-cache replays a stream of key accesses through a bounded LRU
// cache and reports the hit/miss outcome, the evictions, the surviving entries
// and a sliding window of the most recent keys.
//
// It demonstrates three facade paths the other examples leave uncovered:
//   - channels with context: the keys are presented as a context-governed
//     channel. After -limit accesses the context is cancelled, which tears the
//     producing goroutine down and closes the channel even though keys remain
//     unsent — the program then observes the channel is closed.
//   - LRU: a bounded cache promotes keys on access (Get) and evicts the
//     least-recently-used entry on overflow, reported through an eviction
//     callback registered at construction.
//   - a bounded deque: a fixed-capacity ring buffer keeps the last -window keys
//     in arrival order, overwriting the oldest as new keys arrive.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pickeringtech/go-collections/channels"
	"github.com/pickeringtech/go-collections/collections"
	"github.com/pickeringtech/go-collections/collections/deques"
	"github.com/pickeringtech/go-collections/collections/lru"
)

func main() {
	raw := flag.String("keys", "a,b,c,a,d,b,e,a,f", "comma-separated stream of key accesses")
	capacity := flag.Int("cap", 3, "LRU cache capacity")
	window := flag.Int("window", 4, "size of the recent-keys sliding window")
	limit := flag.Int("limit", 6, "stop and cancel the stream after this many accesses")
	flag.Parse()

	keys := split(*raw)
	if *capacity < 1 || *window < 1 {
		fmt.Fprintln(os.Stderr, "stream-cache: -cap and -window must be at least 1")
		os.Exit(1)
	}

	// The context governs the producing goroutine. Cancelling it (below) stops
	// the stream and closes the channel, so a partially consumed source does not
	// leak its goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// LRU: record every eviction the cache makes so we can report it afterwards.
	var evicted []string
	cache := collections.NewLRU[string, int](*capacity, lru.WithOnEvict(func(key string, _ int) {
		evicted = append(evicted, key)
	}))

	// A bounded deque keeps the last -window keys; OverwriteOldest drops the
	// front as each new key is pushed to the back once it is full.
	recent := collections.NewBoundedDeque[string](*window, deques.OverwriteOldest)

	stream := channels.FromSlice(ctx, keys)

	hits, misses, processed := 0, 0, 0
	for key := range stream {
		prev, hit := cache.Get(key)
		if hit {
			hits++
			cache.PutInPlace(key, prev+1)
		} else {
			misses++
			cache.PutInPlace(key, 1)
		}
		recent = recent.PushBack(key)
		processed++
		if processed == *limit {
			break
		}
	}

	// Cancel mid-stream and confirm the producer responded by closing the
	// channel: a receive on a closed, drained channel reports ok == false.
	cancel()
	_, ok := <-stream
	streamClosed := !ok

	fmt.Printf("processed %d of %d accesses (%d hit, %d miss)\n", processed, len(keys), hits, misses)
	fmt.Printf("stream closed after cancel: %t\n", streamClosed)
	fmt.Printf("evicted (oldest first): %s\n", join(evicted))
	fmt.Printf("cache (most→least recent): %s\n", join(cache.Keys()))
	fmt.Printf("recent window (front→back): %s\n", join(recent.AsSlice()))
}

// split turns "a, b ,c" into ["a" "b" "c"], trimming spaces and dropping blanks.
func split(s string) []string {
	out := []string{}
	for _, field := range strings.Split(s, ",") {
		field = strings.TrimSpace(field)
		if field != "" {
			out = append(out, field)
		}
	}
	return out
}

// join renders a slice of keys for printing, with a placeholder when empty.
func join(keys []string) string {
	if len(keys) == 0 {
		return "(none)"
	}
	return strings.Join(keys, ", ")
}
