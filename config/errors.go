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
	"errors"
)

var (
	// ErrFileDoesNotExist is returned when the file we're interacting with does not exist.
	ErrFileDoesNotExist = errors.New("file does not exist")

	// ErrFileMissingExtension is returned when the provided file is missing an extension.
	ErrFileMissingExtension = errors.New("file missing extension")

	// ErrUnsupportedFileExtension is returned when we don't recognize a given file extension.
	ErrUnsupportedFileExtension = errors.New("unsupported file extension")
)
