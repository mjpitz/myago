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
)

// Bytes contains a value to be accepted via paxos.
type Bytes struct {
	Value []byte `json:"value,omitempty"`
}

// Request is used during the PREPARE and OBSERVE phases of the paxos algorithm. Prepare sends along their ID value and
// attempt number, where Observe sends along their last accepted id.
type Request struct {
	ID      uint64 `json:"id,omitempty"`
	Attempt uint64 `json:"attempt,omitempty"`
}

// Proposal is used to propose a log value to system.
type Proposal struct {
	ID    uint64 `json:"id,omitempty"`
	Value []byte `json:"value,omitempty"`
}

// Promise is returned by an accepted prepare. If more than one attempt was made, and accepted value is returned with
// the last accepted proposal so clients can catch up.
type Promise struct {
	ID       uint64    `json:"id,omitempty"`
	Accepted *Proposal `json:"accepted,omitempty"`
}

type ObserveServerStream struct {
	Stream
}

func (s *ObserveServerStream) Recv() (*Request, error) {
	msg := &Request{}

	return msg, s.ReadMsg(msg)
}

func (s *ObserveServerStream) Send(msg *Proposal) error {
	return s.WriteMsg(msg)
}

type ObserveClientStream struct {
	Stream
}

func (s *ObserveClientStream) Recv() (*Proposal, error) {
	msg := &Proposal{}

	return msg, s.ReadMsg(msg)
}

type AcceptorServer interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
}

type AcceptorClient interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
}

type ProposerServer interface {
	Propose(ctx context.Context, value []byte) ([]byte, error)
}

type ProposerClient interface {
	Propose(ctx context.Context, value []byte) ([]byte, error)
}

type ObserverServer interface {
	Observe(call *ObserveServerStream) error
}

type ObserverClient interface {
	Observe(ctx context.Context, request *Request) (*ObserveClientStream, error)
}
