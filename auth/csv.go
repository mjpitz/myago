package auth

import (
	"context"
	"encoding/csv"
	"io"
	"strings"

	"github.com/mjpitz/myago/vfs"
)

// OpenCSV attempts to open the provided csv file and return a parsed index based on the contents.
func OpenCSV(ctx context.Context, fileName string) (Store, error) {
	fs := vfs.Extract(ctx)
	f, err := fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &store{
		idx: make(map[string]*entry),
	}

	reader := csv.NewReader(f)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		c.idx[record[1]] = &entry{
			password: record[0],
			groups:   strings.Split(record[3], ","),
		}
	}

	return c, nil
}
