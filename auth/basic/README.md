# basicauth




```go
import github.com/mjpitz/myago/auth/basic
```

## Usage

```go
var ErrBadRequest = errors.New("bad lookup request")
```
ErrBadRequest is returned when a lookup request does not contain a required
field.

```go
var ErrNotFound = errors.New("not found")
```
ErrNotFound is returned when a credential is not found.

#### func  Basic

```go
func Basic(store Store) auth.HandlerFunc
```
Basic implements a basic access authentication handler function.

#### func  Bearer

```go
func Bearer(store Store) auth.HandlerFunc
```
Bearer returns a handler func that translates bearer tokens into user
information.

#### func  Handler

```go
func Handler(ctx context.Context, cfg Config) (auth.HandlerFunc, error)
```
Handler returns the appropriate handler based on the provided configuration.

#### type AccessToken

```go
type AccessToken struct {
	Token string `json:"token" usage:"the access token used to authenticate requests"`
}
```

AccessToken is used to authenticate a user using a bearer token.

#### type ClientConfig

```go
type ClientConfig struct {
	UsernamePassword
	AccessToken
}
```

ClientConfig defines the options available to a client.

#### func (ClientConfig) Token

```go
func (c ClientConfig) Token() (*oauth2.Token, error)
```

#### type Config

```go
type Config struct {
	PasswordFile string `json:"password_file" usage:"path to the csv file containing usernames and passwords"`
	TokenFile    string `json:"token_file" usage:"path to the csv file containing tokens"`
}
```

Config defines the options available to a server.

#### type LazyStore

```go
type LazyStore struct {
	Provider func() (Store, error)
}
```

LazyStore provides a convenient way to lazily load an underlying store.

#### func (*LazyStore) Lookup

```go
func (c *LazyStore) Lookup(req LookupRequest) (resp LookupResponse, err error)
```

#### type LookupRequest

```go
type LookupRequest struct {
	User  string
	Token string
}
```


#### type LookupResponse

```go
type LookupResponse struct {
	UserID string
	User   string
	Groups []string

	Email         string
	EmailVerified bool

	// one of these will be set based on the LookupRequest
	Password string
	Token    string
}
```


#### type Store

```go
type Store interface {
	// Lookup retrieves the provided user's password and groups.
	Lookup(req LookupRequest) (resp LookupResponse, err error)
}
```

Store defines an abstraction for loading user credentials.

#### func  OpenCSV

```go
func OpenCSV(ctx context.Context, fileName string) (Store, error)
```
OpenCSV attempts to open the provided csv file and return a parsed index based
on the contents.

#### type UsernamePassword

```go
type UsernamePassword struct {
	Username string `json:"username" usage:"the username to login with"`
	Password string `json:"password" usage:"the password associated with the username"`
}
```

UsernamePassword is used to authenticate a user using a username and password.
