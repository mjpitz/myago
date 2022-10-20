// Copyright (C) The AetherFS Authors - All Rights Reserved
// See LICENSE for more information.

package oidcauth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"go.pitz.tech/lib/auth"
	"go.pitz.tech/lib/headers"
	"go.pitz.tech/lib/lazy"
	"go.pitz.tech/lib/zaputil"
)

// OIDC returns a HandlerFunc who authenticates a user with the provided issuer using an access_token attached to the
// request. If provided, this access_token is exchanged for the authenticated user's information. It's important to know
// that this function does not handle authorization and requires an additional HandleFunc to do so.
func OIDC(cfg Issuer) auth.HandlerFunc {
	provider := &lazy.Once{
		Loader: cfg.Provider,
	}

	return func(ctx context.Context) (context.Context, error) {
		header := headers.Extract(ctx)
		logger := zaputil.Extract(ctx)

		accessToken, err := auth.Get(header, "bearer")
		if err != nil {
			logger.Error("failed to obtain bearer token", zap.Error(err))

			return ctx, nil
		}

		provider, err := provider.Get(ctx)
		if err != nil {
			logger.Error("error establishing connection with provider", zap.Error(err))

			return nil, errInternal
		}

		// fetch oidc.UserInfo and put on request
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
		})

		// user is authenticated, but authentication appears to have expired
		// return unauthenticated error here to trigger re-authentication
		userInfo, err := provider.(*oidc.Provider).UserInfo(ctx, tokenSource)
		if err != nil {
			logger.Error("error fetching user information given access token", zap.Error(err))

			return nil, errUnauthorized
		}

		info := auth.UserInfo{}
		err = userInfo.Claims(&info)
		if err != nil {
			return nil, errUnauthorized
		}

		ctx = auth.ToContext(ctx, info)

		// attach user information
		return ctx, nil
	}
}
