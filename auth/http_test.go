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

package auth_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/auth"
	"github.com/mjpitz/myago/headers"
)

type mockStore struct{}

func (m *mockStore) Lookup(username string) (password string, groups []string, err error) {
	if username != "admin" {
		return "", nil, auth.ErrUnauthorized
	}

	return "badadmin", []string{"admin"}, nil
}

var _ auth.Store = &mockStore{}

func TestHTTP(t *testing.T) {
	t.Parallel()

	called := false

	delegate := func(w http.ResponseWriter, r *http.Request) {
		user := auth.Extract(r.Context())
		require.NotNil(t, user)

		require.Equal(t, "admin", user.Subject)
		require.Equal(t, "admin", user.Profile)
		called = true
	}

	handler := auth.HTTP(
		http.HandlerFunc(delegate),
		auth.Basic(&mockStore{}),
		auth.Required(),
	)

	handler = headers.HTTP(handler)

	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	authentication := base64.StdEncoding.EncodeToString([]byte("admin:badadmin"))
	r.Header.Set("Authorization", "Basic "+authentication)

	handler(nil, r)
	require.True(t, called)
}
