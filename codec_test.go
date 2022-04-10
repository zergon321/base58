package base58_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zergon321/base58"
)

func TestEncode(t *testing.T) {
	var buffer [22]byte
	id := uuid.New()
	idBin, err := id.MarshalBinary()
	assert.Nil(t, err)

	res := base58.Encode(idBin)
	err = base58.EncodeToBuffer(idBin, buffer[:])
	assert.Nil(t, err)
	assert.Equal(t, res, string(buffer[:]))
}

func TestDecode(t *testing.T) {
	var buffer [22]byte
	id := uuid.New()
	idBin, err := id.MarshalBinary()
	assert.Nil(t, err)

	res := base58.Encode(idBin)
	err = base58.EncodeToBuffer(idBin, buffer[:])
	assert.Nil(t, err)
	encStr := string(buffer[:])
	assert.Equal(t, res, encStr)

	var decBuffer [34]byte
	start, end, err := base58.DecodeToBuffer(
		encStr, decBuffer[:])
	assert.Positive(t, start)
	assert.Positive(t, end)
	assert.Nil(t, err)

	dec := decBuffer[start:end]

	for i := 0; i < len(idBin); i++ {
		assert.Equal(t, idBin[i], dec[i])
	}
}
