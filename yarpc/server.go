package yarpc

import (
	"context"
	"crypto/tls"
	"net"
	"sync/atomic"

	"github.com/hashicorp/yamux"
	"github.com/jonboulle/clockwork"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
)

const (
	uninitializedState int32 = iota
	stoppedState
	startedState
)

type Server struct {
	// set during initialization
	state    int32
	handlers map[string]Handler

	// set during serve
	options  *options
	listener net.Listener
	pool     *ants.Pool
}

func (s *Server) init() {
	if !atomic.CompareAndSwapInt32(&s.state, uninitializedState, stoppedState) {
		return
	}

	s.handlers = make(map[string]Handler)
}

func (s *Server) handleStream(stream *yamux.Stream) func() {
	return func() {
		var err error

		rpcStream := wrap(stream,
			withContext(s.options.context),
			withEncoding(s.options.encoding))

		defer func() {
			if err != nil {
				// log

				err := stream.Close()
				if err != nil {
					// log
				}
			}
		}()

		invoke := &Invoke{}
		err = rpcStream.ReadMsg(invoke)
		if err != nil {
			return
		}

		handler := s.handlers[invoke.Method]
		if handler == nil {
			return
		}

		err = handler.Handle(rpcStream)
	}
}

func (s *Server) handleSession(session *yamux.Session) func() {
	return func() {
		var err error

		defer func() {
			if err != nil {
				// log
			}

			err = session.Close()
			if err != nil {
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
// todo: figure out how to support opening a stream from the server
//   need to track sessions by host?

func (s *Server) Handle(pattern string, handler Handler) {
	s.init()

	s.handlers[pattern] = handler
}

func (s *Server) Shutdown() error {
	s.init()

	if !atomic.CompareAndSwapInt32(&s.state, startedState, stoppedState) {
		return errors.New("server not started")
	}

	defer s.pool.Release()

	listener := s.listener
	s.listener = nil

	if listener != nil {
		return listener.Close()
	}

	return nil
}

func (s *Server) Serve(listener net.Listener, opts ...Option) error {
	s.init()

	if !atomic.CompareAndSwapInt32(&s.state, stoppedState, startedState) {
		return errors.New("server already started")
	}

	o := &options{
		context:  context.Background(),
		yamux:    yamux.DefaultConfig(),
		encoding: &MSGPackEncoding{},
		clock:    clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.tls != nil {
		listener = tls.NewListener(listener, o.tls)
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
			return errors.Wrap(err, "failed to wrap connection for yamux")
		}

		err = s.pool.Submit(s.handleSession(session))
		if err != nil {
			_ = session.Close()
		}
	}
}

func (s *Server) ListenAndServe(network, address string, opts ...Option) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return errors.Wrap(err, "failed to bind to network")
	}

	return s.Serve(listener, opts...)
}
