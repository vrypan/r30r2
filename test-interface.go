package main

import (
	"fmt"
	"math"

	"github.com/vrypan/rule30rnd/rule30"
)

func main() {
	fmt.Println("Testing Rule30RNG math/rand interface compatibility")
	fmt.Println("====================================================")
	fmt.Println()

	rng := rule30.New(12345)

	// Test Uint32/Uint64
	fmt.Println("Basic integer generation:")
	fmt.Printf("  Uint32():  %d\n", rng.Uint32())
	fmt.Printf("  Uint64():  %d\n", rng.Uint64())
	fmt.Printf("  Int():     %d\n", rng.Int())
	fmt.Printf("  Int31():   %d\n", rng.Int31())
	fmt.Printf("  Int63():   %d\n", rng.Int63())
	fmt.Println()

	// Test bounded integers
	fmt.Println("Bounded integer generation:")
	fmt.Printf("  Intn(100):   %d\n", rng.Intn(100))
	fmt.Printf("  Int31n(100): %d\n", rng.Int31n(100))
	fmt.Printf("  Int63n(100): %d\n", rng.Int63n(100))
	fmt.Println()

	// Test floats
	fmt.Println("Float generation:")
	fmt.Printf("  Float32(): %.6f\n", rng.Float32())
	fmt.Printf("  Float64(): %.6f\n", rng.Float64())
	fmt.Println()

	// Test distributions
	fmt.Println("Distribution samples:")
	fmt.Printf("  NormFloat64(): %.6f\n", rng.NormFloat64())
	fmt.Printf("  ExpFloat64():  %.6f\n", rng.ExpFloat64())
	fmt.Println()

	// Test uniformity of Intn
	fmt.Println("Testing uniformity of Intn(10) over 100,000 samples:")
	counts := make([]int, 10)
	for i := 0; i < 100000; i++ {
		n := rng.Intn(10)
		counts[n]++
	}
	for i, count := range counts {
		percent := float64(count) / 1000.0
		fmt.Printf("  %d: %d (%.2f%%)\n", i, count, percent)
	}
	fmt.Println()

	// Test Float64 distribution
	fmt.Println("Testing Float64() distribution (100,000 samples):")
	buckets := make([]int, 10)
	for i := 0; i < 100000; i++ {
		f := rng.Float64()
		bucket := int(f * 10)
		if bucket >= 10 {
			bucket = 9
		}
		buckets[bucket]++
	}
	for i, count := range buckets {
		percent := float64(count) / 1000.0
		fmt.Printf("  [%.1f-%.1f): %d (%.2f%%)\n", float64(i)/10, float64(i+1)/10, count, percent)
	}
	fmt.Println()

	// Test NormFloat64 mean and stddev
	fmt.Println("Testing NormFloat64() (10,000 samples):")
	var sum, sumSq float64
	n := 10000
	for i := 0; i < n; i++ {
		x := rng.NormFloat64()
		sum += x
		sumSq += x * x
	}
	mean := sum / float64(n)
	variance := (sumSq / float64(n)) - (mean * mean)
	stddev := math.Sqrt(variance)
	fmt.Printf("  Mean:   %.6f (expected: 0.0)\n", mean)
	fmt.Printf("  Stddev: %.6f (expected: 1.0)\n", stddev)
	fmt.Println()

	// Test ExpFloat64 mean
	fmt.Println("Testing ExpFloat64() (10,000 samples):")
	sum = 0
	for i := 0; i < n; i++ {
		sum += rng.ExpFloat64()
	}
	mean = sum / float64(n)
	fmt.Printf("  Mean: %.6f (expected: 1.0)\n", mean)
	fmt.Println()

	fmt.Println("✓ All tests completed successfully")
}
