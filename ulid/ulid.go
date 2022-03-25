// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ulid

import (
	"database/sql/driver"
	"encoding/base32"
	"encoding/binary"
	"fmt"
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
	// PayloadLength varies by number of bits.
)

var crockfordBase32 = base32.
	NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").
	WithPadding(base32.NoPadding)

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

// ULID is a variable-length, generalized, unique lexographical identifier.
type ULID []byte

// Skew returns the current skew used to correct massive time skews.
func (ulid ULID) Skew() byte {
	return ulid[SkewOffset]
}

// Timestamp returns the timestamp portion of the identifier.
func (ulid ULID) Timestamp() time.Time {
	millis := binary.BigEndian.Uint64(append(make([]byte, 2), ulid[UnixOffset:UnixOffset+UnixLength]...))

	return time.UnixMilli(int64(millis))
}

// Payload returns a copy of the payload bytes.
func (ulid ULID) Payload() []byte {
	return append([]byte{}, ulid[PayloadOffset:]...)
}

// Bytes returns a copy of the underlying byte array backing the ulid.
func (ulid ULID) Bytes() []byte {
	return append([]byte{}, ulid...)
}

// String returns a string representation of the payload. It's encoded using a crockford base32 encoding.
func (ulid ULID) String() string {
	return crockfordBase32.EncodeToString(ulid)
}

// Value serializes the ULID so that it can be stored in SQL databases.
func (ulid ULID) Value() (driver.Value, error) {
	return ulid.String(), nil
}

// Scan attempts to parse the provided src value into the ULID.
func (ulid *ULID) Scan(src interface{}) (err error) {
	if src == nil {
		*ulid = ULID{}
		return nil
	}

	if ulid == nil {
		return fmt.Errorf("destination pointer is nil")
	}

	var val ULID
	switch v := src.(type) {
	case *string:
		val, err = Parse(*v)
	case string:
		val, err = Parse(v)
	default:
		err = fmt.Errorf("unsupport source type: %s", v)
	}

	if err != nil {
		return err
	}

	*ulid = val

	return nil
}
