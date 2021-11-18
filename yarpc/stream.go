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

package yarpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hashicorp/yamux"

	"github.com/mjpitz/myago/encoding"
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
		encoding: encoding.MsgPack,
		stream:   ys,
	}

	for _, opt := range opts {
		opt(rs)
	}

	rs.encoder = rs.encoding.Encoder(ys)
	rs.decoder = rs.encoding.Decoder(ys)

	return rs
}

type rpcStream struct {
	context  context.Context
	encoding *encoding.Encoding
	stream   *yamux.Stream

	encoder encoding.Encoder
	decoder encoding.Decoder
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
