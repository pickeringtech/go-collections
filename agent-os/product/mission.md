# Product Mission

## Problem

Go developers constantly rewrite the same standard collection types, the operations on them, and the functional utilities that work over them. The language ships sparse built-ins, so across projects you find many implementations of the same thing — a linked list, a thread-safe map, a set — each painful to work with because they lack:

- **Type-safe generic data structures** — forcing `interface{}` casts or per-project reimplementation.
- **Functional capabilities** — no built-in Filter/Map/Reduce over slices, maps, and channels, so transformation logic is repetitive and error-prone.
- **Clear, consistent APIs and thoughtful design** — the reimplementations that do exist are inconsistent and hard to trust.
- **Real verification** — these reimplementations typically have very little testing, no real benchmarking, and no fuzzing. There's no genuine way to know whether the collection is fit for purpose, because nothing rigorous has been imposed on it to prove it.

go-collections ends the per-repo reinvention with one comprehensive, reliable library.

## Target Users

- **Go application developers** — want ready-made, reliable collections and functional utilities instead of hand-rolling them per project.
- **Teams wanting consistency** — one well-tested, consistently-designed collections library across all their Go projects.
- **Library / framework authors** — need dependency-free, type-safe data structures and constraints to build on.
- **Developers from other ecosystems** — coming from Java/Python/JS and expecting rich collection APIs that Go's stdlib doesn't provide.

## Solution

A comprehensive, type-safe, high-performance collections library distinguished by:

- **Consistent, thoughtful API design** — uniform patterns across every collection: immutable/mutable dual hierarchy, `InPlace` naming, composed capability interfaces. Learn it once, apply everywhere.
- **Readable AND high-performance** — clear default implementations with opt-in `Fast` variants; every function benchmarked across a scaling ladder. No tradeoff between clarity and speed.
- **Exhaustively tested & documented** — every public function ships an Example, a table-driven Test, and a Benchmark; every package has a rich `doc.go`. Trustworthy, not throwaway.
- **Zero dependencies, type-safe** — pure Go, full generics, thread-safe variants (Mutex & RWMutex) for every collection. Nothing external to pull in.

> Design and contribution conventions are codified in `agent-os/standards/`.
