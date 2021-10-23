package yarpc

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

// Encoder writes provided structures to the underlying stream.
type Encoder interface {
	Encode(i interface{}) error
}

// Decoder reads message from the underlying stream.
type Decoder interface {
	Decode(i interface{}) error
}

// Encoding describes a generalization used to create encoders and decoders for new streams.
type Encoding interface {
	NewEncoder(io.Writer) Encoder
	NewDecoder(io.Reader) Decoder
}

// MSGPackEncoding uses msgpack out of box for a better balance of read/write performance. JSON serialization is fast,
// but deserialization is much slower in comparison (over 3x). While msgpack isn't as fast as protobuf, it offers
// reasonable read/write performance.
type MSGPackEncoding struct{}

func (j *MSGPackEncoding) NewEncoder(writer io.Writer) Encoder {
	return msgpack.NewEncoder(writer)
}

func (j *MSGPackEncoding) NewDecoder(reader io.Reader) Decoder {
	return msgpack.NewDecoder(reader)
}

var _ Encoding = &MSGPackEncoding{}
