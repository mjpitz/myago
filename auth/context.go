package auth

import (
	"context"

	"github.com/mjpitz/myago"
)

const contextKey = myago.ContextKey("auth")

// ToContext attaches the provided UserInfo to the context.
func ToContext(ctx context.Context, userInfo UserInfo) context.Context {
	return context.WithValue(ctx, contextKey, &userInfo)
}

// Extract attempts to obtain the UserInfo from the provided context.
func Extract(ctx context.Context) *UserInfo {
	v := ctx.Value(contextKey)
	if v == nil {
		return nil
	}

	userInfo, ok := v.(*UserInfo)
	if !ok {
		return nil
	}

	return userInfo
}
