package paxos

import (
	"context"

	"github.com/mjpitz/myago/yarpc"
)

func RegisterYarpcAcceptorServer(svr *yarpc.Server, impl AcceptorServer) {
	svr.Handle("/paxos.Acceptor/Prepare", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
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

	svr.Handle("/paxos.Acceptor/Accept", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
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

	svr.Handle("/paxos.Acceptor/Observe", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
		return impl.Observe(&ObserveServerStream{stream})
	}))
}

func NewYarpcAcceptorClient(cc *yarpc.ClientConn) AcceptorClient {
	return &yarpcAcceptorClient{
		cc: cc,
	}
}

type yarpcAcceptorClient struct {
	cc *yarpc.ClientConn
}

func (c *yarpcAcceptorClient) Prepare(ctx context.Context, request *Request) (*Promise, error) {
	stream, err := c.cc.OpenStream(ctx, "/paxos.Acceptor/Prepare")
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
	stream, err := c.cc.OpenStream(ctx, "/paxos.Acceptor/Accept")
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

func (c *yarpcAcceptorClient) Observe(ctx context.Context, request *Request) (*ObserveClientStream, error) {
	stream, err := c.cc.OpenStream(ctx, "/paxos.Acceptor/Observe")
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

var _ AcceptorClient = &yarpcAcceptorClient{}

func RegisterYarpcProposerServer(svr *yarpc.Server, impl ProposerServer) {
	svr.Handle("/paxos.Proposer/Propose", yarpc.HandlerFunc(func(stream yarpc.Stream) error {
		req := &Bytes{}
		err := stream.ReadMsg(req)
		if err != nil {
			return err
		}

		bytes, err := impl.Propose(stream.Context(), req.Bytes)
		if err != nil {
			return err
		}

		return stream.WriteMsg(&Bytes{
			Bytes: bytes,
		})
	}))
}

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
		Bytes: value,
	})
	if err != nil {
		return nil, err
	}

	bytes := &Bytes{}
	err = stream.ReadMsg(bytes)
	if err != nil {
		return nil, err
	}

	return bytes.Bytes, nil
}

var _ ProposerClient = &yarpcProposerClient{}
