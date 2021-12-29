# auth

Package auth provides common code for handling user authentication in a rather
implementation agnostic way. Currently, we only provide basic auth backed by a
CSV, but most components contain an interface that _should_ make it rather easy
to swap out implementations.

```go
import github.com/mjpitz/myago/auth
```

## Usage

```go
var ErrUnauthorized = errors.New("unauthorized")
```

ErrUnauthorized is returned when no user information is found on a context.

#### func Get

```go
func Get(header headers.Header, expectedScheme string) (string, error)
```

Get retrieves the current authorization value from the header.

#### func ToContext

```go
func ToContext(ctx context.Context, userInfo UserInfo) context.Context
```

ToContext attaches the provided UserInfo to the context.

#### type Config

```go
type Config struct {
	AuthType string `json:"auth_type" usage:"configure the user authentication type to use"`
}
```

Config defines a general configuration structure used configure which
authentication is enabled.

#### type HandlerFunc

```go
type HandlerFunc func(ctx context.Context) (context.Context, error)
```

HandlerFunc defines a common way to add authentication / authorization to a
Golang context.

#### func Composite

```go
func Composite(handlers ...HandlerFunc) HandlerFunc
```

Composite returns a HandlerFunc that iterates all provided HandlerFunc until the
end or an error occurs.

#### func Required

```go
func Required() HandlerFunc
```

Required returns a HandlerFunc that ensures user information is present on the
context.

#### type UserInfo

```go
type UserInfo struct {
	// Subject is the users ID. CACF1875-7B44-4B77-BF52-51A06E52FFDF
	Subject string `json:"sub"`
	// Profile is the users name. "Jane Doe"
	Profile string `json:"profile"`
	// Email is the users' email address. jane@example.com
	Email string `json:"email"`
	// EmailVerified indicates if the user has verified their email address.
	EmailVerified bool `json:"email_verified"`
	// Groups contains a list of groups that the user belongs to.
	Groups []string `json:"groups"`
}
```

UserInfo represents a minimum set of user information.

#### func Extract

```go
func Extract(ctx context.Context) *UserInfo
```

Extract attempts to obtain the UserInfo from the provided context.

#### func (UserInfo) Claims

```go
func (u UserInfo) Claims(v interface{}) error
```

Claims provides a convenient way to read additional data from the request.

#### func (\*UserInfo) UnmarshalJSON

```go
func (u *UserInfo) UnmarshalJSON(data []byte) error
```

UnmarshalJSON transparently unmarshals the user information structure.
