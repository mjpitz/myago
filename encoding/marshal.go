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
	"io/ioutil"
)

// Unmarshaler defines an interface for unmarshalling data for an interface from a given byte array.
type Unmarshaler func(data []byte, v interface{}) error

// Marshaler defines an interface for marshalling data into its binary/text representation.
type Marshaler func(v interface{}) ([]byte, error)

// FromUnmarshaler returns a Decoder that reads messages from the provided reader and decodes them using the provided
// unmarshaler.
func FromUnmarshaler(reader io.Reader, unmarshaler Unmarshaler) Decoder {
	return &unmarshaledDecoder{
		unmarshaler: unmarshaler,
		reader:      reader,
	}
}

type unmarshaledDecoder struct {
	unmarshaler Unmarshaler
	reader      io.Reader
}

func (d *unmarshaledDecoder) Decode(v interface{}) error {
	data, err := ioutil.ReadAll(d.reader)
	if err != nil {
		return err
	}

	return d.unmarshaler(data, v)
}

// FromMarshaler returns an Encoder that writes messages to the target writer encoded with the results of the provided
// marshaller.
func FromMarshaler(writer io.Writer, marshaler Marshaler) Encoder {
	return &marshaledEncoder{
		marshaler: marshaler,
		writer:    writer,
	}
}

type marshaledEncoder struct {
	marshaler Marshaler
	writer    io.Writer
}

func (e *marshaledEncoder) Encode(v interface{}) error {
	data, err := e.marshaler(v)
	if err != nil {
		return err
	}

	_, err = e.writer.Write(data)

	return err
}
