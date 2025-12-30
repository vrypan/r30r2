package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/vrypan/rule30rnd/rand"
)

func mainRule30() {
	var (
		seed      = flag.Uint64("seed", 0, "RNG seed (default: time-based)")
		bytes     = flag.Int("bytes", 1024, "Number of bytes to generate")
		benchmark = flag.Bool("benchmark", false, "Benchmark mode (measure throughput)")
		help      = flag.Bool("help", false, "Show help")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Rule 30 RND - Random Number Generator using Rule 30 Cellular Automaton

A deterministic RNG based on 1D cellular automata (Rule 30).
Uses a circular 256-bit strip with Rule 30 evolution rules.

Usage:
  rule30 [options]

Seed Format:
  64-bit seed initializes the 256-bit circular strip state

Options:
  --seed N        Seed value (default: current time)
  --bytes N       Number of bytes to generate (default: 1024, 0 = unlimited)
  --benchmark     Benchmark throughput instead of generating output
  --help          Show this help

Examples:
  # Generate 1KB of random data
  rule30 --bytes 1024 > random.bin

  # Use specific seed
  rule30 --seed 12345 --bytes 1048576 > random.bin

  # Generate specific size with dd (piping, not using dd count)
  rule30 --bytes 1073741824 | dd of=test.data bs=1m

  # Unlimited streaming (use with head, pv, or Ctrl+C)
  rule30 --bytes 0 | head -c 1073741824 > test.data

  # Benchmark throughput
  rule30 --benchmark

  # Test randomness with ent
  rule30 --bytes 1048576 | ent

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
	rng := rand.New(seed)

	fmt.Fprintf(os.Stderr, "Rule 30 RNG initialized\n")
	fmt.Fprintf(os.Stderr, "  Seed: 0x%016X (%d)\n", seed, seed)
	fmt.Fprintf(os.Stderr, "  Strip: 256-bit circular\n")
	fmt.Fprintf(os.Stderr, "  Rule: 30 (left XOR (center OR right))\n")
	fmt.Fprintf(os.Stderr, "  Output: 32 bytes per iteration\n")

	if count == 0 {
		fmt.Fprintf(os.Stderr, "Generating unlimited bytes (streaming mode)...\n")
		// Unlimited mode: stream chunks until pipe breaks
		buf := make([]byte, 1024*1024) // 1MB chunks
		for {
			n, err := rng.Read(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Write all bytes, handling partial writes
			written := 0
			for written < n {
				w, err := os.Stdout.Write(buf[written:n])
				if err != nil {
					// Pipe closed (e.g., dd finished) - exit gracefully
					os.Exit(0)
				}
				written += w
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Generating %d bytes...\n", count)

		// Fixed size: stream in chunks to avoid huge allocations
		const chunkSize = 1024 * 1024 // 1MB chunks
		buf := make([]byte, chunkSize)
		remaining := count
		totalWritten := 0

		for remaining > 0 {
			toRead := chunkSize
			if remaining < chunkSize {
				toRead = remaining
			}

			n, err := rng.Read(buf[:toRead])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Write all bytes from this read, handling partial writes
			written := 0
			for written < n {
				w, err := os.Stdout.Write(buf[written:n])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
					os.Exit(1)
				}
				written += w
			}

			totalWritten += n
			remaining -= n
		}

		fmt.Fprintf(os.Stderr, "Generated %d bytes\n", totalWritten)
	}
}

// runBenchmarkRule30 measures RNG throughput
func runBenchmarkRule30(seed uint64) {
	rng := rand.New(seed)

	sizes := []int{1024, 10240, 102400, 1048576} // 1KB, 10KB, 100KB, 1MB

	fmt.Println("Rule 30 RNG Benchmark")
	fmt.Printf("Seed: 0x%016X\n", seed)
	fmt.Println()
	fmt.Printf("%6s    %8s    %12s\n", "Size", "Time", "Throughput")
	fmt.Printf("%6s    %8s    %12s\n", "----", "----", "----------")

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
		timeStr := formatDurationRule30(elapsed)
		throughputStr := formatThroughputRule30(throughput)
		fmt.Printf("%6s    %8s    %12s\n", sizeStr, timeStr, throughputStr)
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

// formatDurationRule30 formats duration in appropriate units with 8-char padding
func formatDurationRule30(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%6d ns", d.Nanoseconds())
	} else if d < time.Millisecond {
		return fmt.Sprintf("%6.2f µs", float64(d.Nanoseconds())/1000.0)
	} else if d < time.Second {
		return fmt.Sprintf("%6.2f ms", float64(d.Nanoseconds())/1000000.0)
	}
	return fmt.Sprintf("%6.2f s", d.Seconds())
}

// formatThroughputRule30 formats throughput in MB/s with padding
func formatThroughputRule30(mbps float64) string {
	return fmt.Sprintf("%9.2f MB/s", mbps)
}
