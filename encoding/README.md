# encoding

Package encoding provides common Encoding and associated interfaces for Encoder
and Decoder logic.

```go
import github.com/mjpitz/myago/encoding
```

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

#### type Encoder

```go
type Encoder interface {
	Encode(v interface{}) error
}
```

Encoder defines how objects are encoded.

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
