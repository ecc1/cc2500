package cc1100

import (
	"bytes"
	"testing"
)

func TestEncoding(t *testing.T) {
	cases := []struct {
		src []byte
		dst []byte
	}{
		{[]byte{
			0xA7, 0x12, 0x34, 0x56, 0x8D, 0x00, 0xA6,
		}, []byte{
			0xA9, 0x6C, 0x72, 0x8F, 0x49, 0x66, 0x68, 0xD5,
			0x55, 0xAA, 0x65,
		}},
		{[]byte{
			0xA7, 0x12, 0x34, 0x56, 0x06, 0x00, 0x03,
		}, []byte{
			0xA9, 0x6C, 0x72, 0x8F, 0x49, 0x66, 0x56, 0x65,
			0x55, 0x56, 0x35,
		}},
		{[]byte{
			0xA7, 0x12, 0x34, 0x56, 0x5D, 0x02, 0x01, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C,
		}, []byte{
			0xA9, 0x6C, 0x72, 0x8F, 0x49, 0x66, 0x94, 0xD5,
			0x72, 0x57, 0x15, 0x71, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x56, 0xC5,
		}},
		{[]byte{
			0xA7, 0x12, 0x34, 0x56, 0x8D, 0x09, 0x03, 0x37,
			0x32, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC2,
		}, []byte{
			0xA9, 0x6C, 0x72, 0x8F, 0x49, 0x66, 0x68, 0xD5,
			0x59, 0x56, 0x38, 0xD6, 0x8F, 0x28, 0xF2, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
			0x55, 0xB3, 0x25,
		}},
	}

	for _, c := range cases {
		result := Encode4b6b(c.src)
		if !bytes.Equal(result, c.dst) {
			t.Errorf("Encode4b6b(%X) == %X, want %X", c.src, result, c.dst)
		}
		result, err := Decode6b4b(c.dst)
		if err != nil {
			t.Errorf("Decode6b4b(%X) == %v, want %X", c.dst, err, c.src)
		} else if !bytes.Equal(result, c.src) {
			t.Errorf("Decode6b4b(%X) == %X, want %X", c.dst, result, c.src)
		}
	}
}