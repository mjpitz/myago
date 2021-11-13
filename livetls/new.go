package livetls

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/mjpitz/myago/clocks"
)

// New construct a tls.Config that will periodically reload the configured certificate.
// nolint:cyclop
func New(ctx context.Context, config Config) (*tls.Config, error) {
	if !config.Enable {
		return nil, nil
	}

	certPool, certPoolErr := LoadCertPool(ctx, &config)
	cert, certErr := LoadCertificate(ctx, &config)

	switch {
	case certPoolErr != nil:
		return nil, certPoolErr
	case certErr != nil:
		return nil, certErr
	case certPool == nil && cert == nil:
		return &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		}, nil
	}

	certCh := make(chan *tls.Certificate, 1)
	certCh <- cert

	clock := clocks.Extract(ctx)
	reloader := clock.NewTicker(config.ReloadInterval)

	getCertificate := func() (*tls.Certificate, error) {
		timeout := clock.NewTicker(time.Second)
		defer timeout.Stop()

		select {
		case <-reloader.Chan():
			reloader.Stop()
			defer func() { reloader = clock.NewTicker(config.ReloadInterval) }()

			newCert, certErr := LoadCertificate(ctx, &config)
			cert := <-certCh

			if certErr == nil && newCert != nil {
				cert = newCert
			}

			certCh <- cert

			return cert, nil

		case <-timeout.Chan():
			return nil, context.DeadlineExceeded

		case cert := <-certCh:
			certCh <- cert

			return cert, nil
		}
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		RootCAs:    certPool,
		ClientCAs:  certPool,
		GetCertificate: func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return getCertificate()
		},
		GetClientCertificate: func(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return getCertificate()
		},
	}, nil
}
