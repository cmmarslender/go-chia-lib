package util

import (
	"encoding/binary"
)

// Uint8ToBytes Converts uint8 to []byte
// Kind of pointless, since byte is uint8, but here for consistency with the other methods
func Uint8ToBytes(num uint8) []byte {
	return []byte{num}
}

// Uint16ToBytes Converts uint16 to []byte
func Uint16ToBytes(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)

	return b
}

// Uint32ToBytes Converts uint32 to []byte
func Uint32ToBytes(num uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)

	return b
}
