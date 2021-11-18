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

package encoding

import (
	"io"
)

// Decoder defines how objects are decoded.
type Decoder interface {
	Decode(v interface{}) error
}

// Encoder defines how objects are encoded.
type Encoder interface {
	Encode(v interface{}) error
}

// Encoding defines the encoding of a file.
type Encoding struct {
	// Encoder produces a new marshaledEncoder that can write messages to the provided io.Writer.
	Encoder func(w io.Writer) Encoder
	// Decoder produces a new decoder that can read messages from the provided io.Reader.
	Decoder func(r io.Reader) Decoder
}
