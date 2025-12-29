package main

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"
	mathrand "math/rand"
	"os"
	"time"

	"github.com/vrypan/rule30rnd/rule30"
)

// mathRandReader wraps math/rand to implement io.Reader
type mathRandReader struct {
	rng *mathrand.Rand
}

func (m *mathRandReader) Read(p []byte) (n int, err error) {
	return m.rng.Read(p)
}

func newMathRandReader(seed int64) io.Reader {
	return &mathRandReader{
		rng: mathrand.New(mathrand.NewSource(seed)),
	}
}

// BenchResult holds benchmark results
type BenchResult struct {
	name       string
	size       int
	throughput float64 // MB/s
}

// runBenchmark tests an io.Reader and returns results
func runBenchmark(name string, r io.Reader, size int, iterations int) BenchResult {
	buf := make([]byte, size)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
			os.Exit(1)
		}
	}
	duration := time.Since(start)

	totalBytes := float64(size * iterations)
	throughput := totalBytes / duration.Seconds() / 1024 / 1024 // MB/s

	return BenchResult{
		name:       name,
		size:       size,
		throughput: throughput,
	}
}

// formatSize formats bytes as KB or MB
func formatSize(bytes int) string {
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%d MB", bytes/(1024*1024))
	}
	return fmt.Sprintf("%d KB", bytes/1024)
}

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Read() Benchmark - Bulk Byte Stream Generation")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Test configurations
	sizes := []int{
		1024,              // 1 KB
		10 * 1024,         // 10 KB
		100 * 1024,        // 100 KB
		1024 * 1024,       // 1 MB
		100 * 1024 * 1024, // 100 MB
	}

	// Adjust iterations based on size for reasonable runtime
	iterations := map[int]int{
		1024:              100000, // 1 KB: 100K iterations
		10 * 1024:         10000,  // 10 KB: 10K iterations
		100 * 1024:        1000,   // 100 KB: 1K iterations
		1024 * 1024:       100,    // 1 MB: 100 iterations
		100 * 1024 * 1024: 10,     // 100 MB: 10 iterations
	}

	// Store results by RNG type and size
	results := make(map[string]map[int]BenchResult)
	results["Rule30RNG"] = make(map[int]BenchResult)
	results["math/rand"] = make(map[int]BenchResult)
	results["crypto/rand"] = make(map[int]BenchResult)

	// Run benchmarks
	for _, size := range sizes {
		iters := iterations[size]
		sizeStr := formatSize(size)

		fmt.Printf("Testing with %s buffers (%d iterations)...\n", sizeStr, iters)

		// Rule30RNG
		rule30rng := rule30.New(12345)
		result := runBenchmark("Rule30RNG", rule30rng, size, iters)
		results["Rule30RNG"][size] = result
		fmt.Printf("  ✓ Rule30RNG:   %7.2f MB/s\n", result.throughput)

		// math/rand
		mathRng := newMathRandReader(12345)
		result = runBenchmark("math/rand", mathRng, size, iters)
		results["math/rand"][size] = result
		fmt.Printf("  ✓ math/rand:   %7.2f MB/s\n", result.throughput)

		// crypto/rand
		result = runBenchmark("crypto/rand", cryptorand.Reader, size, iters)
		results["crypto/rand"][size] = result
		fmt.Printf("  ✓ crypto/rand: %7.2f MB/s\n", result.throughput)

		fmt.Println()
	}

	// Generate summary table
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Summary Table (MB/s)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Table header
	fmt.Printf("%-15s", "RNG")
	for _, size := range sizes {
		fmt.Printf(" │ %-12s", formatSize(size))
	}
	fmt.Println()

	// Separator
	fmt.Print("───────────────")
	for range sizes {
		fmt.Print("─┼──────────────")
	}
	fmt.Println()

	// Table rows
	rngNames := []string{"Rule30RNG", "math/rand", "crypto/rand"}
	for _, rngName := range rngNames {
		fmt.Printf("%-15s", rngName)

		for _, size := range sizes {
			result := results[rngName][size]
			fmt.Printf(" │ %8.0f MB/s", result.throughput)
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Additional info
	fmt.Println("Notes:")
	fmt.Println("  • Rule30RNG:  1D CA (Rule 30), 256-bit state, deterministic")
	fmt.Println("  • math/rand:  Fast PRNG (PCG algorithm), deterministic")
	fmt.Println("  • crypto/rand: Hardware-accelerated CSPRNG")
	fmt.Println()
}
