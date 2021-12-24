package headers

import (
	"context"

	"github.com/mjpitz/myago"
)

const contextKey = myago.ContextKey("headers")

// ToContext attaches the provided headers to the context.
func ToContext(ctx context.Context, header Header) context.Context {
	return context.WithValue(ctx, contextKey, header)
}

// Extract attempts to obtain the headers from the provided context.
func Extract(ctx context.Context) Header {
	v := ctx.Value(contextKey)
	if v == nil {
		return New()
	}

	return v.(Header)
}
