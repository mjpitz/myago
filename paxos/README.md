# paxos
--
    import "github.com/mjpitz/myago/paxos"

Package paxos implements the paxos algorithm. The logic is mostly ported from
mjpitz/paxos, but with a few modifications. First, I didn't continue using gRPC
as the transport as I wanted something a bit less cumbersome. I've tried to
break down the interface in such a way where different transports _could_ be
plugged in. For simplicity, this

## Usage

#### func  RegisterYarpcAcceptorServer

```go
func RegisterYarpcAcceptorServer(svr *yarpc.Server, impl AcceptorServer)
```

#### type Acceptor

```go
type Acceptor struct {
}
```


#### func  NewAcceptor

```go
func NewAcceptor(promiseLog, acceptLog Log) (*Acceptor, error)
```

#### func (*Acceptor) Accept

```go
func (a *Acceptor) Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
```

#### func (*Acceptor) Observe

```go
func (a *Acceptor) Observe(call *ObserveServerStream) error
```

#### func (*Acceptor) Prepare

```go
func (a *Acceptor) Prepare(ctx context.Context, req *Request) (*Promise, error)
```

#### type AcceptorClient

```go
type AcceptorClient interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
	Observe(ctx context.Context, request *Request) (*ObserveClientStream, error)
}
```


#### func  NewYarpcAcceptorClient

```go
func NewYarpcAcceptorClient(cc *yarpc.ClientConn) AcceptorClient
```

#### type AcceptorServer

```go
type AcceptorServer interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
	Observe(call *ObserveServerStream) error
}
```


#### type IDGenerator

```go
type IDGenerator interface {
	Next() (uint64, error)
}
```

IDGenerator defines an interface for generating IDs used internally by paxos.

#### func  ServerIDGenerator

```go
func ServerIDGenerator(serverID uint8, clock clockwork.Clock) IDGenerator
```
ServerIDGenerator returns an IDGenerator that creates an ID using a provided
serverID and clock. It works by by taking a millisecond level timestamp,
shifting it's value left 8 bits, and or'ing it with the server ID. The leading
byte can be used to expand this representation later on.

    	const (
    		wordView  = 0x0000000000000000

    		nowMillis = 0x0000017c96370c09
         shifted   = 0x00017c96370c0900
         withSID   = 0x00017c96370c09XX
    	)

As you can see, there is plenty of space for the IDGenerator to function.
Obviously, there are limitations with this implementation.

    1. 256 max possible instances
    1. Throughput constrained to 1 op/ms

Granted, some of these aren't _huge_ issues for the types of systems that this
_could_ help build.

#### type Log

```go
type Log interface {
	Record(id uint64, msg interface{}) error
	Last(msg interface{}) error
	Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error
}
```


#### type MemoryLog

```go
type MemoryLog struct {
}
```


#### func (*MemoryLog) Last

```go
func (m *MemoryLog) Last(msg interface{}) error
```

#### func (*MemoryLog) Range

```go
func (m *MemoryLog) Range(start, end uint64, proto interface{}, fn func(msg interface{}) error) error
```

#### func (*MemoryLog) Record

```go
func (m *MemoryLog) Record(id uint64, msg interface{}) error
```

#### type ObserveClientStream

```go
type ObserveClientStream struct {
	Stream
}
```


#### func (*ObserveClientStream) Recv

```go
func (s *ObserveClientStream) Recv() (*Proposal, error)
```

#### type ObserveServerStream

```go
type ObserveServerStream struct {
	Stream
}
```


#### func (*ObserveServerStream) Recv

```go
func (s *ObserveServerStream) Recv() (*Request, error)
```

#### func (*ObserveServerStream) Send

```go
func (s *ObserveServerStream) Send(msg *Proposal) error
```

#### type Promise

```go
type Promise struct {
	ID       uint64    `json:"id,omitempty"`
	Accepted *Proposal `json:"accepted,omitempty"`
}
```


#### type Proposal

```go
type Proposal struct {
	ID    uint64 `json:"id,omitempty"`
	Value []byte `json:"value,omitempty"`
}
```


#### type Proposer

```go
type Proposer struct {
	Clock       clockwork.Clock
	IDGenerator IDGenerator
	Acceptor    AcceptorClient
}
```


#### func (*Proposer) Propose

```go
func (p *Proposer) Propose(ctx context.Context, value []byte) (accepted []byte, err error)
```

#### type Request

```go
type Request struct {
	ID      uint64 `json:"id,omitempty"`
	Attempt uint64 `json:"attempt,omitempty"`
}
```


#### type Stream

```go
type Stream interface {
	Context() context.Context
	SetReadDeadline(deadline time.Time) error
	ReadMsg(i interface{}) error
	SetWriteDeadline(deadline time.Time) error
	WriteMsg(i interface{}) error
	Close() error
}
```
