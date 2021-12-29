# lifecycle

Package lifecycle provides common code for hooking into a golang application
lifecycle such as setting up a shutdown hook and deferring functions until
application shutdown.

```go
import github.com/mjpitz/myago/lifecycle
```

## Usage

#### func Defer

```go
func Defer(fn func(ctx context.Context))
```

Defer will enqueue a function that will be invoked by Resolve.

#### func Resolve

```go
func Resolve(ctx context.Context)
```

Resolve will process all functions that have been enqueued by Defer up until
this point.

#### func Setup

```go
func Setup(ctx context.Context) context.Context
```

Setup initializes a shutdown hook that cancels the underlying context.

#### func Shutdown

```go
func Shutdown(ctx context.Context)
```

Shutdown halts the context, stopping any lingering processes.

#### type LifeCycle

```go
type LifeCycle struct {
}
```

LifeCycle hooks into various lifecycle events. It allows functions to be
deferred en masse.

#### func (\*LifeCycle) Defer

```go
func (lc *LifeCycle) Defer(fn func(ctx context.Context))
```

Defer will enqueue a function that will be invoked by Resolve.

#### func (\*LifeCycle) Resolve

```go
func (lc *LifeCycle) Resolve(ctx context.Context)
```

Resolve will process all functions that have been enqueued by Defer up until
this point.

#### func (\*LifeCycle) Setup

```go
func (lc *LifeCycle) Setup(ctx context.Context) context.Context
```

Setup initializes a shutdown hook that cancels the underlying context.

#### func (\*LifeCycle) Shutdown

```go
func (lc *LifeCycle) Shutdown(ctx context.Context)
```
