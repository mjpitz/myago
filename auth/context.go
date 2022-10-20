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

	"go.pitz.tech/lib/libctx"
)

const contextKey = libctx.Key("auth")

// ToContext attaches the provided UserInfo to the context.
func ToContext(ctx context.Context, userInfo UserInfo) context.Context {
	return context.WithValue(ctx, contextKey, &userInfo)
}

// Extract attempts to obtain the UserInfo from the provided context.
func Extract(ctx context.Context) *UserInfo {
	v := ctx.Value(contextKey)
	if v == nil {
		return nil
	}

	return v.(*UserInfo)
}
