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
var DefaultServer = &Server{}
```
DefaultServer is a global server definition that can be leveraged by hosting
program.

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

#### type ClientConn

```go
type ClientConn struct {
}
```

ClientConn.

#### func  DialContext

```go
func DialContext(ctx context.Context, network, target string, opts ...Option) *ClientConn
```
DialContext initializes a new client connection to the target server.

#### func (*ClientConn) OpenStream

```go
func (c *ClientConn) OpenStream(ctx context.Context, method string) (Stream, error)
```
OpenStream starts a stream for a given RPC.

#### type Decoder

```go
type Decoder interface {
	Decode(i interface{}) error
}
```

Decoder reads message from the underlying stream.

#### type Dialer

```go
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}
```

Dialer provides a common interface for obtaining a net.Conn. This makes it easy
to handle TLS transparently.

#### type Encoder

```go
type Encoder interface {
	Encode(i interface{}) error
}
```

Encoder writes provided structures to the underlying stream.

#### type Encoding

```go
type Encoding interface {
	NewEncoder(io.Writer) Encoder
	NewDecoder(io.Reader) Decoder
}
```

Encoding describes a generalization used to create encoders and decoders for new
streams.

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
	Handle(Stream) error
}
```

Handler defines an interface that can be used for handling requests.

#### type HandlerFunc

```go
type HandlerFunc func(Stream) error
```

HandlerFunc provides users with a simple functional interface for a Handler.

#### func (HandlerFunc) Handle

```go
func (fn HandlerFunc) Handle(stream Stream) error
```

#### type Invoke

```go
type Invoke struct {
	Method string `json:"method,omitempty"`
}
```


#### type MSGPackEncoding

```go
type MSGPackEncoding struct{}
```

MSGPackEncoding uses msgpack out of box for a better balance of read/write
performance. JSON serialization is fast, but deserialization is much slower in
comparison (over 3x). While msgpack isn't as fast as protobuf, it offers
reasonable read/write performance.

#### func (*MSGPackEncoding) NewDecoder

```go
func (j *MSGPackEncoding) NewDecoder(reader io.Reader) Decoder
```

#### func (*MSGPackEncoding) NewEncoder

```go
func (j *MSGPackEncoding) NewEncoder(writer io.Writer) Encoder
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
func WithEncoding(encoding Encoding) Option
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

#### type Server

```go
type Server struct {
}
```


#### func (*Server) Handle

```go
func (s *Server) Handle(pattern string, handler Handler)
```

#### func (*Server) ListenAndServe

```go
func (s *Server) ListenAndServe(network, address string, opts ...Option) error
```

#### func (*Server) Serve

```go
func (s *Server) Serve(listener net.Listener, opts ...Option) error
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
