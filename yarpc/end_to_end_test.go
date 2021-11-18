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

package yarpc_test

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/mjpitz/myago/yarpc"
)

const method = "example"

type Stat struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func exampleServer(ctx context.Context, t *testing.T, network, address string) error {
	t.Helper()

	start := time.Now()

	yarpc.HandleFunc(method, func(stream yarpc.Stream) error {
		require.NoError(t, stream.WriteMsg(&Stat{
			Name:  "uptime",
			Value: float64(time.Since(start).Milliseconds()),
		}))

		return nil
	})

	// ignore this error as it's likely a "socket closing" type of thing
	_ = yarpc.ListenAndServe(network, address, yarpc.WithContext(ctx))

	return nil
}

func exampleClient(ctx context.Context, t *testing.T, network, address string) error {
	t.Helper()

	defer func() {
		_ = yarpc.DefaultServer.Shutdown()
	}()

	conn := yarpc.DialContext(ctx, network, address)

	stream, err := conn.OpenStream(ctx, method)
	if err != nil {
		return err
	}

	stat := &Stat{}
	err = stream.ReadMsg(stat)
	if err != nil {
		return err
	}

	require.Equal(t, "uptime", stat.Name)
	require.Greater(t, stat.Value, float64(0.0))

	return stream.Close()
}

func TestEndToEnd(t *testing.T) {
	t.Parallel()

	sock := path.Join(t.TempDir(), t.Name()+".sock")
	defer os.Remove(sock)

	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return exampleServer(ctx, t, "unix", sock)
	})

	group.Go(func() error {
		return exampleClient(ctx, t, "unix", sock)
	})

	require.NoError(t, group.Wait())
}
