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

package scaffold

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"go.uber.org/zap"

	"github.com/mjpitz/myago/cmd/em/internal/scaffold/gitignore"
	"github.com/mjpitz/myago/cmd/em/internal/scaffold/licenses"
	"github.com/mjpitz/myago/zaputil"
)

var (
	//go:embed template
	//go:embed template/.goreleaser.yaml.tmpl
	//go:embed template/.gitignore.tmpl
	templates embed.FS

	// FeatureAliases provides common name associates to their appropriate feature name.
	FeatureAliases = map[string][]string{
		"init": {"make", "git", "legal", "authors"},
		"bin":  {"gobinary", "go", "version"},
	}

	// FilesByFeature contains a map of features to the associated files they render.
	FilesByFeature = map[string][]string{
		"Authors": {
			"AUTHORS.tmpl",
		},
		"Docker": {
			"cmd/{{ .Name }}/docker-compose.yaml.tmpl",
			"cmd/{{ .Name }}/Dockerfile.tmpl",
		},
		"GoBinary": {
			"internal/commands/version.go.tmpl",
			"cmd/{{ .Name }}/main.go.tmpl",
		},
		"GoReleaser": {
			"scripts/dist-go.sh.tmpl",
			".goreleaser.yaml.tmpl",
		},
		"Go": {
			"go.mod.tmpl",
		},
		"Version": {
			"package.json.tmpl",
		},
		"Legal": {
			"legal/header.txt.tmpl",
			"LICENSE.tmpl",
		},
		"Git": {
			".gitignore.tmpl",
			".git/HEAD",
			".git/config",
		},
		"Make": {
			"Makefile.tmpl",
		},
	}

	// functions contains a map of functions for use by templates
	functions = map[string]interface{}{
		"includes": func(features []string, feature string) bool {
			for _, f := range features {
				if strings.EqualFold(f, feature) {
					return true
				}
			}
			return false
		},
		"gitignore": func(id string) string {
			ignore, _ := gitignore.Get(id)
			return ignore
		},
		"license": func(spdx string) string {
			license, ok := licenses.Get(spdx)
			if ok {
				return license
			}

			panic("todo: slow, online path")
		},
		"license_header": func(spdx string) string {
			header, ok := licenses.GetHeader(spdx)
			if ok {
				return header
			}

			license, ok := licenses.Get(spdx)
			if ok {
				return license
			}

			panic("todo: slow, online path")
		},
	}

	gitHeadContents = []byte("ref: refs/heads/main")
)

var gitConfigContents = []byte(strings.TrimSpace(`
[core]
    repositoryformatversion = 0
    filemode = true
    bare = false
    logallrefupdates = true
`))

func init() {
	// automatically setup common aliases
	for featureName := range FilesByFeature {
		FeatureAliases[strings.ToLower(featureName)] = []string{featureName}
	}
}

// Data defines the information needed to render the template.
type Data struct {
	Name     string   `json:"name"`
	License  string   `json:"license"`
	Features []string `json:"features"`
}

func resolveFeatures(features []string) []string {
	resolvedFeatures := append([]string{}, features...)
	resolvedAFeature := true

	for resolvedAFeature {
		resolvedAFeature = false
		next := make([]string, 0)

		for _, feature := range resolvedFeatures {
			resolved, aliased := FeatureAliases[feature]
			_, validFeature := FilesByFeature[feature]

			switch {
			case aliased:
				next = append(next, resolved...)
				resolvedAFeature = true
			case !validFeature:
				// skip invalid, non-aliased features
			default:
				next = append(next, feature)
			}
		}

		resolvedFeatures = next
	}

	return resolvedFeatures
}

// Template returns a new Renderer instance that can be used to render template data based on specified features.
func Template(data Data) *Renderer {
	data.Features = resolveFeatures(data.Features)

	renderer := template.New(data.Name).
		Funcs(sprig.TxtFuncMap()).
		Funcs(functions)

	templateFiles, _ := fs.Sub(templates, "template")

	return &Renderer{
		templates: templateFiles,
		data:      data,
		renderer:  renderer,
	}
}

// File defines a generic file to render.
type File struct {
	Name     string
	Contents []byte
}

// Renderer encapsulates the business logic for rendering template data.
type Renderer struct {
	templates fs.FS
	data      Data
	renderer  *template.Template
}

func (s *Renderer) Render(ctx context.Context) []File {
	logger := zaputil.Extract(ctx)

	filesToRender := make([]string, 0)
	for _, feature := range s.data.Features {
		filesToRender = append(filesToRender, FilesByFeature[feature]...)
	}

	results := make([]File, 0, len(filesToRender))
	renderedFiles := make(map[string]int)

	for _, templateName := range filesToRender {
		logger.Info("processing", zap.String("file", templateName))

		if _, ok := renderedFiles[templateName]; ok {
			logger.Debug("skipping duplicate file", zap.String("file", templateName))
			continue
		}

		if strings.HasPrefix(templateName, ".git/") {
			renderedFiles[templateName] = len(results)
			switch templateName {
			case ".git/HEAD":
				renderedFiles[templateName] = len(results)
				results = append(results, File{
					Name:     templateName,
					Contents: gitHeadContents,
				})
				continue
			case ".git/config":
				renderedFiles[templateName] = len(results)
				results = append(results, File{
					Name:     templateName,
					Contents: gitConfigContents,
				})
				continue
			}
		}

		file, err := s.templates.Open(templateName)
		if err != nil {
			logger.Warn("requested file not found, skipping", zap.String("file", templateName), zap.Error(err))
			continue
		}

		templateContents, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Warn("failed to read template contents", zap.String("file", templateName), zap.Error(err))
			continue
		}

		name := template.Must(s.renderer.Clone())
		contents := template.Must(s.renderer.Clone())

		renderedName := &bytes.Buffer{}
		renderedContents := &bytes.Buffer{}

		err = template.Must(name.Parse(templateName)).Execute(renderedName, s.data)
		if err != nil {
			logger.Warn("requested file not found, skipping", zap.String("file", templateName))
			continue
		}

		err = template.Must(contents.Parse(string(templateContents))).Execute(renderedContents, s.data)
		if err != nil {
			logger.Error("fa")
		}

		renderedFiles[templateName] = len(results)
		results = append(results, File{
			Name:     strings.TrimSuffix(renderedName.String(), ".tmpl"),
			Contents: renderedContents.Bytes(),
		})
	}

	return results
}
