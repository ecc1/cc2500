package cc2500

var (
	transmitterIdChar = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
		'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P',
		'Q', 'R', 'S', 'T', 'U', 'W', 'X', 'Y',
	}
)

// Unmarshal a little-endian uint32.
func unmarshalUint32(v []byte) uint32 {
	n := uint32(0)
	for i := 0; i < 4; i++ {
		n |= uint32(v[i]) << uint(8*i)
	}
	return n
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
