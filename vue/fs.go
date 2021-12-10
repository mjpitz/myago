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

package vue

import (
	"io"
	"os"
)

// File describes what operations an associated File should be able to perform.
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

// FileSystem describes what the file system should look like.
type FileSystem interface {
	Open(string) (File, error)
}

// Wrap creates a new FileSystem that supports server side loading for VueJS applications.
func Wrap(delegate FileSystem) *fs {
	return &fs{delegate}
}

type fs struct {
	delegate FileSystem
}

func (v *fs) Open(name string) (File, error) {
	f, err := v.delegate.Open(name)

	// if it doesn't exist on the server, delegate to the front-end
	// only exception _should_ be favicon.ico which is auto-fetched by browsers
	if name != "/favicon.ico" && f == nil {
		return v.delegate.Open("/index.html")
	}

	return f, err
}
