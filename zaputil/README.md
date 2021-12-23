# zaputil
--
    import "github.com/mjpitz/myago/zaputil"

Package zaputil contains common code for passing a zap logger around. It's a
convenient space for putting custom logger implementations for plugging into
various locations.

## Usage

#### func  Extract

```go
func Extract(ctx context.Context) *zap.Logger
```
Extract pulls the logger from the provided context. If no logger is found, then
the defaultLogger is returned.

#### func  HashiCorpStdLogger

```go
func HashiCorpStdLogger(logger *zap.Logger) *log.Logger
```
HashiCorpStdLogger wraps the provided logger with a golang logger to log
messages at the appropriate level using the HashiCorp log format. Useful for
replacing serf and membership loggers.

#### func  HashicorpStdLogger

```go
func HashicorpStdLogger(logger *zap.Logger) *log.Logger
```
HashicorpStdLogger Deprecated.

#### func  Setup

```go
func Setup(ctx context.Context, cfg Config) context.Context
```
Setup creates a logger given the provided configuration.

#### func  ToContext

```go
func ToContext(ctx context.Context, logger *zap.Logger) context.Context
```
ToContext sets the logger on the provided context.

#### type BadgerLogger

```go
type BadgerLogger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}
```

BadgerLogger is an interface pulled from the badger library. It defines the
functionality needed by the badger system to log messages. It supports a variety
of levels and works similar to the fmt.Printf method.

#### func  Badger

```go
func Badger(log *zap.Logger) BadgerLogger
```
Badger wraps the provided logger so badger can log using zap.

#### type Config

```go
type Config struct {
	Level  string `json:"level"  usage:"adjust the verbosity of the logs" default:"info"`
	Format string `json:"format" usage:"configure the format of the logs" default:"json"`
}
```

Config encapsulates the configurable elements of the logger.

#### func  DefaultConfig

```go
func DefaultConfig() Config
```
DefaultConfig returns the default configuration for zap to use. By default, it
logs at an info level and infers the log format based on the stdout device. If
it looks like a terminal session, then it uses the console format.
