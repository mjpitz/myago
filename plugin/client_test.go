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

package plugin_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/mjpitz/myago/lifecycle"
	"github.com/mjpitz/myago/plugin"
	"github.com/mjpitz/myago/zaputil"
)

//go:generate go install ./examples/myago-plugin-echo
//go:generate go install ./examples/myago-plugin-failure

type message struct {
	Text string
}

// nolint:paralleltest
func TestEchoClient(t *testing.T) {
	ctx := lifecycle.Setup(context.Background())
	defer lifecycle.Resolve(ctx)

	log, _ := zap.NewDevelopment()
	ctx = zaputil.ToContext(ctx, log)

	clientConn := plugin.DialContext(ctx, "myago-plugin-echo")

	stream, err := clientConn.OpenStream(ctx, "/echo")
	require.NoError(t, err)
	defer func() {
		_ = stream.Close()
	}()

	err = stream.WriteMsg(message{
		Text: "hello world",
	})
	require.NoError(t, err)

	msg := &message{}

	err = stream.ReadMsg(msg)
	require.NoError(t, err)

	require.Equal(t, "hello world", msg.Text)
}

// nolint:paralleltest
func TestFailureClient(t *testing.T) {
	ctx := lifecycle.Setup(context.Background())
	defer lifecycle.Resolve(ctx)

	log, _ := zap.NewDevelopment()
	ctx = zaputil.ToContext(ctx, log)

	clientConn := plugin.DialContext(ctx, "myago-plugin-failure")

	// this is a bit of a schrödinger's cat situation so we conditionally check the error at each step
	// I suspect the write path is less-likely to be blocked by the error since yamux uses an async approach

	stream, err := clientConn.OpenStream(ctx, "/echo")
	require.NoError(t, err)

	err = stream.WriteMsg(message{
		Text: "hello world",
	})
	require.NoError(t, err)

	msg := &message{}
	err = stream.ReadMsg(msg)
	require.Error(t, err)

	// yamux returns an EOF when the stream is closed
	// this doesn't seem to be something we can easily work around.
	require.Equal(t, io.EOF, err)
}
