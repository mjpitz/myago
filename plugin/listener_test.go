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
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/mjpitz/myago/yarpc"
)

type message struct {
	Text string
}

type mockDialer struct {
	conn *clientRWC
}

func (m *mockDialer) DialContext(ctx context.Context) (io.ReadWriteCloser, error) {
	return m.conn, nil
}

var _ yarpc.Dialer = &mockDialer{}

func TestListener(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	listener := listen(&serverRWC{
		stdin: stdinReader,
		stdout: stdoutWriter,
	})

	dialer := &mockDialer{
		conn: &clientRWC{
			stdout: stdoutReader,
			stdin: stdinWriter,
		},
	}

	var group errgroup.Group

	defer func() {
		_ = yarpc.DefaultServer.Shutdown()
		_ = group.Wait()
	}()

	group.Go(func() error {
		yarpc.HandleFunc("/echo", func(stream yarpc.Stream) (err error) {
			msg := message{}
			err = stream.ReadMsg(&msg)
			if err != nil {
				return
			}

			err = stream.WriteMsg(msg)
			if err != nil {
				return
			}

			return nil
		})

		return yarpc.Serve(listener, yarpc.WithContext(ctx))
	})

	clientConn := yarpc.NewClientConn(ctx)
	clientConn.Dialer = dialer

	performCall := func() {
		stream, err := clientConn.OpenStream(ctx, "/echo")
		require.NoError(t, err)

		err = stream.WriteMsg(message{
			Text: "hello world",
		})
		require.NoError(t, err)

		msg := &message{}

		err = stream.ReadMsg(msg)
		require.NoError(t, err)

		require.Equal(t, "hello world", msg.Text)
	}

	performCall()

	// make 5 successive calls
	//for callID := 1; callID <= 5; callID++ {
	//	performCall()
	//}
}
