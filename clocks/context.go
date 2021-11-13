package clocks

import (
	"context"

	"github.com/jonboulle/clockwork"

	"github.com/mjpitz/myago"
)

var contextKey = myago.ContextKey("clocks")

var defaultClock = clockwork.NewRealClock()

// Extract pulls the clock from the provided context. If no clock is found, then the defaultClock is returned.
func Extract(ctx context.Context) clockwork.Clock {
	clock := ctx.Value(contextKey)
	if clock == nil {
		return defaultClock
	}

	return clock.(clockwork.Clock)
}

// ToContext sets the clock on the provided context.
func ToContext(ctx context.Context, clock clockwork.Clock) context.Context {
	return context.WithValue(ctx, contextKey, clock)
}

// Setup sets the defaultClock on the provided context. This can always be overridden later.
func Setup(ctx context.Context) context.Context {
	return ToContext(ctx, defaultClock)
}
