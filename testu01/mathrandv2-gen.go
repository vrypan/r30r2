package main

import (
	"encoding/binary"
	"flag"
	mathrandv2 "math/rand/v2"
	"os"
)

func main() {
	bytes := flag.Int64("bytes", 1024*1024*1024, "Number of bytes to generate (0 = unlimited)")
	seed := flag.Uint64("seed", 12345, "Random seed")
	flag.Parse()

	rng := mathrandv2.NewPCG(*seed, 0)

	if *bytes == 0 {
		// Unlimited mode: stream chunks until pipe breaks
		buf := make([]byte, 1024*1024) // 1MB chunks
		for {
			// Fill buffer with uint64 values
			for i := 0; i < len(buf); i += 8 {
				binary.LittleEndian.PutUint64(buf[i:], rng.Uint64())
			}

			// Write, exit gracefully if pipe closes
			_, err := os.Stdout.Write(buf)
			if err != nil {
				os.Exit(0)
			}
		}
	} else {
		// Fixed size mode
		buf := make([]byte, 8192)
		var written int64
		for written < *bytes {
			toWrite := int64(len(buf))
			if written+toWrite > *bytes {
				toWrite = *bytes - written
			}

			// Fill buffer with uint64 values
			for i := 0; i < int(toWrite); i += 8 {
				if i+8 <= int(toWrite) {
					binary.LittleEndian.PutUint64(buf[i:], rng.Uint64())
				} else {
					// Handle remaining bytes
					val := rng.Uint64()
					for j := i; j < int(toWrite); j++ {
						buf[j] = byte(val)
						val >>= 8
					}
					break
				}
			}

			n, err := os.Stdout.Write(buf[:toWrite])
			if err != nil {
				os.Exit(1)
			}
			written += int64(n)
		}
	}
}
