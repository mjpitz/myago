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
	"errors"
)

var (
	// ErrInvalidBitCount is returned when an invalid number of bits is provided to the Generate method of a Generator.
	ErrInvalidBitCount = errors.New("bits must be divisible by 8")

	// ErrNotEnoughBits is returned when fewer than 64 bit ULIDs are requested to be generated.
	ErrNotEnoughBits = errors.New("must be at least 64 bits")

	// ErrInsufficientData is returned when the fill fails to return enough fata for the ULID.
	ErrInsufficientData = errors.New("failed to read sufficient payload data")
)
