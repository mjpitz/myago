package ulid_test

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/ulid"
)

func TestGenerator(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	generator := &ulid.RandomGenerator{
		BaseGenerator: ulid.BaseGenerator{
			Skew:  1,
			Clock: clock,
		},
	}

	testCases := []struct {
		name string
		bits int
		// expectations
		error      bool
		errorMsg   string
		skew       byte
		millis     int64
		payloadLen int
	}{
		{
			name:     "32 bit ulid",
			bits:     32,
			error:    true,
			errorMsg: "must be at least 64 bits",
		},
		{
			name:     "67 bit ulid",
			bits:     67,
			error:    true,
			errorMsg: "bits must be divisible by 8",
		},
		{
			name:       "64 bit ulid",
			bits:       64,
			skew:       1,
			millis:     clock.Now().UnixMilli(),
			payloadLen: 1,
		},
		{
			name:       "96 bit ulid",
			bits:       96,
			skew:       1,
			millis:     clock.Now().UnixMilli(),
			payloadLen: 5,
		},
		{
			name:       "128 bit ulid",
			bits:       128,
			skew:       1,
			millis:     clock.Now().UnixMilli(),
			payloadLen: 9,
		},
		{
			name:       "256 bit ulid",
			bits:       256,
			skew:       1,
			millis:     clock.Now().UnixMilli(),
			payloadLen: 25,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)

		ulid, err := generator.Generate(testCase.bits)

		if testCase.error {
			require.Error(t, err)
			require.Equal(t, testCase.errorMsg, err.Error())
			require.Nil(t, ulid)
		} else {
			require.NoError(t, err)
			require.Equal(t, testCase.skew, ulid.Skew())
			require.Equal(t, testCase.millis, ulid.Timestamp().UnixMilli())
			require.Len(t, ulid.Payload(), testCase.payloadLen)

			t.Log("ulid", testCase.bits, ulid.String())
		}
	}
}
