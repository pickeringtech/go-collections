# Makefile — developer entry points. The benchmark-report pipeline (issue #50)
# lives here so the main-only CI job can stay a thin wrapper over `make
# bench-report`, keeping local and CI generation in lockstep (no logic drift).

# Bench knobs — mirror the CI env in .github/workflows/ci.yml. Override on the
# command line for a fast local preview, e.g.:
#   make bench-report BENCH_TIME=10ms BENCH_COUNT=2
BENCH_TIME  ?= 50ms
BENCH_COUNT ?= 8
# Scope the report to the standardized collections matrix (dicts/lists/sets).
BENCH_PKGS  ?= ./collections/...

# benchstat from PATH by default; CI installs the pinned version first.
BENCHSTAT ?= benchstat

BENCHREPORT_DIR := tools/benchreport
BUILD_DIR       := build

# Provenance, computed once so the generator stays a pure function of its inputs.
GIT_SHA    := $(shell git rev-parse --short HEAD 2>/dev/null)
GEN_DATE   := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go env GOVERSION)

.PHONY: help test bench bench-report

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

test: ## Run the library test suite with race + shuffle
	go test -race -shuffle=on ./...

bench: ## Run the collections benchmarks once (no report), printing results
	go test -run='^$$' -bench=. -benchmem -benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT) $(BENCH_PKGS)

bench-report: ## Run benchmarks and regenerate BENCHMARKS.md, docs/bench.svg, and the README preview
	@mkdir -p $(BUILD_DIR)
	@echo ">> benchmarking $(BENCH_PKGS) (-benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT))"
	go test -run='^$$' -bench=. -benchmem -benchtime=$(BENCH_TIME) -count=$(BENCH_COUNT) $(BENCH_PKGS) \
		| tee $(BUILD_DIR)/bench.txt
	@echo ">> summarising with benchstat"
	$(BENCHSTAT) -format=csv $(BUILD_DIR)/bench.txt > $(BUILD_DIR)/bench.csv
	@echo ">> generating report, chart, and README preview"
	go build -C $(BENCHREPORT_DIR) -o $(CURDIR)/$(BUILD_DIR)/benchreport .
	$(CURDIR)/$(BUILD_DIR)/benchreport \
		-csv $(BUILD_DIR)/bench.csv \
		-commit "$(GIT_SHA)" \
		-date "$(GEN_DATE)" \
		-goversion "$(GO_VERSION)" \
		-benchtime "$(BENCH_TIME)" \
		-count "$(BENCH_COUNT)"
