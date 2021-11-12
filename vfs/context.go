package vfs

import (
	"context"

	"github.com/spf13/afero"
)

// FS provides a file system abstraction.
type FS = afero.Fs

var contextKey = &struct{}{}

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
