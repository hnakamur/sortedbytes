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

func TestTakeNullString(t *testing.T) {
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x02"),
			[]byte("\x02foo"),
			[]byte("\x02\x00\xffa"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeNullString(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendString(t *testing.T) {
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

func TestTakeString(t *testing.T) {
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x02"),
			[]byte("\x02foo"),
			[]byte("\x02\x00\xffa"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeString(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendNullInt32(t *testing.T) {
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

func TestTakeNullInt32(t *testing.T) {
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x19"),
			[]byte("\x19\x01\x02\x03"),
			[]byte("\x19\x80\x00\x00\x00"),
			[]byte("\x0f\x7f\xff\xff\xfe"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeNullInt32(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendInt32(t *testing.T) {
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

func TestTakeInt32(t *testing.T) {
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x19"),
			[]byte("\x0f"),
			[]byte("\x19\x01\x02\x03"),
			[]byte("\x0f\x01\x02\x03"),
			[]byte("\x19\x80\x00\x00\x00"),
			[]byte("\x0f\x7f\xff\xff\xfe"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeInt32(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendNullInt64(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b sql.NullInt64
		}{
			{
				a: sql.NullInt64{Valid: false, Int64: 0},
				b: sql.NullInt64{Valid: true, Int64: 0},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: math.MinInt64},
				b: sql.NullInt64{Valid: true, Int64: math.MinInt64 + 1},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: -2},
				b: sql.NullInt64{Valid: true, Int64: -1},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: -1},
				b: sql.NullInt64{Valid: true, Int64: 0},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: 0},
				b: sql.NullInt64{Valid: true, Int64: 1},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: -1},
				b: sql.NullInt64{Valid: true, Int64: 1},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: -2},
				b: sql.NullInt64{Valid: true, Int64: 1},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: 1},
				b: sql.NullInt64{Valid: true, Int64: 2},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: math.MaxInt64 - 1},
				b: sql.NullInt64{Valid: true, Int64: math.MaxInt64},
			},
			{
				a: sql.NullInt64{Valid: true, Int64: math.MinInt64},
				b: sql.NullInt64{Valid: true, Int64: math.MaxInt64},
			},
		}
		for i, tc := range testCases {
			a := keybytes.AppendNullInt64([]byte(nil), tc.a)
			b := keybytes.AppendNullInt64([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=%+v, b=%+v",
					i, got, want, a, b)
			}
		}
	})
}

func TestTakeNullInt64(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []sql.NullInt64{
			{Valid: false, Int64: 0},
			{Valid: true, Int64: 0},
			{Valid: true, Int64: math.MinInt64},
			{Valid: true, Int64: math.MinInt64 + 1},
			{Valid: true, Int64: -1},
			{Valid: true, Int64: 0},
			{Valid: true, Int64: 1},
			{Valid: true, Int64: math.MaxInt64 - 1},
			{Valid: true, Int64: math.MaxInt64},
		}
		for i, input := range testCases {
			b := keybytes.AppendNullInt64([]byte(nil), input)
			v, rest, err := keybytes.TakeNullInt64(b)
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x0c"),
			[]byte("\x1c"),
			[]byte("\x0c\x01\x02\x03\x04\x05\x06\x07"),
			[]byte("\x1c\x01\x02\x03\x04\x05\x06\x07"),
			[]byte("\x0c\x7f\xff\xff\xff\xff\xff\xff\xfe"),
			[]byte("\x1c\x80\x00\x00\x00\x00\x00\x00\x00"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeNullInt64(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendInt64(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b int64
		}{
			{a: math.MinInt64, b: math.MinInt64 + 1},
			{a: -2, b: -1},
			{a: -1, b: 0},
			{a: 0, b: 1},
			{a: -1, b: 1},
			{a: -2, b: 1},
			{a: 1, b: 2},
			{a: math.MaxInt64 - 1, b: math.MaxInt64},
			{a: math.MinInt64, b: math.MaxInt64},
		}
		for i, tc := range testCases {
			a := keybytes.AppendInt64([]byte(nil), tc.a)
			b := keybytes.AppendInt64([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=%d, b=%d",
					i, got, want, a, b)
			}
		}
	})
}

func TestTakeInt64(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []int64{
			math.MinInt64,
			math.MinInt64 + 1,
			-1,
			0,
			1,
			math.MaxInt64 - 1,
			math.MaxInt64,
		}
		for i, input := range testCases {
			b := keybytes.AppendInt64([]byte(nil), input)
			v, rest, err := keybytes.TakeInt64(b)
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x0c"),
			[]byte("\x1c"),
			[]byte("\x0c\x01\x02\x03\x04\x05\x06\x07"),
			[]byte("\x1c\x01\x02\x03\x04\x05\x06\x07"),
			[]byte("\x0c\x7f\xff\xff\xff\xff\xff\xff\xfe"),
			[]byte("\x1c\x80\x00\x00\x00\x00\x00\x00\x00"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeInt64(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendNullFloat64(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b sql.NullFloat64
		}{
			{
				a: sql.NullFloat64{Valid: false, Float64: 0},
				b: sql.NullFloat64{Valid: true, Float64: 0},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.Inf(-1)},
				b: sql.NullFloat64{Valid: true, Float64: -math.MaxFloat64},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: -math.MaxFloat64},
				b: sql.NullFloat64{Valid: true, Float64: math.Nextafter(-math.MaxFloat64, 0)},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: -2},
				b: sql.NullFloat64{Valid: true, Float64: -1},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: -1},
				b: sql.NullFloat64{Valid: true, Float64: 0},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: 0},
				b: sql.NullFloat64{Valid: true, Float64: 1},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: -1},
				b: sql.NullFloat64{Valid: true, Float64: 1},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: -2},
				b: sql.NullFloat64{Valid: true, Float64: 1},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: 1},
				b: sql.NullFloat64{Valid: true, Float64: 2},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.Nextafter(math.MaxFloat64, 0)},
				b: sql.NullFloat64{Valid: true, Float64: math.MaxFloat64},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.MaxFloat64},
				b: sql.NullFloat64{Valid: true, Float64: math.Inf(1)},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.Inf(-1)},
				b: sql.NullFloat64{Valid: true, Float64: math.Inf(1)},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.MaxFloat64},
				b: sql.NullFloat64{Valid: true, Float64: math.NaN()},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.Inf(-1)},
				b: sql.NullFloat64{Valid: true, Float64: math.NaN()},
			},
			{
				a: sql.NullFloat64{Valid: true, Float64: math.Inf(1)},
				b: sql.NullFloat64{Valid: true, Float64: math.NaN()},
			},
		}
		for i, tc := range testCases {
			a := keybytes.AppendNullFloat64([]byte(nil), tc.a)
			b := keybytes.AppendNullFloat64([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, tc.a=%v, tc.b=%v, a=0x%x, b=0x%x",
					i, got, want, tc.a, tc.b, a, b)
			}
		}
	})
}

func TestTakeNullFloat64(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []sql.NullFloat64{
			{Valid: false, Float64: 0},
			{Valid: true, Float64: math.NaN()},
			{Valid: true, Float64: math.Inf(1)},
			{Valid: true, Float64: -math.MaxFloat64},
			{Valid: true, Float64: -1.5},
			{Valid: true, Float64: -0.1},
			{Valid: true, Float64: -math.SmallestNonzeroFloat64},
			{Valid: true, Float64: math.Float64frombits(0x8000_0000_0000_0000)},
			{Valid: true, Float64: 0},
			{Valid: true, Float64: math.SmallestNonzeroFloat64},
			{Valid: true, Float64: 0.1},
			{Valid: true, Float64: 1.5},
			{Valid: true, Float64: math.MaxFloat64},
			{Valid: true, Float64: math.Inf(0)},
		}
		for i, input := range testCases {
			b := keybytes.AppendNullFloat64([]byte(nil), input)
			v, rest, err := keybytes.TakeNullFloat64(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := v, input; got.Valid && want.Valid {
				if math.Float64bits(got.Float64) != math.Float64bits(want.Float64) {
					t.Errorf("case %d: float64 unmatch: got=%v, want=%v", i, got, want)
				}
			} else if got.Valid != want.Valid {
				t.Errorf("case %d: valid unmatch: got=%v, want=%v", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
}

func TestAppendFloat64(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b float64
		}{
			{a: math.Inf(-1), b: -math.MaxFloat64},
			{a: -math.MaxFloat64, b: math.Nextafter(-math.MaxFloat64, 0)},
			{a: -2, b: -1},
			{a: -1, b: 0},
			{a: 0, b: 1},
			{a: -1, b: 1},
			{a: -2, b: 1},
			{a: 1, b: 2},
			{a: math.Nextafter(math.MaxFloat64, 0), b: math.MaxFloat64},
			{a: math.MaxFloat64, b: math.Inf(1)},
			{a: math.Inf(-1), b: math.Inf(1)},
			{a: math.MaxFloat64, b: math.NaN()},
			{a: math.Inf(1), b: math.NaN()},
			{a: math.Inf(-1), b: math.NaN()},
		}
		for i, tc := range testCases {
			a := keybytes.AppendFloat64([]byte(nil), tc.a)
			b := keybytes.AppendFloat64([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, tc.a=%v, tc.b=%v, a=0x%x, b=0x%x",
					i, got, want, tc.a, tc.b, a, b)
			}
		}
	})
}

func TestTakeFloat64(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []float64{
			math.NaN(),
			math.Inf(1),
			-math.MaxFloat64,
			-1.5,
			-0.1,
			-math.SmallestNonzeroFloat64,
			math.Float64frombits(0x8000_0000_0000_0000),
			0,
			math.SmallestNonzeroFloat64,
			0.1,
			1.5,
			math.MaxFloat64,
			math.Inf(0),
		}
		for i, input := range testCases {
			b := keybytes.AppendFloat64([]byte(nil), input)
			v, rest, err := keybytes.TakeFloat64(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := v, input; math.Float64bits(got) != math.Float64bits(want) {
				t.Errorf("case %d: string unmatch: got=%v, want=%v", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
}

func TestAppendNullBool(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b sql.NullBool
		}{
			{
				a: sql.NullBool{Valid: false, Bool: false},
				b: sql.NullBool{Valid: true, Bool: false},
			},
			{
				a: sql.NullBool{Valid: false, Bool: false},
				b: sql.NullBool{Valid: true, Bool: true},
			},
			{
				a: sql.NullBool{Valid: true, Bool: false},
				b: sql.NullBool{Valid: true, Bool: true},
			},
		}
		for i, tc := range testCases {
			a := keybytes.AppendNullBool([]byte(nil), tc.a)
			b := keybytes.AppendNullBool([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=%+v, b=%+v",
					i, got, want, a, b)
			}
		}
	})
}

func TestTakeNullBool(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []sql.NullBool{
			{Valid: false, Bool: false},
			{Valid: true, Bool: false},
			{Valid: true, Bool: true},
		}
		for i, input := range testCases {
			b := keybytes.AppendNullBool([]byte(nil), input)
			v, rest, err := keybytes.TakeNullBool(b)
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
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\x02"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeNullBool(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendBool(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		testCases := []struct {
			a, b bool
		}{
			{a: false, b: true},
		}
		for i, tc := range testCases {
			a := keybytes.AppendBool([]byte(nil), tc.a)
			b := keybytes.AppendBool([]byte(nil), tc.b)
			if got, want := bytes.Compare(a, b), -1; got != want {
				t.Errorf("case %d: compare result unmatch: got=%d, want=%d, a=%v, b=%v",
					i, got, want, a, b)
			}
		}
	})
}

func TestTakeBool(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		testCases := []bool{
			false,
			true,
		}
		for i, input := range testCases {
			b := keybytes.AppendBool([]byte(nil), input)
			v, rest, err := keybytes.TakeBool(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := v, input; got != want {
				t.Errorf("case %d: string unmatch: got=%v, want=%v", i, got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
	t.Run("invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte("\xff"),
		}
		for i, input := range testCases {
			_, _, err := keybytes.TakeBool(input)
			if err == nil {
				t.Errorf("case %d: got no error", i)
			}
		}
	})
}

func TestAppendCompositeKey(t *testing.T) {
	t.Run("roundtrip", func(t *testing.T) {
		type key struct {
			a sql.NullString
			b sql.NullInt32
			c sql.NullInt64
			d sql.NullFloat64
		}
		testCases := []key{
			{
				a: sql.NullString{Valid: false, String: ""},
				b: sql.NullInt32{Valid: false, Int32: 0},
				c: sql.NullInt64{Valid: false, Int64: 0},
				d: sql.NullFloat64{Valid: false, Float64: 0},
			},
			{
				a: sql.NullString{Valid: true, String: "foo"},
				b: sql.NullInt32{Valid: true, Int32: 1234},
				c: sql.NullInt64{Valid: true, Int64: 5678},
				d: sql.NullFloat64{Valid: true, Float64: 2.3},
			},
		}
		for i, k := range testCases {
			b := keybytes.AppendNullString([]byte(nil), k.a)
			b = keybytes.AppendNullInt32(b, k.b)
			b = keybytes.AppendNullInt64(b, k.c)
			b = keybytes.AppendNullFloat64(b, k.d)

			var k2 key
			var rest []byte
			var err error
			k2.a, rest, err = keybytes.TakeNullString(b)
			if err != nil {
				t.Errorf("case %d: got error: %s", i, err)
			}
			if got, want := k2.a, k.a; got.Valid && want.Valid {
				if got.String != want.String {
					t.Errorf("case %d .a: string unmatch: got=%v, want=%v", i, got, want)
				}
			} else if got.Valid != want.Valid {
				t.Errorf("case %d .a: valid unmatch: got=%v, want=%v", i, got, want)
			}

			k2.b, rest, err = keybytes.TakeNullInt32(rest)
			if err != nil {
				t.Errorf("case %d .b: got error: %s", i, err)
			}
			if got, want := k2.b, k.b; got.Valid && want.Valid {
				if got.Int32 != want.Int32 {
					t.Errorf("case %d .b: int32 unmatch: got=%v, want=%v", i, got, want)
				}
			} else if got.Valid != want.Valid {
				t.Errorf("case %d .b: valid unmatch: got=%v, want=%v", i, got, want)
			}

			k2.c, rest, err = keybytes.TakeNullInt64(rest)
			if err != nil {
				t.Errorf("case %d .c: got error: %s", i, err)
			}
			if got, want := k2.c, k.c; got.Valid && want.Valid {
				if got.Int64 != want.Int64 {
					t.Errorf("case %d .c: int64 unmatch: got=%v, want=%v", i, got, want)
				}
			} else if got.Valid != want.Valid {
				t.Errorf("case %d .c: valid unmatch: got=%v, want=%v", i, got, want)
			}

			k2.d, rest, err = keybytes.TakeNullFloat64(rest)
			if err != nil {
				t.Errorf("case %d .d: got error: %s", i, err)
			}
			if got, want := k2.d, k.d; got.Valid && want.Valid {
				if math.Float64bits(got.Float64) != math.Float64bits(want.Float64) {
					t.Errorf("case %d .d: float64 unmatch: got=%v, want=%v", i, got, want)
				}
			} else if got.Valid != want.Valid {
				t.Errorf("case %d .d: valid unmatch: got=%v, want=%v", i, got, want)
			}

			if got, want := len(rest), 0; got != want {
				t.Errorf("case %d: rest length unmatch: got=%d, want=%d", i, got, want)
			}
		}
	})
}
