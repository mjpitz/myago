# dirset

Package dirset provides discovery of common application directories for things
like caching, locks, and logs.

```go
import github.com/mjpitz/myago/dirset
```

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

#### func Local

```go
func Local(appName string) (DirectorySet, error)
```

Local uses information about the local system to determine a common set of paths
for use by the application.

#### func Must

```go
func Must(appName string) DirectorySet
```

Must panics if Local fails to obtain the directory set.
