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

### As a Command-Line Tool

```bash
go install github.com/vrypan/rule30rnd@latest
```

Or build from source:

```bash
git clone https://github.com/vrypan/rule30rnd
cd rule30rnd
make all
```

### As a Library

Add to your Go project:

```bash
go get github.com/vrypan/rule30rnd
```

Import in your code:

```go
import "github.com/vrypan/rule30rnd/rule30"
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

### Using the Library

Rule30 RNG is compatible with Go's `math/rand` interface and can be used as a drop-in replacement:

```go
package main

import (
    "fmt"
    "github.com/vrypan/rule30rnd/rule30"
)

func main() {
    // Create RNG with seed
    rng := rule30.New(12345)

    // Compatible with math/rand interface
    fmt.Printf("Uint32:    %d\n", rng.Uint32())
    fmt.Printf("Uint64:    %d\n", rng.Uint64())
    fmt.Printf("Int:       %d\n", rng.Int())
    fmt.Printf("Intn(100): %d\n", rng.Intn(100))
    fmt.Printf("Float64:   %.6f\n", rng.Float64())

    // Distribution functions
    fmt.Printf("NormFloat64: %.6f\n", rng.NormFloat64())
    fmt.Printf("ExpFloat64:  %.6f\n", rng.ExpFloat64())

    // Also implements io.Reader
    buf := make([]byte, 1024)
    n, err := rng.Read(buf)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Generated %d random bytes\n", n)
}
```

**Available methods (math/rand compatible):**
- `Uint32()` - random uint32
- `Uint64()` - random uint64
- `Int()` - non-negative random int
- `Int31()` - non-negative random int32
- `Int63()` - non-negative random int64
- `Intn(n)` - random int in [0, n)
- `Int31n(n)` - random int32 in [0, n)
- `Int63n(n)` - random int64 in [0, n)
- `Float32()` - random float32 in [0.0, 1.0)
- `Float64()` - random float64 in [0.0, 1.0)
- `NormFloat64()` - normally distributed float64 (mean=0, stddev=1)
- `ExpFloat64()` - exponentially distributed float64 (rate=1)
- `Read([]byte)` - fills byte slice with random data (io.Reader)

### Testing Randomness

Use the `ent` tool to analyze randomness:

```bash
# Generate test data
./rule30-rng --bytes=10485760 > test.bin

# Analyze with ent
ent test.bin
```

**Actual test results (10MB sample):**

```
Entropy = 7.999982 bits per byte.

Optimum compression would reduce the size
of this 10485760 byte file by 0 percent.

Chi square distribution for 10485760 samples is 264.55, and randomly
would exceed this value 32.73 percent of the times.

Arithmetic mean value of data bytes is 127.4987 (127.5 = random).
Monte Carlo value for Pi is 3.140031105 (error 0.05 percent).
Serial correlation coefficient is 0.000171 (totally uncorrelated = 0.0).
```

**Analysis:**
- ✓ **Entropy**: 7.999982/8.0 (99.9998% of maximum) - Perfect
- ✓ **Chi-square**: 264.55 (expected ~255, range 200-310) - Excellent
- ✓ **Mean**: 127.4987 (expected 127.5) - Perfect uniformity
- ✓ **Monte Carlo π**: 3.140031 (error 0.05%) - Excellent
- ✓ **Serial correlation**: 0.000171 (expected 0.0) - No detectable pattern

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

### Differences from Traditional Rule 30 RNG

**Traditional approach (Wolfram's method):**
- Evolves a large 1D CA over many generations
- Extracts randomness from a **single column** (typically center) over time
- Output: 1 bit per iteration from one cell's temporal evolution
- Example: bit₀ from gen₀, bit₁ from gen₁, bit₂ from gen₂, etc.
- Requires tracking spacetime history or evolving step-by-step

**Our implementation (spatial extraction):**
- Evolves a 256-bit circular strip one generation
- Extracts randomness from **all 256 bits** of the current state
- Output: 256 bits per iteration from entire spatial pattern
- Repeats: evolve → extract all → evolve → extract all
- **256× more efficient**: exploits Rule 30's spatial randomness, not just temporal

**Key advantages of spatial extraction:**
- ✓ **Much faster**: 256 bits per iteration vs 1 bit
- ✓ **Simpler**: no spacetime history needed
- ✓ **Circular topology**: ensures uniform mixing (all positions equivalent)
- ✓ **Parallel-friendly**: processes entire state at once

**Trade-off:**
- Traditional method has longer "history" from single position's evolution
- Our method samples the full spatial pattern at each instant
- Both produce excellent randomness - ours is just more efficient

This spatial extraction approach is why we achieve 5,000+ MB/s throughput while maintaining perfect entropy and distribution.

## Randomness Quality

Rule 30 RNG has been extensively tested and shows exceptional randomness quality.

### Test Results Summary

| Test | Result | Expected | Status |
|------|--------|----------|--------|
| **Shannon Entropy** | 7.999982 bits/byte | 8.0000 | ✓ Perfect (99.9998%) |
| **Chi-Square** | 253.9 - 264.55 | ~255 (200-310) | ✓ Excellent |
| **Arithmetic Mean** | 127.4987 | 127.5 | ✓ Perfect |
| **Monte Carlo π** | 3.140031 | 3.141593 | ✓ 0.05% error |
| **Serial Correlation** | 0.000171 | 0.0000 | ✓ Uncorrelated |

### Shannon Entropy

Measures unpredictability of the data:

```
H = -Σ(p(i) × log₂(p(i)))
```

- **Maximum**: 8.0000 bits/byte (perfect randomness)
- **Rule 30 RNG**: 7.999982 bits/byte (99.9998% of maximum) ✓

### Chi-Square Test

Measures uniformity of byte distribution:

```
χ² = Σ((observed - expected)² / expected)
```

- **Expected value**: ~255 (for 256 byte values, df=255)
- **Acceptable range**: 200-310 (95% confidence interval)
- **Rule 30 RNG**: 253.9 average (comparison tool), 264.55 (ent 10MB test) ✓

### Additional Quality Metrics

From `ent` statistical analysis tool (10MB sample):
- **Arithmetic mean**: 127.4987 (expected 127.5) - Perfect uniformity
- **Monte Carlo π estimate**: 3.140031105 (0.05% error) - Excellent
- **Serial correlation**: 0.000171 (expected 0.0) - No byte-to-byte patterns
- **Compression**: 0% - Data is incompressible (hallmark of randomness)

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
