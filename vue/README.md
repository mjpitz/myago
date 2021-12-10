# vue
--
    import "github.com/mjpitz/myago/vue"

Package vue contains some helper code for VueJS frontends. The FileSystem
constructed by Wrap enables use of the web router, eliminating the need for
fragments in the application layer.

## Usage

#### func  Wrap

```go
func Wrap(delegate http.FileSystem) http.FileSystem
```
Wrap creates a new FileSystem that supports server side loading for VueJS
applications.
