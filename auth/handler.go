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
)

// HandlerFunc defines a common way to add authentication / authorization to a Golang context.
type HandlerFunc func(ctx context.Context) (context.Context, error)

// Composite returns a HandlerFunc that iterates all provided HandlerFunc until the end or an error occurs.
func Composite(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		var err error
		for _, handler := range handlers {
			ctx, err = handler(ctx)
			if err != nil {
				return nil, err
			}
		}

		return ctx, nil
	}
}

// Required returns a HandlerFunc that ensures user information is present on the context.
func Required() HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		userInfo := Extract(ctx)
		if userInfo == nil {
			return nil, ErrUnauthorized
		}

		return ctx, nil
	}
}
