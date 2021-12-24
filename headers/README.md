# headers
--
    import "github.com/mjpitz/myago/headers"

Package headers provides utility code for operating on header values that come
from different sources. A default HTTP middleware handler is provided and will
ensure that the headers are appropriately translated and passed along.

## Usage

#### func  HTTP

```go
func HTTP(delegate http.Handler) http.HandlerFunc
```
HTTP returns an http middleware function that translates HTTP headers into a
context Header.

#### func  ToContext

```go
func ToContext(ctx context.Context, header Header) context.Context
```
ToContext attaches the provided headers to the context.

#### type Header

```go
type Header map[string][]string
```

Header defines an abstract definition of a header.

#### func  Extract

```go
func Extract(ctx context.Context) Header
```
Extract attempts to obtain the headers from the provided context.

#### func  New

```go
func New() Header
```
New constructs a Header for use.

#### func (Header) Get

```go
func (h Header) Get(key string) string
```
Get returns the first possible header value for a key (if present).

#### func (Header) GetAll

```go
func (h Header) GetAll(key string) []string
```
GetAll returns all possible values for a key.

#### func (Header) Set

```go
func (h Header) Set(key, value string)
```
Set sets a single value for the provided key.

#### func (Header) SetAll

```go
func (h Header) SetAll(key string, values []string)
```
SetAll sets the values for the provides key.
