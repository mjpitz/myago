# libctx

Package libctx provides common code for working with contexts and may eventually
write its own.

```go
import go.pitz.tech/lib/libctx
```

## Usage

#### type Key

```go
type Key string
```

Key provides a scoped key used to persist data on contexts.

#### func (Key) String

```go
func (c Key) String() string
```
