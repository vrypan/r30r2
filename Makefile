# Rule 30 RNG Makefile

# Binary names
RULE30_BIN = rule30-rng
COMPARE_BIN = rule30-compare

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
COMPARE_SOURCES = compare.go

.PHONY: all rule30 compare clean fmt help compare-run test-entropy smoke deps

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

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Test randomness with ent - compares all three RNGs
test-entropy: rule30
	@./test-entropy.sh

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
