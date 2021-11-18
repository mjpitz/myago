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
	"context"
	"encoding/binary"

	"github.com/mjpitz/myago/clocks"
)

// NewGenerator constructs a generator using the provided skew and fill.
func NewGenerator(skew byte, fill Fill) *Generator {
	return &Generator{
		skew: skew,
		fill: fill,
	}
}

// Generator is the base interface defines how to generate ULIDs of varying length.
type Generator struct {
	skew byte
	fill Fill
}

// Generate produces a new ULID unless an error occurs. A clock can be provided on the context to override how the
// Generator obtains a timestamp.
func (g *Generator) Generate(ctx context.Context, bits int) (ULID, error) {
	clock := clocks.Extract(ctx)

	switch {
	case bits < 64:
		return nil, ErrNotEnoughBits
	case bits%8 > 0:
		return nil, ErrInvalidBitCount
	}

	unixmillis := uint64(clock.Now().UnixMilli())
	millis := make([]byte, 8)
	binary.BigEndian.PutUint64(millis, unixmillis)

	ulid := make(ULID, bits/8)
	ulid[SkewOffset] = g.skew
	copy(ulid[UnixOffset:PayloadOffset], millis[2:8])

	n, err := g.fill(ulid, ulid[PayloadOffset:])
	if err != nil {
		return nil, err
	} else if n != len(ulid)-PayloadOffset {
		return nil, ErrInsufficientData
	}

	return ulid, nil
}
