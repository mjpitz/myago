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

package headers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/headers"
)

func TestHTTP(t *testing.T) {
	t.Parallel()

	delegate := func(writer http.ResponseWriter, request *http.Request) {
		h := headers.Extract(request.Context())
		require.Len(t, h, 3)

		require.Equal(t, "val-1", h.Get("Test-1"))
		require.Equal(t, "val-2", h.Get("Test-2"))
		require.Equal(t, "val-3", h.Get("Test-3"))
	}

	handler := headers.HTTP(http.HandlerFunc(delegate))

	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	r.Header.Set("Test-1", "val-1")
	r.Header.Set("Test-2", "val-2")
	r.Header.Set("Test-3", "val-3")

	handler(nil, r)
}
