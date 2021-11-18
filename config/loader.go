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

package config

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mjpitz/myago/encoding"
	"github.com/mjpitz/myago/vfs"
)

// DefaultLoader provides a default Loader implementation that supports reading a variety of files.
var DefaultLoader = Loader{
	".json":     encoding.JSON,
	".prototxt": encoding.ProtoText,
	".ptxt":     encoding.ProtoText,
	".toml":     encoding.TOML,
	".yaml":     encoding.YAML,
	".yml":      encoding.YAML,
	".xml":      encoding.XML,
}

// Loader provides functionality for reading a variety of file formats into a struct.
type Loader map[string]*encoding.Encoding

func (l Loader) load(ctx context.Context, v interface{}, filePath string) error {
	fs := vfs.Extract(ctx)

	ext := filepath.Ext(filePath)
	enc, recognized := l[ext]

	exists, err := afero.Exists(fs, filePath)

	switch {
	case err != nil:
		return errors.Wrap(err, "encountered unexpected error")
	case !exists:
		return errors.Wrap(ErrFileDoesNotExist, "file does not exist")
	case ext == "":
		return errors.Wrap(ErrFileMissingExtension, "file missing extensions")
	case !recognized:
		return errors.Wrap(ErrUnsupportedFileExtension, "unsupported file extension")
	}

	file, err := fs.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	return enc.Decoder(file).Decode(v)
}

// Load reads the provided files (if they exist) and unmarshals the data into the provided interface.
func (l Loader) Load(ctx context.Context, v interface{}, filePaths ...string) error {
	for _, filePath := range filePaths {
		err := l.load(ctx, v, filePath)

		switch {
		case errors.Is(err, ErrFileDoesNotExist):
			continue
		case errors.Is(err, ErrFileMissingExtension):
			continue
		case errors.Is(err, ErrUnsupportedFileExtension):
			continue
		case err != nil:
			return errors.Wrapf(err, "failed to load %s", filePath)
		}
	}

	return nil
}

// Load provides a convenience function for being able to load configuration using the DefaultLoader.
func Load(ctx context.Context, v interface{}, filePaths ...string) error {
	return DefaultLoader.Load(ctx, v, filePaths...)
}
