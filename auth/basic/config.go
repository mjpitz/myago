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

package basicauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"go.pitz.tech/lib/auth"
)

// UsernamePassword is used to authenticate a user using a username and password.
type UsernamePassword struct {
	Username string `json:"username" usage:"the username to login with"`
	Password string `json:"password" usage:"the password associated with the username"`
}

// AccessToken is used to authenticate a user using a bearer token.
type AccessToken struct {
	Token string `json:"token" usage:"the access token used to authenticate requests"`
}

// Config defines the options available to a server.
type Config struct {
	PasswordFile   string           `json:"password_file" usage:"path to the csv file containing usernames and passwords"`
	TokenFile      string           `json:"token_file" usage:"path to the csv file containing tokens"`
	StaticUsername string           `json:"static_username" usage:"provide a static username to authenticate the user" hidden:"true"`
	StaticPassword string           `json:"static_password" usage:"provide a static password to authenticate the user" hidden:"true"`
	StaticGroups   *cli.StringSlice `json:"static_groups" usage:"provide a static set of groups to assign to the user" hidden:"true"`
}

// ClientConfig defines the options available to a client.
type ClientConfig struct {
	UsernamePassword
	AccessToken
}

func (c ClientConfig) Token() (*oauth2.Token, error) {
	switch {
	case c.AccessToken.Token != "":
		return &oauth2.Token{
			TokenType:   "bearer",
			AccessToken: c.AccessToken.Token,
		}, nil
	case c.UsernamePassword.Username != "":
		accessToken := c.UsernamePassword.Username + ":" + c.UsernamePassword.Password
		accessToken = base64.StdEncoding.EncodeToString([]byte(accessToken))

		return &oauth2.Token{
			TokenType:   "basic",
			AccessToken: accessToken,
		}, nil
	}

	return nil, nil
}

// Handler returns the appropriate handler based on the provided configuration.
func Handler(ctx context.Context, cfg Config) (auth.HandlerFunc, error) {
	switch {
	case strings.HasSuffix(cfg.PasswordFile, ".csv"):
		return Basic(&LazyStore{
			Provider: func() (Store, error) {
				return OpenCSV(ctx, cfg.PasswordFile)
			},
		}), nil

	case strings.HasSuffix(cfg.TokenFile, ".csv"):
		return Bearer(&LazyStore{
			Provider: func() (Store, error) {
				return OpenCSV(ctx, cfg.TokenFile)
			},
		}), nil
	case cfg.StaticUsername != "" && cfg.StaticPassword != "":
		return Static(cfg.StaticUsername, cfg.StaticPassword, cfg.StaticGroups.Value()...), nil
	}

	return nil, fmt.Errorf("invalid file")
}
