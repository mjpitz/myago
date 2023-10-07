// Copyright (C) 2023 Mya Pitzeruse
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

package build

import (
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// Info defines common build information associated with the binary.
type Info struct {
	OS           string
	Architecture string

	GoVersion  string
	CGoEnabled bool

	Version  string
	VCS      string
	Revision string
	Compiled time.Time
	Modified bool
}

// Metadata formats the underlying information as a generalized key-value map.
func (info Info) Metadata() map[string]any {
	return map[string]any{
		"os":   info.OS,
		"arch": info.Architecture,
		"go":   info.GoVersion,
		"cgo":  strconv.FormatBool(info.CGoEnabled),
		"vcs":  info.VCS,
		"rev":  info.Revision,
		"time": info.Compiled.Format(time.RFC3339),
		"mod":  strconv.FormatBool(info.Modified),
	}
}

// ParseInfo extracts as much build information from the compiled binary as it can.
func ParseInfo() (info Info) {
	info.OS = runtime.GOOS
	info.Architecture = runtime.GOARCH
	info.GoVersion = strings.TrimPrefix(runtime.Version(), "go")
	info.Compiled = time.Now()

	build, ok := debug.ReadBuildInfo()
	if ok {
		info.Version = build.Main.Version

		for _, setting := range build.Settings {
			switch setting.Key {
			case "CGO_ENABLED":
				info.CGoEnabled, _ = strconv.ParseBool(setting.Value)
			case "vcs":
				info.VCS = setting.Value
			case "vcs.revision":
				info.Revision = setting.Value
			case "vcs.time":
				info.Compiled, _ = time.Parse(time.RFC3339, setting.Value)
			case "vcs.modified":
				info.Modified, _ = strconv.ParseBool(setting.Value)
			}
		}
	}

	return info
}
