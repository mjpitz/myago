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
	"errors"
	"net/http"
)

// HTTP returns an http middleware function that invokes the provided auth handlers.
func HTTP(delegate http.Handler, handlers ...HandlerFunc) http.HandlerFunc {
	handler := Composite(handlers...)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := handler(r.Context())

		switch {
		case errors.Is(err, ErrUnauthorized):
			http.Error(w, "", http.StatusUnauthorized)
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			delegate.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
