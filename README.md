# ring30mix

High-performance pseudo-random number generator based on Rule 30 cellular automaton with golden ratio mixing.

- **~2× faster** than Go's math/rand/v2 PCG
- **Perfect BigCrush score** (160/160) with superior p-value distribution

See [ALGORITHM.md](ALGORITHM.md) for details.

## Quick Start

### Installation

```bash
# As a command-line tool
go install github.com/vrypan/ring30mix@latest

# As a library
go get github.com/vrypan/ring30mix
```

### Command Line Usage

```bash
# Generate random data
./ring30mix --bytes=1048576 > random.bin

# Generate specific size with dd
./ring30mix --bytes=1073741824 | dd of=test.data bs=1m

# Unlimited streaming (use with head, pv, or Ctrl+C)
./ring30mix --bytes=0 | head -c 1073741824 > test.data

# Reproducible output with seed
./ring30mix --seed=12345 --bytes=1024 > random.bin
```

### Library Usage

Drop-in replacement for math/rand:

```go
import "github.com/vrypan/ring30mix/rand"

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

Benchmarks on Apple M4, verified 2026-01-03 (run `make bench` to reproduce):

### Absolute Performance

|Algorithm               |  Read() 32KB |   Read() 1KB |     Uint64()|
|------------------------|-------------:|-------------:|------------:|
|math/rand/v2 PCG        |  13249.00 ns |    413.30 ns |      3.23 ns|
|math/rand/v2 ChaCha8    |  11478.00 ns |    360.70 ns |      2.77 ns|
|**ring30mix**           |   **6721.00 ns** |    **214.90 ns** |      **1.62 ns**|
|math/rand               |  21409.00 ns |    683.50 ns |      1.80 ns|
|crypto/rand             |   7448.00 ns |    365.10 ns |     54.58 ns|

### Speed vs math/rand/v2 PCG (baseline = 1.00×)

|Algorithm               | Read() 32KB |  Read() 1KB |    Uint64()|
|------------------------|------------:|------------:|-----------:|
|math/rand/v2 PCG        |       1.00x |       1.00x |       1.00x|
|math/rand/v2 ChaCha8    |       1.15x |       1.15x |       1.16x|
|**ring30mix**           |    **1.97x** |   **1.92x**|   **1.99x**|
|math/rand               |       0.62x |       0.60x |       1.79x|
|crypto/rand             |       1.78x |       1.13x |       0.06x|

## Randomness Quality

**Perfect BigCrush score** - verified 2026-01-03:

| Test Suite | Tests | Passed | Status |
|-----------:|------:|-------:|-------:|
| SmallCrush | 15 | 15/15 ✅ | 100% |
| Crush | 144 | 144/144 ✅ | 100% |
| **BigCrush** | **160** | **160/160 ✅** | **100%** |

**Superior to math/rand/v2 PCG:**
- ring30mix: **0 borderline p-values** (p < 0.01)
- math/rand/v2 PCG: **3 borderline p-values**

Run tests yourself:
```bash
cd testu01
make smallcrush   # 15 tests, ~2 min
make crush        # 144 tests, ~10 min
make bigcrush     # 160 tests, ~1 hour
```

## How It Works

**Core algorithm:**
- Rule 30 cellular automaton: `new_bit = left XOR (center OR right)`
- 256-bit ring (4 × 64-bit words)
- Golden ratio mixing on output (rotation + φ multiply + shift-XOR)

**Key optimizations:**
- Bit-parallel processing (64 bits per operation)
- Fully unrolled loops
- Pre-computed border bits
- Amortized state evolution (1 step per 4 outputs)

See [ALGORITHM.md](ALGORITHM.md) for complete technical details, optimizations, and analysis.

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
