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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/plugin"
)

//go:generate go install ./examples/myago-plugin-echo

type message struct {
	Text string
}

func TestClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

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
