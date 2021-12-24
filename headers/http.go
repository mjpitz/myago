package headers

import (
	"net/http"
)

// HTTP returns an http middleware function that translates HTTP headers into a context Header.
func HTTP(delegate http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		header := Extract(ctx)
		for key, values := range r.Header {
			header.SetAll(key, values)
		}

		delegate.ServeHTTP(w, r.WithContext(ToContext(ctx, header)))
	}
}
