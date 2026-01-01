package rand

import "math"

// This file contains math/rand compatibility methods that build on top of
// the core Uint64() and Read() methods defined in r30r2.go.

// Uint32 returns a random uint32
func (r *RNG) Uint32() uint32 {
	return uint32(r.Uint64())
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
