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
	"sync/atomic"

	"github.com/hashicorp/yamux"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"

	"github.com/mjpitz/myago/encoding"
)

const (
	uninitializedState int32 = iota
	stoppedState
	startedState
)

type Listener interface {
	Accept() (io.ReadWriteCloser, error)
	Close() error
}

type Server struct {
	Handler Handler

	// set during initialization
	state int32

	// set during serve
	options  options
	listener Listener
	pool     *ants.Pool
}

func (s *Server) init() {
	if !atomic.CompareAndSwapInt32(&s.state, uninitializedState, stoppedState) {
		return
	}
}

func (s *Server) handleStream(stream *yamux.Stream) func() {
	return func() {
		var err error

		rpcStream := Wrap(stream, withOptions(&s.options))

		defer func() {
			if err != nil {
				// log
				if err := stream.Close(); err != nil {
					// log
				}
			}
		}()

		err = s.Handler.ServeYARPC(rpcStream)
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
	s.init()

	if !atomic.CompareAndSwapInt32(&s.state, startedState, stoppedState) {
		return errors.New("server not started")
	}

	defer func() {
		s.pool.Release()
		s.listener = nil
	}()

	if listener := s.listener; listener != nil {
		return listener.Close()
	}

	return nil
}

func (s *Server) Serve(listener Listener, opts ...Option) error {
	s.init()

	if !atomic.CompareAndSwapInt32(&s.state, stoppedState, startedState) {
		return errors.New("server already started")
	}

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

	s.options = o
	s.listener = listener

	if s.pool == nil {
		pool, err := ants.NewPool(3000, ants.WithOptions(ants.Options{
			// Logger: zap.NewStdLog(logger),
		}))
		if err != nil {
			return errors.Wrap(err, "failed to construct ant pool")
		}

		s.pool = pool
	} else {
		s.pool.Reboot()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "failed to accept connection")
		}

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
