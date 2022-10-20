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

package livetls_test

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/clocks"
	"go.pitz.tech/lib/livetls"
)

func TestNew(t *testing.T) {
	t.Parallel()

	frozen := time.Now()
	clock := clockwork.NewFakeClockAt(frozen)
	ctx := clocks.ToContext(context.Background(), clock)

	{
		tlsConfig, err := livetls.New(ctx, livetls.Config{
			Enable: false,
		})

		require.Nil(t, tlsConfig)
		require.NoError(t, err)
	}

	{
		tlsConfig, err := livetls.New(ctx, livetls.Config{
			Enable:         true,
			CertPath:       "sslconf",
			CAFile:         "ca.crt",
			CertFile:       "tls.crt",
			KeyFile:        "tls.key",
			ReloadInterval: time.Second,
		})
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)

		require.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
		require.Equal(t, uint16(tls.VersionTLS13), tlsConfig.MaxVersion)
		require.NotNil(t, tlsConfig.RootCAs)
		require.NotNil(t, tlsConfig.ClientCAs)
		require.NotNil(t, tlsConfig.GetCertificate)
		require.NotNil(t, tlsConfig.GetClientCertificate)

		for i := 0; i < 10; i++ {
			clock.Advance(time.Second)

			cert, err := tlsConfig.GetCertificate(nil)
			require.NoError(t, err)
			require.NotNil(t, cert)

			clientCert, err := tlsConfig.GetClientCertificate(nil)
			require.NoError(t, err)
			require.NotNil(t, cert)

			require.Equal(t, cert, clientCert)
		}
	}
}
