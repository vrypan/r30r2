# Rule 30 RNG Makefile

# Binary names
RULE30_BIN = rule30-rng
COMPARE_BIN = rule30-compare

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

# Build flags
LDFLAGS = -s -w
BUILD_FLAGS = -ldflags "$(LDFLAGS)"

# Source files
RULE30_SOURCES = rule30-main.go rule30-cli.go rule30.go
COMPARE_SOURCES = compare.go rule30.go
TEST_SOURCES = benchmark_test.go

.PHONY: all rule30 compare test bench clean fmt help install compare-run

# Default target
all: rule30 compare

# Build the Rule 30 CLI tool
rule30: $(RULE30_BIN)

$(RULE30_BIN): $(RULE30_SOURCES)
	@echo "Building $(RULE30_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(RULE30_BIN) $(RULE30_SOURCES)
	@echo "✓ Built $(RULE30_BIN)"

# Build the comparison tool
compare: $(COMPARE_BIN)

$(COMPARE_BIN): $(COMPARE_SOURCES)
	@echo "Building $(COMPARE_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_BIN) $(COMPARE_SOURCES)
	@echo "✓ Built $(COMPARE_BIN)"

# Run comparison benchmarks
compare-run: $(COMPARE_BIN)
	@echo "Running performance comparison..."
	./$(COMPARE_BIN)

# Run Go tests
test:
	@echo "Running tests..."
	$(GOTEST) -v

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem

# Run benchmarks with CPU profiling
bench-profile:
	@echo "Running benchmarks with profiling..."
	$(GOTEST) -bench=. -benchmem -cpuprofile=cpu.prof
	@echo "View profile with: go tool pprof cpu.prof"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(RULE30_BIN)
	rm -f $(COMPARE_BIN)
	rm -f *.prof
	rm -f *.test
	rm -f *.bin
	rm -f *.dat
	@echo "✓ Cleaned"

# Install binaries to GOPATH/bin
install: rule30 compare
	@echo "Installing binaries..."
	cp $(RULE30_BIN) $(GOPATH)/bin/
	cp $(COMPARE_BIN) $(GOPATH)/bin/
	@echo "✓ Installed to $(GOPATH)/bin/"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Generate random test data
testdata: $(RULE30_BIN)
	@echo "Generating test data (1MB)..."
	./$(RULE30_BIN) --bytes=1048576 > testdata.bin
	@echo "✓ Generated testdata.bin (1MB)"

# Test randomness with ent (if available)
test-entropy: testdata
	@if command -v ent >/dev/null 2>&1; then \
		echo "Testing entropy with ent..."; \
		ent testdata.bin; \
	else \
		echo "ent not installed. Install with: brew install ent"; \
	fi

# Quick smoke test
smoke: rule30
	@echo "Running smoke test..."
	@./$(RULE30_BIN) --seed=12345 --bytes=1024 > /dev/null
	@echo "✓ Smoke test passed"

# Show help
help:
	@echo "Rule 30 RNG Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            Build all binaries (default)"
	@echo "  rule30         Build rule30-rng CLI tool"
	@echo "  compare        Build rule30-compare tool"
	@echo "  compare-run    Run performance comparison"
	@echo "  test           Run Go tests"
	@echo "  bench          Run benchmarks"
	@echo "  bench-profile  Run benchmarks with CPU profiling"
	@echo "  fmt            Format code with gofmt"
	@echo "  clean          Remove build artifacts"
	@echo "  install        Install binaries to GOPATH/bin"
	@echo "  deps           Download and tidy dependencies"
	@echo "  testdata       Generate 1MB test file"
	@echo "  test-entropy   Test randomness with ent tool"
	@echo "  smoke          Quick smoke test"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make rule30"
	@echo "  make compare-run"
	@echo "  make bench"
	@echo "  make clean rule30"
