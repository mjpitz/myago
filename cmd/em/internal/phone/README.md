# phone




```go
import go.pitz.tech/lib/cmd/em/internal/phone
```

## Usage

#### type Encoder

```go
type Encoder struct {
}
```


#### func  NewEncoder

```go
func NewEncoder(writer io.Writer) *Encoder
```
NewEncoder returns an encoder that translates data into a phone code.

#### func (*Encoder) Write

```go
func (e *Encoder) Write(p []byte) (n int, err error)
```
