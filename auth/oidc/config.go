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

package oidcauth

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/urfave/cli/v2"

	"go.pitz.tech/lib/livetls"
)

// Issuer defines data needed to establish a connection to an issuer.
type Issuer struct {
	ServerURL            string `json:"server_url"            usage:"the address of the server where user authentication is performed"`
	CertificateAuthority string `json:"certificate_authority" usage:"path pointing to a file containing the certificate authority data for the server"`
}

func (i Issuer) Provider(ctx context.Context) (*oidc.Provider, error) {
	tlsConfig, err := livetls.New(ctx, livetls.Config{
		Enable: len(i.CertificateAuthority) > 0,
		CAFile: i.CertificateAuthority,
	})

	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		ctx = oidc.ClientContext(ctx, &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       tlsConfig,
			},
		})
	}

	return oidc.NewProvider(ctx, i.ServerURL)
}

// Config defines the information needed for an application to obtain an identity token from a provider.
type Config struct {
	Issuer       Issuer           `json:"issuer"`
	ClientID     string           `json:"client_id"     usage:"the client_id associated with this service"`
	ClientSecret string           `json:"client_secret" usage:"the client_secret associated with this service"`
	RedirectURL  string           `json:"redirect_url"  usage:"the redirect_url used by this service to obtain a token"`
	Scopes       *cli.StringSlice `json:"scopes"        usage:"specify the scopes that this authorization requires"     default:"openid,profile,email"`
}

// ClientConfig encapsulates the information needed to establish a client connection to an identity provider.
type ClientConfig struct {
	Issuer Issuer `json:"issuer"`
}
