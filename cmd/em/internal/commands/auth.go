// Copyright (C) 2022 Mya Pitzeruse
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

package commands

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"

	oidcauth "go.pitz.tech/lib/auth/oidc"
	"go.pitz.tech/lib/browser"
	"go.pitz.tech/lib/cmd/em/internal/storjauth"
	"go.pitz.tech/lib/flagset"
	"go.pitz.tech/lib/zaputil"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

var (
	oidcAuthConfig = &oidcauth.Config{
		Scopes: cli.NewStringSlice(),
	}

	storjAuthConfig = &oidcauth.Config{
		Scopes: cli.NewStringSlice(),
	}

	Auth = &cli.Command{
		Name:  "auth",
		Usage: "Authenticate using common mechanisms.",
		Subcommands: []*cli.Command{
			{
				Name:  "oidc",
				Usage: "Authenticate with an OIDC provider.",
				Flags: flagset.ExtractPrefix("em", oidcAuthConfig),
				Action: func(ctx *cli.Context) error {
					uri, err := url.Parse(oidcAuthConfig.RedirectURL)
					if err != nil {
						return err
					}

					svr := &http.Server{
						Addr: uri.Host,
					}

					if len(oidcAuthConfig.Scopes.Value()) == 0 {
						oidcAuthConfig.Scopes = cli.NewStringSlice("openid", "profile", "email")
					}

					cctx, cancel := context.WithCancel(ctx.Context)
					defer cancel()

					svr.Handler = oidcauth.ServeMux(*oidcAuthConfig, func(token *oauth2.Token) {
						defer cancel()

						enc := json.NewEncoder(ctx.App.Writer)
						enc.SetIndent("", "  ")
						_ = enc.Encode(token)
					})

					group := &errgroup.Group{}

					group.Go(func() error {
						time.Sleep(time.Second)
						return browser.Open(ctx.Context, uri.Scheme+"://"+uri.Host+"/login")
					})

					group.Go(svr.ListenAndServe)

					<-cctx.Done()
					err = svr.Shutdown(ctx.Context)
					_ = group.Wait()

					return nil
				},
				HideHelpCommand: true,
			},
			{
				Name:  "storj",
				Usage: "Authenticate with a Storj OIDC provider.",
				Flags: flagset.ExtractPrefix("em", storjAuthConfig),
				Action: func(ctx *cli.Context) error {
					uri, err := url.Parse(storjAuthConfig.RedirectURL)
					if err != nil {
						return err
					}

					svr := &http.Server{
						Addr: uri.Host,
						BaseContext: func(_ net.Listener) context.Context {
							return ctx.Context
						},
					}

					if len(storjAuthConfig.Scopes.Value()) == 0 {
						storjAuthConfig.Scopes = cli.NewStringSlice("openid", "profile", "email", "object:list", "object:read", "object:write", "object:delete")
					}

					cctx, cancel := context.WithCancel(ctx.Context)
					defer cancel()

					svr.Handler = storjauth.ServeMux(*storjAuthConfig, func(token *oauth2.Token, rootKey []byte) {
						defer cancel()

						enc := json.NewEncoder(ctx.App.Writer)
						enc.SetIndent("", "  ")
						_ = enc.Encode(struct {
							Token   *oauth2.Token `json:"token"`
							RootKey []byte        `json:"root_key"`
						}{
							Token:   token,
							RootKey: rootKey,
						})
					})

					group := &errgroup.Group{}

					group.Go(func() error {
						time.Sleep(time.Second)
						url := uri.Scheme+"://"+uri.Host+"/login"

						zaputil.Extract(ctx.Context).Info("Opening " + url)
						return browser.Open(ctx.Context, url)
					})

					group.Go(svr.ListenAndServe)

					<-cctx.Done()
					_ = svr.Shutdown(ctx.Context)
					_ = group.Wait()

					return nil
				},
				HideHelpCommand: true,
			},
		},
		HideHelpCommand: true,
	}
)
