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