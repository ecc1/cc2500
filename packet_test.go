package cc2500

import (
	"reflect"
	"testing"
	"time"
)

var (
	p1 = parseBytes("FF FF FF FF AE D1 63 00 3F 03 76 59 8D 12 49 D5 00 E4")
	p2 = parseBytes("FF FF FF FF AE D1 63 00 3F 03 6E 39 4D 89 C9 D5 00 DE")
	p3 = parseBytes("FF FF FF FF F2 58 68 00 3F 03 D6 FC 1D ED 19 D7 00 CE")
	p4 = parseBytes("FF FF FF FF F2 58 68 00 3F 03 DA 03 ED 5A 19 D7 00 F4")
	p5 = parseBytes("FF FF FF FF B0 CD 61 00 3F 03 AF 11 39 B3 7E D9 00 4E")
)

func TestUnmarshalPacket(t *testing.T) {
	cases := []struct {
		data   []byte
		packet Packet
	}{
		{p1, Packet{
			Data:          p1,
			TransmitterID: "67LDE",
			Raw:           144192,
			Filtered:      149760,
			Battery:       213,
		}},
		{p2, Packet{
			Data:          p2,
			TransmitterID: "67LDE",
			Raw:           152448,
			Filtered:      160288,
			Battery:       213,
		}},
		{p3, Packet{
			Data:          p3,
			TransmitterID: "6GN7J",
			Raw:           198624,
			Filtered:      202464,
			Battery:       215,
		}},
		{p4, Packet{
			Data:          p4,
			TransmitterID: "6GN7J",
			Raw:           194560,
			Filtered:      199488,
			Battery:       215,
		}},
		{p5, Packet{
			Data:          p5,
			TransmitterID: "63KDG",
			Raw:           116864,
			Filtered:      126160,
			Battery:       217,
		}},
	}
	t0 := time.Time{}
	for _, c := range cases {
		packet := unmarshalPacket(t0, 0, c.data, 0)
		if !reflect.DeepEqual(packet, &c.packet) {
			t.Errorf("makePacket(% X) == %+v, want %+v", c.data, *packet, c.packet)
		}
	}
}
