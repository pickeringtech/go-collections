# Product Roadmap

The core packages already exist: `slices`, `maps`, `channels`, `constraints`, `concurrency`, and `collections/{dicts,lists,sets}`. Phase 1 is about hardening these into a trustworthy, stable v1.0.

## Phase 1: v1.0 Launch

- **Full test/bench/example coverage** — every public function across all packages has the Example + Test + Benchmark trio (per `standards/testing/`); no gaps before tagging v1.0.
- **API consistency audit** — sweep all packages for adherence to the interface, naming, and functional standards (composed capability interfaces, `InPlace` suffix, `Func` aliases, param ordering, zero-value output). Fix deviations, including the two known gaps:
  - [#1](https://github.com/pickeringtech/go-collections/issues/1) — `slices` functions return `nil` on empty/nil input; standardize on non-nil empty.
  - [#2](https://github.com/pickeringtech/go-collections/issues/2) — concurrent immutable ops return a plain type; return the same concurrent type.
- **Stable v1.0 API + semver tag** — lock the public API, document any breaking changes, and cut a tagged `v1.0.0` release with a CHANGELOG. CI already surfaces a per-PR API-compatibility verdict via `gorelease` ([#29](https://github.com/pickeringtech/go-collections/issues/29), report-only pre-1.0); flip it to blocking once v1.0 is tagged.

## Phase 2: Post-Launch

- **Fast variants for hot paths** — add benchmark-proven `FooFast` optimized variants where the readable default is a bottleneck (per `standards/global/readability-and-performance.md`).
- **More data structures** — expand the catalog: trees/heaps/priority queues, ordered maps, deques, ring buffers, graphs.
- **Richer channel pipelines** — grow the `channels` package: fan-out/fan-in, batching, throttling, more composable stream operators.
- **Iterator (range-over-func) support** — add Go 1.23+ `iter.Seq`/`iter.Seq2` support so collections integrate with range-over-func and the stdlib `iter` ecosystem.
