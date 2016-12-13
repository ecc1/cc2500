package cc2500

import (
	"time"
)

type Packet struct {
	Timestamp     time.Time
	Channel       int
	TransmitterID string
	Raw           uint32
	Filtered      uint32
	Battery       uint8
	RSSI          int
}

func makePacket(t time.Time, n int, data []byte, rssi int) *Packet {
	return &Packet{
		Timestamp:     t,
		Channel:       n,
		TransmitterID: transmitterID(data[4:8]),
		Raw:           unmarshalReading(data[11:13]),
		Filtered:      2 * unmarshalReading(data[13:15]),
		Battery:       data[15],
		RSSI:          rssi,
	}
}

var transmitterIdChar = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P',
	'Q', 'R', 'S', 'T', 'U', 'W', 'X', 'Y',
}

// Unmarshal a little-endian uint16.
func unmarshalUint16(v []byte) uint16 {
	return uint16(v[0]) | uint16(v[1])<<8
}

// Unmarshal a little-endian uint32.
func unmarshalUint32(v []byte) uint32 {
	return uint32(unmarshalUint16(v[0:2])) | uint32(unmarshalUint16(v[2:4]))<<16
}

// The 5-character transmitter ID is encoded as a sequence of 5-bit symbols
// in a left-padded, little-endian 32-bit integer.
func transmitterID(v []byte) string {
	u := unmarshalUint32(v)
	id := make([]byte, 5)
	for i := 0; i < 5; i++ {
		n := byte(u>>uint(20-5*i)) & 0x1F
		id[i] = transmitterIdChar[n]
	}
	return string(id)
}

// Unmarshal a 16-bit float (13-bit mantissa, 3-bit exponent) as a uint32.
func unmarshalReading(v []byte) uint32 {
	u0, u1 := reverseBits[v[0]], reverseBits[v[1]]
	u := uint16(u0) | uint16(u1)<<8
	return uint32(u&0x1FFF) << (u >> 13)
}
