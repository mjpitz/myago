package ulid256

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"time"
)

// Empty is an empty ULID.
var Empty = ULID{}

// Parse converts a string based ULID to it's binary representation.
func Parse(ulid string) (ULID, error) {
	val, err := base64.RawURLEncoding.DecodeString(ulid)
	if err != nil {
		return Empty, err
	}

	out := ULID{}
	copy(out[:], val)

	if err := out.Validate(); err != nil {
		return Empty, err
	}

	return out, nil
}

// ULID is a unique, lexigraphic identifier. Unlike the the canonical ULID implementation, this version is 256 bits,
// holds a version identifier, a programmable the payload, and a CRC32 checksum. The binary format is as follows:
//
// `[ skew - 2 bytes ][ sec - 6 bytes ][ nsec - 3 bytes ][ version - byte ][ payload - 16 bytes ][ checksum - 4 bytes ]`
//
type ULID [32]byte

// Skew returns the clock skew factor.
func (u ULID) Skew() uint16 {
	return binary.BigEndian.Uint16(u[0:2])
}

// Time returns the timestamp associated with this ULID.
func (u ULID) Time() time.Time {
	sec := binary.BigEndian.Uint64(append(make([]byte, 2), u[2:8]...))
	nsec := binary.BigEndian.Uint32(append(make([]byte, 1), u[8:11]...))
	return time.Unix(int64(sec), int64(nsec))
}

// Version returns the version of the payload data.
func (u ULID) Version() byte {
	return u[11]
}

// Payload returns the data portion for the ULID.
func (u ULID) Payload() []byte {
	return append([]byte{}, u[12:28]...)
}

// Checksum returns and IEEE CRC32 checksum portion of the ULID.
func (u ULID) Checksum() uint32 {
	return binary.BigEndian.Uint32(u[28:32])
}

// Validate returns an error if the checksums do not match.
func (u ULID) Validate() error {
	checksum := crc32.NewIEEE()
	_, err := checksum.Write(u[:28])
	if err != nil {
		return err
	}

	if checksum.Sum32() != u.Checksum() {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}

// String serializes the ULID to a string representation.
func (u ULID) String() string {
	return base64.RawURLEncoding.EncodeToString(u[:])
}
