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
