package cc2500

import (
	"math"
	"testing"
)

func TestReverseBits(t *testing.T) {
	cases := []struct {
		b   byte
		rev byte
	}{
		{0x00, 0x00},
		{0x01, 0x80},
		{0x0F, 0xF0},
		{0x05, 0xA0},
		{0x55, 0xAA},
		{0xFF, 0xFF},
	}
	for _, c := range cases {
		rev := reverseBits[c.b]
		if rev != c.rev {
			t.Errorf("reverseBits[%08b] == %08b, want %08b", c.b, rev, c.rev)
		}
	}
	// Test that bit-reversal is self-inverse.
	for i := 0; i < 256; i++ {
		b := byte(i)
		j := reverseBits[reverseBits[b]]
		if j != b {
			t.Errorf("reverseBits[reverseBits[%08b]] == %08b, want %08b", b, j, b)
		}
	}
}

func TestUnarshalUint32(t *testing.T) {
	cases := []struct {
		val uint32
		rep []byte
	}{
		{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		{0, []byte{0, 0, 0, 0}},
		{math.MaxUint32, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := unmarshalUint32(c.rep)
		if val != c.val {
			t.Errorf("unmarshalUint32(% X) == %08X, want %08X", c.rep, val, c.val)
		}
	}
}

func TestTransmitterID(t *testing.T) {
	cases := []struct {
		rep []byte
		id  string
	}{
		{[]byte{0xCA, 0xC1, 0x61, 0x00}, "63GEA"},
		{[]byte{0xCA, 0x4C, 0x62, 0x00}, "64K6A"},
		{[]byte{0xAE, 0xD1, 0x63, 0x00}, "67LDE"},
	}
	for _, c := range cases {
		id := transmitterID(c.rep)
		if id != c.id {
			t.Errorf("transmitterID(% X) == %q, want %q", c.rep, id, c.id)
		}
	}
}
