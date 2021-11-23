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

package paxos_test

import (
	"context"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/paxos"
	"github.com/mjpitz/myago/zaputil"
)

// nolint:funlen // idc about length for tests
func testLog(ctx context.Context, t *testing.T, root paxos.Log) {
	t.Helper()

	promiseLog := root.WithPrefix("promised/")
	acceptedLog := root.WithPrefix("accepted/")

	acceptor, err := paxos.NewAcceptor(promiseLog, acceptedLog)
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

		lastPromise := &paxos.Promise{}
		err = promiseLog.Last(lastPromise)
		require.NoError(t, err)
		require.Equal(t, uint64(1), lastPromise.ID)
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

		lastAccept := &paxos.Proposal{}
		err = acceptedLog.Last(lastAccept)
		require.NoError(t, err)
		require.Equal(t, uint64(1), lastAccept.ID)
		require.Equal(t, "hello-paxos", string(lastAccept.Value))
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

		lastPromise := &paxos.Promise{}
		err = promiseLog.Last(lastPromise)
		require.NoError(t, err)
		require.Equal(t, uint64(2), lastPromise.ID)
		require.NotNil(t, lastPromise.Accepted)

		require.Equal(t, uint64(1), lastPromise.Accepted.ID)
		require.Equal(t, "hello-paxos", string(lastPromise.Accepted.Value))
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

		lastAccept := &paxos.Proposal{}
		err = acceptedLog.Last(lastAccept)
		require.NoError(t, err)
		require.Equal(t, uint64(2), lastAccept.ID)
		require.Equal(t, "hello-paxos", string(lastAccept.Value))
	}

	// nolint:dupl
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

func TestBadger(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := zaputil.Extract(ctx)

	opts := badger.DefaultOptions(t.TempDir()).
		WithSyncWrites(true).
		WithLogger(zaputil.BadgerLogger(log))

	db, err := badger.Open(opts)
	require.NoError(t, err)

	defer db.Close()

	root := &paxos.Badger{DB: db}

	testLog(ctx, t, root)
}

func TestMemory(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testLog(ctx, t, &paxos.Memory{})
}
