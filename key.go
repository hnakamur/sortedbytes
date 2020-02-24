// Package keybytes provides feature for encoding and decoding key bytes.
// Two encoded key bytes of two different keys keeps the order
// so that you can use encoded key bytes for keys in a key-value-store
// which is capable to do range scans.
//
// The encoding in this package is a subset of the FDB Tuple layer typecodes encoding.
// https://github.com/apple/foundationdb/blob/92b41e3562e639e16dbe0142cc479a3304e9c08a/design/tuple.md
// https://activesphere.com/blog/2018/08/17/order-preserving-serialization
package keybytes

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strings"
)

const (
	typeCodeNull          = 0x00
	typeCodeUTF8String    = 0x02
	typeCodeNegativeInt64 = 0x0C
	typeCodeNegativeInt32 = 0x0F
	typeCodeIntZero       = 0x14
	typeCodePositiveInt32 = 0x19
	typeCodePositiveInt64 = 0x1C
	typeCodeFloat64       = 0x21
	typeCodeFalse         = 0x26
	typeCodeTrue          = 0x27
)

var errUnpexptedTypeCode = errors.New("unexpected type code")
var errValueOutOfRange = errors.New("value out of range")

// AppendNullString appends a sql.NullString value to dst.
func AppendNullString(dst []byte, value sql.NullString) []byte {
	if value.Valid {
		return AppendString(dst, value.String)
	}
	return append(dst, typeCodeNull)
}

// AppendString appends a string value to dst.
func AppendString(dst []byte, value string) []byte {
	dst = append(dst, typeCodeUTF8String)
	for {
		i := strings.IndexByte(value, '\x00')
		if i == -1 {
			return append(append(dst, value...), '\x00')
		}

		dst = append(append(dst, value[:i+1]...), '\xFF')
		value = value[i+1:]
	}
}

// TakeString takes a sql.NullString value from b and returns it and the rest of b.
func TakeNullString(b []byte) (value sql.NullString, rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return value, b, err
	}
	switch c {
	case typeCodeNull:
		return value, b[1:], nil
	case typeCodeUTF8String:
		s, rest, err := takeStringValue(rest)
		if err != nil {
			return value, b, err
		}
		return sql.NullString{Valid: true, String: s}, rest, nil
	default:
		return value, b, errUnpexptedTypeCode
	}
}

// TakeString takes a string value from b and returns it and the rest of b.
func TakeString(b []byte) (value string, rest []byte, err error) {
	rest, err = expectTypeCode(b, typeCodeUTF8String)
	if err != nil {
		return "", b, err
	}
	value, rest, err = takeStringValue(rest)
	if err != nil {
		return "", b, err
	}
	return value, rest, nil
}

func takeStringValue(src []byte) (value string, rest []byte, err error) {
	var out []byte
	for {
		i := bytes.IndexByte(src, '\x00')
		if i == -1 {
			return "", nil, io.ErrUnexpectedEOF
		}

		if i+1 < len(src) && src[i+1] == '\xFF' {
			if out == nil {
				out = []byte{}
			}
			out = append(out, src[:i+1]...)
			src = src[i+2:]
			continue
		}

		if out == nil {
			return string(src[:i]), src[i+1:], nil
		}
		return string(append(out, src[:i]...)), src[i+1:], nil
	}
}

// AppendNullInt32 appends a NullInt32 value to dst.
func AppendNullInt32(dst []byte, value sql.NullInt32) []byte {
	if value.Valid {
		return AppendInt32(dst, value.Int32)
	}
	return append(dst, typeCodeNull)
}

// AppendInt32 appends an int32 value to dst.
func AppendInt32(dst []byte, value int32) []byte {
	if value == 0 {
		return append(dst, typeCodeIntZero)
	}

	var b [4]byte
	if value > 0 {
		binary.BigEndian.PutUint32(b[:], uint32(value))
		return append(append(dst, typeCodePositiveInt32), b[:]...)
	}

	binary.BigEndian.PutUint32(b[:], math.MaxUint32-uint32(-value))
	return append(append(dst, typeCodeNegativeInt32), b[:]...)
}

// TakeNullInt32 takes a sql.NullInt32 value from b and returns it and the rest of b.
func TakeNullInt32(b []byte) (value sql.NullInt32, rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return value, b, err
	}
	if c == typeCodeNull {
		return value, b[1:], nil
	}
	var v int32
	v, rest, err = takeInt32Value(c, rest)
	if err != nil {
		return value, b, err
	}
	return sql.NullInt32{Valid: true, Int32: v}, rest, nil
}

// TakeInt32 takes an int32 value from b and returns it and the rest of b.
func TakeInt32(b []byte) (value int32, rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return 0, b, err
	}
	value, rest, err = takeInt32Value(c, rest)
	if err != nil {
		return 0, b, err
	}
	return value, rest, nil
}

func takeInt32Value(c byte, b []byte) (value int32, rest []byte, err error) {
	switch c {
	case typeCodeIntZero:
		return 0, b, nil
	case typeCodePositiveInt32:
		if len(b) < 4 {
			return 0, nil, io.ErrUnexpectedEOF
		}
		v := binary.BigEndian.Uint32(b[:4])
		if v > math.MaxInt32 {
			return 0, nil, errValueOutOfRange
		}
		return int32(v), b[4:], nil
	case typeCodeNegativeInt32:
		if len(b) < 4 {
			return 0, nil, io.ErrUnexpectedEOF
		}
		v := math.MaxUint32 - binary.BigEndian.Uint32(b[:4])
		if v > -math.MinInt32 {
			return 0, nil, errValueOutOfRange
		}
		return -int32(v), b[4:], nil
	default:
		return value, nil, errUnpexptedTypeCode
	}
}

// AppendNullInt64 appends a NullInt64 value to dst.
func AppendNullInt64(dst []byte, value sql.NullInt64) []byte {
	if value.Valid {
		return AppendInt64(dst, value.Int64)
	}
	return append(dst, typeCodeNull)
}

// AppendInt64 appends an int64 value to dst.
func AppendInt64(dst []byte, value int64) []byte {
	if value == 0 {
		return append(dst, typeCodeIntZero)
	}

	var b [8]byte
	if value > 0 {
		binary.BigEndian.PutUint64(b[:], uint64(value))
		return append(append(dst, typeCodePositiveInt64), b[:]...)
	}

	binary.BigEndian.PutUint64(b[:], math.MaxUint64-uint64(-value))
	return append(append(dst, typeCodeNegativeInt64), b[:]...)
}

// TakeNullInt64 takes a sql.NullInt64 value from b and returns it and the rest of b.
func TakeNullInt64(b []byte) (value sql.NullInt64, rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return value, b, err
	}
	if c == typeCodeNull {
		return value, b[1:], nil
	}
	var v int64
	v, rest, err = takeInt64Value(c, rest)
	if err != nil {
		return value, b, err
	}
	return sql.NullInt64{Valid: true, Int64: v}, rest, nil
}

// TakeInt64 takes an int64 value from b and returns it and the rest of b.
func TakeInt64(b []byte) (value int64, rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return 0, b, err
	}
	value, rest, err = takeInt64Value(c, rest)
	if err != nil {
		return 0, b, err
	}
	return value, rest, nil
}

func takeInt64Value(c byte, b []byte) (value int64, rest []byte, err error) {
	switch c {
	case typeCodeIntZero:
		return 0, b, nil
	case typeCodePositiveInt64:
		if len(b) < 8 {
			return 0, nil, io.ErrUnexpectedEOF
		}
		v := binary.BigEndian.Uint64(b[:8])
		if v > math.MaxInt64 {
			return 0, nil, errValueOutOfRange
		}
		return int64(v), b[8:], nil
	case typeCodeNegativeInt64:
		if len(b) < 8 {
			return 0, nil, io.ErrUnexpectedEOF
		}
		v := math.MaxUint64 - binary.BigEndian.Uint64(b[:8])
		if v > -math.MinInt64 {
			return 0, nil, errValueOutOfRange
		}
		return -int64(v), b[8:], nil
	default:
		return value, nil, errUnpexptedTypeCode
	}
}

// AppendNullFloat64 appends a NullFloat64 value to dst.
func AppendNullFloat64(dst []byte, value sql.NullFloat64) []byte {
	if value.Valid {
		return AppendFloat64(dst, value.Float64)
	}
	return append(dst, typeCodeNull)
}

// AppendFloat64 appends a float64 value to dst.
func AppendFloat64(dst []byte, value float64) []byte {
	v := math.Float64bits(value)
	if v&0x8000_0000_0000_0000 == 0 {
		v ^= 0x8000_0000_0000_0000
	} else {
		v ^= 0xffff_ffff_ffff_ffff
	}
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], v)
	return append(append(dst, typeCodeFloat64), b[:]...)
}

// TakeNullFloat64 takes a sql.NullFloat64 value from b and returns it and the rest of b.
func TakeNullFloat64(b []byte) (value sql.NullFloat64, rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return value, b, err
	}
	if c == typeCodeNull {
		return value, b[1:], nil
	}
	var v float64
	v, rest, err = takeFloat64Value(rest)
	if err != nil {
		return value, b, err
	}
	return sql.NullFloat64{Valid: true, Float64: v}, rest, nil
}

// TakeFloat64 takes an int64 value from b and returns it and the rest of b.
func TakeFloat64(b []byte) (value float64, rest []byte, err error) {
	rest, err = expectTypeCode(b, typeCodeFloat64)
	if err != nil {
		return 0, b, err
	}
	value, rest, err = takeFloat64Value(rest)
	if err != nil {
		return 0, b, err
	}
	return value, rest, nil
}

func takeFloat64Value(b []byte) (value float64, rest []byte, err error) {
	if len(b) < 8 {
		return 0, nil, io.ErrUnexpectedEOF
	}
	v := binary.BigEndian.Uint64(b[:8])
	if v&0x8000_0000_0000_0000 != 0 {
		v ^= 0x8000_0000_0000_0000
	} else {
		v ^= 0xffff_ffff_ffff_ffff
	}
	return math.Float64frombits(v), b[8:], nil
}

func expectTypeCode(b []byte, typeCode byte) (rest []byte, err error) {
	var c byte
	c, rest, err = takeTypeCode(b)
	if err != nil {
		return nil, err
	}
	if c != typeCode {
		return nil, errUnpexptedTypeCode
	}
	return rest, nil
}

func takeTypeCode(b []byte) (typeCode byte, rest []byte, err error) {
	if len(b) < 1 {
		return 0, nil, io.ErrUnexpectedEOF
	}
	return b[0], b[1:], nil
}
