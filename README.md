# Rule 30 RNG

A high-performance random number generator based on Rule 30 cellular automaton, implemented in Go.

## Overview

Rule 30 RNG generates pseudo-random numbers using a 1D cellular automaton (Rule 30) on a circular 256-bit strip. Rule 30 is known for producing high-quality randomness and is famously used in Mathematica's default random number generator.

**Key Features:**
- **Extremely fast**: 5,213 MB/s average throughput (4.2× faster than math/rand, 1.5× faster than crypto/rand)
- **Perfect entropy**: 8.0000 bits/byte (maximum possible)
- **Excellent distribution**: Chi-square 253.9 (nearly ideal uniform distribution)
- **Deterministic**: Same seed always produces same output
- **Simple implementation**: ~110 lines of optimized Go code

## Performance

Benchmark results on Apple Silicon (M-series):

| RNG         | Throughput | vs Rule30 | Entropy | Chi-Square |
|-------------|------------|-----------|---------|------------|
| Rule30RNG   | 5,213 MB/s | 1.0×      | 8.0000  | 253.9      |
| crypto/rand | 3,515 MB/s | 0.67×     | 8.0000  | 245.7      |
| math/rand   | 1,234 MB/s | 0.24×     | 8.0000  | 272.3      |

*All values show excellent randomness quality (entropy = 8.0, chi-square ≈ 255)*

## Installation

```bash
go install github.com/yourusername/rule30-rng@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/rule30-rng
cd rule30-rng
make all
```

## Usage

### Command Line

Generate random bytes:

```bash
# Generate 1MB of random data
./rule30-rng --bytes=1048576 > random.bin

# Use specific seed for reproducibility
./rule30-rng --seed=12345 --bytes=1024 > random.bin

# Benchmark throughput
./rule30-rng --benchmark
```

### As a Library

```go
package main

import (
    "fmt"
    "io"
)

func main() {
    // Create RNG with seed
    rng := NewRule30(12345)

    // Read random bytes (implements io.Reader)
    buf := make([]byte, 1024)
    n, err := rng.Read(buf)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated %d random bytes\n", n)
}
```

### Testing Randomness

Use the `ent` tool to analyze randomness:

```bash
# Generate test data
./rule30-rng --bytes=1048576 > test.bin

# Analyze with ent
ent test.bin
```

Expected output should show:
- Entropy: ~8.000000 bits per byte
- Chi-square: ~255 (acceptable range: 200-310)
- Monte Carlo π approximation: ~3.14

## How It Works

### Rule 30 Cellular Automaton

Rule 30 is a 1D cellular automaton where each cell evolves based on itself and its two neighbors:

```
new_bit = left XOR (center OR right)
```

### Implementation Details

1. **256-bit circular strip**: Represented as 4 × 64-bit words (`[4]uint64`)
2. **Parallel processing**: Processes 64 bits simultaneously using bitwise operations
3. **Fully unrolled loop**: Each of 4 words updated in single operation (no loops)
4. **Single iteration extraction**: Extracts 32 bytes after each evolution step

**Architecture-specific optimizations:**
- Uses native 64-bit word size on 64-bit architectures
- Single-cycle bitwise operations (XOR, OR, shifts)
- Minimized memory allocations with pre-allocated buffers

### Why It's Fast

1. **Bit-level parallelism**: 64 bits updated per operation (vs 1 bit serially)
2. **Loop unrolling**: Zero loop overhead for state evolution
3. **Immediate extraction**: No wasted iterations (generates 32 bytes per step)
4. **Simple operations**: Only XOR, OR, and bit shifts (1 CPU cycle each)

## Randomness Quality

### Shannon Entropy

Measures unpredictability of the data:

```
H = -Σ(p(i) × log₂(p(i)))
```

- **Maximum**: 8.0000 bits/byte (perfect randomness)
- **Rule 30 RNG**: 8.0000 bits/byte ✓

### Chi-Square Test

Measures uniformity of byte distribution:

```
χ² = Σ((observed - expected)² / expected)
```

- **Expected value**: ~255 (for 256 byte values, df=255)
- **Acceptable range**: 200-310 (95% confidence interval)
- **Rule 30 RNG**: 253.9 ✓

## Building

```bash
# Build all binaries
make all

# Build just the RNG CLI
make rule30

# Build comparison tool
make compare

# Run comparison benchmarks
make compare-run

# Run tests
make test

# Run Go benchmarks
make bench

# Generate test data
make testdata

# Clean build artifacts
make clean
```

## Comparison Tool

Compare Rule 30 RNG against Go's standard libraries:

```bash
./rule30-compare
```

This runs comprehensive benchmarks measuring:
- Throughput (MB/s) at different buffer sizes
- Shannon entropy (bits per byte)
- Chi-square distribution test

## Use Cases

**Good for:**
- Monte Carlo simulations
- Procedural generation (games, graphics)
- Non-cryptographic random sampling
- Reproducible random sequences (deterministic from seed)
- High-throughput random data generation

**Not recommended for:**
- Cryptographic applications (use crypto/rand instead)
- Security-critical random number generation
- Applications requiring cryptographically secure randomness

## Algorithm Background

Rule 30 was discovered by Stephen Wolfram in 1983 as one of 256 elementary cellular automata. Despite its simple definition, it exhibits chaotic behavior and generates seemingly random patterns.

**Key properties:**
- Class III cellular automaton (chaotic behavior)
- Sensitivity to initial conditions
- No known compact mathematical description
- Passes many statistical randomness tests

**Historical note**: Wolfram Research uses Rule 30 as the basis for Mathematica's default random number generator.

## References

- Wolfram, Stephen (1983). "Statistical mechanics of cellular automata"
- Wolfram, Stephen (2002). "A New Kind of Science"
- [Rule 30 on Wikipedia](https://en.wikipedia.org/wiki/Rule_30)

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please open an issue or pull request.

Areas for contribution:
- Additional randomness tests (Diehard, TestU01)
- SIMD optimizations (AVX2, NEON)
- Platform-specific builds
- Extended seed formats
- Additional CLI features
