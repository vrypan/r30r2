# R30R2

High-performance pseudo-random number generator based on Rule 30 (radius-2) cellular automaton. 
- **2× faster** than Go's math/rand/v2 
- **160/160 BigCrush** tests pass perfectly.

## Quick Start

### Installation

```bash
# As a command-line tool
go install github.com/vrypan/r30r2@latest

# As a library
go get github.com/vrypan/r30r2
```

### Command Line Usage

```bash
# Generate random data
./r30r2 --bytes=1048576 > random.bin

# Generate specific size with dd
./r30r2 --bytes=1073741824 | dd of=test.data bs=1m

# Unlimited streaming (use with head, pv, or Ctrl+C)
./r30r2 --bytes=0 | head -c 1073741824 > test.data

# Reproducible output with seed
./r30r2 --seed=12345 --bytes=1024 > random.bin
```

### Library Usage

Drop-in replacement for math/rand:

```go
import "github.com/vrypan/r30r2/rand"

rng := rand.New(12345)

// Full math/rand interface
rng.Uint64()           // Random uint64
rng.Intn(100)          // Random int in [0, 100)
rng.Float64()          // Random float64 in [0.0, 1.0)
rng.NormFloat64()      // Normal distribution (mean=0, stddev=1)

// Also implements io.Reader
buf := make([]byte, 1024)
rng.Read(buf)
```

**API**: Compatible with `math/rand` - all methods supported (Uint32/64, Int/Intn, Float32/64, NormFloat64, ExpFloat64, Read).

## Performance

Benchmarks on Apple M4, verified 2026-01-01 (run `make bench` to reproduce):

### Throughput Comparison

|Algorithm       |  Read() 32KB |   Read() 1KB |     Uint64()|
|----------------|--------------|--------------|-------------|
|math/rand       |  21316.00 ns |    674.30 ns |      1.81 ns|
|math/rand/v2    |  13369.00 ns |    423.50 ns |      3.22 ns|
|**R30R2**       |   **5516.00 ns** |    **183.90 ns** |      **1.75 ns**|
|crypto/rand     |   7009.00 ns |    367.90 ns |     56.29 ns|

**Relative to math/rand:**

|Algorithm       | Read() 32KB |  Read() 1KB |    Uint64()|
|----------------|-------------|-------------|------------|
|math/rand       |       1.00x |       1.00x |       1.00x|
|math/rand/v2    |       1.59x |       1.59x |       0.56x|
|**R30R2**       |    **3.86x** |   **3.67x** |   **1.03x**|
|crypto/rand     |       3.04x |       1.83x |       0.03x|

## Randomness Quality

**Perfect scores on complete TestU01 test suite** - exceptional statistical quality verified through rigorous testing:

| Test Battery | Tests | R30R2 | Pass Rate | Status |
|--------------|-------|--------|-----------|------|
| **SmallCrush** | 15 | **15/15** ✓ | **100%** | ✅ Verified 2026-01-01 |
| **Crush** | 144 | **144/144** ✓ | **100%**  | ✅ Verified 2026-01-01 |
| **BigCrush** | 160 | **160/160** ✓ | **100%** | ✅ Verified 2026-01-01 |

**Historic Achievement:** The implementation demonstrates exceptional quality across all test categories: serial correlation, birthday spacing, collision, permutation, matrix rank, spectral, string, compression, random walk, Fourier analysis, linear complexity, and autocorrelation tests.

### Comparison with Go Standard Library

BigCrush results comparing Rule 30 with Go's stdlib RNGs (verified 2026-01-01):

| Generator | BigCrush Score | Performance vs math/rand |
|-----------|---------------|--------------------------|
| **math/rand** | 159/160* | 1.00× (baseline) |
| **math/rand/v2** | **160/160** ✓ | 0.56× |
| **Rule 30** | **160/160** ✓ | **1.03×** |

*math/rand fails test #37 (Gap, r = 20) with p-value 5.6e-16

**Key insight:** Rule 30 delivers **equivalent statistical quality** to math/rand/v2 (both perfect 160/160) while being **2x faster**.

To run TestU01 tests:

```bash
cd testu01
make smallcrush   # Quick test (15 tests, 2 min)
make crush        # Medium test (144 tests, 10 min)
make bigcrush     # Comprehensive test (160 tests, 1 hour)
```

## How It Works

**Algorithm**: 
- Radius-2 cellular automaton based on Rule 30 on a circular 256-bit strip.

  `new_bit = (left2 XOR left1) XOR ((center OR right1) OR right2)`
- Hybrid rotation + multiply mixing applied on output

```
Evolution rule: new_bit = (left2 XOR left1) XOR ((center OR right1) OR right2)
```

**Implementation**:
- 256-bit state as 4 × 64-bit words
- Radius-2 neighborhood: each bit depends on 5 neighboring cells
- Processes 64 bits in parallel per word (bitwise operations)
- Fully unrolled loop - all 4 words updated simultaneously
- Pure CA evolution stored in state
- Hybrid rotation + multiply mixing applied at output time for diffusion
  - Combines rotation, golden ratio multiplication, and shift-XOR
  - Provides strong avalanche effect and non-linearity
  - Verified through rigorous TestU01 testing (144/144 Crush tests)
- Each iteration generates 256 bits (4 × 64-bit words)

## Building & Testing

```bash
# Build all binaries
make all

# Run Go benchmarks
make bench

# Run comparison tools
make compare-run
./misc/compare-urandom.sh

# Run TestU01 statistical tests
cd testu01
make smallcrush
```

## Use Cases

**Good for:**
- Monte Carlo simulations
- Procedural generation (games, graphics)
- High-throughput random sampling
- Reproducible sequences (deterministic from seed)

**Not for:**
- Cryptographic applications (use `crypto/rand`)
- Security-critical randomness

## Background

Rule 30 is a Class III cellular automaton discovered by Stephen Wolfram in 1983. Despite its simple definition, it exhibits chaotic behavior and generates high-quality randomness. Wolfram Research used it as the basis for Mathematica's random number generator.

## License

MIT License - See LICENSE file for details

## References

- Wolfram, Stephen (1983). "Statistical mechanics of cellular automata"
- Wolfram, Stephen (2002). "A New Kind of Science"
- [Rule 30 on Wikipedia](https://en.wikipedia.org/wiki/Rule_30)
