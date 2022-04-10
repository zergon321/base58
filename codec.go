package base58

import (
	"fmt"
	"unsafe"
)

var (
	alphabetEncode = [58]byte{
		49, 50, 51, 52, 53, 54, 55, 56, 57, 65, 66,
		67, 68, 69, 70, 71, 72, 74, 75, 76, 77, 78,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
		97, 98, 99, 100, 101, 102, 103, 104, 105, 106,
		107, 109, 110, 111, 112, 113, 114, 115, 116,
		117, 118, 119, 120, 121, 122}

	alphabetDecode = [128]int8{
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, 0, 1, 2, 3, 4, 5, 6, 7, 8,
		-1, -1, -1, -1, -1, -1, -1, 9,
		10, 11, 12, 13, 14, 15, 16, -1,
		17, 18, 19, 20, 21, -1, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, -1, -1, -1, -1, -1, -1, 33,
		34, 35, 36, 37, 38, 39, 40, 41,
		42, 43, -1, 44, 45, 46, 47, 48,
		49, 50, 51, 52, 53, 54, 55, 56,
		57, -1, -1, -1, -1, -1}
)

// byteSlice2String converts the byte
// slice into a string without copying.
func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// EncodeToBuffer encodes the binary data
// into base58 using the pre-allocated buffer.
func EncodeToBuffer(bin, buffer []byte) error {
	size := len(bin)
	zcount := 0

	for zcount < size && bin[zcount] == 0 {
		zcount++
	}

	// It is crucial to make this as short as possible, especially for
	// the usual case of bitcoin addrs
	size = zcount +
		// This is an integer simplification of
		// ceil(log(256)/log(58))
		(size-zcount)*555/406 + 1

	if len(buffer) < size {
		return fmt.Errorf("insufficient out buffer size")
	}

	var i, high int
	var carry uint32
	high = size - 1

	for _, b := range bin {
		i = size - 1

		for carry = uint32(b); i > high || carry != 0; i-- {
			carry = carry + 256*uint32(buffer[i])
			buffer[i] = byte(carry % 58)
			carry /= 58
		}

		high = i
	}

	// Determine the additional "zero-gap" in the buffer (aside from zcount)
	for i = zcount; i < size && buffer[i] == 0; i++ {
	}

	// Now encode the values with actual alphabet in-place
	val := buffer[i-zcount:]
	size = len(val)

	for i = 0; i < size; i++ {
		buffer[i] = alphabetEncode[val[i]]
	}

	return nil
}

// Encode encodes the byte slice
// into a base58 string.
func Encode(bin []byte) string {
	size := len(bin)
	zcount := 0

	for zcount < size && bin[zcount] == 0 {
		zcount++
	}

	// It is crucial to make this as short as possible, especially for
	// the usual case of bitcoin addrs
	size = zcount +
		// This is an integer simplification of
		// ceil(log(256)/log(58))
		(size-zcount)*555/406 + 1
	out := make([]byte, size)

	var i, high int
	var carry uint32
	high = size - 1

	for _, b := range bin {
		i = size - 1

		for carry = uint32(b); i > high || carry != 0; i-- {
			carry = carry + 256*uint32(out[i])
			out[i] = byte(carry % 58)
			carry /= 58
		}

		high = i
	}

	// Determine the additional "zero-gap" in the buffer (aside from zcount)
	for i = zcount; i < size && out[i] == 0; i++ {
	}

	// Now encode the values with actual alphabet in-place
	val := out[i-zcount:]
	size = len(val)

	for i = 0; i < size; i++ {
		out[i] = alphabetEncode[val[i]]
	}

	return byteSlice2String(out[:size])
}

// DecodeToBuffer decodes the base58 string
// and writes the result to the pre-allocated buffer.
func DecodeToBuffer(str string, buffer []byte) (int, int, error) {
	if len(str) == 0 {
		return -1, -1, fmt.Errorf("zero length string")
	}

	zero := alphabetEncode[0]
	b58sz := len(str)
	var zcount int

	for i := 0; i < b58sz && str[i] == zero; i++ {
		zcount++
	}

	var t, c uint64
	// the 32bit algo stretches the result up to 2 times
	outi := make([]uint32, (b58sz+3)/4)

	if len(buffer) < 2*((b58sz*406/555)+1) {
		return -1, -1, fmt.Errorf("insufficient out buffer size")
	}

	for _, r := range str {
		if r > 127 {
			return -1, -1, fmt.Errorf("high-bit set on invalid digit")
		}

		if alphabetDecode[r] == -1 {
			return -1, -1, fmt.Errorf("invalid base58 digit (%q)", r)
		}

		c = uint64(alphabetDecode[r])

		for j := len(outi) - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + c
			c = t >> 32
			outi[j] = uint32(t & 0xffffffff)
		}
	}

	// initial mask depends on b58sz, on further loops it always starts at 24 bits
	mask := (uint(b58sz%4) * 8)

	if mask == 0 {
		mask = 32
	}

	mask -= 8
	outLen := 0

	for j := 0; j < len(outi); j++ {
		for mask < 32 { // loop relies on uint overflow
			buffer[outLen] = byte(outi[j] >> mask)
			mask -= 8
			outLen++
		}

		mask = 24
	}

	// find the most significant byte post-decode, if any
	for msb := zcount; msb < len(buffer); msb++ {
		if buffer[msb] > 0 {
			return msb - zcount, outLen, nil
		}
	}

	// it's all zeroes
	return 0, outLen, nil
}

// Decode decodes the given string from
// the base58 format into a byte slice.
func Decode(str string) ([]byte, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("zero length string")
	}

	zero := alphabetEncode[0]
	b58sz := len(str)
	var zcount int

	for i := 0; i < b58sz && str[i] == zero; i++ {
		zcount++
	}

	var t, c uint64
	// the 32bit algo stretches the result up to 2 times
	binu := make([]byte, 2*((b58sz*406/555)+1))
	outi := make([]uint32, (b58sz+3)/4)

	for _, r := range str {
		if r > 127 {
			return nil, fmt.Errorf("high-bit set on invalid digit")
		}

		if alphabetDecode[r] == -1 {
			return nil, fmt.Errorf("invalid base58 digit (%q)", r)
		}

		c = uint64(alphabetDecode[r])

		for j := len(outi) - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + c
			c = t >> 32
			outi[j] = uint32(t & 0xffffffff)
		}
	}

	// initial mask depends on b58sz, on further loops it always starts at 24 bits
	mask := (uint(b58sz%4) * 8)

	if mask == 0 {
		mask = 32
	}

	mask -= 8
	outLen := 0

	for j := 0; j < len(outi); j++ {
		for mask < 32 { // loop relies on uint overflow
			binu[outLen] = byte(outi[j] >> mask)
			mask -= 8
			outLen++
		}

		mask = 24
	}

	// find the most significant byte post-decode, if any
	for msb := zcount; msb < len(binu); msb++ {
		if binu[msb] > 0 {
			return binu[msb-zcount : outLen], nil
		}
	}

	// it's all zeroes
	return binu[:outLen], nil
}
