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

package headers

import (
	"strings"
)

// New constructs a Header for use.
func New() Header {
	return make(Header)
}

// Header defines an abstract definition of a header.
type Header map[string][]string

// SetAll sets the values for the provides key.
func (h Header) SetAll(key string, values []string) {
	h[strings.ToLower(key)] = values
}

// Set sets a single value for the provided key.
func (h Header) Set(key, value string) {
	h.SetAll(key, []string{value})
}

// GetAll returns all possible values for a key.
func (h Header) GetAll(key string) []string {
	return h[strings.ToLower(key)]
}

// Get returns the first possible header value for a key (if present).
func (h Header) Get(key string) string {
	all := h.GetAll(key)
	if len(all) > 0 {
		return all[0]
	}

	return ""
}
