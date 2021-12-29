# vfs

Package vfs provides utilities for managing virtual file systems on contexts to
avoid direct calls to the built-in `os` interface. This is particularly useful
for testing. Currently, this wraps the afero virtual file system which provides
OS and in memory implementations.

```go
import github.com/mjpitz/myago/vfs
```

## Usage

#### func ToContext

```go
func ToContext(ctx context.Context, fs FS) context.Context
```

ToContext sets the file system on the provided context.

#### type FS

```go
type FS = afero.Fs
```

FS provides a file system abstraction.

#### func Extract

```go
func Extract(ctx context.Context) FS
```

Extract pulls the file system from the provided context. If no file system is
found, then the defaultFS is returned.
