# ulid
--
    import "github.com/mjpitz/myago/ulid"

Package ulid provides code for generating variable length unique, lexigraphic
identifiers (ULID) with programmable fills. Currently, there is a
RandomGenerator that can be used to generate ULIDs with a randomized payload. To
provide a custom payload, simply extend the BaseGenerator, and override the
Generate method. It's important to call the BaseGenerator's Generate method,
otherwise the skew and timestamp bits won't be set properly.

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
	ErrInvalidBitCount = fmt.Errorf("bits must be divisible by 8")

	// ErrNotEnoughBits is returned when fewer than 64 bit ULIDs are requested to be generated.
	ErrNotEnoughBits = fmt.Errorf("must be at least 64 bits")
)
```

#### type BaseGenerator

```go
type BaseGenerator struct {
	Skew  byte
	Clock clockwork.Clock
}
```

BaseGenerator is the common core generator for ULIDs. All ULIDs should extend
the base Generator. Extensions are free to format the `payload` portion of the
ULID however they like.

#### func (*BaseGenerator) Generate

```go
func (g *BaseGenerator) Generate(bits int) (ULID, error)
```

#### type Generator

```go
type Generator interface {
	Generate(bits int) (ULID, error)
}
```

Generator is the base interface defines how to generate ULIDs of varying length.

#### type RandomGenerator

```go
type RandomGenerator struct {
	BaseGenerator

	Reader io.Reader
}
```

RandomGenerator produces a randomly generated ULID. It fills the payload portion
of the ULID with random data.

#### func (*RandomGenerator) Generate

```go
func (g *RandomGenerator) Generate(bits int) (ULID, error)
```

#### type ULID

```go
type ULID []byte
```

ULID is a generalized, unique lexographical identifier. The format is as
follows:

`[ skew ][ sec ][ payload ]`

- `skew` - 1 byte used to handle major clock skews (reserved, unused) - `sec` -
6 bytes of a unix timestamp (should give us until the year 10k or so) -
`payload` - N bytes for the payload

#### func (ULID) Payload

```go
func (ulid ULID) Payload() []byte
```
Payload returns a copy of the payload bytes.

#### func (ULID) Skew

```go
func (ulid ULID) Skew() byte
```
Skew returns the current skew used to correct massive time skews.

#### func (ULID) Timestamp

```go
func (ulid ULID) Timestamp() time.Time
```
Timestamp returns the timestamp portion of the identifier.
