package paxos

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/cenkalti/backoff/v4"

	"github.com/mjpitz/myago/cluster"
)

type MultiAcceptorClient struct {
	Dialer   func(ctx context.Context, member string) (AcceptorClient, error)
	cache    *sync.Map
	size     int32
	majority int32
}

func sendPrepare(ctx context.Context, member string, client AcceptorClient, request *Request, ch chan *Vote) {
	promise, _ := client.Prepare(ctx, request)
	ch <- &Vote{Member: member, Payload: promise}
}

func (m *MultiAcceptorClient) Prepare(ctx context.Context, request *Request) (*Promise, error) {
	majority := int(atomic.LoadInt32(&(m.majority)))
	size := int(atomic.LoadInt32(&(m.size)))

	if size == 0 || size < majority {
		return &Promise{}, nil
	}

	votes := make(chan *Vote, size)
	defer close(votes)

	m.cache.Range(func(key, value interface{}) bool {
		member := key.(string)
		client := value.(AcceptorClient)

		go sendPrepare(ctx, member, client, request, votes)
		return true
	})

	promises := make([]*Promise, 0, size)

	for i := 0; i < size; i++ {
		vote := <-votes
		if vote.Payload == nil {
			continue
		}

		promise := vote.Payload.(*Promise)
		idx := sort.Search(len(promises), func(i int) bool {
			return promise.ID < promises[i].ID
		})

		if idx == len(promises) {
			promises = append(promises, promise)
		} else {
			promises = append(promises[:idx], append([]*Promise{promise}, promises[idx:]...)...)
		}
	}

	if len(promises) < majority {
		return &Promise{}, nil
	}

	var greatest *Proposal
	votesForRequest := 0

	for _, promise := range promises {
		if promise.Accepted != nil {
			if greatest == nil {
				greatest = promise.Accepted
			} else if greatest.ID < promise.Accepted.ID {
				greatest = promise.Accepted
			}
		}

		if promise.ID == request.ID {
			votesForRequest++
		}
	}

	if majority <= votesForRequest {
		return &Promise{
			ID:       request.ID,
			Accepted: greatest,
		}, nil
	}

	return &Promise{}, nil
}

func sendAccept(ctx context.Context, member string, client AcceptorClient, proposal *Proposal, ch chan *Vote) {
	proposal, _ = client.Accept(ctx, proposal)
	ch <- &Vote{Member: member, Payload: proposal}
}

func (m *MultiAcceptorClient) Accept(ctx context.Context, in *Proposal) (*Proposal, error) {
	majority := int(atomic.LoadInt32(&(m.majority)))
	size := int(atomic.LoadInt32(&(m.size)))

	if size == 0 {
		return &Proposal{}, nil
	}

	votes := make(chan *Vote, size)
	defer close(votes)

	sent := 0
	m.cache.Range(func(key, value interface{}) bool {
		member := key.(string)
		client := value.(AcceptorClient)

		go sendAccept(ctx, member, client, in, votes)
		sent++
		return true
	})

	proposals := make([]*Proposal, 0, size)

	for i := 0; i < sent; i++ {
		vote := <-votes
		if vote.Payload == nil {
			continue
		}

		proposal := vote.Payload.(*Proposal)
		idx := sort.Search(len(proposals), func(i int) bool {
			return proposal.ID < proposals[i].ID
		})

		if idx == len(proposals) {
			proposals = append(proposals, proposal)
		} else {
			proposals = append(proposals[:idx], append([]*Proposal{proposal}, proposals[idx:]...)...)
		}
	}

	if len(proposals) < majority {
		return &Proposal{}, nil
	}

	votesForRequest := 0
	for _, proposal := range proposals {
		if in.ID == proposal.ID {
			votesForRequest++
		}
	}

	if majority <= votesForRequest {
		return in, nil
	}

	return &Proposal{}, nil
}

func (m *MultiAcceptorClient) add(ctx context.Context, member string) {
	var client AcceptorClient

	var err error
	err = backoff.Retry(func() error {
		client, err = m.Dialer(ctx, member)
		return err
	}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))

	if err != nil {
		// log
		return
	}

	_, loaded := m.cache.LoadOrStore(member, client)
	if !loaded {
		atomic.AddInt32(&(m.size), 1)
	}
}

func (m *MultiAcceptorClient) remove(ctx context.Context, member string) {
	_, loaded := m.cache.LoadAndDelete(member)
	if loaded {
		atomic.AddInt32(&(m.size), -1)
	}
}

func (m *MultiAcceptorClient) handleMembershipChange(ctx context.Context, change cluster.MembershipChange) {
	for _, active := range change.Active {
		go m.add(ctx, active)
	}

	for _, removed := range change.Removed {
		go m.remove(ctx, removed)
	}
}

func (m *MultiAcceptorClient) Start(ctx context.Context, membership *cluster.Membership) error {
	changes, cancel := membership.Watch()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil
		case change := <-changes:
			atomic.SwapInt32(&(m.majority), int32(membership.Majority()))
			m.handleMembershipChange(ctx, change)
		}
	}
}

var (
	_ cluster.Discovery = &MultiAcceptorClient{}
	_ AcceptorClient    = &MultiAcceptorClient{}
)
