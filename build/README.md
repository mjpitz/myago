# build

```go
import go.pitz.tech/lib/build
```

## Usage

#### type Info

```go
type Info struct {
	OS           string
	Architecture string

	GoVersion  string
	CGoEnabled bool

	Version  string
	VCS      string
	Revision string
	Compiled time.Time
	Modified bool
}
```

Info defines common build information associated with the binary.

#### func ParseInfo

```go
func ParseInfo() (info Info)
```

ParseInfo extracts as much build information from the compiled binary as it can.

#### func (Info) Metadata

```go
func (info Info) Metadata() map[string]any
```

Metadata formats the underlying information as a generalized key-value map.
