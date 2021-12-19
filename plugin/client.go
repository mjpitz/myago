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
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/mjpitz/myago/yarpc"
)

// DialContext returns a ClientConn whose dialer forks a process for the specified binary.
func DialContext(ctx context.Context, binary string, args ...string) *yarpc.ClientConn {
	clientConn := yarpc.NewClientConn(ctx)
	clientConn.Dialer = &dialer{
		Binary: binary,
		Args:   args,
	}

	return clientConn
}

// dialer provides a yarpc.Dialer implementation that forks the specified binary and enabled bidirectional communication
// via its stdin and stdout.
type dialer struct {
	Binary string
	Args   []string
}

func (d *dialer) DialContext(ctx context.Context) (io.ReadWriteCloser, error) {
	stdin := NewBlockingReadWriteCloser()
	stdout := NewBlockingReadWriteCloser()
	stderr := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, d.Binary, d.Args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("%v\n%s", err, stderr.String())
	}

	return &clientRWC{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		cmd:    cmd,
	}, nil
}

var _ yarpc.Dialer = &dialer{}

// clientRWC implements an io.ReadWriteCloser that mimics a net.Conn over a forked programs stdin and stdout.
type clientRWC struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr *bytes.Buffer
	cmd    *exec.Cmd
}

func (c *clientRWC) Read(p []byte) (n int, err error) {
	return c.stdout.Read(p)
}

func (c *clientRWC) Write(p []byte) (n int, err error) {
	return c.stdin.Write(p)
}

func (c *clientRWC) Close() (err error) {
	_ = c.stdin.Close()
	err = c.cmd.Wait()
	_ = c.stdout.Close()

	if err != nil {
		return fmt.Errorf("%v\n%s", err, c.stderr.String())
	}

	return nil
}

var _ io.ReadWriteCloser = &clientRWC{}
