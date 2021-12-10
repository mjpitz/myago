# vue
--
    import "github.com/mjpitz/myago/vue"

Package vue contains some helper code for VueJS frontends. The FileSystem
constructed by Wrap enables use of the web router, eliminating the need for
fragments in the application layer.

## Usage

#### func  Wrap

```go
func Wrap(delegate FileSystem) *fs
```
Wrap creates a new FileSystem that supports server side loading for VueJS
applications.

#### type File

```go
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}
```

File describes what operations an associated File should be able to perform.

#### type FileSystem

```go
type FileSystem interface {
	Open(string) (File, error)
}
```

FileSystem describes what the file system should look like.
