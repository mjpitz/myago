package lifecycle

import (
	"context"
)

var systemLifeCycle = &LifeCycle{}

// Defer will enqueue a function that will be invoked by Resolve.
func Defer(fn func(ctx context.Context)) {
	systemLifeCycle.Defer(fn)
}

// Resolve will process all functions that have been enqueued by Defer up until this point.
func Resolve(ctx context.Context) {
	systemLifeCycle.Resolve(ctx)
}

// Setup initializes a shutdown hook that cancels the underlying context.
func Setup(ctx context.Context) context.Context {
	return systemLifeCycle.Setup(ctx)
}
