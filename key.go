// Package keybytes provides feature for encoding and decoding key bytes.
// Two encoded key bytes of two different keys keeps the order
// so that you can use encoded key bytes for keys in a key-value-store
// which is capable to do range scans.
package keybytes

import (
	"bytes"
	"encoding/binary"
	"math"
)

// AppendRaw appends a raw bytes value to dst.
func AppendRaw(dst, value []byte) []byte {
	return append(dst, value...)
}

// AppendByte appends a byte value to dst.
func AppendByte(dst []byte, value byte) []byte {
	return append(dst, value)
}

// AppendByte appends a byte value to dst for descending order range scan.
func AppendByteDesc(dst []byte, value byte) []byte {
	return append(dst, value^math.MaxUint8)
}

// AppendUint16 appends a uint16 value to dst.
func AppendUint16(dst []byte, value uint16) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], value)
	return append(dst, b[:]...)
}

// AppendUint16Desc appends a uint16 value to dst for descending order range scan.
func AppendUint16Desc(dst []byte, value uint16) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], value^math.MaxUint16)
	return append(dst, b[:]...)
}

// AppendUint32 appends a uint32 value to dst.
func AppendUint32(dst []byte, value uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], value)
	return append(dst, b[:]...)
}

// AppendUint32Desc appends a uint32 value to dst for descending order range scan.
func AppendUint32Desc(dst []byte, value uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], value^math.MaxUint32)
	return append(dst, b[:]...)
}

// AppendUint64 appends a uint64 value to dst.
func AppendUint64(dst []byte, value uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], value)
	return append(dst, b[:]...)
}

// AppendUint64Desc appends a uint64 value to dst for descending order range scan.
func AppendUint64Desc(dst []byte, value uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], value^math.MaxUint64)
	return append(dst, b[:]...)
}

// AppendStringNul appends a string value and a nul byte '\x00' to dst.
// The value must not contain a nul byte in supposed usage of this function,
// but ensuring that is caller's responsibility.
func AppendStringNul(dst []byte, value string) []byte {
	dst = append(dst, []byte(value)...)
	return append(dst, '\x00')
}

// TakeRaw takes the first n bytes from the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than n.
func TakeRaw(key []byte, n int) (value []byte, rest []byte) {
	return key[:n], key[n:]
}

// TakeByte takes the first byte from the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 1.
func TakeByte(key []byte) (value byte, rest []byte) {
	return key[0], key[1:]
}

// TakeByteDesc takes the first byte for descending order range scan from the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 1.
func TakeByteDesc(key []byte) (value byte, rest []byte) {
	return key[0] ^ math.MaxUint8, key[1:]
}

// TakeUint16 takes a uint16 value from the beginning of the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 2.
func TakeUint16(key []byte) (value uint16, rest []byte) {
	return binary.BigEndian.Uint16(key[:2]), key[2:]
}

// TakeUint16Desc takes a uint16 value for descending order range scan from the beginning of the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 2.
func TakeUint16Desc(key []byte) (value uint16, rest []byte) {
	return binary.BigEndian.Uint16(key[:2]) ^ math.MaxUint16, key[2:]
}

// TakeUint32 takes a uint32 value from the beginning of the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 4.
func TakeUint32(key []byte) (value uint32, rest []byte) {
	return binary.BigEndian.Uint32(key[:4]), key[4:]
}

// TakeUint32Desc takes a uint32 value for descending order range scan from the beginning of the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 4.
func TakeUint32Desc(key []byte) (value uint32, rest []byte) {
	return binary.BigEndian.Uint32(key[:4])^math.MaxUint32, key[4:]
}

// TakeUint64 takes a uint64 value from the beginning of the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 8.
func TakeUint64(key []byte) (value uint64, rest []byte) {
	return binary.BigEndian.Uint64(key[:8]), key[8:]
}

// TakeUint64Desc takes a uint64 value for descending order range scan from the beginning of the key and returns it and the rest of the key.
// It panics if the length of the key is smaller than 4.
func TakeUint64Desc(key []byte) (value uint64, rest []byte) {
	return binary.BigEndian.Uint64(key[:8])^math.MaxUint64, key[8:]
}

// TakeStringNul takes a string value before the first nul byte '\x00' and the first nul byte from the beginning of the key and returns the string and the rest of the key.
// It panics if the key does not contain a nul byte '\x00'
func TakeStringNul(key []byte) (value string, rest []byte) {
	i := bytes.IndexByte(key, '\x00')
	if i == -1 {
		panic("a nul byte not found in key")
	}
	return string(key[:i]), key[i+1:]
}
