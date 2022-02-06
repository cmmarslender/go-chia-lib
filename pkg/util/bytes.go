package util

import (
	"fmt"
)

// ShiftNBytes returns the specified number of bytes from the start of the provided []byte
// and removes them from the original byte slice
// First returned value is the requested number of bytes from the beginning of the original byte slice
// Second returned value is the new original byte slice with the requested number of bytes removed from the front of it
func ShiftNBytes(numBytes uint, bytes []byte) ([]byte, []byte, error) {
	if uint(len(bytes)) < numBytes {
		return nil, bytes, fmt.Errorf("requested more bytes than available")
	}

	requestedBytes := bytes[:numBytes]
	bytes = bytes[numBytes:]

	return requestedBytes, bytes, nil
}
