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
	"fmt"
	"io"
)

// Pipe returns a pseudo-async io.ReadWriteCloser.
//
//nolint:revive
func Pipe() *pipe {
	done := make(chan bool, 1)
	done <- false

	data := make(chan *bytes.Buffer, 1)
	data <- bytes.NewBuffer(nil)

	return &pipe{
		done: done,
		data: data,
	}
}

type pipe struct {
	done chan bool
	data chan *bytes.Buffer
}

func (pip *pipe) Len() int {
	data := <-pip.data
	defer func() { pip.data <- data }()

	return data.Len()
}

func (pip *pipe) Read(p []byte) (n int, err error) {
	data := <-pip.data
	done := <-pip.done

	defer func() {
		pip.done <- done
		pip.data <- data
	}()

	n, _ = data.Read(p)

	if done && data.Len() == 0 {
		err = io.EOF
	}

	return
}

func (pip *pipe) Write(p []byte) (n int, err error) {
	data := <-pip.data
	done := <-pip.done

	defer func() {
		pip.done <- done
		pip.data <- data
	}()

	if done {
		err = io.ErrClosedPipe

		return
	}

	return data.Write(p)
}

func (pip *pipe) Close() error {
	data := <-pip.data
	done := <-pip.done

	defer func() {
		pip.done <- done
		pip.data <- data
	}()

	if done {
		return fmt.Errorf("pipe already closed")
	}

	done = true

	return nil
}

var _ io.ReadWriteCloser = &pipe{}
