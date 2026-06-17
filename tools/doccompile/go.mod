// doccompile lives in its OWN module (like examples/ and tools/benchreport) so
// it never touches the library's root module: it adds no dependency to go.mod,
// sits outside the root's 100% coverage floor / lint / gosec gates, and is
// invoked only via `make doc-compile` and the PR-time CI job. It is
// intentionally dependency-free — it parses the library's package APIs and the
// godoc code blocks with the standard library alone, and shells out to the
// installed `go` toolchain to compile the parseable blocks. See issue #151.
module github.com/pickeringtech/go-collections/tools/doccompile

go 1.24
