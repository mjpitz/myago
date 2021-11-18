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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/pelletier/go-toml"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

var (
	// JSON defines a common structure for handling JSON encoding.
	JSON = &Encoding{
		Encoder: func(w io.Writer) Encoder {
			return json.NewEncoder(w)
		},
		Decoder: func(r io.Reader) Decoder {
			return json.NewDecoder(r)
		},
	}

	// MsgPack defines a common structure for handling MsgPack encoding.
	MsgPack = &Encoding{
		Encoder: func(w io.Writer) Encoder {
			return msgpack.NewEncoder(w)
		},
		Decoder: func(r io.Reader) Decoder {
			return msgpack.NewDecoder(r)
		},
	}

	// ProtoText defines a common structure for handling protobuf text encoding.
	ProtoText = &Encoding{
		Encoder: func(w io.Writer) Encoder {
			return FromMarshaler(w, func(v interface{}) ([]byte, error) {
				m, ok := v.(proto.Message)
				if !ok {
					return nil, fmt.Errorf("value is not a protobuf")
				}

				return prototext.Marshal(m)
			})
		},
		Decoder: func(r io.Reader) Decoder {
			return FromUnmarshaler(r, func(data []byte, v interface{}) error {
				m, ok := v.(proto.Message)
				if !ok {
					return fmt.Errorf("value is not a protobuf")
				}

				return prototext.Unmarshal(data, m)
			})
		},
	}

	// TOML defines a common structure for handling TOML encoding.
	TOML = &Encoding{
		Encoder: func(w io.Writer) Encoder {
			return toml.NewEncoder(w)
		},
		Decoder: func(r io.Reader) Decoder {
			return toml.NewDecoder(r)
		},
	}

	// YAML defines a common structure for handling YAML encoding.
	YAML = &Encoding{
		Encoder: func(w io.Writer) Encoder {
			return yaml.NewEncoder(w)
		},
		Decoder: func(r io.Reader) Decoder {
			return yaml.NewDecoder(r)
		},
	}

	// XML defines a common structure for handling XML encoding.
	XML = &Encoding{
		Encoder: func(w io.Writer) Encoder {
			return xml.NewEncoder(w)
		},
		Decoder: func(r io.Reader) Decoder {
			return xml.NewDecoder(r)
		},
	}
)
