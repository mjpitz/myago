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
	"sync"

	"github.com/mjpitz/myago/yarpc"
)

type Acceptor interface {
	AcceptorServer
	ObserverServer
}

func NewAcceptor(promiseLog, acceptedLog Log) (Acceptor, error) {
	lastPromise, lastAccept := &Promise{}, &Proposal{}

	// read the last entries
	if err := promiseLog.Last(lastPromise); err != nil {
		return nil, err
	} else if err := acceptedLog.Last(lastAccept); err != nil {
		return nil, err
	}

	return &acceptor{
		lastPromise: lastPromise,
		lastAccept:  lastAccept,
		promiseLog:  promiseLog,
		acceptedLog: acceptedLog,
		updates:     make(map[yarpc.Stream]chan *Proposal),
	}, nil
}

type acceptor struct {
	mu          sync.Mutex
	lastPromise *Promise
	lastAccept  *Proposal
	promiseLog  Log
	acceptedLog Log
	updates     map[yarpc.Stream]chan *Proposal
}

func (a *acceptor) Prepare(ctx context.Context, req *Request) (*Promise, error) {
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

func (a *acceptor) Accept(ctx context.Context, proposal *Proposal) (*Proposal, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if proposal.ID < a.lastPromise.ID {
		return &Proposal{}, nil
	}

	err := a.acceptedLog.Record(proposal.ID, proposal)
	if err != nil {
		return nil, err
	}

	a.lastAccept = proposal

	for _, stream := range a.updates {
		stream <- proposal
	}

	return proposal, nil
}

func (a *acceptor) Observe(call *ObserveServerStream) error {
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

	err = a.acceptedLog.Range(req.ID, lastAcceptID, Proposal{}, func(msg interface{}) error {
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

var (
	_ AcceptorServer = &acceptor{}
	_ ObserverServer = &acceptor{}
)
