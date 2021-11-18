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

package paxos

import (
	"context"
	"reflect"
	"time"

	"github.com/mjpitz/myago/yarpc"
)

// Stream provides an abstract definition of the functionality the underlying stream needs to provide.
type Stream interface {
	Context() context.Context
	SetReadDeadline(deadline time.Time) error
	ReadMsg(i interface{}) error
	SetWriteDeadline(deadline time.Time) error
	WriteMsg(i interface{}) error
	Close() error
}

// NewMockStream provides a mock Stream implementation useful for testing. This could be yarpc or paxos related.
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
