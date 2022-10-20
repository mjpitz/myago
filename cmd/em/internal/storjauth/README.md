# storjauth




```go
import go.pitz.tech/lib/cmd/em/internal/storjauth
```

## Usage

#### func  ServeMux

```go
func ServeMux(cfg oidcauth.Config, callback TokenCallback) *http.ServeMux
```
ServeMux is some rough code that should allow a command line tool to receive a
token and invoke the provided callback function when a successful exchange is
performed.

#### type TokenCallback

```go
type TokenCallback func(token *oauth2.Token, encryptionKey []byte)
```

TokenCallback is invoked by the OIDCServeMux endpoint when we've successfully
received and validated the authenticated user session.
