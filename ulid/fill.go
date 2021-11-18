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
