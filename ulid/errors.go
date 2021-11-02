package ulid

import "fmt"

var (
	// ErrInvalidBitCount is returned when an invalid number of bits is provided to the Generate method of a Generator.
	ErrInvalidBitCount = fmt.Errorf("bits must be divisible by 8")

	// ErrNotEnoughBits is returned when fewer than 64 bit ULIDs are requested to be generated.
	ErrNotEnoughBits = fmt.Errorf("must be at least 64 bits")
)
