package paxos

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/jonboulle/clockwork"
)

// Proposer can be run either as an embedded client, or as part of a standalone server. Proposers Propose additions to
// the paxos log and uses the acceptors to get consensus on if the proposed value was accepted.
type Proposer struct {
	Clock       clockwork.Clock
	IDGenerator IDGenerator
	Acceptor    AcceptorClient
}

func (p *Proposer) prepare(ctx context.Context) (*Promise, error) {
	for attempt := uint64(1); ; attempt++ {
		nextID, err := p.IDGenerator.Next()
		if err != nil {
			return nil, err
		}

		promise, err := p.Acceptor.Prepare(ctx, &Request{
			ID:      nextID,
			Attempt: attempt,
		})
		if err != nil {
			return nil, err
		}

		if promise.ID == nextID {
			return promise, nil
		}
	}
}

func (p *Proposer) Propose(ctx context.Context, value []byte) (accepted []byte, err error) {
	err = backoff.Retry(func() error {
		promise, err := p.prepare(ctx)
		if err != nil {
			return err
		}

		accepted = value
		if promise.Accepted != nil {
			accepted = promise.Accepted.Value
		}

		proposal, err := p.Acceptor.Accept(ctx, &Proposal{
			ID:    promise.ID,
			Value: accepted,
		})
		if err != nil {
			return err
		}

		accepted = proposal.Value
		return nil
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return nil, err
	}

	return accepted, nil
}

var _ ProposerClient = &Proposer{}
var _ ProposerServer = &Proposer{}
