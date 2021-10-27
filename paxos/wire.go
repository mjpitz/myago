package paxos

import (
	"context"
	"time"
)

type Stream interface {
	Context() context.Context
	SetReadDeadline(deadline time.Time) error
	ReadMsg(i interface{}) error
	SetWriteDeadline(deadline time.Time) error
	WriteMsg(i interface{}) error
	Close() error
}

type Bytes struct {
	Bytes []byte `json:"bytes,omitempty"`
}

type Request struct {
	ID      uint64 `json:"id,omitempty"`
	Attempt uint64 `json:"attempt,omitempty"`
}

type Proposal struct {
	ID    uint64 `json:"id,omitempty"`
	Value []byte `json:"value,omitempty"`
}

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
	Observe(call *ObserveServerStream) error
}

type AcceptorClient interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
	Observe(ctx context.Context, request *Request) (*ObserveClientStream, error)
}

type ProposerServer interface {
	Propose(ctx context.Context, value []byte) ([]byte, error)
}

type ProposerClient interface {
	Propose(ctx context.Context, value []byte) ([]byte, error)
}
