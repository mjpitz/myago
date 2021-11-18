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

// Handler defines an interface that can be used for handling requests.
type Handler interface {
	Handle(Stream) error
}

// HandlerFunc provides users with a simple functional interface for a Handler.
type HandlerFunc func(Stream) error

func (fn HandlerFunc) Handle(stream Stream) error {
	if fn == nil {
		return nil
	}

	return fn(stream)
}

// DefaultServer is a global server definition that can be leveraged by hosting program.
var DefaultServer = &Server{}

// Handle adds the provided handler to the default server.
func Handle(pattern string, handler Handler) {
	DefaultServer.Handle(pattern, handler)
}

// HandleFunc adds the provided handler function to the default server.
func HandleFunc(pattern string, handler func(Stream) error) {
	Handle(pattern, HandlerFunc(handler))
}

// ListenAndServe starts the default server on the provided network and address.
func ListenAndServe(network, address string, opts ...Option) error {
	return DefaultServer.ListenAndServe(network, address, opts...)
}
