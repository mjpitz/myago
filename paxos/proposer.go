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
		if promise.Accepted != nil && len(promise.Accepted.Value) > 0 {
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

var (
	_ ProposerClient = &Proposer{}
	_ ProposerServer = &Proposer{}
)
