# crockford32

```go
import go.pitz.tech/lib/crockford32
```

## Usage

```go
var (
	// Encoding provides a common implementation of a crockford base32 encoding.
	Encoding = base32.
		NewEncoding("0123456789abcdefghjkmnpqrstvwxyz").
		WithPadding(base32.NoPadding)
)
```
