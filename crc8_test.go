package cc2500

import (
	"strconv"
	"strings"
	"testing"
)

func TestCRC8(t *testing.T) {
	cases := []struct {
		msg []byte
		sum byte
	}{
		{parseBytes("19 FD A5 F9 D8 00"), 0xD2},
		{parseBytes("1B FD 4E F9 D7 00"), 0x2F},
		{parseBytes("45 FD CC F9 D8 00"), 0x28},
		{parseBytes("4E 7D AB B9 D7 00"), 0xD4},
		{parseBytes("54 FD 5D 79 D8 00"), 0x60},
		{parseBytes("5D BD 7C 39 D7 00"), 0x33},
		{parseBytes("68 7D 08 B9 D8 00"), 0x13},
		{parseBytes("8A DD 29 59 D7 00"), 0x0C},
		{parseBytes("A2 DD D1 99 D7 00"), 0xD0},
		{parseBytes("A2 FD 1F 79 D7 00"), 0x5C},
		{parseBytes("BD FD 03 F9 D7 00"), 0xF8},
		{parseBytes("EB 7D C6 79 D7 00"), 0x2F},
		{parseBytes("EF 3D 76 D9 D7 00"), 0x3C},
		{parseBytes("FC 1D ED 19 D7 00"), 0xCE},
	}
	for _, c := range cases {
		sum := CRC8(c.msg)
		if sum != c.sum {
			t.Errorf("CRC8(% X) == %02X, want %02X", c.msg, sum, c.sum)
		}
	}
}

func parseBytes(hex string) []byte {
	fields := strings.Fields(hex)
	data := make([]byte, len(fields))
	for i, s := range fields {
		b, err := strconv.ParseUint(string(s), 16, 8)
		if err != nil {
			panic(err)
		}
		data[i] = byte(b)
	}
	return data
}
