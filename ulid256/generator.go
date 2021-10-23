package ulid256

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"sync"

	"github.com/jonboulle/clockwork"
)

const defaultSkew = 0

var mu = sync.Mutex{}
var registry = make(map[byte]*Generator)

// Random provides a default random generator that fills the payload with random data from a cryptographically secure
// source. Since it is the default generate, it uses version 1 (max of 255)
var Random, _ = NewGenerator(1, RandomFill())

// NewGenerator produces a generator responsible for constructing ULIDs. It maintains the core wire format for the ULID.
func NewGenerator(version byte, fill Fill) (*Generator, error) {
	mu.Lock()
	defer mu.Unlock()

	if existing, ok := registry[version]; ok {
		return existing, fmt.Errorf("version already registered")
	}

	registry[version] = &Generator{
		version: version,
		clock:   clockwork.NewRealClock(),
		fill:    fill,
	}

	return registry[version], nil
}

// Generator is used to generate ULIDs.
type Generator struct {
	version byte
	clock   clockwork.Clock
	fill    Fill
}

// WithClock allows callers to override the clock implementation being used.
func (g *Generator) WithClock(clock clockwork.Clock) *Generator {
	return &Generator{
		version: g.version,
		clock:   clock,
		fill:    g.fill,
	}
}

// New constructs a new ULID with the configured metadata.
func (g *Generator) New() (ULID, error) {
	timestamp := g.clock.Now()

	sec := make([]byte, 8)
	nsec := make([]byte, 4)
	binary.BigEndian.PutUint64(sec, uint64(timestamp.Unix()))
	binary.BigEndian.PutUint32(nsec, uint32(timestamp.Nanosecond()))

	ulid := ULID{}
	binary.BigEndian.PutUint16(ulid[:2], defaultSkew)
	copy(ulid[2:8], sec[2:])
	copy(ulid[8:11], nsec[1:])

	ulid[11] = g.version
	n, err := g.fill(ulid[12:28])
	switch {
	case err != nil:
		return Empty, err
	case n != 16:
		return Empty, fmt.Errorf("failed to fill 15 bytes of data")
	}

	checksum := crc32.NewIEEE()
	_, err = checksum.Write(ulid[:28])
	if err != nil {
		return Empty, err
	}

	binary.BigEndian.PutUint32(ulid[28:32], checksum.Sum32())

	return ulid, nil
}

// Must ensures that a ULID was returned.
func Must(ulid ULID, err error) ULID {
	if err != nil {
		panic(err)
	}
	return ulid
}
