package ulid

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"sync"

	"github.com/jonboulle/clockwork"
)

// Generator is the base interface defines how to generate ULIDs of varying length.
type Generator interface {
	Generate(bits int) (ULID, error)
}

// BaseGenerator is the common core generator for ULIDs. All ULIDs should extend the base Generator. Extensions are free
// to format the `payload` portion of the ULID however they like.
type BaseGenerator struct {
	init  sync.Once
	Skew  byte
	Clock clockwork.Clock
}

func (g *BaseGenerator) Generate(bits int) (ULID, error) {
	g.init.Do(func() {
		if g.Clock == nil {
			g.Clock = clockwork.NewRealClock()
		}
	})

	switch {
	case bits < 64:
		return nil, ErrNotEnoughBits
	case bits%8 > 0:
		return nil, ErrInvalidBitCount
	}

	unix := uint64(g.Clock.Now().Unix())
	seconds := make([]byte, 8)
	binary.BigEndian.PutUint64(seconds, unix)

	ulid := make(ULID, bits/8)
	ulid[SkewOffset] = g.Skew
	copy(ulid[UnixOffset:PayloadOffset], seconds[2:8])

	return ulid, nil
}

var _ Generator = &BaseGenerator{}

// RandomGenerator produces a randomly generated ULID. It fills the payload portion of the ULID with random data.
type RandomGenerator struct {
	BaseGenerator

	init   sync.Once
	Reader io.Reader
}

func (g *RandomGenerator) Generate(bits int) (ULID, error) {
	g.init.Do(func() {
		if g.Reader == nil {
			g.Reader = rand.Reader
		}
	})

	base, err := g.BaseGenerator.Generate(bits)
	if err != nil {
		return nil, err
	}

	randomBytes := len(base) - PayloadOffset
	random := make([]byte, randomBytes)
	n, err := io.ReadFull(g.Reader, random)
	if err != nil {
		return nil, err
	} else if n != randomBytes {
		return nil, io.ErrUnexpectedEOF
	}

	copy(base[PayloadOffset:], random)
	return base, nil
}
