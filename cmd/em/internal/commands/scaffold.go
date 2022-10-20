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

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"go.pitz.tech/lib/cmd/em/internal/scaffold"
	"go.pitz.tech/lib/flagset"
	"go.pitz.tech/lib/vfs"
	"go.pitz.tech/lib/zaputil"
)

const scaffoldHelpTemplate = `
  Features:
    {{- range $feature, $files := .features }}
    - {{ $feature }}
    {{- end }}

  Aliases:
    {{- range $alias, $targets := .aliases }}
    - {{ $alias }}: {{ join $targets ", " }}
    {{- end }}

`

type ScaffoldConfig struct {
	Mkdir    bool             `json:"mkdir"    usage:"specify if we should make the target project directory"`
	License  string           `json:"license"  usage:"specify which license should be applied to the project" default:"agpl3"`
	Features *cli.StringSlice `json:"features" usage:"specify the features to generate"`
}

var (
	scaffoldConfig = &ScaffoldConfig{}

	Scaffold = &cli.Command{
		Name:  "scaffold",
		Usage: "Scaffold out a new project or add onto an existing one.",
		UsageText: flagset.ExampleString(
			"em scaffold [options] <name>",
			"em scaffold features    # will output a list of features and aliases",
			"em scaffold --mkdir --license mpl --features init <name>",
			"em scaffold --mkdir --license mpl --features init --features bin <name>",
		),
		Flags: flagset.ExtractPrefix("em", scaffoldConfig),
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return fmt.Errorf("name not specified")
			}

			name := ctx.Args().Get(0)
			if name == "features" {
				return template.Must(
					template.New("scaffold-help").
						Funcs(map[string]interface{}{
							"join": func(elems []string, sep string) string {
								return strings.Join(elems, sep)
							},
						}).
						Parse(scaffoldHelpTemplate),
				).Execute(ctx.App.Writer, map[string]interface{}{
					"features": scaffold.FilesByFeature,
					"aliases":  scaffold.FeatureAliases,
				})
			}

			if scaffoldConfig.Mkdir {
				zaputil.Extract(ctx.Context).Info("making directory")
				if err := os.MkdirAll(name, 0755); err != nil {
					return errors.Wrap(err, "failed to make project directory")
				}

				if err := os.Chdir(name); err != nil {
					return errors.Wrap(err, "failed to change into directory")
				}
			}

			zaputil.Extract(ctx.Context).Info("rendering files")
			files := scaffold.Template(
				scaffold.Data{
					Name:     name,
					License:  scaffoldConfig.License,
					Features: scaffoldConfig.Features.Value(),
				},
			).Render(ctx.Context)

			zaputil.Extract(ctx.Context).Info("writing files")
			afs := vfs.Extract(ctx.Context)
			for _, file := range files {
				dir := filepath.Dir(file.Name)
				_ = afs.MkdirAll(dir, 0755)

				if exists, _ := afero.Exists(afs, file.Name); exists {
					// don't overwrite existing files
					continue
				}

				zaputil.Extract(ctx.Context).Info("writing file", zap.String("file", file.Name))
				err := afero.WriteFile(afs, file.Name, file.Contents, 0644)
				if err != nil {
					return err
				}
			}

			if scaffoldConfig.Mkdir {
				if exists, _ := afero.Exists(afs, "go.mod"); exists {
					_, err := exec.Command("go", "mod", "tidy").CombinedOutput()
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
		HideHelpCommand: true,
	}
)
