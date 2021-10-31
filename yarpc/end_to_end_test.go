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

func exampleServer(t *testing.T, ctx context.Context, network, address string) error {
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

func exampleClient(t *testing.T, ctx context.Context, network, address string) error {
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
	sock := path.Join(t.TempDir(), t.Name()+".sock")
	defer os.Remove(sock)

	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return exampleServer(t, ctx, "unix", sock)
	})

	group.Go(func() error {
		return exampleClient(t, ctx, "unix", sock)
	})

	require.NoError(t, group.Wait())
}
