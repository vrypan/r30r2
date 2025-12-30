package rand

import (
	"encoding/binary"
	"math"
	"math/bits"
)

// RNG implements a 1D cellular automaton (Rule 30) on a circular 256-bit strip
// Optimized for 64-bit architectures using uint64 words
type RNG struct {
	state [4]uint64 // 256 bits as 4 × 64-bit words
	pos   int       // position in state for Uint64() extraction (0-3)
}

// New creates a new Rule 30 RNG from a seed
func New(seed uint64) *RNG {
	rng := &RNG{
		pos: 4, // Force Step() on first Uint64() call
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
	// Cache state words locally (cheaper than repeated array indexing)
	s0 := r.state[0]
	s1 := r.state[1]
	s2 := r.state[2]
	s3 := r.state[3]

	// Word 0: left neighbor from word 3, right neighbor from word 1
	left0 := (s0 >> 1) | (s3 << 63)
	right0 := (s0 << 1) | (s1 >> 63)
	new0 := left0 ^ (s0 | right0)

	// Word 1: left neighbor from word 0, right neighbor from word 2
	left1 := (s1 >> 1) | (s0 << 63)
	right1 := (s1 << 1) | (s2 >> 63)
	new1 := left1 ^ (s1 | right1)

	// Word 2: left neighbor from word 1, right neighbor from word 3
	left2 := (s2 >> 1) | (s1 << 63)
	right2 := (s2 << 1) | (s3 >> 63)
	new2 := left2 ^ (s2 | right2)

	// Word 3: left neighbor from word 2, right neighbor from word 0 (circular)
	left3 := (s3 >> 1) | (s2 << 63)
	right3 := (s3 << 1) | (s0 >> 63)
	new3 := left3 ^ (s3 | right3)

	// Update state with XOR rotation mixing (use primes for good diffusion)
	r.state[0] = new0 ^ bits.RotateLeft64(s0, 13)
	r.state[1] = new1 ^ bits.RotateLeft64(s1, 17)
	r.state[2] = new2 ^ bits.RotateLeft64(s2, 23)
	r.state[3] = new3 ^ bits.RotateLeft64(s3, 29)
}

// Read implements io.Reader interface
// Reads in 8-byte (uint64) chunks. Each Step() generates 32 bytes (4 uint64s),
// so every 4 Read() calls of 8 bytes fully utilizes one Step() with no waste.
func (r *RNG) Read(buf []byte) (n int, err error) {
	i := 0
	limit := len(buf)

	// Fill full uint64 chunks directly into the destination buffer
	for limit-i >= 8 {
		if r.pos >= 4 {
			r.Step()
			r.pos = 0
		}
		binary.LittleEndian.PutUint64(buf[i:], r.state[r.pos])
		r.pos++
		i += 8
	}

	// Handle the remaining tail bytes, if any, without creating a temporary buffer
	if rem := limit - i; rem > 0 {
		if r.pos >= 4 {
			r.Step()
			r.pos = 0
		}
		val := r.state[r.pos]
		r.pos++
		for j := 0; j < rem; j++ {
			buf[i+j] = byte(val)
			val >>= 8
		}
	}

	return limit, nil
}

// CopyState returns a copy of the current state
func (r *RNG) CopyState() [4]uint64 {
	return r.state
}

// Uint32 returns a random uint32
func (r *RNG) Uint32() uint32 {
	return uint32(r.Uint64())
}

// Uint64 returns a random uint64
// Optimized to extract directly from state without byte conversion
func (r *RNG) Uint64() uint64 {
	// Generate new state if we've exhausted all 4 uint64 values
	if r.pos >= 4 {
		r.Step()
		r.pos = 0
	}

	// Extract uint64 directly from state
	val := r.state[r.pos]
	r.pos++
	return val
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
