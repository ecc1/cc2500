package cc2500

import (
	"bytes"
	"testing"
)

func TestFrequency(t *testing.T) {
	cases := []struct {
		f       uint32
		b       []byte
		fApprox uint32 // 0 => equal to f
	}{
		{2418000000, []byte{0x5D, 0x00, 0x00}, 0},
		{2434250000, []byte{0x5D, 0xA0, 0x00}, 0},
		// some that can't be represented exactly:
		{2400000000, []byte{0x5C, 0x4E, 0xC5}, 2400000030},
		{2403470000, []byte{0x5C, 0x70, 0xEF}, 2403469818},
		{2425000000, []byte{0x5D, 0x44, 0xEC}, 2424999877},
	}

	for _, c := range cases {
		b := frequencyToRegisters(c.f)
		if !bytes.Equal(b, c.b) {
			t.Errorf("frequencyToRegisters(%d) == % X, want % X", c.f, b, c.b)
		}
		f := registersToFrequency(c.b)
		if c.fApprox == 0 {
			if f != c.f {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.f)
			}
		} else {
			if f != c.fApprox {
				t.Errorf("registersToFrequency(% X) == %d, want %d", c.b, f, c.fApprox)
			}
		}
	}
}
