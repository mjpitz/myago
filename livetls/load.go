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
	"crypto/x509"
	"path/filepath"

	"github.com/spf13/afero"

	"go.pitz.tech/lib/vfs"
)

// LoadCertPool loads the x509 certificate authority pool.
func LoadCertPool(ctx context.Context, cfg *Config) (*x509.CertPool, error) {
	if cfg.CAFile == "" {
		return nil, nil
	}

	fs := vfs.Extract(ctx)
	caPath := filepath.Join(cfg.CertPath, cfg.CAFile)

	ok, err := afero.Exists(fs, caPath)
	if !ok || err != nil {
		return nil, err
	}

	caData, err := afero.ReadFile(fs, caPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)

	return certPool, nil
}

// LoadCertificate loads the certificate from the configured public/private key.
func LoadCertificate(ctx context.Context, cfg *Config) (*tls.Certificate, error) {
	if cfg.CertFile == "" || cfg.KeyFile == "" {
		return nil, nil
	}

	fs := vfs.Extract(ctx)
	certPath := filepath.Join(cfg.CertPath, cfg.CertFile)
	keyPath := filepath.Join(cfg.CertPath, cfg.KeyFile)

	certOK, certErr := afero.Exists(fs, certPath)
	keyOK, keyErr := afero.Exists(fs, keyPath)

	switch {
	case certErr != nil || !certOK:
		return nil, certErr
	case keyErr != nil || !keyOK:
		return nil, keyErr
	}

	certData, certErr := afero.ReadFile(fs, certPath)
	keyData, keyErr := afero.ReadFile(fs, keyPath)

	switch {
	case certErr != nil:
		return nil, certErr
	case keyErr != nil:
		return nil, keyErr
	}

	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}
