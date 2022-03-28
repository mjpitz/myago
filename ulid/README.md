# ulid

Package ulid provides code for generating variable length unique, lexigraphic
identifiers (ULID) with programmable fills. Currently, there is a
RandomGenerator that can be used to generate ULIDs with a randomized payload. To
provide a custom payload, simply extend the BaseGenerator, and override the
Generate method. It's important to call the BaseGenerator's Generate method,
otherwise the skew and timestamp bits won't be set properly.

Unlike the canonical [ULID](https://github.com/ulid/spec), this version holds a
placeholder byte for major clock skews which can often occur in distributed
systems. The wire format is as follows: `[ skew ][ millis ][ payload ]`

    - `skew` - 1 byte used to handle major clock skews (reserved, unused)
    - `millis` - 6 bytes of a unix timestamp (should give us until the year 10k or so)
    - `payload` - N bytes for the payload

```go
import github.com/mjpitz/myago/ulid
```

## Usage

```go
const (
	// SkewOffset is the starting byte position for the skew data.
	SkewOffset = 0
	// SkewLength is the number of bytes representing the skew.
	SkewLength = 1

	// UnixOffset is the starting byte position for the unix timestamp data.
	UnixOffset = 1
	// UnixLength is the number of bytes representing the unix timestamp data.
	UnixLength = 6

	// PayloadOffset is the starting byte position for the payload data.
	PayloadOffset = 7
)
```

```go
var (
	// ErrInvalidBitCount is returned when an invalid number of bits is provided to the Generate method of a Generator.
	ErrInvalidBitCount = errors.New("bits must be divisible by 8")

	// ErrNotEnoughBits is returned when fewer than 64 bit ULIDs are requested to be generated.
	ErrNotEnoughBits = errors.New("must be at least 64 bits")

	// ErrInsufficientData is returned when the fill fails to return enough fata for the ULID.
	ErrInsufficientData = errors.New("failed to read sufficient payload data")
)
```

#### func RandomFill

```go
func RandomFill(_ ULID, data []byte) (int, error)
```

RandomFill is a fill that populates the data payload with random data.

#### func ToContext

```go
func ToContext(ctx context.Context, generator *Generator) context.Context
```

ToContext appends the provided generator to the provided context.

#### type Fill

```go
type Fill func(ulid ULID, data []byte) (int, error)
```

Fill provides an abstraction for filling the data payload of a ULID.

#### type Generator

```go
type Generator struct {
}
```

Generator is the base interface defines how to generate ULIDs of varying length.

#### func Extract

```go
func Extract(ctx context.Context) *Generator
```

Extract is used to obtain the generator from a context. If none is present, the
system generator is used.

#### func NewGenerator

```go
func NewGenerator(skew byte, fill Fill) *Generator
```

NewGenerator constructs a generator using the provided skew and fill.

#### func (\*Generator) Generate

```go
func (g *Generator) Generate(ctx context.Context, bits int) (ULID, error)
```

Generate produces a new ULID unless an error occurs. A clock can be provided on
the context to override how the Generator obtains a timestamp.

#### type ULID

```go
type ULID []byte
```

ULID is a variable-length, generalized, unique lexographical identifier.

#### func Parse

```go
func Parse(ulid string) (ULID, error)
```

Parse accepts a ULID string and attempts to extract a ULID from the provided
string.

#### func (ULID) Bytes

```go
func (ulid ULID) Bytes() []byte
```

Bytes returns a copy of the underlying byte array backing the ulid.

#### func (ULID) MarshalJSON

```go
func (ulid ULID) MarshalJSON() ([]byte, error)
```

#### func (ULID) Payload

```go
func (ulid ULID) Payload() []byte
```

Payload returns a copy of the payload bytes.

#### func (\*ULID) Scan

```go
func (ulid *ULID) Scan(src interface{}) (err error)
```

Scan attempts to parse the provided src value into the ULID.

#### func (ULID) Skew

```go
func (ulid ULID) Skew() byte
```

Skew returns the current skew used to correct massive time skews.

#### func (ULID) String

```go
func (ulid ULID) String() string
```

String returns a string representation of the payload. It's encoded using a
crockford base32 encoding.

#### func (ULID) Timestamp

```go
func (ulid ULID) Timestamp() time.Time
```

Timestamp returns the timestamp portion of the identifier.

#### func (\*ULID) UnmarshalJSON

```go
func (ulid *ULID) UnmarshalJSON(bytes []byte) error
```

#### func (ULID) Value

```go
func (ulid ULID) Value() (driver.Value, error)
```

Value serializes the ULID so that it can be stored in SQL databases.
