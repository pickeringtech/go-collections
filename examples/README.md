# Examples

Small, focused, runnable apps that exercise the public API in realistic,
cross-package flows. They double as **living documentation** and as a
**downstream-consumer smoke test**: each app is built, run, and asserted against
a golden output in CI, so a change that breaks the public API — even one the
library's own internal tests would miss — fails the build.

## Why a separate module?

These examples live in their **own module** (`examples/go.mod`) so they consume
the library as a genuine *outside* package — only the exported API is reachable,
and they compile as a separate unit instead of recompiling in lockstep with the
library.

The module `replace`s its dependency with the local checkout (`../`), so CI
builds and tests the examples against the **current tree**. That validates the
public API still compiles and runs for a downstream consumer and powers the
end-to-end tests. (API compatibility for *released* versions is a separate
concern, covered by the gorelease gate.)

## The apps

| App                                            | Scenario                                                | Packages                       |
| ---------------------------------------------- | ------------------------------------------------------- | ------------------------------ |
| [`word-frequency`](./cmd/word-frequency)       | Tokenise text, count words, print the top-N             | `slices`, `maps`               |
| [`set-algebra`](./cmd/set-algebra)             | Union / intersection / difference / subset of two sets  | `collections`, `sets`, `slices`|
| [`worker-pipeline`](./cmd/worker-pipeline)     | Fan-out/fan-in a stream through a bounded worker pool   | `channels`, `concurrency`      |
| [`ordered-processing`](./cmd/ordered-processing)| Reverse (stack), replay (queue) and sort a list         | `lists`, `slices`              |

## Running an app

```bash
cd examples

echo "the quick brown fox the lazy dog the fox" | go run ./cmd/word-frequency -n 3
go run ./cmd/set-algebra -a apple,banana,cherry -b banana,cherry,fig
go run ./cmd/worker-pipeline -n 8 -workers 3
go run ./cmd/ordered-processing -nums 5,3,8,1,9,2
```

## Running the E2E tests

Each app produces deterministic output that is asserted against a checked-in
golden file in [`testdata/`](./testdata):

```bash
cd examples
go test ./...                # build, run and assert every app's stdout
go test ./... -update        # regenerate the golden files after an intentional change
```
