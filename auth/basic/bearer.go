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

	"github.com/mjpitz/myago/auth"
	"github.com/mjpitz/myago/headers"
)

// Bearer returns a handler func that translates bearer tokens into user information.
func Bearer(store Store) auth.HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		header := headers.Extract(ctx)

		token, err := auth.Get(header, "bearer")
		if err != nil {
			return ctx, nil
		}

		resp, err := store.Lookup(LookupRequest{
			Token: token,
		})

		if err != nil {
			return ctx, nil
		}

		userInfo := auth.UserInfo{
			Subject:       resp.UserID,
			Profile:       resp.User,
			Email:         resp.Email,
			EmailVerified: resp.EmailVerified,
			Groups:        resp.Groups,
		}

		return auth.ToContext(ctx, userInfo), nil
	}
}
