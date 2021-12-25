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

package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/mjpitz/myago/headers"
)

// Basic implements a basic access authentication handler function. It parses values from the headers to obtain info
// about the authenticated user.
func Basic(store Store) HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		header := headers.Extract(ctx)
		authentication, err := Get(header, "basic")
		if err != nil {
			return ctx, nil
		}

		decoded, err := base64.StdEncoding.DecodeString(authentication)
		if err != nil {
			return ctx, nil
		}

		parts := strings.Split(string(decoded), ":")
		if len(parts) < 2 {
			return ctx, nil
		}

		username := parts[0]
		provided := parts[1]

		password, groups, err := store.Lookup(username)
		if err != nil {
			return ctx, nil
		}

		if provided != password {
			return ctx, nil
		}

		userInfo := &UserInfo{
			Subject: username,
			Profile: username,
		}

		err = userInfo.WithExtra(&Claims{
			Groups: groups,
		})
		if err != nil {
			return ctx, nil
		}

		return ToContext(ctx, *userInfo), nil
	}
}

// Claims defines some additional data that can be found on the user object.
type Claims struct {
	Groups []string `json:"groups"`
}
