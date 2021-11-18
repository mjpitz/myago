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
	"net"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/yamux"
	"github.com/jonboulle/clockwork"

	"github.com/mjpitz/myago/encoding"
)

// Dialer provides a common interface for obtaining a net.Conn. This makes it easy to handle TLS transparently.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// DialContext initializes a new client connection to the target server.
func DialContext(ctx context.Context, network, target string, opts ...Option) *ClientConn {
	o := &options{
		context:  ctx,
		yamux:    yamux.DefaultConfig(),
		encoding: encoding.MsgPack,
		clock:    clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		opt(o)
	}

	var dialer Dialer = &net.Dialer{}
	if o.tls != nil {
		dialer = &tls.Dialer{
			NetDialer: &net.Dialer{},
			Config:    o.tls,
		}
	}

	return &ClientConn{
		dialer:  dialer,
		network: network,
		target:  target,
		options: o,
		mu:      sync.Mutex{},
	}
}

// ClientConn.
type ClientConn struct {
	dialer  Dialer
	network string
	target  string

	options *options
	mu      sync.Mutex
	session *yamux.Session
}

func (c *ClientConn) obtainSession(ctx context.Context) (*yamux.Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session == nil || c.session.IsClosed() {
		backoffConfig := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)

		err := backoff.Retry(
			func() error {
				conn, err := c.dialer.DialContext(ctx, c.network, c.target)
				if err != nil {
					return err
				}

				c.session, err = yamux.Client(conn, c.options.yamux)
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

// OpenStream starts a stream for a given RPC.
func (c *ClientConn) OpenStream(ctx context.Context, method string) (Stream, error) {
	session, err := c.obtainSession(ctx)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}

	rpcStream := wrap(stream,
		withContext(c.options.context),
		withEncoding(c.options.encoding))

	err = rpcStream.WriteMsg(&Invoke{
		Method: method,
	})

	return rpcStream, err
}
