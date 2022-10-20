// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ulid_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/encoding"
	"go.pitz.tech/lib/ulid"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		ulid string
		// expectations
		skew       byte
		millis     int64
		payloadLen int
	}{
		{
			name:       "64 bit ulid",
			ulid:       "0406HFSS8G0AT",
			skew:       1,
			millis:     449884800000,
			payloadLen: 1,
		},
		{
			name:       "96 bit ulid",
			ulid:       "0406HFSS8G0DEMF9ENSG",
			skew:       1,
			millis:     449884800000,
			payloadLen: 5,
		},
		{
			name:       "128 bit ulid",
			ulid:       "0406HFSS8G0AGWPHSM7EFW8ZH4",
			skew:       1,
			millis:     449884800000,
			payloadLen: 9,
		},
		{
			name:       "256 bit ulid",
			ulid:       "0406HFSS8G04GAQGVET8C35S7DS8E28QCZ9AKRPS2X0NFPN1E9M0",
			skew:       1,
			millis:     449884800000,
			payloadLen: 25,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)

		u, err := ulid.Parse(testCase.ulid)
		require.NoError(t, err)

		require.Equal(t, testCase.skew, u.Skew())
		require.Equal(t, testCase.millis, u.Timestamp().UnixMilli())
		require.Len(t, u.Payload(), testCase.payloadLen)

		v, err := u.Value()
		require.NoError(t, err)

		w := ulid.ULID{}
		err = (&w).Scan(v)
		require.NoError(t, err)

		require.Equal(t, testCase.skew, w.Skew())
		require.Equal(t, testCase.millis, w.Timestamp().UnixMilli())
		require.Len(t, w.Payload(), testCase.payloadLen)

		buf := bytes.NewBuffer(nil)
		err = encoding.JSON.Encoder(buf).Encode(u)
		require.NoError(t, err)

		require.Equal(t, "\""+testCase.ulid+"\"\n", buf.String())

		x := ulid.ULID{}
		err = encoding.JSON.Decoder(bytes.NewBuffer(buf.Bytes())).Decode(&x)
		require.NoError(t, err)

		require.Equal(t, testCase.skew, x.Skew())
		require.Equal(t, testCase.millis, x.Timestamp().UnixMilli())
		require.Len(t, x.Payload(), testCase.payloadLen)
	}
}
