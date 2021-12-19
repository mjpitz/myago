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

package plugin

import (
	"context"
	"io"
	"os"

	"github.com/mjpitz/myago/yarpc"
)

// Listen returns a yarpc.Listener that treats a processes stdin and stdout as a connection.
func Listen() yarpc.Listener {
	return listen(&serverRWC{
		stdin:  os.Stdin,
		stdout: os.Stdout,
	})
}

func listen(conn io.ReadWriteCloser) yarpc.Listener {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan io.ReadWriteCloser, 1)
	ch <- conn

	return &listener{
		context: ctx,
		cancel:  cancel,
		conn:    ch,
	}
}

// listener implements a single use yarpc.Listener that will only ever return one io.ReadWriteCloser.
type listener struct {
	context context.Context
	cancel  context.CancelFunc
	conn    <-chan io.ReadWriteCloser
}

func (l *listener) Accept() (io.ReadWriteCloser, error) {
	select {
	case <-l.context.Done():
		return nil, l.context.Err()
	case conn := <-l.conn:
		return conn, nil
	}
}

func (l *listener) Close() error {
	l.cancel()
	return nil
}

var _ yarpc.Listener = &listener{}

// serverRWC implements an io.ReadWriteCloser that mimics a net.Conn using a process's stdin and stdout.
type serverRWC struct {
	stdin  io.ReadCloser
	stdout io.WriteCloser
}

func (s *serverRWC) Read(p []byte) (n int, err error) {
	return s.stdin.Read(p)
}

func (s *serverRWC) Write(p []byte) (n int, err error) {
	return s.stdout.Write(p)
}

func (s *serverRWC) Close() error {
	stdout, ok := s.stdout.(*os.File)
	if ok {
		// force stdout to sync before closing
		_ = stdout.Sync()
	}

	_ = s.stdin.Close()
	_ = s.stdout.Close()
	return nil
}

var _ io.ReadWriteCloser = &serverRWC{}
