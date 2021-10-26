# lifecycle
--
    import "github.com/mjpitz/myago/lifecycle"

Package lifecycle provides common code for hooking into a golang application
lifecycle such as setting up a shutdown hook and deferring functions until
application shutdown.

## Usage

#### func  Defer

```go
func Defer(fn func(ctx context.Context))
```
Defer will enqueue a function that will be invoked by Resolve.

#### func  Resolve

```go
func Resolve(ctx context.Context)
```
Resolve will process all functions that have been enqueued by Defer up until
this point.

#### func  Setup

```go
func Setup(ctx context.Context) context.Context
```
Setup initializes a shutdown hook that cancels the underlying context.
