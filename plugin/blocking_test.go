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
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/plugin"
)

type message struct {
	Text string
}

func TestBlocking(t *testing.T) {
	t.Parallel()

	type message struct {
		Text string
	}

	channel := plugin.NewBlockingReadWriteCloser()
	reader := json.NewDecoder(channel)
	writer := json.NewEncoder(channel)

	msg := message{
		Text: "hello world",
	}
	err := writer.Encode(msg)
	require.NoError(t, err)

	msg = message{}
	err = reader.Decode(&msg)
	require.NoError(t, err)

	require.Equal(t, "hello world", msg.Text)

	err = channel.Close()
	require.NoError(t, err)

	_, err = channel.Read(make([]byte, 0))
	require.Equal(t, io.EOF, err)

	_, err = channel.Write(make([]byte, 0))
	require.Equal(t, io.ErrClosedPipe, err)
}
