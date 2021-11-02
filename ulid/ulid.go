package ulid

import (
	"encoding/binary"
	"time"
)

const (
	// SkewOffset is the starting byte position for the skew data.
	SkewOffset = 0
	// SkewLength is the number of bytes representing the skew.
	SkewLength = 1

	// UnixOffset is the starting byte position for the unix timestamp data.
	UnixOffset = 1
	// UnixLength is the number of bytes representing the unix timestamp data.
	UnixLength = 6

	// PayloadOffset is the starting byte position for the payload data.
	PayloadOffset = 7
	// PayloadLength varies by number of bits
)

// ULID is a generalized, unique lexographical identifier. The format is as follows:
//
// `[ skew ][ sec ][ payload ]`
//
// - `skew` - 1 byte used to handle major clock skews (reserved, unused)
// - `sec` - 6 bytes of a unix timestamp (should give us until the year 10k or so)
// - `payload` - N bytes for the payload
//
type ULID []byte

// Skew returns the current skew used to correct massive time skews.
func (ulid ULID) Skew() byte {
	return ulid[SkewOffset]
}

// Timestamp returns the timestamp portion of the identifier.
func (ulid ULID) Timestamp() time.Time {
	seconds := binary.BigEndian.Uint64(append(make([]byte, 2), ulid[UnixOffset:UnixOffset+UnixLength]...))
	return time.Unix(int64(seconds), 0)
}

// Payload returns a copy of the payload bytes.
func (ulid ULID) Payload() []byte {
	return append([]byte{}, ulid[PayloadOffset:]...)
}
