package auth_test

import (
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

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	authentication := base64.StdEncoding.EncodeToString([]byte("admin:badadmin"))
	r.Header.Set("Authorization", "Basic "+authentication)

	handler(nil, r)
	require.True(t, called)
}
