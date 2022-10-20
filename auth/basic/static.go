// Copyright (C) 2022 Mya Pitzeruse
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
	"crypto/sha256"
	"encoding/base32"
	"strings"

	"go.pitz.tech/lib/auth"
)

// Static returns an auth.HandlerFunc that uses a static username/password for the system.
func Static(username, password string, groups ...string) auth.HandlerFunc {
	return Basic(&static{
		username: username,
		password: password,
		groups:   groups,
	})
}

type static struct {
	username string
	password string
	groups   []string
}

func (s *static) Lookup(req LookupRequest) (resp LookupResponse, err error) {
	if req.User != s.username {
		err = ErrNotFound
	} else {
		hash := sha256.Sum256([]byte(s.username))

		resp = LookupResponse{
			UserID:   base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:]),
			User:     s.username,
			Password: s.password,
			Groups:   s.groups,
		}

		if strings.Contains(s.username, "@") {
			resp.Email = s.username
		}
	}

	return
}

var _ Store = &static{}
