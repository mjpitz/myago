package ulid_test

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/mjpitz/myago/ulid"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	clock := clockwork.NewFakeClock()
	generator := &ulid.RandomGenerator{
		BaseGenerator: ulid.BaseGenerator{
			Skew:  1,
			Clock: clock,
		},
	}

	{
		ulid, err := generator.Generate(32)
		require.Error(t, err)
		require.Equal(t, "must be at least 64 bits", err.Error())
		require.Nil(t, ulid)
	}

	{
		ulid, err := generator.Generate(67)
		require.Error(t, err)
		require.Equal(t, "bits must be divisible by 8", err.Error())
		require.Nil(t, ulid)
	}

	{
		ulid, err := generator.Generate(64)
		require.NoError(t, err)
		require.Equal(t, byte(1), ulid.Skew())
		require.Equal(t, clock.Now().Unix(), ulid.Timestamp().Unix())
		require.Len(t, ulid.Payload(), 1)

		t.Log("ulid(64): ", ulid.String())
	}

	{
		ulid, err := generator.Generate(96)
		require.NoError(t, err)
		require.Equal(t, byte(1), ulid.Skew())
		require.Equal(t, clock.Now().Unix(), ulid.Timestamp().Unix())
		require.Len(t, ulid.Payload(), 5)

		t.Log("ulid(96): ", ulid.String())
	}

	{
		ulid, err := generator.Generate(128)
		require.NoError(t, err)
		require.Equal(t, byte(1), ulid.Skew())
		require.Equal(t, clock.Now().Unix(), ulid.Timestamp().Unix())
		require.Len(t, ulid.Payload(), 9)

		t.Log("ulid(128): ", ulid.String())
	}

	{
		ulid, err := generator.Generate(256)
		require.NoError(t, err)
		require.Equal(t, byte(1), ulid.Skew())
		require.Equal(t, clock.Now().Unix(), ulid.Timestamp().Unix())
		require.Len(t, ulid.Payload(), 25)

		t.Log("ulid(256): ", ulid.String())
	}
}
