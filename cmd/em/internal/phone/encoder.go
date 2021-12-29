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

var encoding = map[byte]byte{}

var keyPad = map[byte][]byte{
	'2': []byte("abcABC"),
	'3': []byte("defDEF"),
	'4': []byte("ghiGHI"),
	'5': []byte("jklJKL"),
	'6': []byte("mnoMNO"),
	'7': []byte("pqrsPQRS"),
	'8': []byte("tuvTUV"),
	'9': []byte("wxyzWXYZ"),
	'0': []byte("@&%?,=[]_:-+*$#!'^~;()/."),
}

func init() {
	for num, values := range keyPad {
		for _, value := range values {
			encoding[value] = num
		}
	}
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
