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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/mjpitz/myago/yarpc"
	"github.com/mjpitz/myago/zaputil"
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
	stdin := Pipe()
	stdout := Pipe()
	stderr := Pipe()

	// nolint:gosec
	cmd := exec.CommandContext(ctx, d.Binary, d.Args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	rwc := &clientRWC{
		stdin:  stdin,
		stdout: stdout,
		cmd:    cmd,
		err:    make(chan error, 1),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to locate plugin: %s", d.Binary)
	}

	go func() {
		_ = cmd.Wait()
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		var err error
		for err == nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			case <-ticker.C:
				if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
					err = errors.New("plugin exited")
				}
				ticker.Reset(time.Second)
			}
		}

		// close the pipe to allow io.EOF to be returned
		_ = stderr.Close()
		rwc.err <- err
		_ = rwc.Close()

		body, _ := ioutil.ReadAll(stderr)
		if len(body) > 0 {
			zaputil.Extract(ctx).Error("plugin execution failed",
				zap.String("plugin", d.Binary),
				zap.String("error", strings.TrimSpace(string(body))))
		}
	}()

	return rwc, nil
}

var _ yarpc.Dialer = &dialer{}

// clientRWC implements an io.ReadWriteCloser that mimics a net.Conn over a forked programs stdin and stdout.
type clientRWC struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser

	cmd *exec.Cmd
	err chan error
}

func (c *clientRWC) readError() error {
	err := <-c.err
	c.err <- err

	return err
}

func (c *clientRWC) Read(p []byte) (n int, err error) {
	n, err = c.stdout.Read(p)

	if err != nil {
		read := c.readError()
		if read != nil {
			err = read
		}
	}

	return
}

func (c *clientRWC) Write(p []byte) (n int, err error) {
	n, err = c.stdin.Write(p)

	if err != nil {
		read := c.readError()
		if read != nil {
			err = read
		}
	}

	return
}

func (c *clientRWC) Close() (err error) {
	_ = c.stdin.Close()

	if c.cmd != nil {
		_ = c.cmd.Wait()
	}

	_ = c.stdout.Close()

	return nil
}

var _ io.ReadWriteCloser = &clientRWC{}
