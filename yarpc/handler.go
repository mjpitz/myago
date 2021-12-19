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
	"fmt"
	"sync"
)

var (
	// DefaultServeMux provides a default request multiplexer (router).
	DefaultServeMux = &ServeMux{}

	// DefaultServer is a global server definition that can be leveraged by hosting program.
	DefaultServer = &Server{
		Handler: DefaultServeMux,
	}
)

// Handle adds the provided handler to the default server.
func Handle(pattern string, handler Handler) {
	DefaultServeMux.Handle(pattern, handler)
}

// HandleFunc adds the provided handler function to the default server.
func HandleFunc(pattern string, handler func(Stream) error) {
	DefaultServeMux.Handle(pattern, HandlerFunc(handler))
}

// Handler defines an interface that can be used for handling requests.
type Handler interface {
	ServeYARPC(Stream) error
}

// HandlerFunc provides users with a simple functional interface for a Handler.
type HandlerFunc func(Stream) error

func (fn HandlerFunc) ServeYARPC(stream Stream) error {
	if fn == nil {
		return nil
	}

	return fn(stream)
}

// ServeMux provides a router implementation for yarpc calls.
type ServeMux struct {
	once     sync.Once
	handlers map[string]Handler
}

func (s *ServeMux) init() {
	s.once.Do(func() {
		s.handlers = make(map[string]Handler)
	})
}

func (s *ServeMux) Handle(pattern string, handler Handler) {
	s.init()

	s.handlers[pattern] = handler
}

func (s *ServeMux) ServeYARPC(stream Stream) (err error) {
	s.init()

	invoke := &Invoke{}
	err = stream.ReadMsg(invoke)
	if err != nil {
		return
	}

	handler := s.handlers[invoke.Method]
	if handler == nil {
		err = fmt.Errorf("handler not found")
		return
	}

	return handler.ServeYARPC(stream)
}

var _ Handler = &ServeMux{}

// ListenAndServe starts the default server on the provided network and address.
func ListenAndServe(network, address string, opts ...Option) error {
	return DefaultServer.ListenAndServe(network, address, opts...)
}

// Serve starts the default server using the provided listener.
func Serve(listener Listener, opts ...Option) error {
	return DefaultServer.Serve(listener, opts...)
}
