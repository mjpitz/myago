# oidcauth
--
    import "github.com/mjpitz/myago/auth/oidc"


## Usage

#### func  OIDC

```go
func OIDC(cfg Issuer) auth.HandlerFunc
```
OIDC returns a HandlerFunc who authenticates a user with the provided issuer
using an access_token attached to the request. If provided, this access_token is
exchanged for the authenticated user's information. It's important to know that
this function does not handle authorization and requires an additional
HandleFunc to do so.

#### func  ServeMux

```go
func ServeMux(cfg Config, callback TokenCallback) *http.ServeMux
```
ServeMux is some rough code that should allow a command line tool to receive a
token and invoke the provided callback function when a successful exchange is
performed.

#### type ClientConfig

```go
type ClientConfig struct {
	Issuer Issuer `json:"issuer"`
}
```

ClientConfig encapsulates the information needed to establish a client
connection to an identity provider.

#### type Config

```go
type Config struct {
	Issuer       Issuer           `json:"issuer"`
	ClientID     string           `json:"client_id"     usage:"the client_id associated with this service"`
	ClientSecret string           `json:"client_secret" usage:"the client_secret associated with this service"`
	RedirectURL  string           `json:"redirect_url"  usage:"the redirect_url used by this service to obtain a token"`
	Scopes       *cli.StringSlice `json:"scopes"        usage:"specify the scopes that this authorization requires"     default:"openid,profile,email"`
}
```

Config defines the information needed for an application to obtain an identity
token from a provider.

#### type Issuer

```go
type Issuer struct {
	ServerURL            string `json:"server_url"            usage:"the address of the server where user authentication is performed"`
	CertificateAuthority string `json:"certificate_authority" usage:"path pointing to a file containing the certificate authority data for the server"`
}
```

Issuer defines data needed to establish a connection to an issuer.

#### func (Issuer) Provider

```go
func (i Issuer) Provider(ctx context.Context) (*oidc.Provider, error)
```

#### type TokenCallback

```go
type TokenCallback func(token *oauth2.Token)
```

TokenCallback is invoked by the OIDCServeMux endpoint when we've successfully
received and validated the authenticated user session.
