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

package plugin

import (
	"bytes"
	"context"
	"io"
)

// NewBlockingReadWriteCloser returns a new io.ReadWriteCloser that supports a blocking Read operation. Write will push
// new data into the channel to be added to the buffer after.
func NewBlockingReadWriteCloser() *blockingReadWriteCloser {
	ctx, cancel := context.WithCancel(context.Background())
	buf := make(chan *bytes.Buffer, 1)
	buf <- bytes.NewBuffer(nil)

	return &blockingReadWriteCloser{
		context: ctx,
		cancel:  cancel,
		buf:     buf,
	}
}

type blockingReadWriteCloser struct {
	context context.Context
	cancel  context.CancelFunc
	buf     chan *bytes.Buffer
}

func (b *blockingReadWriteCloser) Read(p []byte) (n int, err error) {
	// buf has priority on reads
	buf := <-b.buf
	n, _ = buf.Read(p)

	err = b.context.Err()
	if buf.Len() == 0 && err != nil {
		err = io.EOF
	}

	b.buf <- buf

	return n, err
}

func (b *blockingReadWriteCloser) Write(p []byte) (n int, err error) {
	// context has priority on writes
	select {
	case <-b.context.Done():
		err = io.ErrClosedPipe
	default:
		select {
		case buf := <-b.buf:
			defer func() {
				b.buf <- buf
			}()

			return buf.Write(p)
		case <-b.context.Done():
			err = io.ErrClosedPipe
		}
	}

	return n, err
}

func (b *blockingReadWriteCloser) Close() error {
	b.cancel()
	return nil
}

var _ io.ReadWriteCloser = &blockingReadWriteCloser{}
