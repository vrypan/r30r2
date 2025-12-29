# Turmite RNG Makefile

# Binary names
MAIN_BIN = turmite-rng
COMPARE_BIN = turmite-compare

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
MAIN_SOURCES = main.go turmite.go rng.go
COMPARE_SOURCES = compare.go turmite.go rng.go
TEST_SOURCES = benchmark_test.go

.PHONY: all build compare test bench clean fmt help install compare-run

# Default target
all: build compare

# Build the main CLI tool
build: $(MAIN_BIN)

$(MAIN_BIN): $(MAIN_SOURCES)
	@echo "Building $(MAIN_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(MAIN_BIN) $(MAIN_SOURCES)
	@echo "✓ Built $(MAIN_BIN)"

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
	rm -f $(MAIN_BIN)
	rm -f $(COMPARE_BIN)
	rm -f *.prof
	rm -f *.test
	rm -f *.bin
	rm -f *.dat
	@echo "✓ Cleaned"

# Install binaries to GOPATH/bin
install: build compare
	@echo "Installing binaries..."
	cp $(MAIN_BIN) $(GOPATH)/bin/
	cp $(COMPARE_BIN) $(GOPATH)/bin/
	@echo "✓ Installed to $(GOPATH)/bin/"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Generate random test data
testdata: $(MAIN_BIN)
	@echo "Generating test data (1MB)..."
	./$(MAIN_BIN) --bytes=1048576 > testdata.bin
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
smoke: build
	@echo "Running smoke test..."
	@./$(MAIN_BIN) --seed=12345 --bytes=1024 > /dev/null
	@echo "✓ Smoke test passed"

# Show help
help:
	@echo "Turmite RNG Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            Build all binaries (default)"
	@echo "  build          Build turmite-rng CLI tool"
	@echo "  compare        Build turmite-compare tool"
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
	@echo "  make build"
	@echo "  make compare-run"
	@echo "  make bench"
	@echo "  make clean build"
