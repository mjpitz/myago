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

package licenses

import (
	"embed"
	"fmt"
)

//go:generate sh sync.sh

//go:embed *.txt
var licenses embed.FS

// Get retrieves the associated license text.
func Get(spdxID string) (string, bool) {
	body, err := licenses.ReadFile(fmt.Sprintf("%s.txt", spdxID))
	if err != nil {
		return "", false
	}

	return string(body), true
}

// GetHeader retrieves the header associated with the license text.
func GetHeader(spdxID string) (string, bool) {
	body, err := licenses.ReadFile(fmt.Sprintf("%s-header.txt", spdxID))
	if err != nil {
		return "", false
	}

	return string(body), true
}
