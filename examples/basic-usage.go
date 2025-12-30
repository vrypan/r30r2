package main

import (
	"fmt"

	"github.com/vrypan/rule30rnd/rand"
)

func main() {
	// Create a new Rule 30 RNG with a seed
	rng := rand.New(12345)

	fmt.Println("Rule 30 RNG - Basic Usage Example")
	fmt.Println("==================================")
	fmt.Println()

	// Generate basic integers
	fmt.Println("Basic Integer Generation:")
	fmt.Printf("  Random uint32:  %d\n", rng.Uint32())
	fmt.Printf("  Random uint64:  %d\n", rng.Uint64())
	fmt.Printf("  Random int:     %d\n", rng.Int())
	fmt.Println()

	// Generate bounded integers
	fmt.Println("Bounded Integer Generation:")
	fmt.Printf("  Random number [0, 100):  %d\n", rng.Intn(100))
	fmt.Printf("  Random number [0, 10):   %d\n", rng.Intn(10))
	fmt.Println()

	// Generate floats
	fmt.Println("Float Generation:")
	fmt.Printf("  Random float64 [0.0, 1.0):  %.6f\n", rng.Float64())
	fmt.Printf("  Random float32 [0.0, 1.0):  %.6f\n", rng.Float32())
	fmt.Println()

	// Statistical distributions
	fmt.Println("Statistical Distributions:")
	fmt.Printf("  Normal distribution (μ=0, σ=1):  %.6f\n", rng.NormFloat64())
	fmt.Printf("  Exponential distribution (λ=1): %.6f\n", rng.ExpFloat64())
	fmt.Println()

	// Generate random bytes (implements io.Reader)
	fmt.Println("Byte Generation (io.Reader):")
	buf := make([]byte, 16)
	n, _ := rng.Read(buf)
	fmt.Printf("  Generated %d random bytes: %x\n", n, buf)
}
