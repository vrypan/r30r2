# Rule 30 RNG Makefile

# Binary names
RULE30_BIN = rule30
COMPARE_READ_BIN = misc/compare-read
COMPARE_UINT64_BIN = misc/compare-uint64
VISUALIZE_BIN = misc/visualize-rule30

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

# Build flags
LDFLAGS = -s -w
BUILD_FLAGS = -ldflags "$(LDFLAGS)"

# Source files
RULE30_SOURCES = rule30-main.go rule30-cli.go
COMPARE_READ_SOURCES = misc/compare-read.go
COMPARE_UINT64_SOURCES = misc/compare-uint64.go
VISUALIZE_SOURCES = misc/visualize-rule30.go

.PHONY: all rule30 compare compare-read compare-uint64 visualize clean fmt help compare-run test-entropy smoke deps bench

# Default target
all: rule30 compare

# Build the Rule 30 CLI tool
rule30: $(RULE30_BIN)

$(RULE30_BIN): $(RULE30_SOURCES)
	@echo "Building $(RULE30_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(RULE30_BIN) $(RULE30_SOURCES)
	@echo "✓ Built $(RULE30_BIN)"

# Build both comparison tools
compare: compare-read compare-uint64

# Build the Read() comparison tool
compare-read: $(COMPARE_READ_BIN)

$(COMPARE_READ_BIN): $(COMPARE_READ_SOURCES)
	@echo "Building $(COMPARE_READ_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_READ_BIN) $(COMPARE_READ_SOURCES)
	@echo "✓ Built $(COMPARE_READ_BIN)"

# Build the Uint64() comparison tool
compare-uint64: $(COMPARE_UINT64_BIN)

$(COMPARE_UINT64_BIN): $(COMPARE_UINT64_SOURCES)
	@echo "Building $(COMPARE_UINT64_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(COMPARE_UINT64_BIN) $(COMPARE_UINT64_SOURCES)
	@echo "✓ Built $(COMPARE_UINT64_BIN)"

# Build the visualization tool
visualize: $(VISUALIZE_BIN)

$(VISUALIZE_BIN): $(VISUALIZE_SOURCES)
	@echo "Building $(VISUALIZE_BIN)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(VISUALIZE_BIN) $(VISUALIZE_SOURCES)
	@echo "✓ Built $(VISUALIZE_BIN)"

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
	rm -f $(RULE30_BIN)
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
test-entropy: rule30
	@./misc/test-entropy.sh

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
	@echo "  rule30         Build rule30 CLI tool"
	@echo "  compare        Build comparison tools (read + uint64)"
	@echo "  compare-read   Build compare-read tool (MB/s benchmark)"
	@echo "  compare-uint64 Build compare-uint64 tool (ns/call benchmark)"
	@echo "  compare-run    Run both comparison benchmarks"
	@echo "  visualize      Build Rule 30 visualization tool"
	@echo "  bench          Run go test benchmarks (table format)"
	@echo "  fmt            Format code with gofmt"
	@echo "  clean          Remove build artifacts"
	@echo "  deps           Download and tidy dependencies"
	@echo "  test-entropy   Test randomness with ent tool"
	@echo "  smoke          Quick smoke test"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make rule30"
	@echo "  make compare-run"
	@echo "  make test-entropy"
	@echo "  make clean rule30"
