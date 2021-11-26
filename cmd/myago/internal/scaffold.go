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

package internal

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/mjpitz/myago/flagset"
	"github.com/mjpitz/myago/zaputil"
)

//go:embed templates/*
var templates embed.FS

func renderTemplateFile(templateFile string) func(t *Template, out io.Writer) error {
	return func(t *Template, out io.Writer) error {
		body, err := templates.ReadFile(filepath.Join("templates", templateFile))
		if err != nil {
			return err
		}

		tem, err := template.New(templateFile).Parse(string(body))
		if err != nil {
			return err
		}

		return tem.Execute(out, t)
	}
}

type ScaffoldConfig struct {
	License string `json:"license" usage:"what license the project should use" default:"agpl3"`
}

type Template struct {
	Config *ScaffoldConfig
	Name   string
}

type Rendering struct {
	OutputFile string
	Render     func(t *Template, out io.Writer) error
}

const licenseTemplate = `https://raw.githubusercontent.com/licenses/license-templates/master/templates/%s.txt`

var (
	renderings = []*Rendering{
		{
			OutputFile: "cmd/{{ .Name }}/docker-compose.yaml",
			Render:     renderTemplateFile("docker-compose.yaml.tmpl"),
		},
		{
			OutputFile: "cmd/{{ .Name }}/Dockerfile",
			Render:     renderTemplateFile("Dockerfile.tmpl"),
		},
		{
			OutputFile: "cmd/{{ .Name }}/main.go",
			Render:     renderTemplateFile("main.go.tmpl"),
		},
		{
			OutputFile: "internal/commands/version.go",
			Render:     renderTemplateFile("version.go.tmpl"),
		},
		{
			OutputFile: "legal/header.txt",
			Render:     renderTemplateFile("header.txt.tmpl"),
		},
		{
			OutputFile: "scripts/dist-go.sh",
			Render:     renderTemplateFile("dist-go.sh.tmpl"),
		},
		{
			OutputFile: ".gitignore",
			Render: renderTemplateFile("gitignore.tmpl"),
		},
		{
			OutputFile: ".goreleaser.yaml",
			Render:     renderTemplateFile("goreleaser.yaml.tmpl"),
		},
		{
			OutputFile: "go.mod",
			Render:     renderTemplateFile("go.mod.tmpl"),
		},
		{
			OutputFile: "LICENSE",
			Render: func(t *Template, out io.Writer) error {
				resp, err := http.Get(fmt.Sprintf(licenseTemplate, t.Config.License))
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				_, err = io.Copy(out, resp.Body)
				return err
			},
		},
		{
			OutputFile: "Makefile",
			Render: renderTemplateFile("Makefile.tmpl"),
		},
		{
			OutputFile: "package.json",
			Render: renderTemplateFile("package.json.tmpl"),
		},
	}

	scaffoldConfig = &ScaffoldConfig{}

	ScaffoldCommand = &cli.Command{
		Name:      "scaffold",
		Usage:     "Build out a default go repository",
		UsageText: "myago scaffold [options] <name>",
		Flags:     flagset.ExtractPrefix("myago", scaffoldConfig),
		Action: func(ctx *cli.Context) error {
			name := ctx.Args().Get(0)

			if err := os.MkdirAll(name, 0755); err != nil {
				return errors.Wrap(err, "failed to make project directory")
			}

			t := &Template{
				Config: scaffoldConfig,
				Name:   name,
			}

			idx := make(map[string]string)

			for _, rendering := range renderings {
				tem, err := template.New(rendering.OutputFile).Parse(rendering.OutputFile)
				if err != nil {
					return err
				}

				rendered := bytes.NewBuffer(nil)
				err = tem.Execute(rendered, t)
				if err != nil {
					return err
				}

				outputFile := filepath.Join(name, rendered.String())
				outputDir := filepath.Dir(outputFile)

				err = os.MkdirAll(outputDir, 0755)
				if err != nil {
					return err
				}

				idx[rendering.OutputFile] = outputFile
			}

			render := func(rendering *Rendering, outputFile string) error {
				f, err := os.Create(outputFile)
				if err != nil {
					return err
				}
				defer f.Close()

				return rendering.Render(t, f)
			}

			log := zaputil.Extract(ctx.Context)
			for _, rendering := range renderings {
				outputFile := idx[rendering.OutputFile]
				log.Info("rendering file " + outputFile)

				err := render(rendering, outputFile)
				if err != nil {
					return err
				}
			}

			return nil
		},
		HideHelpCommand: true,
	}
)
