package rule30

import (
	"encoding/binary"
	"math"
)

// RNG implements a 1D cellular automaton (Rule 30) on a circular 256-bit strip
// Optimized for 64-bit architectures using uint64 words
type RNG struct {
	state  [4]uint64 // 256 bits as 4 × 64-bit words
	buffer []byte
}

// New creates a new Rule 30 RNG from a seed
func New(seed uint64) *RNG {
	rng := &RNG{
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
func (r *RNG) Step() {
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
func (r *RNG) Read(buf []byte) (n int, err error) {
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
func (r *RNG) CopyState() [4]uint64 {
	return r.state
}

// Uint32 returns a random uint32
func (r *RNG) Uint32() uint32 {
	var buf [4]byte
	r.Read(buf[:])
	return binary.LittleEndian.Uint32(buf[:])
}

// Uint64 returns a random uint64
func (r *RNG) Uint64() uint64 {
	var buf [8]byte
	r.Read(buf[:])
	return binary.LittleEndian.Uint64(buf[:])
}

// Int63 returns a non-negative random int64 (0 to 2^63-1)
func (r *RNG) Int63() int64 {
	return int64(r.Uint64() & 0x7FFFFFFFFFFFFFFF)
}

// Int31 returns a non-negative random int32 (0 to 2^31-1)
func (r *RNG) Int31() int32 {
	return int32(r.Uint32() >> 1)
}

// Int returns a non-negative random int
func (r *RNG) Int() int {
	u := uint(r.Int63())
	return int(u << 1 >> 1) // clear sign bit if int is 32 bits
}

// Int63n returns a random int64 in [0, n)
// Panics if n <= 0
func (r *RNG) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	if n&(n-1) == 0 { // n is power of two
		return r.Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := r.Int63()
	for v > max {
		v = r.Int63()
	}
	return v % n
}

// Int31n returns a random int32 in [0, n)
// Panics if n <= 0
func (r *RNG) Int31n(n int32) int32 {
	if n <= 0 {
		panic("invalid argument to Int31n")
	}
	if n&(n-1) == 0 { // n is power of two
		return r.Int31() & (n - 1)
	}
	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := r.Int31()
	for v > max {
		v = r.Int31()
	}
	return v % n
}

// Intn returns a random int in [0, n)
// Panics if n <= 0
func (r *RNG) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}
	if n <= 1<<31-1 {
		return int(r.Int31n(int32(n)))
	}
	return int(r.Int63n(int64(n)))
}

// Float64 returns a random float64 in [0.0, 1.0)
func (r *RNG) Float64() float64 {
	// Use 53 bits of precision (same as math/rand)
	return float64(r.Int63()>>11) / (1 << 52)
}

// Float32 returns a random float32 in [0.0, 1.0)
func (r *RNG) Float32() float32 {
	// Use 24 bits of precision
	return float32(r.Int31()>>7) / (1 << 24)
}

// NormFloat64 returns a normally distributed float64 with mean 0 and stddev 1
// Uses the Box-Muller transform
func (r *RNG) NormFloat64() float64 {
	for {
		u := 2*r.Float64() - 1
		v := 2*r.Float64() - 1
		s := u*u + v*v
		if s < 1 && s != 0 {
			return u * math.Sqrt(-2*math.Log(s)/s)
		}
	}
}

// ExpFloat64 returns an exponentially distributed float64 with rate 1
func (r *RNG) ExpFloat64() float64 {
	for {
		u := r.Float64()
		if u > 0 {
			return -math.Log(u)
		}
	}
}
