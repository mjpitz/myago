# dirset
--
    import "github.com/mjpitz/myago/dirset"

Package dirset provides discovery of common application directories for things
like caching, locks, and logs.

## Usage

#### type DirectorySet

```go
type DirectorySet struct {
	CacheDir      string
	StateDir      string
	LocalStateDir string
	LockDir       string
	LogDir        string
}
```

DirectorySet defines a common set of paths that applications can use for a
variety of reasons.

#### func  Local

```go
func Local(appName string) (DirectorySet, error)
```
Local uses information about the local system to determine a common set of paths
for use by the application.
