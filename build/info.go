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
	"flag"
	"fmt"
	"os"
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

	build, ok := debug.ReadBuildInfo()
	if ok {
		info.Version = build.Main.Version

		devNull := os.NewFile(0, os.DevNull)
		defer func() { _ = devNull.Close() }()

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
			case "-ldflags":
				ldflags := flag.NewFlagSet("ldflags", flag.ErrorHandling(-1))
				ldflags.SetOutput(devNull)
				ldflags.Usage = func() {}

				stringValues := &KVSlice{}
				ldflags.Var(stringValues, "X", "")

				// not actively tracked, but added to avoid parsing errors since they appear in logs, despite /dev/null
				ldflags.String("B", "", "")
				ldflags.String("E", "", "")
				ldflags.String("H", "", "")
				ldflags.String("I", "", "")
				ldflags.String("L", "", "")
				ldflags.String("R", "", "")
				ldflags.String("T", "", "")
				ldflags.Bool("V", false, "")
				ldflags.Bool("a", false, "")
				ldflags.Bool("asan", false, "")
				ldflags.String("buildid", "", "")
				ldflags.String("buildmode", "", "")
				ldflags.Bool("c", false, "")
				ldflags.Bool("compressdwarf", false, "")
				ldflags.String("cpuprofile", "", "")
				ldflags.Bool("d", false, "")
				ldflags.Int("debugtramp", 0, "")
				ldflags.Bool("dumpdep", false, "")
				ldflags.String("extar", "", "")
				ldflags.String("extld", "", "")
				ldflags.String("extldflags", "", "")
				ldflags.Bool("f", false, "")
				ldflags.Bool("g", false, "")
				ldflags.String("importcfg", "", "")
				ldflags.String("installsuffix", "", "")
				ldflags.String("k", "", "")
				ldflags.String("libgcc", "", "")
				ldflags.String("linkmode", "", "")
				ldflags.Bool("linkshared", false, "")
				ldflags.String("memprofile", "", "")
				ldflags.String("memprofilerate", "", "")
				ldflags.Bool("msan", false, "")
				ldflags.Bool("n", false, "")
				ldflags.String("o", "", "")
				ldflags.String("pluginpath", "", "")
				ldflags.String("r", "", "")
				ldflags.Bool("race", false, "")
				ldflags.Bool("s", false, "")
				ldflags.Bool("shared", false, "")
				ldflags.String("tmpdir", "", "")
				ldflags.Bool("u", false, "")
				ldflags.Bool("v", false, "")
				ldflags.Bool("w", false, "")

				_ = ldflags.Parse(strings.Split(setting.Value, " "))

				for _, kv := range *stringValues {
					switch kv.Key {
					case "main.version":
						if info.Version == "" {
							info.Version = kv.Value
						}
					case "main.commit":
						if info.Revision == "" {
							info.Revision = kv.Value
						}
					case "main.date":
						if info.Compiled.IsZero() {
							info.Compiled, _ = time.Parse(time.RFC3339, kv.Value)
						}
					}
				}
			}
		}
	}

	return info
}

type KV struct {
	Key   string
	Value string
}

type KVSlice []KV

func (s *KVSlice) Set(value string) error {
	if value == "" {
		return nil
	}

	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("unexpected  number of parts")
	}

	*s = append(*s, KV{
		Key:   parts[0],
		Value: parts[1],
	})

	return nil
}

func (s *KVSlice) String() string {
	str := ""

	if s != nil {
		for _, kv := range *s {
			if str != "" {
				str += ", "
			}

			str += kv.Key + "=" + kv.Value
		}
	}

	return str
}
