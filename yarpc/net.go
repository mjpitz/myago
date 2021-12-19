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
	"io"
	"net"
)

// NetDialer provides a common interface for obtaining a net.Conn. This makes it easy to handle TLS transparently.
type NetDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// NetDialerAdapter adapts the provided NetDialer to support io.ReadWriteCloser.
type NetDialerAdapter struct {
	Dialer  NetDialer
	Network string
	Target  string
}

// DialContext returns a creates a new network connection.
func (a *NetDialerAdapter) DialContext(ctx context.Context) (io.ReadWriteCloser, error) {
	return a.Dialer.DialContext(ctx, a.Network, a.Target)
}

var _ Dialer = &NetDialerAdapter{}

// NetListenerAdapter adapts the provided net.Listener to support io.ReadWriteCloser.
type NetListenerAdapter struct {
	Listener net.Listener
}

func (n *NetListenerAdapter) Accept() (io.ReadWriteCloser, error) {
	return n.Listener.Accept()
}

func (n *NetListenerAdapter) Close() error {
	return n.Listener.Close()
}

var _ Listener = &NetListenerAdapter{}
