package keybytes_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/hnakamur/keybytes"
)

func TestAppendRaw(t *testing.T) {
	t.Run("emptyDst", func(t *testing.T) {
		key := []byte("foo")
		var dst []byte
		dst = keybytes.AppendRaw(dst, key)

		key2, rest := keybytes.TakeRaw(dst, len(key))
		if got, want := key2, key; !bytes.Equal(got, want) {
			t.Errorf("decoded key unmatch, got=%s, want=%s", string(got), string(want))
		}
		if got, want := len(rest), 0; got != want {
			t.Errorf("rest length unmatch, got=%v, want=%v", got, want)
		}
	})
	t.Run("shortCapDst", func(t *testing.T) {
		key := []byte("foo")
		dst := make([]byte, 0, 2)
		dst = keybytes.AppendRaw(dst, key)

		key2, rest := keybytes.TakeRaw(dst, len(key))
		if got, want := key2, key; !bytes.Equal(got, want) {
			t.Errorf("decoded key unmatch, got=%s, want=%s", string(got), string(want))
		}
		if got, want := len(rest), 0; got != want {
			t.Errorf("rest length unmatch, got=%v, want=%v", got, want)
		}
	})
	t.Run("longCapDst", func(t *testing.T) {
		key := []byte("foo")
		dst := make([]byte, 0, 3)
		dst = keybytes.AppendRaw(dst, key)

		key2, rest := keybytes.TakeRaw(dst, len(key))
		if got, want := key2, key; !bytes.Equal(got, want) {
			t.Errorf("decoded key unmatch, got=%s, want=%s", string(got), string(want))
		}
		if got, want := len(rest), 0; got != want {
			t.Errorf("rest length unmatch, got=%v, want=%v", got, want)
		}
	})
}

func TestAppendByte(t *testing.T) {
	t.Run("encodeDecode", func(t *testing.T) {
		for k := byte(0); ; k++ {
			dst := keybytes.AppendByte([]byte{}, k)
			k2, rest := keybytes.TakeByte(dst)
			if got, want := k2, k; got != want {
				t.Errorf("decoded key unmatch, got=%v, want=%v", got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("rest length unmatch, got=%v, want=%v", got, want)
			}

			if k == math.MaxUint8 {
				break
			}
		}
	})
	t.Run("keepsOrder", func(t *testing.T) {
		for k1 := byte(0); k1 < math.MaxUint8; k1++ {
			k2 := k1 + 1
			kb1 := keybytes.AppendByte([]byte{}, k1)
			kb2 := keybytes.AppendByte([]byte{}, k2)
			if got, want := bytes.Compare(kb1, kb2), -1; got != want {
				t.Errorf("unexpected compare result, got=%d, want=%d, k1=%v, k2=%v, kb1=%v, kb2=%v",
					got, want, k1, k2, kb1, kb2)
			}
		}
	})
}

func TestAppendUint16(t *testing.T) {
	t.Run("encodeDecode", func(t *testing.T) {
		for k := uint16(0); ; k++ {
			dst := keybytes.AppendUint16([]byte{}, k)
			k2, rest := keybytes.TakeUint16(dst)
			if got, want := k2, k; got != want {
				t.Errorf("decoded key unmatch, got=%v, want=%v", got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("rest length unmatch, got=%v, want=%v", got, want)
			}

			if k == math.MaxUint16 {
				break
			}
		}
	})
	t.Run("keepsOrder", func(t *testing.T) {
		for k1 := uint16(0); k1 < math.MaxUint16; k1++ {
			k2 := k1 + 1
			kb1 := keybytes.AppendUint16([]byte{}, k1)
			kb2 := keybytes.AppendUint16([]byte{}, k2)
			if got, want := bytes.Compare(kb1, kb2), -1; got != want {
				t.Errorf("unexpected compare result, got=%d, want=%d, k1=%v, k2=%v, kb1=%v, kb2=%v",
					got, want, k1, k2, kb1, kb2)
			}
		}
	})
}

func TestAppendUint32(t *testing.T) {
	t.Run("encodeDecode", func(t *testing.T) {
		testCases := []uint32{0, 1, math.MaxUint32 - 1, math.MaxUint32}
		for _, k := range testCases {
			dst := keybytes.AppendUint32([]byte{}, k)
			k2, rest := keybytes.TakeUint32(dst)
			if got, want := k2, k; got != want {
				t.Errorf("decoded key unmatch, got=%v, want=%v", got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("rest length unmatch, got=%v, want=%v, k=%d", got, want, k)
			}
		}
	})
	t.Run("keepsOrder", func(t *testing.T) {
		testCases := []uint32{0, 1, math.MaxUint32 - 2, math.MaxUint32 - 1}
		for _, k1 := range testCases {
			k2 := k1 + 1
			kb1 := keybytes.AppendUint32([]byte{}, k1)
			kb2 := keybytes.AppendUint32([]byte{}, k2)
			if got, want := bytes.Compare(kb1, kb2), -1; got != want {
				t.Errorf("unexpected compare result, got=%d, want=%d, k1=%v, k2=%v, kb1=%v, kb2=%v",
					got, want, k1, k2, kb1, kb2)
			}
		}
	})
}

func TestAppendUint64(t *testing.T) {
	t.Run("encodeDecode", func(t *testing.T) {
		testCases := []uint64{0, 1, math.MaxUint64 - 1, math.MaxUint64}
		for _, k := range testCases {
			dst := keybytes.AppendUint64([]byte{}, k)
			k2, rest := keybytes.TakeUint64(dst)
			if got, want := k2, k; got != want {
				t.Errorf("decoded key unmatch, got=%v, want=%v", got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("rest length unmatch, got=%v, want=%v, k=%d", got, want, k)
			}
		}
	})
	t.Run("keepsOrder", func(t *testing.T) {
		testCases := []uint64{0, 1, math.MaxUint64 - 2, math.MaxUint64 - 1}
		for _, k1 := range testCases {
			k2 := k1 + 1
			kb1 := keybytes.AppendUint64([]byte{}, k1)
			kb2 := keybytes.AppendUint64([]byte{}, k2)
			if got, want := bytes.Compare(kb1, kb2), -1; got != want {
				t.Errorf("unexpected compare result, got=%d, want=%d, k1=%v, k2=%v, kb1=%v, kb2=%v",
					got, want, k1, k2, kb1, kb2)
			}
		}
	})
}

func TestAppendStringNul(t *testing.T) {
	t.Run("encodeDecode", func(t *testing.T) {
		testCases := []string{"", "f", "foo"}
		for _, k := range testCases {
			dst := keybytes.AppendStringNul([]byte{}, k)
			k2, rest := keybytes.TakeStringNul(dst)
			if got, want := k2, k; got != want {
				t.Errorf("decoded key unmatch, got=%v, want=%v", got, want)
			}
			if got, want := len(rest), 0; got != want {
				t.Errorf("rest length unmatch, got=%v, want=%v, k=%s", got, want, k)
			}
		}
	})
	t.Run("keepsOrder", func(t *testing.T) {
		testCases := []struct {
			k1, k2 string
		}{
			{k1: "", k2: "b"},
			{k1: "bar", k2: "bars"},
			{k1: "bar", k2: "baz"},
			{k1: "bar", k2: "bb"},
			{k1: "bar", k2: "c"},
		}
		for _, tc := range testCases {
			k1 := tc.k1
			k2 := tc.k2
			kb1 := keybytes.AppendStringNul([]byte{}, k1)
			kb2 := keybytes.AppendStringNul([]byte{}, k2)
			if got, want := bytes.Compare(kb1, kb2), -1; got != want {
				t.Errorf("unexpected compare result, got=%d, want=%d, k1=%v, k2=%v, kb1=%v, kb2=%v",
					got, want, k1, k2, kb1, kb2)
			}
		}
	})
}
