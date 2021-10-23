package ulid256

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// Fill defines an arbitrary way to fill a slice of data with information. This information can be structureless or
// semi-structured. In the end, it should not be used for any control logic.
type Fill func(data []byte) (int, error)

// RandomFill fills the data array with random data from a cryptographically secure source.
func RandomFill() Fill {
	return func(data []byte) (int, error) {
		return io.ReadFull(rand.Reader, data)
	}
}

// ServerIDFill prefixes the data payload with a serverID. It delegates filling the remaining portion to the provided
// fill.
func ServerIDFill(serverId uint16, fill Fill) Fill {
	return func(data []byte) (int, error) {
		binary.BigEndian.PutUint16(data[:2], serverId)
		n, err := fill(data[2:])
		return n + 2, err
	}
}
