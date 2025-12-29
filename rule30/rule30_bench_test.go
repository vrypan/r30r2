package rule30

import "testing"

func BenchmarkRead32KB(b *testing.B) {
	rng := New(12345)
	buf := make([]byte, 32<<10)
	for i := 0; i < b.N; i++ {
		if _, err := rng.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRead1KB(b *testing.B) {
	rng := New(67890)
	buf := make([]byte, 1<<10)
	for i := 0; i < b.N; i++ {
		if _, err := rng.Read(buf); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUint64(b *testing.B) {
	rng := New(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.Uint64()
	}
}
