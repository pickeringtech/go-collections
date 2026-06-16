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

## Distribution & Docs

- **Module path:** `github.com/pickeringtech/go-collections`, distributed via the Go module proxy.
- **API reference:** published on [pkg.go.dev](https://pkg.go.dev/github.com/pickeringtech/go-collections) / GoDoc, driven by per-package `doc.go` files.

## Planned (not yet in repo)

- **CI:** GitHub Actions running build, test (`-race`), lint, complexity, fuzz, and vulnerability checks on push/PR (#7).
- **Linting:** committed `.golangci.yml` configuration for development.
- **Fuzzing:** `testing.F` fuzz targets across the collections and utility packages (#10).
