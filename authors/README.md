# authors
--
    import "github.com/mjpitz/myago/authors"


## Usage

#### func  Parse

```go
func Parse(contents string) []*cli.Author
```
Parse parses the contents of an AUTHORS file. An AUTHORS file is a plaintext
file whose contents details the primary contributors to the project. Each line
in the file contains either a comment (denoted by the pound symbol, "#") or an
author. Each author line should contain the name of the contributor and an
optional email address. The format for and author line is "name <email>". For
more information, see the AUTHORS file for this project.
