package main

import (
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	mathrand "math/rand"
	"os"
	"time"
)

func main() {
	var (
		rngType = flag.String("type", "math", "RNG type: 'math' or 'crypto'")
		bytes   = flag.Int("bytes", 10485760, "Number of bytes to generate")
		seed    = flag.Int64("seed", 0, "Seed for math/rand (0 = time-based)")
	)

	flag.Parse()

	buf := make([]byte, *bytes)

	switch *rngType {
	case "math":
		if *seed == 0 {
			*seed = time.Now().UnixNano()
		}
		rng := mathrand.New(mathrand.NewSource(*seed))
		n, err := rng.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Generated %d bytes using math/rand (seed: %d)\n", n, *seed)

	case "crypto":
		n, err := cryptorand.Read(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Generated %d bytes using crypto/rand\n", n)

	default:
		fmt.Fprintf(os.Stderr, "Unknown RNG type: %s\n", *rngType)
		os.Exit(1)
	}

	os.Stdout.Write(buf)
}
