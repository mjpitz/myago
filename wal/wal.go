package wal

import (
	"encoding/binary"
)

// RecordLength returns the computed length for the underlying record.
func RecordLength(data []byte) int {
	buffer := make([]byte, 10)
	n := binary.PutUvarint(buffer, uint64(len(data)))
	return n + len(data) + 4
}
