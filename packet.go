package cc2500

import (
	"math/bits"
	"time"
)

// Packet represents a Dexcom G4 packet.
type Packet struct {
	Timestamp     time.Time
	Channel       int
	Data          []byte
	TransmitterID string
	Raw           uint32
	Filtered      uint32
	Battery       uint8
	RSSI          int
}

// Wire format of Dexcom G4 packet:
//	0..3: destination address (always FF FF FF FF = broadcast)
//	4..7: transmitter ID
//	8: port? (always 3F)
//	9: device info? (always 03)
//	10: sequence number
//	11..12: raw reading
//	13..14: filtered reading
//	15: battery level
//	16: unknown
//	17: checksum
const packetLength = 18

func unmarshalPacket(t time.Time, n int, data []byte, rssi int) *Packet {
	return &Packet{
		Timestamp:     t,
		Channel:       n,
		Data:          data,
		TransmitterID: unmarshalTransmitterID(data[4:8]),
		Raw:           unmarshalReading(data[11:13]),
		Filtered:      2 * unmarshalReading(data[13:15]),
		Battery:       data[15],
		RSSI:          rssi,
	}
}

// Unmarshal a little-endian uint16.
func unmarshalUint16(v []byte) uint16 {
	return uint16(v[0]) | uint16(v[1])<<8
}

// Unmarshal a little-endian uint32.
func unmarshalUint32(v []byte) uint32 {
	return uint32(unmarshalUint16(v[0:2])) | uint32(unmarshalUint16(v[2:4]))<<16
}

var transmitterIDChar = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P',
	'Q', 'R', 'S', 'T', 'U', 'W', 'X', 'Y',
}

// The 5-character transmitter ID is encoded as a sequence of 5-bit symbols
// in a left-padded, little-endian 32-bit integer.
func unmarshalTransmitterID(v []byte) string {
	u := unmarshalUint32(v)
	id := make([]byte, 5)
	for i := 0; i < 5; i++ {
		n := byte(u>>uint(20-5*i)) & 0x1F
		id[i] = transmitterIDChar[n]
	}
	return string(id)
}

// Unmarshal a 16-bit float (13-bit mantissa, 3-bit exponent) as a uint32.
func unmarshalReading(v []byte) uint32 {
	u0, u1 := bits.Reverse8(v[0]), bits.Reverse8(v[1])
	u := uint16(u0) | uint16(u1)<<8
	return uint32(u&0x1FFF) << (u >> 13)
}
