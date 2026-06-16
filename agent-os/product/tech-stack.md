# Tech Stack

go-collections is a pure-Go library with zero runtime dependencies.

## Frontend

N/A — this is a library, not an application.

## Backend

- **Language:** Go 1.24
- **Generics:** full use of type parameters for type-safe collections and utilities
- **Concurrency:** `sync.Mutex` / `sync.RWMutex` for thread-safe variants; `sync.WaitGroup` + buffered-channel semaphores for work limiters
- **Runtime dependencies:** none (standard library only)

## Database

N/A.

## Testing & Tooling

- **Testing:** standard library `testing` — table-driven tests, GoDoc `Example` functions, `Benchmark` functions (scaling ladder), and fuzz targets (`testing.F`). Run with `-race` for concurrent types.
- **Formatting / vetting:** `gofmt`, `go vet` (standard Go toolchain).

These are development/CI dependencies only — they do not affect the zero **runtime** dependency guarantee:

- **Linting / static analysis:** `golangci-lint` (bundles staticcheck, cyclop, gocognit, etc.).
- **Security:** `govulncheck` (vulnerability scanning), `gosec` (security lint).
- **API compatibility:** `gorelease` (`golang.org/x/exp/cmd/gorelease`) compares the working tree's exported API against the last released semver tag and reports the implied version bump plus any breaking changes. Report-only pre-1.0 (breaking changes are permitted while major version is 0); planned to gate post-v1.0. See [#29](https://github.com/pickeringtech/go-collections/issues/29).

## Distribution & Docs

- **Module path:** `github.com/pickeringtech/go-collections`, distributed via the Go module proxy.
- **API reference:** published on [pkg.go.dev](https://pkg.go.dev/github.com/pickeringtech/go-collections) / GoDoc, driven by per-package `doc.go` files.

## CI / verification (in repo)

GitHub Actions on push/PR, fronted by a single stable `CI Gate` aggregator that is
the only required branch-protection check (#41), so matrix/job changes never wedge PRs:

- **Blocking** (in `CI Gate`'s `needs:`): **build & module hygiene**, **race + coverage tests** across an OS × Go-version matrix (#33), **lint/complexity** (`golangci-lint`, committed `.golangci.yml`), **security** (`govulncheck` + `gosec`), **examples E2E** (separate-module apps with golden-output assertions, #30), **cross-arch** (386 build+vet+test, arm64/s390x build+vet, #33), and **fuzzing** (`testing.F` targets, deterministic count-based smoke run, #10/#25). Cross-arch and fuzz were promoted from report-only to gates (#72) — both are deterministic, so a failure is a real signal, not runner noise.
- **Report-only** (not in the gate; `continue-on-error`, each with a documented reason): **Go tip** (intentionally unstable toolchain, #33), **API compatibility** via `gorelease` (pre-1.0, #29), **benchmark-regression** via `benchstat` (noisy shared runners, #31), **mutation testing** via `gremlins` (threshold being triaged, #32), and **Codecov** upload (needs `CODECOV_TOKEN`, #14). These ratchet toward blocking as they prove out (#72).

## Planned (not yet in repo)

- **Iterator support:** Go 1.23+ `iter.Seq`/`iter.Seq2` accessors across the collections (#53).
- **Serialization:** `MarshalJSON`/`UnmarshalJSON` on the concrete collection types.
