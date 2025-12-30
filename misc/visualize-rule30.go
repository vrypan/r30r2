package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vrypan/rule30rnd/rand"
)

func main() {
	var (
		seed        = flag.Uint64("seed", 1, "RNG seed")
		generations = flag.Int("generations", 50, "Number of generations to display")
		width       = flag.Int("width", 256, "Width in bits (max 256)")
		char0       = flag.String("char0", "░", "Character for 0 bits")
		char1       = flag.String("char1", "█", "Character for 1 bits")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Rule 30 Visualization - Cellular Automaton Evolution

Shows the evolution of Rule 30's 256-bit state over multiple generations.
Each row represents one generation, with each character showing a single bit.

Usage:
  visualize-rule30 [options]

Options:
  --seed N          RNG seed (default: 1)
  --generations N   Number of generations to show (default: 50)
  --width N         Width in bits to display, max 256 (default: 256)
  --char0 C         Character for 0 bits (default: ░)
  --char1 C         Character for 1 bits (default: █)

Examples:
  # Default visualization (50 generations, full width)
  visualize-rule30

  # Narrower view for terminals
  visualize-rule30 --width=128 --generations=100

  # ASCII characters
  visualize-rule30 --char0=" " --char1="*"

  # Different seed
  visualize-rule30 --seed=12345

  # Compact 0/1 display
  visualize-rule30 --char0="0" --char1="1"
`)
	}

	flag.Parse()

	if *width < 1 || *width > 256 {
		fmt.Fprintf(os.Stderr, "Error: width must be between 1 and 256\n")
		os.Exit(1)
	}

	// Create RNG
	rng := rand.New(*seed)

	// Print header
	fmt.Printf("Rule 30 Visualization\n")
	fmt.Printf("Seed: %d | Generations: %d | Width: %d bits\n", *seed, *generations, *width)
	fmt.Printf("Evolution rule: new_bit = left XOR (center OR right)\n")
	fmt.Println()

	// Display generations
	for gen := 0; gen < *generations; gen++ {
		// Get current state
		state := rng.CopyState()

		// Print generation number (padded)
		fmt.Printf("%4d │ ", gen)

		// Print bits
		bitsDisplayed := 0
		for wordIdx := 0; wordIdx < 4 && bitsDisplayed < *width; wordIdx++ {
			word := state[wordIdx]
			bitsInThisWord := 64
			if bitsDisplayed+bitsInThisWord > *width {
				bitsInThisWord = *width - bitsDisplayed
			}

			for bit := 0; bit < bitsInThisWord; bit++ {
				if word&1 == 1 {
					fmt.Print(*char1)
				} else {
					fmt.Print(*char0)
				}
				word >>= 1
			}
			bitsDisplayed += bitsInThisWord
		}
		fmt.Println()

		// Advance to next generation by consuming the state
		// Read 32 bytes (256 bits) to force evolution to next state
		buf := make([]byte, 32)
		rng.Read(buf)
	}

	fmt.Println()
	fmt.Printf("Displayed %d generations of Rule 30 evolution\n", *generations)
}
