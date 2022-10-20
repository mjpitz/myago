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

package vfs

import (
	"context"

	"github.com/spf13/afero"
	"go.pitz.tech/lib/libctx"
)

// FS provides a file system abstraction.
type FS = afero.Fs

var contextKey = libctx.Key("vfs")

var defaultFS = afero.NewOsFs()

// Extract pulls the file system from the provided context. If no file system is found, then the defaultFS is returned.
func Extract(ctx context.Context) FS {
	val := ctx.Value(contextKey)
	if val == nil {
		return defaultFS
	}

	return val.(FS)
}

// ToContext sets the file system on the provided context.
func ToContext(ctx context.Context, fs FS) context.Context {
	return context.WithValue(ctx, contextKey, fs)
}
