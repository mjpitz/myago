// Copyright (C) The AetherFS Authors - All Rights Reserved
// See LICENSE for more information.

package oidcauth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/mjpitz/myago/lazy"
	"github.com/mjpitz/myago/ulid"
	"github.com/mjpitz/myago/zaputil"
)

// TokenCallback is invoked by the OIDCServeMux endpoint when we've successfully received and validated the
// authenticated user session.
type TokenCallback func(token *oauth2.Token)

// ServeMux is some rough code that should allow a command line tool to receive a token and invoke the provided
// callback function when a successful exchange is performed.
func ServeMux(cfg Config, callback TokenCallback) *http.ServeMux {
	providerOnce := &lazy.Once{
		Loader: cfg.Issuer.Provider,
	}

	configOnce := &lazy.Once{
		Loader: func(ctx context.Context) (*oauth2.Config, error) {
			instance, err := providerOnce.Get(ctx)
			if err != nil {
				return nil, err
			}

			return &oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				Endpoint:     instance.(*oidc.Provider).Endpoint(),
				RedirectURL:  cfg.RedirectURL,
				Scopes:       cfg.Scopes.Value(),
			}, nil
		},
	}

	verifierOnce := &lazy.Once{
		Loader: func(ctx context.Context) (*oidc.IDTokenVerifier, error) {
			provider, err := providerOnce.Get(ctx)
			if err != nil {
				return nil, err
			}

			config, err := configOnce.Get(ctx)
			if err != nil {
				return nil, err
			}

			return provider.(*oidc.Provider).Verifier(&oidc.Config{
				ClientID: config.(*oauth2.Config).ClientID,
			}), nil
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		generator := ulid.Extract(ctx)

		state, stateErr := generator.Generate(ctx, 128)

		config, configErr := configOnce.Get(ctx)
		if configErr != nil || stateErr != nil {
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		url := config.(*oauth2.Config).AuthCodeURL(state.String())
		http.Redirect(w, r, url, http.StatusFound)
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zaputil.Extract(ctx)

		config, configErr := configOnce.Get(ctx)
		verifier, verifierErr := verifierOnce.Get(ctx)

		if configErr != nil || verifierErr != nil {
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		query := r.URL.Query()
		errKind := query.Get("error")
		errDescription := query.Get("error_description")
		if errKind != "" {
			http.Error(w, errKind+": "+errDescription, http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.(*oauth2.Config).Exchange(ctx, query.Get("code"))
		if err != nil {
			logger.Error("failed to exchange code for auth info", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		rawIDToken := oauth2Token.Extra("id_token")

		idToken, err := verifier.(*oidc.IDTokenVerifier).Verify(ctx, rawIDToken.(string))
		if err != nil {
			logger.Error("failed to verify id token", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		err = idToken.VerifyAccessToken(oauth2Token.AccessToken)
		if err != nil {
			logger.Error("failed to verify access token token", zap.Error(err))
			http.Error(w, "", http.StatusUnauthorized)

			return
		}

		reader := strings.NewReader("You have successfully logged in. You may now close this tab in your browser.")
		http.ServeContent(w, r, "", time.Now(), reader)

		go callback(oauth2Token)
	})

	return mux
}
