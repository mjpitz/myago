# lazy

```go
import github.com/mjpitz/myago/lazy
```

## Usage

#### type Once

```go
type Once struct {
	// Loader is a function that returns an object and optional error. It conditionally accepts a context value.
	Loader interface{}
}
```

Once will attempt to load a value until one is loaded.

#### func (\*Once) Get

```go
func (o *Once) Get(ctx context.Context) (interface{}, error)
```

Get returns the loaded value if set or an error should one occur.
