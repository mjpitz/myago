package ulid256_test

import (
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/ulid256"
)

func TestULID(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-10-23T06:07:14.043Z")
	require.NoError(t, err)

	fill := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, fill)
	require.Equal(t, 16, n)
	require.NoError(t, err)

	gen, err := ulid256.NewGenerator(0, func(data []byte) (int, error) {
		return copy(data, fill), nil
	})
	require.NoError(t, err)

	gen = gen.WithClock(clockwork.NewFakeClockAt(timestamp))

	u, err := gen.New()
	require.NoError(t, err)
	require.Equal(t, "AAAAAGFzppKQIM", u.String()[:14])

	require.Equal(t, uint16(0), u.Skew())
	require.Equal(t, "2021-10-23T06:07:14Z", u.Time().UTC().Format(time.RFC3339))
	require.Equal(t, byte(0), u.Version())
	require.Equal(t, fill, u.Payload())
	require.NoError(t, u.Validate())

	{
		ulidString := u.String()

		parse, err := ulid256.Parse(ulidString)
		require.NoError(t, err)

		require.Equal(t, u.Skew(), parse.Skew())
		require.Equal(t, u.Time().UnixNano(), parse.Time().UnixNano())
		require.Equal(t, u.Version(), parse.Version())
		require.Equal(t, u.Payload(), parse.Payload())
		require.Equal(t, u.Checksum(), parse.Checksum())
		require.NoError(t, parse.Validate())
	}

	{ // test checksum, ideally end users shouldn't be mutating the bytes directly
		u[0] = 200
		ulidString := u.String()

		_, err := ulid256.Parse(ulidString)
		require.Error(t, err)
	}
}
