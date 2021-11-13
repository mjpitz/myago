package livetls

import (
	"time"
)

// Config defines common configuration that can be used to LoadCertificates for encrypted communication.
type Config struct {
	Enable          bool          `json:"enable"           usage:"whether or not TLS should be enabled"`
	CertPath        string        `json:"cert_path"        usage:"where to locate certificates for communication"`
	CAFile          string        `json:"ca_file"          usage:"override the ca file name"   default:"ca.pem"`
	CertFile        string        `json:"cert_file"        usage:"override the cert file name" default:"cert.pem"`
	KeyFile         string        `json:"key_file"         usage:"override the key file name"  default:"key.pem"`
}

type ReloadingConfig struct {
	Config
	ReloadInterval time.Duration `json:"reload_interval" usage:"configures how often the config is checked for reload"`
}
