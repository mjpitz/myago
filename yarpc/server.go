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
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"

	"go.pitz.tech/lib/encoding"
	"go.pitz.tech/lib/logger"
)

type Listener interface {
	Accept() (io.ReadWriteCloser, error)
	Close() error
}

type Server struct {
	Handler Handler

	// set during serve
	mu       sync.Mutex
	options  options
	listener Listener
	pool     *ants.Pool
}

func (s *Server) handleStream(stream *yamux.Stream) func() {
	return func() {
		rpcStream := Wrap(stream, withOptions(&s.options))
		var err error

		defer func() {
			if err != nil {
				// log
			}

			if err := stream.Close(); err != nil {
				// log
			}
		}()

		err = s.Handler.ServeYARPC(rpcStream)
		if err != nil {
			err = rpcStream.WriteMsg(&Status{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}
}

func (s *Server) handleSession(session *yamux.Session) func() {
	return func() {
		var err error

		defer func() {
			if err != nil {
				// log
			}

			if err = session.Close(); err != nil {
				// log
			}
		}()

		var stream *yamux.Stream

		for {
			stream, err = session.AcceptStream()
			if err != nil {
				return
			}

			err = s.pool.Submit(s.handleStream(stream))
			if err != nil {
				// log
				_ = stream.Close()
			}
		}
	}
}

/* all public functions must start with s.once.Do(s.init) */

func (s *Server) Shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if listener := s.listener; listener != nil {
		s.listener = nil
		_ = listener.Close()
	}

	if s.pool != nil {
		s.pool.Release()
	}

	return nil
}

func (s *Server) Serve(listener Listener, opts ...Option) error {
	o := options{
		context:  context.Background(),
		yamux:    yamux.DefaultConfig(),
		encoding: encoding.MsgPack,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.tls != nil {
		ntl, ok := listener.(*NetListenerAdapter)
		if !ok {
			return fmt.Errorf("tls not supported on non-net listeners")
		}

		listener = &NetListenerAdapter{
			Listener: tls.NewListener(ntl.Listener, o.tls),
		}
	}

	err := func() (err error) {
		s.mu.Lock()
		defer s.mu.Unlock()

		pool := s.pool
		if pool == nil {
			pool, err = ants.NewPool(3000, ants.WithOptions(ants.Options{
				// Logger: zap.NewStdLog(logger),
			}))
			if err != nil {
				return errors.Wrap(err, "failed to construct ant pool")
			}
		} else {
			pool.Reboot()
		}

		s.options = o
		s.listener = listener
		s.pool = pool

		return nil
	}()
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "failed to accept connection")
		}

		yamuxcfg := *o.yamux
		yamuxcfg.Logger = logger.HashiCorpStdLogger(logger.Extract(o.context))
		yamuxcfg.LogOutput = nil

		session, err := yamux.Server(conn, o.yamux)
		if err != nil {
			return errors.Wrap(err, "failed to Wrap connection for yamux")
		}

		err = s.pool.Submit(s.handleSession(session))
		if err != nil {
			_ = session.Close()
		}
	}
}

func (s *Server) ListenAndServe(network, address string, opts ...Option) error {
	netListener, err := net.Listen(network, address)
	if err != nil {
		return errors.Wrap(err, "failed to bind to network")
	}

	listener := &NetListenerAdapter{
		Listener: netListener,
	}

	return s.Serve(listener, opts...)
}
