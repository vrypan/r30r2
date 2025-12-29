package main

import (
	"encoding/binary"
)

// Rule30RNG implements a 1D cellular automaton (Rule 30) on a circular 256-bit strip
// Optimized for 64-bit architectures using uint64 words
type Rule30RNG struct {
	state  [4]uint64 // 256 bits as 4 × 64-bit words
	buffer []byte
}

// NewRule30 creates a new Rule 30 RNG from a seed
func NewRule30(seed uint64) *Rule30RNG {
	rng := &Rule30RNG{
		buffer: make([]byte, 0, 32),
	}

	// Initialize state from seed
	// Use seed to create varied initial patterns
	rng.state[0] = seed
	rng.state[1] = seed ^ 0x5555555555555555
	rng.state[2] = seed ^ 0xAAAAAAAAAAAAAAAA
	rng.state[3] = seed ^ 0x3333333333333333

	return rng
}

// Step applies Rule 30 to all 256 bits in parallel (64-bit word-wise)
// Rule 30: new_bit = left XOR (center OR right)
// Optimized for 64-bit architecture: processes 64 bits at once, fully unrolled
func (r *Rule30RNG) Step() {
	// Fully unrolled loop for maximum performance
	// Each operation processes 64 bits in parallel

	// Word 0: left neighbor from word 3, right neighbor from word 1
	left0 := (r.state[0] >> 1) | (r.state[3] << 63)
	right0 := (r.state[0] << 1) | (r.state[1] >> 63)
	new0 := left0 ^ (r.state[0] | right0)

	// Word 1: left neighbor from word 0, right neighbor from word 2
	left1 := (r.state[1] >> 1) | (r.state[0] << 63)
	right1 := (r.state[1] << 1) | (r.state[2] >> 63)
	new1 := left1 ^ (r.state[1] | right1)

	// Word 2: left neighbor from word 1, right neighbor from word 3
	left2 := (r.state[2] >> 1) | (r.state[1] << 63)
	right2 := (r.state[2] << 1) | (r.state[3] >> 63)
	new2 := left2 ^ (r.state[2] | right2)

	// Word 3: left neighbor from word 2, right neighbor from word 0 (circular)
	left3 := (r.state[3] >> 1) | (r.state[2] << 63)
	right3 := (r.state[3] << 1) | (r.state[0] >> 63)
	new3 := left3 ^ (r.state[3] | right3)

	// Update state
	r.state[0] = new0
	r.state[1] = new1
	r.state[2] = new2
	r.state[3] = new3
}

// Read implements io.Reader interface
func (r *Rule30RNG) Read(buf []byte) (n int, err error) {
	needed := len(buf)
	copied := 0

	// Use leftover buffer first
	if len(r.buffer) > 0 {
		n := copy(buf, r.buffer)
		r.buffer = r.buffer[n:]
		copied += n
		if copied >= needed {
			return copied, nil
		}
	}

	// Generate more bytes as needed
	for copied < needed {
		// Run single Rule 30 iteration (produces 256 new bits = 32 bytes)
		r.Step()

		// Extract all 32 bytes from current state
		// Convert uint64 words to bytes (little-endian)
		var output [32]byte
		binary.LittleEndian.PutUint64(output[0:8], r.state[0])
		binary.LittleEndian.PutUint64(output[8:16], r.state[1])
		binary.LittleEndian.PutUint64(output[16:24], r.state[2])
		binary.LittleEndian.PutUint64(output[24:32], r.state[3])

		n := copy(buf[copied:], output[:])
		copied += n

		// Save leftover bytes
		if n < len(output) {
			r.buffer = append(r.buffer[:0], output[n:]...)
		}
	}

	return copied, nil
}

// CopyState returns a copy of the current state
func (r *Rule30RNG) CopyState() [4]uint64 {
	return r.state
}
