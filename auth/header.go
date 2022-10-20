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
	"strings"

	"go.pitz.tech/lib/headers"
)

const authorization = "authorization"

// Get retrieves the current authorization value from the header.
func Get(header headers.Header, expectedScheme string) (string, error) {
	value := header.Get(authorization)
	if value == "" {
		return "", ErrUnauthorized
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) < 2 {
		return "", ErrUnauthorized
	}

	if !strings.EqualFold(parts[0], expectedScheme) {
		return "", ErrUnauthorized
	}

	return parts[1], nil
}
