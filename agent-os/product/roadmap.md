# Product Roadmap

The core packages already exist: `slices`, `maps`, `channels`, `constraints`,
`concurrency`, and `collections/{dicts,lists,sets}`, and a comprehensive
verification harness is now in place (see "Verification harness" below). Phase 1
is about hardening these into a trustworthy, stable v1.0.

## Phase 1: v1.0 Launch

- **Full test/bench/example coverage** — every public function has the Example +
  Test + Benchmark trio (per `standards/testing/`). Remaining gap: backfill the
  missing **Benchmarks** in `maps`, `channels`, `sets`, `concurrency`
  ([#52](https://github.com/pickeringtech/go-collections/issues/52)).
- **API consistency audit** — sweep all packages for adherence to the interface,
  naming, and functional standards. The two originally-known gaps are **done**
  ([#1](https://github.com/pickeringtech/go-collections/issues/1) nil→empty and
  [#2](https://github.com/pickeringtech/go-collections/issues/2) concurrent
  return type, both closed). Remaining contract deviations to fix before the tag:
  - [#15](https://github.com/pickeringtech/go-collections/issues/15) — `SortInPlace` leaks into the immutable `List` interface.
  - [#16](https://github.com/pickeringtech/go-collections/issues/16) — reconcile `Get` semantics and `AsSlice` naming.
  - [#17](https://github.com/pickeringtech/go-collections/issues/17) — unify the `Searchable` contract (dicts missing `AllMatch`/`AnyMatch`; add `NoneMatch`).
  - [#18](https://github.com/pickeringtech/go-collections/issues/18) — lists parity (`IsEmpty`/`Clear`) + removal/membership decision for `[T any]`.
- **Iterator-native API (pre-1.0)** — add Go 1.23+ `iter.Seq`/`iter.Seq2`
  accessors to the collections *before* locking the API, so v1 is iterator-first
  rather than bolting it on later
  ([#53](https://github.com/pickeringtech/go-collections/issues/53)).
- **Documentation pass** — plain language, no ancillary emojis, enough context for
  all reader levels ([#34](https://github.com/pickeringtech/go-collections/issues/34)).
- **Stable v1.0 API + semver tag** — lock the public API, write a CHANGELOG, and
  cut `v1.0.0` ([#54](https://github.com/pickeringtech/go-collections/issues/54)).
  CI already surfaces a per-PR API-compatibility verdict via `gorelease`
  ([#29](https://github.com/pickeringtech/go-collections/issues/29), report-only
  pre-1.0); flip it to **blocking** once v1.0 is tagged.

## Verification harness (done)

The "real verification" mission pillar is largely delivered. **Blocking** via the
`CI Gate` aggregator ([#41](https://github.com/pickeringtech/go-collections/issues/41),
its `needs:` list): build, race+coverage tests across an OS × Go-version matrix
([#33](https://github.com/pickeringtech/go-collections/issues/33)), lint/complexity,
security (govulncheck + gosec), and the examples golden-output E2E
([#30](https://github.com/pickeringtech/go-collections/issues/30)).
**Report-only** (not in the gate; `continue-on-error`): fuzz
([#10](https://github.com/pickeringtech/go-collections/issues/10), fast count-based
run [#25](https://github.com/pickeringtech/go-collections/issues/25)), API
compatibility ([#29](https://github.com/pickeringtech/go-collections/issues/29)),
benchmark regression ([#31](https://github.com/pickeringtech/go-collections/issues/31)),
mutation testing ([#32](https://github.com/pickeringtech/go-collections/issues/32)),
and Codecov ([#14](https://github.com/pickeringtech/go-collections/issues/14)). The
last four ratchet toward blocking as they prove out.

## Phase 2: Post-Launch

- **More data structures** — expand the catalog (one issue per type):
  heap / priority queue ([#55](https://github.com/pickeringtech/go-collections/issues/55)),
  LRU cache ([#56](https://github.com/pickeringtech/go-collections/issues/56)),
  deque / ring buffer ([#57](https://github.com/pickeringtech/go-collections/issues/57)),
  ordered map + sorted set on the existing tree
  ([#58](https://github.com/pickeringtech/go-collections/issues/58)),
  multimap ([#59](https://github.com/pickeringtech/go-collections/issues/59)).
  Heavier/later: graphs, tries, persistent/immutable structures.
- **Fast variants for hot paths** — benchmark-proven `FooFast` variants where the
  readable default is a bottleneck (gated by the regression guard, #31).
- **Serialization & capacity hints** — `MarshalJSON`/`UnmarshalJSON` on the
  concrete types and `NewWithCapacity`-style preallocation constructors.
- **Richer channel pipelines** — fan-out/fan-in, batching, throttling, more
  composable stream operators.
