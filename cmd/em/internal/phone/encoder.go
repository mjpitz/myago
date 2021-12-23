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

package phone

import (
	"io"
)

var encoding = map[byte]byte{
	// upper
	65: 50, 66: 50, 67: 50,
	68: 51, 69: 51, 70: 51,
	71: 52, 72: 52, 73: 52,
	74: 53, 75: 53, 76: 53,
	77: 54, 78: 54, 79: 54,
	80: 55, 81: 55, 82: 55, 83: 55,
	84: 56, 85: 56, 86: 56,
	87: 57, 88: 57, 89: 57, 90: 57,

	// lower
	97: 50, 98: 50, 99: 50,
	100: 51, 101: 51, 102: 51,
	103: 52, 104: 52, 105: 52,
	106: 53, 107: 53, 108: 53,
	109: 54, 110: 54, 111: 54,
	112: 55, 113: 55, 114: 55, 115: 55,
	116: 56, 117: 56, 118: 56,
	119: 57, 120: 57, 121: 57, 122: 57,
}

// NewEncoder returns an encoder that translates data into a phone code.
func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		writer: writer,
	}
}

type Encoder struct {
	writer io.Writer
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	encoded := make([]byte, 0, len(p))

	for i := 0; i < len(p); i++ {
		v, ok := encoding[p[i]]
		if ok {
			encoded = append(encoded, v)
		}
	}

	_, err = e.writer.Write(encoded)
	return len(p), err
}

var _ io.Writer = &Encoder{}
