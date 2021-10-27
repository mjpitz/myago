package paxos

import (
	"context"
	"sync"

	"github.com/mjpitz/myago/yarpc"
)

func NewAcceptor(promiseLog, acceptLog Log) (*Acceptor, error) {
	lastPromise, lastAccept := &Promise{}, &Proposal{}

	// read the last entries
	if err := promiseLog.Last(lastPromise); err != nil {
		return nil, err
	} else if err := acceptLog.Last(lastAccept); err != nil {
		return nil, err
	}

	return &Acceptor{
		lastPromise: lastPromise,
		lastAccept:  lastAccept,
		promiseLog:  promiseLog,
		acceptLog:   acceptLog,
		updates:     make(map[yarpc.Stream]chan *Proposal),
	}, nil
}

type Acceptor struct {
	mu          sync.Mutex
	lastPromise *Promise
	lastAccept  *Proposal
	promiseLog  Log
	acceptLog   Log
	updates     map[yarpc.Stream]chan *Proposal
}

func (a *Acceptor) Prepare(ctx context.Context, req *Request) (*Promise, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if req.ID <= a.lastPromise.ID {
		return &Promise{}, nil
	}

	promise := &Promise{}
	promise.ID = req.ID

	if req.Attempt > 1 {
		promise.Accepted = a.lastAccept
	}

	err := a.promiseLog.Record(promise.ID, promise)
	if err != nil {
		return nil, err
	}

	a.lastPromise = promise

	return promise, nil
}

func (a *Acceptor) Accept(ctx context.Context, proposal *Proposal) (*Proposal, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if proposal.ID < a.lastPromise.ID {
		return &Proposal{}, nil
	}

	err := a.acceptLog.Record(proposal.ID, proposal)
	if err != nil {
		return nil, err
	}

	a.lastAccept = proposal

	for _, stream := range a.updates {
		stream <- proposal
	}

	return proposal, nil
}

func (a *Acceptor) Observe(call *ObserveServerStream) error {
	var lastAcceptID uint64

	a.mu.Lock()
	lastAcceptID = a.lastAccept.ID
	subscription := make(chan *Proposal, 5)
	a.updates[call] = subscription
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		delete(a.updates, call)
		a.mu.Unlock()
	}()

	req, err := call.Recv()
	if err != nil {
		return err
	}

	err = a.acceptLog.Range(req.ID, lastAcceptID, Proposal{}, func(msg interface{}) error {
		return call.WriteMsg(msg)
	})

	for err == nil {
		select {
		case <-call.Context().Done():
			return nil
		case proposal := <-subscription:
			err = call.Send(proposal)
		}
	}

	return err
}

var _ AcceptorServer = &Acceptor{}
