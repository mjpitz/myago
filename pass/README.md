# pass

Package pass provides password derivation functions backing solutions like
Spectre. There are three steps in the process. First, you need to derive an
Identity key based your name and primary password. This key is unique to you
(assuming name / password combinations are unique). The second step is to
generate a SiteKey. This key is unique to you for the site that you're
generating the key for. Finally, the last step is to generate a SitePassword
using the derived SiteKey and associated password format.

```go
import github.com/mjpitz/myago/pass
```

## Usage

#### func Identity

```go
func Identity(scope Scope, password []byte, name string) ([]byte, error)
```

Identity computes your identity which is defined by your root password. This key
unlocks all doors. The result is a cryptographic key that is derived from the
scope of the operation, your password (authentication), and your name
(identification).

#### func SiteKey

```go
func SiteKey(scope Scope, identity []byte, site string, counter uint32) []byte
```

SiteKey derives a site specific key from your identity key. Use of your identity
key ensures only your identity has access to this key and your site name scopes
the key to the site. The site counter ensures you can easily create new keys for
the site should a key become compromised.

#### func SitePassword

```go
func SitePassword(siteKey []byte, class TemplateClass) []byte
```

SitePassword is an identifier derived from your site key in compliance with the
site's password policy. This step renders the sites cryptographic key into a
format that the site's password input will accept.

#### type Scope

```go
type Scope string
```

Scope defines an enumeration of possible scopes used in key derivation.

```go
const (
	// Authentication is used when generating a key that is used for authenticating the user, such as a password.
	Authentication Scope = "com.lyndir.masterpassword"

	// Identification is used when generating a key that is intended for the purpose of identifying the user.
	// Identification keys are not necessarily private.
	Identification Scope = "com.lyndir.masterpassword.login"

	// Recovery is used for generating fallback identifiers for use in access recovery when the primary mechanism has
	// failed.
	Recovery Scope = "com.lyndir.masterpassword.answer"
)
```

#### type TemplateClass

```go
type TemplateClass string
```

TemplateClass defines an enumeration of password templates to choose from.

```go
const (
	// MaximumSecurity defines a set of templates used to generate passwords with the strongest security.
	MaximumSecurity TemplateClass = "max"
	// Long defines a set of templates used to generate long passwords.
	Long TemplateClass = "long"
	// Medium defines a set of templates used to generate medium-length passwords.
	Medium TemplateClass = "medium"
	// Short defines a set of templates used to generate short-length passwords.
	Short TemplateClass = "short"
	// Basic defines a set of templates used to generate basic passwords.
	Basic TemplateClass = "basic"
	// PIN generates a pin.
	PIN TemplateClass = "pin"
	// VerificationCode provides a template for generating a 6-digit verification code.
	VerificationCode TemplateClass = "code"
)
```
