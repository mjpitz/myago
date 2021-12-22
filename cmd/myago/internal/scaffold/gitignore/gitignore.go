package gitignore

import (
	"embed"
	"fmt"
)

//go:generate sh sync.sh

//go:embed *.gitignore
var licenses embed.FS

// Get returns the gitignore associated with the ID.
func Get(id string) (string, bool) {
	body, err := licenses.ReadFile(fmt.Sprintf("%s.gitignore", id))
	if err != nil {
		return "", false
	}

	return string(body), true
}
