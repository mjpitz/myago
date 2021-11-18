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
