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
			return
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delegate.ServeHTTP(w, r.WithContext(ctx))
	}
}
