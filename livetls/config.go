// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
