# ulid256
--
    import "github.com/mjpitz/myago/ulid256"

Package ulid256 provides functionality for generating unique, lexigraphic
identifiers (ULID). There are many differences between this implementation and
it's 128bit counterpart, but their goals are roughly the same.

## Usage

```go
var Empty = ULID{}
```
Empty is an empty ULID.

```go
var Random, _ = NewGenerator(1, RandomFill())
```
Random provides a default random generator that fills the payload with random
data from a cryptographically secure source. Since it is the default generate,
it uses version 1 (max of 255)

#### type Fill

```go
type Fill func(data []byte) (int, error)
```

Fill defines an arbitrary way to fill a slice of data with information. This
information can be structureless or semi-structured. In the end, it should not
be used for any control logic.

#### func  RandomFill

```go
func RandomFill() Fill
```
RandomFill fills the data array with random data from a cryptographically secure
source.

#### func  ServerIDFill

```go
func ServerIDFill(serverId uint16, fill Fill) Fill
```
ServerIDFill prefixes the data payload with a serverID. It delegates filling the
remaining portion to the provided fill.

#### type Generator

```go
type Generator struct {
}
```

Generator is used to generate ULIDs.

#### func  NewGenerator

```go
func NewGenerator(version byte, fill Fill) (*Generator, error)
```
NewGenerator produces a generator responsible for constructing ULIDs. It
maintains the core wire format for the ULID.

#### func (*Generator) New

```go
func (g *Generator) New() (ULID, error)
```
New constructs a new ULID with the configured metadata.

#### func (*Generator) WithClock

```go
func (g *Generator) WithClock(clock clockwork.Clock) *Generator
```
WithClock allows callers to override the clock implementation being used.

#### type ULID

```go
type ULID [32]byte
```

ULID is a unique, lexigraphic identifier. Unlike the the canonical ULID
implementation, this version is 256 bits, holds a version identifier, a
programmable the payload, and a CRC32 checksum. The binary format is as follows:

`[ skew - 2 bytes ][ sec - 6 bytes ][ nsec - 3 bytes ][ version - byte ][
payload - 16 bytes ][ checksum - 4 bytes ]`

#### func  Must

```go
func Must(ulid ULID, err error) ULID
```
Must ensures that a ULID was returned.

#### func  Parse

```go
func Parse(ulid string) (ULID, error)
```
Parse converts a string based ULID to it's binary representation.

#### func (ULID) Checksum

```go
func (u ULID) Checksum() uint32
```
Checksum returns and IEEE CRC32 checksum portion of the ULID.

#### func (ULID) Payload

```go
func (u ULID) Payload() []byte
```
Payload returns the data portion for the ULID.

#### func (ULID) Skew

```go
func (u ULID) Skew() uint16
```
Skew returns the clock skew factor.

#### func (ULID) String

```go
func (u ULID) String() string
```
String serializes the ULID to a string representation.

#### func (ULID) Time

```go
func (u ULID) Time() time.Time
```
Time returns the timestamp associated with this ULID.

#### func (ULID) Validate

```go
func (u ULID) Validate() error
```
Validate returns an error if the checksums do not match.

#### func (ULID) Version

```go
func (u ULID) Version() byte
```
Version returns the version of the payload data.
