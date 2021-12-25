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
	"context"
	"encoding/csv"
	"errors"
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
		if errors.Is(err, io.EOF) {
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
