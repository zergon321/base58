package base58_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/zergon321/base58"
)

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := uuid.New()
		idBin, _ := id.MarshalBinary()
		base58.Encode(idBin)
	}
}

func BenchmarkSlice(b *testing.B) {
	slice := make([]byte, 10_000)
	_ = slice[:]
}
