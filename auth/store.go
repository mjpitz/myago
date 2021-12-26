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

package auth

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

func (c *store) Lookup(req LookupRequest) (resp LookupResponse, err error) {
	switch {
	case len(req.Token) > 0:
		entry := c.idx[req.Token]
		if entry == nil {
			err = errNotFound
			return
		} else {
			resp = LookupResponse{
				UserID: entry.userID,
				User:   entry.f0,
				Groups: entry.groups,
				Token:  req.Token,
			}
		}
	case len(req.User) > 0:
		entry := c.idx[req.User]
		if entry == nil {
			err = errNotFound
		} else {
			resp = LookupResponse{
				UserID:   entry.userID,
				Password: entry.f0,
				Groups:   entry.groups,
				User:     req.User,
			}
		}
	default:
		err = errBadRequest
	}

	return
}

var _ Store = &store{}

type entry struct {
	f0     string // username for tokens, password for basic auth
	f1     string // token for tokens, username for basic auth
	userID string
	groups []string
}
