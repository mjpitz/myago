package ulid

import (
	"encoding/base32"
	"encoding/binary"
	"strings"
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

var (
	crockfordBase32 = base32.
		NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").
		WithPadding(base32.NoPadding)
)

// Parse accepts a ULID string and attempts to extract a ULID from the provided string.
func Parse(ulid string) (ULID, error) {
	ulid = strings.ToUpper(ulid)

	bytes, err := crockfordBase32.DecodeString(ulid)
	switch {
	case err != nil:
		return nil, err
	case len(bytes) < 8:
		return nil, ErrNotEnoughBits
	}

	parsed := make(ULID, len(bytes))
	copy(parsed[:], bytes)

	return parsed, nil
}

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

// String returns a string representation of the payload. It's encoded using a crockford base32 encoding.
func (ulid ULID) String() string {
	return crockfordBase32.EncodeToString(ulid)
}
