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
