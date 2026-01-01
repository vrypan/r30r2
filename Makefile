# R30R2 RNG Makefile

# Binary names
R30R2_BIN = r30r2
COMPARE_READ_BIN = misc/compare-read
COMPARE_UINT64_BIN = misc/compare-uint64
VISUALIZE_BIN = misc/visualize-r30r2

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

# Build flags
LDFLAGS = -s -w
BUILD_FLAGS = -ldflags "$(LDFLAGS)"

# Source files for dependency tracking
R30R2_SOURCES = r30r2-main.go r30r2-cli.go rand/r30r2.go
COMPARE_READ_SOURCES = misc/compare-read.go rand/r30r2.go
COMPARE_UINT64_SOURCES = misc/compare-uint64.go rand/r30r2.go
VISUALIZE_SOURCES = misc/visualize-r30r2.go rand/r30r2.go

.PHONY: all compare visualize clean fmt help compare-run test-entropy smoke deps bench

# Default target
all: $(R30R2_BIN) compare

# Build the R30R2 CLI tool
$(R30R2_BIN): $(R30R2_SOURCES)
	@echo "Building $(R30R2_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(R30R2_BIN) r30r2-main.go r30r2-cli.go
	@echo "✓ Built $(R30R2_BIN)"

# Build both comparison tools
compare: $(COMPARE_READ_BIN) $(COMPARE_UINT64_BIN)

# Build the Read() comparison tool
$(COMPARE_READ_BIN): $(COMPARE_READ_SOURCES)
	@echo "Building $(COMPARE_READ_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_READ_BIN) misc/compare-read.go
	@echo "✓ Built $(COMPARE_READ_BIN)"

# Build the Uint64() comparison tool
$(COMPARE_UINT64_BIN): $(COMPARE_UINT64_SOURCES)
	@echo "Building $(COMPARE_UINT64_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_UINT64_BIN) misc/compare-uint64.go
	@echo "✓ Built $(COMPARE_UINT64_BIN)"

# Build the visualization tool
$(VISUALIZE_BIN): $(VISUALIZE_SOURCES)
	@echo "Building $(VISUALIZE_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(VISUALIZE_BIN) misc/visualize-r30r2.go
	@echo "✓ Built $(VISUALIZE_BIN)"

# Convenience alias for building visualization tool
visualize: $(VISUALIZE_BIN)

# Run comparison benchmarks
compare-run: compare
	@echo "Running Read() benchmark..."
	./$(COMPARE_READ_BIN)
	@echo ""
	@echo "Running Uint64() benchmark..."
	./$(COMPARE_UINT64_BIN)

# Run go test benchmarks with table output
bench:
	@./misc/bench-table.sh

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(R30R2_BIN)
	rm -f $(COMPARE_READ_BIN)
	rm -f $(COMPARE_UINT64_BIN)
	rm -f $(VISUALIZE_BIN)
	rm -f misc/stdlib-rng
	rm -f *.prof
	rm -f *.test
	rm -f *.bin
	rm -f *.dat
	@echo "✓ Cleaned"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Test randomness with ent - compares all three RNGs
test-entropy: r30r2
	@./misc/test-entropy.sh

# Quick smoke test
smoke: r30r2
	@echo "Running smoke test..."
	@./$(R30R2_BIN) --seed=12345 --bytes=1024 > /dev/null
	@echo "✓ Smoke test passed"

# Show help
help:
	@echo "R30R2 RNG Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            Build all binaries (default)"
	@echo "  r30r2          Build r30r2 CLI tool"
	@echo "  compare        Build comparison tools (read + uint64)"
	@echo "  compare-read   Build compare-read tool (MB/s benchmark)"
	@echo "  compare-uint64 Build compare-uint64 tool (ns/call benchmark)"
	@echo "  compare-run    Run both comparison benchmarks"
	@echo "  visualize      Build R30R2 visualization tool"
	@echo "  bench          Run go test benchmarks (table format)"
	@echo "  fmt            Format code with gofmt"
	@echo "  clean          Remove build artifacts"
	@echo "  deps           Download and tidy dependencies"
	@echo "  test-entropy   Test randomness with ent tool"
	@echo "  smoke          Quick smoke test"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make r30r2"
	@echo "  make compare-run"
	@echo "  make test-entropy"
	@echo "  make clean r30r2"
