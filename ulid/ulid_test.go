package ulid_test

import (
	"testing"

	"github.com/mjpitz/myago/ulid"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	{ // 64bit
		u, err := ulid.Parse("040006pgpa034")
		require.NoError(t, err)

		require.Equal(t, byte(1), u.Skew())
		require.Equal(t, int64(449884800), u.Timestamp().Unix())
		require.Len(t, u.Payload(), 1)
	}

	{ // 96bit
		u, err := ulid.Parse("040006pgpa08w7xa02s0")
		require.NoError(t, err)

		require.Equal(t, byte(1), u.Skew())
		require.Equal(t, int64(449884800), u.Timestamp().Unix())
		require.Len(t, u.Payload(), 5)
	}

	{ // 128bit
		u, err := ulid.Parse("040006pgpa0avg9vkqdast68zc")
		require.NoError(t, err)

		require.Equal(t, byte(1), u.Skew())
		require.Equal(t, int64(449884800), u.Timestamp().Unix())
		require.Len(t, u.Payload(), 9)
	}

	{ // 256bit
		u, err := ulid.Parse("040006pgpa072nrkstejwat4cq3swwd88xh62afnckn7qw0wzdng")
		require.NoError(t, err)

		require.Equal(t, byte(1), u.Skew())
		require.Equal(t, int64(449884800), u.Timestamp().Unix())
		require.Len(t, u.Payload(), 25)
	}
}
