# encoding
--
    import "github.com/mjpitz/myago/encoding"

Package encoding provides common Encoding and associated interfaces for Encoder
and Decoder logic.

## Usage

```go
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
```

#### type Decoder

```go
type Decoder interface {
	Decode(v interface{}) error
}
```

Decoder defines how objects are decoded.

#### func  FromUnmarshaler

```go
func FromUnmarshaler(reader io.Reader, unmarshaler Unmarshaler) Decoder
```
FromUnmarshaler returns a Decoder that reads messages from the provided reader
and decodes them using the provided unmarshaler.

#### type Encoder

```go
type Encoder interface {
	Encode(v interface{}) error
}
```

Encoder defines how objects are encoded.

#### func  FromMarshaler

```go
func FromMarshaler(writer io.Writer, marshaler Marshaler) Encoder
```
FromMarshaler returns an Encoder that writes messages to the target writer
encoded with the results of the provided marshaller.

#### type Encoding

```go
type Encoding struct {
	// Encoder produces a new marshaledEncoder that can write messages to the provided io.Writer.
	Encoder func(w io.Writer) Encoder
	// Decoder produces a new decoder that can read messages from the provided io.Reader.
	Decoder func(r io.Reader) Decoder
}
```

Encoding defines the encoding of a file.

#### type Marshaler

```go
type Marshaler func(v interface{}) ([]byte, error)
```

Marshaler defines an interface for marshalling data into its binary/text
representation.

#### type Unmarshaler

```go
type Unmarshaler func(data []byte, v interface{}) error
```

Unmarshaler defines an interface for unmarshalling data for an interface from a
given byte array.
