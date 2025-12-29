package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func mainRule30() {
	var (
		seed      = flag.Uint64("seed", 0, "RNG seed (default: time-based)")
		bytes     = flag.Int("bytes", 1024, "Number of bytes to generate")
		benchmark = flag.Bool("benchmark", false, "Benchmark mode (measure throughput)")
		help      = flag.Bool("help", false, "Show help")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Rule 30 RNG - Random Number Generator using Rule 30 Cellular Automaton

A deterministic RNG based on 1D cellular automata (Rule 30).
Uses a circular 256-bit strip with Rule 30 evolution rules.

Usage:
  rule30-rng [options]

Seed Format:
  64-bit seed initializes the 256-bit circular strip state

Options:
  --seed N        Seed value (default: current time)
  --bytes N       Number of bytes to generate (default: 1024)
  --benchmark     Benchmark throughput instead of generating output
  --help          Show this help

Examples:
  # Generate 1KB of random data
  rule30-rng --bytes 1024 > random.bin

  # Use specific seed
  rule30-rng --seed 12345 --bytes 1048576 > random.bin

  # Benchmark throughput
  rule30-rng --benchmark

  # Test randomness with ent
  rule30-rng --bytes 1048576 | ent

Rule 30:
  A 1D cellular automaton where each cell evolves based on itself
  and its two neighbors according to Rule 30:
    new_bit = left XOR (center OR right)

  Known for generating high-quality pseudo-randomness.
  Used in Mathematica's random number generator.
`)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Use time-based seed if not specified
	if *seed == 0 {
		*seed = uint64(time.Now().UnixNano())
	}

	if *benchmark {
		runBenchmarkRule30(*seed)
	} else {
		generateBytesRule30(*seed, *bytes)
	}
}

// generateBytesRule30 generates and writes random bytes to stdout
func generateBytesRule30(seed uint64, count int) {
	rng := NewRule30(seed)

	fmt.Fprintf(os.Stderr, "Rule 30 RNG initialized\n")
	fmt.Fprintf(os.Stderr, "  Seed: 0x%016X (%d)\n", seed, seed)
	fmt.Fprintf(os.Stderr, "  Strip: 256-bit circular\n")
	fmt.Fprintf(os.Stderr, "  Rule: 30 (left XOR (center OR right))\n")
	fmt.Fprintf(os.Stderr, "  Output: 32 bytes per iteration\n")
	fmt.Fprintf(os.Stderr, "Generating %d bytes...\n", count)

	buf := make([]byte, count)
	n, err := rng.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write to stdout
	written, err := os.Stdout.Write(buf[:n])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Generated %d bytes\n", written)
}

// runBenchmarkRule30 measures RNG throughput
func runBenchmarkRule30(seed uint64) {
	rng := NewRule30(seed)

	sizes := []int{1024, 10240, 102400, 1048576} // 1KB, 10KB, 100KB, 1MB

	fmt.Println("Rule 30 RNG Benchmark")
	fmt.Printf("Seed: 0x%016X\n", seed)
	fmt.Println()
	fmt.Println("Size\t\tTime\t\tThroughput")
	fmt.Println("----\t\t----\t\t----------")

	for _, size := range sizes {
		buf := make([]byte, size)

		start := time.Now()
		n, err := rng.Read(buf)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		throughput := float64(n) / elapsed.Seconds() / 1024 / 1024 // MB/s

		sizeStr := formatSizeRule30(size)
		fmt.Printf("%s\t\t%v\t%.2f MB/s\n", sizeStr, elapsed.Round(time.Millisecond), throughput)
	}
}

// formatSizeRule30 formats byte count for display
func formatSizeRule30(bytes int) string {
	if bytes >= 1048576 {
		return fmt.Sprintf("%d MB", bytes/1048576)
	} else if bytes >= 1024 {
		return fmt.Sprintf("%d KB", bytes/1024)
	}
	return fmt.Sprintf("%d B", bytes)
}
