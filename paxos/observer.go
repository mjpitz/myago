package paxos

import (
	"context"

	"github.com/cenkalti/backoff/v4"

	"github.com/mjpitz/myago/cluster"
)

// Observer watches the Acceptors to learn about what values have been accepted.
type Observer struct {
	Dialer func(ctx context.Context, member string) (ObserverClient, error)
	Log    Log
}

// nolint:cyclop
func (o *Observer) observe(ctx context.Context, member string, lastAccepted *Proposal, votes chan *Vote) {
	var client ObserverClient
	var observations *ObserveClientStream

	var err error
	err = backoff.Retry(func() error {
		client, err = o.Dialer(ctx, member)

		return err
	}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))

	if err != nil {
		// log
		return
	}

	run := func() error {
		err = backoff.Retry(func() error {
			observations, err = client.Observe(ctx, &Request{
				ID: lastAccepted.ID,
			})

			return err
		}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))

		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				proposal, err := observations.Recv()
				if err != nil {
					return err
				}

				select {
				case votes <- &Vote{Member: member, Payload: proposal}:
				case <-ctx.Done():
					return nil
				}
			}
		}
	}

	for {
		err = run()
		if err == nil {
			return
		}
	}
}

// nolint:gocognit,cyclop
func (o *Observer) Start(ctx context.Context, membership *cluster.Membership) error {
	lastAccepted := &Proposal{}
	err := o.Log.Last(lastAccepted)
	if err != nil {
		return err
	}

	changes, cancel := membership.Watch()
	defer cancel()

	idx := make(map[string]context.CancelFunc)
	votes := make(chan *Vote, 16)
	tallies := make(map[uint64]map[string]bool)

	majority := membership.Majority()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case vote := <-votes:
			proposal := vote.Payload.(*Proposal)
			id := proposal.ID
			if _, ok := tallies[id]; !ok {
				tallies[id] = make(map[string]bool)
			}

			tallies[id][vote.Member] = true

			if majority <= len(tallies[id]) {
				err = o.Log.Record(id, proposal)
				if err != nil {
					continue
				}

				if lastAccepted.ID < id {
					lastAccepted.ID = id
					lastAccepted.Value = proposal.Value
				}
			}
		case change := <-changes:
			for _, active := range change.Active {
				if _, ok := idx[active]; !ok {
					child, childCancel := context.WithCancel(ctx)

					go o.observe(child, active, lastAccepted, votes)

					idx[active] = childCancel
				}
			}

			for _, left := range change.Left {
				if cancel, ok := idx[left]; ok {
					cancel()
					delete(idx, left)
				}
			}

			for _, removed := range change.Removed {
				if cancel, ok := idx[removed]; ok {
					cancel()
					delete(idx, removed)
				}
			}

			majority = membership.Majority()
		}
	}
}

var _ cluster.Discovery = &Observer{}
