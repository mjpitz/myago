package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/mjpitz/myago/headers"
)

// Basic implements a basic access authentication handler function. It parses values from the headers to obtain info
// about the authenticated user.
func Basic(store Store) HandlerFunc {
	basicPrefix := "Basic "

	return func(ctx context.Context) (context.Context, error) {
		header := headers.Extract(ctx)

		authentication := header.Get("Authorization")
		if !strings.HasPrefix(authentication, basicPrefix) {
			return ctx, nil
		}

		authentication = strings.TrimPrefix(authentication, basicPrefix)

		decoded, err := base64.StdEncoding.DecodeString(authentication)
		if err != nil {
			return ctx, nil
		}

		parts := strings.Split(string(decoded), ":")
		if len(parts) < 2 {
			return ctx, nil
		}

		username := parts[0]
		provided := parts[1]

		password, groups, err := store.Lookup(username)
		if err != nil {
			return ctx, nil
		}

		if provided != password {
			return ctx, nil
		}

		userInfo := &UserInfo{
			Subject: username,
			Profile: username,
		}

		err = userInfo.WithExtra(&Claims{
			Groups: groups,
		})
		if err != nil {
			return ctx, nil
		}

		return ToContext(ctx, *userInfo), nil
	}
}

// Claims defines some additional data that can be found on the user object.
type Claims struct {
	Groups []string `json:"groups"`
}
