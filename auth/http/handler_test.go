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

package httpauth_test

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/auth"
	"github.com/mjpitz/myago/auth/basic"
	httpauth "github.com/mjpitz/myago/auth/http"
	"github.com/mjpitz/myago/headers"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	called := false

	delegate := func(w http.ResponseWriter, r *http.Request) {
		user := auth.Extract(r.Context())
		require.NotNil(t, user)

		require.Equal(t, "userID", user.Subject)
		require.Equal(t, "username", user.Profile)
		require.Equal(t, "", user.Email)
		require.Equal(t, false, user.EmailVerified)
		require.Equal(t, []string{"group1", "group2"}, user.Groups)
		called = true
	}

	authHandler, err := basicauth.Handler(context.Background(), basicauth.Config{
		PasswordFile: filepath.Join("..", "basic", "testdata", "basic.csv"),
	})
	require.NoError(t, err)

	handler := httpauth.Handler(
		http.HandlerFunc(delegate),
		authHandler,
		auth.Required(),
	)

	handler = headers.HTTP(handler)

	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	r.Header.Set("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")

	handler(nil, r)
	require.True(t, called)
}
