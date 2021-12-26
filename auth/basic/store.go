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

package basicauth

import (
	"sync"
)

type LookupRequest struct {
	User  string
	Token string
}

type LookupResponse struct {
	UserID string
	User   string
	Groups []string

	Email         string
	EmailVerified bool

	// one of these will be set based on the LookupRequest
	Password string
	Token    string
}

// Store defines an abstraction for loading user credentials.
type Store interface {
	// Lookup retrieves the provided user's password and groups.
	Lookup(req LookupRequest) (resp LookupResponse, err error)
}

// store provides an in-memory index for looking up passwords and groups for a named user.
type store struct {
	idx map[string]*entry
}

func (s *store) Lookup(req LookupRequest) (resp LookupResponse, err error) {
	switch {
	case len(req.Token) > 0:
		entry := s.idx[req.Token]
		if entry == nil {
			err = ErrNotFound
		} else {
			resp = LookupResponse{
				UserID: entry.userID,
				User:   entry.f0,
				Groups: entry.groups,
				Token:  req.Token,
			}
		}
	case len(req.User) > 0:
		entry := s.idx[req.User]
		if entry == nil {
			err = ErrNotFound
		} else {
			resp = LookupResponse{
				UserID:   entry.userID,
				Password: entry.f0,
				Groups:   entry.groups,
				User:     req.User,
			}
		}
	default:
		err = ErrBadRequest
	}

	return resp, err
}

var _ Store = &store{}

type entry struct {
	f0     string // username for tokens, password for basic auth
	f1     string // token for tokens, username for basic auth
	userID string
	groups []string
}

// LazyStore provides a convenient way to lazily load an underlying store.
type LazyStore struct {
	mu       sync.Mutex
	inst     Store
	Provider func() (Store, error)
}

func (c *LazyStore) Lookup(req LookupRequest) (resp LookupResponse, err error) {
	inst, err := func() (Store, error) {
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.inst == nil {
			c.inst, err = c.Provider()
			if err != nil {
				return nil, err
			}
		}

		return c.inst, nil
	}()

	if err != nil {
		return resp, err
	}

	return inst.Lookup(req)
}
