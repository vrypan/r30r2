package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math"
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
	duration   time.Duration
	throughput float64 // MB/s
	entropy    float64 // Shannon entropy (bits per byte)
	chiSquare  float64 // Chi-square statistic
}

// calculateEntropy computes Shannon entropy in bits per byte
func calculateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	// Count frequency of each byte value
	freq := make([]int, 256)
	for _, b := range data {
		freq[b]++
	}

	// Calculate Shannon entropy
	entropy := 0.0
	dataLen := float64(len(data))

	for _, count := range freq {
		if count == 0 {
			continue
		}
		p := float64(count) / dataLen
		entropy -= p * math.Log2(p)
	}

	return entropy
}

// calculateChiSquare computes chi-square statistic for uniform distribution test
func calculateChiSquare(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	// Count frequency of each byte value
	freq := make([]int, 256)
	for _, b := range data {
		freq[b]++
	}

	// Expected frequency for uniform distribution
	expected := float64(len(data)) / 256.0

	// Calculate chi-square: χ² = Σ((observed - expected)² / expected)
	chiSquare := 0.0
	for _, count := range freq {
		observed := float64(count)
		diff := observed - expected
		chiSquare += (diff * diff) / expected
	}

	return chiSquare
}

// runBenchmark tests an io.Reader and returns results
func runBenchmark(name string, r io.Reader, size int, iterations int) BenchResult {
	buf := make([]byte, size)
	allData := make([]byte, 0, size*iterations)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
			os.Exit(1)
		}
		// Collect data for entropy calculation
		allData = append(allData, buf...)
	}
	duration := time.Since(start)

	totalBytes := float64(size * iterations)
	throughput := totalBytes / duration.Seconds() / 1024 / 1024 // MB/s

	// Calculate entropy and chi-square on collected data
	entropy := calculateEntropy(allData)
	chiSquare := calculateChiSquare(allData)

	return BenchResult{
		name:       name,
		size:       size,
		duration:   duration,
		throughput: throughput,
		entropy:    entropy,
		chiSquare:  chiSquare,
	}
}

// runUint64Benchmark tests Uint64() generation and returns results
func runUint64Benchmark(name string, iterations int, genFunc func() uint64) BenchResult {
	// Collect data for entropy calculation
	allData := make([]byte, 0, iterations*8)
	buf := make([]byte, 8)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		val := genFunc()
		binary.LittleEndian.PutUint64(buf, val)
		allData = append(allData, buf...)
	}
	duration := time.Since(start)

	totalBytes := float64(iterations * 8)
	throughput := totalBytes / duration.Seconds() / 1024 / 1024 // MB/s

	// Calculate entropy and chi-square on collected data
	entropy := calculateEntropy(allData)
	chiSquare := calculateChiSquare(allData)

	return BenchResult{
		name:       name,
		size:       iterations * 8,
		duration:   duration,
		throughput: throughput,
		entropy:    entropy,
		chiSquare:  chiSquare,
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
	// Parse command line flags
	mode := flag.String("mode", "both", "Benchmark mode: read, uint64, or both")
	flag.Parse()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Random Number Generator Performance Comparison")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Test configurations
	sizes := []int{
		1024,          // 1 KB
		10 * 1024,     // 10 KB
		100 * 1024,    // 100 KB
		1024 * 1024,   // 1 MB
	}

	// Adjust iterations based on size for reasonable runtime
	iterations := map[int]int{
		1024:          10000,  // 1 KB: 10K iterations
		10 * 1024:     5000,   // 10 KB: 5K iterations
		100 * 1024:    1000,   // 100 KB: 1K iterations
		1024 * 1024:   100,    // 1 MB: 100 iterations
	}

	// Validate mode
	if *mode != "read" && *mode != "uint64" && *mode != "both" {
		fmt.Fprintf(os.Stderr, "Error: Invalid mode '%s'. Use 'read', 'uint64', or 'both'\n", *mode)
		os.Exit(1)
	}

	// Store results by RNG type and size
	results := make(map[string]map[int]BenchResult)
	results["Rule30RNG"] = make(map[int]BenchResult)
	results["math/rand"] = make(map[int]BenchResult)
	results["crypto/rand"] = make(map[int]BenchResult)

	// Run Read() benchmarks
	if *mode == "read" || *mode == "both" {
		if *mode == "both" {
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("  Benchmark 1: Read() - Bulk Byte Generation")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()
			fmt.Println("This benchmark measures bulk byte stream generation performance.")
			fmt.Println("Rule30 is optimized for this (32 bytes/iteration).")
			fmt.Println("math/rand.Read() is a convenience wrapper over Int63().")
			fmt.Println()
		}

		for _, size := range sizes {
			iters := iterations[size]
			sizeStr := formatSize(size)

			fmt.Printf("Testing with %s buffers (%d iterations)...\n", sizeStr, iters)

			// Rule30RNG
			rule30rng := rule30.New(12345)
			result := runBenchmark("Rule30RNG", rule30rng, size, iters)
			results["Rule30RNG"][size] = result
			fmt.Printf("  ✓ Rule30RNG:   %7.2f MB/s  (entropy: %.4f, χ²: %.1f)\n", result.throughput, result.entropy, result.chiSquare)

			// math/rand
			mathRng := newMathRandReader(12345)
			result = runBenchmark("math/rand", mathRng, size, iters)
			results["math/rand"][size] = result
			fmt.Printf("  ✓ math/rand:   %7.2f MB/s  (entropy: %.4f, χ²: %.1f)\n", result.throughput, result.entropy, result.chiSquare)

			// crypto/rand
			result = runBenchmark("crypto/rand", cryptorand.Reader, size, iters)
			results["crypto/rand"][size] = result
			fmt.Printf("  ✓ crypto/rand: %7.2f MB/s  (entropy: %.4f, χ²: %.1f)\n", result.throughput, result.entropy, result.chiSquare)

			fmt.Println()
		}
	}

	// Run Uint64() benchmarks
	if *mode == "uint64" || *mode == "both" {
		if *mode == "both" {
			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("  Benchmark 2: Uint64() - Fair RNG Algorithm Comparison")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()
			fmt.Println("This benchmark compares the underlying RNG algorithms fairly.")
			fmt.Println("All RNGs generate 64-bit values using their primitive operations.")
			fmt.Println()
		}

		// Uint64 iterations (generating 8 bytes each)
		uint64Iterations := []struct {
			count int
			desc  string
		}{
			{10000000, "10M values (80 MB)"},
			{50000000, "50M values (400 MB)"},
		}

		// crypto/rand Uint64 wrapper
		cryptoUint64 := func() uint64 {
			var buf [8]byte
			cryptorand.Read(buf[:])
			return binary.LittleEndian.Uint64(buf[:])
		}

		for _, test := range uint64Iterations {
			fmt.Printf("Testing %s...\n", test.desc)

			// Rule30RNG
			rule30rng := rule30.New(12345)
			result := runUint64Benchmark("Rule30RNG", test.count, rule30rng.Uint64)
			// Store in results using size as key (for table generation)
			results["Rule30RNG"][test.count*8] = result
			fmt.Printf("  ✓ Rule30RNG:   %7.2f MB/s  (entropy: %.4f, χ²: %.1f)\n", result.throughput, result.entropy, result.chiSquare)

			// math/rand
			mathRng := mathrand.New(mathrand.NewSource(12345))
			result = runUint64Benchmark("math/rand", test.count, mathRng.Uint64)
			results["math/rand"][test.count*8] = result
			fmt.Printf("  ✓ math/rand:   %7.2f MB/s  (entropy: %.4f, χ²: %.1f)\n", result.throughput, result.entropy, result.chiSquare)

			// crypto/rand
			result = runUint64Benchmark("crypto/rand", test.count, cryptoUint64)
			results["crypto/rand"][test.count*8] = result
			fmt.Printf("  ✓ crypto/rand: %7.2f MB/s  (entropy: %.4f, χ²: %.1f)\n", result.throughput, result.entropy, result.chiSquare)

			fmt.Println()
		}

		// Update sizes for table generation in Uint64 mode
		if *mode == "uint64" {
			sizes = []int{80000000, 400000000} // 10M and 50M * 8 bytes
		}
	}

	// Generate summary table
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Summary Table")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Table header
	fmt.Printf("%-15s", "RNG")
	for _, size := range sizes {
		fmt.Printf("│ %8s ", formatSize(size))
	}
	fmt.Printf("│ Avg Speed\n")

	fmt.Println("───────────────┼──────────┼──────────┼──────────┼──────────┼──────────")

	// Table rows
	rngNames := []string{"Rule30RNG", "math/rand", "crypto/rand"}
	for _, rngName := range rngNames {
		fmt.Printf("%-15s", rngName)

		var totalThroughput float64
		for _, size := range sizes {
			result := results[rngName][size]
			fmt.Printf("│ %6.1f MB ", result.throughput)
			totalThroughput += result.throughput
		}
		avgThroughput := totalThroughput / float64(len(sizes))
		fmt.Printf("│ %6.1f MB\n", avgThroughput)
	}

	fmt.Println()

	// Entropy comparison table
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Shannon Entropy (bits per byte)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Printf("%-15s", "RNG")
	for _, size := range sizes {
		fmt.Printf("│ %8s ", formatSize(size))
	}
	fmt.Printf("│ Average\n")

	fmt.Println("───────────────┼──────────┼──────────┼──────────┼──────────┼──────────")

	for _, rngName := range rngNames {
		fmt.Printf("%-15s", rngName)

		var totalEntropy float64
		for _, size := range sizes {
			result := results[rngName][size]
			fmt.Printf("│  %7.5f ", result.entropy)
			totalEntropy += result.entropy
		}
		avgEntropy := totalEntropy / float64(len(sizes))
		fmt.Printf("│  %7.5f\n", avgEntropy)
	}

	fmt.Println()
	fmt.Println("Note: Maximum entropy = 8.000000 bits/byte (perfect randomness)")
	fmt.Println()

	// Chi-square comparison table
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Chi-Square Distribution Test")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Printf("%-15s", "RNG")
	for _, size := range sizes {
		fmt.Printf("│ %8s ", formatSize(size))
	}
	fmt.Printf("│ Average\n")

	fmt.Println("───────────────┼──────────┼──────────┼──────────┼──────────┼──────────")

	for _, rngName := range rngNames {
		fmt.Printf("%-15s", rngName)

		var totalChiSquare float64
		for _, size := range sizes {
			result := results[rngName][size]
			fmt.Printf("│  %7.1f ", result.chiSquare)
			totalChiSquare += result.chiSquare
		}
		avgChiSquare := totalChiSquare / float64(len(sizes))
		fmt.Printf("│  %7.1f\n", avgChiSquare)
	}

	fmt.Println()
	fmt.Println("Note: Expected value ≈ 255 for uniform distribution (df=255)")
	fmt.Println("      Acceptable range: ~200-310 (within 95% confidence)")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Additional info
	fmt.Println("Notes:")
	fmt.Println("  • Rule30RNG:  1D CA (Rule 30), 256-bit state, deterministic")
	fmt.Println("  • math/rand:  Fast PRNG, deterministic")
	fmt.Println("  • crypto/rand: Hardware-accelerated, cryptographically secure")
	fmt.Println()
	fmt.Println("Shannon Entropy Interpretation:")
	fmt.Println("  7.990-8.000: Excellent randomness")
	fmt.Println("  7.900-7.990: Good randomness")
	fmt.Println("  7.500-7.900: Fair randomness")
	fmt.Println("  < 7.500:     Poor randomness")
	fmt.Println()
	fmt.Println("Chi-Square Interpretation (df=255):")
	fmt.Println("  200-310:     Excellent (within 95% confidence interval)")
	fmt.Println("  180-330:     Good (within 99% confidence interval)")
	fmt.Println("  < 180/> 330: Poor (distribution not uniform)")
	fmt.Println()
}
