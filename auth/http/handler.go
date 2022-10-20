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

package httpauth

import (
	"errors"
	"net/http"

	"go.pitz.tech/lib/auth"
)

// Handler returns an HTTP middleware function that invokes the provided auth handlers.
func Handler(delegate http.Handler, handlers ...auth.HandlerFunc) http.HandlerFunc {
	handler := auth.Composite(handlers...)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := handler(r.Context())

		switch {
		case errors.Is(err, auth.ErrUnauthorized):
			http.Error(w, "", http.StatusUnauthorized)
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			delegate.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
