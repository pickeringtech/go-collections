#!/usr/bin/env bash
# memcap-exec.sh — `go test -exec` wrapper that runs a test binary under a
# per-process address-space cap (ulimit -v), then execs it unchanged.
#
# Why: mutation testing (.github/workflows/mutation.yml) deliberately breaks the
# code. A mutant that turns a bounded loop into an unbounded one — e.g. a
# linked-list traversal step or a slice-growth bound — can allocate without
# limit. Several such test processes running in parallel exhaust the runner's
# RAM, and the kernel then takes down the whole runner ("The runner has received
# a shutdown signal"), losing the entire run. Capped, the runaway test instead
# dies with "fatal error: out of memory" and gremlins records the mutant as
# caught, so the run completes.
#
# Only the test binaries are capped: the cap is applied via `go test -exec`, so
# the compiler, linker and gremlins itself (which needs address space for many
# threads) run unconstrained.
#
# Usage:  GOFLAGS="-exec=$PWD/.github/scripts/memcap-exec.sh" gremlins unleash .
# Tune:   MEMCAP_KB (default 3 GiB) — comfortably above the suite's real peak.
set -euo pipefail
ulimit -v "${MEMCAP_KB:-3145728}"
exec "$@"
