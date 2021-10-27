package paxos_test

import (
	"context"
	"testing"

	"github.com/mjpitz/myago/paxos"
	"github.com/stretchr/testify/require"
)

func TestAcceptor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	promiseLog := &paxos.MemoryLog{}
	acceptLog := &paxos.MemoryLog{}

	acceptor, err := paxos.NewAcceptor(promiseLog, acceptLog)
	require.NoError(t, err)

	{
		t.Log("running phase 1")

		promise, err := acceptor.Prepare(ctx, &paxos.Request{
			ID:      1,
			Attempt: 1,
		})

		require.NoError(t, err)
		require.Equal(t, uint64(1), promise.ID)
		require.Nil(t, promise.Accepted)
	}

	{
		t.Log("running phase 2")

		proposal, err := acceptor.Accept(ctx, &paxos.Proposal{
			ID:    1,
			Value: []byte("hello-paxos"),
		})

		require.NoError(t, err)
		require.Equal(t, uint64(1), proposal.ID)
		require.Equal(t, "hello-paxos", string(proposal.Value))
	}

	{
		t.Log("running phase 1 - learn")

		promise, err := acceptor.Prepare(ctx, &paxos.Request{
			ID:      1,
			Attempt: 1,
		})

		require.NoError(t, err)
		require.Equal(t, uint64(0), promise.ID)

		promise, err = acceptor.Prepare(ctx, &paxos.Request{
			ID:      2,
			Attempt: 2,
		})

		require.NoError(t, err)
		require.Equal(t, uint64(2), promise.ID)
		require.NotNil(t, promise.Accepted)

		require.Equal(t, uint64(1), promise.Accepted.ID)
		require.Equal(t, "hello-paxos", string(promise.Accepted.Value))
	}

	{
		t.Log("running phase 2 - learn")

		proposal, err := acceptor.Accept(ctx, &paxos.Proposal{
			ID:    2,
			Value: []byte("hello-paxos"),
		})

		require.NoError(t, err)
		require.Equal(t, uint64(2), proposal.ID)
		require.Equal(t, "hello-paxos", string(proposal.Value))
	}

	{
		t.Log("verifying observations")
		observeStream := paxos.NewMockStream(5)
		observeStream.Ctx = ctx
		observeStream.Incoming <- &paxos.Request{}

		go func() {
			err := acceptor.Observe(&paxos.ObserveServerStream{Stream: observeStream})
			require.NoError(t, err)
		}()

		{ // learn
			msg := <-observeStream.Outgoing
			proposal, ok := msg.(*paxos.Proposal)
			require.True(t, ok, "Outgoing message was not a proposal")
			require.Equal(t, uint64(1), proposal.ID)
			require.Equal(t, "hello-paxos", string(proposal.Value))
		}

		{ // re-learn
			msg := <-observeStream.Outgoing
			proposal, ok := msg.(*paxos.Proposal)
			require.True(t, ok, "Outgoing message was not a proposal")
			require.Equal(t, uint64(2), proposal.ID)
			require.Equal(t, "hello-paxos", string(proposal.Value))
		}
	}
}
