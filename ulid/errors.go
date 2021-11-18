package ulid

import (
	"errors"
)

var (
	// ErrInvalidBitCount is returned when an invalid number of bits is provided to the Generate method of a Generator.
	ErrInvalidBitCount = errors.New("bits must be divisible by 8")

	// ErrNotEnoughBits is returned when fewer than 64 bit ULIDs are requested to be generated.
	ErrNotEnoughBits = errors.New("must be at least 64 bits")

	// ErrInsufficientData is returned when the fill fails to return enough fata for the ULID.
	ErrInsufficientData = errors.New("failed to read sufficient payload data")
)
