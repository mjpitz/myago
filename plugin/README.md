# plugin
--
    import "github.com/mjpitz/myago/plugin"

Package plugin provides a simple plugin interface by forking processes and using
their stdout/stdin to enable communication between the parent process
(main-component) and the child (plugin). This is inspired by how protoc and its
various plugins work. Applications can read arguments, flags, and environment
variables provided to the program to configure its behaviour, but then stream
data from stdin to issue RPCs and write their responses to stdout.

## Usage

#### func  DialContext

```go
func DialContext(ctx context.Context, binary string, args ...string) *yarpc.ClientConn
```
DialContext returns a ClientConn whose dialer forks a process for the specified
binary.

#### func  Listen

```go
func Listen() yarpc.Listener
```
Listen returns a yarpc.Listener that treats a processes stdin and stdout as a
connection.

#### func  Pipe

```go
func Pipe() *pipe
```
Pipe returns a pseudo-async io.ReadWriteCloser. nolint:revive
