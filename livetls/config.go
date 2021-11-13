package livetls

import (
	"time"
)

// Config defines common configuration that can be used to LoadCertificates for encrypted communication.
type Config struct {
	Enable         bool          `json:"enable"          usage:"whether or not TLS should be enabled"`
	CertPath       string        `json:"cert_path"       usage:"where to locate certificates for communication"`
	CAFile         string        `json:"ca_file"         usage:"override the ca file name"      default:"ca.crt"`
	CertFile       string        `json:"cert_file"       usage:"override the cert file name"    default:"tls.crt"`
	KeyFile        string        `json:"key_file"        usage:"override the key file name"     default:"tls.key"`
	ReloadInterval time.Duration `json:"reload_interval" usage:"how often to reload the config" default:"5m"`
}
