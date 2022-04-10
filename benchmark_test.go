package base58_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/zergon321/base58"
)

func BenchmarkEncode(b *testing.B) {
	id := uuid.New()
	idBin, _ := id.MarshalBinary()

	for i := 0; i < b.N; i++ {
		base58.Encode(idBin)
	}
}

func BenchmarkEncodeToBuffer(b *testing.B) {
	id := uuid.New()
	idBin, _ := id.MarshalBinary()

	for i := 0; i < b.N; i++ {
		var buffer [22]byte
		base58.EncodeToBuffer(idBin, buffer[:])
	}
}

func BenchmarkDecode(b *testing.B) {
	str := "YaUV8Ysvm9EjpT8mmbTbwU"

	for i := 0; i < b.N; i++ {
		base58.Decode(str)
	}
}

func BenchmarkDecodeToBuffer(b *testing.B) {
	str := "YaUV8Ysvm9EjpT8mmbTbwU"

	for i := 0; i < b.N; i++ {
		var buffer [34]byte
		base58.DecodeToBuffer(str, buffer[:])
	}
}
