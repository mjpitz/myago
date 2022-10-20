# httpauth




```go
import go.pitz.tech/lib/auth/http
```

## Usage

#### func  Handler

```go
func Handler(delegate http.Handler, handlers ...auth.HandlerFunc) http.HandlerFunc
```
Handler returns an HTTP middleware function that invokes the provided auth
handlers.
