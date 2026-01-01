# Rule 30 RND

High-performance pseudo-random number generator based on Rule 30 cellular automaton. **Up to 3× faster** than Go's standard library for bulk operations, with **exceptional statistical quality** - passes all 160 BigCrush tests.

## Quick Start

### Installation

```bash
# As a command-line tool
go install github.com/vrypan/rule30rnd@latest

# As a library
go get github.com/vrypan/rule30rnd
```

### Command Line Usage

```bash
# Generate random data
./rule30 --bytes=1048576 > random.bin

# Generate specific size with dd
./rule30 --bytes=1073741824 | dd of=test.data bs=1m

# Unlimited streaming (use with head, pv, or Ctrl+C)
./rule30 --bytes=0 | head -c 1073741824 > test.data

# Benchmark throughput
./rule30 --benchmark

# Reproducible output with seed
./rule30 --seed=12345 --bytes=1024 > random.bin
```

### Library Usage

Drop-in replacement for math/rand:

```go
import "github.com/vrypan/rule30rnd/rand"

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

Benchmarks on Apple M4 (run `make bench` to reproduce):

### Throughput Comparison

|Algorithm       |     Read32KB |      Read1KB |       Uint64|
|----------------|--------------|--------------|-------------|
|MathRand        |  21238.00 ns |    667.90 ns |      1.85 ns|
|MathRandV2      |  13042.00 ns |    427.00 ns |      3.29 ns|
|**Rule30**      |   **7592.00 ns** |    **242.30 ns** |      **1.67 ns**|
|CryptoRand      |   7019.00 ns |    368.70 ns |     56.10 ns|

**Relative to math/rand:**

|Algorithm       |    Read32KB |     Read1KB |      Uint64|
|----------------|-------------|-------------|------------|
|MathRand        |       1.00x |       1.00x |       1.00x|
|MathRandV2      |       1.63x |       1.56x |       0.56x|
|**Rule30**      |       **2.80x** |       **2.76x** |       **1.11x**|
|CryptoRand      |       3.03x |       1.81x |       0.03x|

### vs /dev/urandom

Rule30 vs kernel CSPRNG (run `./misc/compare-urandom.sh`):

|Size    |       Rule30 |  /dev/urandom |   Speedup|
|--------|--------------|---------------|----------|
|10 MB   |     388 MB/s |      433 MB/s |    0.89x|
|100 MB  |    2220 MB/s |      508 MB/s |    4.36x|
|1 GB    |    2490 MB/s |      498 MB/s |    4.99x|

R30R2 excels at bulk generation: **~5× faster** than the kernel CSPRNG for files 100 MB and larger.

## Randomness Quality

**Perfect scores across all TestU01 batteries** - passes every test in the most comprehensive PRNG test suite:

| Test Battery | Tests | Rule30 | Pass Rate | Time |
|--------------|-------|--------|-----------|------|
| **SmallCrush** | 15 | **15/15** ✓ | **100%** | ~2 min |
| **Crush** | 144 | **144/144** ✓ | **100%** | ~10 min |
| **BigCrush** | 160 | **160/160** ✓ | **100%** | ~1 hour |

Rule30 achieves perfect results on BigCrush, the most rigorous statistical randomness test suite. This demonstrates exceptional quality that matches or exceeds established PRNGs, testing across serial correlation, birthday spacing, collision, permutation, matrix rank, spectral, string, compression, and random walk tests.

Run tests yourself:

```bash
cd testu01
make smallcrush   # Quick test (15 tests, 2 min)
make crush        # Medium test (144 tests, 10 min)
make bigcrush     # Comprehensive test (160 tests, 1 hour)
```

## How It Works

**Algorithm**: Radius-2 cellular automaton based on Rule 30 on a circular 256-bit strip.

```
Evolution rule: new_bit = (left2 XOR left1) XOR ((center OR right1) OR right2)
```

**Implementation**:
- 256-bit state as 4 × 64-bit words
- Radius-2 neighborhood: each bit depends on 5 neighboring cells
- Processes 64 bits in parallel per word (bitwise operations)
- Fully unrolled loop - all 4 words updated simultaneously
- Pure CA evolution stored in state
- Multi-rotation XOR mixing (angles 13, 17, 23) applied at output time for exceptional diffusion
- Each iteration generates 256 bits (4 × 64-bit words)

**Why it's fast**:
- Bit-level parallelism (64 bits per operation vs 1 bit serially)
- Single-cycle operations (XOR, OR, shifts)
- Zero loop overhead (fully unrolled)
- Minimal memory allocations

**Spatial vs Temporal extraction**: Unlike Wolfram's traditional approach (1 bit per generation from a single column), we extract all 256 bits from the spatial pattern each iteration - **256× more efficient** while maintaining excellent randomness.

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
