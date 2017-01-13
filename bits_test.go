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

func TestUnmarshalUint16(t *testing.T) {
	cases := []struct {
		val uint16
		rep []byte
	}{
		{0x1234, []byte{0x34, 0x12}},
		{0, []byte{0, 0}},
		{math.MaxUint16, []byte{0xFF, 0xFF}},
	}
	for _, c := range cases {
		val := unmarshalUint16(c.rep)
		if val != c.val {
			t.Errorf("unmarshalUint16(% X) == %04X, want %04X", c.rep, val, c.val)
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

func TestUnmarshalTransmitterID(t *testing.T) {
	cases := []struct {
		rep []byte
		id  string
	}{
		{[]byte{0xCA, 0xC1, 0x61, 0x00}, "63GEA"},
		{[]byte{0xCA, 0x4C, 0x62, 0x00}, "64K6A"},
		{[]byte{0xAE, 0xD1, 0x63, 0x00}, "67LDE"},
	}
	for _, c := range cases {
		id := unmarshalTransmitterID(c.rep)
		if id != c.id {
			t.Errorf("unmarshalTransmitterID(% X) == %q, want %q", c.rep, id, c.id)
		}
	}
}

func TestUnmarshalReadings(t *testing.T) {
	cases := []struct {
		v []byte
		f uint32
		u uint32
	}{
		{[]byte{0xB8, 0xCD, 0x2E, 0x29}, 156576, 167552},
		{[]byte{0x39, 0x4D, 0x89, 0xC9}, 152448, 160288},
		{[]byte{0x59, 0x8D, 0x12, 0x49}, 144192, 149760},
	}
	for _, c := range cases {
		f := unmarshalReading(c.v[0:2])
		if f != c.f {
			t.Errorf("unmarshalReading(% X) == %d, want %d", c.v[0:2], f, c.f)
		}
		u := 2 * unmarshalReading(c.v[2:4])
		if u != c.u {
			t.Errorf("unmarshalReading(% X) == %d, want %d", c.v[2:4], u, c.u)
		}
	}
}
