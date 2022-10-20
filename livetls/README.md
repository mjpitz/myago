# livetls

```go
import go.pitz.tech/lib/livetls
```

## Usage

#### func LoadCertPool

```go
func LoadCertPool(ctx context.Context, cfg *Config) (*x509.CertPool, error)
```

LoadCertPool loads the x509 certificate authority pool.

#### func LoadCertificate

```go
func LoadCertificate(ctx context.Context, cfg *Config) (*tls.Certificate, error)
```

LoadCertificate loads the certificate from the configured public/private key.

#### func New

```go
func New(ctx context.Context, config Config) (*tls.Config, error)
```

New construct a tls.Config that will periodically reload the configured
certificate.

#### type Config

```go
type Config struct {
	Enable         bool          `json:"enable"          usage:"whether or not TLS should be enabled"`
	CertPath       string        `json:"cert_path"       usage:"where to locate certificates for communication"`
	CAFile         string        `json:"ca_file"         usage:"override the ca file name"      default:"ca.crt"`
	CertFile       string        `json:"cert_file"       usage:"override the cert file name"    default:"tls.crt"`
	KeyFile        string        `json:"key_file"        usage:"override the key file name"     default:"tls.key"`
	ReloadInterval time.Duration `json:"reload_interval" usage:"how often to reload the config" default:"5m"`
}
```

Config defines common configuration that can be used to LoadCertificates for
encrypted communication.
