package paxos_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/mjpitz/myago/cluster"
	"github.com/mjpitz/myago/paxos"
	"github.com/mjpitz/myago/yarpc"
)

// nolint:funlen // idc about length for tests
func TestPaxos(t *testing.T) {
	t.Parallel()

	network := "unix"
	clock := clockwork.NewFakeClock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newPaxos := func(id uint8) (*paxos.Paxos, error) {
		root := &paxos.Memory{}

		return paxos.New(&paxos.Config{
			Clock:       clock,
			IDGenerator: paxos.ServerIDGenerator(id, clock),
			PromiseLog:  root.WithPrefix("promised/"),
			AcceptedLog: root.WithPrefix("accepted/"),
			RecordedLog: root.WithPrefix("recorded/"),
			AcceptorDialer: func(ctx context.Context, member string) (paxos.AcceptorClient, error) {
				return paxos.NewYarpcAcceptorClient(yarpc.DialContext(ctx, network, member)), nil
			},
			ObserverDialer: func(ctx context.Context, member string) (paxos.ObserverClient, error) {
				return paxos.NewYarpcObserverClient(yarpc.DialContext(ctx, network, member)), nil
			},
		})
	}

	numServers := uint8(3)

	socks := make([]string, 0, numServers)
	defer func() {
		for _, sock := range socks {
			_ = os.Remove(sock)
		}
	}()

	paxi := make([]*paxos.Paxos, 0, numServers)
	svrs := make([]*yarpc.Server, 0, numServers)

	defer func() {
		for _, svr := range svrs {
			_ = svr.Shutdown()
		}
	}()

	for i := uint8(0); i < numServers; i++ {
		sock := path.Join(t.TempDir(), fmt.Sprintf("%s-%d.sock", t.Name(), i))
		socks = append(socks, sock)

		pax, err := newPaxos(i)
		require.NoError(t, err)

		svr := &yarpc.Server{}
		paxos.RegisterYarpcAcceptorServer(svr, pax)
		paxos.RegisterYarpcObserverServer(svr, pax)

		svrContext := yarpc.WithContext(ctx)

		go func() {
			_ = svr.ListenAndServe(network, sock, svrContext)
		}()

		paxi = append(paxi, pax)
		svrs = append(svrs, svr)
	}

	membership := new(cluster.Membership)
	membership.Add(socks)

	waitForStartup := sync.WaitGroup{}
	waitForStartup.Add(len(paxi) + 1)

	// spin up observers and acceptors
	ctx, shutdown := context.WithCancel(ctx)
	defer shutdown()

	group, ctx := errgroup.WithContext(ctx)
	submitToGroup := func(pax *paxos.Paxos) {
		group.Go(func() error {
			waitForStartup.Done()

			return pax.Start(ctx, membership)
		})
	}

	for _, pax := range paxi {
		submitToGroup(pax)
	}

	go func() {
		waitForStartup.Done()
		_ = group.Wait()
	}()

	t.Log("waiting for startup")
	waitForStartup.Wait()

	t.Log("picking random proposer")
	data := make([]byte, 1)
	_, err := io.ReadFull(rand.Reader, data)
	require.NoError(t, err)
	idx := int(data[0]) % len(paxi)

	t.Log("proposing value")
	request := []byte("hello paxos")
	accepted, err := paxi[idx].Propose(ctx, request)
	require.NoError(t, err)

	t.Log("asserting state")
	require.True(t, bytes.Equal(request, accepted), string(accepted))

	// verify logs?

	proposal := &paxos.Proposal{}
	attempt := 1

	for proposal.ID == 0 {
		t.Log("awaiting observer log attempt: ", attempt)
		time.Sleep(time.Second)

		_ = paxi[idx].Observer.Log.Last(proposal)
		attempt++
	}

	require.Equal(t, uint64(0x68bf39440001>>1<<1|idx), proposal.ID)
	require.True(t, bytes.Equal(request, proposal.Value), string(proposal.Value))
}
