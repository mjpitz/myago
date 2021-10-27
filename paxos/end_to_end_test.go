package paxos_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/mjpitz/myago/paxos"
	"github.com/mjpitz/myago/yarpc"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func exampleServer(t *testing.T, ctx context.Context, network, address string) error {
	promiseLog := &paxos.MemoryLog{}
	acceptedLog := &paxos.MemoryLog{}

	acceptor, err := paxos.NewAcceptor(promiseLog, acceptedLog)
	if err != nil {
		return err
	}

	paxos.RegisterYarpcAcceptorServer(yarpc.DefaultServer, acceptor)

	// ignore this error as it's likely a "socket closing" type of thing
	_ = yarpc.ListenAndServe(network, address, yarpc.WithContext(ctx))
	return nil
}

func exampleClient(t *testing.T, ctx context.Context, network, address string) error {
	defer yarpc.DefaultServer.Shutdown()

	clock := clockwork.NewFakeClock()
	conn := yarpc.DialContext(ctx, network, address)

	proposer := &paxos.Proposer{
		Clock:       clock,
		IDGenerator: paxos.ServerIDGenerator(1, clock),
		Acceptor:    paxos.NewYarpcAcceptorClient(conn),
	}

	accepted, err := proposer.Propose(ctx, []byte("hello-world"))
	if err != nil {
		return err
	}

	require.Equal(t, "hello-world", string(accepted))
	return nil
}

func TestYarpc(t *testing.T) {
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
