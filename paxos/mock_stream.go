package paxos

import (
	"context"
	"reflect"
	"time"

	"github.com/mjpitz/myago/yarpc"
)

func NewMockStream(size int) *MockStream {
	return &MockStream{
		Ctx:      context.Background(),
		Incoming: make(chan interface{}, size),
		Outgoing: make(chan interface{}, size),
	}
}

type MockStream struct {
	Ctx      context.Context
	Incoming chan interface{}
	Outgoing chan interface{}
}

func (m *MockStream) Context() context.Context {
	return m.Ctx
}

func (m *MockStream) SetReadDeadline(deadline time.Time) error {
	return nil
}

func (m *MockStream) ReadMsg(i interface{}) error {
	msg := <-m.Incoming
	reflect.Indirect(reflect.ValueOf(i)).Set(reflect.Indirect(reflect.ValueOf(msg)))
	return nil
}

func (m *MockStream) SetWriteDeadline(deadline time.Time) error {
	return nil
}

func (m *MockStream) WriteMsg(i interface{}) error {
	m.Outgoing <- i
	return nil
}

func (m *MockStream) Close() error {
	return nil
}

var _ yarpc.Stream = &MockStream{}
