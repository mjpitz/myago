# index




```go
import go.pitz.tech/lib/cmd/em/internal/index
```

## Usage

#### type Config

```go
type Config struct {
	DatabaseDSN string `json:"database_dsn" usage:"specify the connection string for database" default:"file:db.sqlite"`
}
```


#### type Index

```go
type Index struct {
}
```


#### func  Open

```go
func Open(cfg Config) (*Index, error)
```

#### func (*Index) Close

```go
func (i *Index) Close() error
```

#### func (*Index) Index

```go
func (i *Index) Index(docs ...interface{}) error
```

#### func (*Index) Migrate

```go
func (i *Index) Migrate(schema ...interface{}) error
```
