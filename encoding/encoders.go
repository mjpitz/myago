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
