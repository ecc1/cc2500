package cc2500

//go:generate crcgen -size 8 -poly 0x2F

// CRC8 computes the 8-bit CRC of the given data.
func CRC8(msg []byte) byte {
	res := byte(0)
	for _, b := range msg {
		res = crc8Table[res^b]
	}
	return res
}
