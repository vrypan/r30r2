package rand

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"io"
	mathrand "math/rand"
	mathrandv2 "math/rand/v2"
	"testing"
)

// ====================
// Rule30RNG Benchmarks
// ====================

func BenchmarkRule30_Read32KB(b *testing.B) {
	rng := New(12345)
	buf := make([]byte, 32<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := rng.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRule30_Read1KB(b *testing.B) {
	rng := New(67890)
	buf := make([]byte, 1<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := rng.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRule30_Uint64(b *testing.B) {
	rng := New(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.Uint64()
	}
}

// ====================
// math/rand Benchmarks
// ====================

func BenchmarkMathRand_Read32KB(b *testing.B) {
	rng := mathrand.New(mathrand.NewSource(12345))
	buf := make([]byte, 32<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := rng.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMathRand_Read1KB(b *testing.B) {
	rng := mathrand.New(mathrand.NewSource(67890))
	buf := make([]byte, 1<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := rng.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMathRand_Uint64(b *testing.B) {
	rng := mathrand.New(mathrand.NewSource(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.Uint64()
	}
}

// =======================
// math/rand/v2 Benchmarks
// =======================

func BenchmarkMathRandV2_Read32KB(b *testing.B) {
	rng := mathrandv2.NewPCG(12345, 67890)
	buf := make([]byte, 32<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// math/rand/v2 doesn't have Read(), so read uint64s
		for j := 0; j < len(buf); j += 8 {
			binary.LittleEndian.PutUint64(buf[j:], rng.Uint64())
		}
	}
}

func BenchmarkMathRandV2_Read1KB(b *testing.B) {
	rng := mathrandv2.NewPCG(12345, 67890)
	buf := make([]byte, 1<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// math/rand/v2 doesn't have Read(), so read uint64s
		for j := 0; j < len(buf); j += 8 {
			binary.LittleEndian.PutUint64(buf[j:], rng.Uint64())
		}
	}
}

func BenchmarkMathRandV2_Uint64(b *testing.B) {
	rng := mathrandv2.NewPCG(42, 12345)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.Uint64()
	}
}

// ====================
// crypto/rand Benchmarks
// ====================

func BenchmarkCryptoRand_Read32KB(b *testing.B) {
	buf := make([]byte, 32<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := cryptorand.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCryptoRand_Read1KB(b *testing.B) {
	buf := make([]byte, 1<<10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := cryptorand.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCryptoRand_Uint64(b *testing.B) {
	buf := make([]byte, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := io.ReadFull(cryptorand.Reader, buf); err != nil {
			b.Fatal(err)
		}
		_ = binary.LittleEndian.Uint64(buf)
	}
}
