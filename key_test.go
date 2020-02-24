package keybytes_test

import (
	"bytes"
	"database/sql"
	"math"
	"reflect"
	"testing"

	"github.com/hnakamur/keybytes"
)

func TestAppendNullString(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []sql.NullString{
			{Valid: false, String: ""},
			{Valid: true, String: ""},
			{Valid: true, String: "foo"},
			{Valid: true, String: "F\u00d4O\u0000bar"},
			{Valid: true, String: "\u0000foo"},
			{Valid: true, String: "foo\u0000"},
			{Valid: true, String: "f\x00\x00oo"},
			{Valid: true, String: "\x00"},
			{Valid: true, String: "\x00\x00"},
			{Valid: true, String: "\xff"},
			{Valid: true, String: "\xff\xff"},
			{Valid: true, String: "\x00\xff"},
			{Valid: true, String: "\x00\x00\xff\xff"},
		}
		for i, input := range testCases {
			b := keybytes.AppendNullString([]byte(nil), input)
			s, rest, err := keybytes.TakeNullString(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := s, input; !reflect.DeepEqual(got, want) {
				t.Errorf("case %d: string unmatch: got=%+v, want=%+v", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b sql.NullString
		}{
			{
				a: sql.NullString{Valid: false, String: ""},
				b: sql.NullString{Valid: true, String: ""},
			},
			{
				a: sql.NullString{Valid: true, String: ""},
				b: sql.NullString{Valid: true, String: "a"},
			},
			{
				a: sql.NullString{Valid: true, String: "a"},
				b: sql.NullString{Valid: true, String: "a\x00"},
			},
			{
				a: sql.NullString{Valid: true, String: "bar"},
				b: sql.NullString{Valid: true, String: "bb"},
			},
		}
		for i, tc := range testCases {
			a := keybytes.AppendNullString([]byte(nil), tc.a)
			b := keybytes.AppendNullString([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=0x%x, b=0x%x",
					i, got, want, a, b)
			}
		}
	})
}

func TestAppendString(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []string{
			"",
			"foo",
			"F\u00d4O\u0000bar",
			"\u0000foo",
			"foo\u0000",
			"f\x00\x00oo",
			"\x00",
			"\x00\x00",
			"\xff",
			"\xff\xff",
			"\x00\xff",
			"\x00\x00\xff\xff",
		}
		for i, input := range testCases {
			b := keybytes.AppendString([]byte(nil), input)
			s, rest, err := keybytes.TakeString(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := s, input; got != want {
				t.Errorf("case %d: string unmatch: got=%q, want=%q", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b string
		}{
			{a: "", b: "\x00"},
			{a: "", b: "a"},
			{a: "a", b: "a\x00"},
			{a: "bar", b: "bb"},
		}
		for i, tc := range testCases {
			a := keybytes.AppendString([]byte(nil), tc.a)
			b := keybytes.AppendString([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=0x%x, b=0x%x",
					i, got, want, a, b)
			}
		}
	})
}

func TestAppendNullInt32(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []sql.NullInt32{
			{Valid: false, Int32: 0},
			{Valid: true, Int32: 0},
			{Valid: true, Int32: math.MinInt32},
			{Valid: true, Int32: math.MinInt32 + 1},
			{Valid: true, Int32: -1},
			{Valid: true, Int32: 0},
			{Valid: true, Int32: 1},
			{Valid: true, Int32: math.MaxInt32 - 1},
			{Valid: true, Int32: math.MaxInt32},
		}
		for i, input := range testCases {
			b := keybytes.AppendNullInt32([]byte(nil), input)
			v, rest, err := keybytes.TakeNullInt32(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := v, input; !reflect.DeepEqual(got, want) {
				t.Errorf("case %d: string unmatch: got=%+v, want=%+v", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b sql.NullInt32
		}{
			{
				a: sql.NullInt32{Valid: false, Int32: 0},
				b: sql.NullInt32{Valid: true, Int32: 0},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: math.MinInt32},
				b: sql.NullInt32{Valid: true, Int32: math.MinInt32 + 1},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: -2},
				b: sql.NullInt32{Valid: true, Int32: -1},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: -1},
				b: sql.NullInt32{Valid: true, Int32: 0},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: 0},
				b: sql.NullInt32{Valid: true, Int32: 1},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: -1},
				b: sql.NullInt32{Valid: true, Int32: 1},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: -2},
				b: sql.NullInt32{Valid: true, Int32: 1},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: 1},
				b: sql.NullInt32{Valid: true, Int32: 2},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: math.MaxInt32 - 1},
				b: sql.NullInt32{Valid: true, Int32: math.MaxInt32},
			},
			{
				a: sql.NullInt32{Valid: true, Int32: math.MinInt32},
				b: sql.NullInt32{Valid: true, Int32: math.MaxInt32},
			},
		}
		for i, tc := range testCases {
			a := keybytes.AppendNullInt32([]byte(nil), tc.a)
			b := keybytes.AppendNullInt32([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=%+v, b=%+v",
					i, got, want, a, b)
			}
		}
	})
}

func TestAppendInt32(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []int32{
			math.MinInt32,
			math.MinInt32 + 1,
			-1,
			0,
			1,
			math.MaxInt32 - 1,
			math.MaxInt32,
		}
		for i, input := range testCases {
			b := keybytes.AppendInt32([]byte(nil), input)
			v, rest, err := keybytes.TakeInt32(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := v, input; got != want {
				t.Errorf("case %d: string unmatch: got=%q, want=%q", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b int32
		}{
			{a: math.MinInt32, b: math.MinInt32 + 1},
			{a: -2, b: -1},
			{a: -1, b: 0},
			{a: 0, b: 1},
			{a: -1, b: 1},
			{a: -2, b: 1},
			{a: 1, b: 2},
			{a: math.MaxInt32 - 1, b: math.MaxInt32},
			{a: math.MinInt32, b: math.MaxInt32},
		}
		for i, tc := range testCases {
			a := keybytes.AppendInt32([]byte(nil), tc.a)
			b := keybytes.AppendInt32([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=%d, b=%d",
					i, got, want, a, b)
			}
		}
	})
}
