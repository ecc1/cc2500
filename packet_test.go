package cc2500

import (
	"testing"
	"time"
)

func TestMakePacket(t *testing.T) {
	cases := []struct {
		data   []byte
		packet Packet
	}{
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xAE, 0xD1, 0x63, 0x00, 0x3F, 0x03, 0x76, 0x59, 0x8D, 0x12, 0x49, 0xD5, 0x00, 0xE4},
			Packet{
				TransmitterID: "67LDE",
				Raw:           144192,
				Filtered:      149760,
				Battery:       213,
			}},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xAE, 0xD1, 0x63, 0x00, 0x3F, 0x03, 0x6E, 0x39, 0x4D, 0x89, 0xC9, 0xD5, 0x00, 0xDE},
			Packet{
				TransmitterID: "67LDE",
				Raw:           152448,
				Filtered:      160288,
				Battery:       213,
			}},
	}
	t0 := time.Time{}
	for _, c := range cases {
		packet := makePacket(t0, 0, c.data, 0)
		if *packet != c.packet {
			t.Errorf("makePacket(% X) == %+v, want %+v", c.data, packet, c.packet)
		}
	}
}
