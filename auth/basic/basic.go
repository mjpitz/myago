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
	"strings"

	"github.com/mjpitz/myago/auth"
	"github.com/mjpitz/myago/headers"
)

// Basic implements a basic access authentication handler function.
func Basic(store Store) auth.HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		header := headers.Extract(ctx)
		authentication, err := auth.Get(header, "basic")
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

		resp, err := store.Lookup(LookupRequest{
			User: username,
		})
		if err != nil {
			return ctx, nil
		}

		if provided != resp.Password {
			return ctx, nil
		}

		userInfo := auth.UserInfo{
			Subject:       resp.UserID,
			Profile:       username,
			Email:         resp.Email,
			EmailVerified: resp.EmailVerified,
			Groups:        resp.Groups,
		}

		return auth.ToContext(ctx, userInfo), nil
	}
}
