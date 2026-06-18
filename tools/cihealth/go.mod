// cihealth lives in its OWN module (like examples/ and tools/benchreport/) so it
// never touches the library's root module: it adds no dependency to go.mod, sits
// outside the root's 100% coverage floor / lint / gosec gates, and is invoked
// only via `make ci-health-report` and the scheduled ci-health-badges workflow.
// It is intentionally dependency-free — it consumes the CI run history GitHub
// hands it and emits shields.io endpoint JSON using only the standard library.
module github.com/pickeringtech/go-collections/tools/cihealth

go 1.24
