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

- **Build & module hygiene**, **race + coverage tests** across an OS × Go-version matrix (#33), **lint/complexity** (`golangci-lint`, committed `.golangci.yml`), **security** (`govulncheck` + `gosec`) — all blocking.
- **Fuzzing:** `testing.F` targets across collections/utilities (#10), run as a fast count-based smoke step (#25).
- **Examples E2E:** separate-module example apps with golden-output assertions (#30) — blocking.
- **Report-only (ratcheting toward blocking):** API compatibility via `gorelease` (#29), benchmark-regression via `benchstat` (#31), mutation testing via `gremlins` (#32), and Codecov upload (#14).

## Planned (not yet in repo)

- **Iterator support:** Go 1.23+ `iter.Seq`/`iter.Seq2` accessors across the collections (#53).
- **Serialization:** `MarshalJSON`/`UnmarshalJSON` on the concrete collection types.
