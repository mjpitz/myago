# clocks
--
    import "github.com/mjpitz/myago/clocks"

Package clocks provides code for setting up and managing clocks on contexts.

## Usage

#### func  Extract

```go
func Extract(ctx context.Context) clockwork.Clock
```
Extract pulls the clock from the provided context. If no clock is found, then
the defaultClock is returned.

#### func  Setup

```go
func Setup(ctx context.Context) context.Context
```
Setup sets the defaultClock on the provided context. This can always be
overridden later.

#### func  ToContext

```go
func ToContext(ctx context.Context, clock clockwork.Clock) context.Context
```
ToContext sets the clock on the provided context.
