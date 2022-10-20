# config

```go
import go.pitz.tech/lib/config
```

## Usage

```go
var (
	// ErrFileDoesNotExist is returned when the file we're interacting with does not exist.
	ErrFileDoesNotExist = errors.New("file does not exist")

	// ErrFileMissingExtension is returned when the provided file is missing an extension.
	ErrFileMissingExtension = errors.New("file missing extension")

	// ErrUnsupportedFileExtension is returned when we don't recognize a given file extension.
	ErrUnsupportedFileExtension = errors.New("unsupported file extension")
)
```

```go
var DefaultLoader = Loader{
	".json": encoding.JSON,
	".toml": encoding.TOML,
	".yaml": encoding.YAML,
	".yml":  encoding.YAML,
	".xml":  encoding.XML,
}
```

DefaultLoader provides a default Loader implementation that supports reading a
variety of files.

#### func Load

```go
func Load(ctx context.Context, v interface{}, filePaths ...string) error
```

Load provides a convenience function for being able to load configuration using
the DefaultLoader.

#### type Loader

```go
type Loader map[string]*encoding.Encoding
```

Loader provides functionality for reading a variety of file formats into a
struct.

#### func (Loader) Load

```go
func (l Loader) Load(ctx context.Context, v interface{}, filePaths ...string) error
```

Load reads the provided files (if they exist) and unmarshals the data into the
provided interface.
