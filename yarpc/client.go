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

package yarpc

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/yamux"

	"github.com/mjpitz/myago/encoding"
	"github.com/mjpitz/myago/zaputil"
)

// DialContext initializes a new client connection to the target server.
func DialContext(ctx context.Context, network, target string, opts ...Option) *ClientConn {
	c := NewClientConn(ctx).WithOptions(opts...)

	dialer := &NetDialerAdapter{
		Dialer:  &net.Dialer{},
		Network: network,
		Target:  target,
	}

	if c.options.tls != nil {
		dialer.Dialer = &tls.Dialer{
			NetDialer: &net.Dialer{},
			Config:    c.options.tls,
		}
	}

	c.Dialer = dialer
	return c
}

// NewClientConn creates a default ClientConn with an empty dialer implementation. The Dialer must be configured before
// use. This function is intended to be used in initializer functions such as DialContext.
func NewClientConn(ctx context.Context) *ClientConn {
	return &ClientConn{
		Dialer: &emptyDialer{},
		options: options{
			context:  ctx,
			yamux:    yamux.DefaultConfig(),
			encoding: encoding.MsgPack,
		},
		mu: sync.Mutex{},
	}
}

// ClientConn defines an abstract connection yarpc clients to use.
type ClientConn struct {
	Dialer  Dialer
	options options
	mu      sync.Mutex
	session *yamux.Session
}

// WithOptions configures the options for the underlying client connection.
func (c *ClientConn) WithOptions(opts ...Option) *ClientConn {
	for _, opt := range opts {
		opt(&(c.options))
	}

	return c
}

func (c *ClientConn) obtainSession(ctx context.Context) (*yamux.Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session == nil || c.session.IsClosed() {
		backoffConfig := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)

		err := backoff.Retry(
			func() error {
				conn, err := c.Dialer.DialContext(ctx)
				if err != nil {
					return err
				}

				yamuxcfg := *c.options.yamux
				yamuxcfg.Logger = zaputil.HashiCorpStdLogger(zaputil.Extract(ctx))
				yamuxcfg.LogOutput = nil

				c.session, err = yamux.Client(conn, &yamuxcfg)
				if err != nil {
					return err
				}

				return nil
			},
			backoffConfig,
		)
		if err != nil {
			return nil, err
		}
	}

	return c.session, nil
}

// OpenStream starts a stream for the named RPC.
func (c *ClientConn) OpenStream(ctx context.Context, method string) (Stream, error) {
	session, err := c.obtainSession(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}

	rpcStream := Wrap(stream, withOptions(&c.options))

	err = rpcStream.WriteMsg(&Invoke{
		Method: method,
	})

	return rpcStream, err
}

// Dialer provides a minimal interface needed to establish a client.
type Dialer interface {
	DialContext(ctx context.Context) (io.ReadWriteCloser, error)
}

// emptyDialer is used to signal that the Dialer implementation needs to be provided on the ClientConn.
type emptyDialer struct{}

func (d *emptyDialer) DialContext(ctx context.Context) (io.ReadWriteCloser, error) {
	panic("Dialer not provided")
}
