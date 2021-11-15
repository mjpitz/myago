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
