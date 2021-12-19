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

	"github.com/mjpitz/myago/yarpc"
)

// RegisterYarpcAcceptorServer registers the provided AcceptorServer implementation with the yarpc.Server to handle
// requests.
func RegisterYarpcAcceptorServer(svr *yarpc.ServeMux, impl AcceptorServer) {
	svr.Handle("/paxos.acceptor/Prepare", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
		req := &Request{}
		err := stream.ReadMsg(req)
		if err != nil {
			return err
		}

		promise, err := impl.Prepare(stream.Context(), req)
		if err != nil {
			return err
		}

		return stream.WriteMsg(promise)
	}))

	svr.Handle("/paxos.acceptor/Accept", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
		req := &Proposal{}
		err := stream.ReadMsg(req)
		if err != nil {
			return err
		}

		proposal, err := impl.Accept(stream.Context(), req)
		if err != nil {
			return err
		}

		return stream.WriteMsg(proposal)
	}))
}

// NewYarpcAcceptorClient wraps the provided yarpc.ClientConn with an AcceptorClient implementation.
func NewYarpcAcceptorClient(cc *yarpc.ClientConn) AcceptorClient {
	return &yarpcAcceptorClient{
		cc: cc,
	}
}

type yarpcAcceptorClient struct {
	cc *yarpc.ClientConn
}

func (c *yarpcAcceptorClient) Prepare(ctx context.Context, request *Request) (*Promise, error) {
	stream, err := c.cc.OpenStream(ctx, "/paxos.acceptor/Prepare")
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	err = stream.WriteMsg(request)
	if err != nil {
		return nil, err
	}

	promise := &Promise{}

	return promise, stream.ReadMsg(promise)
}

func (c *yarpcAcceptorClient) Accept(ctx context.Context, proposal *Proposal) (*Proposal, error) {
	stream, err := c.cc.OpenStream(ctx, "/paxos.acceptor/Accept")
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	err = stream.WriteMsg(proposal)
	if err != nil {
		return nil, err
	}

	proposal = &Proposal{}

	return proposal, stream.ReadMsg(proposal)
}

var _ AcceptorClient = &yarpcAcceptorClient{}

// RegisterYarpcProposerServer registers the provided ProposerServer implementation with the yarpc.Server to handle
// requests. Typically, proposers aren't embedded as a server and are instead run as client side code.
func RegisterYarpcProposerServer(svr *yarpc.ServeMux, impl ProposerServer) {
	svr.Handle("/paxos.Proposer/Propose", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
		req := &Bytes{}
		err := stream.ReadMsg(req)
		if err != nil {
			return err
		}

		bytes, err := impl.Propose(stream.Context(), req.Value)
		if err != nil {
			return err
		}

		return stream.WriteMsg(&Bytes{
			Value: bytes,
		})
	}))
}

// NewYarpcProposerClient wraps the provided yarpc.ClientConn with an ProposerClient implementation.
func NewYarpcProposerClient(cc *yarpc.ClientConn) ProposerClient {
	return &yarpcProposerClient{
		cc: cc,
	}
}

type yarpcProposerClient struct {
	cc *yarpc.ClientConn
}

func (c *yarpcProposerClient) Propose(ctx context.Context, value []byte) ([]byte, error) {
	stream, err := c.cc.OpenStream(ctx, "/paxos.Proposer/Propose")
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	err = stream.WriteMsg(&Bytes{
		Value: value,
	})
	if err != nil {
		return nil, err
	}

	bytes := &Bytes{}
	err = stream.ReadMsg(bytes)
	if err != nil {
		return nil, err
	}

	return bytes.Value, nil
}

var _ ProposerClient = &yarpcProposerClient{}

// RegisterYarpcObserverServer registers the provided ObserverServer implementation with the yarpc.Server to handle
// requests. Acceptors should implement the observer server, otherwise other members of the cluster cannot determine
// what records have been accepted.
func RegisterYarpcObserverServer(svr *yarpc.ServeMux, impl ObserverServer) {
	svr.Handle("/paxos.Observer/Observe", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
		return impl.Observe(&ObserveServerStream{stream})
	}))
}

// NewYarpcObserverClient wraps the provided yarpc.ClientConn with an ObserverClient implementation.
func NewYarpcObserverClient(cc *yarpc.ClientConn) ObserverClient {
	return &yarpcObserverClient{
		cc: cc,
	}
}

type yarpcObserverClient struct {
	cc *yarpc.ClientConn
}

func (c *yarpcObserverClient) Observe(ctx context.Context, request *Request) (*ObserveClientStream, error) {
	stream, err := c.cc.OpenStream(ctx, "/paxos.Observer/Observe")
	if err != nil {
		return nil, err
	}

	err = stream.WriteMsg(request)
	if err != nil {
		return nil, err
	}

	return &ObserveClientStream{
		Stream: stream,
	}, nil
}

var _ ObserverClient = &yarpcObserverClient{}
