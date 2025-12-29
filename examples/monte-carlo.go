package main

import (
	"fmt"
	"math"

	"github.com/vrypan/rule30rnd/rule30"
)

// Estimate π using Monte Carlo method
func estimatePi(rng *rule30.RNG, samples int) float64 {
	inside := 0

	for i := 0; i < samples; i++ {
		x := rng.Float64()
		y := rng.Float64()

		// Check if point is inside unit circle
		if x*x+y*y <= 1.0 {
			inside++
		}
	}

	// π ≈ 4 × (points inside circle / total points)
	return 4.0 * float64(inside) / float64(samples)
}

func main() {
	rng := rule30.New(42)

	fmt.Println("Monte Carlo π Estimation using Rule 30 RNG")
	fmt.Println("==========================================")
	fmt.Println()

	samples := []int{1000, 10000, 100000, 1000000}

	for _, n := range samples {
		pi := estimatePi(rng, n)
		error := math.Abs(pi - math.Pi)
		errorPct := (error / math.Pi) * 100

		fmt.Printf("%8d samples: π ≈ %.6f (error: %.6f, %.3f%%)\n",
			n, pi, error, errorPct)
	}

	fmt.Println()
	fmt.Printf("Actual π:         %.6f\n", math.Pi)
}
