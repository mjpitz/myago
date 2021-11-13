package yarpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/yamux"
)

func nonce() string {
	nonce := make([]byte, 16)

	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	return hex.EncodeToString(nonce)
}

type Stream interface {
	Context() context.Context
	SetReadDeadline(deadline time.Time) error
	ReadMsg(i interface{}) error
	SetWriteDeadline(deadline time.Time) error
	WriteMsg(i interface{}) error
	Close() error
}

func wrap(ys *yamux.Stream, opts ...streamOption) *rpcStream {
	rs := &rpcStream{
		context:  context.Background(),
		encoding: &MSGPackEncoding{},
		stream:   ys,
	}

	for _, opt := range opts {
		opt(rs)
	}

	rs.encoder = rs.encoding.NewEncoder(ys)
	rs.decoder = rs.encoding.NewDecoder(ys)

	return rs
}

type rpcStream struct {
	context  context.Context
	encoding Encoding
	stream   *yamux.Stream

	encoder Encoder
	decoder Decoder
}

func (j *rpcStream) Context() context.Context {
	return j.context
}

func (j *rpcStream) SetReadDeadline(deadline time.Time) error {
	return j.stream.SetReadDeadline(deadline)
}

func (j *rpcStream) ReadMsg(i interface{}) error {
	frame := &Frame{
		Body: i,
	}

	if err := j.decoder.Decode(frame); err != nil {
		return err
	}

	if frame.Status != nil {
		return fmt.Errorf("%d: %s", frame.Status.Code, frame.Status.Message)
	}

	return nil
}

func (j *rpcStream) SetWriteDeadline(deadline time.Time) error {
	return j.stream.SetWriteDeadline(deadline)
}

func (j *rpcStream) WriteMsg(i interface{}) error {
	frame := &Frame{
		Nonce: nonce(),
		Body:  i,
	}

	err, ok := i.(*Status)
	if ok {
		frame.Status = err
		frame.Body = nil
	}

	return j.encoder.Encode(frame)
}

func (j *rpcStream) Close() error {
	return j.stream.Close()
}

var _ Stream = &rpcStream{}
