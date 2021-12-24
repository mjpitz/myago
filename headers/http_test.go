package headers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/headers"
)

func TestHTTP(t *testing.T) {
	delegate := func(writer http.ResponseWriter, request *http.Request) {
		h := headers.Extract(request.Context())
		require.Len(t, h, 3)

		require.Equal(t, "val-1", h.Get("Test-1"))
		require.Equal(t, "val-2", h.Get("Test-2"))
		require.Equal(t, "val-3", h.Get("Test-3"))
	}

	handler := headers.HTTP(http.HandlerFunc(delegate))

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	r.Header.Set("Test-1", "val-1")
	r.Header.Set("Test-2", "val-2")
	r.Header.Set("Test-3", "val-3")

	handler(nil, r)
}
