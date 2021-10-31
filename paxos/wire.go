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
