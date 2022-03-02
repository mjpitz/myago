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

package storjauth

import (
	"context"
	"crypto/aes"
	"crypto/sha256"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	oidcauth "github.com/mjpitz/myago/auth/oidc"
	"github.com/mjpitz/myago/lazy"
	"github.com/mjpitz/myago/ulid"
	"github.com/mjpitz/myago/zaputil"

	"storj.io/common/base58"
)

// TokenCallback is invoked by the OIDCServeMux endpoint when we've successfully received and validated the
// authenticated user session.
type TokenCallback func(token *oauth2.Token, encryptionKey []byte)

// ServeMux is some rough code that should allow a command line tool to receive a token and invoke the provided
// callback function when a successful exchange is performed.
func ServeMux(cfg oidcauth.Config, callback TokenCallback) *http.ServeMux {
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

	//verifierOnce := &lazy.Once{
	//	Loader: func(ctx context.Context) (*oidc.IDTokenVerifier, error) {
	//		provider, err := providerOnce.Get(ctx)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		config, err := configOnce.Get(ctx)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		return provider.(*oidc.Provider).Verifier(&oidc.Config{
	//			ClientID: config.(*oauth2.Config).ClientID,
	//		}), nil
	//	},
	//}

	holder := make(chan ulid.ULID, 1)
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		generator := ulid.Extract(ctx)

		config, err := configOnce.Get(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		state, err := generator.Generate(ctx, 128)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		key, err := generator.Generate(ctx, 256)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		holder <- key

		url := config.(*oauth2.Config).AuthCodeURL(state.String()) + "#" + key.String()
		http.Redirect(w, r, url, http.StatusFound)
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zaputil.Extract(ctx)

		config, configErr := configOnce.Get(ctx)
		provider, providerErr := providerOnce.Get(ctx)
		//verifier, verifierErr := verifierOnce.Get(ctx)

		if configErr != nil || providerErr != nil {
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

		source := config.(*oauth2.Config).TokenSource(ctx, oauth2Token)

		userInfo, err := provider.(*oidc.Provider).UserInfo(ctx, source)
		if err != nil {
			logger.Error("failed to exchange token for user info", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		claims := make(map[string]interface{})
		err = userInfo.Claims(&claims)
		if err != nil {
			logger.Error("failed to obtain claims for user", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		cubbyhole := claims["cubbyhole"].(string)
		if cubbyhole == "" {
			logger.Error("missing cubbyhole")
			http.Error(w, "missing cubbyhole", http.StatusBadRequest)

			return
		}

		encryptedRootKey := base58.Decode(cubbyhole)

		key := <- holder
		aesKey := sha256.Sum256(key.Bytes())

		cipher, err := aes.NewCipher(aesKey[:])
		if err != nil {
			logger.Error("failed to construct aes cipher", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		rootKey := make([]byte, len(encryptedRootKey))
		cipher.Decrypt(rootKey, encryptedRootKey)

		// not yet supported
		//
		//rawIDToken := oauth2Token.Extra("id_token")
		//
		//idToken, err := verifier.(*oidc.IDTokenVerifier).Verify(ctx, rawIDToken.(string))
		//if err != nil {
		//	logger.Error("failed to verify id token", zap.Error(err))
		//	http.Error(w, "", http.StatusInternalServerError)
		//
		//	return
		//}
		//
		//err = idToken.VerifyAccessToken(oauth2Token.AccessToken)
		//if err != nil {
		//	logger.Error("failed to verify access token token", zap.Error(err))
		//	http.Error(w, "", http.StatusUnauthorized)
		//
		//	return
		//}

		reader := strings.NewReader("You have successfully logged in. You may now close this tab in your browser.")
		http.ServeContent(w, r, "", time.Now(), reader)

		go callback(oauth2Token, rootKey)
	})

	return mux
}
