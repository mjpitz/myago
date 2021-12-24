package auth

import (
	"context"
)

// HandlerFunc defines a common way to add authentication / authorization to a Golang context.
type HandlerFunc func(ctx context.Context) (context.Context, error)

// Composite returns a HandlerFunc that iterates all provided HandlerFunc until the end or an error occurs.
func Composite(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		var err error
		for _, handler := range handlers {
			ctx, err = handler(ctx)
			if err != nil {
				return nil, err
			}
		}

		return ctx, nil
	}
}

// Required returns a HandlerFunc that ensures user information is present on the context.
func Required() HandlerFunc {
	return func(ctx context.Context) (context.Context, error) {
		userInfo := Extract(ctx)
		if userInfo == nil {
			return nil, ErrUnauthorized
		}
		return ctx, nil
	}
}
