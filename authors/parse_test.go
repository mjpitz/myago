package authors_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/authors"
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
	emptyAuthors := authors.Parse("")
	require.Len(t, emptyAuthors, 0)

	exampleAuthors := authors.Parse(exampleAuthorsFile)
	require.Len(t, exampleAuthors, 2)

	require.Equal(t, "Name", exampleAuthors[0].Name)
	require.Equal(t, "", exampleAuthors[0].Email)

	require.Equal(t, "First Last", exampleAuthors[1].Name)
	require.Equal(t, "noreply@example.com", exampleAuthors[1].Email)
}
