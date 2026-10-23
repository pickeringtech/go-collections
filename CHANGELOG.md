# Changelog

All notable changes to this project are documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and
the project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
While the major version is 0, a minor release (0.x.0) may contain breaking
changes. Each release lists them under **Breaking changes**, with what to do
instead. The API freezes at v1.0.0 (#54), from which point the
`release.yml` API-compatibility check gates releases.

## [Unreleased]

## [0.1.0] - 2026-10-23

The first release since v0.0.2 (January 2024). Over that time the library grew
from five utility packages into a full collections, statistics and
data-engineering toolkit. The release ends with a hardening pass: an
independent audit of every package confirmed around 60 defects, each fixed with
a regression test that failed beforehand.

**Requirements:** Go 1.24 or later. Zero dependencies.

### Highlights

- **Collections.** The library now has deques (ring buffer, optionally
  bounded), binary heaps (priority queues), an LRU cache, multimaps (list- or
  set-backed) and tree-backed ordered dicts and sets, alongside the expanded
  dicts, sets and lists. Every structure has `Concurrent*` (mutex) and
  `ConcurrentRW*` (read-write mutex) variants. Callbacks never run under the
  lock, so they may call back into the collection.
- **Iterators.** Range-over-func `iter.Seq`/`iter.Seq2` accessors work across
  every collection, with `FromSeq` constructors for the inbound direction.
- **New packages:**
  - `stats`: sums, means, quantiles, variance, covariance, correlation,
    regression, shape, entropy and transforms, with Kahan/Welford numerics.
  - `relational`: GroupBy and aggregate, the four joins, pivot/unpivot and
    partition.
  - `collections/sketches`: MinHash, Bloom, Count-Min, HyperLogLog and
    t-digest.
  - `collections/streaming`: top-k, reservoir and weighted-reservoir
    sampling, bootstrap resampling, and online mean, variance, EWMA and
    min/max.
  - `ml/distance` and `ml/similarity`.
  - `ml/metrics`: classification, clustering, ranking and regression metrics.
  - `ml/preprocessing`: scalers, encoders, imputers, binners and dataset
    splits.
- **Concurrency.** Order-preserving, bounded-concurrency `Map`, `ForEach` and
  `Batch`, with error policies, alongside the blocking and background work
  limiters.
- **Channels.** Every pipeline stage is now context-aware, and there are new
  tumbling, sliding and session windows plus `WindowedReduce`.
- **Verification.** Every merge is gated on:
  - 100% statement coverage;
  - `-race` tests on Linux, macOS and Windows;
  - Go 1.24 and the latest stable Go;
  - cross-architecture builds (386, arm64 and big-endian s390x);
  - fuzzing;
  - golangci-lint and gosec;
  - doc examples that compile;
  - a benchmark regression gate.

  Mutation testing, deep fuzzing and govulncheck run on a schedule.

### Breaking changes

Since v0.0.2:

- **Go 1.24 is the minimum version.**
- **`channels`:** every stage takes a `context.Context` as its first parameter.
  This covers `FromSlice`, `FromMap`, `Map`, `Filter`, `Reduce`, `NewPipeline`
  and `PipelineCreationFunc`. Cancelling it stops the stage and closes its
  output. `Map` returns a receive-only channel.
- **`collections` facade:** the constructors take initial values and return the
  mutable interfaces: `MutableList`, `MutableQueue`, `MutableStack`,
  `MutableDict`, `MutableSet`, `MutableHeap`, `MutableDeque` and so on.
  Immutable operations are still available through them. Code that reassigned
  an immutable result to the same variable, such as `h = h.Push(x)`, should use
  the `…InPlace` method instead.
- **`collections/dicts`:** the interfaces are split into small role interfaces
  composed into `Dict`/`MutableDict`. `Tree` requires `constraints.Ordered` keys,
  and `NewTree` returns `*Tree`.
- **`collections/lists`:** immutable operations return a `List[T]` rather than a
  `[]T` (call `AsSlice` for a slice). `Get` returns `(T, bool)`, and
  `GetAsSlice` is removed.
- **`collections/sets`:** `Set[T comparable]` (was `any`).
- **`slices`:** `Sum` and `Avg` are removed in favour of `stats.Sum` and
  `stats.Mean`, the single home for numeric summaries. `Max`, `Min` and
  `NumericSlice.Max`/`Min`/`Sum` return `(T, bool)`, with `false` for empty
  input.

From this release's hardening pass, mostly behaviour corrections. Callers who
relied on the old behaviour should check these:

- **`slices`:**
  - `AllMatch` on empty input returns `true` (vacuous truth), consistent with
    the lists and sets.
  - Appending to a result can no longer overwrite the caller's slice.
    `Concatenate`, `Push`, `PushFront`, `Delete`, `Fill*` and `Insert` return
    fresh storage on every path. `Pop`, `PopFront`, `SubSlice` and `Paginate`
    stay O(1): they return views whose capacity is clipped to their length, so
    an `append` always allocates. Writing to an element of a view still writes
    through, as with any Go sub-slice; use `slices.Copy` for an independent
    copy.
  - `Max`/`Min` propagate NaN as the standard library does.
- **Zero-value outputs:** these return an empty, non-nil slice rather than
  `nil`: `maps.Keys`, `Values`, `Items`, `GetMany` and `GetManyOrDefault`;
  `channels.CollectAsSlice` and `CollectNAsSlice`; `slices.Pop`, `PopFront`,
  `SubSlice` and `Paginate`; and `lists.Linked.AsSlice`.
- **`collections/streaming`:** `RunningMean.Count()` and
  `RunningVariance.Count()` return `int64`. `RunningMean` reports `(0, false)`
  after a non-finite input, matching `stats.Mean`.
- **`collections/sketches`:**
  - `EstimatedJaccard` of two empty MinHash sketches returns `(0, false)`
    (undefined) instead of `(1, true)`.
  - The shared hash function changed (see Fixed). Bloom, Count-Min and HLL
    outputs for a given key and seed differ from earlier versions. The sketches
    have no serialised form, so nothing persisted needs migrating.
- **`collections/lists`:** every list sorts stably.
- **`ml/similarity`:** `Jaccard`, `Dice` and `Overlap` return
  `(float64, bool)`, with `ok == false` when the result is undefined (empty
  sets), matching `DotProduct`/`CosineSimilarity` (#235).
- **`ml/distance`:** `Levenshtein` is generic over `[]T` for any comparable
  `T`. Use the new `StringLevenshtein(a, b string)` for strings; it compares
  runes, as before (#237).
- **`ml/preprocessing`:**
  - A rejected `Fit` resets the estimator to unfitted (`Transform` reports
    `ok == false`) instead of silently keeping the previous fit's parameters.
  - `StratifiedSplit` sizes its test set exactly as `TrainTestSplit` does, so
    seeded splits differ from earlier versions.
- **NaN labels are rejected (`ok == false`)** by the classification metrics
  and `StratifiedSplit`, as are negative relevances by `NDCG`.

### Added

- `collections/deques`: `RingBuffer` deque, unbounded or bounded with an
  overflow policy (`OverwriteOldest` or `RejectWhenFull`).
- `collections/heaps`: generic binary heap / priority queue (min, max or
  comparator), with `Drain` and `PushMany`.
- `collections/lru`: bounded LRU cache with eviction callbacks
  (`WithOnEvict`).
- `collections/multimaps`: `ListMultimap` and `SetMultimap`.
- `collections/dicts`: `Tree`, an ordered AVL-tree dict with
  `Floor`/`Ceiling`/`Range`/`Backward`, plus `UpdateInPlace`.
- `collections/sets`: `TreeSet`, an ordered set, and the full set algebra.
- Builders: `ListBuilder`, `QueueBuilder`, `StackBuilder` and `DequeBuilder`.
- `stats`:
  - `Sum`, `Product`, `Mean`, `WeightedMean`, `GeometricMean` and
    `HarmonicMean`;
  - `Median`, `Mode`, `Range` and `MinMax`;
  - `Quantile`/`Percentile`/`QuantileWith`, `IQR` and `PercentileOfScore`;
  - sample and population variance and standard deviation, and covariance;
  - `Correlation`, `LinearRegression`, `Skewness`, `Kurtosis`, `Entropy` and
    `Gini`;
  - `Histogram`, `Normalize`, `Standardize`, `MovingAverage`,
    `CumulativeSum` and `Clamp`;
  - `Dot`, `Norm`, `CosineSimilarity` and `EuclideanDistance`.
- `relational`: `GroupBy`, `GroupBySeq`, `CountBy`, `Aggregate`,
  `AggregateBy`, the inner, left, right and full outer joins, `Pivot`,
  `Unpivot` and `Partition`.
- `collections/sketches`: `MinHash`; `bloom` (`ApproxCount`); `countmin`;
  `hll`; and `tdigest` (mergeable streaming quantiles, `ConcurrentDigest`).
- `collections/streaming`: `TopK`, `Reservoir`, `WeightedReservoir`,
  `Bootstrap`, `BootstrapN`, `RunningMean`, `RunningVariance`, `EWMA` and
  `RunningMinMax`.
- `concurrency`: `Map`, `ForEach` and `Batch`, with `WithConcurrency` and
  `WithErrorPolicy` (`StopOnError`, `CollectErrors`, `ContinueOnError`).
- `channels`: `TumblingWindow`, `SlidingWindow`, `SessionWindow`,
  `WindowedReduce` and `CollectNAsSlice`.
- `slices`: `Chunk`, `Window`, `Zip`/`ZipWith`, `Unique`/`UniqueBy`,
  `FlatMap`, `Paginate` and `PushCopy`.
- `ml/distance`: `Euclidean`, `Manhattan`, `Minkowski`, `CosineDistance`,
  `Hamming` and `Levenshtein`.
- `ml/similarity`: `DotProduct`, `CosineSimilarity` and the set similarities
  `Jaccard`, `Dice` and `Overlap`.
- `ml/metrics`:
  - classification: accuracy, the confusion matrix, precision, recall and F1
    (macro, micro and weighted), ROC/AUC and log-loss;
  - clustering: silhouette;
  - ranking: DCG, NDCG, AP and MAP;
  - regression: MSE, RMSE, MAE, MAPE and R².
- `ml/preprocessing`:
  - `StandardScaler`, `MinMaxScaler` and `RobustScaler`;
  - the one-hot, ordinal, label and target encoders;
  - the mean, median, mode and constant imputers;
  - the fixed-width and quantile binners;
  - `TrainTestSplit`, `StratifiedSplit`, `KFold` and `Shuffle`, all seedable.
- Package READMEs and runnable `Example`s for every package. Every `doc.go`
  example is compile-checked in CI.

### Fixed

This release's hardening pass. Each fix ships with a regression test.

- **`collections/lists`:**
  - A circular `Linked` list's `DequeueInPlace` never emptied a one-element
    list, so drain loops hung.
  - Its `FilterInPlace` dropped a kept tail element and called the predicate
    twice on the head.
  - Builders branched from a common base overwrote each other's values.
  - `Array`'s in-place push, pop, enqueue and dequeue stay O(1) amortised.
  - Iterators range over a snapshot taken when ranging begins, so popping
    inside the loop no longer yields phantom zero values.
- **Concurrent collections:**
  - `FilterInPlace` is now exact when no write races the predicate.
    Previously, rejected NaN-like values were never removed, the concurrent
    list multimaps removed the wrong deep-equal binding, and the concurrent
    lists removed by value instead of position. When a write does race the
    predicate, the write is still preserved.
  - The concurrent sets' binary operations honour the other set's `Contains`.
- **`collections/dicts`, `collections/sets`, `collections/multimaps`:**
  - Mutating a `Tree`/`TreeSet` from inside `ForEach`/`All` no longer skips
    entries.
  - `ListMultimap.RemoveInPlace` no longer corrupts an in-flight iteration.
  - Every iterator follows one rule: it snapshots when ranging begins.
- **`collections/sketches`:**
  - `+0` and `-0` hash alike, ending Bloom false negatives and Count-Min
    under-counts.
  - Distinct seeds are independent hash functions; seeds 0–255 used to
    coincide.
  - `any`-typed keys of different dynamic types no longer collide.
  - Count-Min rows are independent for power-of-two widths, which restores
    the (ε, δ) bound.
  - Oversized configs return `ErrInvalidConfig` on 32-bit platforms instead
    of panicking.
- **`collections/sketches/tdigest`:**
  - `CDF` is correct at repeated values and at the minimum.
  - `Merge(self)` no longer triple-counts.
  - Infinite or overflowing weights are rejected (new `ErrWeightOverflow`).
  - Estimates stay finite near ±MaxFloat64.
- **`collections/streaming`:**
  - `WeightedReservoir` is unbiased at any weight scale.
  - One NaN no longer blocks `TopKOrdered`.
  - Counters no longer overflow on 32-bit platforms.
- **`concurrency`:** under `StopOnError`, the reported error is the failure
  that stopped the run, not a sibling's induced `context.Canceled`.
- **`channels`:**
  - Cancellation takes priority over a ready input in every stage. A cancelled
    stage used to consume and process one more element about half the time.
  - A huge window size no longer crashes the process.
- **`slices`:** `Paginate` no longer overflows for huge page arguments.
- **`stats`:**
  - `Skewness`/`Kurtosis` reject constant data that isn't exactly
    representable, such as `{0.1, 0.1, 0.1}`.
  - The means are overflow-safe and always lie within the data's [min, max].
  - The shared compensated accumulator (now Neumaier's variant) keeps an
    infinite total infinite instead of turning it into NaN.
  - These stay correct at extreme magnitudes: `Correlation`,
    `LinearRegression`'s R², `Norm` and `CosineSimilarity` (now clamped to
    [−1, 1]).
  - `Histogram` counts agree with its edges, and it no longer panics on an
    overflowing span.
  - `MovingAverage` no longer loses precision after a large value leaves the
    window.
  - `Normalize` keeps a NaN local to its position.
  - `Mode`/`Entropy`/`Gini` reject non-finite named float types.
  - `Dot` and `Norm` are 20–30% faster for vectors of 100 or more elements.
- **`ml`:**
  - `StandardScaler` now matches `stats.Standardize` exactly.
  - The binners reject non-finite values.
  - `StratifiedSplit` honours `testFrac` overall.
  - `CosineDistance` is never negative.
  - Fixed a heap regression in `SilhouetteScore`.
  - `SilhouetteSamplesWith` rejects a nil or misbehaving `DistanceFunc` rather
    than breaking its [−1, 1] guarantee (#231).
- **Cross-platform determinism:** the t-digest's rank arithmetic is protected
  against fused multiply-add contraction, so results are identical on amd64 and
  arm64. They used to differ for subnormal weights.

### Documentation

- `BackgroundWorkLimiter` documents its back-pressure and that calling `Add`
  from a work function on the same limiter can deadlock.
- `channels.CollectNAsSlice` documents that the caller owns upstream
  cancellation (#234).
- The linked lists, deques, LRU and builders document their real complexities
  and contracts.

### CI and tooling

- The test matrix covers both ends of the supported range: the Go 1.24 floor
  and the latest stable Go. The old "Go 1.23" cell never actually ran 1.23.
- Every pinned tool installs under `setup-go`'s `GOTOOLCHAIN=local`, and the
  GitHub Actions moved to their Node 24 majors.
- Weekly mutation testing now completes, with a per-test memory cap and
  calibrated timeouts. Baseline efficacy is 87.8% and mutator coverage 99.5%.
- New weekly deep-fuzz job with a persisted corpus.
- govulncheck runs against stable Go.
- The local `make ci` mirror ignores agent worktrees, and golangci-lint is
  v2.14.0.

## [0.0.2] - 2024-01-02

## [0.0.1] - 2023-08-07

Initial release: `channels`, `collections` (dicts, lists, sets), `constraints`,
`maps` and `slices`.

[Unreleased]: https://github.com/pickeringtech/go-collections/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/pickeringtech/go-collections/compare/v0.0.2...v0.1.0
[0.0.2]: https://github.com/pickeringtech/go-collections/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/pickeringtech/go-collections/releases/tag/v0.0.1
