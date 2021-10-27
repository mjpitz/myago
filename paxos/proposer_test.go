package paxos_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/paxos"
)

// this is quite similar to the yarpc client. would be good to generalize the yarpc.ClientConn definition too...
type mockAcceptor struct {
	mockStream *paxos.MockStream
}

func (m *mockAcceptor) Prepare(ctx context.Context, request *paxos.Request) (*paxos.Promise, error) {
	err := m.mockStream.WriteMsg(request)
	if err != nil {
		return nil, err
	}

	promise := &paxos.Promise{}
	return promise, m.mockStream.ReadMsg(promise)
}

func (m *mockAcceptor) Accept(ctx context.Context, proposal *paxos.Proposal) (*paxos.Proposal, error) {
	err := m.mockStream.WriteMsg(proposal)
	if err != nil {
		return nil, err
	}

	proposal = &paxos.Proposal{}
	return proposal, m.mockStream.ReadMsg(proposal)
}

func (m *mockAcceptor) Observe(ctx context.Context, request *paxos.Request) (*paxos.ObserveClientStream, error) {
	err := m.mockStream.WriteMsg(request)
	if err != nil {
		return nil, err
	}

	return &paxos.ObserveClientStream{Stream: m.mockStream}, nil
}

var _ paxos.AcceptorClient = &mockAcceptor{}

// TestProposer_Simple runs a typical paxos run where the value proposed is the value
func TestProposer_Simple(t *testing.T) {
	ctx := context.Background()
	clock := clockwork.NewFakeClockAt(time.Now())

	acceptorStream := paxos.NewMockStream(5)
	acceptorStream.Ctx = ctx

	proposer := &paxos.Proposer{
		Clock:       clock,
		IDGenerator: paxos.ServerIDGenerator(1, clock),
		Acceptor: &mockAcceptor{
			mockStream: acceptorStream,
		},
	}

	id, err := proposer.IDGenerator.Next()
	require.NoError(t, err)

	// remember...
	// - proposer sends prepare 0..n times
	// - acceptor responds with promise
	// - proposer sends accept
	// - acceptor response with promise
	acceptorStream.Incoming <- &paxos.Promise{
		ID: id,
	}

	acceptorStream.Incoming <- &paxos.Proposal{
		ID:    id,
		Value: []byte("alice"),
	}

	accepted, err := proposer.Propose(ctx, []byte("alice"))
	require.NoError(t, err)

	{ // verify prepare messages sent
		request, ok := (<-acceptorStream.Outgoing).(*paxos.Request)
		require.True(t, ok, "message was not a *paxosv1.Request")
		require.Equal(t, id, request.ID)
	}

	{ // verify accept messages sent
		request, ok := (<-acceptorStream.Outgoing).(*paxos.Proposal)
		require.True(t, ok, "message was not a *paxosv1.Proposal")
		require.Equal(t, id, request.ID)
		require.Equal(t, "alice", string(request.Value))
	}

	require.Equal(t, "alice", string(accepted))
}
