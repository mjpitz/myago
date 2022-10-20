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

package authors_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/authors"
)

const exampleAuthorsFile = `
# This is the list of significant contributors.
#
# This does not necessarily list everyone who has contributed code,
# especially since many employees of one corporation may be contributing.
# To see the full list of contributors, see the revision history in
# source control.
Name
First Last <noreply@example.com>
`

func TestParse(t *testing.T) {
	t.Parallel()

	emptyAuthors := authors.Parse("")
	require.Len(t, emptyAuthors, 0)

	exampleAuthors := authors.Parse(exampleAuthorsFile)
	require.Len(t, exampleAuthors, 2)

	require.Equal(t, "Name", exampleAuthors[0].Name)
	require.Equal(t, "", exampleAuthors[0].Email)

	require.Equal(t, "First Last", exampleAuthors[1].Name)
	require.Equal(t, "noreply@example.com", exampleAuthors[1].Email)
}
