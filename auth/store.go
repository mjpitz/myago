package auth

import (
	"fmt"
)

// Store defines an abstraction for loading user credentials.
type Store interface {
	// Lookup retrieves the provided user's password and groups.
	Lookup(username string) (password string, groups []string, err error)
}

// store provides an in-memory index for looking up passwords and groups for a named user.
type store struct {
	idx map[string]*entry
}

func (c *store) Lookup(username string) (password string, groups []string, err error) {
	entry := c.idx[username]
	if entry == nil {
		return "", nil, fmt.Errorf("not found")
	}

	return entry.password, entry.groups, nil
}

var _ Store = &store{}

type entry struct {
	password string
	groups   []string
}
