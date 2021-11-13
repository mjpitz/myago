package livetls_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/clocks"
	"github.com/mjpitz/myago/livetls"
)

func TestNew(t *testing.T) {
	frozen := time.Now()
	clock := clockwork.NewFakeClockAt(frozen)
	ctx := clocks.ToContext(context.Background(), clock)

	{
		tlsConfig, err := livetls.New(ctx, livetls.ReloadingConfig{
			Config: livetls.Config{
				Enable:   false,
			},
		})

		require.Nil(t, tlsConfig)
		require.NoError(t, err)
	}

	{
		tlsConfig, err := livetls.New(ctx, livetls.ReloadingConfig{
			Config: livetls.Config{
				Enable:   true,
				CertPath: "sslconf",
				CAFile:   "ca.pem",
				CertFile: "cert.pem",
				KeyFile:  "key.pem",
			},
			ReloadInterval: time.Second,
		})
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)

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
