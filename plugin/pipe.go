package plugin

import (
	"bytes"
	"fmt"
	"io"
)

// Pipe returns a psuedo-async io.ReadWriteCloser.
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
