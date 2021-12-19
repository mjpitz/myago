# yarpc
--
    import "github.com/mjpitz/myago/yarpc"

Package yarpc implements "yet another RPC framework" on top of HashiCorp's yamux
library. I wanted something with the simplicity of Go's HTTP library and the
ability to easily manage connections like gRPC.

Why? gRPC comes with a rather large foot print and in many of these cases, I
wanted a slimmer package for passing messages between processes.

Example Server:

    type Stat struct {
    	Name string
    	Value int
    }

    start := time.Now()

    yarpc.HandleFunc("admin.stats", func(stream yarpc.Stream) error {
    	for {
    		err = stream.SendMsg(&Stat{ "uptime", time.Since(start).Seconds() })
    		if err != nil {
    			return err
    		}
    		time.Sleep(5 * time.Second)
    	}
    })

    yarpc.ListenAndServe("tcp", "0.0.0.0:8080")

Example ClientConn:

    ctx := context.Background()
    conn := yarpc.Dial("tcp", "localhost:8080")

    stream := conn.openStream(ctx, "admin.stats")

    stat := Stat{}
    for {
    	err = stream.RecvMsg(&stat)
    	if err != nil {
    		break
    	}

    	stat.Name // "uptime"
    	stat.Name // "uptime"
    }

## Usage

```go
var (
	// DefaultServeMux provides a default request multiplexer (router).
	DefaultServeMux = &ServeMux{}

	// DefaultServer is a global server definition that can be leveraged by hosting program.
	DefaultServer = &Server{
		Handler: DefaultServeMux,
	}
)
```

#### func  Handle

```go
func Handle(pattern string, handler Handler)
```
Handle adds the provided handler to the default server.

#### func  HandleFunc

```go
func HandleFunc(pattern string, handler func(Stream) error)
```
HandleFunc adds the provided handler function to the default server.

#### func  ListenAndServe

```go
func ListenAndServe(network, address string, opts ...Option) error
```
ListenAndServe starts the default server on the provided network and address.

#### func  Serve

```go
func Serve(listener Listener, opts ...Option) error
```
Serve starts the default server using the provided listener.

#### type ClientConn

```go
type ClientConn struct {
	Dialer Dialer
}
```

ClientConn defines an abstract connection yarpc clients to use.

#### func  DialContext

```go
func DialContext(ctx context.Context, network, target string, opts ...Option) *ClientConn
```
DialContext initializes a new client connection to the target server.

#### func  NewClientConn

```go
func NewClientConn(ctx context.Context) *ClientConn
```
NewClientConn creates a default ClientConn with an empty dialer implementation.
The Dialer must be configured before use. This function is intended to be used
in initializer functions such as DialContext.

#### func (*ClientConn) OpenStream

```go
func (c *ClientConn) OpenStream(ctx context.Context, method string) (Stream, error)
```
OpenStream starts a stream for the named RPC.

#### func (*ClientConn) WithOptions

```go
func (c *ClientConn) WithOptions(opts ...Option) *ClientConn
```
WithOptions configures the options for the underlying client connection.

#### type Dialer

```go
type Dialer interface {
	DialContext(ctx context.Context) (io.ReadWriteCloser, error)
}
```

Dialer provides a minimal interface needed to establish a client.

#### type Frame

```go
type Frame struct {
	Nonce  string      `json:"nonce,omitempty"`
	Status *Status     `json:"status,omitempty"`
	Body   interface{} `json:"body"`
}
```

Frame is the generalized structure passed along the wire.

#### type Handler

```go
type Handler interface {
	ServeYARPC(Stream) error
}
```

Handler defines an interface that can be used for handling requests.

#### type HandlerFunc

```go
type HandlerFunc func(Stream) error
```

HandlerFunc provides users with a simple functional interface for a Handler.

#### func (HandlerFunc) ServeYARPC

```go
func (fn HandlerFunc) ServeYARPC(stream Stream) error
```

#### type Invoke

```go
type Invoke struct {
	Method string `json:"method,omitempty"`
}
```


#### type Listener

```go
type Listener interface {
	Accept() (io.ReadWriteCloser, error)
	Close() error
}
```


#### type NetDialer

```go
type NetDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}
```

NetDialer provides a common interface for obtaining a net.Conn. This makes it
easy to handle TLS transparently.

#### type NetDialerAdapter

```go
type NetDialerAdapter struct {
	Dialer  NetDialer
	Network string
	Target  string
}
```

NetDialerAdapter adapts the provided NetDialer to support io.ReadWriteCloser.

#### func (*NetDialerAdapter) DialContext

```go
func (a *NetDialerAdapter) DialContext(ctx context.Context) (io.ReadWriteCloser, error)
```
DialContext returns a creates a new network connection.

#### type NetListenerAdapter

```go
type NetListenerAdapter struct {
	Listener net.Listener
}
```

NetListenerAdapter adapts the provided net.Listener to support
io.ReadWriteCloser.

#### func (*NetListenerAdapter) Accept

```go
func (n *NetListenerAdapter) Accept() (io.ReadWriteCloser, error)
```

#### func (*NetListenerAdapter) Close

```go
func (n *NetListenerAdapter) Close() error
```

#### type Option

```go
type Option func(opt *options)
```

Option defines an generic way to configure clients and servers.

#### func  WithContext

```go
func WithContext(ctx context.Context) Option
```
WithContext provides a custom context to the underlying system. Mostly used on
servers.

#### func  WithEncoding

```go
func WithEncoding(encoding *encoding.Encoding) Option
```
WithEncoding configures how messages are serialized.

#### func  WithTLS

```go
func WithTLS(config *tls.Config) Option
```
WithTLS enables TLS.

#### func  WithYamux

```go
func WithYamux(config *yamux.Config) Option
```
WithYamux configures yamux using the provided configuration.

#### type ServeMux

```go
type ServeMux struct {
}
```

ServeMux provides a router implementation for yarpc calls.

#### func (*ServeMux) Handle

```go
func (s *ServeMux) Handle(pattern string, handler Handler)
```

#### func (*ServeMux) ServeYARPC

```go
func (s *ServeMux) ServeYARPC(stream Stream) (err error)
```

#### type Server

```go
type Server struct {
	Handler Handler
}
```


#### func (*Server) ListenAndServe

```go
func (s *Server) ListenAndServe(network, address string, opts ...Option) error
```

#### func (*Server) Serve

```go
func (s *Server) Serve(listener Listener, opts ...Option) error
```

#### func (*Server) Shutdown

```go
func (s *Server) Shutdown() error
```

#### type Status

```go
type Status struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
```

Status reports an optional code and message along with the request.

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

Stream provides an interface for reading and writing message structures from a
stream.

#### func  Wrap

```go
func Wrap(ys *yamux.Stream, opts ...Option) Stream
```
Wrap converts the provided yamux stream into a yarpc Stream.
