package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/vrypan/rule30rnd/rand"
)

// BenchResult holds benchmark results
type BenchResult struct {
	name      string
	calls     int
	nsPerCall float64 // nanoseconds per call
}

// runUint64Benchmark tests Uint64() generation and returns results
func runUint64Benchmark(name string, iterations int, genFunc func() uint64) BenchResult {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = genFunc()
	}
	duration := time.Since(start)

	nsPerCall := float64(duration.Nanoseconds()) / float64(iterations)

	return BenchResult{
		name:      name,
		calls:     iterations,
		nsPerCall: nsPerCall,
	}
}

// formatCalls formats call counts
func formatCalls(calls int) string {
	if calls >= 1000000 {
		return fmt.Sprintf("%dM calls", calls/1000000)
	} else if calls >= 1000 {
		return fmt.Sprintf("%dk calls", calls/1000)
	}
	return fmt.Sprintf("%d calls", calls)
}

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Uint64() Benchmark - Latency per Call")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Test configurations - call counts
	callCounts := []int{
		10000000, // 10M calls
	}

	// Store results by RNG type and call count
	results := make(map[string]map[int]BenchResult)
	results["Rule30RNG"] = make(map[int]BenchResult)
	results["math/rand"] = make(map[int]BenchResult)
	results["crypto/rand"] = make(map[int]BenchResult)

	// crypto/rand Uint64 wrapper
	cryptoUint64 := func() uint64 {
		var buf [8]byte
		cryptorand.Read(buf[:])
		return binary.LittleEndian.Uint64(buf[:])
	}

	// Run benchmarks
	for _, calls := range callCounts {
		fmt.Printf("Testing %s...\n", formatCalls(calls))

		// Rule30RNG
		rule30rng := rand.New(12345)
		result := runUint64Benchmark("Rule30RNG", calls, rule30rng.Uint64)
		results["Rule30RNG"][calls] = result
		fmt.Printf("  ✓ Rule30RNG:   %6.1f ns/call\n", result.nsPerCall)

		// math/rand
		mathRng := mathrand.New(mathrand.NewSource(12345))
		result = runUint64Benchmark("math/rand", calls, mathRng.Uint64)
		results["math/rand"][calls] = result
		fmt.Printf("  ✓ math/rand:   %6.1f ns/call\n", result.nsPerCall)

		// crypto/rand
		result = runUint64Benchmark("crypto/rand", calls, cryptoUint64)
		results["crypto/rand"][calls] = result
		fmt.Printf("  ✓ crypto/rand: %6.1f ns/call\n", result.nsPerCall)

		fmt.Println()
	}

	// Generate summary table
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Summary Table (ns/call)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Table header
	fmt.Printf("%-15s │ %-12s │ %-10s\n", "RNG", "10M ns/call", "Relative")
	fmt.Println("────────────────┼──────────────┼────────────")

	// Table rows with Rule30RNG as baseline
	rngNames := []string{"Rule30RNG", "math/rand", "crypto/rand"}
	baseline := results["Rule30RNG"][callCounts[0]].nsPerCall
	for _, rngName := range rngNames {
		result := results[rngName][callCounts[0]]
		relative := result.nsPerCall / baseline
		fmt.Printf("%-15s │ %9.1f ns │ %8.2f×\n", rngName, result.nsPerCall, relative)
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Additional info
	fmt.Println("Notes:")
	fmt.Println("  • Rule30RNG:  1D CA (Rule 30), 256-bit state, deterministic")
	fmt.Println("  • math/rand:  Fast PRNG (PCG algorithm), deterministic")
	fmt.Println("  • crypto/rand: Hardware-accelerated CSPRNG")
	fmt.Println("  • Lower ns/call is better (faster)")
	fmt.Println()
}
