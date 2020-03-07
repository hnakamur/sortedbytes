package fuzztest

import (
	"bytes"
	"database/sql"
	"reflect"

	"github.com/hnakamur/sortedbytes"
)

func FuzzTakeString(data []byte) int {
	v, rest, err := sortedbytes.TakeString(data)
	if err != nil {
		if v != "" {
			panic("v != \"\" on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeNullString(data []byte) int {
	v, rest, err := sortedbytes.TakeNullString(data)
	if err != nil {
		if !reflect.DeepEqual(v, sql.NullString{Valid: false, String: ""}) {
			panic("v != sql.NullString{Valid: false, String: \"\"} on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeInt32(data []byte) int {
	v, rest, err := sortedbytes.TakeInt32(data)
	if err != nil {
		if v != 0 {
			panic("v != 0 on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeNullInt32(data []byte) int {
	v, rest, err := sortedbytes.TakeNullInt32(data)
	if err != nil {
		if !reflect.DeepEqual(v, sql.NullInt32{Valid: false, Int32: 0}) {
			panic("v != sql.NullInt32{Valid: false, Int32: 0} on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeInt64(data []byte) int {
	v, rest, err := sortedbytes.TakeInt64(data)
	if err != nil {
		if v != 0 {
			panic("v != 0 on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeNullInt64(data []byte) int {
	v, rest, err := sortedbytes.TakeNullInt64(data)
	if err != nil {
		if !reflect.DeepEqual(v, sql.NullInt64{Valid: false, Int64: 0}) {
			panic("v != sql.NullInt64{Valid: false, Int64: 0} on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeFloat64(data []byte) int {
	v, rest, err := sortedbytes.TakeFloat64(data)
	if err != nil {
		if v != 0 {
			panic("v != 0 on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeNullFloat64(data []byte) int {
	v, rest, err := sortedbytes.TakeNullFloat64(data)
	if err != nil {
		if !reflect.DeepEqual(v, sql.NullFloat64{Valid: false, Float64: 0}) {
			panic("v != sql.NullFloat64{Valid: false, Float64: 0} on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeBool(data []byte) int {
	v, rest, err := sortedbytes.TakeBool(data)
	if err != nil {
		if v != false {
			panic("v != false on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}

func FuzzTakeNullBool(data []byte) int {
	v, rest, err := sortedbytes.TakeNullBool(data)
	if err != nil {
		if !reflect.DeepEqual(v, sql.NullBool{Valid: false, Bool: false}) {
			panic("v != sql.NullBool{Valid: false, Bool: false} on error")
		}
		if !bytes.Equal(rest, data) {
			panic("!bytes.Equal(rest, data) on error")
		}
		return 0
	}
	if len(rest) >= len(data) {
		panic("len(rest) >= len(data) on success")
	}
	return 1
}
