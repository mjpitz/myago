package ulid

import (
	"crypto/rand"
	"io"
)

// Fill provides an abstraction for filling the data payload of a ULID.
type Fill func(ulid ULID, data []byte) (int, error)

// RandomFill is a fill that populates the data payload with random data.
func RandomFill(_ ULID, data []byte) (int, error) {
	random := make([]byte, len(data))
	n, err := io.ReadFull(rand.Reader, random)
	if err != nil {
		return n, err
	}

	return copy(data, random[:n]), nil
}
