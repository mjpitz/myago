# paxos
--
    import "github.com/mjpitz/myago/paxos"

Package paxos implements the paxos algorithm. The logic is mostly ported from
mjpitz/paxos, but with a few modifications. First, I didn't continue using gRPC
as the transport as I wanted something a bit less cumbersome. I've tried to
break down the interface in such a way where different transports _could_ be
plugged in. More on that later.

This package is (and likely will be for a while) a work in progress. As it
stands, it _should_ support simple paxos.

## Usage

#### func  RegisterYarpcAcceptorServer

```go
func RegisterYarpcAcceptorServer(svr *yarpc.Server, impl AcceptorServer)
```
RegisterYarpcAcceptorServer registers the provided AcceptorServer implementation
with the yarpc.Server to handle requests.

#### func  RegisterYarpcObserverServer

```go
func RegisterYarpcObserverServer(svr *yarpc.Server, impl ObserverServer)
```
RegisterYarpcObserverServer registers the provided ObserverServer implementation
with the yarpc.Server to handle requests. Acceptors should implement the
observer server, otherwise other members of the cluster cannot determine what
records have been accepted.

#### func  RegisterYarpcProposerServer

```go
func RegisterYarpcProposerServer(svr *yarpc.Server, impl ProposerServer)
```
RegisterYarpcProposerServer registers the provided ProposerServer implementation
with the yarpc.Server to handle requests. Typically, proposers aren't embedded
as a server and are instead run as client side code.

#### type Acceptor

```go
type Acceptor interface {
	AcceptorServer
	ObserverServer
}
```


#### func  NewAcceptor

```go
func NewAcceptor(promiseLog, acceptedLog Log) (Acceptor, error)
```

#### type AcceptorClient

```go
type AcceptorClient interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
}
```


#### func  NewYarpcAcceptorClient

```go
func NewYarpcAcceptorClient(cc *yarpc.ClientConn) AcceptorClient
```
NewYarpcAcceptorClient wraps the provided yarpc.ClientConn with an
AcceptorClient implementation.

#### type AcceptorServer

```go
type AcceptorServer interface {
	Prepare(ctx context.Context, request *Request) (*Promise, error)
	Accept(ctx context.Context, proposal *Proposal) (*Proposal, error)
}
```


#### type Badger

```go
type Badger struct {
	DB *badger.DB
}
```

Badger implements a Log that wraps an underlying badgerdb instance.

#### func (*Badger) Last

```go
func (l *Badger) Last(msg interface{}) error
```

#### func (*Badger) Range

```go
func (l *Badger) Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error
```

#### func (*Badger) Record

```go
func (l *Badger) Record(id uint64, msg interface{}) error
```

#### func (*Badger) WithPrefix

```go
func (l *Badger) WithPrefix(prefix string) Log
```

#### type Bytes

```go
type Bytes struct {
	Value []byte `json:"value,omitempty"`
}
```

Bytes contains a value to be accepted via paxos.

#### type Config

```go
type Config struct {
	Clock          clockwork.Clock
	IDGenerator    IDGenerator
	PromiseLog     Log
	AcceptedLog    Log
	RecordedLog    Log
	AcceptorDialer func(ctx context.Context, member string) (AcceptorClient, error)
	ObserverDialer func(ctx context.Context, member string) (ObserverClient, error)
}
```

Config contains configurable elements of Paxos.

#### func (*Config) Validate

```go
func (c *Config) Validate() error
```
Validate ensures the configuration is valid.

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
	WithPrefix(str string) Log
	Record(id uint64, msg interface{}) error
	Last(msg interface{}) error
	Range(start, stop uint64, proto interface{}, fn func(msg interface{}) error) error
}
```


#### type Memory

```go
type Memory struct {
}
```


#### func (*Memory) Last

```go
func (m *Memory) Last(msg interface{}) error
```

#### func (*Memory) Range

```go
func (m *Memory) Range(start, end uint64, proto interface{}, fn func(msg interface{}) error) error
```

#### func (*Memory) Record

```go
func (m *Memory) Record(id uint64, msg interface{}) error
```

#### func (*Memory) WithPrefix

```go
func (m *Memory) WithPrefix(prefix string) Log
```

#### type MockStream

```go
type MockStream struct {
	Ctx      context.Context
	Incoming chan interface{}
	Outgoing chan interface{}
}
```


#### func  NewMockStream

```go
func NewMockStream(size int) *MockStream
```
NewMockStream provides a mock Stream implementation useful for testing. This
could be yarpc or paxos related.

#### func (*MockStream) Close

```go
func (m *MockStream) Close() error
```

#### func (*MockStream) Context

```go
func (m *MockStream) Context() context.Context
```

#### func (*MockStream) ReadMsg

```go
func (m *MockStream) ReadMsg(i interface{}) error
```

#### func (*MockStream) SetReadDeadline

```go
func (m *MockStream) SetReadDeadline(deadline time.Time) error
```

#### func (*MockStream) SetWriteDeadline

```go
func (m *MockStream) SetWriteDeadline(deadline time.Time) error
```

#### func (*MockStream) WriteMsg

```go
func (m *MockStream) WriteMsg(i interface{}) error
```

#### type MultiAcceptorClient

```go
type MultiAcceptorClient struct {
	Dialer func(ctx context.Context, member string) (AcceptorClient, error)
}
```


#### func (*MultiAcceptorClient) Accept

```go
func (m *MultiAcceptorClient) Accept(ctx context.Context, in *Proposal) (*Proposal, error)
```

#### func (*MultiAcceptorClient) Prepare

```go
func (m *MultiAcceptorClient) Prepare(ctx context.Context, request *Request) (*Promise, error)
```
nolint:cyclop

#### func (*MultiAcceptorClient) Start

```go
func (m *MultiAcceptorClient) Start(ctx context.Context, membership *cluster.Membership) error
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

#### type Observer

```go
type Observer struct {
	Dialer func(ctx context.Context, member string) (ObserverClient, error)
	Log    Log
}
```

Observer watches the Acceptors to learn about what values have been accepted.

#### func (*Observer) Start

```go
func (o *Observer) Start(ctx context.Context, membership *cluster.Membership) error
```
nolint:gocognit,cyclop

#### type ObserverClient

```go
type ObserverClient interface {
	Observe(ctx context.Context, request *Request) (*ObserveClientStream, error)
}
```


#### func  NewYarpcObserverClient

```go
func NewYarpcObserverClient(cc *yarpc.ClientConn) ObserverClient
```
NewYarpcObserverClient wraps the provided yarpc.ClientConn with an
ObserverClient implementation.

#### type ObserverServer

```go
type ObserverServer interface {
	Observe(call *ObserveServerStream) error
}
```


#### type Paxos

```go
type Paxos struct {
	// Proposer contains the logic required to propose changes to the paxos state machine. Any member in paxos can act
	// as a proposer. Proposers communicate with all acceptor to propose changes to the log.
	Proposer

	// Observer contains the logic required to be an observer of the paxos protocol. Every member in paxos _must_ be an
	// observer. Observers watch all acceptor to learn about the records they've accepted.
	Observer

	// Acceptor must implement the functionality of an AcceptorServer and an ObserverServer. The ObserverServer is how
	// other members of the cluster learn about changes.
	Acceptor
}
```

Paxos defines the core elements of a paxos participant.

#### func  New

```go
func New(cfg *Config) (*Paxos, error)
```
New constructs a new instance of paxos given the provided configuration. It
returns an error should the provided configuration be invalid.

#### func (*Paxos) Start

```go
func (p *Paxos) Start(ctx context.Context, membership *cluster.Membership) error
```

#### type Promise

```go
type Promise struct {
	ID       uint64    `json:"id,omitempty"`
	Accepted *Proposal `json:"accepted,omitempty"`
}
```

Promise is returned by an accepted prepare. If more than one attempt was made,
and accepted value is returned with the last accepted proposal so clients can
catch up.

#### type Proposal

```go
type Proposal struct {
	ID    uint64 `json:"id,omitempty"`
	Value []byte `json:"value,omitempty"`
}
```

Proposal is used to propose a log value to system.

#### type Proposer

```go
type Proposer struct {
	IDGenerator IDGenerator
	Acceptor    AcceptorClient
}
```

Proposer can be run either as an embedded client, or as part of a standalone
server. Proposers Propose additions to the paxos log and uses the acceptors to
get consensus on if the proposed value was accepted.

#### func (*Proposer) Propose

```go
func (p *Proposer) Propose(ctx context.Context, value []byte) (accepted []byte, err error)
```

#### type ProposerClient

```go
type ProposerClient interface {
	Propose(ctx context.Context, value []byte) ([]byte, error)
}
```


#### func  NewYarpcProposerClient

```go
func NewYarpcProposerClient(cc *yarpc.ClientConn) ProposerClient
```
NewYarpcProposerClient wraps the provided yarpc.ClientConn with an
ProposerClient implementation.

#### type ProposerServer

```go
type ProposerServer interface {
	Propose(ctx context.Context, value []byte) ([]byte, error)
}
```


#### type Request

```go
type Request struct {
	ID      uint64 `json:"id,omitempty"`
	Attempt uint64 `json:"attempt,omitempty"`
}
```

Request is used during the PREPARE and OBSERVE phases of the paxos algorithm.
Prepare sends along their ID value and attempt number, where Observe sends along
their last accepted id.

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

Stream provides an abstract definition of the functionality the underlying
stream needs to provide.

#### type Vote

```go
type Vote struct {
	// Member contains which member of the cluster cast this vote.
	Member string
	// Payload contains the payload of the message we're voting on. This is usually a Promise or Proposal.
	Payload interface{}
}
```

Vote is an internal structure used by multiple components to cast votes on
behalf of the acceptor that they're communicating with.
