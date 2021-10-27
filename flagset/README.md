# flagset
--
    import "github.com/mjpitz/myago/flagset"

Package flagset provides logic for creating a flagset from a struct.

## Usage

#### func  Extract

```go
func Extract(v interface{}) []cli.Flag
```
Extract parses the provided object to create a flagset.

#### func  ExtractPrefix

```go
func ExtractPrefix(prefix string, v interface{}) []cli.Flag
```
ExtractPrefix parses the provided to create a flagset with the provided prefix.
